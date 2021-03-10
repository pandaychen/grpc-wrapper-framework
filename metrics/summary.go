package metrics

import (
	"errors"

	"github.com/prometheus/client_golang/prometheus"
)

type SummaryVecOption struct {
	Namespace string
	Subsystem string
	Name      string
	Help      string
	Labels    []string
}

type SummaryVec struct {
	*SummaryVecOption
	*prometheus.SummaryVec
}

func (opts SummaryVecOption) Build() *SummaryVec {
	if opts.Namespace == "" {
		opts.Namespace = DefaultNamespace
	}

	if opts.Subsystem == "" {
		opts.Subsystem = DefaultSubsystem
	}

	if opts.Name == "" {
		opts.Name = DefaultName
	}
	vec := prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace: opts.Namespace,
			Subsystem: opts.Subsystem,
			Name:      opts.Name,
			Help:      opts.Help,
		}, opts.Labels)
	prometheus.MustRegister(vec)
	return &SummaryVec{
		SummaryVecOption: &opts,
		SummaryVec:       vec,
	}
}

func (summary *SummaryVec) Observe(v float64, labels ...string) error {
	if len(labels) != len(summary.SummaryVecOption.Labels) {
		return errors.New("labels count not match")
	}
	summary.WithLabelValues(labels...).Observe(v)

	return nil
}
