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
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/purpleclay/imds-mock/pkg/imds/token"
)

const (
	unauthorised = `<?xml version="1.0" encoding="iso-8859-1"?>
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN"
	"http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html xmlns="http://www.w3.org/1999/xhtml" xml:lang="en" lang="en">
  <head>
    <title>401 - Unauthorized</title>
  </head>
  <body>
    <h1>401 - Unauthorized</h1>
  </body>
</html>`

	// V2TokenHeader defines the HTTP header used by the IMDS service for extracting a
	// a session token from an HTTP request
	V2TokenHeader = "X-aws-ec2-metadata-token"
)

// StrictV2 provides middleware that explicitly enables IMDSv2 authorisation
// through the use of session tokens. Any requests without a valid session
// token are immediately rejected. A session token is presented to this middleware
// through the use of an HTTP header:
//
//	X-aws-ec2-metadata-token: TOKEN
func StrictV2() gin.HandlerFunc {
	return func(c *gin.Context) {
		tkn := c.Request.Header.Get(V2TokenHeader)
		if tkn == "" {
			abortUnauthorised(c)
			return
		}

		decoded, err := base64.StdEncoding.DecodeString(tkn)
		if err != nil {
			abortUnauthorised(c)
			return
		}

		var st token.SessionToken
		if err := json.Unmarshal(decoded, &st); err != nil {
			abortUnauthorised(c)
			return
		}

		if st.Expired() {
			abortUnauthorised(c)
			return
		}

		// Safe to proceed
		c.Next()
	}
}

func abortUnauthorised(c *gin.Context) {
	c.Writer.Header().Add("Content-Type", "text/html")
	c.String(http.StatusUnauthorized, unauthorised)
	c.Abort()
}
