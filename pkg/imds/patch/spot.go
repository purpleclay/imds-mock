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
	},
	{
		"op": "add",
		"path": "/events",
		"value": {
			"recommendations": {
				"rebalance": {
					"noticeTime": "{{ .RebalanceTime }}"
				}
			}
		}
	}
]`

var spotPatch = template.Must(template.New("SpotPatch").
	Parse(spotTemplate))

// SpotInstanceAction is used to represent the lifecycle event (or action) of a spot instance
type SpotInstanceAction string

const (
	TerminateSpotInstanceAction SpotInstanceAction = "terminate"
	StopSpotInstanceAction      SpotInstanceAction = "stop"
	HibernateSpotInstanceAction SpotInstanceAction = "hibernate"
)

// Spot is used to patch a JSON document and replicate the use and lifecycle of a spot EC2 instance.
// The lifecycle of a spot instance is exposed through the use of an instance action category,
// see: https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/spot-instance-termination-notices.html#instance-action-metadata
type Spot struct {
	InstanceAction SpotInstanceAction
}

// Patch the document based on the provided instance action. The resulting JSON
// document will conform to the IMDS specification and expose spot details within
// the IMDS metadata
func (p Spot) Patch(in []byte) ([]byte, error) {
	// A spot interruption notice can be issued two minutes in advance during a best case scenario
	nowTime := time.Now().UTC()
	if p.InstanceAction != HibernateSpotInstanceAction {
		nowTime = nowTime.Add(2 * time.Minute)
	}

	now := nowTime.Format(time.RFC3339)

	spotDetails := struct {
		Action          SpotInstanceAction
		ActionTime      string
		TerminationTime string
		RebalanceTime   string
	}{
		Action:        p.InstanceAction,
		ActionTime:    now,
		RebalanceTime: now,
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
