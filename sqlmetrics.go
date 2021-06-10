package sqlmetrics

import (
	"database/sql"
	"strings"
	"time"

	"github.com/VictoriaMetrics/metrics"
)

const namespace = "go_sql"

// Statser is an interface that gets sql.DBStats.
// Most of the DB clients support this.
type Statser interface {
	Stats() sql.DBStats
}

// Collector for the sql.DBStats for Prometheus (via VictoriaMetrics client).
type Collector struct {
	maxOpen           *metrics.Gauge
	open              *metrics.Gauge
	inUse             *metrics.Gauge
	idle              *metrics.Gauge
	waitCount         *metrics.Counter
	waitDuration      *metrics.Counter
	maxIdleClosed     *metrics.Counter
	maxIdleTimeClosed *metrics.Counter
	maxLifetimeClosed *metrics.Counter
}

// NewCollector creates a new Collector.
func NewCollector(stats Statser, every time.Duration, labels ...string) *Collector {
	allLabels := buildLabels(labels...)

	c := &Collector{
		maxOpen: newGauge("max_open", allLabels, func() float64 {
			return float64(stats.Stats().MaxOpenConnections)
		}),
		open: newGauge("open", allLabels, func() float64 {
			return float64(stats.Stats().OpenConnections)
		}),
		inUse: newGauge("in_use", allLabels, func() float64 {
			return float64(stats.Stats().InUse)
		}),
		idle: newGauge("idle", allLabels, func() float64 {
			return float64(stats.Stats().Idle)
		}),

		waitCount:         newCounter("wait_count", allLabels),
		waitDuration:      newCounter("wait_duration_seconds", allLabels),
		maxIdleClosed:     newCounter("max_idle_closed", allLabels),
		maxIdleTimeClosed: newCounter("max_idletime_closed", allLabels),
		maxLifetimeClosed: newCounter("max_lifetime_closed", allLabels),
	}

	go func() {
		for range time.NewTicker(every).C {
			s := stats.Stats()
			c.waitCount.Set(uint64(s.WaitCount))
			c.waitDuration.Set(uint64(s.WaitDuration.Seconds()))
			c.maxIdleClosed.Set(uint64(s.MaxIdleClosed))
			c.maxIdleTimeClosed.Set(uint64(s.MaxIdleTimeClosed))
			c.maxLifetimeClosed.Set(uint64(s.MaxLifetimeClosed))
		}
	}()
	return c
}

func newGauge(name, labels string, f func() float64) *metrics.Gauge {
	return metrics.NewGauge(buildName(name, labels), f)
}

func newCounter(name, labels string) *metrics.Counter {
	return metrics.NewCounter(buildName(name, labels))
}

func buildName(name, labels string) string {
	return namespace + "_" + name + labels
}

func buildLabels(labels ...string) string {
	if len(labels) == 0 {
		return ""
	}
	if len(labels)%2 != 0 {
		panic("dbpstats: incorrect label pairs")
	}

	var b strings.Builder
	b.WriteByte('{')
	for i := 0; i < len(labels); i += 2 {
		if i != 0 {
			b.WriteString(`",`)
		}
		b.WriteString(labels[i])
		b.WriteString(`="`)
		b.WriteString(labels[i+1])
	}
	b.WriteString(`"}`)
	return b.String()
}
