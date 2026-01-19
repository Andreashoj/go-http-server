package router

import "fmt"

type Request string

const (
	Post   Request = "POST"
	Get    Request = "GET"
	Put    Request = "PUT"
	Delete Request = "DELETE"
)

type HTTPRequest interface {
	Params() map[string]string
	Body() string
	GetQueryParam(key string) (string, error)
	Url() string
}

type httpRequest struct {
	startLine string
	headers   []string
	body      string
	params    map[string]string
	url       string
}

func NewHTTPRequest() HTTPRequest {
	return &httpRequest{
		params: make(map[string]string),
	}
}

func (r *httpRequest) Params() map[string]string {
	return r.params
}

func (r *httpRequest) GetQueryParam(key string) (string, error) {
	value, exists := r.params[key]
	if !exists {
		return "", fmt.Errorf("query param not present: %s in parameter list: %s", key, r.params)
	}

	return value, nil
}

func (r *httpRequest) Body() string {
	return r.body
}

func (r *httpRequest) Url() string {
	return r.url
}
