package cmd

import "github.com/spf13/cobra"

var (
	tsConfigLocation string
	appDir           string
	targetName       string
	rootDetection    bool
	debugDetection   bool
	isCordova        bool
)

func init() {
	RootCmd.AddCommand(androidCmd)
	RootCmd.AddCommand(iosCmd)
	RootCmd.PersistentFlags().StringVarP(&tsConfigLocation, "config", "c", "", "Path to Token Server config zip file")
	RootCmd.PersistentFlags().StringVarP(&appDir, "app-dir", "a", ".", "Path to application project root directory")
	RootCmd.PersistentFlags().StringVarP(&targetName, "target-name", "t", "", "Name of the target (for iOS)")
	RootCmd.PersistentFlags().BoolVarP(&debugDetection, "debugDetection", "d", true, "Enable or disable debug detection: --debugDetection=false")
	RootCmd.PersistentFlags().BoolVarP(&rootDetection, "rootDetection", "r", true, "Enable or disable root detection: --rootDetection=false")
	RootCmd.PersistentFlags().BoolVarP(&isCordova, "cordova", "o", false, "Configure as Cordova project")
}

var RootCmd = &cobra.Command{
	Use:   "onegini-sdk-configurator [platform]",
	Short: "Configure your onegini SDK",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}
