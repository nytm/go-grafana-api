# grafana-api-golang-client

[![Build Status](https://cloud.drone.io/api/badges/nytm/go-grafana-api/status.svg)](https://cloud.drone.io/nytm/go-grafana-api)

Grafana HTTP API Client for Go

## Tests

To run unit tests:

```
make test
```

To run integration tests:

Start a `localhost:3000` Grafana:

```
make serve-grafana
```

Run the integration tests:

```
make integration-test:
```
