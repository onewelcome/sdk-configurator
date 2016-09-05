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
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
)

func ReadCordovaSecurityPreferences(config *Config) (rootDetection bool, debugDetection bool) {
	rootDetectionSet := false
	debugDetectionSet := false

	for _, pref := range config.Cordova.Preferences {
		if pref.Name == "OneginiRootDetectionEnabled" {
			rootDetectionSet = true
			var err error
			rootDetection, err = strconv.ParseBool(pref.Value)
			if err != nil {
				os.Stderr.WriteString(fmt.Sprintf("ERROR: could not parse 'OneginiRootDetectionEnabled' preference: %v\n", err.Error()))
				os.Exit(1)
			}
		}
		if pref.Name == "OneginiDebugDetectionEnabled" {
			debugDetectionSet = true
			var err error
			debugDetection, err = strconv.ParseBool(pref.Value)
			if err != nil {
				os.Stderr.WriteString(fmt.Sprintf("ERROR: could not parse 'OneginiDebugDetectionEnabled' preference: %v\n", err.Error()))
				os.Exit(1)
			}
		}
	}

	if !rootDetectionSet {
		rootDetection = true
	}
	if !debugDetectionSet {
		debugDetection = true
	}
	return
}

func WriteAndroidSecurityController(config *Config, debugDetection bool, rootDetection bool) {
	fileContents := `package %s;

@SuppressWarnings({ "unused", "WeakerAccess" })
public final class SecurityController {
  public static final boolean debugDetection = %s;
  public static final boolean rootDetection = %s;
}`
	packageId := getPackageIdentifierFromConfig(config)
	fileContents = fmt.Sprintf(fileContents, packageId, strconv.FormatBool(debugDetection), strconv.FormatBool(rootDetection))
	storePath := config.getAndroidSecurityControllerPath()

	if rootDetection && debugDetection {
		os.Remove(storePath)
	} else {
		if err := ioutil.WriteFile(storePath, []byte(fileContents), os.ModePerm); err != nil {
			log.Fatal("WARNING! Could not update security controller. This might be dangerous!")
		}
	}
}

func WriteIOSSecurityController(config *Config, debugDetection bool, rootDetection bool) {
	group := "Configuration"
	headerContents := `#import <Foundation/Foundation.h>

@interface SecurityController : NSObject
+ (bool)rootDetection;
+ (bool)debugDetection;
@end
`

	modelContents := `#import "SecurityController.h"

@implementation SecurityController
+(bool)rootDetection{
    return %s;
}
+(bool)debugDetection{
    return %s;
}
@end
`
	var (
		sDebugDetection string
		sRootDetection  string
	)

	if debugDetection {
		sDebugDetection = "YES"
	} else {
		sDebugDetection = "NO"
	}

	if rootDetection {
		sRootDetection = "YES"
	} else {
		sRootDetection = "NO"
	}

	modelContents = fmt.Sprintf(modelContents, sRootDetection, sDebugDetection)
	xcodeProjPath := config.getIosXcodeProjPath()
	configModelPath := config.getIosConfigModelPath()

	headerStorePath := path.Join(configModelPath, "SecurityController.h")
	modelStorePath := path.Join(configModelPath, "SecurityController.m")

	if rootDetection && debugDetection {
		removeFileFromXcodeProj(headerStorePath, xcodeProjPath, group)
		removeFileFromXcodeProj(modelStorePath, xcodeProjPath, group)
		os.Remove(headerStorePath)
		os.Remove(modelStorePath)
	} else {
		ioutil.WriteFile(headerStorePath, []byte(headerContents), os.ModePerm)
		ioutil.WriteFile(modelStorePath, []byte(modelContents), os.ModePerm)
		addFileToXcodeProj(headerStorePath, xcodeProjPath, config.AppTarget, group)
		addFileToXcodeProj(modelStorePath, xcodeProjPath, config.AppTarget, group)
	}
}
