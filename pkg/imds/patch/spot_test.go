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
	"encoding/json"
	"testing"
	"time"

	"github.com/purpleclay/imds-mock/pkg/imds/patch"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Only used for serialising the patch into a struct for assertions
type spotPatchJSON struct {
	InstanceLifeCycle string `json:"instance-life-cycle"`
	Spot              struct {
		InstanceAction struct {
			Action     string    `json:"action"`
			ActionTime time.Time `json:"time"`
		} `json:"instance-action"`

		TerminationTime *time.Time `json:"termination-time"`
	} `json:"spot"`
	Events struct {
		Recommendations struct {
			Rebalance struct {
				NoticeTime time.Time `json:"noticeTime"`
			} `json:"rebalance"`
		} `json:"recommendations"`
	} `json:"events"`
}

func TestSpotPatch_TerminateAction(t *testing.T) {
	tests := []struct {
		name   string
		action patch.SpotInstanceAction
	}{
		{
			name:   "TerminateAction",
			action: patch.TerminateSpotInstanceAction,
		},
		{
			name:   "StopAction",
			action: patch.StopSpotInstanceAction,
		},
		{
			name:   "HibernateAction",
			action: patch.HibernateSpotInstanceAction,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spotPatch := patch.Spot{
				InstanceAction: tt.action,
			}

			out, err := spotPatch.Patch([]byte(`{"instance-life-cycle":"on-demand"}`))
			require.NoError(t, err)

			var spotJSON spotPatchJSON
			json.Unmarshal(out, &spotJSON)

			now := time.Now().UTC().Add(2 * time.Minute)

			assert.Equal(t, "spot", spotJSON.InstanceLifeCycle)
			assert.Equal(t, string(tt.action), spotJSON.Spot.InstanceAction.Action)
			assert.WithinDuration(t, now, spotJSON.Spot.InstanceAction.ActionTime, 1*time.Second)

			if tt.action == patch.TerminateSpotInstanceAction {
				require.NotNil(t, spotJSON.Spot.TerminationTime)
				assert.WithinDuration(t, now, *spotJSON.Spot.TerminationTime, 1*time.Second)
			} else {
				require.Nil(t, spotJSON.Spot.TerminationTime)
			}

			assert.WithinDuration(t, now, spotJSON.Events.Recommendations.Rebalance.NoticeTime, 1*time.Second)
		})
	}
}

func TestSpotPatch_InvalidInputJSON(t *testing.T) {
	spotPatch := patch.Spot{
		InstanceAction: patch.HibernateSpotInstanceAction,
	}

	_, err := spotPatch.Patch([]byte(`{`))
	require.Error(t, err)
}
