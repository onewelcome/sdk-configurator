//Copyright 2016 Onegini B.V.
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
	"github.com/spf13/cobra"
	"github.com/Onegini/sdk-configurator/util"
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

