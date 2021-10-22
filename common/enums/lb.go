package enums

type LbType string

const (
	LB_TYPE_RR   LbType = "rr"
	LB_TYPE_WRR  LbType = "wrr"
	LB_TYPE_RAND LbType = "random"
	LB_TYPE_P2C  LbType = "p2c"
	LB_TYPE_CH   LbType = "ch"
)
