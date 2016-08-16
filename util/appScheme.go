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
