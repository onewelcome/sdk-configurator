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

import "os"

func PrepareIosPaths(config *Config) {
	configModelPath := config.getIosConfigModelPath()
	if _, err := os.Stat(configModelPath); os.IsNotExist(err) {
		os.Mkdir(configModelPath, os.ModePerm)
	}

	certificatePath := config.getIosXcodeCertificatePath()
	if _, err := os.Stat(certificatePath); os.IsNotExist(err) {
		os.Mkdir(certificatePath, os.ModePerm)
	}
}

func exists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return true
}
