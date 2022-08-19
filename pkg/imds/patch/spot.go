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
	"time"

	jsonpatch "github.com/evanphx/json-patch/v5"
)

const spotTemplate = `[
	{
		"op": "replace",
		"path": "/instance-life-cycle",
		"value": "spot"
	},
	{
		"op": "add", 
		"path": "/spot", 
		"value": {
			"instance-action": {
				"action": "{{ .Action }}", 
				"time": "{{ .ActionTime }}"
			}{{ if .TerminationTime }},
			"termination-time": "{{ .TerminationTime }}"{{ end }}
		}
	}
]`

var spotPatch = template.Must(template.New("SpotPatch").
	Parse(spotTemplate))

// SpotInstanceAction ...
type SpotInstanceAction string

const (
	TerminateSpotInstanceAction SpotInstanceAction = "terminate"
	StopSpotInstanceAction      SpotInstanceAction = "stop"
	HibernateSpotInstanceAction SpotInstanceAction = "hibernate"
)

// Spot ...
type Spot struct {
	InstanceAction SpotInstanceAction
}

// Patch ...
func (p Spot) Patch(in []byte) ([]byte, error) {
	now := time.Now().UTC().Format(time.RFC3339)

	spotDetails := struct {
		Action          SpotInstanceAction
		ActionTime      string
		TerminationTime string
	}{
		Action:     p.InstanceAction,
		ActionTime: now,
	}

	if spotDetails.Action == TerminateSpotInstanceAction {
		// For backwards compatibility set the termination time
		spotDetails.TerminationTime = now
	}

	var buf bytes.Buffer
	spotPatch.Execute(&buf, spotDetails)

	patch, _ := jsonpatch.DecodePatch(buf.Bytes())

	out, err := patch.Apply(in)
	if err != nil {
		return in, err
	}

	return out, nil
}
