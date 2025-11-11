package grpc

import (
	"context"
	"fmt"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	"github.com/encomos/api-encomos/customer-service/internal/config"
	"github.com/encomos/api-encomos/customer-service/internal/domain/service"
	"github.com/encomos/api-encomos/customer-service/internal/infrastructure/logger"
	"github.com/encomos/api-encomos/customer-service/internal/infrastructure/middleware"
	customerpb "github.com/encomos/api-encomos/customer-service/proto/customer"
)

// Server represents the gRPC server
type Server struct {
	server   *grpc.Server
	listener net.Listener
	config   *config.GRPCConfig
	logger   *logger.Logger
}

// NewServer creates a new gRPC server
func NewServer(cfg *config.GRPCConfig) (*Server, error) {
	// Create logger
	logger := logger.NewWithService("customer-service")

	// Create listener
	address := fmt.Sprintf(":%d", cfg.Port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("failed to listen on %s: %w", address, err)
	}

	// Create gRPC server with middleware
	serverOptions := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			middleware.TenantInterceptor(logger),
			middleware.LoggingInterceptor(logger),
			middleware.RecoveryInterceptor(logger),
			// TODO: Add authentication interceptor when auth service is ready
		),
		grpc.ChainStreamInterceptor(
			middleware.StreamTenantInterceptor(logger),
			middleware.StreamLoggingInterceptor(logger),
			middleware.StreamRecoveryInterceptor(logger),
		),
	}

	// Add TLS if configured
	if !cfg.Insecure {
		// TODO: Add TLS configuration when needed
		logger.WithFields(map[string]interface{}{"tls": "not_implemented"}).Warn("TLS is configured but not implemented yet")
	}

	server := grpc.NewServer(serverOptions...)

	return &Server{
		server:   server,
		listener: listener,
		config:   cfg,
		logger:   logger,
	}, nil
}

// RegisterServices registers all gRPC services
func (s *Server) RegisterServices(
	customerService *service.CustomerService,
	vehicleService *service.VehicleService,
) {
	// Create handlers
	customerHandler := NewCustomerHandler(customerService, vehicleService)

	// Register services
	customerpb.RegisterCustomerServiceServer(s.server, customerHandler)

	// Register health service
	healthServer := health.NewServer()
	healthServer.SetServingStatus("customer-service", grpc_health_v1.HealthCheckResponse_SERVING)
	grpc_health_v1.RegisterHealthServer(s.server, healthServer)

	// Enable reflection for development
	if s.config.Insecure {
		reflection.Register(s.server)
		s.logger.WithFields(map[string]interface{}{"reflection": "enabled"}).Info("gRPC reflection enabled (development mode)")
	}

	s.logger.WithFields(map[string]interface{}{"status": "registered"}).Info("All gRPC services registered successfully")
}

// Start starts the gRPC server
func (s *Server) Start() error {
	s.logger.WithFields(map[string]interface{}{"address": s.listener.Addr().String()}).Info("Starting gRPC server")

	go func() {
		if err := s.server.Serve(s.listener); err != nil {
			s.logger.WithError(err).Error("gRPC server failed")
		}
	}()

	s.logger.WithFields(map[string]interface{}{"port": s.config.Port}).Info("gRPC server started successfully")
	return nil
}

// Stop stops the gRPC server gracefully
func (s *Server) Stop(ctx context.Context) error {
	s.logger.WithFields(map[string]interface{}{"action": "stopping"}).Info("Stopping gRPC server...")

	// Channel to signal when graceful stop is complete
	stopped := make(chan struct{})

	go func() {
		s.server.GracefulStop()
		close(stopped)
	}()

	// Wait for graceful stop or context timeout
	select {
	case <-stopped:
		s.logger.WithFields(map[string]interface{}{"status": "graceful"}).Info("gRPC server stopped gracefully")
		return nil
	case <-ctx.Done():
		s.logger.WithFields(map[string]interface{}{"timeout": true}).Warn("gRPC server stop timeout, forcing shutdown")
		s.server.Stop()
		return ctx.Err()
	}
}

// GetPort returns the port the server is listening on
func (s *Server) GetPort() int {
	if s.listener != nil {
		if addr, ok := s.listener.Addr().(*net.TCPAddr); ok {
			return addr.Port
		}
	}
	return s.config.Port
}

// Healthcheck checks if the server is healthy
func (s *Server) Healthcheck() error {
	if s.server == nil {
		return fmt.Errorf("gRPC server is not initialized")
	}

	// Create a simple connection to test the server
	conn, err := grpc.Dial(
		s.listener.Addr().String(),
		grpc.WithInsecure(),
		grpc.WithTimeout(5*time.Second),
	)
	if err != nil {
		return fmt.Errorf("failed to connect to gRPC server: %w", err)
	}
	defer conn.Close()

	// Test health check
	client := grpc_health_v1.NewHealthClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	resp, err := client.Check(ctx, &grpc_health_v1.HealthCheckRequest{
		Service: "customer-service",
	})
	if err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}

	if resp.Status != grpc_health_v1.HealthCheckResponse_SERVING {
		return fmt.Errorf("service is not serving, status: %v", resp.Status)
	}

	return nil
}

// GetGRPCServer returns the underlying gRPC server (for testing)
func (s *Server) GetGRPCServer() *grpc.Server {
	return s.server
}

// GetListener returns the underlying listener (for testing)
func (s *Server) GetListener() net.Listener {
	return s.listener
}

// AddHealthCheck adds a health check for a specific service
func (s *Server) AddHealthCheck(serviceName string, check func() error) {
	// TODO: Implement custom health checks if needed
	s.logger.WithFields(map[string]interface{}{"service": serviceName}).Info("Health check added")
}

// SetServingStatus sets the serving status for health checks
func (s *Server) SetServingStatus(serviceName string, status grpc_health_v1.HealthCheckResponse_ServingStatus) {
	// TODO: Get health server and update status if needed
	s.logger.WithFields(map[string]interface{}{
		"service": serviceName,
		"status":  status.String(),
	}).Info("Service status updated")
}
