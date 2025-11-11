package middleware

import (
	"context"
	"fmt"
	"runtime/debug"
	"time"

	"github.com/encomos/api-encomos/customer-service/internal/infrastructure/logger"
	"github.com/encomos/api-encomos/customer-service/internal/infrastructure/persistence/postgres"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// LoggingInterceptor logs gRPC requests and responses
func LoggingInterceptor(logger *logger.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()

		logger.WithFields(map[string]interface{}{
			"method": info.FullMethod,
			"type":   "unary",
		}).Info("gRPC request started")

		resp, err := handler(ctx, req)

		duration := time.Since(start)
		logEntry := logger.WithFields(map[string]interface{}{
			"method":   info.FullMethod,
			"duration": duration.String(),
			"type":     "unary",
		})

		if err != nil {
			logEntry.WithError(err).Error("gRPC request failed")
		} else {
			logEntry.Info("gRPC request completed")
		}

		return resp, err
	}
}

// RecoveryInterceptor recovers from panics in gRPC handlers
func RecoveryInterceptor(logger *logger.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				logger.WithFields(map[string]interface{}{
					"method": info.FullMethod,
					"panic":  r,
					"stack":  string(debug.Stack()),
				}).Error("gRPC handler panicked")

				err = status.Errorf(codes.Internal, "internal server error")
			}
		}()

		return handler(ctx, req)
	}
}

// StreamLoggingInterceptor logs gRPC stream requests
func StreamLoggingInterceptor(logger *logger.Logger) grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		start := time.Now()

		logger.WithFields(map[string]interface{}{
			"method": info.FullMethod,
			"type":   "stream",
		}).Info("gRPC stream started")

		err := handler(srv, stream)

		duration := time.Since(start)
		logEntry := logger.WithFields(map[string]interface{}{
			"method":   info.FullMethod,
			"duration": duration.String(),
			"type":     "stream",
		})

		if err != nil {
			logEntry.WithError(err).Error("gRPC stream failed")
		} else {
			logEntry.Info("gRPC stream completed")
		}

		return err
	}
}

// StreamRecoveryInterceptor recovers from panics in gRPC stream handlers
func StreamRecoveryInterceptor(logger *logger.Logger) grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		defer func() {
			if r := recover(); r != nil {
				logger.WithFields(map[string]interface{}{
					"method": info.FullMethod,
					"panic":  r,
					"stack":  string(debug.Stack()),
				}).Error("gRPC stream handler panicked")

				err = status.Errorf(codes.Internal, "internal server error")
			}
		}()

		return handler(srv, stream)
	}
}

// TenantInterceptor extracts tenant_id from metadata and adds it to context
func TenantInterceptor(logger *logger.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Extract metadata from context
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			logger.WithFields(map[string]interface{}{
				"method": info.FullMethod,
			}).Error("no metadata found in request")
			return nil, status.Errorf(codes.InvalidArgument, "tenant_id is required")
		}

		// DEBUG: Log all metadata
		logger.WithFields(map[string]interface{}{
			"method":   info.FullMethod,
			"metadata": md,
		}).Info("DEBUG: All incoming metadata")

		// Get tenant_id from metadata (sent by API Gateway)
		tenantIDValues := md.Get("x-tenant-id")
		if len(tenantIDValues) == 0 {
			logger.WithFields(map[string]interface{}{
				"method":   info.FullMethod,
				"metadata": md,
			}).Error("x-tenant-id not found in metadata")
			return nil, status.Errorf(codes.InvalidArgument, "tenant_id is required")
		}

		tenantID := tenantIDValues[0]
		if tenantID == "" {
			logger.WithFields(map[string]interface{}{
				"method": info.FullMethod,
			}).Error("x-tenant-id is empty")
			return nil, status.Errorf(codes.InvalidArgument, "tenant_id is required")
		}

		// Add tenant_id to context using the correct postgres helper function
		ctx = postgres.WithTenantID(ctx, tenantID)

		logger.WithFields(map[string]interface{}{
			"method":    info.FullMethod,
			"tenant_id": tenantID,
		}).Info("DEBUG TENANT INTERCEPTOR: tenant_id extracted and added to context")

		// Verify it was added correctly
		verifyValue, verifyOK := postgres.GetTenantID(ctx)
		logger.WithFields(map[string]interface{}{
			"method":       info.FullMethod,
			"verify_value": verifyValue,
			"verify_ok":    verifyOK,
			"verify_type":  fmt.Sprintf("%T", verifyValue),
		}).Info("DEBUG TENANT INTERCEPTOR: Verifying context value")

		return handler(ctx, req)
	}
}

// StreamTenantInterceptor extracts tenant_id from metadata for stream requests
func StreamTenantInterceptor(logger *logger.Logger) grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := stream.Context()

		// Extract metadata from context
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			logger.WithFields(map[string]interface{}{
				"method": info.FullMethod,
			}).Error("no metadata found in stream request")
			return status.Errorf(codes.InvalidArgument, "tenant_id is required")
		}


		// Get tenant_id from metadata
		tenantIDValues := md.Get("x-tenant-id")
		if len(tenantIDValues) == 0 {
			logger.WithFields(map[string]interface{}{
				"method": info.FullMethod,
			}).Error("x-tenant-id not found in metadata")
			return status.Errorf(codes.InvalidArgument, "tenant_id is required")
		}

		tenantID := tenantIDValues[0]
		if tenantID == "" {
			logger.WithFields(map[string]interface{}{
				"method": info.FullMethod,
			}).Error("x-tenant-id is empty")
			return status.Errorf(codes.InvalidArgument, "tenant_id is required")
		}

		// Add tenant_id to context using the correct postgres helper function
		ctx = postgres.WithTenantID(ctx, tenantID)

		logger.WithFields(map[string]interface{}{
			"method":    info.FullMethod,
			"tenant_id": tenantID,
		}).Debug("tenant_id extracted from metadata for stream")

		// Wrap the stream with the new context
		wrapped := &wrappedServerStream{
			ServerStream: stream,
			ctx:          ctx,
		}

		return handler(srv, wrapped)
	}
}

// wrappedServerStream wraps a grpc.ServerStream with a custom context
type wrappedServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

// Context returns the custom context
func (w *wrappedServerStream) Context() context.Context {
	return w.ctx
}
