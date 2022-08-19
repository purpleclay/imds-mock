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

package patch_test

import (
	"testing"

	"github.com/purpleclay/imds-mock/pkg/imds/patch"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/pretty"
)

func TestInstanceTagPatch(t *testing.T) {
	tagPatch := patch.InstanceTag{
		Tags: map[string]string{
			"Name":        "testing",
			"Environment": "dev",
		},
	}

	out, err := tagPatch.Patch([]byte(`{}`))
	require.NoError(t, err)

	// Pretty print and sort keys for predictable matching
	opts := pretty.DefaultOptions
	opts.SortKeys = true

	assert.Equal(t, `{
  "tags": {
    "instance": {
      "Environment": "dev",
      "Name": "testing"
    }
  }
}
`, string(pretty.PrettyOptions(out, opts)))
}

func TestInstanceTagPatch_NoTags(t *testing.T) {
	tagPatch := patch.InstanceTag{
		Tags: map[string]string{},
	}

	out, err := tagPatch.Patch([]byte(`{"testing":"123"}`))
	require.NoError(t, err)

	assert.Equal(t, `{"testing":"123"}`, string(out))
}

func TestInstanceTagPatch_InvalidInputJSON(t *testing.T) {
	tagPatch := patch.InstanceTag{
		Tags: map[string]string{
			"Name":        "testing",
			"Environment": "dev",
		},
	}

	_, err := tagPatch.Patch([]byte(`{`))
	require.Error(t, err)
}
