package cmd

import (
	"fmt"

	"os"

	"github.com/spf13/cobra"
)

const version = "v2.0.0"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print onegini-sdk-configurator version and exit",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("onegini-sdk-configurator %s\n", version)
		os.Exit(0)
	},
}
