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

//go:embed basic.json
var basicResponse []byte

/*
TODO: if unknown path used

<?xml version="1.0" encoding="iso-8859-1"?>
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN"
	"http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html xmlns="http://www.w3.org/1999/xhtml" xml:lang="en" lang="en">
 <head>
  <title>404 - Not Found</title>
 </head>
 <body>
  <h1>404 - Not Found</h1>
 </body>
</html>
*/

// StartAndListen ...
func StartAndListen() {
	r := gin.Default()

	// see: https://pkg.go.dev/github.com/gin-gonic/gin#readme-don-t-trust-all-proxies
	r.SetTrustedProxies(nil)

	// TODO: return XML if path doesn't exist within the JSON

	r.GET("/latest/meta-data", func(c *gin.Context) {
		keys := gjson.GetBytes(basicResponse, "@keys").Array()

		categories := make([]string, 0, len(keys))
		for _, key := range keys {
			categories = append(categories, key.String())
		}

		c.String(http.StatusOK, strings.Join(categories, "\n"))
	})

	r.GET("/latest/meta-data/*category", func(c *gin.Context) {
		// Convert param into gjson path query
		path := strings.ReplaceAll(c.Param("category"), "/", ".")[1:]
		res := gjson.GetBytes(basicResponse, path)

		if res.IsObject() {
			// TODO: duplicate code
			keys := gjson.GetBytes(basicResponse, path+".@keys").Array()

			categories := make([]string, 0, len(keys))
			for _, key := range keys {
				categories = append(categories, key.String())
			}

			c.String(http.StatusOK, strings.Join(categories, "\n"))
		} else {
			c.String(http.StatusOK, res.String())
		}
	})

	r.Run()
}
