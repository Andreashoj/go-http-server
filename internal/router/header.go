package router

type HeaderType string

const (
	ContentLength HeaderType = "Content-Length"
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
