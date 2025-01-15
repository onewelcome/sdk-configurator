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

import "github.com/spf13/cobra"

var (
	tsConfigLocation string
	appDir           string
	targetName       string
	moduleName       string
	flavorName       string
	isCordova        bool
	isNativeScript   bool
)

func init() {
	RootCmd.AddCommand(androidCmd)
	RootCmd.AddCommand(iosCmd)
	RootCmd.AddCommand(versionCmd)
	RootCmd.PersistentFlags().StringVarP(&tsConfigLocation, "config", "c", "", "Path to Token Server config zip file")
	RootCmd.PersistentFlags().StringVarP(&appDir, "app-dir", "a", ".", "Path to application project root directory")
	RootCmd.PersistentFlags().StringVarP(&targetName, "target-name", "t", "", "The target name in your Xcode project for which you want to configure the SDK (for iOS). More info can be found at https://developer.apple.com/library/ios/documentation/IDEs/Conceptual/AppDistributionGuide/ConfiguringYourApp/ConfiguringYourApp.html")
	RootCmd.PersistentFlags().StringVarP(&moduleName, "module-name", "m", "", "The Gradle module name that contains your application sources (for Android). More info can be found at https://developer.android.com/studio/projects/index.html")
	RootCmd.PersistentFlags().StringVarP(&flavorName, "flavor-name", "f", "", "The optional flavor name for Android project (or destination subfolder for iOS). More info can be found at https://developer.android.com/studio/build/build-variants#product-flavors")
	RootCmd.PersistentFlags().BoolVarP(&isCordova, "cordova", "o", false, "Configure as Cordova project")
	RootCmd.PersistentFlags().BoolVarP(&isNativeScript, "nativescript", "n", false, "Configure as NativeScript project")
	_ = RootCmd.PersistentFlags().MarkHidden("tamperingProtection")
}

var RootCmd = &cobra.Command{
	Use:   "sdk-configurator [platform]",
	Short: "Configure your mobile SDK",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}
