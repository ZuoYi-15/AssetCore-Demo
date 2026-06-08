# Asset-Core Curl Requests

Import the Postman collection file `docs/asset-core.postman_collection.json` directly into Postman. The commands below are equivalent curl requests.

```bash
curl --location 'http://127.0.0.1:8080/health'

curl --location 'http://127.0.0.1:8080/api/v1/assets' \
  --header 'Content-Type: application/json' \
  --data '{"asset_name":"edge-gateway-01","asset_type":"gateway","vendor":"example-vendor","model":"gw-1000","serial_number":"SN-001","mac_address":"00:11:22:33:44:55","ip_address":"192.168.1.10","hostname":"edge-gateway-01","owner_department":"security","owner_user":"admin","location":"shanghai","source":"manual"}'

curl --location 'http://127.0.0.1:8080/api/v1/assets?page=1&page_size=20'

curl --location 'http://127.0.0.1:8080/api/v1/assets/1'

curl --location --request PUT 'http://127.0.0.1:8080/api/v1/assets/1' \
  --header 'Content-Type: application/json' \
  --data '{"asset_name":"edge-gateway-01-updated","ip_address":"192.168.1.11","status":"registered"}'

curl --location --request DELETE 'http://127.0.0.1:8080/api/v1/assets/1'

curl --location 'http://127.0.0.1:8080/api/v1/assets/1/status' \
  --header 'Content-Type: application/json' \
  --data '{"status":"verified"}'

curl --location 'http://127.0.0.1:8080/api/v1/assets/1/changes'

curl --location --request POST 'http://127.0.0.1:8080/api/v1/assets/1/verify'

curl --location 'http://127.0.0.1:8080/api/v1/assets/1/verification-result'

curl --location 'http://127.0.0.1:8080/api/v1/identities/generate' \
  --header 'Content-Type: application/json' \
  --data '{"tenant_id":"default","serial_number":"SN-001","vendor":"example-vendor","model":"gw-1000","mac_address":"00:11:22:33:44:55","ip_address":"192.168.1.10","source":"manual"}'

curl --location 'http://127.0.0.1:8080/api/v1/identities/did:asset:replace_with_generated_identity'

curl --location 'http://127.0.0.1:8080/api/v1/identities/did:asset:replace_with_generated_identity/bind' \
  --header 'Content-Type: application/json' \
  --data '{"asset_id":1}'

curl --location --request POST 'http://127.0.0.1:8080/api/v1/identities/did:asset:replace_with_generated_identity/unbind'

curl --location 'http://127.0.0.1:8080/api/v1/identities/did:asset:replace_with_generated_identity/features'

curl --location 'http://127.0.0.1:8080/api/v1/verifications' \
  --header 'Content-Type: application/json' \
  --data '{"asset_id":1}'

curl --location 'http://127.0.0.1:8080/api/v1/verifications/1'

curl --location 'http://127.0.0.1:8080/api/v1/data/import' \
  --header 'Content-Type: application/json' \
  --data '{"file_name":"assets.csv","file_url":"local://assets.csv","operator_id":"admin"}'

curl --location 'http://127.0.0.1:8080/api/v1/data/import-tasks?page=1&page_size=20'

curl --location 'http://127.0.0.1:8080/api/v1/data/import-tasks/1'

curl --location 'http://127.0.0.1:8080/api/v1/data/import-tasks/1/errors'

curl --location 'http://127.0.0.1:8080/api/v1/data/export/assets'
```
