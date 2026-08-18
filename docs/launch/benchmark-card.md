# AegisFlow benchmark card

Measured on Apple M1, macOS 26.6, Go 1.26.6, 2026-08-18.

```text
AegisFlow gateway benchmark

Zero-latency mock provider, 30,000 requests, concurrency 50
  Throughput       54,860 req/s
  p50 / p95 / p99 0.6 / 2.4 / 3.5 ms
  Errors           0

25 ms mock provider, cache disabled, 300 requests, concurrency 20
  Throughput       611.75 req/s
  p50 / p95 / p99 27.9 / 34.0 / 36.4 ms
  Error rate       0.00%

Governance microbenchmarks
  Policy allow     666 ns/op
  Policy + evidence 2.634 us/op
  Full allow path  2.905 us/op
```

Reproduce:

```bash
./scripts/loadtest_e2e.sh
BENCH_RESULTS_FILE=benchmark-results.txt bash scripts/benchmark.sh
go test ./scripts/benchgovern/ -bench=Benchmark -benchmem -count=1 -run='^$'
```

HTTP tests use local mock providers. They do not predict external provider latency. See [performance notes](../performance.md) for method and limits.
