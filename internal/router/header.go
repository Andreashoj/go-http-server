package router

type HeaderType string

const (
	// Response/Request Headers
	ContentLength HeaderType = "Content-Length"
	ContentType   HeaderType = "Content-Type"
	Host          HeaderType = "Host"

	// Authentication
	Authorization   HeaderType = "Authorization"
	WWWAuthenticate HeaderType = "WWW-Authenticate"

	// Caching
	CacheControl HeaderType = "Cache-Control"
	ETag         HeaderType = "ETag"
	Expires      HeaderType = "Expires"
	LastModified HeaderType = "Last-Modified"

	// CORS
	AccessControlAllowOrigin      HeaderType = "Access-Control-Allow-Origin"
	AccessControlAllowMethods     HeaderType = "Access-Control-Allow-Methods"
	AccessControlAllowHeaders     HeaderType = "Access-Control-Allow-Headers"
	AccessControlAllowCredentials HeaderType = "Access-Control-Allow-Credentials"
	AccessControlExposeHeaders    HeaderType = "Access-Control-Expose-Headers"
	Origin                        HeaderType = "Origin"

	// Redirects
	Location HeaderType = "Location"

	// Content Encoding
	ContentEncoding HeaderType = "Content-Encoding"
	Accept          HeaderType = "Accept"
	AcceptEncoding  HeaderType = "Accept-Encoding"

	// Connection
	KeepAlive        HeaderType = "Keep-Alive"
	TransferEncoding HeaderType = "Transfer-Encoding"

	// Security
	StrictTransportSecurity HeaderType = "Strict-Transport-Security"
	XContentTypeOptions     HeaderType = "X-Content-Type-Options"
	XFrameOptions           HeaderType = "X-Frame-Options"
	XSSProtection           HeaderType = "X-XSS-Protection"
	ContentSecurityPolicy   HeaderType = "Content-Security-Policy"

	// Other Common Headers
	Server     HeaderType = "Server"
	UserAgent  HeaderType = "User-Agent"
	Referer    HeaderType = "Referer"
	Cookie     HeaderType = "Cookie"
	SetCookie  HeaderType = "Set-Cookie"
	Date       HeaderType = "Date"
	Allow      HeaderType = "Allow"
	RetryAfter HeaderType = "Retry-After"
)

type Header interface {
	Add(header HeaderType, value interface{})
	Get() map[HeaderType]string
}

type header struct {
	writer HTTPWriter
	value  map[HeaderType]string
}

func (h *httpWriter) Header() Header {
	return &header{
		writer: h,
		value:  make(map[HeaderType]string, 1),
	}
}

func (h *header) Add(headerType HeaderType, value interface{}) {
	h.value[headerType] = value.(string)
	h.writer.addHeader(h)
}

func (h *header) Get() map[HeaderType]string {
	return h.value
}
