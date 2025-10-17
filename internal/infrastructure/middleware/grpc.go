package middleware

import (
	"context"
	"runtime/debug"
	"time"

	"github.com/encomos/api-encomos/customer-service/internal/infrastructure/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
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
