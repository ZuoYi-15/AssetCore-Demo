# Asset-Core

Asset-Core is a trusted device asset management system providing comprehensive asset lifecycle management, identity generation and verification, and cross-verification capabilities. It enables organizations to maintain a centralized registry of IoT and edge devices with rich metadata, automatic identity assignment, and verification workflows.

## Features

- **Asset Management**: Create, retrieve, update, and delete device assets with rich metadata
- **Identity Generation**: Automatically generate unique DIDs (Decentralized Identifiers) for assets
- **Verification Workflow**: Verify asset authenticity through cross-verification mechanisms
- **Audit Trail**: Complete change history for all assets
- **Bulk Operations**: Support for CSV import and asset data export
- **Status Tracking**: Comprehensive asset status lifecycle management

## Tech Stack

| Component | Technology |
|-----------|------------|
| Backend API | Go 1.22 + Gin |
| Database | MySQL 8.0 + GORM |
| Cache | Redis |
| Message Queue | Kafka (optional, disabled by default) |
| Logging | Zap |
| Config Management | Viper |
| Frontend | Vue 3 + TypeScript + Vite + Element Plus |

## Project Structure

```
asset-core/
├── cmd/
│   ├── server/           # Main API server
│   └── migrate/          # Database migration tool
├── internal/
│   ├── api/              # HTTP handlers & middleware
│   ├── infrastructure/   # Database, cache, message queue
│   ├── module/           # Business logic (asset, identity, verification, audit)
│   └── pkg/              # Shared utilities (logger, errors, pagination, response)
├── web/                  # Vue 3 frontend application
├── deployments/          # Docker Compose configuration
├── docs/                 # API documentation & Postman collection
└── config.example.yaml   # Configuration template
```

## Prerequisites

- Go 1.22+
- Node.js 18+
- Docker & Docker Compose
- MySQL 8.0 (via Docker)
- Redis (via Docker)

## Quick Start

### Setup Backend

```powershell
# Copy configuration template
Copy-Item config.example.yaml config.yaml

# Start databases
docker compose -f deployments/docker-compose.yml up -d mysql redis

# Run database migrations
go run ./cmd/migrate

# Start API server
go run ./cmd/server
```

The API will be available at `http://127.0.0.1:8080`

### Setup Frontend

```powershell
Set-Location web
npm install
npm run dev
```

The frontend will be available at `http://127.0.0.1:5173`

### Health Check

```powershell
# Check API health
Invoke-RestMethod http://127.0.0.1:8080/health
```

## Troubleshooting

### MySQL Access Denied

If MySQL returns `Error 1045 (28000): Access denied`, the existing Docker volume was initialized with an older password. MySQL only applies `MYSQL_ROOT_PASSWORD`, `MYSQL_USER`, and `MYSQL_PASSWORD` when the data directory is first created.

**Solution**: Reset the volume:

```powershell
docker compose -f deployments/docker-compose.yml down -v
docker compose -f deployments/docker-compose.yml up -d mysql redis
go run ./cmd/migrate
```

## API Documentation

### Asset Management

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/assets` | Create new asset |
| GET | `/api/v1/assets` | List assets (paginated) |
| GET | `/api/v1/assets/{id}` | Get asset details |
| PUT | `/api/v1/assets/{id}` | Update asset |
| DELETE | `/api/v1/assets/{id}` | Delete asset |
| GET | `/api/v1/assets/{id}/changes` | Get asset change history |
| PUT | `/api/v1/assets/{id}/status` | Update asset status |

### Identity Management

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/identities/generate` | Generate identity for asset |
| GET | `/api/v1/identities/{identity_id}` | Get identity details |
| PUT | `/api/v1/identities/{identity_id}/bind` | Bind identity to asset |
| POST | `/api/v1/identities/{identity_id}/unbind` | Unbind identity from asset |
| GET | `/api/v1/identities/{identity_id}/features` | Get identity features |

### Verification

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/assets/{id}/verify` | Verify asset |
| GET | `/api/v1/assets/{id}/verification-result` | Get verification result |
| POST | `/api/v1/verifications` | Create verification record |
| GET | `/api/v1/verifications/{id}` | Get verification details |

### Data Operations

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/data/import` | Initiate CSV data import task |
| GET | `/api/v1/data/import-tasks` | List import tasks (paginated) |
| GET | `/api/v1/data/import-tasks/{id}` | Get import task details |
| GET | `/api/v1/data/import-tasks/{id}/errors` | Get import task errors |
| GET | `/api/v1/data/export/assets` | Export all assets as CSV |

### Example: Create Asset

```bash
curl --location 'http://127.0.0.1:8080/api/v1/assets' \
  --header 'Content-Type: application/json' \
  --data '{
    "asset_name": "edge-gateway-01",
    "asset_type": "gateway",
    "vendor": "example-vendor",
    "model": "gw-1000",
    "serial_number": "SN-001",
    "mac_address": "00:11:22:33:44:55",
    "ip_address": "192.168.1.10",
    "hostname": "edge-gateway-01",
    "owner_department": "security",
    "owner_user": "admin",
    "location": "shanghai",
    "source": "manual"
  }'
```

### Example: Generate Identity

```bash
curl --location 'http://127.0.0.1:8080/api/v1/identities/generate' \
  --header 'Content-Type: application/json' \
  --data '{
    "tenant_id": "default",
    "serial_number": "SN-001",
    "vendor": "example-vendor",
    "model": "gw-1000",
    "mac_address": "00:11:22:33:44:55",
    "ip_address": "192.168.1.10",
    "source": "manual"
  }'
```

See [docs/asset-core-curl.md](docs/asset-core-curl.md) for complete API examples.

## Configuration

Copy `config.example.yaml` to `config.yaml` and customize:

```yaml
server:
  port: 8080
  
database:
  dsn: "user:password@tcp(localhost:3306)/asset_core?charset=utf8mb4&parseTime=True"
  
redis:
  addr: "localhost:6379"
  
kafka:
  enabled: false
  brokers:
    - "localhost:9092"
```

## Development

### Run Tests

```bash
go test ./...
```

### Build Frontend

```powershell
Set-Location web
npm run build
```

### Format Code

```bash
go fmt ./...
```

## Database Schema

The system maintains the following core entities:

- **Assets**: Device inventory with metadata
- **Identities**: DIDs assigned to assets for cross-verification
- **Verifications**: Verification records and audit trails
- **Audit Logs**: Complete change history for compliance

## License

Proprietary
