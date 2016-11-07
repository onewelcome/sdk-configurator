package util

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

func WriteAndroidAppScheme(appDir string, config *Config) {
	manifestPath := getAndroidManifestPath(appDir, config)
	manifest, err := ioutil.ReadFile(manifestPath)
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("ERROR: Cannot read the Android Manifest: %v.\n", err))
		os.Exit(1)
	}

	scheme := strings.Split(config.Options.RedirectUrl, "://")[0]

	if isCordova(config) {
		re := regexp.MustCompile(`(?s)<activity\s+.*android:name="MainActivity".*>.*<intent-filter>.*android:scheme="([^"]*)".*</intent-filter>.*</activity>`)
		rem := regexp.MustCompile(`android:scheme="[^"]*"`)
		manifest = re.ReplaceAllFunc(manifest, func(m []byte) (r []byte) {
			r = rem.ReplaceAll(m, []byte("android:scheme=\""+scheme+"\""))
			return
		})
		ioutil.WriteFile(manifestPath, manifest, os.ModePerm)
	}
}
