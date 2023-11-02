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
	"os/exec"
	"path"

	"crypto/sha256"
	"fmt"

	"crypto/rand"
	"encoding/base64"

	"path/filepath"

	"github.com/onewelcome/sdk-configurator/data"
)

func CreateKeystore(config *Config) {
	storePath := config.getAndroidKeystorePath()

	_, err := os.Stat(storePath)
	if err == nil {
		os.Remove(storePath)
	}

	keystorePassword := generateKeystorePassword(2048)
	bcprovPath := restoreBcprov()
	keytoolPath := findKeytool()

	for certName, certContents := range config.Certs {
		cmdKeytool := exec.Command(
			keytoolPath,
			"-import",
			"-alias", certName,
			"-keystore", storePath,
			"-storepass", keystorePassword,
			"-providerpath", bcprovPath,
			"-storetype", "BKS",
			"-provider", "org.bouncycastle.jce.provider.BouncyCastleProvider",
			"-noprompt",
		)
		cmdStdinPipe, _ := cmdKeytool.StdinPipe()
		cmdStdinPipe.Write([]byte(certContents))
		cmdStdinPipe.Close()

		_, err := cmdKeytool.CombinedOutput()
		if err != nil {
			os.Stderr.WriteString(fmt.Sprintf("ERROR: could not execute keytool: %v\n", err.Error()))
			os.Exit(1)
		}
	}
}

func findKeytool() (keyToolPath string) {
	keyToolPath, lookErr := exec.LookPath("keytool")
	if lookErr != nil {
		keyToolInJavaHome := path.Join(os.Getenv("JAVA_HOME"), "bin")
		keyToolPath, _ = exec.LookPath(keyToolInJavaHome + string(filepath.Separator) + "keytool")
	}

	if _, err := os.Stat(keyToolPath); err != nil {
		os.Stderr.WriteString(fmt.Sprintln("ERROR: Could not find keytool utility in your $PATH or $JAVA_HOME/bin.\n\nSee https://docs.oracle.com/cd/E19182-01/820-7851/inst_cli_jdk_javahome_t for instructions on how to set $JAVA_HOME"))
		os.Exit(1)
	}
	return keyToolPath
}

func CalculateKeystoreHash(keystorePath string) (hash string) {
	keystore, err := ioutil.ReadFile(keystorePath)
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("ERROR: Could not calculate keystore hash: %v\n", err))
		os.Exit(1)
	}
	rawHash := sha256.Sum256(keystore)
	hash = fmt.Sprintf("%x", string(rawHash[:]))

	return
}

func generateKeystorePassword(n int) string {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("ERROR: Could not generate random data: %v\n", err.Error()))
		os.Exit(1)
	}

	return base64.URLEncoding.EncodeToString(b)
}

func restoreBcprov() (filePath string) {
	tempPath := path.Join(os.TempDir(), "sdk-configurator")
	filePath = path.Join(tempPath, "lib", "bcprov-jdk15on-1.46.jar")

	if err := data.RestoreAsset(tempPath, "lib/bcprov-jdk15on-1.46.jar"); err != nil {
		os.Stderr.WriteString(fmt.Sprintf("ERROR: Could not restore required asset: %v\n", err))
		os.Exit(1)
	}

	return
}
