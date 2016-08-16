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
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"

	"fmt"

	"github.com/Onegini/onegini-sdk-configurator/data"
)

func WriteIOSConfigModel(appDir string, appName string, config *Config) {
	xcodeProjPath := getIosXcodeProjPath(appDir, appName, config)

	modelMFilePath := getIosConfigModelPathMFile(appDir, appName, config)
	modelHFilePath := getIosConfigModelPathHFile(appDir, appName, config)
	modelMFile := readIosConfigModelFromAssetsOrProject(modelMFilePath, "lib/OneginiConfigModel.m")
	modelHFile := readIosConfigModelFromAssetsOrProject(modelHFilePath, "lib/OneginiConfigModel.h")

	base64Certs := getBase64Certs(config)
	modelMFile = overrideIosConfigModelValues(config, base64Certs, modelMFile)

	ioutil.WriteFile(modelMFilePath, modelMFile, os.ModePerm)
	ioutil.WriteFile(modelHFilePath, modelHFile, os.ModePerm)

	iosAddConfigModelFileToXcodeProj(modelMFilePath, xcodeProjPath, appName)
	iosAddConfigModelFileToXcodeProj(modelHFilePath, xcodeProjPath, appName)
}

func readIosConfigModelFromAssetsOrProject(modelPath string, assetPath string) []byte {
	_, errFileNotFoundInAppProject := os.Stat(modelPath)
	if errFileNotFoundInAppProject == nil {
		appProjectModel, err := ioutil.ReadFile(modelPath)
		if err != nil {
			os.Stderr.WriteString(fmt.Sprintf("ERROR: could not read Config model in Project: %v\n", err.Error()))
			os.Exit(1)
		}
		return appProjectModel
	} else {
		modelFromTmp, errFileNotFoundInTmp := data.Asset(assetPath)
		if errFileNotFoundInTmp != nil {
			os.Stderr.WriteString(fmt.Sprintf("ERROR: could not read Config model in assets: %v\n", errFileNotFoundInTmp))
			os.Exit(1)
		}

		return modelFromTmp
	}
}

func overrideIosConfigModelValues(config *Config, base64Certs []string, model []byte) []byte {
	configMap := map[string]string{
		"ONGAppIdentifier":   config.Options.AppID,
		"ONGAppVersion":      config.Options.AppVersion,
		"ONGAppBaseURL":      config.Options.TokenServerUri,
		"ONGResourceBaseURL": config.Options.ResourceGatewayUris[0],
		"ONGRedirectURL":     config.Options.RedirectUrl,
	}

	for preference, value := range configMap {
		newPref := `@"` + preference + `" : @"` + value + `"`
		re := regexp.MustCompile(`@"` + preference + `"\s*:\s*@".*"`)
		model = re.ReplaceAll(model, []byte(newPref))
	}

	newDef := "return @[@\"" + strings.Join(base64Certs, "\", @\"") + "\"]; //Base64Certificates"

	re := regexp.MustCompile(`return @\[.*\];.*`)
	model = re.ReplaceAll(model, []byte(newDef))

	return model
}

func WriteAndroidConfigModel(config *Config, appDir string, keystorePath string) {
	modelPath := getAndroidConfigModelPath(appDir, config)

	model := readAndroidConfigModelFromAssetsOrProject(modelPath)
	model = overrideAndroidConfigModelValues(config, keystorePath, model)
	ioutil.WriteFile(modelPath, model, os.ModePerm)
}

func overrideAndroidConfigModelValues(config *Config, keystorePath string, model []byte) []byte {
	stringConfigMap := map[string]string{
		"appIdentifier":   config.Options.AppID,
		"redirectionUri":  config.Options.RedirectUrl,
		"appVersion":      config.Options.AppVersion,
		"baseURL":         config.Options.TokenServerUri,
		"resourceBaseURL": config.Options.ResourceGatewayUris[0],
		"keystoreHash":    CalculateKeystoreHash(keystorePath),
	}
	intConfigMap := map[string]string{
		"maxPinFailures": strconv.Itoa(config.Options.MaxPinFailures),
	}

	newPackage := "package " + getPackageIdentifierFromConfig(config) + ";"
	packageRe := regexp.MustCompile(`package\s.*;`)
	model = packageRe.ReplaceAll(model, []byte(newPackage))

	for preference, value := range stringConfigMap {
		newPref := preference + ` = "` + value + `";`
		re := regexp.MustCompile(preference + `\s=\s".*";`)
		model = re.ReplaceAll(model, []byte(newPref))
	}

	for preference, value := range intConfigMap {
		newPref := preference + ` = ` + value + `;`
		re := regexp.MustCompile(preference + `\s=\s.*;`)
		model = re.ReplaceAll(model, []byte(newPref))
	}

	return model
}

func readAndroidConfigModelFromAssetsOrProject(modelPath string) []byte {
	_, errFileNotFoundInAppProject := os.Stat(modelPath)
	if errFileNotFoundInAppProject == nil {
		appProjectModel, err := ioutil.ReadFile(modelPath)
		if err != nil {
			os.Stderr.WriteString(fmt.Sprintf("ERROR: Could not read config model in Project: %v\n", err.Error()))
			os.Exit(1)
		}
		return appProjectModel
	} else {
		modelFromTmp, errFileNotFoundInTmp := data.Asset("lib/OneginiConfigModel.java")
		if errFileNotFoundInTmp != nil {
			os.Stderr.WriteString(fmt.Sprintf("ERROR: Could not read config model in assets: %v\n", errFileNotFoundInTmp))
			os.Exit(1)
		}

		return modelFromTmp
	}
}
