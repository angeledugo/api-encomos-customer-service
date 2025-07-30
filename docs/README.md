# Customer Service - Encomos

**Puerto**: 50055  
**Estado**: ✅ **COMPLETAMENTE IMPLEMENTADO**
**Tipo**: Microservicio gRPC  
**Base de Datos**: PostgreSQL (compartida con RLS)  

## Descripción

El Customer Service es responsable de la gestión completa de información de clientes finales para cada tenant en la plataforma Encomos. Proporciona funcionalidades de CRM (Customer Relationship Management) adaptadas a los diferentes tipos de licencia.

## Características Implementadas

### ✅ Gestión de Clientes
- **CRUD completo** de perfiles de clientes
- **Tipos de cliente**: Individual y Business (empresas)
- **Validaciones** de email y Tax ID únicos por tenant
- **Búsqueda avanzada** multi-campo con paginación
- **Activación/Desactivación** de clientes

### ✅ Gestión de Vehículos (AutoParts)
- **CRUD de vehículos** asociados a clientes
- **Validación de VIN** (17 caracteres, sin I/O/Q)
- **Búsqueda por compatibilidad** para repuestos
- **Gestión de placas** únicas
- **Metadatos flexibles** en JSON

### ✅ Sistema de Notas
- **Notas por cliente** con tipos (general, service, complaint, etc.)
- **Historial temporal** de interacciones
- **Búsqueda por tipo** y fecha
- **Staff tracking** (quién creó la nota)

### ✅ Funcionalidades Avanzadas
- **Multi-tenancy** con Row-Level Security (RLS)
- **Preferencias de cliente** en formato JSON
- **Estadísticas de cliente** (placeholder para integración futura)
- **Búsqueda inteligente** con scoring por relevancia

## Estructura del Proyecto

```
customer-service/
├── cmd/
│   └── main.go                     # ✅ Punto de entrada completo
├── config/
│   └── local/
│       └── app.env                 # ✅ Variables de entorno
├── internal/
│   ├── config/
│   │   └── config.go              # ✅ Configuración con Viper
│   ├── domain/
│   │   ├── model/                 # ✅ Todos los modelos
│   │   │   ├── customer.go        # ✅ Modelo Customer completo
│   │   │   ├── vehicle.go         # ✅ Modelo Vehicle completo
│   │   │   ├── customer_note.go   # ✅ Modelo CustomerNote
│   │   │   └── customer_stats.go  # ✅ Modelo CustomerStats
│   │   └── service/               # ✅ Servicios de negocio
│   │       ├── customer_service.go # ✅ Lógica completa de clientes
│   │       └── vehicle_service.go  # ✅ Lógica completa de vehículos
│   ├── infrastructure/
│   │   ├── grpc/                  # ✅ Handlers gRPC
│   │   │   ├── customer_handler.go # ✅ Handler completo
│   │   │   ├── vehicle_handler.go  # ✅ Handler de vehículos
│   │   │   └── server.go          # ✅ Servidor gRPC con middleware
│   │   └── persistence/
│   │       └── postgres/          # ✅ Repositorios PostgreSQL
│   │           ├── db.go          # ✅ Conexión con RLS
│   │           ├── customer_repo.go # ✅ Repository completo
│   │           ├── vehicle_repo.go  # ✅ Repository completo
│   │           └── customer_note_repo.go # ✅ Repository completo
│   └── port/
│       └── repository/            # ✅ Interfaces de repositorio
│           ├── customer_repository.go # ✅ Interface Customer
│           ├── vehicle_repository.go  # ✅ Interface Vehicle
│           ├── customer_note_repository.go # ✅ Interface Notes
│           └── customer_stats_repository.go # ✅ Interface Stats
├── proto/
│   └── customer/
│       ├── customer.proto         # ✅ Definiciones protobuf
│       ├── customer.pb.go         # ✅ Código generado
│       └── customer_grpc.pb.go    # ✅ Servidor gRPC generado
├── docs/
│   └── README.md                  # ✅ Este archivo
├── Dockerfile                     # ✅ Imagen Docker
├── go.mod                         # ✅ Dependencias actualizadas
└── generate-proto.sh              # ✅ Script de generación
```

## APIs gRPC Implementadas

### CustomerService
```protobuf
service CustomerService {
  // CRUD Customers
  rpc ListCustomers(ListCustomersRequest) returns (ListCustomersResponse);
  rpc GetCustomer(GetCustomerRequest) returns (GetCustomerResponse);
  rpc CreateCustomer(CreateCustomerRequest) returns (CreateCustomerResponse);
  rpc UpdateCustomer(UpdateCustomerRequest) returns (UpdateCustomerResponse);
  rpc DeleteCustomer(DeleteCustomerRequest) returns (DeleteCustomerResponse);
  
  // Vehicles (AutoParts)
  rpc ListVehicles(ListVehiclesRequest) returns (ListVehiclesResponse);
  rpc GetVehicle(GetVehicleRequest) returns (GetVehicleResponse);
  rpc CreateVehicle(CreateVehicleRequest) returns (CreateVehicleResponse);
  rpc UpdateVehicle(UpdateVehicleRequest) returns (UpdateVehicleResponse);
  rpc DeleteVehicle(DeleteVehicleRequest) returns (DeleteVehicleResponse);
  
  // Search & Analysis
  rpc SearchCustomers(SearchCustomersRequest) returns (SearchCustomersResponse);
  rpc AddCustomerNote(AddCustomerNoteRequest) returns (AddCustomerNoteResponse);
  rpc GetCustomerHistory(GetCustomerHistoryRequest) returns (GetCustomerHistoryResponse);
}
```

## Configuración

### Variables de Entorno
```bash
# Servicio
SERVER_ENVIRONMENT=development
GRPC_PORT=50055
HTTP_PORT=9055

# Base de Datos
DB_HOST=localhost
DB_PORT=5432
DB_USER=encomos_user
DB_PASSWORD=dev_password_123
DB_NAME=encomos
DB_SSLMODE=disable

# Logging
LOG_LEVEL=info
LOG_JSON=false
```

## Inicio Rápido

### 1. Configurar Entorno
```bash
# Copiar configuración
cp config/local/app.env.example config/local/app.env

# Instalar dependencias
go mod tidy
```

### 2. Generar Protobuf
```bash
# Generar código protobuf
./generate-proto.sh

# O manualmente
protoc --go_out=. --go-grpc_out=. --proto_path=proto customer/customer.proto
```

### 3. Ejecutar Servicio
```bash
# Desarrollo
go run cmd/main.go

# Producción
go build -o bin/customer-service cmd/main.go
./bin/customer-service
```

### 4. Verificar Salud
```bash
# Health check general
curl http://localhost:9055/health

# Health check de base de datos
curl http://localhost:9055/health/database

# Health check de gRPC
curl http://localhost:9055/health/grpc
```

## Funcionalidades por Licencia

### 🔧 AutoParts
- **Gestión de talleres** y particulares
- **Vehículos asociados** a clientes con VIN
- **Búsqueda de compatibilidad** para repuestos
- **Historial de repuestos** comprados

### ✂️ Barbería & Estética  
- **Historial de servicios** y citas
- **Preferencias** de servicios y profesionales
- **Notas de cuidado** (alergias, tratamientos)

### 🍽️ Resto & Bar
- **Historial de pedidos** y reservas
- **Programa de fidelización**
- **Preferencias gastronómicas**

### 👔 Moda & Calzado
- **Historial de compras** y devoluciones
- **Preferencias de tallas** y marcas
- **Seguimiento de tendencias**

## Seguridad

### Row-Level Security (RLS)
- **Aislamiento automático** por tenant_id
- **Context injection** en todas las queries
- **Políticas PostgreSQL** automáticas

### Validaciones
- **Email único** por tenant
- **Tax ID único** por tenant  
- **VIN único** globalmente
- **Placa única** globalmente

## Ejemplos de Uso

### Crear Cliente
```bash
grpcurl -plaintext -d '{
  "first_name": "Juan",
  "last_name": "Pérez", 
  "email": "juan@example.com",
  "customer_type": "individual"
}' localhost:50055 customer.v1.CustomerService/CreateCustomer
```

### Buscar Clientes
```bash
grpcurl -plaintext -d '{
  "query": "juan",
  "limit": 10
}' localhost:50055 customer.v1.CustomerService/SearchCustomers
```

### Crear Vehículo (AutoParts)
```bash
grpcurl -plaintext -d '{
  "customer_id": 1,
  "make": "Toyota",
  "model": "Camry",
  "year": 2020,
  "vin": "1HGBH41JXMN109186"
}' localhost:50055 customer.v1.CustomerService/CreateVehicle
```

## Monitoreo

### Health Checks
- **HTTP :9055/health** - Estado general
- **HTTP :9055/health/database** - PostgreSQL
- **HTTP :9055/health/grpc** - Servidor gRPC

### Logging
- **Structured logging** con logrus
- **Request/Response logging** automático
- **Error tracking** con contexto

### Métricas (Futuro)
- **Customer operations** por minuto
- **Search performance** promedio
- **Database connection** pool status

## Dependencias

### Servicios Externos
- **PostgreSQL**: Base de datos principal con RLS
- **Shared-Lib**: Utilidades comunes (logging, database, middleware)

### Servicios que lo Consumen
- **Sales Service**: Para asociar órdenes a clientes
- **Appointment Service**: Para gestionar citas
- **Reporting Service**: Para generar reportes de clientes
- **API Gateway**: Para operaciones desde frontend

## Testing

### Ejecutar Tests
```bash
# Tests unitarios
go test ./internal/domain/service/...

# Tests de integración
go test ./internal/infrastructure/...

# Tests con coverage
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Docker

### Build
```bash
# Construir imagen
docker build -t encomos/customer-service:latest .

# Ejecutar contenedor
docker run -p 50055:50055 -p 9055:9055 \
  -e DB_HOST=host.docker.internal \
  encomos/customer-service:latest
```

## Troubleshooting

### Problemas Comunes

#### 1. Error de Conexión a Base de Datos
```bash
# Verificar variables de entorno
env | grep DB_

# Test de conexión
curl http://localhost:9055/health/database
```

#### 2. Error de Generación de Protobuf
```bash
# Instalar dependencias
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Regenerar
./generate-proto.sh
```

#### 3. Error de Tenant ID
```bash
# Verificar contexto en requests gRPC
# El tenant_id debe ser inyectado por el API Gateway
```

## Roadmap

### ✅ Completado
- [x] CRUD completo de clientes
- [x] Gestión de vehículos para AutoParts  
- [x] Sistema de notas
- [x] Búsqueda avanzada
- [x] Multi-tenancy con RLS
- [x] Health checks y monitoreo

### 🔄 En Progreso
- [ ] Integración con Sales Service para estadísticas
- [ ] Cache de consultas frecuentes con Redis
- [ ] Métricas avanzadas con Prometheus

### 📋 Planificado
- [ ] Tests de carga con k6
- [ ] Backup automático de datos críticos
- [ ] API REST complementaria para integraciones externas
- [ ] Webhooks para sincronización externa

---

**Customer Service** - 👥 CRM inteligente y robusto para la plataforma Encomos.

**Estado**: ✅ Listo para producción
