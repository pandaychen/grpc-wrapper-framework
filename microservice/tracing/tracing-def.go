package tracing

// 公共定义

/*
    [Span A]  ←←←(the root span)
            |
     +------+------+
     |             |
 [Span B]      [Span C] ←←←(Span C is a `ChildOf` Span A)
     |             |
 [Span D]      +---+-------+
               |           |
           [Span E]    [Span F] >>> [Span G] >>> [Span H]
                                       ↑
                                       ↑
                                       ↑
                         (Span G `FollowsFrom` Span F)

*/

/*
   Client SpanTrace                                            Server SpanTrace
┌──────────────────┐                                       ┌──────────────────┐
│                  │                                       │                  │
│    SpanContext   │           Http Request Headers        │   SpanContext    │
│ ┌──────────────┐ │          ┌───────────────────┐        │ ┌──────────────┐ │
│ │ TraceId      │ │          │ X-B3-TraceId      │        │ │ TraceId      │ │
│ │              │ │          │                   │        │ │              │ │
│ │ ParentSpanId │ │ Inject   │ X-B3-ParentSpanId │Extract │ │ ParentSpanId │ │
│ │              ├─┼─────────>│                   ├────────┼>│              │ │
│ │ SpanId       │ │          │ X-B3-SpanId       │        │ │ SpanId       │ │
│ │              │ │          │                   │        │ │              │ │
│ │ Sampled      │ │          │ X-B3-Sampled      │        │ │ Sampled      │ │
│ └──────────────┘ │          └───────────────────┘        │ └──────────────┘ │
│                  │                                       │                  │
└──────────────────┘                                       └──────────────────┘
*/

// SpanTrace 为 Span 的公共接口
type SpanTrace interface {
	// traceid包含 traceid+spanid+parentid的组合
	TraceID() string

	// Fork fork a trace with client trace.
	Fork(serviceName, operationName string) SpanTrace

	// Follow
	Follow(serviceName, operationName string) SpanTrace

	// Finish when trace finish call it.
	Finish(err *error)

	// Adds a tag to the trace.
	//
	// If there is a pre-existing tag set for `key`, it is overwritten.
	//
	// Tag values can be numeric types, strings, or bools. The behavior of
	// other tag value types is undefined at the OpenTracing level. If a
	// tracing system does not know how to handle a particular value type, it
	// may ignore the tag, but shall not panic.
	// NOTE current only support legacy tag: TagAnnotation TagAddress TagComment
	// other will be ignore
	SetTag(tags ...Tag) SpanTrace

	// LogFields is an efficient and type-checked way to record key:value
	// NOTE current unsupport
	SetLog(logs ...LogField) SpanTrace

	// Visit visits the k-v pair in trace, calling fn for each.
	Visit(fn func(k, v string))

	// SetTitle reset trace title
	SetTitle(title string)
}
