# Financial Crime Detector

## Getting Started

### Prerequisites
- Docker
- Docker Compose
- Basic knowledge of Go, React, and Docker

### Installation

1. **Clone the repository:**
    ```bash
    git clone https://github.com/Abasi-ifreke/fincrime-detector.git
    cd fincrime-detector
    ```

2. **Set up the `.env` file:**
    - Create a `.env` file in the root directory.
    - Copy and modify the provided `.env` with your environment-specific values.
    - Ensure you set secure passwords for `ELASTIC_PASSWORD` and `KIBANA_PASSWORD`.

    ```bash
    cp .env .env
    # Modify .env with your actual environment variables
    ```

3. **Start the services:**
    ```bash
    docker compose up -d
    ```

This command will build the necessary images and start the containers in the correct order.

### Architecture



### Access the application
- Frontend: [http://localhost:80](http://localhost:80)
- Kibana: [http://localhost:5601](http://localhost:5601)

## Important Notes on Certificates
- The system uses TLS certificates for secure communication with Elasticsearch.
- Certificates are generated during the setup phase of the `docker-compose` process.
- The `setup` service generates a Certificate Authority (CA) and server certificates.
- These certificates are stored in the `config/certs` directory and are used by:
  - Elasticsearch
  - Kibana
  - Backend service

The backend trusts the CA certificate using the `RootCAs` field in the `tls.Config` when creating the Elasticsearch client (`initElasticsearch()` function in `main.go`).

> Failure to properly configure certificates will result in communication errors between the backend and Elasticsearch.

## API Documentation

### Transactions

#### `POST /api/transactions`
- Accepts a JSON payload representing a transaction.
- Requires `Authorization` header (Basic Auth).

**Example:**
```bash
curl -X POST \
  http://localhost:8080/api/transactions \
  -H 'Authorization: Basic dGVzdHVzZXI6dGVzdHBhc3N3b3Jk' \
  -H 'Content-Type: application/json' \
  -d '{
    "id": "transaction_id_123",
    "account_id": "account_id_456",
    "transaction_type": "transfer",
    "amount": 1000.00,
    "timestamp": "2024-07-24T12:00:00Z"
  }'
```

## Alerts

### `GET /api/alerts`
Retrieves a list of triggered alerts.

- Requires `Authorization` header (Basic Authentication).

**Example:**
```bash
curl -X GET \
  http://localhost:8080/api/alerts \
  -H 'Authorization: Basic dGVzdHVzZXI6dGVzdHBhc3N3b3Jk'
```

## Design Considerations

- **Scalability**: Designed to handle high volumes of transactions using Elasticsearch’s distributed architecture and Go’s concurrency features.
- **Real-time Performance**: Elasticsearch enables fast indexing and querying of transactions and alerts.
- **Extensibility**: New detection rules can be easily added to the PostgreSQL database without requiring code changes.
- **Maintainability**: Modular architecture with comprehensive logging simplifies maintenance and debugging.

---