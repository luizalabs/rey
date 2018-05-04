package metric

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	mu              = new(sync.Mutex)
	registredGauges = map[string]prometheus.Gauge{}
)

type Gauge struct {
	pg prometheus.Gauge
}

func (g *Gauge) Set(n float64) {
	g.pg.Set(n)
}

func NewGauge(compName string) *Gauge {
	return &Gauge{pg: getOrCreatePG(compName)}
}

func getOrCreatePG(compName string) prometheus.Gauge {
	mu.Lock()
	defer mu.Unlock()

	if g, found := registredGauges[compName]; found {
		return g
	}
	g := prometheus.NewGauge(prometheus.GaugeOpts{
		Name:        "rey",
		Help:        "Doctor Rey the Health Checker",
		ConstLabels: prometheus.Labels{"component_name": compName},
	})
	prometheus.Register(g)
	registredGauges[compName] = g

	return g
}
