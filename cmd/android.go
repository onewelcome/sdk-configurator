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
	"fmt"
	"os"
	"path"

	"github.com/Onegini/onegini-sdk-configurator/util"
	"github.com/spf13/cobra"
)

var androidCmd = &cobra.Command{
	Use:   "android",
	Short: "Configure an Android project",
	Run: func(cmd *cobra.Command, args []string) {
		config := util.ParseConfig(appDir, tsConfigLocation)

		verifyAppModuleName(moduleName)
		util.SetAppTarget(moduleName, config)

		if isCordova {
			config.ConfigureForCordova = true
			util.ParseCordovaConfig(config)
			rootDetection, debugDetection, debugLogs = util.ReadCordovaSecurityPreferences(config)
			verifyAndroidPlatformInstalled("ERROR: Your project does not seem to have the Android platform added. Please try `cordova platform add android`")
		} else if isNativeScript {
			config.ConfigureForNativeScript = true
			util.SetAppTarget("", config)
			util.ParseNativeScriptConfig(config)
			rootDetection, debugDetection, debugLogs = util.ReadNativeScriptSecurityPreferences(config)
			verifyAndroidPlatformInstalled("ERROR: Your project does not seem to have the Android platform added. Please try `tns platform add android`")
		}
		util.ParseAndroidManifest(config)

		util.WriteAndroidSecurityController(config, debugDetection, rootDetection, debugLogs)
		util.WriteAndroidAppScheme(config)
		util.CreateKeystore(config)
		util.WriteAndroidConfigModel(config)
		util.PrintSuccessMessage(config, debugDetection, rootDetection, debugLogs)
		util.PrintAndroidManifestUpdateHint(config)
	},
}

func verifyAndroidPlatformInstalled(errorMessage string) {
	_, err := os.Stat(path.Join(appDir, "platforms", "android"))
	if os.IsNotExist(err) {
		os.Stderr.WriteString(fmt.Sprintln(errorMessage))
		os.Exit(1)
	}
}

func verifyAppModuleName(moduleName string) {
	if isCordova || isNativeScript {
		if len(moduleName) != 0 {
			fmt.Println("WARNING: Ignoring the module name parameter for Cordova or NativeScript")
		}
	} else {
		if len(moduleName) == 0 {
			fmt.Print("ERROR: No module name provided. Provide one using 'onegini-sdk-configurator android -m <module-name>'\n")
			fmt.Print("ERROR: More info on the module name can be found here: https://developer.android.com/studio/projects/index.html\n\n")
			fmt.Print("execute 'onegini-sdk-configurator --help' to see how to use the configurator\n")
			os.Exit(1)
		}
	}
}
