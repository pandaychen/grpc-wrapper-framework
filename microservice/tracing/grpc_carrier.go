package tracing

type grpcCarrier map[string][]string //equal to metadta.MD

func (g grpcCarrier) Get(key string) string {
	if v, ok := g[key]; ok && len(v) > 0 {
		return v[0]
	}
	return ""
}

func (g grpcCarrier) Set(key, val string) {
	g[key] = append(g[key], val)
}
