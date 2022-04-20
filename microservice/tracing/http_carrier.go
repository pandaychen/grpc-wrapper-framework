package tracing

import "net/http"

type HttpCarrier http.Header

//type Header map[string][]string

func (h HttpCarrier) Set(key, val string) {
	http.Header(h).Set(key, val)
}

func (h HttpCarrier) Get(key string) string {
	return http.Header(h).Get(key)
}
