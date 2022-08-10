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
	"bytes"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Capture any write to the response with an in memory buffer
type cacheWriter struct {
	gin.ResponseWriter
	body bytes.Buffer
}

func (w *cacheWriter) Write(data []byte) (n int, err error) {
	w.body.Write(data)
	return w.ResponseWriter.Write(data)
}

// Cache provides middleware caching any request to the IMDS mock using
// an in memory map. The IMDS response is cached using a lookup query based
// on the IMDS instance category path
func Cache() gin.HandlerFunc {
	cache := map[string]string{}

	return func(c *gin.Context) {
		if res, hit := cache[c.Request.URL.Path]; hit {
			c.Writer.Header().Add("Content-Type", "text/plain")
			c.String(http.StatusOK, res)

			c.Abort()
			return
		}

		c.Writer = &cacheWriter{ResponseWriter: c.Writer}
		c.Next()

		// Grab the contents of the cache
		cache[c.Request.URL.Path] = c.Writer.(*cacheWriter).body.String()
	}
}
