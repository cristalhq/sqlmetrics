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
	_ = makeCollector(t, &mockStatser{}, time.Nanosecond, "db", "mydb", "table", "mytable")

	time.Sleep(time.Second) // get some time to collect metrics

	b := &bytes.Buffer{}
	metrics.WritePrometheus(b, false)

	got := b.String()
	want := `go_sql_idle{db="mydb",table="mytable"} 4
go_sql_in_use{db="mydb",table="mytable"} 3
go_sql_max_idle_closed{db="mydb",table="mytable"} 7
go_sql_max_idletime_closed{db="mydb",table="mytable"} 8
go_sql_max_lifetime_closed{db="mydb",table="mytable"} 9
go_sql_max_open{db="mydb",table="mytable"} 1
go_sql_open{db="mydb",table="mytable"} 2
go_sql_wait_count{db="mydb",table="mytable"} 5
go_sql_wait_duration_seconds{db="mydb",table="mytable"} 6
`
	if got != want {
		t.Fatalf("got %s\n want %s", got, want)
	}
}

func TestPassSQL(t *testing.T) {
	_ = makeCollector(t, &sql.DB{}, time.Second, "sql", "best", "label", "value")
}

func TestBadLabels(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("must panic")
		}
	}()

	_ = makeCollector(t, &mockStatser{}, time.Second, "mock", "stub", "onlyone")
}

type mockStatser struct{}

func (m *mockStatser) Stats() sql.DBStats {
	return sql.DBStats{
		MaxOpenConnections: 1,
		OpenConnections:    2,
		InUse:              3,
		Idle:               4,
		WaitCount:          5,
		WaitDuration:       6 * time.Second,
		MaxIdleClosed:      7,
		MaxIdleTimeClosed:  8,
		MaxLifetimeClosed:  9,
	}
}

func makeCollector(t testing.TB, db Statser, every time.Duration, labels ...string) *Collector {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	c := NewCollector(ctx, db, every, labels...)

	t.Cleanup(func() {
		cancel()
		unregisterCollectorMetrics(labels...)
	})
	return c
}

func unregisterCollectorMetrics(labels ...string) {
	names := []string{
		"max_open",
		"open",
		"in_use",
		"idle",
		"wait_count",
		"wait_duration_seconds",
		"max_idle_closed",
		"max_idletime_closed",
		"max_lifetime_closed",
	}

	allLabels := buildLabels(labels...)
	for _, name := range names {
		metrics.UnregisterMetric(buildName(name, allLabels))
	}
}
