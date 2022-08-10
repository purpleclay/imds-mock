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
	"github.com/gin-gonic/gin"
	"github.com/tidwall/pretty"
)

// JSONFormatter defines an interface for formatting JSON text
type JSONFormatter interface {
	Format(in []byte) []byte
}

// A crude wrapper around a gin.ResponseWriter for supporting the custom
// injection of a JSON formatter
type jsonRewriter struct {
	gin.ResponseWriter
	Formatter JSONFormatter
}

func (w *jsonRewriter) Write(data []byte) (n int, err error) {
	return w.ResponseWriter.Write(w.Formatter.Format(data))
}

type compactJSONFormatter struct{}

func (f compactJSONFormatter) Format(in []byte) []byte {
	if len(in) == 0 {
		return in
	}

	if in[0] != '{' {
		return in
	}

	return pretty.Ugly(in)
}

// CompactJSON provides middleware capable of rewriting all JSON responses
// into a compact format
func CompactJSON() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer = &jsonRewriter{
			ResponseWriter: c.Writer,
			Formatter:      compactJSONFormatter{},
		}

		c.Next()
	}
}

var prettyOptions = &pretty.Options{
	Width:    pretty.DefaultOptions.Width,
	Prefix:   pretty.DefaultOptions.Prefix,
	Indent:   pretty.DefaultOptions.Indent,
	SortKeys: true,
}

type prettyJSONFormatter struct{}

func (f prettyJSONFormatter) Format(in []byte) []byte {
	if len(in) == 0 {
		return in
	}

	if in[0] != '{' {
		return in
	}

	return pretty.PrettyOptions(in, prettyOptions)
}

// PrettyJSON provides middleware capable of rewriting all JSON responses
// into a pretty format
func PrettyJSON() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer = &jsonRewriter{
			ResponseWriter: c.Writer,
			Formatter:      prettyJSONFormatter{},
		}

		c.Next()
	}
}
