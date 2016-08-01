#import "OneginiConfigModel.h"

@implementation OneginiConfigModel

+ (NSArray *)certificates
{
    return @[@""]; //Base64Certificates
}

+ (NSDictionary *)configuration
{
    return @{
             @"kOGAppIdentifier" : @"",
             @"kOGAppScheme" : @"",
             @"kOGAppPlatform" : @"ios",
             @"kOGAppVersion" : @"",
             @"kOGAppBaseURL" : @"",
             @"kOGMaxPinFailures" : @"",
             @"kOGResourceBaseURL" : @"",
             @"kOGRedirectURL" : @"",
             @"kOGStoreCookies" : @(YES),
             @"kOGUseEmbeddedWebview" : @(YES),
             @"kOGDeviceName" : @"",
             };
}

@end