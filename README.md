# sqlmetrics

[![build-img]][build-url]
[![pkg-img]][pkg-url]
[![reportcard-img]][reportcard-url]
[![coverage-img]][coverage-url]

Prometheus metrics for Go `database/sql` via [VictoriaMetrics/metrics](https://github.com/VictoriaMetrics/metrics)

## Features

* Simple API.
* Easy to integrate.

## Install

Go version 1.16+

```
go get github.com/cristalhq/sqlmetrics
```

## Example

```go
import (
    "github.com/VictoriaMetrics/metrics"
    "github.com/cristalhq/sqlmetrics"
)

// ...

db, err := sql.Open("<some-connection-string>")
if err != nil {
    panic(err)
}

ctx := context.Background() // or any other context you have
every := 3 * time.Second

sqlmetrics.NewCollector(ctx, db, every, "label1", "value1", "another", "etc")

// done, db metrics are registered
// you can see them here
w := &bytes.Buffer{}
metrics.WritePrometheus(w, true)
```

See this examples: [example_test.go](https://github.com/cristalhq/sqlmetrics/blob/main/example_test.go).

## Documentation

See [these docs][pkg-url].

## License

[MIT License](LICENSE).

[build-img]: https://github.com/cristalhq/sqlmetrics/workflows/build/badge.svg
[build-url]: https://github.com/cristalhq/sqlmetrics/actions
[pkg-img]: https://pkg.go.dev/badge/cristalhq/sqlmetrics
[pkg-url]: https://pkg.go.dev/github.com/cristalhq/sqlmetrics
[reportcard-img]: https://goreportcard.com/badge/cristalhq/sqlmetrics
[reportcard-url]: https://goreportcard.com/report/cristalhq/sqlmetrics
[coverage-img]: https://codecov.io/gh/cristalhq/sqlmetrics/branch/main/graph/badge.svg
[coverage-url]: https://codecov.io/gh/cristalhq/sqlmetrics
