package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// Config representa la configuración del servicio
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	GRPC     GRPCConfig
	HTTP     HTTPConfig
	Log      LogConfig
}

// ServerConfig representa la configuración del servidor
type ServerConfig struct {
	Environment  string
	ShutdownTime time.Duration
}

// DatabaseConfig representa la configuración de la base de datos
type DatabaseConfig struct {
	Host         string
	Port         int
	User         string
	Password     string
	Name         string
	SSLMode      string
	MaxOpenConns int
	MaxIdleConns int
	MaxLifetime  time.Duration
}

// PostgresURL devuelve la URL de conexión a PostgreSQL
func (c *DatabaseConfig) PostgresURL() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode,
	)
}

// GRPCConfig representa la configuración del servidor gRPC
type GRPCConfig struct {
	Port        int
	Insecure    bool
	TLSCertFile string
	TLSKeyFile  string
}

// HTTPConfig representa la configuración del servidor HTTP
type HTTPConfig struct {
	Port               int
	CORSAllowedOrigins []string
	TLSEnabled         bool
	TLSCertFile        string
	TLSKeyFile         string
}

// LogConfig representa la configuración de logging
type LogConfig struct {
	Level string
	JSON  bool
	File  string
}

// LoadConfig carga la configuración desde archivos y variables de entorno
func LoadConfig(path string) (*Config, error) {
	v := viper.New()

	// Configuraciones por defecto
	setDefaults(v)

	// Leer desde archivo de configuración
	if path != "" {
		v.AddConfigPath(path)
		v.SetConfigName("app")
		v.SetConfigType("env")
	}

	// Leer desde variables de entorno
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error leyendo archivo de configuración: %w", err)
		}
		// Es ok si el archivo no existe, usaremos env vars y defaults
	}

	// Mapear variables de entorno explícitamente a las claves de Viper
	bindEnvVars(v)

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshal config: %w", err)
	}

	// Validar configuración
	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("configuración inválida: %w", err)
	}

	return &config, nil
}

// bindEnvVars mapea las variables de entorno a las claves de Viper
func bindEnvVars(v *viper.Viper) {
	// Server
	v.BindEnv("server.environment", "SERVER_ENVIRONMENT")
	v.BindEnv("server.shutdowntime", "SERVER_SHUTDOWNTIME")

	// Database
	v.BindEnv("database.host", "DB_HOST")
	v.BindEnv("database.port", "DB_PORT")
	v.BindEnv("database.user", "DB_USER")
	v.BindEnv("database.password", "DB_PASSWORD")
	v.BindEnv("database.name", "DB_NAME")
	v.BindEnv("database.sslmode", "DB_SSLMODE")
	v.BindEnv("database.maxopenconns", "DB_MAX_OPEN_CONNS")
	v.BindEnv("database.maxidleconns", "DB_MAX_IDLE_CONNS")
	v.BindEnv("database.maxlifetime", "DB_MAX_LIFETIME_SECONDS")

	// GRPC
	v.BindEnv("grpc.port", "GRPC_PORT")
	v.BindEnv("grpc.insecure", "GRPC_INSECURE")
	v.BindEnv("grpc.tlscertfile", "GRPC_TLS_CERT_FILE")
	v.BindEnv("grpc.tlskeyfile", "GRPC_TLS_KEY_FILE")

	// HTTP
	v.BindEnv("http.port", "HTTP_PORT")
	v.BindEnv("http.corsallowedorigins", "HTTP_CORS_ALLOWED_ORIGINS")
	v.BindEnv("http.tlsenabled", "HTTP_TLS_ENABLED")
	v.BindEnv("http.tlscertfile", "HTTP_TLS_CERT_FILE")
	v.BindEnv("http.tlskeyfile", "HTTP_TLS_KEY_FILE")

	// Log
	v.BindEnv("log.level", "LOG_LEVEL")
	v.BindEnv("log.json", "LOG_JSON")
	v.BindEnv("log.file", "LOG_FILE")
}

// setDefaults establece valores por defecto para la configuración
func setDefaults(v *viper.Viper) {
	// Server defaults
	v.SetDefault("server.environment", "development")
	v.SetDefault("server.shutdowntime", 10*time.Second)

	// Database defaults
	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", 5432)
	v.SetDefault("database.user", "encomos_user")
	v.SetDefault("database.password", "dev_password_123")
	v.SetDefault("database.name", "encomos")
	v.SetDefault("database.sslmode", "disable")
	v.SetDefault("database.maxopenconns", 25)
	v.SetDefault("database.maxidleconns", 5)
	v.SetDefault("database.maxlifetime", 30*time.Minute)

	// GRPC defaults
	v.SetDefault("grpc.port", 50055) // Puerto específico para customer-service
	v.SetDefault("grpc.insecure", true)
	v.SetDefault("grpc.tlscertfile", "")
	v.SetDefault("grpc.tlskeyfile", "")

	// HTTP defaults
	v.SetDefault("http.port", 9055) // Puerto específico para customer-service
	v.SetDefault("http.corsallowedorigins", []string{"*"})
	v.SetDefault("http.tlsenabled", false)
	v.SetDefault("http.tlscertfile", "")
	v.SetDefault("http.tlskeyfile", "")

	// Log defaults
	v.SetDefault("log.level", "info")
	v.SetDefault("log.json", false)
	v.SetDefault("log.file", "")
}

// validateConfig valida la configuración cargada
func validateConfig(config *Config) error {
	// Validar puertos
	if config.GRPC.Port <= 0 || config.GRPC.Port > 65535 {
		return fmt.Errorf("puerto gRPC inválido: %d", config.GRPC.Port)
	}

	if config.HTTP.Port <= 0 || config.HTTP.Port > 65535 {
		return fmt.Errorf("puerto HTTP inválido: %d", config.HTTP.Port)
	}

	// Validar configuración de base de datos
	if config.Database.Host == "" {
		return fmt.Errorf("host de la base de datos es requerido")
	}

	if config.Database.Name == "" {
		return fmt.Errorf("nombre de la base de datos es requerido")
	}

	return nil
}

// IsDevelopment verifica si estamos en modo desarrollo
func (c *Config) IsDevelopment() bool {
	return c.Server.Environment == "development"
}

// IsProduction verifica si estamos en modo producción
func (c *Config) IsProduction() bool {
	return c.Server.Environment == "production"
}

// GetGRPCAddress devuelve la dirección completa del servidor gRPC
func (c *Config) GetGRPCAddress() string {
	return fmt.Sprintf(":%d", c.GRPC.Port)
}

// GetHTTPAddress devuelve la dirección completa del servidor HTTP
func (c *Config) GetHTTPAddress() string {
	return fmt.Sprintf(":%d", c.HTTP.Port)
}
