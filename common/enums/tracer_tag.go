package enums

//定义了各类标准的tag名字

// Standard Span tags https://github.com/opentracing/specification/blob/master/semantic_conventions.md#span-tags-table
const (
	// The software package, framework, library, or module that generated the associated Span.
	// E.g., "grpc", "django", "JDBI".
	// type string
	TagComponent = "component"

	// Database instance name.
	// E.g., In java, if the jdbc.url="jdbc:mysql://127.0.0.1:3306/customers", the instance name is "customers".
	// type string
	TagDBInstance = "db.instance"

	// A database statement for the given database type.
	// E.g., for db.type="sql", "SELECT * FROM wuser_table"; for db.type="redis", "SET mykey 'WuValue'".
	TagDBStatement = "db.statement"

	// Database type. For any SQL database, "sql". For others, the lower-case database category,
	// e.g. "cassandra", "hbase", or "redis".
	// type string
	TagDBType = "db.type"

	// Username for accessing database. E.g., "readonly_user" or "reporting_user"
	// type string
	TagDBUser = "db.user"

	// true if and only if the application considers the operation represented by the Span to have failed
	// type bool
	TagError = "error"

	// HTTP method of the request for the associated Span. E.g., "GET", "POST"
	// type string
	TagHTTPMethod = "http.method"

	// HTTP response status code for the associated Span. E.g., 200, 503, 404
	// type integer
	TagHTTPStatusCode = "http.status_code"

	// URL of the request being handled in this segment of the trace, in standard URI format.
	// E.g., "https://domain.net/path/to?resource=here"
	// type string
	TagHTTPURL = "http.url"

	// An address at which messages can be exchanged.
	// E.g. A Kafka record has an associated "topic name" that can be extracted by the instrumented producer or consumer and stored using this tag.
	// type string
	TagMessageBusDestination = "message_bus.destination"

	// Remote "address", suitable for use in a networking client library.
	// This may be a "ip:port", a bare "hostname", a FQDN, or even a JDBC substring like "mysql://prod-db:3306"
	// type string
	TagPeerAddress = "peer.address"

	// 	Remote hostname. E.g., "opentracing.io", "internal.dns.name"
	// type string
	TagPeerHostname = "peer.hostname"

	// Remote IPv4 address as a .-separated tuple. E.g., "127.0.0.1"
	// type string
	TagPeerIPv4 = "peer.ipv4"

	// Remote IPv6 address as a string of colon-separated 4-char hex tuples.
	// E.g., "2001:0db8:85a3:0000:0000:8a2e:0370:7334"
	// type string
	TagPeerIPv6 = "peer.ipv6"

	// Remote port. E.g., 80
	// type integer
	TagPeerPort = "peer.port"

	// Remote service name (for some unspecified definition of "service").
	// E.g., "elasticsearch", "a_custom_microservice", "memcache"
	// type string
	TagPeerService = "peer.service"

	// If greater than 0, a hint to the Tracer to do its best to capture the trace.
	// If 0, a hint to the trace to not-capture the trace. If absent, the Tracer should use its default sampling mechanism.
	// type string
	TagSamplingPriority = "sampling.priority"

	// Either "client" or "server" for the appropriate roles in an RPC,
	// and "producer" or "consumer" for the appropriate roles in a messaging scenario.
	// type string
	TagSpanKind = "span.kind"

	// legacy tag
	TagAnnotation = "legacy.annotation"
	TagAddress    = "legacy.address"
	TagComment    = "legacy.comment"
)

//user define
const (
	TagDBExecuteCosts = "db.costs"
)