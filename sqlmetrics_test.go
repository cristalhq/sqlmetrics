package sqlmetrics

import (
	"bytes"
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/VictoriaMetrics/metrics"
)

func TestCollector(t *testing.T) {
	NewCollector(context.Background(), &mockStatser{}, time.Nanosecond, "db", "table", "key", "value")

	time.Sleep(2 * time.Second) // get some time to collect metrics

	b := &bytes.Buffer{}
	metrics.WritePrometheus(b, false)

	got := b.String()
	want := `go_sql_idle{db="table",key="value"} 4
go_sql_in_use{db="table",key="value"} 3
go_sql_max_idle_closed{db="table",key="value"} 7
go_sql_max_idletime_closed{db="table",key="value"} 8
go_sql_max_lifetime_closed{db="table",key="value"} 9
go_sql_max_open{db="table",key="value"} 1
go_sql_open{db="table",key="value"} 2
go_sql_wait_count{db="table",key="value"} 5
go_sql_wait_duration_seconds{db="table",key="value"} 6
`

	if want != got {
		t.Fatalf("want %q, got %q", want, got)
	}
}

func TestPassSQL(t *testing.T) {
	NewCollector(context.Background(), &sql.DB{}, 1, "sql", "best", "a", "b")
}

func TestBadLabels(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("must panic")
		}
	}()

	NewCollector(context.Background(), &mockStatser{}, 1, "mock", "stub", "onlyone")
}

type mockStatser struct{}

func (m *mockStatser) Stats() sql.DBStats {
	return sql.DBStats{
		MaxOpenConnections: 1,
		OpenConnections:    2,
		InUse:              3,
		Idle:               4,
		WaitCount:          5,
		WaitDuration:       6,
		MaxIdleClosed:      7,
		MaxIdleTimeClosed:  8,
		MaxLifetimeClosed:  9,
	}
}
