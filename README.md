# Onegini Mobile SDK Configurator

## About the tool

The tool was created to help developers to setup their apps to use Onegini Mobile SDKs. The main responsibilities of the tool is to generate config models and perform certificate pinning. In order to perform both operations a ZIP file containing certificates and app config is needed. The file can be downloaded from the Token Server's Admin Panel.
Tool can be used with Android and iOS projects.

## Requirements

- bash : for Windows bash can be installed with Cygwin. More info http://cygwin.com/install.html
- jq : JSON processor. Can be downloaded from: https://stedolan.github.io/jq/download/

#### Android specific:
- Only Android Studio's projects structure is supported

#### iOS specific:
- Ruby : for more info go to https://www.ruby-lang.org/en/documentation/installation/
- Xcodeproj which can be installed with $ [sudo] gem install xcodeproj. For more info go to https://github.com/CocoaPods/Xcodeproj

## How it works

Tool is a bash script. After execution, script will ask user for all required data one by one:
1. Choose between Android and iOS platform.
2. Provide path to ZIP file containing app's config and certificates.
The next steps are platform specific.

### Android
3. Provide path to a main "src" directory of Android project

The ZIP file is being unzipped and the certificates are stored in a new secure AndroidKeystore. During this process the user has to confirm, that he trusts shown certificates. After that the keystore is being copied to the project's `/res/raw` asset directory and a hash of the keystore is computed.
`OneginiConfigModel` (that implements `OneginiClientConfigModel` interface) is being generated with values provided in config.json stored within the ZIP package. To assure tampering detection the tool also computes hash of the generated keystore and puts it into the output config model class. Afterwards the config model is being copied to the project's main sources package, for example `/src/main/java/com/onegini/mobile/demo/`.

### iOS
3. Provide path to `.xcodeproj` file
4. Provide name of the target to which generated files will be added

OneginiConfigModel.h and OneginiConfigModel.m are generated along with certificate files in .pem and .cer formats. Files are added to the project into a group named "Configuration". At runtime SDK will detect those files and use them for configuration. When `Onegini SDK Configurator` tool is used OGOneginiClient should initialized through #initWithDelegate:delegate method.

