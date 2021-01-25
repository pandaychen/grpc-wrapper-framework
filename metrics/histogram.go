package metrics

import (
	"errors"

	"github.com/prometheus/client_golang/prometheus"
)

type HistogramVecOption struct {
	Namespace string
	Subsystem string
	Name      string
	Help      string
	Labels    []string
	Buckets   []float64
}

type HistogramVec struct {
	*HistogramVecOption
	*prometheus.HistogramVec
}

func (opts HistogramVecOption) Build() *HistogramVec {
	if opts.Namespace == "" {
		opts.Namespace = DefaultNamespace
	}

	if opts.Subsystem == "" {
		opts.Subsystem = DefaultSubsystem
	}

	if opts.Name == "" {
		opts.Name = DefaultName
	}
	vec := prometheus.NewHistogramVec(
		prometheus.HistogramVecOption{
			Namespace: opts.Namespace,
			Subsystem: opts.Subsystem,
			Name:      opts.Name,
			Help:      opts.Help,
			Buckets:   opts.Buckets,
		}, opts.Labels)
	prometheus.MustRegister(vec)
	return &HistogramVec{
		HistogramVec: vec,
	}
}

func NewHistogramVec(name, help string, labels []string, buckets []float64) *HistogramVec {
	opts := HistogramVecOption{
		Namespace: DefaultNamespace,
		Name:      name,
		Help:      help,
		Labels:    labels,
		Buckets:   buckets,
	}

	return opts.Build()
}

func (histogram *HistogramVec) Observe(v float64, labels ...string) error {
	if len(labels) != len(counter.HistogramVecOption.Labels) {
		return errors.New("labels count not match")
	}
	histogram.WithLabelValues(labels...).Observe(v)
	return nil
}
