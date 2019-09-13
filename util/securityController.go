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
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
)

func ReadCordovaSecurityPreference(preferences []cordovaPreference, preferenceName string, defaultValue bool) bool {
	for _, pref := range preferences {
		if pref.Name == preferenceName {
			var err error
			var value bool
			value, err = strconv.ParseBool(pref.Value)
			if err != nil {
				_, _ = os.Stderr.WriteString(fmt.Sprintf("ERROR: could not parse '%s' preference: %v\n", preferenceName, err.Error()))
				os.Exit(1)
			} else {
				fmt.Printf("WARNING: config.xml contains %s=%t, this value will be used in the SecurityController\n", preferenceName, value)
				return value
			}
		}
	}
	return defaultValue
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
	_ = os.Remove(storePath) // always remove old file

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
	xcodeProjPath := config.getIosXcodeProjPath()
	configModelPath := config.getIosConfigModelPath()
	headerStorePath := path.Join(configModelPath, "SecurityController.h")
	modelStorePath := path.Join(configModelPath, "SecurityController.m")

	// always remove old files
	removeFileFromXcodeProj(headerStorePath, xcodeProjPath, group)
	removeFileFromXcodeProj(modelStorePath, xcodeProjPath, group)
	_ = os.Remove(headerStorePath)
	_ = os.Remove(modelStorePath)

	if rootDetection && debugDetection && !debugLogs && tamperingProtection {
		return
	}

	headerContents := `#import <Foundation/Foundation.h>

@interface SecurityController : NSObject
%s@end
`
	headerContents = fmt.Sprintf(headerContents, PrepareHeaderFlagsForIos(debugDetection, rootDetection, debugLogs, tamperingProtection))

	modelContents := `#import "SecurityController.h"

@implementation SecurityController
%s@end
`
	modelContents = fmt.Sprintf(modelContents, PrepareModelFlagsForIos(debugDetection, rootDetection, debugLogs, tamperingProtection))

	_ = ioutil.WriteFile(headerStorePath, []byte(headerContents), os.ModePerm)
	_ = ioutil.WriteFile(modelStorePath, []byte(modelContents), os.ModePerm)
	addFileToXcodeProj(headerStorePath, xcodeProjPath, config.AppTarget, group)
	addFileToXcodeProj(modelStorePath, xcodeProjPath, config.AppTarget, group)
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

func PrepareHeaderFlagsForIos(debugDetection bool, rootDetection bool, debugLogs bool, tamperingProtection bool) string {
	// don't print unnecessary (default) flags
	var sb strings.Builder
	if !rootDetection {
		sb.WriteString("+ (bool)rootDetection;\n")
	}
	if !debugDetection {
		sb.WriteString("+ (bool)debugDetection;\n")
	}
	if debugLogs {
		sb.WriteString("+ (bool)debugLogs;\n")
	}
	if !tamperingProtection {
		sb.WriteString("+ (bool)tamperingProtection;\n")
	}
	return sb.String()
}

func PrepareModelFlagsForIos(debugDetection bool, rootDetection bool, debugLogs bool, tamperingProtection bool) string {
	// don't print unnecessary (default) flags
	var sb strings.Builder
	if !rootDetection {
		sb.WriteString("+(bool)rootDetection{\n  return NO;\n}\n")
	}
	if !debugDetection {
		sb.WriteString("+(bool)debugDetection{\n  return NO;\n}\n")
	}
	if debugLogs {
		sb.WriteString("+(bool)debugLogs{\n  return YES;\n}\n")
	}
	if !tamperingProtection {
		sb.WriteString("+(bool)tamperingProtection{\n  return NO;\n}\n")
	}
	return sb.String()
}
