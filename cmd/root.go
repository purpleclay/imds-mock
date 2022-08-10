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
	"io"

	"github.com/purpleclay/imds-mock/pkg/imds"
	"github.com/spf13/cobra"
)

func Execute(out io.Writer) error {
	opts := imds.DefaultOptions

	rootCmd := &cobra.Command{
		Use:          "imds-mock",
		Short:        "Mock the Amazon Instance Metadata Service (IMDS) for EC2",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
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

	rootCmd.AddCommand(newVersionCmd(out))
	rootCmd.AddCommand(newManPagesCmd(out))
	rootCmd.AddCommand(newCompletionCmd(out))

	return rootCmd.ExecuteContext(ctx.Background())
}
