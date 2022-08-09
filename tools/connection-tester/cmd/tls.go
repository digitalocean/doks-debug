/*
Copyright 2022 Wayne Warren wayne.warren.s@gmail.com

*/
package cmd

import (
	"github.com/spf13/cobra"

	"github.com/digitalocean/doks-debug/tools/connection-tester/pkg/connections"
)

// tlsCmd represents the tls subcommand
var tlsCmd = &cobra.Command{
	Use:   "tls",
	Short: "Open (and promptly close) a barrage of TLS connections to the target host.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var (
			targetHost = args[0]
			ctx        = cmd.Context()
			c          = connections.NewTLSConnector(targetHost)
		)
		connections.TryLoops(ctx, numTryers, c)
	},
}

func init() {
	rootCmd.AddCommand(tlsCmd)
}
