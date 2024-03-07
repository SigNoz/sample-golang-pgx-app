# Todo App REST API using Gin and Pgx

```
go mod tidy
```
to download the required packages and populate the go.sum file.

Set up database and table using the following commands:

```bash
psql -U postgres -h localhost -d postgres -c "CREATE DATABASE tododb"
psql -U postgres -h localhost -d tododb -f ./db.sql
```

Set the environment variables:

```bash
export DATABASE_URL="postgres://user:password@localhost:5432/tododb"
```

## To run and configure app to send data to SigNoz

### For SigNoz Cloud

```bash
SERVICE_NAME=gopgxApp INSECURE_MODE=false OTEL_EXPORTER_OTLP_HEADERS=signoz-access-token=<SIGNOZ-INGESTION-TOKEN> OTEL_EXPORTER_OTLP_ENDPOINT=ingest.{region}.signoz.cloud:443 go run .
```

- Update `<SIGNOZ-INGESTION-TOKEN>` with the ingestion token provided by SigNoz
- Update `ingest.{region}.signoz.cloud:443` with the ingestion endpoint of your region. Refer to the table below for the same.

| Region | Endpoint                   |
| ------ | -------------------------- |
| US     | ingest.us.signoz.cloud:443 |
| IN     | ingest.in.signoz.cloud:443 |
| EU     | ingest.eu.signoz.cloud:443 |

### For SigNoz OSS

``` bash
SERVICE_NAME=gopgxApp INSECURE_MODE=true OTEL_EXPORTER_OTLP_ENDPOINT=<IP of SigNoz backend>:4317 go run .
```

- `<IP of SigNoz backend:4317>` should be without http/https scheme. Eg `localhost:4317`.

---

This runs the gin application at port `8080`. Try accessing API at `http://localhost:8080/todos`

Below are the apis available to play around. The API calls will generate telemetry data which will be sent to SigNoz which can be viewed at `<IP of SigNoz backend>:3000`

```
GET    /todos                    
POST   /todos                    
PATCH  /todos/:id                
DELETE /todos/:id                
```