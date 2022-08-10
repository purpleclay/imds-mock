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
	"bytes"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"testing"
	"text/template"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/purpleclay/imds-mock/pkg/imds/middleware"
	"github.com/stretchr/testify/assert"
)

const (
	// Token was generated (base64) from the following JSON: {"expire":"2022-08-10T06:49:05.415516+01:00"}
	expiredToken = "eyJleHBpcmUiOiIyMDIyLTA4LTEwVDA2OjQ5OjA1LjQxNTUxNiswMTowMCJ9"

	// Token was generated (base64) from the following text: hello, world
	invalidSessionToken = "aGVsbG8sIHdvcmxk"

	// Token is not base64 encoded
	notBase64Token = "hello, world"

	// Template for generating a valid token
	tokenTemplate = `{"expire":"{{ . }}"}`
)

var instanceTagPatch = template.Must(template.New("TokenTemplate").Parse(tokenTemplate))

func TestStrictV2(t *testing.T) {
	tests := []struct {
		name  string
		token string
	}{
		{
			name:  "MissingToken",
			token: "",
		},
		{
			name:  "ExpiredToken",
			token: expiredToken,
		},
		{
			name:  "UnsupportedToken",
			token: invalidSessionToken,
		},
		{
			name:  "CannotDecodeToken",
			token: notBase64Token,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := v2Router(t)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/", http.NoBody)
			if tt.token != "" {
				req.Header.Add(middleware.V2TokenHeader, tt.token)
			}

			r.ServeHTTP(w, req)

			assert.Equal(t, http.StatusUnauthorized, w.Code)
			assert.Equal(t, "text/html", w.Result().Header["Content-Type"][0])
			assert.Equal(t, `<?xml version="1.0" encoding="iso-8859-1"?>
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN"
	"http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html xmlns="http://www.w3.org/1999/xhtml" xml:lang="en" lang="en">
  <head>
    <title>401 - Unauthorized</title>
  </head>
  <body>
    <h1>401 - Unauthorized</h1>
  </body>
</html>`, w.Body.String())
		})
	}
}

func TestStrictV2_ValidToken(t *testing.T) {
	var buf bytes.Buffer
	instanceTagPatch.Execute(&buf, time.Now().Add(2*time.Second).Format(time.RFC3339))

	r := v2Router(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/", http.NoBody)
	req.Header.Add(middleware.V2TokenHeader, base64.StdEncoding.EncodeToString(buf.Bytes()))

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "ok", w.Body.String())
}

func v2Router(t *testing.T) *gin.Engine {
	t.Helper()
	r := gin.Default()
	r.GET("/", middleware.StrictV2(), func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	return r
}
