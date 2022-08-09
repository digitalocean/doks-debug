/*
Copyright © 2022 Wayne Warren wayne.warren.s@gmail.com

*/
package cmd

import (
	"github.com/spf13/cobra"

	"github.com/digitalocean/doks-debug/tools/connection-tester/pkg/connections"
)

var (
	httpRequestSizeBytes int
	httpRequestMethod    string
)

// httpCmd represents the http command
var httpCmd = &cobra.Command{
	Use:   "http",
	Short: "Repeatedly make HTTP requests to target host.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var (
			targetHost = args[0]
			ctx        = cmd.Context()
			c          = connections.NewHTTPConnector(targetHost,
				connections.HTTPConnectorWithRandomBody(httpRequestSizeBytes),
				connections.HTTPRequestMethod(httpRequestMethod),
			)
		)

		connections.TryLoops(ctx, numTryers, c)
	},
}

func init() {
	rootCmd.AddCommand(httpCmd)

	httpCmd.Flags().IntVar(&httpRequestSizeBytes, "request-size-bytes", 1024*1024, "size of randomly-generated request body size")
	httpCmd.Flags().StringVar(&httpRequestMethod, "request-method", "PATCH", "http request method to use")
}
