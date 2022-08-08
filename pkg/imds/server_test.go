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

package imds_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/purpleclay/imds-mock/pkg/imds"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/pretty"
)

var testOptions = imds.Options{
	AutoStart: false,
	Pretty:    false,
}

func TestMain(m *testing.M) {
	// Ensure Gin runs in test mode
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

func TestLatestMetadataKeys(t *testing.T) {
	r, _ := imds.ServeWith(testOptions)

	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name: "ForRoot",
			path: "/latest/meta-data",
			expected: `ami-id
ami-launch-index
ami-manifest-path
block-device-mapping/
events/
hostname
iam/
instance-action
instance-id
instance-life-cycle
instance-type
local-hostname
local-ipv4
mac
network/
placement/
profile
public-keys/
reservation-id
security-groups
services/`,
		},
		{
			name: "ForRootTrailingSlash",
			path: "/latest/meta-data/",
			expected: `ami-id
ami-launch-index
ami-manifest-path
block-device-mapping/
events/
hostname
iam/
instance-action
instance-id
instance-life-cycle
instance-type
local-hostname
local-ipv4
mac
network/
placement/
profile
public-keys/
reservation-id
security-groups
services/`,
		},
		{
			name: "ForSubCategory",
			path: "/latest/meta-data/iam",
			expected: `info/
security-credentials/`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, tt.path, http.NoBody)

			r.ServeHTTP(w, req)

			require.Equal(t, http.StatusOK, w.Code)
			require.Equal(t, tt.expected, w.Body.String())
		})
	}
}

func TestCategoryValueIsString(t *testing.T) {
	r, _ := imds.ServeWith(testOptions)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/latest/meta-data/ami-id", http.NoBody)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "text/plain", w.Result().Header["Content-Type"][0])
	assert.Equal(t, "ami-0e34bbddc66def5ac", w.Body.String())
}

func TestCategoryValueIsCompactJSON(t *testing.T) {
	r, _ := imds.ServeWith(testOptions)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/latest/meta-data/iam/info", http.NoBody)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "text/plain", w.Result().Header["Content-Type"][0])
	assert.Equal(t, string(pretty.Ugly([]byte(`{
	"Code": "Success",
	"LastUpdated": "2022-08-08T04:25:36Z",
	"InstanceProfileArn": "arn:aws:iam::112233445566:instance-profile/ssm-access",
	"InstanceProfileId": "AIPAYUKXDENX4ZNCZWHF6"
}`))), w.Body.String())
}

func TestCategoryValueIsPrettyJSON(t *testing.T) {
	opts := testOptions
	opts.Pretty = true

	r, _ := imds.ServeWith(opts)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/latest/meta-data/iam/info", http.NoBody)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "text/plain", w.Result().Header["Content-Type"][0])
	assert.Equal(t, `{
  "Code": "Success",
  "InstanceProfileArn": "arn:aws:iam::112233445566:instance-profile/ssm-access",
  "InstanceProfileId": "AIPAYUKXDENX4ZNCZWHF6",
  "LastUpdated": "2022-08-08T04:25:36Z"
}
`, w.Body.String())
}

func TestCategoryPathDoesNotExist(t *testing.T) {
	r, _ := imds.ServeWith(testOptions)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/latest/meta-data/unknown", http.NoBody)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, "text/html", w.Result().Header["Content-Type"][0])
	assert.Equal(t, `<?xml version="1.0" encoding="iso-8859-1"?>
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN"
	"http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html xmlns="http://www.w3.org/1999/xhtml" xml:lang="en" lang="en">
 <head>
  <title>404 - Not Found</title>
 </head>
 <body>
  <h1>404 - Not Found</h1>
 </body>
</html>`, w.Body.String())
}
