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
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/onewelcome/sdk-configurator/version"

	"fmt"

	"github.com/onewelcome/sdk-configurator/data"
)

func WriteIOSConfigModel(config *Config) {
	cleanupOldIosConfigModel(config)

	modelMFile := overrideIosConfigModelValues(config)
	modelHFile := readIosConfigModelFromAssetsOrProject(config.getIosConfigModelPathHFile(), "lib/OneginiConfigModel.h")

	WriteIosConfigModel(modelMFile, modelHFile, config)
}
func WriteIosConfigModel(modelMFile []byte, modelHFile []byte, config *Config) {
	xcodeProjPath := config.getIosXcodeProjPath()
	modelMFilePath := config.getIosConfigModelPathMFile()
	modelHFilePath := config.getIosConfigModelPathHFile()

	ioutil.WriteFile(modelMFilePath, modelMFile, os.ModePerm)
	ioutil.WriteFile(modelHFilePath, modelHFile, os.ModePerm)

	iosAddConfigModelFileToXcodeProj(modelMFilePath, xcodeProjPath, config.AppTarget, config.FlavorName)
	iosAddConfigModelFileToXcodeProj(modelHFilePath, xcodeProjPath, config.AppTarget, config.FlavorName)
}

func cleanupOldIosConfigModel(config *Config) {
	modelMFilePath := config.getIosConfigModelPathMFile()
	modelHFilePath := config.getIosConfigModelPathHFile()

	deleteFileIfExists(modelMFilePath, "ERROR: Could not delete old config model M file in Project")
	deleteFileIfExists(modelHFilePath, "ERROR: Could not delete old config model H file in Project")

	iosRemoveConfigModelFileFromXcodeProj(modelMFilePath, config.getIosXcodeProjPath(), config.FlavorName)
	iosRemoveConfigModelFileFromXcodeProj(modelHFilePath, config.getIosXcodeProjPath(), config.FlavorName)
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

func overrideIosConfigModelValues(config *Config) (modelMFile []byte) {
	modelMFile = readIosConfigModelFromAssetsOrProject(config.getIosConfigModelPathMFile(), "lib/OneginiConfigModel.m")

	base64Certs := getBase64Certs(config)

	configMap := map[string]string{
		"ONGServerType":      config.Options.ServerType,
		"ONGServerVersion":   config.Options.ServerVersion,
		"ONGAppIdentifier":   config.Options.AppID,
		"ONGAppVersion":      config.Options.AppVersion,
		"ONGAppBaseURL":      config.Options.TokenServerUri,
		"ONGResourceBaseURL": config.Options.ResourceGatewayUris[0],
		"ONGRedirectURL":     config.Options.RedirectUrl,
	}

	for preference, value := range configMap {
		newPref := `@"` + preference + `" : @"` + value + `"`
		re := regexp.MustCompile(`@"` + preference + `"\s*:\s*@".*"`)
		modelMFile = re.ReplaceAll(modelMFile, []byte(newPref))
	}

	newDef := "certificates\n{\n	return @[@\"" + strings.Join(base64Certs, "\", @\"") + "\"]; //Base64Certificates"

	re := regexp.MustCompile(`certificates\s*{\s*return @\[.*\];.*`)
	modelMFile = re.ReplaceAll(modelMFile, []byte(newDef))

	serverPublicKeyNewDef := "serverPublicKey\n{\n	return @\"" + config.Options.ServerPublicKey.Encoded + "\";"
	reServerPublicKey := regexp.MustCompile(`serverPublicKey\s*{\s*return @\".*\";`)
	modelMFile = reServerPublicKey.ReplaceAll(modelMFile, []byte(serverPublicKeyNewDef))

	versionRe := regexp.MustCompile(`CONFIGURATOR_VERSION`)
	modelMFile = versionRe.ReplaceAll(modelMFile, []byte(version.Version))

	return
}

func WriteAndroidConfigModel(config *Config, generateJavaConfigModel bool) {
	modelJavaPath := config.getAndroidConfigModelJavaPath()
	modelKotlinPath := config.getAndroidConfigModelKotlinPath()
	keyStorePath := config.getAndroidKeystorePath()

	deleteFileIfExists(modelJavaPath, "ERROR: Could not delete old java config model in Project")
	deleteFileIfExists(modelKotlinPath, "ERROR: Could not delete old kotlin config model in Project")

	if generateJavaConfigModel {
		model := readAndroidJavaConfigModelFromAssets()
		model = overrideAndroidConfigJavaModelValues(config, keyStorePath, model)
		os.WriteFile(modelJavaPath, model, os.ModePerm)
	} else {
		model := readAndroidKotlinConfigModelFromAssets()
		model = overrideAndroidConfigKotlinModelValues(config, keyStorePath, model)
		os.WriteFile(modelKotlinPath, model, os.ModePerm)
	}
}

func deleteFileIfExists(filePath string, errorDescription string) {
	fmt.Sprintf("%v, %v", 1, 2)
	if exists(filePath) {
		err := os.Remove(filePath)

		if err != nil {
			os.Stderr.WriteString(fmt.Sprintf("%v: %v\n", errorDescription, err.Error()))
			os.Exit(1)
		}
	}
}

func readAndroidKotlinConfigModelFromAssets() []byte {
	model, errFileNotFoundInTmp := data.Asset("lib/OneginiConfigModel.kt")
	if errFileNotFoundInTmp != nil {
		os.Stderr.WriteString(fmt.Sprintf("ERROR: Could not read config model in assets: %v\n", errFileNotFoundInTmp))
		os.Exit(1)
	}

	return model
}

func readAndroidJavaConfigModelFromAssets() []byte {
	model, errFileNotFoundInTmp := data.Asset("lib/OneginiConfigModel.java")
	if errFileNotFoundInTmp != nil {
		os.Stderr.WriteString(fmt.Sprintf("ERROR: Could not read config model in assets: %v\n", errFileNotFoundInTmp))
		os.Exit(1)
	}

	return model
}

func overrideAndroidConfigKotlinModelValues(config *Config, keystorePath string, model []byte) []byte {
	stringConfigMap := map[string]string{
		"appIdentifier":   config.Options.AppID,
		"redirectUri":     config.Options.RedirectUrl,
		"appVersion":      config.Options.AppVersion,
		"baseUrl":         config.Options.TokenServerUri,
		"resourceBaseUrl": config.Options.ResourceGatewayUris[0],
		"serverPublicKey": config.Options.ServerPublicKey.Encoded,
		"keyStoreHash":    CalculateKeystoreHash(keystorePath),
		"serverType":      config.Options.ServerType,
		"serverVersion":   config.Options.ServerVersion,
	}

	newPackage := "package " + getPackageIdentifierFromConfig(config)
	packageRe := regexp.MustCompile(`package\s.*`)
	model = packageRe.ReplaceAll(model, []byte(newPackage))

	for preference, value := range stringConfigMap {
		newPref := preference + ` = "` + value + `"`
		if preference == "serverPublicKey" && len(value) == 0 {
			newPref = preference + ` = null;`
		}

		re := regexp.MustCompile(preference + `\s=\s.*`)
		model = re.ReplaceAll(model, []byte(newPref))
	}

	re := regexp.MustCompile(`CONFIGURATOR_VERSION`)
	model = re.ReplaceAll(model, []byte(version.Version))

	return model
}

func overrideAndroidConfigJavaModelValues(config *Config, keystorePath string, model []byte) []byte {
	stringConfigMap := map[string]string{
		"appIdentifier":   config.Options.AppID,
		"redirectionUri":  config.Options.RedirectUrl,
		"appVersion":      config.Options.AppVersion,
		"baseURL":         config.Options.TokenServerUri,
		"resourceBaseURL": config.Options.ResourceGatewayUris[0],
		"serverPublicKey": config.Options.ServerPublicKey.Encoded,
		"keystoreHash":    CalculateKeystoreHash(keystorePath),
		"serverType":      config.Options.ServerType,
		"serverVersion":   config.Options.ServerVersion,
	}

	newPackage := "package " + getPackageIdentifierFromConfig(config) + ";"
	packageRe := regexp.MustCompile(`package\s.*;`)
	model = packageRe.ReplaceAll(model, []byte(newPackage))

	for preference, value := range stringConfigMap {
		newPref := preference + ` = "` + value + `";`
		if preference == "serverPublicKey" && len(value) == 0 {
			newPref = preference + ` = null;`
		}

		re := regexp.MustCompile(preference + `\s=\s.*;`)
		model = re.ReplaceAll(model, []byte(newPref))
	}

	re := regexp.MustCompile(`CONFIGURATOR_VERSION`)
	model = re.ReplaceAll(model, []byte(version.Version))

	return model
}
