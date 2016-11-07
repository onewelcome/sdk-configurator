# SDK Configurator

The SDK Configurator configures the Onegini SDK in your application project.

It currently supports iOS, Android and Cordova projects. For Cordova it supports both the iOS and Android platforms.

## About the tool

The tool was created to help developers to setup their apps to use the Onegini Mobile SDKs. The main responsibilities of the tool is to generate config models 
and perform certificate pinning. In order to perform both operations a ZIP file containing certificates and app config is needed. The file can be downloaded 
from the Token Server's Admin Panel.


## Installing

You can download the latest binary for you platform from the [Release page](https://github.com/Onegini/sdk-configurator/releases).

### Requirements

**Android specific:**

- Only Android Studio's projects structure is supported

**iOS specific:**

- Ruby : for more info go to https://www.ruby-lang.org/en/documentation/installation/
- Xcodeproj which can be installed with $ [sudo] gem install xcodeproj. For more info go to https://github.com/CocoaPods/Xcodeproj

## Usage

Use the `--help` flag for up to date help:
```sh
onegini-sdk-configurator --help
```

### iOS example
 
Example for configuring an iOS project:
```sh
onegini-sdk-configurator ios --config ~/path/to/tokenserver-app-config.zip --app-dir ~/onegini/cordova-app/ --debugDetection=true --rootDetection=true
```

### Android Example
Example for configuring an Android project:
```sh
onegini-sdk-configurator android --config ~/path/to/tokenserver-app-config.zip --app-dir ~/onegini/cordova-app/ --debugDetection=true --rootDetection=true
```

### Cordova example
Example for configuring a Cordova Android project:
```sh
onegini-sdk-configurator android --config /path/to/tokenserver-app-config.zip --app-dir /path/to/cordova-app/ --cordova
```
Make sure you have `onegini-cordova-plugin` installed before running the configurator.

## Building from source

The Onegini Cordova plugin actually contains a hook that will automatically trigger the configurator. 

Make sure you have go installed:
```sh
brew install go
```

Read more about setting up go in the [official docs](https://golang.org/doc/install)

Install dependencies:
```sh
go get github.com/spf13/cobra
go get -u github.com/jteeuwen/go-bindata/...
```

Clone project:
```sh
mkdir -p $GOPATH/src/github.com/Onegini
git clone ssh://git@gitlab.onegini.com:10022/mobile-platform/sdk-config-importer.git $GOPATH/src/github.com/Onegini/onegini-sdk-configurator

Build project with:
```sh
go build
```

Or run without export a binary using:
```sh
go run main.go
```

Update binary assets using
```sh
go-bindata -pkg data -o data/bindata.go lib/
```