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
	"net/textproto"

	"github.com/gin-gonic/gin"
)

// V1OptionalV2 provides middleware that enables both IMDSv1 and IMDSv2 support.
// While using V1 all requests will pass straight through without any authorisation
// checks. V2 checking will only be carried out on the presence of the HTTP header:
//
//	X-aws-ec2-metadata-token: TOKEN
func V1OptionalV2() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Headers are stored in a canonical format
		if _, exists := c.Request.Header[textproto.CanonicalMIMEHeaderKey(V2TokenHeader)]; exists {
			// Treat this exactly like a V2 request
			if !validV2Token(c.Request.Header.Get(V2TokenHeader)) {
				abortUnauthorised(c)
				return
			}
		}

		c.Next()
	}
}
