# isso - Iterative Sampling Schedule Optimization

[![Test status](https://img.shields.io/github/actions/workflow/status/mlange-42/isso/tests.yml?branch=main&label=Tests&logo=github)](https://github.com/mlange-42/isso/actions/workflows/tests.yml)
[![Coverage Status](https://img.shields.io/coverallsCoverage/github/mlange-42/isso?logo=coveralls)](https://badge.coveralls.io/github/mlange-42/isso?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/mlange-42/isso)](https://goreportcard.com/report/github.com/mlange-42/isso)
[![Go Reference](https://img.shields.io/badge/reference-%23007D9C?logo=go&logoColor=white&labelColor=gray)](https://pkg.go.dev/github.com/mlange-42/isso)
[![GitHub](https://img.shields.io/badge/github-repo-blue?logo=github)](https://github.com/mlange-42/isso)
[![MIT license](https://img.shields.io/badge/MIT-brightgreen?label=license)](https://github.com/mlange-42/isso/blob/main/LICENSE)

isso is a Go library and CLI app for optimizing sampling schedules under constrained capacity and with potential sample re-use.

## CLI usage

Run the included examples like this...

The default test problem:

```
go run ./cmd/isso -i data/problem.json
```

A pareto optimization example:

```
go run ./cmd/isso -i data/pareto.json --pareto --format fitness
```

See folder `data` for problem definition examples.

## License

This project is distributed under the [MIT license](./LICENSE).
