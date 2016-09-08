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
	"encoding/base64"
	"encoding/pem"
	"io/ioutil"
	"os"
	"path"

	"fmt"
	"path/filepath"
	"strings"
)

func ConfigureIOSCertificates(config *Config) {
	storeDir := config.getIosXcodeCertificatePath()
	xcodeProjPath := config.getIosXcodeProjPath()

	removeOldCerts(storeDir, xcodeProjPath)

	for certName, certContents := range config.Certs {

		certPath := path.Join(storeDir, certName+".cer")

		block, _ := pem.Decode([]byte(certContents))
		ioutil.WriteFile(certPath, block.Bytes, os.ModePerm)

		iosAddCertFilesToXcodeProj(certPath, xcodeProjPath, config.AppTarget)
	}
}

func removeOldCerts(storeDir string, xcodeProjPath string) {
	d, err := os.Open(storeDir)
	if err != nil {
		if os.IsNotExist(err) {
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

func getBase64Certs(config *Config) []string {
	var base64Certs []string

	for cert, certContents := range config.Certs {
		if !strings.HasPrefix(certContents, "-----BEGIN CERTIFICATE-----") {
			os.Stderr.WriteString(fmt.Sprintf("ERROR: The '%v' certificate file provided in the Token Server configuration zip does not have the correct format.\n", cert))
			os.Stderr.WriteString(fmt.Sprint("ERROR: Make sure that it is a PEM encoded certificate. All cert files should start with '-----BEGIN CERTIFICATE-----'\n\n"))
			os.Exit(1)
		}

		block, _ := pem.Decode([]byte(certContents))

		base64Cert := base64.StdEncoding.EncodeToString(block.Bytes)
		base64Certs = append(base64Certs, base64Cert)
	}

	return base64Certs
}
