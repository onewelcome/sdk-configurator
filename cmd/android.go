package cmd

import (
	"github.com/spf13/cobra"
	"gitlab.onegini.com/mobile-platform/onegini-sdk-configurator/util"
	"os"
	"path"
	"fmt"
)

var androidCmd = &cobra.Command{
	Use:   "android",
	Short: "Configure an Android project",
	Run: func(cmd *cobra.Command, args []string) {
		var keystorePath string

		config := util.ParseConfig(tsConfigLocation, cmd)

		if isCordova {
			util.ParseCordovaConfig(appDir, config)
			rootDetection, debugDetection = util.ReadCordovaSecurityPreferences(config)
			verifyCordovaAndroidPlatformInstalled()

			util.WriteAndroidSecurityController(appDir, config, debugDetection, rootDetection)
		} else {
			util.ParseAndroidManifest(appDir, config)

			util.WriteAndroidSecurityController(appDir, config, debugDetection, rootDetection)
		}

		keystorePath = util.GetAndroidKeystorePath(appDir, config)

		util.WriteAndroidAppScheme(appDir, config)
		util.CreateKeystore(config, keystorePath)
		util.WriteAndroidConfigModel(config, appDir, keystorePath)
		util.PrintSuccessMessage(config, debugDetection, rootDetection)
		util.PrintAndroidManifestUpdateHint(config)
	},
}

func verifyCordovaAndroidPlatformInstalled() {
	_, err := os.Stat(path.Join(appDir, "platforms", "android"))
	if os.IsNotExist(err) {
		os.Stderr.WriteString(fmt.Sprintln("ERROR: Your project does not seem to have the Android platform added. Please try `cordova platform add android`"))
		os.Exit(1)
	}
}

