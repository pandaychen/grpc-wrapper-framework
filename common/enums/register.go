package enums

type RegType string

const (
	RET_TYPE_DIRECT RegType = "direct"
	REG_TYPE_ETCD   RegType = "etcd"
	REG_TYPE_CONSUL RegType = "consul"
	REG_TYPE_DNS    RegType = "dns"
)
