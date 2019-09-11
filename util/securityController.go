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

package util

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
)

func ReadCordovaSecurityPreferences(config *Config) (rootDetection bool, debugDetection bool, debugLogs bool) {
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
		if pref.Name == "OneginiDebugLogsEnabled" {
			var err error
			debugLogs, err = strconv.ParseBool(pref.Value)
			if err != nil {
				os.Stderr.WriteString(fmt.Sprintf("ERROR: could not parse 'OneginiDebugLogsEnabled' preference: %v\n", err.Error()))
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

func ReadNativeScriptSecurityPreferences(config *Config) (rootDetection bool, debugDetection bool, debugLogs bool) {
	rootDetectionSet := false
	debugDetectionSet := false

	if config.NativeScript.OneginiPreferences.DebugDetectionEnabled != nil {
		debugDetectionSet = true
		debugDetection = *config.NativeScript.OneginiPreferences.DebugDetectionEnabled
	}
	if config.NativeScript.OneginiPreferences.RootDetectionEnabled != nil {
		rootDetectionSet = true
		rootDetection = *config.NativeScript.OneginiPreferences.RootDetectionEnabled
	}
	if config.NativeScript.OneginiPreferences.DebugLogsEnabled {
		debugLogs = config.NativeScript.OneginiPreferences.DebugLogsEnabled
	}

	if !rootDetectionSet {
		rootDetection = true
	}
	if !debugDetectionSet {
		debugDetection = true
	}
	return
}

func WriteAndroidSecurityController(config *Config, debugDetection bool, rootDetection bool, debugLogs bool, tamperingProtection bool) {
	packageId := getPackageIdentifierFromConfig(config)
	storePath := config.getAndroidSecurityControllerPath()
	os.Remove(storePath) // always remove old file

	if rootDetection && debugDetection && !debugLogs && tamperingProtection {
		return
	}

	fileContents := `package %s;

@SuppressWarnings({ "unused", "WeakerAccess" })
public final class SecurityController {
%s}`

	flagsContents := PrepareFlagsForAndroid(debugDetection, rootDetection, debugLogs, tamperingProtection)
	fileContents = fmt.Sprintf(fileContents, packageId, flagsContents)

	if err := ioutil.WriteFile(storePath, []byte(fileContents), os.ModePerm); err != nil {
		log.Fatal("WARNING! Could not update security controller. This might be dangerous!")
	}
}

func WriteIOSSecurityController(config *Config, debugDetection bool, rootDetection bool, debugLogs bool, tamperingProtection bool) {
	group := "Configuration"
	headerContents := `#import <Foundation/Foundation.h>

@interface SecurityController : NSObject
+ (bool)rootDetection;
+ (bool)debugDetection;
+ (bool)debugLogs;
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
+(bool)debugLogs{
    return %s;
}
@end
`
	var (
		sDebugDetection string
		sRootDetection  string
		sDebugLogs      string
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

	if debugLogs {
		sDebugLogs = "YES"
	} else {
		sDebugLogs = "NO"
	}

	modelContents = fmt.Sprintf(modelContents, sRootDetection, sDebugDetection, sDebugLogs)
	xcodeProjPath := config.getIosXcodeProjPath()
	configModelPath := config.getIosConfigModelPath()

	headerStorePath := path.Join(configModelPath, "SecurityController.h")
	modelStorePath := path.Join(configModelPath, "SecurityController.m")

	if rootDetection && debugDetection && !debugLogs {
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

func PrepareFlagsForAndroid(debugDetection bool, rootDetection bool, debugLogs bool, tamperingProtection bool) string {
	// don't print unnecessary (default) flags
	var sb strings.Builder
	if !rootDetection {
		sb.WriteString("  public static final boolean rootDetection = false;\n")
	}
	if !debugDetection {
		sb.WriteString("  public static final boolean debugDetection = false;\n")
	}
	if debugLogs {
		sb.WriteString("  public static final boolean debugLogs = true;\n")
	}
	if !tamperingProtection {
		sb.WriteString("  public static final boolean tamperingProtection = false;\n")
	}
	return sb.String()
}
