package metrics

import (
	"errors"

	"github.com/prometheus/client_golang/prometheus"
)

type GaugeVecOption struct {
	Namespace string
	Subsystem string
	Name      string
	Help      string
	Labels    []string
}

type GaugeVec struct {
	*GaugeVecOption
	*prometheus.GaugeVec
}

func (opts GaugeVecOption) Build() *GaugeVec {
	if opts.Namespace == "" {
		opts.Namespace = DefaultNamespace
	}

	if opts.Subsystem == "" {
		opts.Subsystem = DefaultSubsystem
	}

	if opts.Name == "" {
		opts.Name = DefaultName
	}

	vec := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: opts.Namespace,
			Subsystem: opts.Subsystem,
			Name:      opts.Name,
			Help:      opts.Help,
		}, opts.Labels)
	prometheus.MustRegister(vec)
	return &GaugeVec{
		GaugeVecOption: &opts,
		GaugeVec:       vec,
	}
}

func NewGaugeVec(name, help string, labels []string) *GaugeVec {
	opts := GaugeVecOption{
		Namespace: DefaultNamespace,
		Name:      name,
		Help:      help,
		Labels:    labels,
	}

	return opts.Build()
}

func (gauge *GaugeVec) Inc(labels ...string) error {
	if len(labels) != len(gauge.GaugeVecOption.Labels) {
		return errors.New("labels count not match")
	}
	gauge.WithLabelValues(labels...).Inc()
	return nil
}

func (gauge *GaugeVec) Add(v float64, labels ...string) error {
	if len(labels) != len(gauge.GaugeVecOption.Labels) {
		return errors.New("labels count not match")
	}
	gauge.WithLabelValues(labels...).Add(v)
	return nil
}

func (gauge *GaugeVec) Set(v float64, labels ...string) error {
	if len(labels) != len(gauge.GaugeVecOption.Labels) {
		return errors.New("labels count not match")
	}
	gauge.WithLabelValues(labels...).Set(v)
	return nil
}
