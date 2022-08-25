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
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/purpleclay/imds-mock/pkg/imds/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func zapObservedRouter(t *testing.T) (*gin.Engine, *observer.ObservedLogs) {
	t.Helper()
	r := gin.New()

	core, obs := observer.New(zap.InfoLevel)
	logger := zap.New(core)

	r.Use(middleware.ZapRecovery(logger), middleware.ZapLogger(logger))

	return r, obs
}

func TestZapLogger(t *testing.T) {
	r, observedLogs := zapObservedRouter(t)

	r.GET("/testing", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/testing?a=b", http.NoBody)
	r.ServeHTTP(w, req)

	require.Len(t, observedLogs.All(), 1)
	logEntry := observedLogs.All()[0]

	require.Len(t, logEntry.Context, 8)
	assert.Equal(t, zapcore.InfoLevel, logEntry.Level)
	assert.Equal(t, "GET /testing", logEntry.Message)
	assert.Equal(t, int64(http.StatusOK), logEntry.ContextMap()["status"])
	assert.Equal(t, "GET", logEntry.ContextMap()["method"])
	assert.Equal(t, "/testing", logEntry.ContextMap()["path"])
	assert.Equal(t, "a=b", logEntry.ContextMap()["query"])
	assert.Equal(t, "", logEntry.ContextMap()["ip"])
	assert.Equal(t, "", logEntry.ContextMap()["user-agent"])

	_, err := time.Parse(time.RFC3339, logEntry.ContextMap()["time"].(string))
	assert.NoError(t, err)

	dur, err := time.ParseDuration(logEntry.ContextMap()["latency"].(string))
	assert.NoError(t, err)
	assert.True(t, dur > 0)
}

func TestZapLoggerError(t *testing.T) {
	r, observedLogs := zapObservedRouter(t)

	r.GET("/testing", func(c *gin.Context) {
		c.AbortWithError(http.StatusBadRequest, errors.New("an error has occurred"))
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/testing?a=b", http.NoBody)
	r.ServeHTTP(w, req)

	require.Len(t, observedLogs.All(), 1)
	logEntry := observedLogs.All()[0]

	require.Len(t, logEntry.Context, 9)
	assert.Equal(t, zapcore.ErrorLevel, logEntry.Level)
	assert.Equal(t, "GET /testing", logEntry.Message)
	assert.Equal(t, "an error has occurred", logEntry.ContextMap()["error"])
	assert.Equal(t, int64(http.StatusBadRequest), logEntry.ContextMap()["status"])
	assert.Equal(t, "GET", logEntry.ContextMap()["method"])
	assert.Equal(t, "/testing", logEntry.ContextMap()["path"])
	assert.Equal(t, "a=b", logEntry.ContextMap()["query"])
	assert.Equal(t, "", logEntry.ContextMap()["ip"])
	assert.Equal(t, "", logEntry.ContextMap()["user-agent"])

	_, err := time.Parse(time.RFC3339, logEntry.ContextMap()["time"].(string))
	assert.NoError(t, err)

	dur, err := time.ParseDuration(logEntry.ContextMap()["latency"].(string))
	assert.NoError(t, err)
	assert.True(t, dur > 0)
}

func TestZapRecovery(t *testing.T) {
	r, observedLogs := zapObservedRouter(t)

	r.GET("/testing", func(c *gin.Context) {
		panic("panic")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/testing?a=b", http.NoBody)
	r.ServeHTTP(w, req)

	require.Len(t, observedLogs.All(), 1)
	logEntry := observedLogs.All()[0]

	assert.Equal(t, zapcore.ErrorLevel, logEntry.Level)
	assert.Equal(t, "panic recovery", logEntry.Message)
	assert.Equal(t, "panic", logEntry.ContextMap()["error"])
	assert.Equal(t, "GET", logEntry.ContextMap()["method"])
	assert.Equal(t, "/testing", logEntry.ContextMap()["path"])
	assert.Equal(t, "a=b", logEntry.ContextMap()["query"])

	_, err := time.Parse(time.RFC3339, logEntry.ContextMap()["time"].(string))
	assert.NoError(t, err)
}

func TestZapRecoveryBrokenPipe(t *testing.T) {
	r, observedLogs := zapObservedRouter(t)

	r.GET("/testing", func(c *gin.Context) {
		// Horrible nesting of errors
		panic(&net.OpError{Err: &os.SyscallError{Err: errors.New("broken pipe")}})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/testing?a=b", http.NoBody)
	r.ServeHTTP(w, req)

	require.Len(t, observedLogs.All(), 1)
	logEntry := observedLogs.All()[0]

	assert.Equal(t, zapcore.ErrorLevel, logEntry.Level)
	assert.Equal(t, "GET /testing", logEntry.Message)
	assert.Contains(t, logEntry.ContextMap()["error"], "broken pipe")
	assert.Equal(t, "GET", logEntry.ContextMap()["method"])
	assert.Equal(t, "/testing", logEntry.ContextMap()["path"])
	assert.Equal(t, "a=b", logEntry.ContextMap()["query"])
}
