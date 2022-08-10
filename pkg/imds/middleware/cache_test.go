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

package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/purpleclay/imds-mock/pkg/imds/middleware"
	"github.com/stretchr/testify/assert"
)

func TestCache_Miss(t *testing.T) {
	count := 0

	r := gin.Default()
	r.GET("/cache", middleware.Cache(), func(c *gin.Context) {
		count++
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/cache", http.NoBody)

	r.ServeHTTP(w, req)

	assert.Equal(t, 1, count)
}

func TestCache_Hit(t *testing.T) {
	count := 0

	r := gin.Default()
	r.GET("/cache", middleware.Cache(), func(c *gin.Context) {
		count++
		c.String(http.StatusOK, "ok")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/cache", http.NoBody)

	r.ServeHTTP(w, req)
	// Its important to reset the buffer here
	w.Body.Reset()
	r.ServeHTTP(w, req)

	assert.Equal(t, 1, count)
	assert.Equal(t, "ok", w.Body.String())
}
