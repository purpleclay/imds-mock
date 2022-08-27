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
	ctx "context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/purpleclay/imds-mock/pkg/imds"
	"github.com/purpleclay/imds-mock/pkg/imds/patch"
	"github.com/spf13/cobra"
)

// Custom flag for parsing a spot action
type spotActionFlag struct {
	event *imds.SpotActionEvent
}

func (e *spotActionFlag) String() string {
	return "terminate=0s"
}

func (e *spotActionFlag) Set(value string) error {
	action, duration, found := strings.Cut(value, "=")
	if !found {
		return fmt.Errorf("%s must be formatted as key=value e.g. terminate=2s", value)
	}

	spotAction := patch.SpotInstanceAction(action)
	switch spotAction {
	case patch.TerminateSpotInstanceAction:
		fallthrough
	case patch.HibernateSpotInstanceAction:
		fallthrough
	case patch.StopSpotInstanceAction:
		break
	default:
		return fmt.Errorf("%s is not a supported spot action expecting (%s, %s or %s)",
			action, patch.TerminateSpotInstanceAction, patch.StopSpotInstanceAction, patch.HibernateSpotInstanceAction)
	}

	eventTime, err := time.ParseDuration(duration)
	if err != nil {
		return fmt.Errorf("%s is not a supported duration format e.g. 10m30s, see: https://pkg.go.dev/time#ParseDuration", duration)
	}

	e.event = &imds.SpotActionEvent{
		Action:   spotAction,
		Duration: eventTime,
	}

	return nil
}

func (e *spotActionFlag) Type() string {
	return "stringToString"
}

func Execute(out io.Writer) error {
	opts := imds.DefaultOptions

	// flag to parse custom spot action
	var spotAction spotActionFlag

	rootCmd := &cobra.Command{
		Use:          "imds-mock",
		Short:        "Easy mocking of the Amazon EC2 Instance Metadata Service (IMDS)",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if spotAction.event != nil {
				// Overwrite the default spot action, since a custom one has been provided
				opts.SpotAction = imds.SpotActionEvent{
					Action:   spotAction.event.Action,
					Duration: spotAction.event.Duration,
				}
			}

			_, err := imds.ServeWith(opts)
			return err
		},
	}

	flags := rootCmd.Flags()
	flags.BoolVar(&opts.ExcludeInstanceTags, "exclude-instance-tags", imds.DefaultOptions.ExcludeInstanceTags, "exclude access to instance tags associated with the instance")
	flags.BoolVar(&opts.IMDSv2, "imdsv2", imds.DefaultOptions.IMDSv2, "enforce IMDSv2 requiring all requests to contain a valid metadata token")
	flags.StringToStringVar(&opts.InstanceTags, "instance-tags", imds.DefaultOptions.InstanceTags, "a list of instance tags (key pairs) to expose as metadata")
	flags.IntVar(&opts.Port, "port", imds.DefaultOptions.Port, "the port to be used at startup")
	flags.BoolVar(&opts.Pretty, "pretty", imds.DefaultOptions.Pretty, "if instance categories should return pretty printed JSON")
	flags.BoolVar(&opts.Spot, "spot", imds.DefaultOptions.Spot, "enable simulation of a spot instance and interruption notice")
	flags.Var(&spotAction, "spot-action", "configure the type and delay of the spot interruption notice")

	rootCmd.AddCommand(newVersionCmd(out))
	rootCmd.AddCommand(newManPagesCmd(out))
	rootCmd.AddCommand(newCompletionCmd(out))

	return rootCmd.ExecuteContext(ctx.Background())
}
