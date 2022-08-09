/*
Copyright 2022 Wayne Warren wayne.warren.s@gmail.com

*/
package cmd

import (
	"context"
	"os"

	"github.com/spf13/cobra"
)

var (
	numTryers int
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "connection-tester",
	Short: "a tool for testing a variety of network connection behaviors",
}

// Execute adds all child commands to the root command and sets flags
// appropriately. This is called by main.main(). It only needs to happen once
// to the rootCmd.
func Execute(ctx context.Context) {
	err := rootCmd.ExecuteContext(ctx)
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	httpCmd.PersistentFlags().IntVar(&numTryers, "num-tryers", 5, "number of goroutines attempting connections")
}
