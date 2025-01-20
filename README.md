# SDK Configurator

The SDK Configurator configures the Onegini SDK in your application project.

It currently supports iOS, Android, Cordova and NativeScript projects. For Cordova & NativeScript it supports both the iOS and Android platforms.

## About the tool

The tool was created to help developers to setup their apps to use the Onegini Mobile SDKs. The main responsibilities of the tool is to generate config models 
and perform certificate pinning. In order to perform both operations a ZIP file containing certificates and app config is needed. The file can be downloaded 
from the Token Server's Admin Panel.

## Installing

You can download the latest compatible binary for your platform and SDK version from the [Release page](https://github.com/Onegini/sdk-configurator/releases). 
Please have a look at the SDK compatibility in the next paragraph to choose the correct SDK configurator version.

### SDK compatibility

The configurator binaries differ per SDK version. Please see the compatibility matrix below to choose the correct SDK configurator version:

**Android SDK**

| SDK version           | Configurator version   |
|-----------------------|------------------------|
| Android SDK >= 12.3.0 | 6.x
| Android SDK 11.x & 12.x | 5.x                    |
| Android SDK 8.x & 9.x & 10.x | 4.x                    |
| Android SDK 6.x & 7.x | 3.x                    |
| Android SDK < 6.00.00 | 2.x                    |
 
**iOS SDK**

| SDK version                | Configurator version   |
|----------------------------|------------------------|
| iOS SDK >= 12.4.0 | 6.x
| iOS SDK 10.x                 | 5.x                    |
| iOS SDK 7.x & 8.x & 9.x    | 4.x                    |
| iOS SDK 5.x & 6.x          | 3.x                    |
| iOS SDK < 5.00.00          | 2.x                    |

**Cordova plugin**

| Plugin version                 | Configurator version   |
|--------------------------------|------------------------|
| Cordova plugin 9.x           | 5.0                  |
| Cordova plugin 6.x & 7.x & 8.x           | 4.2.1                  |
| Cordova plugin 5.0.x           | 4.2.1                  |
| Cordova plugin 5.0.0           | 4.2.0                  |
| Cordova plugin 4.3.x           | 4.2.0                  |
| Cordova plugin 4.2.1 - 4.2.x   | 4.2.1                  |
| Cordova plugin 4.0.0 - 4.2.0   | 4.2.0                  |
| Cordova plugin 3.x & 2.x       | 3.x                    |
| Cordova plugin < 2.0.0         | 2.x                    |

**NativeScript plugin**

| Plugin version             | Configurator version   |
|----------------------------|------------------------|
| NativeScript plugin 1.x    | 4.1.x                  |

### Requirements

**Android specific:**

- Only Android Studio's projects structure is supported

**iOS specific:**

- Ruby : for more info go to https://www.ruby-lang.org/en/documentation/installation/
- Xcodeproj which can be installed with $ [sudo] gem install xcodeproj. For more info go to https://github.com/CocoaPods/Xcodeproj

## Assumptions

Please read the following assumptions **carefully** if you wish to use the SDK configurator for an **existing** project in which the Onegini SDK is already 
configured.

#### Android

- **Config model:** The configurator tries to look for an existing config model class with the following name `OneginiConfigModel`. The location in which the 
SDK configurator searches is the package that is mentioned in your `AndroidManifest.xml`. The package can be found in the `package` attribute of the 
`<manifest>` element. If it cannot be found, as a second step it looks for `namespace` attribute in the gradle file.
You must remove the existing config model if you have named it differently or if it is placed in a different location before running the 
SDK configurator.

#### iOS

- **Config model:** The configurator tries to look for an existing config model class in the `Configuration` group in the root of your Xcode project. You must 
remove the existing config model if it is located in a different group before running the SDK configurator.
- **Certificates:** [Applies only to Configurator versions up to and including 4.x.x] The configurator will remove any existing certificates located in the `Resources` group in the root of your Xcode project. You must remove 
any certificates located in another location before running the SDK configurator.

## Usage

Use the `--help` flag for up to date help:
```sh
./sdk-configurator --help
```

### iOS example
 
Example for configuring an iOS project:
```sh
./sdk-configurator ios --config ~/path/to/tokenserver-app-config.zip --app-dir ~/path/to/ios-app/ --target-name myTarget -f mySubfolder
```

Replace the `myTarget` value with the application target located in your Xcode project. See the [Apple documentation](https://developer.apple.com/library/ios/documentation/IDEs/Conceptual/AppDistributionGuide/ConfiguringYourApp/ConfiguringYourApp.html) for more information on the app target.

Optionally use and replace `mySubfolder` for `-f` flag with proper subfolder name which is useful for many targets with different configurations each.

### Android Example
Example for configuring an Android project:
```sh
./sdk-configurator android --config ~/path/to/tokenserver-app-config.zip --module-name app --app-dir ~/path/to/android-app/ 
```

Replace the `app` value with the name of the Gradle module that contains your application sources. See the [Android documentation](https://developer.android.com/studio/projects/index.html) for more info.

### Cordova example
The Onegini Cordova plugin contains a hook that will automatically trigger the configurator when you run `cordova platform add`. You can still choose to run the configurator manually (e.g. for updating an existing platform).

Example for configuring a Cordova Android project manually:
```sh
./sdk-configurator android --config /path/to/tokenserver-app-config.zip --app-dir /path/to/cordova-app/ --cordova
```

### NativeScript example
The Onegini NativeScript plugin contains a hook that will automatically trigger the configurator when you run `tns prepare`. You can still choose to run the configurator manually (e.g. for updating an existing platform).

Example for configuring a NativeScript Android project manually:
```sh
./sdk-configurator android --config /path/to/tokenserver-app-config.zip --app-dir /path/to/cordova-app/ --nativescript
```

Make sure you have the `nativescript-onegini` plugin installed before running the configurator. You will need to rerun the configurator for each installed platform 
in your NativeScript project.
