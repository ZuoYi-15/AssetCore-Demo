# Asset-Core

Asset-Core is the backend service for trusted device asset management. It provides asset ledger APIs, identity ID generation, cross verification, and data management task APIs.

## Tech Stack

- Go + Gin
- MySQL + GORM
- Redis
- Kafka, optional and disabled by default
- Zap
- Viper

## Quick Start

```powershell
Copy-Item config.example.yaml config.yaml
docker compose -f deployments/docker-compose.yml up -d mysql redis
go run ./cmd/migrate
go run ./cmd/server
```

Frontend:

```powershell
Set-Location web
npm install
npm run dev
```

Open:

```text
http://127.0.0.1:5173
```

If MySQL returns `Error 1045 (28000): Access denied`, the existing Docker volume was likely initialized with an older password. MySQL only applies `MYSQL_ROOT_PASSWORD`, `MYSQL_USER`, and `MYSQL_PASSWORD` when the data directory is created for the first time. For a clean local reset:

```powershell
docker compose -f deployments/docker-compose.yml down -v
docker compose -f deployments/docker-compose.yml up -d mysql redis
go run ./cmd/migrate
```

Health check:

```powershell
Invoke-RestMethod http://127.0.0.1:8080/health
```

## Main APIs

- `POST /api/v1/assets`
- `GET /api/v1/assets`
- `GET /api/v1/assets/{id}`
- `PUT /api/v1/assets/{id}`
- `DELETE /api/v1/assets/{id}`
- `POST /api/v1/assets/{id}/verify`
- `POST /api/v1/identities/generate`
- `GET /api/v1/identities/{identity_id}`
- `POST /api/v1/verifications`
- `POST /api/v1/data/import`

## Example Asset Payload

```json
{
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
}
```
