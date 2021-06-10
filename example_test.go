package sqlmetrics_test

import (
	"bytes"
	"context"
	"database/sql"
	"time"

	"github.com/VictoriaMetrics/metrics"
	"github.com/cristalhq/sqlmetrics"
)

func ExampleCollector() {
	db, err := sql.Open("driver", "<some-connection-string>")
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
}
