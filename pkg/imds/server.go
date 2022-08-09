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
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/purpleclay/imds-mock/pkg/imds/patch"
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

//go:embed on-demand.json
var onDemandResponse []byte

// Options provides a set of options for configuring the behaviour
// of the IMDS mock
type Options struct {
	// AutoStart determines whether the IMDS mock immediately starts
	// after initialisation
	AutoStart bool

	// InstanceTags ...
	InstanceTags bool

	// Port controls the port that is used by the IMDS mock. By default
	// it will use port 1338
	Port int

	// Pretty controls if the JSON outputted by any instance category
	// is pretty printed. By default all JSON will be compacted
	Pretty bool
}

// DefaultOptions defines the default set of options that will be applied
// to the IMDS mock upon startup
var DefaultOptions = Options{
	AutoStart:    true,
	InstanceTags: false,
	Port:         1338,
	Pretty:       false,
}

// Used as a hashset for quick lookups. Any matched path will just return its value
// and not be used to perform a key lookup
var reservedPaths = map[string]struct{}{
	"iam.info":                 {},
	"iam.security-credentials": {},
}

// Serve configures the IMDS mock using default options to handle HTTP requests
// in the exact same way as the IMDS service accessible from any EC2 instance
func Serve() (*gin.Engine, error) {
	return ServeWith(DefaultOptions)
}

// ServeWith configures the IMDS mock based on the incoming options to handle HTTP requests
// in the exact same way as the IMDS service accessible from any EC2 instance
func ServeWith(opts Options) (*gin.Engine, error) {
	r := gin.Default()

	// Inject middleware based on the provided options
	if opts.Pretty {
		r.Use(PrettyJSON())
	} else {
		r.Use(CompactJSON())
	}

	// see: https://pkg.go.dev/github.com/gin-gonic/gin#readme-don-t-trust-all-proxies
	r.SetTrustedProxies(nil)

	// Dynamically build the JSON response message used by the mock by using JSON
	// patching strategies defined by the input options
	mockResponse, err := patchResponseJSON(onDemandResponse, opts)
	if err != nil {
		return nil, err
	}

	r.GET("/latest/meta-data", func(c *gin.Context) {
		c.String(http.StatusOK, keys(mockResponse, ""))
	})

	r.GET("/latest/meta-data/*category", func(c *gin.Context) {
		categoryPath := c.Param("category")
		if categoryPath == "/" {
			// Exact same behaviour as /latest/meta-data
			c.String(http.StatusOK, keys(mockResponse, ""))
			return
		}

		// Convert param into gjson path query
		categoryPath = strings.TrimSuffix(categoryPath, "/")
		categoryPath = strings.ReplaceAll(categoryPath, "/", ".")[1:]

		res := gjson.GetBytes(mockResponse, categoryPath)

		if !res.Exists() {
			c.Writer.Header().Add("Content-Type", "text/html")
			c.String(http.StatusNotFound, notFound)
			return
		}

		c.Writer.Header().Add("Content-Type", "text/plain")

		// If the path returns a JSON object, then return a set of keys
		if res.IsObject() && notReservedPath(categoryPath) {
			c.String(http.StatusOK, keys(mockResponse, categoryPath))
		} else {
			c.String(http.StatusOK, res.String())
		}
	})

	if opts.AutoStart {
		err = r.Run(":" + strconv.Itoa(opts.Port))
	}

	return r, err
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

func patchResponseJSON(in []byte, opts Options) ([]byte, error) {
	// Build up a list of ordered patching strategies and apply
	patches := make([]patch.JSONPatcher, 0)

	if opts.InstanceTags {
		patches = append(patches, patch.InstanceTag{
			Tags: map[string]string{
				"Name": "ec2-imds-mock",
			},
		})
	}

	var err error
	for _, p := range patches {
		in, err = p.Patch(in)
		if err != nil {
			return in, err
		}
	}

	return in, nil
}
