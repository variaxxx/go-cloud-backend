package obs_prometheus

import (
	file_app "cloud/internal/features/file/application"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type FileMetrics struct {
	UploadsTotal           prometheus.Counter
	UploadBytesTotal       prometheus.Counter
	ProcessingStartedTotal prometheus.Counter
	ProcessingFinished     *prometheus.CounterVec
	ProcessingDuration     *prometheus.HistogramVec
}

func NewAPIFileMetrics() *FileMetrics {
	return &FileMetrics{
		UploadsTotal: prometheus.NewCounter(
			prometheus.CounterOpts{
				Namespace: "cloud",
				Name:      "files_uploaded_total",
				Help:      "Total number of successfully uploaded files",
			},
		),
		UploadBytesTotal: prometheus.NewCounter(
			prometheus.CounterOpts{
				Namespace: "cloud",
				Name:      "file_upload_bytes_total",
				Help:      "Total size of successfully uploaded files in bytes",
			},
		),
	}
}

func NewWorkerFileMetrics() *FileMetrics {
	return &FileMetrics{
		ProcessingStartedTotal: prometheus.NewCounter(
			prometheus.CounterOpts{
				Namespace: "cloud",
				Name:      "file_processing_started_total",
				Help:      "Total number of file processing jobs started by the worker",
			},
		),
		ProcessingFinished: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "cloud",
				Name:      "file_processing_finished_total",
				Help:      "Total number of completed file processing jobs by result",
			},
			[]string{"result"},
		),
		ProcessingDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: "cloud",
				Name:      "file_processing_duration_seconds",
				Help:      "Duration of completed file processing jobs in seconds by result",
				Buckets:   prometheus.DefBuckets,
			},
			[]string{"result"},
		),
	}
}

func (m *FileMetrics) UploadSucceeded(sizeBytes int64) {
	if m.UploadsTotal != nil {
		m.UploadsTotal.Inc()
	}
	if m.UploadBytesTotal != nil {
		m.UploadBytesTotal.Add(float64(sizeBytes))
	}
}

func (m *FileMetrics) ProcessingStarted() {
	if m.ProcessingStartedTotal != nil {
		m.ProcessingStartedTotal.Inc()
	}
}

func (m *FileMetrics) ObserveProcessingFinished(
	result file_app.FileProcessingResult,
	duration time.Duration,
) {
	if m.ProcessingFinished == nil || m.ProcessingDuration == nil {
		return
	}

	resultLabel := string(result)
	m.ProcessingFinished.WithLabelValues(resultLabel).Inc()
	m.ProcessingDuration.WithLabelValues(resultLabel).Observe(duration.Seconds())
}