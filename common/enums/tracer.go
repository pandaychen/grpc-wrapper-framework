package enums

type TracerType string

const (
	TRACER_TYPE_JAEGER TracerType = "jaeger"
	TRACER_TYPE_ZIPKIN TracerType = "zipkin"
)
