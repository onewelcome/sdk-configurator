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

package util

import (
	"fmt"
	"strconv"
	"strings"
)

func PrintSuccessMessage(config *Config, debugDetection bool, rootDetection bool) {
	fmt.Print("SUCCESS! Your application is now configured.\n\n")
	fmt.Println("CONFIGURATION")

	fmt.Printf("App Identifier:		%v\n", config.Options.AppID)
	fmt.Printf("App Platform:		%v\n", config.Options.AppPlatform)
	fmt.Printf("App Version:		%v\n", config.Options.AppVersion)
	fmt.Printf("Redirect URI:		%v\n", config.Options.RedirectUrl)
	fmt.Printf("Max PIN failures:	%v\n", strconv.Itoa(config.Options.MaxPinFailures))
	fmt.Printf("Debug detection:	%v\n", debugDetection)
	fmt.Printf("Root detection:		%v\n", rootDetection)
	fmt.Printf("Token Server URI:	%v\n", config.Options.TokenServerUri)
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
	if !isCordova(config) {
		fmt.Println("")
		fmt.Println("INFO: Don't forget to update your android manifest to let Android handle the custom URL scheme")
		fmt.Println("INFO: The scheme that you must add: " + strings.Split(config.Options.RedirectUrl, "://")[0])
		fmt.Println("INFO: More info is provided here: https://docs.onegini.com/public/android-sdk/topics/authenticate-user-with-pin.html#handling-the-authentication-callback-during-registration")
	}
}
