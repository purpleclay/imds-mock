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
}

func TestSpotPatch_TerminateAction(t *testing.T) {
	spotPatch := patch.Spot{
		InstanceAction: patch.TerminateSpotInstanceAction,
	}

	out, err := spotPatch.Patch([]byte(`{"instance-life-cycle":"on-demand"}`))
	require.NoError(t, err)

	var spotJSON spotPatchJSON
	json.Unmarshal(out, &spotJSON)

	now := time.Now().UTC()

	assert.Equal(t, "spot", spotJSON.InstanceLifeCycle)
	assert.Equal(t, "terminate", spotJSON.Spot.InstanceAction.Action)
	assert.WithinDuration(t, now, spotJSON.Spot.InstanceAction.ActionTime, 1*time.Second)

	require.NotNil(t, spotJSON.Spot.TerminationTime)
	assert.WithinDuration(t, now, *spotJSON.Spot.TerminationTime, 1*time.Second)
}

func TestSpotPatch_StopAction(t *testing.T) {
	spotPatch := patch.Spot{
		InstanceAction: patch.StopSpotInstanceAction,
	}

	out, err := spotPatch.Patch([]byte(`{"instance-life-cycle":"on-demand"}`))
	require.NoError(t, err)

	var spotJSON spotPatchJSON
	json.Unmarshal(out, &spotJSON)

	now := time.Now().UTC()

	assert.Equal(t, "spot", spotJSON.InstanceLifeCycle)
	assert.Equal(t, "stop", spotJSON.Spot.InstanceAction.Action)
	assert.WithinDuration(t, now, spotJSON.Spot.InstanceAction.ActionTime, 1*time.Second)
	assert.Nil(t, spotJSON.Spot.TerminationTime)
}

func TestSpotPatch_HibernateAction(t *testing.T) {
	spotPatch := patch.Spot{
		InstanceAction: patch.HibernateSpotInstanceAction,
	}

	out, err := spotPatch.Patch([]byte(`{"instance-life-cycle":"on-demand"}`))
	require.NoError(t, err)

	var spotJSON spotPatchJSON
	json.Unmarshal(out, &spotJSON)

	now := time.Now().UTC()

	assert.Equal(t, "spot", spotJSON.InstanceLifeCycle)
	assert.Equal(t, "hibernate", spotJSON.Spot.InstanceAction.Action)
	assert.WithinDuration(t, now, spotJSON.Spot.InstanceAction.ActionTime, 1*time.Second)
	assert.Nil(t, spotJSON.Spot.TerminationTime)
}
