---
title: Performance
description: Reproduce AegisFlow HTTP load tests and governance microbenchmarks.
---

# Performance

Results below were measured on Apple M1 running macOS 26.6 and Go 1.26.6 on 2026-08-18. Run same commands on target hardware before capacity planning.

## Zero-latency provider

`scripts/loadtest_e2e.sh` drives actual HTTP server, JSON handling, routing, input policy, usage tracking, and zero-latency mock provider.

Command: `hey -n 30000 -c 50`

| Metric | Result |
|---|---:|
| Throughput | 54,860 requests per second |
| p50 | 0.6 ms |
| p95 | 2.4 ms |
| p99 | 3.5 ms |
| Errors | 0 of 30,000 |

```bash
go install github.com/rakyll/hey@v0.1.5
./scripts/loadtest_e2e.sh
```

## Provider-latency test

`scripts/benchmark.sh` disables cache and sends every request to mock provider with fixed 25 ms delay.

Command: `hey -n 300 -c 20`

| Metric | Result |
|---|---:|
| Throughput | 611.75 requests per second |
| p50 | 27.9 ms |
| p95 | 34.0 ms |
| p99 | 36.4 ms |
| Error rate | 0.00% |

```bash
BENCH_RESULTS_FILE=benchmark-results.txt bash scripts/benchmark.sh
```

## Governance microbenchmarks

Go benchmark results:

| Scenario | Time per operation | Allocations |
|---|---:|---:|
| Envelope creation | 402 ns | 4 |
| Policy allow, 20 rules | 666 ns | 4 |
| Policy block, no match | 577 ns | 4 |
| Policy plus evidence | 2.634 us | 19 |
| Full allow with benchmark credential | 2.905 us | 22 |
| Review queue submission | 1.173 us | 5 |

```bash
go test ./scripts/benchgovern/ -bench=Benchmark -benchmem -count=1 -run='^$'
```

Standalone percentile runner reported 0.67 us p50 for policy evaluation, 2.62 us p50 for policy plus evidence, and 2.58 us p50 for full allow path across 10,000 iterations:

```bash
go run ./scripts/benchmark_governance.go
```

Microbenchmark operations per second and HTTP throughput measure different paths. Provider latency, policy count, evidence backend, telemetry exporters, TLS, and concurrency affect deployment result.
