# AegisFlow

Runtime policy gateway for coding agents and tool-using agents.

AegisFlow evaluates routed MCP and API actions as `allow`, `review`, or `block`. It supports human approval, scoped credentials, and signed session evidence.

## Run

```bash
docker run --rm \
  -p 8080:8080 \
  -p 8081:8081 \
  -p 8082:8082 \
  saivedant169/aegisflow:0.9.0
```

Default image uses local configuration and mock provider, so startup does not require a provider key.

## Health Check

The image includes a built-in health probe to monitor the gateway's status. It probes `http://127.0.0.1:8080/health` with the following configuration:
- **Interval**: 30s
- **Timeout**: 5s
- **Start Period**: 10s
- **Retries**: 3

## Ports

| Port | Service |
|---|---|
| `8080` | Gateway API |
| `8081` | Admin API, dashboard, metrics |
| `8082` | MCP gateway |

## Mount configuration

```bash
docker run --rm \
  -p 8080:8080 \
  -p 8081:8081 \
  -p 8082:8082 \
  -v "$PWD/aegisflow.yaml:/app/configs/aegisflow.yaml:ro" \
  saivedant169/aegisflow:0.9.0
```

## Documentation

- [Quickstart](https://github.com/saivedant169/AegisFlow#start-here)
- [Configuration](https://github.com/saivedant169/AegisFlow/blob/main/configs/aegisflow.example.yaml)
- [Production checklist](https://github.com/saivedant169/AegisFlow/blob/main/docs/production-checklist.md)
- [Release verification](https://github.com/saivedant169/AegisFlow/releases/latest)

Images also publish at `ghcr.io/saivedant169/aegisflow`.
