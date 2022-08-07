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

//go:embed basic.json
var basicResponse []byte

// StartAndListen ...
func StartAndListen() {
	r := gin.Default()

	// see: https://pkg.go.dev/github.com/gin-gonic/gin#readme-don-t-trust-all-proxies
	r.SetTrustedProxies(nil)

	r.GET("/latest/meta-data", func(c *gin.Context) {
		c.String(http.StatusOK, keys(basicResponse, ""))
	})

	r.GET("/latest/meta-data/*category", func(c *gin.Context) {
		categoryPath := c.Param("category")
		if categoryPath == "/" {
			// Exact same behaviour as /latest/meta-data
			c.String(http.StatusOK, keys(basicResponse, ""))
			return
		}

		// Convert param into gjson path query
		categoryPath = strings.ReplaceAll(categoryPath, "/", ".")[1:]
		res := gjson.GetBytes(basicResponse, categoryPath)

		if !res.Exists() {
			c.Writer.Header().Add("Content-Type", "text/html")
			c.String(http.StatusNotFound, notFound)
			return
		}

		// If the path returns a JSON object, then return a set of keys
		if res.IsObject() {
			c.String(http.StatusOK, keys(basicResponse, categoryPath))
		} else {
			c.String(http.StatusOK, res.String())
		}
	})

	r.Run()
}

func keys(json []byte, path string) string {
	query := "@keys"
	if path != "" {
		query = path + ".@keys"
	}

	keys := gjson.GetBytes(json, query).Array()

	categories := make([]string, 0, len(keys))
	for _, key := range keys {
		categories = append(categories, key.String())
	}

	return strings.Join(categories, "\n")
}
