package obs_prometheus

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Observability struct {
	Registry *prometheus.Registry

	HTTP *HTTPMetrics
}

func Register() *Observability {
	reg := prometheus.NewRegistry()

	http := NewHTTPMetrics()

	reg.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),

		http.RequestsTotal,
		http.RequestDuration,
	)

	return &Observability{
		Registry: reg,
		HTTP:     http,
	}
}

func (o *Observability) Handler() http.Handler {
	return promhttp.HandlerFor(o.Registry, promhttp.HandlerOpts{})
}
