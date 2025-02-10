package util

import (
	"os"
	"path"
	"strings"
)

func RemoveAndroidSecurityController(config *Config) {
	storePath := config.getAndroidSecurityControllerPath()
	_ = os.Remove(storePath)
}

func RemoveIOSSecurityController(config *Config) {
	group := "Configuration"
	xcodeProjPath := config.getIosXcodeProjPath()
	configModelPath := config.getIosConfigModelPath()
	headerStorePath := path.Join(configModelPath, "SecurityController.h")
	modelStorePath := path.Join(configModelPath, "SecurityController.m")

	removeFileFromXcodeProj(headerStorePath, xcodeProjPath, group, config.FlavorName)
	removeFileFromXcodeProj(modelStorePath, xcodeProjPath, group, config.FlavorName)
	_ = os.Remove(headerStorePath)
	_ = os.Remove(modelStorePath)

}

func (config *Config) getAndroidSecurityControllerPath() string {
	modelPath := path.Join(getPlatformSpecificAndroidClasspathPath(config), "SecurityController.java")
	// if modelPath has no package name, check namespace property in build.gradle
	if strings.HasSuffix(modelPath, "java/SecurityController.java") {
		modelPath = strings.TrimSuffix(modelPath, "SecurityController.java")
		modelPath = path.Join(modelPath, strings.ReplaceAll(config.getAndroidNamespacePath(), ".", "/"), "/SecurityController.java")
	}
	return modelPath
}
