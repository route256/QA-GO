# Ozon Marketplace Template API

---

## Build project

### Local

For local assembly you need to perform

```zsh
$ make deps # Installation of dependencies
$ make build # Build project
```
## Running

### For local development
 - up
   ```zsh
   make dc-up
   ```
 - down
   ```zsh
   make dc-down
   ```
 - rebuild app image and reUp compose
   ```shell
   make dc-rebuild-reup
   ```
---

## Services

### Swagger UI

The Swagger UI is an open source project to visually render documentation for an API defined with the OpenAPI (Swagger) Specification

- http://localhost:8081
- http://localhost:8080/swagger

### Grafana:

- http://localhost:3000
- - login `admin`
- - password `MYPASSWORT`

### gRPC:

- http://localhost:8082

```sh
[I] ➜ grpc_cli ls localhost:8082
response:
grpc.reflection.v1alpha.ServerReflection
ozonmp.act_device_api.v1.ActDeviceApiService

grpc_cli ls localhost:8082 ozonmp.act_device_api.v1.ActDeviceApiService -l     


[I] ➜ grpc_cli call localhost:8082 DescribeDeviceV1 "device_id: 1"
connecting to localhost:8082
value {
  id: 1
  platform: "Ios"
  user_id: 24154
  entered_at {
    seconds: 1636465909
    nanos: 870041000
  }
}
Rpc succeeded with OK status
```

```sh
[I] ➜ grpcurl -plaintext localhost:8082 list
grpc.reflection.v1alpha.ServerReflection
ozonmp.act_device_api.v1.ActDeviceApiService

[I] ➜ grpcurl -plaintext localhost:8082 list ozonmp.act_device_api.v1.ActDeviceApiService

grpcurl -plaintext -emit-defaults -d '{"device_id": 1}' localhost:8082 ozonmp.act_device_api.v1.ActDeviceApiService.DescribeDeviceV1

response
{
  "value": {
    "id": "1",
    "platform": "Ios",
    "userId": "24154",
    "enteredAt": "2021-11-09T13:51:49.870041Z"
  }
}

```

### Gateway:

It reads protobuf service definitions and generates a reverse-proxy server which translates a RESTful HTTP API into gRPC

- http://localhost:8080

```sh
[I] ➜ curl -X 'GET' \
  'http://localhost:8080/api/v1/devices/1' \
  -H 'accept: application/json' | jq .
{
  "value": {
    "id": "1",
    "platform": "Ios",
    "userId": "24154",
    "enteredAt": "2021-11-09T13:51:49.870041Z"
  }
}
```

### Metrics:

Metrics GRPC Server

- http://localhost:9100/metrics

### Status:

Service condition and its information

- http://localhost:8000
- - `/live`- Layed whether the server is running
- - `/ready` - Is it ready to accept requests
- - `/version` - Version and assembly information

### Prometheus:

Prometheus is an open-source systems monitoring and alerting toolkit

- http://localhost:9090

### Kafka

Apache Kafka is an open-source distributed event streaming platform used by thousands of companies for high-performance data pipelines, streaming analytics, data integration, and mission-critical applications.

- http://localhost:9094

### Kafka UI

UI for Apache Kafka is a simple tool that makes your data flows observable, helps find and troubleshoot issues faster and deliver optimal performance. Its lightweight dashboard makes it easy to track key metrics of your Kafka clusters - Brokers, Topics, Partitions, Production, and Consumption.

- http://localhost:9001

### Jaeger UI

Monitor and troubleshoot transactions in complex distributed systems.

- http://localhost:16686

### Graylog

Graylog is a leading centralized log management solution for capturing, storing, and enabling real-time analysis of terabytes of machine data.

- http://localhost:9000
- - login `admin`
- - password `admin`

### PostgreSQL

For the convenience of working with the database, you can use the [pgcli](https://github.com/dbcli/pgcli) utility. Migrations are rolled out when the service starts. migrations are located in the **./migrations** directory and are created using the [goose](https://github.com/pressly/goose) tool.

```sh
$ pgcli "postgresql://docker:docker@localhost:5432/omp_template_api"

goose -dir migrations postgres "postgres://postgres:password@localhost:5432/act_device_api?sslmode=disable" status   
```

### Thanks

- [Evald Smalyakov](https://github.com/evald24)
- [Michael Morgoev](https://github.com/zerospiel)
