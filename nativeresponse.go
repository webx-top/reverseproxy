package reverseproxy

import (
	"io"
	"net/http"
)

var _ Context = &NativeResponse{}

type NativeResponse struct {
	http.ResponseWriter
	*http.Request
}

func (n *NativeResponse) SetBody(body []byte) {
	n.ResponseWriter.Write(body)
}

func (n *NativeResponse) SetStatusCode(code int) {
	n.ResponseWriter.WriteHeader(code)
}

func (n *NativeResponse) Redirect(url string, code int) {
	http.Redirect(n.ResponseWriter, n.Request, url, code)
}

func (n *NativeResponse) SetHeader(key string, value string) {
	n.ResponseWriter.Header().Set(key, value)
}

func (n *NativeResponse) GetHeader(key string) string {
	return n.ResponseWriter.Header().Get(key)
}

func (n *NativeResponse) RequestURI() string {
	return n.Request.URL.RequestURI()
}

func (n *NativeResponse) RequestPath() string {
	return n.Request.URL.Path
}

func (n *NativeResponse) RequestMethod() string {
	return n.Request.Method
}

func (n *NativeResponse) RemoteAddr() string {
	return n.Request.RemoteAddr
}

func (n *NativeResponse) QueryValue(key string) string {
	return n.Request.URL.Query().Get(key)
}

func (n *NativeResponse) QueryValues(key string) []string {
	q := n.Request.URL.Query()
	if v, ok := q[key]; ok {
		return v
	}
	return []string{}
}

func (n *NativeResponse) RespWriter() io.Writer {
	return n.ResponseWriter
}

func (n *NativeResponse) RequestHost() string {
	return n.Request.Host
}
