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
	"strings"
)

type Config struct {
	Options         *options
	Certs           map[string]string
	Cordova         cordovaConfig
	AndroidManifest androidManifest
	AppDir          string
	AppTarget       string
}

type options struct {
	MaxPinFailures      int      `json:"max_pin_failures"`
	TokenServerUri      string   `json:"token_server_uri"`
	AppID               string   `json:"application_identifier"`
	AppPlatform         string   `json:"application_platform"`
	AppVersion          string   `json:"application_version"`
	ResourceGatewayUris []string `json:"resource_gateway_uri"`
	RedirectUrl         string   `json:"redirect_url"`
}

type cordovaPreference struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

type cordovaConfig struct {
	Preferences []cordovaPreference `xml:"preference"`
	ID          string              `xml:"id,attr"`
	AppName     string              `xml:"name"`
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

	config.AppDir = appDir
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

func ParseAndroidManifest(config *Config) {
	values := androidManifest{}

	manifestXml, err := ioutil.ReadFile(path.Join(config.AppDir, "app", "src", "main", "AndroidManifest.xml"))
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

func parseTsZip(path string, config *Config) {
	readCloser, err := zip.OpenReader(path)
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("ERROR: could not read Token Server configuration zip: %v\n", err.Error()))
		os.Exit(1)
	}

	defer readCloser.Close()

	for _, file := range readCloser.File {
		readCloser, err := file.Open()
		if err != nil {
			os.Stderr.WriteString(fmt.Sprintf("ERROR: could not contents of Token Server configuration zip: %v\n", err.Error()))
			os.Exit(1)
		}

		if file.Name == "config.json" {
			config.Options, _ = parseJson(readCloser)
			// Don't use the filepath.Separator in the statement below because the filename always contains the forward / regardless of the
			// platform the configurator is run on
		} else if strings.HasPrefix(file.Name, "certificates/") {
			config.Certs[strings.Replace(file.Name, "certificates"+string(filepath.Separator), "", -1)] = readCert(readCloser)
		}
	}
	VerifyTsZipContents(config)
}

func parseJson(reader io.Reader) (v *options, err error) {
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
	if isCordova(config) {
		return config.Cordova.ID
	} else {
		return config.AndroidManifest.PackageID
	}
}

func isCordova(config *Config) bool {
	var cordovaConfig = &config.Cordova
	return len(cordovaConfig.ID) > 0
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

func (config *Config) getAndroidKeystorePath() string {
	androidPlatformPath := ""
	if isCordova(config) {
		androidPlatformPath = path.Join(config.AppDir, "platforms", "android")
	} else {
		androidPlatformPath = path.Join(config.AppDir, config.AppTarget, "src", "main")
	}

	androidRawPath := path.Join(androidPlatformPath, "res", "raw")
	if exists(androidRawPath) == false {
		os.MkdirAll(androidRawPath, os.ModePerm)
	}

	return path.Join(androidRawPath, "keystore.bks")
}

func (config *Config) getAndroidSecurityControllerPath() string {
	if isCordova(config) {
		return path.Join(config.AppDir, "platforms", "android", "src", path.Join(strings.Split(config.Cordova.ID, ".")...), "SecurityController.java")
	} else {
		return path.Join(config.AppDir, config.AppTarget, "src", "main", "java", path.Join(strings.Split(config.AndroidManifest.PackageID, ".")...), "SecurityController.java")
	}
}

func (config *Config) getAndroidManifestPath() string {
	if isCordova(config) {
		return path.Join(config.AppDir, "platforms", "android", "AndroidManifest.xml")
	} else {
		return path.Join(config.AppDir, config.AppTarget, "src", "main", "AndroidManifest.xml")
	}
}

func (config *Config) getAndroidConfigModelPath() string {
	if isCordova(config) {
		return path.Join(config.AppDir, "platforms", "android", "src", path.Join(strings.Split(config.Cordova.ID, ".")...), "OneginiConfigModel.java")
	} else {
		return path.Join(config.AppDir, config.AppTarget, "src", "main", "java", path.Join(strings.Split(config.AndroidManifest.PackageID, ".")...), "OneginiConfigModel.java")
	}
}

func (config *Config) getIosXcodeProjPath() string {
	files, err := filepath.Glob(path.Join(config.getIosSrcPath(), "*.xcodeproj"))

	if err != nil || len(files) == 0 {
		os.Stderr.WriteString(fmt.Sprintf("ERROR: Could not find an Xcode project directory (.xcodeproj). Are you sure that '%v' contains one?\n", config.getIosSrcPath()))
		os.Exit(1)
	}

	if len(files) > 1 {
		os.Stderr.WriteString(fmt.Sprint("ERROR: Found multiple Xcode project directories (.xcodeproj). The SDK configurator currently only support a " +
			"single xcodeproj directory."))
		os.Exit(1)
	}

	return files[0]
}

func (config *Config) getIosSrcPath() string {
	if isCordova(config) {
		return path.Join(config.AppDir, "platforms", "ios")
	} else {
		return config.AppDir
	}
}

func (config *Config) getIosConfigModelPath() string {
	if isCordova(config) {
		return path.Join(config.AppDir, "platforms", "ios", config.AppTarget, "Configuration")
	} else {
		return path.Join(config.AppDir, "Configuration")
	}
}

func (config *Config) getIosXcodeCertificatePath() string {
	if isCordova(config) {
		return path.Join(config.getIosSrcPath(), config.AppTarget, "Resources")
	} else {
		return path.Join(config.getIosSrcPath(), "Resources")
	}
}

func (config *Config) getIosConfigModelPathMFile() string {
	return path.Join(config.getIosConfigModelPath(), "OneginiConfigModel.m")
}

func (config *Config) getIosConfigModelPathHFile() string {
	return path.Join(config.getIosConfigModelPath(), "OneginiConfigModel.h")
}
