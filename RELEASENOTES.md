# Release notes

## 4.3.0

### Improvements
* Add support for Universal links in Cordova

## 4.2.2

### Improvements
* Provide a way to more easily recover apps using the SDK with the configurator, if potential future unannounced changes in iOS or Android prevent them from working correctly.

## 4.2.1

Please note that this release is required if you are using the following Cordova Plugin Onegini versions:
* Cordova Plugin Onegini 4.2.1
* Cordova Plugin Onegini 5.0.1 and higher

It is backward compatible with the previous Cordova Plugin Onegini and SDK versions.

### Bug fixes

* Fixed user registration flow with JavaScript callback on Android platform using Cordova Onegini Plugin

## 4.2.0

### Improvements

* Add support for the Cordova Android Platform 7.0.0

## 4.1.0

### Features

* Add support for NativeScript projects
* Print a comment in the config model to show which Configurator version was used to generate it

## 4.0.0

Please note that this release is required if you are using the following SDK versions:
* Android SDK 8.0.0 and higher
* iOS SDK 7.0.0 and higher

It is backwards compatible with the previous SDK versions.

### Features

* Made the SDK configurator compatible with Android SDK 8.0.0 and iOS SDK 7.0.0 

## 3.1.1

### Bug fixes

* Fixed path resolving issues when using a relative path as `app-dir`
* Auto locate the Xcode project file
* Prevent the configurator from creating multiple Xcode references wen run multiple times

## 3.1.0

### Features

* The Max PIN failures property in the ConfigModel for Android is no longer a required property.

## 3.0.0

Please note that this release is only compatible with the following SDK versions:
* Android SDK 6.00.00 and higher
* iOS SDK 5.00.00 and higher

### Features

* Add support for Android SDK versions 6.00.00 and higher
* Add support for iOS SDK versions 5.00.00 and higher
* Add a version flag
 
### Bug fixes

* Fixed a bug that forced a specific Gradle project layout for Android

## 2.0.2

### Bug fixes
* Fixed a bug that caused an error in the iOS Config model.

## 2.0.1

### Bug fixes

* The configurator now parses the value "OneginiStoreCookies" from config.xml for Cordova projects.

## 2.0.0

### Features

* Complete rebuild of the SDK configurator in go
* CLI api using flags

## 1.0.0

### Features

* Initial release