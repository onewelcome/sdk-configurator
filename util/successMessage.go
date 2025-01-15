//Copyright 2019 Onegini B.V.
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

package util

import (
	"fmt"
	"strings"
)

func PrintSuccessMessage(config *Config) {
	fmt.Print("SUCCESS! Your application ")
	if len(config.FlavorName) > 0 {
		fmt.Printf("(\"%v\" flavor) ", config.FlavorName)
	}
	fmt.Print("is now configured.\n\n")

	fmt.Println("CONFIGURATION")

	fmt.Printf("App Identifier:		%v\n", config.Options.AppID)
	fmt.Printf("App Platform:		%v\n", config.Options.AppPlatform)
	fmt.Printf("App Version:		%v\n", config.Options.AppVersion)
	fmt.Printf("Redirect URI:		%v\n", config.Options.RedirectUrl)
	fmt.Printf("Token Server URI:	%v\n", config.Options.TokenServerUri)
	fmt.Printf("Server type:		%v\n", config.Options.ServerType)
	fmt.Printf("Server version:		%v\n", config.Options.ServerVersion)
	rgUris := config.Options.ResourceGatewayUris
	for i := 0; i < len(rgUris); i++ {
		if i == 0 {
			fmt.Printf("Resource Gateways:	%v\n", rgUris[i])
		} else {
			fmt.Printf("			%v\n", rgUris[i])
		}
	}
}

func PrintAndroidManifestUpdateHint(config *Config) {
	if config.ConfigureForCordova {
		return
	}
	if config.ConfigureForNativeScript {
		fmt.Println("")
		fmt.Println("INFO: Don't forget to update your android manifest to let Android handle the custom URL scheme")
		fmt.Println("INFO: The scheme that you must add: " + strings.Split(config.Options.RedirectUrl, "://")[0])
		fmt.Println("INFO: More info is provided here: https://docs.onegini.com/public/nativescript-plugin/topics/configuration.html#configuring-a-custom-url-scheme-for-authentication")
	} else {
		fmt.Println("")
		fmt.Println("INFO: Don't forget to update your android manifest to let Android handle the custom URL scheme")
		fmt.Println("INFO: The scheme that you must add: " + strings.Split(config.Options.RedirectUrl, "://")[0])
		fmt.Println("INFO: More info is provided here: https://docs.onegini.com/public/android-sdk/topics/authenticate-user-with-pin.html#handling-the-authentication-callback-during-registration")
	}
}

func PrintIosInfoPlistUpdateHint(config *Config) {
	if config.ConfigureForCordova {
		return
	}
	if config.ConfigureForNativeScript {
		fmt.Println("")
		fmt.Println("INFO: If you are using the system browser for user registration, don't forget to update your Info.plist to let iOS handle the custom URL scheme")
		fmt.Println("INFO: The scheme that you must add: " + strings.Split(config.Options.RedirectUrl, "://")[0])
		fmt.Println("INFO: More info is provided here: https://docs.onegini.com/public/nativescript-plugin/topics/configuration.html#configuring-a-custom-url-scheme-for-authentication")
	} else {
		fmt.Println("")
		fmt.Println("INFO: If you are using the system browser for user registration, don't forget to update your Info.plist to let iOS handle the custom URL scheme")
		fmt.Println("INFO: The scheme that you must add: " + strings.Split(config.Options.RedirectUrl, "://")[0])
		fmt.Println("INFO: More info is provided here: https://docs.onegini.com/public/ios-sdk/topics/user-authentication.html#handling-registration-request-url-with-external-web-browser")
	}
}
