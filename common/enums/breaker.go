package enums

type BreakerType string

const (
	GOOGLE_SRE_BREAKER_TYPE BreakerType = "google"
	GOBREAKER_BREAKER_TYPE  BreakerType = "gobreaker"
)
