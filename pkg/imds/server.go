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

package imds

import (
	_ "embed"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"github.com/tidwall/pretty"
)

const (
	notFound = `<?xml version="1.0" encoding="iso-8859-1"?>
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN"
	"http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html xmlns="http://www.w3.org/1999/xhtml" xml:lang="en" lang="en">
 <head>
  <title>404 - Not Found</title>
 </head>
 <body>
  <h1>404 - Not Found</h1>
 </body>
</html>`
)

//go:embed on-demand.json
var onDemandResponse []byte

// Options ...
type Options struct {
	AutoStart bool
}

// DefaultOptions ...
var DefaultOptions = Options{
	AutoStart: true,
}

// Used as a hashset for quick lookups. Any matched path will just return its value
// and not be used to perform a key lookup
var reservedPaths = map[string]struct{}{
	"iam.info":                 {},
	"iam.security-credentials": {},
}

// Serve configures the IMDS mock based on the incoming options to handle HTTP requests
// in the exact same way as the IMDS service accessible from an EC2 instance
func Serve(with Options) *gin.Engine {
	r := gin.Default()

	// see: https://pkg.go.dev/github.com/gin-gonic/gin#readme-don-t-trust-all-proxies
	r.SetTrustedProxies(nil)

	r.GET("/latest/meta-data", func(c *gin.Context) {
		c.String(http.StatusOK, keys(onDemandResponse, ""))
	})

	r.GET("/latest/meta-data/*category", func(c *gin.Context) {
		categoryPath := c.Param("category")
		if categoryPath == "/" {
			// Exact same behaviour as /latest/meta-data
			c.String(http.StatusOK, keys(onDemandResponse, ""))
			return
		}

		// Convert param into gjson path query
		categoryPath = strings.TrimSuffix(categoryPath, "/")
		categoryPath = strings.ReplaceAll(categoryPath, "/", ".")[1:]

		res := gjson.GetBytes(onDemandResponse, categoryPath)

		if !res.Exists() {
			c.Writer.Header().Add("Content-Type", "text/html")
			c.String(http.StatusNotFound, notFound)
			return
		}

		c.Writer.Header().Add("Content-Type", "text/plain")

		// If the path returns a JSON object, then return a set of keys
		if res.IsObject() && notReservedPath(categoryPath) {
			c.String(http.StatusOK, keys(onDemandResponse, categoryPath))
		} else {
			out := res.String()
			if res.IsObject() {
				out = string(pretty.Ugly([]byte(out)))
			}

			c.String(http.StatusOK, out)
		}
	})

	if with.AutoStart {
		r.Run(":1338")
	}

	return r
}

func keys(json []byte, path string) string {
	// Scan the JSON document, retrieving all of the top-level fields as keys
	query := "@keys"
	if path != "" {
		query = path + ".@keys"
	}
	keys := gjson.GetBytes(json, query).Array()

	categories := make([]string, 0, len(keys))
	for _, key := range keys {
		k := key.String()

		// IMDS service returns a category with a trailing slash, if it is a parent category
		query = k
		if path != "" {
			query = path + "." + k
		}

		if gjson.GetBytes(json, query).IsObject() {
			k = k + "/"
		}

		categories = append(categories, k)
	}

	return strings.Join(categories, "\n")
}

func notReservedPath(path string) bool {
	_, ok := reservedPaths[path]
	return !ok
}
