package metrics

// Counter 计数器

import (
	"errors"

	"github.com/prometheus/client_golang/prometheus"
)

// CounterVecOpts 选项
type CounterVecOption struct {
	Namespace string
	Subsystem string
	Name      string
	Help      string
	Labels    []string
}

// 封装prometheus.CounterVec
type CounterVec struct {
	*CounterVecOption
	*prometheus.CounterVec
}

func (opts CounterVecOption) Build() *CounterVec {
	if opts.Namespace == "" {
		opts.Namespace = DefaultNamespace
	}

	if opts.Subsystem == "" {
		opts.Subsystem = DefaultSubsystem
	}

	if opts.Name == "" {
		opts.Name = DefaultName
	}

	vec := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: opts.Namespace,
			Subsystem: opts.Subsystem,
			Name:      opts.Name,
			Help:      opts.Help,
		}, opts.Labels)

	prometheus.MustRegister(vec)
	return &CounterVec{
		CounterVecOption: &opts,
		CounterVec:       vec,
	}
}

// 创建vec
func NewCounterVec(name, help string, labels []string) *CounterVec {
	opts := CounterVecOption{
		Namespace: DefaultNamespace,
		Name:      name,
		Help:      help,
		Labels:    labels,
	}

	return opts.Build()
}

func (counter *CounterVec) Inc(labels ...string) error {
	if len(labels) != len(counter.CounterVecOption.Labels) {
		return errors.New("labels count not match")
	}
	counter.WithLabelValues(labels...).Inc()

	return nil
}

func (counter *CounterVec) Add(v float64, labels ...string) error {
	if len(labels) != len(counter.CounterVecOption.Labels) {
		return errors.New("labels count not match")
	}
	counter.WithLabelValues(labels...).Add(v)

	return nil
}
