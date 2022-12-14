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

package patch

import (
	"bytes"
	"text/template"

	jsonpatch "github.com/evanphx/json-patch/v5"
)

const instanceTagTemplate = `[
	{
		"op": "add", 
		"path": "/tags", 
		"value": {
			"instance": { {{ jsonPairs .Tags }} }
		}
	}
]`

var instanceTagPatch = template.Must(template.New("InstanceTagPatch").
	Funcs(template.FuncMap{"jsonPairs": JSONPairs}).
	Parse(instanceTagTemplate))

// InstanceTag is used to patch a JSON document and replicate the inclusion
// of instance tags within the IMDS service. The same behaviour can be achieved
// by enabling access to the EC2 tags through the CLI, see:
// https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/Using_Tags.html#allow-access-to-tags-in-IMDS
type InstanceTag struct {
	Tags map[string]string
}

// Patch the JSON document with any provided instance tags. The resulting JSON
// document will conform to the IMDS specification and return EC2 tags within the
// IMDS metadata
func (p InstanceTag) Patch(in []byte) ([]byte, error) {
	if len(p.Tags) == 0 {
		return in, nil
	}

	var buf bytes.Buffer
	instanceTagPatch.Execute(&buf, p)

	patch, _ := jsonpatch.DecodePatch(buf.Bytes())

	out, err := patch.Apply(in)
	if err != nil {
		return in, err
	}

	return out, nil
}
