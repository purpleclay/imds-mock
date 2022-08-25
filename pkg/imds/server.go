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
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/purpleclay/imds-mock/pkg/imds/cache"
	"github.com/purpleclay/imds-mock/pkg/imds/event"
	"github.com/purpleclay/imds-mock/pkg/imds/middleware"
	"github.com/purpleclay/imds-mock/pkg/imds/patch"
	"github.com/purpleclay/imds-mock/pkg/imds/token"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
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

	badRequest = `<?xml version="1.0" encoding="iso-8859-1"?>
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN"
	"http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html xmlns="http://www.w3.org/1999/xhtml" xml:lang="en" lang="en">
 <head>
  <title>400 - Bad Request</title>
 </head>
 <body>
  <h1>400 - Bad Request</h1>
 </body>
</html>`

	// V2TokenTTLHeader defines the HTTP header used by the IMDS service for
	// generating a new IMDS session token
	V2TokenTTLHeader = "X-aws-ec2-metadata-token-ttl-seconds"
)

//go:embed on-demand.json
var onDemandInstance []byte

// Really crude attempt to protect a byte array from concurrency issues during
// event driven patches
type patchedJSON struct {
	data []byte
	mu   sync.RWMutex
}

func (p *patchedJSON) Bytes() []byte {
	p.mu.RLock()
	copy := p.data
	p.mu.RUnlock()

	return copy
}

func (p *patchedJSON) Patch(patcher patch.JSONPatcher) error {
	p.mu.Lock()
	var err error
	p.data, err = patcher.Patch(p.data)
	p.mu.Unlock()

	return err
}

// Options provides a set of options for configuring the behaviour
// of the IMDS mock
type Options struct {
	// AutoStart determines whether the IMDS mock immediately starts
	// after initialisation
	AutoStart bool

	// ExcludeInstanceTags controls if the IMDS mock excludes instance
	// tags from its supported list of metadata categories
	ExcludeInstanceTags bool

	// IMDSv2 enables exclusive V2 support only. All requests must contain
	// a valid metadata token, otherwise they will be rejected. By default
	// the mock will run with both V1 and V2 support
	IMDSv2 bool

	// InstanceTags contains a map of key value pairs that are to be
	// exposed as instance tags through the IMDS mock
	InstanceTags map[string]string

	// Port controls the port that is used by the IMDS mock. By default
	// it will use port 1338
	Port int

	// Pretty controls if the JSON outputted by any instance category
	// is pretty printed. By default all JSON will be compacted
	Pretty bool

	// Spot enables the simulation of a spot instance and interruption notice
	// through the IMDS mock. By default this will set to false and an on-demand
	// instance will be simulated
	Spot bool

	// SpotAction defines the type of spot interruption event to be raised when
	// simulating a spot instance. By default the spot interruption will be
	// immediate, but can be delayed by a pre-configured interval
	SpotAction SpotActionEvent
}

// SpotActionEvent defines a spot interruption event
type SpotActionEvent struct {
	Action   patch.SpotInstanceAction
	Duration time.Duration
}

// DefaultOptions defines the default set of options that will be applied
// to the IMDS mock upon startup
var DefaultOptions = Options{
	AutoStart:           true,
	ExcludeInstanceTags: false,
	IMDSv2:              false,
	InstanceTags: map[string]string{
		"Name": "imds-mock-ec2",
	},
	Port:   1338,
	Pretty: false,
	Spot:   false,
	SpotAction: SpotActionEvent{
		Action:   patch.TerminateSpotInstanceAction,
		Duration: 0 * time.Second,
	},
}

// Used as a hashset for quick lookups. Any matched path will just return its value
// and not be used to perform a key lookup
var reservedPaths = map[string]struct{}{
	"iam.info":                         {},
	"iam.security-credentials":         {},
	"spot.instance-action":             {},
	"events.recommendations.rebalance": {},
}

// Serve configures the IMDS mock using default options to handle HTTP requests
// in the exact same way as the IMDS service accessible from any EC2 instance
func Serve() (*gin.Engine, error) {
	return ServeWith(DefaultOptions)
}

// ServeWith configures the IMDS mock based on the incoming options to handle HTTP requests
// in the exact same way as the IMDS service accessible from any EC2 instance
func ServeWith(opts Options) (*gin.Engine, error) {
	r := gin.New()
	injectMiddleware(r, opts)

	// see: https://pkg.go.dev/github.com/gin-gonic/gin#readme-don-t-trust-all-proxies
	r.SetTrustedProxies(nil)

	// Locally managed cache
	memcache := cache.New()

	// Manage the patching of the underlying JSON that is served by the IMDS mock
	mockResponse := patchedJSON{data: onDemandInstance}

	if !opts.ExcludeInstanceTags {
		if err := mockResponse.Patch(patch.InstanceTag{Tags: opts.InstanceTags}); err != nil {
			return nil, err
		}
	}

	// Event based patching of spot instance
	if opts.Spot {
		if opts.SpotAction.Duration > 0 {
			event.Once(opts.SpotAction.Duration, func() {
				if patchErr := mockResponse.Patch(patch.Spot{InstanceAction: opts.SpotAction.Action}); patchErr != nil {
					return
				}

				// Invalidate the cache to ensure the mock returns the new spot instance categories
				memcache.Remove("/latest/meta-data", "/latest/meta-data/")
			})
		} else {
			if err := mockResponse.Patch(patch.Spot{InstanceAction: opts.SpotAction.Action}); err != nil {
				return nil, err
			}
		}
	}

	r.GET("/latest/meta-data", middleware.Cache(memcache), func(c *gin.Context) {
		c.String(http.StatusOK, keys(mockResponse.Bytes(), ""))
	})

	r.GET("/latest/meta-data/*category", middleware.Cache(memcache), func(c *gin.Context) {
		categoryPath := c.Param("category")
		if categoryPath == "/" {
			// Exact same behaviour as /latest/meta-data
			c.String(http.StatusOK, keys(mockResponse.Bytes(), ""))
			return
		}

		// Convert param into gjson path query
		categoryPath = strings.TrimSuffix(categoryPath, "/")
		categoryPath = strings.ReplaceAll(categoryPath, "/", ".")[1:]

		res := gjson.GetBytes(mockResponse.Bytes(), categoryPath)

		if !res.Exists() {
			c.Writer.Header().Add("Content-Type", "text/html")
			c.String(http.StatusNotFound, notFound)
			return
		}

		c.Writer.Header().Add("Content-Type", "text/plain")

		// If the path returns a JSON object, then return a set of keys
		if res.IsObject() && notReservedPath(categoryPath) {
			c.String(http.StatusOK, keys(mockResponse.Bytes(), categoryPath))
		} else {
			c.String(http.StatusOK, res.String())
		}
	})

	r.PUT("/latest/api/token", func(c *gin.Context) {
		ttl, err := strconv.Atoi(c.Request.Header.Get(V2TokenTTLHeader))

		if err == nil && (ttl > 0 && ttl <= token.MaxTTLInSeconds) {
			tkn := token.NewV2(ttl)
			out, _ := json.Marshal(&tkn)

			c.Writer.Header().Add("Content-Type", "text/plain")
			c.String(http.StatusOK, base64.StdEncoding.EncodeToString(out))
			return
		}

		c.Writer.Header().Add("Content-Type", "text/html")
		c.String(http.StatusBadRequest, badRequest)
	})

	var err error
	if opts.AutoStart {
		err = r.Run(":" + strconv.Itoa(opts.Port))
	}

	return r, err
}

func injectMiddleware(r *gin.Engine, opts Options) {
	logger, _ := zap.NewProduction()
	r.Use(middleware.ZapLogger(logger), middleware.ZapRecovery(logger))

	if opts.IMDSv2 {
		r.Use(middleware.StrictV2())
	} else {
		r.Use(middleware.V1OptionalV2())
	}

	if opts.Pretty {
		r.Use(middleware.PrettyJSON())
	} else {
		r.Use(middleware.CompactJSON())
	}
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
