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

package token_test

import (
	"testing"
	"time"

	"github.com/purpleclay/imds-mock/pkg/imds/token"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewV2(t *testing.T) {
	token := token.NewV2(10)

	assert.WithinDuration(t, token.Expire, time.Now(), 10*time.Second)
}

func TestExpired(t *testing.T) {
	tests := []struct {
		name    string
		ttl     int
		expired bool
	}{
		{
			name:    "WithFutureTTL",
			ttl:     10,
			expired: false,
		},
		{
			name:    "WithExpiredTTL",
			ttl:     -1,
			expired: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token := token.V2{
				Expire: time.Now().Add(time.Duration(tt.ttl) * time.Second),
			}

			require.Equal(t, tt.expired, token.Expired())
		})
	}
}
