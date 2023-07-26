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
	"archive/zip"
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

type Config struct {
	Options                  *options
	Certs                    map[string]string
	Cordova                  cordovaConfig
	NativeScript             nativeScriptConfig
	AndroidManifest          androidManifest
	AppDir                   string
	AppTarget                string
	FlavorName               string
	ConfigureForCordova      bool
	ConfigureForNativeScript bool
}

type options struct {
	MaxPinFailures      int             `json:"max_pin_failures"`
	TokenServerUri      string          `json:"token_server_uri"`
	AppID               string          `json:"application_identifier"`
	AppPlatform         string          `json:"application_platform"`
	AppVersion          string          `json:"application_version"`
	ResourceGatewayUris []string        `json:"resource_gateway_uri"`
	RedirectUrl         string          `json:"redirect_url"`
	ServerPublicKey     serverPublicKey `json:"server_public_key"`
}

type serverPublicKey struct {
	Encoded   string `json:"encoded"`
	Algorithm string `json:"algorithm"`
}

type cordovaPreference struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

type cordovaConfig struct {
	Preferences []cordovaPreference `xml:"preference"`
	AppName     string              `xml:"name"`
}

type nativeScriptConfig struct {
	NS                 NS                 `json:"nativescript"`
	OneginiPreferences OneginiPreferences `json:"onegini"`
}

type NS struct {
	ID string `json:"id"`
}

type OneginiPreferences struct {
	RootDetectionEnabled  *bool `json:"root-detection-enabled,omitempty"`
	DebugDetectionEnabled *bool `json:"debug-detection-enabled,omitempty"`
	DebugLogsEnabled      bool  `json:"debug-logs-enabled"`
}

type androidManifest struct {
	PackageID string `xml:"package,attr"`
}

func ParseConfig(appDir string, configPath string) (config *Config) {
	config = new(Config)
	config.Certs = make(map[string]string)

	if len(configPath) == 0 {
		fmt.Print("ERROR: No Token Server configuration provided. Provide one using 'onegini-sdk-configurator <platform> -c <config-zip-location>'\n\n")
		fmt.Print("execute 'onegini-sdk-configurator --help' to see how to use the configurator\n")
		os.Exit(1)
	}

	config.AppDir = config.resolveAppDirPath(appDir)
	parseTsZip(configPath, config)

	return
}

func ParseCordovaConfig(config *Config) {
	values := cordovaConfig{}

	configXml, err := ioutil.ReadFile(path.Join(config.AppDir, "config.xml"))
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("ERROR: Cannot read the Cordova config.xml: %v\n", err.Error()))
		os.Exit(1)
	}

	err = xml.Unmarshal(configXml, &values)
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("ERROR: Cannot read the Cordova config.xml: %v\n", err.Error()))
		os.Exit(1)
	}

	config.Cordova = values
}

func ParseNativeScriptConfig(config *Config) {
	jsonFile, err := os.Open("package.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("ERROR: Cannot read the NativeScript package.json: %v\n", err.Error()))
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var nsConfig nativeScriptConfig
	json.Unmarshal(byteValue, &nsConfig)

	config.NativeScript = nsConfig
}

func ParseAndroidManifest(config *Config) {
	values := androidManifest{}

	manifestPath := config.getAndroidManifestPath()
	manifestXml, err := ioutil.ReadFile(manifestPath)
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("ERROR: Cannot read the Android Manifest: %v\n", err.Error()))
		os.Exit(1)
	}

	err = xml.Unmarshal(manifestXml, &values)
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("ERROR: Cannot read the Android Manifest: %v\n", err.Error()))
		os.Exit(1)
	}

	config.AndroidManifest = values
}

func SetAppTarget(appTarget string, config *Config) {
	config.AppTarget = appTarget
}

func SetFlavorName(flavorName string, config *Config) {
	config.FlavorName = flavorName
}

func parseTsZip(path string, config *Config) {
	readCloser, err := zip.OpenReader(path)
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("ERROR: could not read Token Server configuration zip: %v\n", err.Error()))
		os.Exit(1)
	}

	defer readCloser.Close()

	for _, file := range readCloser.File {
		openedFile, err := file.Open()
		if err != nil {
			os.Stderr.WriteString(fmt.Sprintf("ERROR: could not read the contents of Token Server configuration zip: %v\n", err.Error()))
			os.Exit(1)
		}

		if file.Name == "config.json" {
			config.Options, _ = parseTsJson(openedFile)
			// Don't use the filepath.Separator in the statement below because the filename always contains the forward / regardless of the
			// platform the configurator is run on
		} else if strings.HasPrefix(file.Name, "certificates/") {
			config.Certs[strings.Replace(file.Name, "certificates"+string(filepath.Separator), "", -1)] = readCert(openedFile)
		}
	}
	VerifyTsZipContents(config)
}

func parseTsJson(reader io.Reader) (v *options, err error) {
	v = new(options)
	err = json.NewDecoder(reader).Decode(v)
	return
}

func readCert(reader io.Reader) (contents string) {
	buffer := new(bytes.Buffer)
	buffer.ReadFrom(reader)
	contents = buffer.String()
	return
}

func getPackageIdentifierFromConfig(config *Config) string {
	if config.AndroidManifest.PackageID != "" {
		fmt.Println("PKPK I am inside empty package id", config.getAndroidNamespacePath())

		return config.AndroidManifest.PackageID
	} else {
		fmt.Println("PKPK NAMESPACE!!!", config.getAndroidNamespacePath())
		return config.getAndroidNamespacePath()
	}
}

func VerifyTsZipContents(config *Config) {
	if config.Options == nil {
		os.Stderr.WriteString(fmt.Sprintln("ERROR: The provided configuration zip does not contain the required information. Is the supplied archive a valid Token " +
			"Server configuration zip?"))
		os.Exit(1)
	}

	if config.Options.ResourceGatewayUris == nil || len(config.Options.ResourceGatewayUris) == 0 {
		os.Stderr.WriteString(fmt.Sprint("ERROR: No resource gateway URI is specified in the configuration zip. Please check the Token Server configuration.\n" +
			"See the following link for more info: https://docs.onegini.com/public/token-server/topics/general-app-config/resource-gateway/resource-gateway.html\n"))
		os.Exit(1)
	}

	if config.Certs == nil || len(config.Certs) == 0 {
		os.Stderr.WriteString(fmt.Sprintln("ERROR: Does the Token Server configuration zip contain certificates?"))
		os.Exit(1)
	}
}

func (config *Config) resolveAppDirPath(appDir string) string {
	absAppDirPath, err := filepath.Abs(appDir)

	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("ERROR: Could resolve App dir '%v' into absolute path", config.AppDir))
		os.Exit(1)
	}
	return absAppDirPath
}

// Android Paths
func getCordovaAndroidPlatformPath(config *Config) string {
	var android7PlatformPath = getCordovaAndroid7PlatformPath(config)
	if exists(android7PlatformPath) {
		return android7PlatformPath
	}
	return getCordovaAndroid6PlatformPath(config)
}

func getCordovaAndroid7PlatformPath(config *Config) string {
	return path.Join(config.AppDir, "platforms", "android", "app", "src", "main")
}

func getCordovaAndroid6PlatformPath(config *Config) string {
	return path.Join(config.AppDir, "platforms", "android")
}

func getCordovaAndroidClasspath(config *Config) string {
	var cordovaAndroid7Classpath = getCordovaAndroid7Classpath(config)
	if exists(cordovaAndroid7Classpath) {
		return cordovaAndroid7Classpath
	}
	return getCordovaAndroid6Classpath(config)
}

func getCordovaAndroid7Classpath(config *Config) string {
	return path.Join(getCordovaAndroidPlatformPath(config), "java", path.Join(strings.Split(config.AndroidManifest.PackageID, ".")...))
}

func getCordovaAndroid6Classpath(config *Config) string {
	return path.Join(getCordovaAndroidPlatformPath(config), "src", path.Join(strings.Split(config.AndroidManifest.PackageID, ".")...))
}

func getNativeScriptAndroidPlatformPath(config *Config) string {
	return path.Join(config.AppDir, "platforms", "android", "app", "src", "main")
}

func getNativeScriptAndroidClasspath(config *Config) string {
	return path.Join(getNativeScriptAndroidPlatformPath(config), "java", path.Join(strings.Split(config.AndroidManifest.PackageID, ".")...))
}

func getDefaultAndroidPlatformPath(config *Config, useFlavor bool) string {
	srcPath := path.Join(config.AppDir, config.AppTarget, "src")
	if useFlavor && len(config.FlavorName) > 0 {
		return path.Join(srcPath, config.FlavorName)
	} else {
		return path.Join(srcPath, "main")
	}
}

func getDefaultAndroidClasspath(config *Config, useFlavor bool) string {
	return path.Join(getDefaultAndroidPlatformPath(config, useFlavor), "java", path.Join(strings.Split(config.AndroidManifest.PackageID, ".")...))
}

func getPlatformSpecificAndroidPlatformPath(config *Config, useFlavor bool) string {
	androidPlatformPath := ""
	if config.ConfigureForCordova {
		androidPlatformPath = getCordovaAndroidPlatformPath(config)
	} else if config.ConfigureForNativeScript {
		androidPlatformPath = getNativeScriptAndroidPlatformPath(config)
	} else {
		androidPlatformPath = getDefaultAndroidPlatformPath(config, useFlavor)
	}

	return androidPlatformPath
}

func getPlatformSpecificAndroidClasspathPath(config *Config) string {
	androidClasspathPath := ""
	if config.ConfigureForCordova {
		androidClasspathPath = getCordovaAndroidClasspath(config)
	} else if config.ConfigureForNativeScript {
		androidClasspathPath = getNativeScriptAndroidClasspath(config)
	} else {
		androidClasspathPath = getDefaultAndroidClasspath(config, true)
	}

	return androidClasspathPath
}

func (config *Config) getAndroidKeystorePath() string {
	androidRawPath := path.Join(getPlatformSpecificAndroidPlatformPath(config, true), "res", "raw")
	if exists(androidRawPath) == false {
		os.MkdirAll(androidRawPath, os.ModePerm)
	}

	return path.Join(androidRawPath, "keystore.bks")
}

func (config *Config) getAndroidManifestPath() string {
	return path.Join(getPlatformSpecificAndroidPlatformPath(config, false), "AndroidManifest.xml")
}

func (config *Config) getAndroidConfigModelPath() string {
	modelPath := path.Join(getPlatformSpecificAndroidClasspathPath(config), "OneginiConfigModel.java")
	// if package name wasn't found in AndroidManifest.xml file, check namespace property in build.gradle
	if strings.HasSuffix(modelPath, "java/OneginiConfigModel.java") {
		modelPath = strings.TrimSuffix(modelPath, "OneginiConfigModel.java")
		modelPath = path.Join(modelPath, strings.ReplaceAll(config.getAndroidNamespacePath(), ".", "/"), "/OneginiConfigModel.java")
	}
	return modelPath
}

func (config *Config) getAndroidSecurityControllerPath() string {
	modelPath := path.Join(getPlatformSpecificAndroidClasspathPath(config), "SecurityController.java")
	if strings.HasSuffix(modelPath, "java/SecurityController.java") {
		modelPath = strings.TrimSuffix(modelPath, "SecurityController.java")
		modelPath = path.Join(modelPath, strings.ReplaceAll(config.getAndroidNamespacePath(), ".", "/"), "/SecurityController.java")
	}
	return modelPath
}

func (config *Config) getAndroidClasspathPath() string {
	return path.Join(getPlatformSpecificAndroidClasspathPath(config))
}

func (config *Config) getAndroidNamespacePath() string {
	gradleFilePath := path.Join(config.AppDir, config.AppTarget, "build.gradle")
	gradleContent, err := ioutil.ReadFile(gradleFilePath)
	if err != nil {
		fmt.Println("Error during reading gradle file", err)
	}
	namespaceRegexMatches := regexp.MustCompile(`(?m)^\s*namespace\s+'([^']+)'\s*$`).FindStringSubmatch(string(gradleContent))
	if len(namespaceRegexMatches) == 2 {
		return namespaceRegexMatches[1]
	} else {
		fmt.Println("Namespace property not found in build.gradle file")
		return ""
	}
}

// iOS Paths

func getCordovaIosProjPath(config *Config) string {
	return path.Join(config.AppDir, "platforms", "ios")
}

func getCordovaIosSrcPath(config *Config) string {
	return path.Join(getCordovaIosProjPath(config), config.AppTarget)
}

func getNativeScriptIosProjPath(config *Config) string {
	return path.Join(config.AppDir, "platforms", "ios")
}

func getNativeScriptIosSrcPath(config *Config) string {
	return path.Join(getNativeScriptIosProjPath(config), config.AppTarget)
}

func getNativeIosProjPath(config *Config) string {
	return config.AppDir
}

func getPlatformSpecificIosProjPath(config *Config) string {
	if config.ConfigureForCordova {
		return getCordovaIosProjPath(config)
	} else if config.ConfigureForNativeScript {
		return getNativeScriptIosProjPath(config)
	} else {
		return getNativeIosProjPath(config)
	}
}

func getPlatformSpecificIosSrcPath(config *Config) string {
	if config.ConfigureForCordova {
		return getCordovaIosSrcPath(config)
	} else if config.ConfigureForNativeScript {
		return getNativeScriptIosSrcPath(config)
	} else {
		return getNativeIosProjPath(config)
	}
}

func (config *Config) getIosXcodeProjPath() string {
	files, err := filepath.Glob(path.Join(getPlatformSpecificIosProjPath(config), "*.xcodeproj"))

	if err != nil || len(files) == 0 {
		os.Stderr.WriteString(fmt.Sprintf("ERROR: Could not find an Xcode project directory (.xcodeproj). Are you sure that '%v' contains one?\n", getPlatformSpecificIosProjPath(config)))
		os.Exit(1)
	}

	if len(files) > 1 {
		os.Stderr.WriteString(fmt.Sprint("ERROR: Found multiple Xcode project directories (.xcodeproj). The SDK configurator currently only support a " +
			"single xcodeproj directory."))
		os.Exit(1)
	}

	return files[0]
}

func (config *Config) getIosConfigModelPath() string {
	subfolder := config.FlavorName
	srcPath := path.Join(getPlatformSpecificIosSrcPath(config), "Configuration")
	if len(subfolder) > 0 {
		return path.Join(srcPath, subfolder)
	} else {
		return srcPath
	}
}

func (config *Config) getIosXcodeCertificatePath() string {
	// Certs need to be stored in a different place NS a full build will override the resources in the src path.
	if config.ConfigureForNativeScript {
		return path.Join(config.AppDir, "app", "App_Resources", "iOS")
	} else {
		return path.Join(getPlatformSpecificIosSrcPath(config), "Resources")
	}
}

func (config *Config) getIosConfigModelPathMFile() string {
	return path.Join(config.getIosConfigModelPath(), "OneginiConfigModel.m")
}

func (config *Config) getIosConfigModelPathHFile() string {
	return path.Join(config.getIosConfigModelPath(), "OneginiConfigModel.h")
}
