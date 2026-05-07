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
	File *FileMetrics
}

func RegisterAPI() *Observability {
	reg := prometheus.NewRegistry()

	http := NewHTTPMetrics()
	file := NewAPIFileMetrics()

	reg.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),

		http.RequestsTotal,
		http.RequestDuration,
		file.UploadsTotal,
		file.UploadBytesTotal,
	)

	return &Observability{
		Registry: reg,
		HTTP:     http,
		File:     file,
	}
}

func RegisterWorker() *Observability {
	reg := prometheus.NewRegistry()

	file := NewWorkerFileMetrics()

	reg.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),

		file.ProcessingStartedTotal,
		file.ProcessingFinished,
		file.ProcessingDuration,
	)

	return &Observability{
		Registry: reg,
		File:     file,
	}
}

func (o *Observability) Handler() http.Handler {
	return promhttp.HandlerFor(o.Registry, promhttp.HandlerOpts{})
}
