package tracing

const (
	// The type or "kind" of an error (only for event="error" logs). E.g., "Exception", "OSError"
	// type string
	LogErrorKind = "error.kind"

	// For languages that support such a thing (e.g., Java, Python),
	// the actual Throwable/Exception/Error object instance itself.
	// E.g., A java.lang.UnsupportedOperationException instance, a python exceptions.NameError instance
	// type string
	LogErrorObject = "error.object"

	// A stable identifier for some notable moment in the lifetime of a Span. For instance, a mutex lock acquisition or release or the sorts of lifetime events in a browser page load described in the Performance.timing specification. E.g., from Zipkin, "cs", "sr", "ss", or "cr". Or, more generally, "initialized" or "timed out". For errors, "error"
	// type string
	LogEvent = "event"

	// A concise, human-readable, one-line message explaining the event.
	// E.g., "Could not connect to backend", "Cache invalidation succeeded"
	// type string
	LogMessage = "message"

	// A stack trace in platform-conventional format; may or may not pertain to an error. E.g., "File \"example.py\", line 7, in \<module\>\ncaller()\nFile \"example.py\", line 5, in caller\ncallee()\nFile \"example.py\", line 2, in callee\nraise Exception(\"Yikes\")\n"
	// type string
	LogStack = "stack"
)

// LogField LogField
type LogField struct {
	Key   string
	Value string
}

// Log new log.
func Log(key string, val string) LogField {
	return LogField{Key: key, Value: val}
}
