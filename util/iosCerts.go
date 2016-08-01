package util

import (
	"encoding/base64"
	"encoding/pem"
	"io/ioutil"
	"os"
	"path"

	"strings"
	"fmt"
	"path/filepath"
)

func ConfigureIOSCertificates(config *Config, appName string, appDir string) {
	storeDir := getIosXcodeCertificatePath(appDir, appName, config)
	xcodeProjPath := getIosXcodeProjPath(appDir, appName, config)

	removeOldCerts(storeDir, xcodeProjPath)

	for certName, certContents := range config.Certs {

		certPath := path.Join(storeDir, certName + ".cer")

		block, _ := pem.Decode([]byte(certContents))
		ioutil.WriteFile(certPath, block.Bytes, os.ModePerm)

		iosAddCertFilesToXcodeProj(certPath, xcodeProjPath, appName)
	}
}

func removeOldCerts(storeDir string, xcodeProjPath string) {
	d, err := os.Open(storeDir)
	if err != nil {
		if (os.IsNotExist(err)) {
			return
		} else {
			os.Stderr.WriteString(fmt.Sprintf("ERROR: Cannot open certificate Store dir: %v\n", err.Error()))
			os.Exit(1)
		}
	}
	defer d.Close()

	fileInfo, err := d.Readdir(-1)
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("ERROR: Cannot remove old certs: %v\n", err.Error()))
		os.Exit(1)
	}

	for _, file := range fileInfo {
		if file.Mode().IsRegular() && strings.HasSuffix(file.Name(), ".cer") {
			filePath := storeDir + string(filepath.Separator) + file.Name()
			os.Remove(filePath)
			iosRemoveCertFilesFromXcodeProj(filePath, xcodeProjPath)
		}
	}
}

func getBase64Certs(config *Config) ([]string) {
	var base64Certs []string

	for _, certContents := range config.Certs {
		block, _ := pem.Decode([]byte(certContents))

		base64Cert := base64.StdEncoding.EncodeToString(block.Bytes)
		base64Certs = append(base64Certs, base64Cert)
	}

	return base64Certs
}
