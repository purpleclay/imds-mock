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
	"fmt"
	"net"
	"net/http"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ZapLogger provides middleware for logging all requests using Zap based structured logs
func ZapLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		end := time.Now()
		latency := end.Sub(start)

		if len(c.Errors) > 0 {
			for _, e := range c.Errors.Errors() {
				logger.Error(fmt.Sprintf("%s %s", c.Request.Method, c.Request.URL.Path),
					[]zapcore.Field{
						zap.Any("error", e),
						zap.String("ip", c.ClientIP()),
						zap.String("latency", latency.String()),
						zap.String("method", c.Request.Method),
						zap.String("path", c.Request.URL.Path),
						zap.String("query", c.Request.URL.RawQuery),
						zap.Int("status", c.Writer.Status()),
						zap.String("time", end.Format(time.RFC3339)),
						zap.String("user-agent", c.Request.UserAgent()),
					}...)
			}
		} else {
			// Replicates the existing Gin logger as close as possible
			logger.Info(fmt.Sprintf("%s %s", c.Request.Method, c.Request.URL.Path),
				[]zapcore.Field{
					zap.String("ip", c.ClientIP()),
					zap.String("latency", latency.String()),
					zap.String("method", c.Request.Method),
					zap.String("path", c.Request.URL.Path),
					zap.String("query", c.Request.URL.RawQuery),
					zap.Int("status", c.Writer.Status()),
					zap.String("time", end.Format(time.RFC3339)),
					zap.String("user-agent", c.Request.UserAgent()),
				}...)
		}
	}
}

// ZapRecovery provides middleware that recovers from any panics raised within requests
// and logs them using Zap based structured logs
func ZapRecovery(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Pinched from recovery.go#58 (gin)
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a condition that
				// warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				if brokenPipe {
					logger.Error(fmt.Sprintf("%s %s", c.Request.Method, c.Request.URL.Path),
						zap.String("time", time.Now().Format(time.RFC3339)),
						zap.Any("error", err),
						zap.String("method", c.Request.Method),
						zap.String("path", c.Request.URL.Path),
						zap.String("query", c.Request.URL.RawQuery),
					)
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				logger.Error("panic recovery",
					zap.String("time", time.Now().Format(time.RFC3339)),
					zap.Any("error", err),
					zap.String("method", c.Request.Method),
					zap.String("path", c.Request.URL.Path),
					zap.String("query", c.Request.URL.RawQuery),
					zap.String("stack", string(debug.Stack())),
				)

				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
