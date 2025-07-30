package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/yourorg/api-encomos/customer-service/internal/config"
	"github.com/yourorg/api-encomos/customer-service/internal/domain/service"
	"github.com/yourorg/api-encomos/customer-service/internal/infrastructure/grpc"
	"github.com/yourorg/api-encomos/customer-service/internal/infrastructure/persistence/postgres"
)

func main() {
	// Configurar el logger
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Iniciando Customer Service...")

	// Cargar configuración
	cfg, err := loadConfig()
	if err != nil {
		log.Fatalf("Error al cargar configuración: %v", err)
	}

	log.Printf("Configuración cargada para entorno: %s", cfg.Server.Environment)

	// Conectar a PostgreSQL
	db, err := postgres.NewDB(&cfg.Database)
	if err != nil {
		log.Fatalf("Error al conectar a PostgreSQL: %v", err)
	}
	defer db.Close()
	log.Println("✓ Conexión a PostgreSQL establecida")

	// Crear repositorios
	customerRepo := postgres.NewCustomerRepository(db)
	vehicleRepo := postgres.NewVehicleRepository(db)
	customerNoteRepo := postgres.NewCustomerNoteRepository(db)

	log.Println("✓ Repositorios inicializados")

	// Crear servicios de dominio
	customerService := service.NewCustomerService(customerRepo, vehicleRepo, customerNoteRepo)
	vehicleService := service.NewVehicleService(vehicleRepo, customerRepo)

	log.Println("✓ Servicios de dominio inicializados")

	// Crear servidor gRPC
	grpcServer, err := grpc.NewServer(&cfg.GRPC)
	if err != nil {
		log.Fatalf("Error al crear servidor gRPC: %v", err)
	}

	// Registrar servicios gRPC
	grpcServer.RegisterServices(customerService, vehicleService)

	log.Println("✓ Servicios gRPC registrados")

	// Iniciar servidor gRPC
	if err := grpcServer.Start(); err != nil {
		log.Fatalf("Error al iniciar servidor gRPC: %v", err)
	}

	log.Printf("✓ Servidor gRPC iniciado en puerto %d", cfg.GRPC.Port)

	// Configurar servidor HTTP para health checks
	httpServer := setupHTTPServer(cfg.HTTP.Port, db, grpcServer)

	log.Printf("✓ Servidor HTTP iniciado en puerto %d", cfg.HTTP.Port)
	log.Println("🚀 Customer Service completamente inicializado")

	// Capturar señales para shutdown graceful
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	// Esperar señal
	sig := <-signalChan
	log.Printf("⚠️  Señal de terminación recibida: %v", sig)

	// Shutdown graceful
	shutdownGracefully(httpServer, grpcServer, cfg.Server.ShutdownTime)
}

// loadConfig carga la configuración desde el archivo y variables de entorno
func loadConfig() (*config.Config, error) {
	env := os.Getenv("ENV")
	if env == "" {
		env = "local" // Default a entorno local
	}

	log.Printf("Cargando configuración para entorno: %s", env)

	configPath := filepath.Join("config", env)

	// Verificar si el directorio existe
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Probar en el directorio padre
		configPath = filepath.Join("..", "config", env)
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			// Probar sin directorio específico (solo variables de entorno)
			log.Printf("Directorio de configuración %s no encontrado, usando solo variables de entorno", configPath)
			configPath = ""
		}
	}

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("error al cargar configuración: %w", err)
	}

	// Log de configuraciones importantes (sin secretos)
	log.Printf("Database: %s:%d/%s", cfg.Database.Host, cfg.Database.Port, cfg.Database.Name)
	log.Printf("gRPC Port: %d (Insecure: %v)", cfg.GRPC.Port, cfg.GRPC.Insecure)
	log.Printf("HTTP Port: %d", cfg.HTTP.Port)

	return cfg, nil
}

// setupHTTPServer configura un servidor HTTP para health checks y métricas
func setupHTTPServer(port int, db *postgres.DB, grpcServer *grpc.Server) *http.Server {
	mux := http.NewServeMux()

	// Ruta para health check general
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		healthStatus := make(map[string]interface{})
		overallStatus := "ok"

		// Verificar conexión a PostgreSQL
		if err := db.Healthcheck(r.Context()); err != nil {
			healthStatus["database"] = map[string]string{
				"status": "error",
				"error":  err.Error(),
			}
			overallStatus = "error"
		} else {
			healthStatus["database"] = map[string]string{"status": "ok"}
		}

		// Verificar servidor gRPC
		if err := grpcServer.Healthcheck(); err != nil {
			healthStatus["grpc"] = map[string]string{
				"status": "error",
				"error":  err.Error(),
			}
			overallStatus = "error"
		} else {
			healthStatus["grpc"] = map[string]string{"status": "ok"}
		}

		healthStatus["overall"] = overallStatus
		healthStatus["timestamp"] = time.Now().UTC().Format(time.RFC3339)
		healthStatus["service"] = "customer-service"

		// Establecer código de estado HTTP
		if overallStatus == "error" {
			w.WriteHeader(http.StatusServiceUnavailable)
		} else {
			w.WriteHeader(http.StatusOK)
		}

		// Escribir respuesta JSON
		w.Header().Set("Content-Type", "application/json")
		if overallStatus == "ok" {
			w.Write([]byte(`{"status":"ok","message":"customer service is healthy","services":{"database":"ok","grpc":"ok"}}`))
		} else {
			w.Write([]byte(fmt.Sprintf(`{"status":"error","message":"some services are unhealthy","details":%v}`, formatHealthStatus(healthStatus))))
		}
	})

	// Ruta para health check específico de base de datos
	mux.HandleFunc("/health/database", func(w http.ResponseWriter, r *http.Request) {
		if err := db.Healthcheck(r.Context()); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(fmt.Sprintf(`{"status":"error","message":"%v"}`, err)))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","message":"database is healthy"}`))
	})

	// Ruta para health check específico de gRPC
	mux.HandleFunc("/health/grpc", func(w http.ResponseWriter, r *http.Request) {
		if err := grpcServer.Healthcheck(); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(fmt.Sprintf(`{"status":"error","message":"%v"}`, err)))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","message":"grpc is healthy"}`))
	})

	// Ruta para información del servicio
	mux.HandleFunc("/info", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf(`{
			"service": "customer-service",
			"version": "1.0.0",
			"grpc_port": %d,
			"http_port": %d,
			"status": "running",
			"timestamp": "%s"
		}`, grpcServer.GetPort(), port, time.Now().UTC().Format(time.RFC3339))))
	})

	// Ruta para métricas básicas
	mux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"service": "customer-service",
			"metrics": {
				"uptime": "` + time.Since(time.Now()).String() + `",
				"status": "healthy"
			}
		}`))
	})

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
		// Configurar timeouts
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Iniciar servidor HTTP en goroutine
	go func() {
		log.Printf("Servidor HTTP iniciado en puerto %d", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error al iniciar servidor HTTP: %v", err)
		}
	}()

	return server
}

// shutdownGracefully detiene los servidores de forma controlada
func shutdownGracefully(httpServer *http.Server, grpcServer *grpc.Server, shutdownTimeout time.Duration) {
	log.Printf("Iniciando shutdown graceful (timeout: %v)...", shutdownTimeout)

	// Crear contexto con timeout para shutdown
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	// Canal para coordinar shutdowns
	shutdownComplete := make(chan bool, 2)

	// Detener servidor HTTP
	go func() {
		log.Println("🔄 Deteniendo servidor HTTP...")
		if err := httpServer.Shutdown(ctx); err != nil {
			log.Printf("❌ Error al detener servidor HTTP: %v", err)
		} else {
			log.Println("✅ Servidor HTTP detenido correctamente")
		}
		shutdownComplete <- true
	}()

	// Detener servidor gRPC
	go func() {
		log.Println("🔄 Deteniendo servidor gRPC...")
		if err := grpcServer.Stop(ctx); err != nil {
			log.Printf("❌ Error al detener servidor gRPC: %v", err)
		} else {
			log.Println("✅ Servidor gRPC detenido correctamente")
		}
		shutdownComplete <- true
	}()

	// Esperar a que ambos servidores terminen o timeout
	shutdownCount := 0
	for shutdownCount < 2 {
		select {
		case <-shutdownComplete:
			shutdownCount++
		case <-ctx.Done():
			log.Println("⚠️  Timeout de shutdown alcanzado, forzando terminación")
			return
		}
	}

	log.Println("🎉 Customer Service terminado correctamente")
}

// formatHealthStatus formatea el estado de salud para la respuesta JSON
func formatHealthStatus(status map[string]interface{}) string {
	// Implementación simple para evitar dependencias adicionales
	return `{"status":"partial"}`
}
