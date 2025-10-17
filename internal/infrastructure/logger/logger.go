package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Logger wraps logrus logger
type Logger struct {
	*logrus.Logger
	service string
}

// NewWithService creates a new logger instance for a specific service
func NewWithService(serviceName string) *Logger {
	logger := logrus.New()

	// Set log level from environment or default to Info
	level := os.Getenv("LOG_LEVEL")
	if level == "" {
		level = "info"
	}

	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		logLevel = logrus.InfoLevel
	}
	logger.SetLevel(logLevel)

	// Set JSON formatter
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
	})

	return &Logger{
		Logger:  logger,
		service: serviceName,
	}
}

// WithFields adds fields to the logger
func (l *Logger) WithFields(fields map[string]interface{}) *logrus.Entry {
	fields["service"] = l.service
	return l.Logger.WithFields(fields)
}

// WithError adds error field to the logger
func (l *Logger) WithError(err error) *logrus.Entry {
	return l.Logger.WithError(err).WithField("service", l.service)
}

// WithService adds service field to the logger
func (l *Logger) WithService(service string) *logrus.Entry {
	return l.Logger.WithField("service", service)
}

// WithTenant adds tenant field to the logger
func (l *Logger) WithTenant(tenantID int64) *logrus.Entry {
	return l.Logger.WithField("tenant_id", tenantID).WithField("service", l.service)
}

// WithUser adds user field to the logger
func (l *Logger) WithUser(userID string) *logrus.Entry {
	return l.Logger.WithField("user_id", userID).WithField("service", l.service)
}

// WithRequestID adds request ID field to the logger
func (l *Logger) WithRequestID(requestID string) *logrus.Entry {
	return l.Logger.WithField("request_id", requestID).WithField("service", l.service)
}
