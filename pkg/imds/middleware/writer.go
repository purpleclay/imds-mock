/*
Copyright (c) 2022 Purple Clay

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package middleware

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
