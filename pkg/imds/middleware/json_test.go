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

func TestCompactJSON(t *testing.T) {
	r := gin.Default()
	r.GET("/", middleware.CompactJSON(), func(c *gin.Context) {
		c.String(http.StatusOK, `{
	"a": "1",
	"b": "2"
}`)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/", http.NoBody)

	r.ServeHTTP(w, req)

	assert.Equal(t, `{"a":"1","b":"2"}`, w.Body.String())
}

func TestCompactJSON_IgnoreNonJSON(t *testing.T) {
	r := gin.Default()
	r.GET("/", middleware.CompactJSON(), func(c *gin.Context) {
		c.String(http.StatusOK, `no json`)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/", http.NoBody)

	r.ServeHTTP(w, req)

	assert.Equal(t, "no json", w.Body.String())
}

func TestCompactJSON_IgnoreEmptyString(t *testing.T) {
	r := gin.Default()
	r.GET("/", middleware.CompactJSON(), func(c *gin.Context) {
		c.String(http.StatusOK, "")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/", http.NoBody)

	r.ServeHTTP(w, req)

	assert.Empty(t, w.Body.String())
}

func TestPrettyJSON(t *testing.T) {
	r := gin.Default()
	r.GET("/", middleware.PrettyJSON(), func(c *gin.Context) {
		c.String(http.StatusOK, `{"c":"3","b":"2","a":"1"}`)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/", http.NoBody)

	r.ServeHTTP(w, req)

	assert.Equal(t, `{
  "a": "1",
  "b": "2",
  "c": "3"
}
`, w.Body.String())
}

func TestPrettyJSON_IgnoreNonJSON(t *testing.T) {
	r := gin.Default()
	r.GET("/", middleware.PrettyJSON(), func(c *gin.Context) {
		c.String(http.StatusOK, "not json")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/", http.NoBody)

	r.ServeHTTP(w, req)

	assert.Equal(t, "not json", w.Body.String())
}

func TestPrettyJSON_IgnoreEmptyString(t *testing.T) {
	r := gin.Default()
	r.GET("/", middleware.PrettyJSON(), func(c *gin.Context) {
		c.String(http.StatusOK, "")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/", http.NoBody)

	r.ServeHTTP(w, req)

	assert.Empty(t, w.Body.String())
}
