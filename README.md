# Customer Service (CRM)

Microservicio responsable de la gestión de información de clientes finales para cada tenant en la plataforma Encomos.

## Responsabilidades

- **CRUD de perfiles de clientes** (contacto, historial de compras/servicios)
- **Notas, preferencias, etiquetas/segmentación**
- **AutoParts**: Vehículos asociados a clientes
- **Barbería & Estética**: Historial de citas, preferencias de servicios/profesionales
- **Programas de fidelización** (opcional)

## Puerto de Servicio

- **gRPC**: 50055

## Modelos de Datos

### Customer
```go
type Customer struct {
    ID           int64     `json:"id"`
    TenantID     int64     `json:"tenant_id"`
    FirstName    string    `json:"first_name"`
    LastName     string    `json:"last_name"`
    Email        *string   `json:"email"`
    Phone        *string   `json:"phone"`
    CustomerType string    `json:"customer_type"` // individual, business
    CompanyName  *string   `json:"company_name"`
    TaxID        *string   `json:"tax_id"`
    Address      *string   `json:"address"`
    Birthday     *time.Time `json:"birthday"`
    Notes        *string   `json:"notes"`
    Preferences  json.RawMessage `json:"preferences"`
    IsActive     bool      `json:"is_active"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
}
```

### Vehicle (AutoParts)
```go
type Vehicle struct {
    ID         int64     `json:"id"`
    CustomerID int64     `json:"customer_id"`
    Make       string    `json:"make"`
    Model      string    `json:"model"`
    Year       int       `json:"year"`
    VIN        *string   `json:"vin"`
    LicensePlate *string `json:"license_plate"`
    Notes      *string   `json:"notes"`
    IsActive   bool      `json:"is_active"`
    CreatedAt  time.Time `json:"created_at"`
}
```

## APIs gRPC

```protobuf
service CustomerService {
  rpc ListCustomers(ListCustomersRequest) returns (ListCustomersResponse);
  rpc GetCustomer(GetCustomerRequest) returns (GetCustomerResponse);
  rpc CreateCustomer(CreateCustomerRequest) returns (CreateCustomerResponse);
  rpc UpdateCustomer(UpdateCustomerRequest) returns (UpdateCustomerResponse);
  rpc DeleteCustomer(DeleteCustomerRequest) returns (DeleteCustomerResponse);
  
  // Vehicles (AutoParts)
  rpc ListVehicles(ListVehiclesRequest) returns (ListVehiclesResponse);
  rpc CreateVehicle(CreateVehicleRequest) returns (CreateVehicleResponse);
  rpc UpdateVehicle(UpdateVehicleRequest) returns (UpdateVehicleResponse);
  
  // Search
  rpc SearchCustomers(SearchCustomersRequest) returns (SearchCustomersResponse);
}
```

---

**Customer Service** - CRM inteligente para gestión de clientes por tipo de negocio.
