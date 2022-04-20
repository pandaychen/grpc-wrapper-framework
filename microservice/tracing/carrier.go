package tracing

import "fmt"

// Carrier 为spanContext存储结构公共接口，类似http-Header、grpc-metadata

// BuiltinFormat is used to demarcate the values within package `trace`
// that are intended for use with the Tracer.Inject() and Tracer.Extract()
// methods.
type BuiltinFormat byte

// support format list
const (
	// HTTPFormat represents Trace as HTTP header string pairs.
	//
	// the HTTPFormat format requires that the keys and values
	// be valid as HTTP headers as-is (i.e., character casing may be unstable
	// and special characters are disallowed in keys, values should be
	// URL-escaped, etc).
	//
	// the carrier must be a `http.Header`.
	HTTPFormat BuiltinFormat = iota
	// GRPCFormat represents Trace as gRPC metadata.
	//
	// the carrier must be a `google.golang.org/grpc/metadata.MD`.
	GRPCFormat
)

// Carrier propagator must convert generic interface{} to something this
// implement Carrier interface, Trace can use Carrier to represents itself.
type Carrier interface {
	Set(key, val string)
	Get(key string) string
}

func main() {
	h := make(HttpCarrier) //must use make！
	h.Set("a", "b")

	fmt.Println(h)
}
