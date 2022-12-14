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

package cmd

import (
	"testing"

	"github.com/purpleclay/imds-mock/pkg/imds/patch"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSpotActionFlagString(t *testing.T) {
	flag := spotActionFlag{}

	assert.Equal(t, "terminate=0s", flag.String())
}

func TestSpotFlagSet(t *testing.T) {
	flag := spotActionFlag{}
	err := flag.Set("terminate=1s")

	require.NoError(t, err)
	require.NotNil(t, flag.event)
	assert.Equal(t, patch.TerminateSpotInstanceAction, flag.event.Action)
	assert.Equal(t, "1s", flag.event.Duration.String())
}

func TestSpotFlagSetError(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		errMsg string
	}{
		{
			name:   "InvalidValue",
			input:  "invalid",
			errMsg: "invalid must be formatted as key=value e.g. terminate=2s",
		},
		{
			name:   "UnsupportedAction",
			input:  "restart=1s",
			errMsg: "restart is not a supported spot action expecting (terminate, stop or hibernate)",
		},
		{
			name:   "UnsupportedDuration",
			input:  "terminate=2days",
			errMsg: "2days is not a supported duration format e.g. 10m30s, see: https://pkg.go.dev/time#ParseDuration",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag := spotActionFlag{}
			err := flag.Set(tt.input)

			require.EqualError(t, err, tt.errMsg)
		})
	}
}

func TestSpotActionFlagType(t *testing.T) {
	flag := spotActionFlag{}

	assert.Equal(t, "stringToString", flag.Type())
}
