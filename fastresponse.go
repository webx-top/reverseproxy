package reverseproxy

import (
	"github.com/valyala/fasthttp"
)

var _ Context = &FastResponse{}

type FastResponse struct {
	*fasthttp.RequestCtx
}

func (f *FastResponse) SetHeader(key string, value string) {
	f.RequestCtx.Response.Header.Set(key, value)
}

func (f *FastResponse) GetHeader(key string) string {
	b := f.RequestCtx.Response.Header.Peek(key)
	return string(b)
}

func (f *FastResponse) RequestURI() string {
	return f.RequestURI()
}

func (f *FastResponse) RequestPath() string {
	return string(f.Request.URI().Path())
}

func (f *FastResponse) RequestMethod() string {
	return string(f.Method())
}

func (f *FastResponse) RemoteAddr() string {
	return f.RequestCtx.RemoteAddr().String()
}

func (f *FastResponse) QueryValue(key string) string {
	return string(f.RequestCtx.QueryArgs().Peek(key))
}

func (f *FastResponse) QueryValues(key string) []string {
	b := f.RequestCtx.QueryArgs().PeekMulti(key)
	r := make([]string, len(b))
	for k, v := range b {
		r[k] = string(v)
	}
	return r
}