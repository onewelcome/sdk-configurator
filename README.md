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

## Assumptions

Please read the following assumptions **carefully** if you start to use the SDK configurator for an **existing** project in which the Onegini SDK is already 
configured.

####Android

- **Config model:** The configurator tries to look for an existing config model class with the following name `OneginiConfigModel`. The location in which the 
SDK configurator searches is the package that is mentioned in your `AndroidManifest.xml`. The package can be found in the `package` attribute of the 
`<manifest>` element. You must remove the existing config model if you have named it differently or it it's placed in a different location before running the 
SDK configurator.

####iOS

- **Config model:** The configurator tries to look for an existing config model class in the `Configuration` group in the root of your Xcode project. You must 
remove the existing config model if it's located in a different group before running the SDK configurator.
- **Certificates:** The configurator will remove any existing certificates located in the `Resources` group in the root of your Xcode project. You must remove 
any certificates located in another location before running the SDK configurator

## Usage

Use the `--help` flag for up to date help:
```sh
./onegini-sdk-configurator --help
```

### iOS example
 
Example for configuring an iOS project:
```sh
./onegini-sdk-configurator ios --config ~/path/to/tokenserver-app-config.zip --app-dir ~/onegini/ios-app/ --target-name myTarget --debugDetection=true --rootDetection=true
```

Replace the `myTarget` value with the application target located in your Xcode project. See the [Apple documentation](https://developer.apple.com/library/ios/documentation/IDEs/Conceptual/AppDistributionGuide/ConfiguringYourApp/ConfiguringYourApp.html) 
for more information on the app target.

### Android Example
Example for configuring an Android project:
```sh
./onegini-sdk-configurator android --config ~/path/to/tokenserver-app-config.zip --app-dir ~/onegini/android-app/ --debugDetection=true --rootDetection=true
```

### Cordova example
The Onegini Cordova plugin contains a hook that will automatically trigger the configurator when you run `cordova platform add`.
You can still choose to run the configurator manually (e.g for updating an existing platform).
 
Example for configuring a Cordova Android project manually:
```sh
./onegini-sdk-configurator android --config /path/to/tokenserver-app-config.zip --app-dir /path/to/cordova-app/ --cordova
```

Make sure you have the `onegini-cordova-plugin` installed before running the configurator. You will need to rerun the configurator for each installed platform 
in your Cordova project.