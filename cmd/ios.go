//Copyright 2017 Onegini B.V.
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.

package cmd

import (
	"os"
	"path"

	"fmt"

	"github.com/Onegini/onegini-sdk-configurator/util"
	"github.com/spf13/cobra"
)

var iosCmd = &cobra.Command{
	Use:   "ios",
	Short: "Configure an iOS project",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		config := util.ParseConfig(appDir, tsConfigLocation)
		var appTarget string

		if isCordova {
			config.ConfigureForCordova = true
			util.ParseCordovaConfig(config)
			rootDetection, debugDetection, debugLogs = util.ReadCordovaSecurityPreferences(config)
			appTarget = config.Cordova.AppName

			verifyIosPlatformInstalled("ERROR: Your project does not seem to have the iOS platform added. Please try `cordova platform add ios`")
		} else if isNativeScript {
			config.ConfigureForNativeScript = true
			util.ParseNativeScriptConfig(config)
			rootDetection, debugDetection, debugLogs = util.ReadNativeScriptSecurityPreferences(config)
			appTarget = targetName

			verifyIosPlatformInstalled("ERROR: Your project does not seem to have the iOS platform added. Please try `tns platform add ios`")
		} else {
			appTarget = targetName
		}
		verifyAppTarget(appTarget)
		util.SetAppTarget(appTarget, config)

		util.PrepareIosPaths(config)
		util.WriteIOSConfigModel(config)
		util.WriteIOSSecurityController(config, debugDetection, rootDetection, debugLogs, tamperingProtection)
		util.ConfigureIOSCertificates(config)

		util.PrintSuccessMessage(config, debugDetection, rootDetection, debugLogs, tamperingProtection)
		util.PrintIosInfoPlistUpdateHint(config)
	},
}

func verifyAppTarget(appTarget string) {
	if len(appTarget) == 0 {
		if isCordova {
			os.Stderr.WriteString(fmt.Sprintln("ERROR: No application identifier found in your 'config.xml'. Please make sure that you have set one."))
			os.Exit(1)
		} else {
			fmt.Print("ERROR: No target name provided. Provide one using 'onegini-sdk-configurator ios -t <target-name>'\n")
			fmt.Print("ERROR: More info on the target name can be found here: https://developer.apple.com/library/ios/documentation/IDEs/Conceptual/AppDistributionGuide/ConfiguringYourApp/ConfiguringYourApp.html\n\n")
			fmt.Print("execute 'onegini-sdk-configurator --help' to see how to use the configurator\n")
			os.Exit(1)
		}
	}
}

func verifyIosPlatformInstalled(errorMessage string) {
	_, err := os.Stat(path.Join(appDir, "platforms", "ios"))
	if os.IsNotExist(err) {
		os.Stderr.WriteString(fmt.Sprintln(errorMessage))
		os.Exit(1)
	}
}
