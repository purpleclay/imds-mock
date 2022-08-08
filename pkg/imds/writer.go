package imds

import (
	"bufio"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
)

// JSONFormatter defines an interface for formatting JSON text
type JSONFormatter interface {
	Format(in []byte) []byte
	FormatString(in string) string
}

// A crude wrapper around a gin.ResponseWriter for supporting the custom
// injection of a JSON formatter
type jsonRewriter struct {
	ResponseWriter gin.ResponseWriter
	Formatter      JSONFormatter
}

func (w *jsonRewriter) Header() http.Header  { return w.ResponseWriter.Header() }
func (w *jsonRewriter) WriteHeader(code int) { w.ResponseWriter.WriteHeader(code) }
func (w *jsonRewriter) WriteHeaderNow()      { w.ResponseWriter.WriteHeaderNow() }

func (w *jsonRewriter) Write(data []byte) (n int, err error) {
	return w.ResponseWriter.Write(w.Formatter.Format(data))
}

func (w *jsonRewriter) WriteString(s string) (n int, err error) {
	return w.ResponseWriter.WriteString(w.Formatter.FormatString(s))
}

func (w *jsonRewriter) Status() int   { return w.ResponseWriter.Status() }
func (w *jsonRewriter) Size() int     { return w.ResponseWriter.Size() }
func (w *jsonRewriter) Written() bool { return w.ResponseWriter.Written() }

func (w *jsonRewriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return w.ResponseWriter.Hijack()
}

func (w *jsonRewriter) CloseNotify() <-chan bool     { return w.ResponseWriter.CloseNotify() }
func (w *jsonRewriter) Flush()                       { w.ResponseWriter.Flush() }
func (w *jsonRewriter) Pusher() (pusher http.Pusher) { return w.ResponseWriter.Pusher() }
