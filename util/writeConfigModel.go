package util

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"

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
	var storeCookiesStr string

	if shouldStoreCookies(config) == true {
		storeCookiesStr = "YES"
	} else {
		storeCookiesStr = "NO"
	}

	stringConfigMap := map[string]string{
		"kOGAppIdentifier":   config.Options.AppID,
		"kOGAppScheme":       strings.Split(config.Options.RedirectUrl, "://")[0],
		"kOGAppVersion":      config.Options.AppVersion,
		"kOGAppBaseURL":      config.Options.TokenServerUri,
		"kOGMaxPinFailures":  strconv.Itoa(config.Options.MaxPinFailures),
		"kOGResourceBaseURL": config.Options.ResourceGatewayUris[0],
		"kOGRedirectURL":     config.Options.RedirectUrl,
	}
	nonStringConfigMap := map[string]string{
		"kOGStoreCookies": storeCookiesStr,
	}

	for preference, value := range stringConfigMap {
		newPref := `@"` + preference + `" : @"` + value + `"`
		re := regexp.MustCompile(`@"` + preference + `"\s*:\s*@".*"`)
		model = re.ReplaceAll(model, []byte(newPref))
	}

	for preference, value := range nonStringConfigMap {
		newPref := `@"` + preference + `" : @(` + value + `)`
		re := regexp.MustCompile(`@"` + preference + `"\s*:\s*@(.*)`)
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
		"appScheme":       strings.Split(config.Options.RedirectUrl, "://")[0],
		"appVersion":      config.Options.AppVersion,
		"baseURL":         config.Options.TokenServerUri,
		"resourceBaseURL": config.Options.ResourceGatewayUris[0],
		"keystoreHash":    CalculateKeystoreHash(keystorePath),
	}
	nonStringConfigMap := map[string]string{
		"maxPinFailures": strconv.Itoa(config.Options.MaxPinFailures),
		"storeCookies":   strconv.FormatBool(shouldStoreCookies(config)),
	}

	newPackage := "package " + getPackageIdentifierFromConfig(config) + ";"
	packageRe := regexp.MustCompile(`package\s.*;`)
	model = packageRe.ReplaceAll(model, []byte(newPackage))

	for preference, value := range stringConfigMap {
		newPref := preference + ` = "` + value + `";`
		re := regexp.MustCompile(preference + `\s=\s".*";`)
		model = re.ReplaceAll(model, []byte(newPref))
	}

	for preference, value := range nonStringConfigMap {
		newPref := preference + ` = ` + value + `;`
		re := regexp.MustCompile(preference + `\s=\s.*;`)
		model = re.ReplaceAll(model, []byte(newPref))
	}

	return model
}

func shouldStoreCookies(config *Config) bool {
	storeCookies := true

	if config.Cordova.Preferences == nil {
		return storeCookies
	}

	for _, pref := range config.Cordova.Preferences {
		if pref.Name == "OneginiStoreCookies" {
			var err error
			storeCookies, err = strconv.ParseBool(pref.Value)
			if err != nil {
				os.Stderr.WriteString(fmt.Sprintf("ERROR: could not parse 'OneginiStoreCookies' preference: %v\n", err.Error()))
				os.Exit(1)
			}
		}
	}

	return storeCookies
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
