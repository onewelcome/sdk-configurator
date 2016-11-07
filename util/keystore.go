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

	"github.com/Onegini/onegini-sdk-configurator/data"
)

func CreateKeystore(config *Config, storePath string) {
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

func findKeytool() (keytoolPath string) {
	keytoolPath, lookErr := exec.LookPath("keytool")
	if lookErr != nil {
		keyToolInJavaHome := path.Join(os.Getenv("JAVA_HOME"), "bin")
		keytoolPath, _ = exec.LookPath(keyToolInJavaHome + string(filepath.Separator) + "keytool")
	}

	if _, err := os.Stat(keytoolPath); err != nil {
		os.Stderr.WriteString(fmt.Sprintln("ERROR: Could not find keytool utility in your $PATH or $JAVA_HOME/bin.\n\nSee https://docs.oracle.com/cd/E19182-01/820-7851/inst_cli_jdk_javahome_t for istructions on how to set $JAVA_HOME"))
		os.Exit(1)
	}
	return keytoolPath
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
	tempPath := path.Join(os.TempDir(), "onegini-sdk-configurator")
	filePath = path.Join(tempPath, "lib", "bcprov-jdk15on-1.46.jar")

	if err := data.RestoreAsset(tempPath, "lib/bcprov-jdk15on-1.46.jar"); err != nil {
		os.Stderr.WriteString(fmt.Sprintf("ERROR: Could not restore required asset: %v\n", err))
		os.Exit(1)
	}

	return
}
