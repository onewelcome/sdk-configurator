package util

import (
	"os"
	"path"
	"strings"
)

func GetAndroidKeystorePath(appDir string, config *Config) string {
	androidPlatformPath := ""
	if isCordova(config) {
		androidPlatformPath = path.Join(appDir, "platforms", "android")
	} else {
		androidPlatformPath = path.Join(appDir, "app", "src", "main")
	}

	androidRawPath := path.Join(androidPlatformPath, "res", "raw")
	if exists(androidRawPath) == false {
		os.MkdirAll(androidRawPath, os.ModePerm)
	}

	return path.Join(androidRawPath, "keystore.bks")
}

func getAndroidSecurityControllerPath(appDir string, config *Config) string {
	if isCordova(config) {
		return path.Join(appDir, "platforms", "android", "src", path.Join(strings.Split(config.Cordova.ID, ".")...), "SecurityController.java")
	} else {
		return path.Join(appDir, "app", "src", "main", "java", path.Join(strings.Split(config.AndroidManifest.PackageID, ".")...), "SecurityController.java")
	}
}

func getAndroidManifestPath(appDir string, config *Config) string {
	if isCordova(config) {
		return path.Join(appDir, "platforms", "android", "AndroidManifest.xml")
	} else {
		return path.Join(appDir, "app", "src", "main", "AndroidManifest.xml")
	}
}

func getAndroidConfigModelPath(appDir string, config *Config) string {
	if isCordova(config) {
		return path.Join(appDir, "platforms", "android", "src", path.Join(strings.Split(config.Cordova.ID, ".")...), "OneginiConfigModel.java")
	} else {
		return path.Join(appDir, "app", "src", "main", "java", path.Join(strings.Split(config.AndroidManifest.PackageID, ".")...), "OneginiConfigModel.java")
	}
}

func PrepareIosPaths(appDir string, appName string, config *Config) {
	configModelPath := getIosConfigModelPath(appDir, appName, config)
	if _, err := os.Stat(configModelPath); os.IsNotExist(err) {
		os.Mkdir(configModelPath, os.ModePerm)
	}

	certificatePath := getIosXcodeCertificatePath(appDir, appName, config)
	if _, err := os.Stat(certificatePath); os.IsNotExist(err) {
		os.Mkdir(certificatePath, os.ModePerm)
	}
}

func getIosXcodeProjPath(appDir string, appName string, config *Config) string {
	return path.Join(getIosSrcPath(appDir, config), appName+".xcodeproj")
}

func getIosXcodeCertificatePath(appDir string, appName string, config *Config) string {
	if isCordova(config) {
		return path.Join(getIosSrcPath(appDir, config), appName, "Resources")
	} else {
		return path.Join(getIosSrcPath(appDir, config), "Resources")
	}
}

func getIosSrcPath(appDir string, config *Config) string {
	if isCordova(config) {
		return path.Join(appDir, "platforms", "ios")
	} else {
		return appDir
	}
}

func getIosConfigModelPathMFile(appDir string, appName string, config *Config) string {
	return path.Join(getIosConfigModelPath(appDir, appName, config), "OneginiConfigModel.m")
}

func getIosConfigModelPathHFile(appDir string, appName string, config *Config) string {
	return path.Join(getIosConfigModelPath(appDir, appName, config), "OneginiConfigModel.h")
}

func getIosConfigModelPath(appDir string, appName string, config *Config) string {
	if isCordova(config) {
		return path.Join(appDir, "platforms", "ios", appName, "Configuration")
	} else {
		return path.Join(appDir, "Configuration")
	}
}

func iosAddCertFilesToXcodeProj(certPath string, xcodeProjPath string, appName string) {
	addFileToXcodeProj(certPath, xcodeProjPath, appName, "Resources")
}

func iosRemoveCertFilesFromXcodeProj(certPath string, xcodeProjPath string) {
	removeFileFromXcodeProj(certPath, xcodeProjPath, "Resources")
}

func iosAddConfigModelFileToXcodeProj(modelFile string, xcodeProjPath string, appName string) {
	addFileToXcodeProj(modelFile, xcodeProjPath, appName, "Configuration")
}

func exists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return true
}
