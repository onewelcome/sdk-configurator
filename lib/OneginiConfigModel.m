#import "OneginiConfigModel.h"

@implementation OneginiConfigModel

+ (NSArray *)certificates
{
    return @[@""]; //Base64Certificates
}

+ (NSDictionary *)configuration
{
    return @{
             @"ONGAppIdentifier" : @"",
             @"ONGAppPlatform" : @"ios",
             @"ONGAppVersion" : @"",
             @"ONGAppBaseURL" : @"",
             @"ONGResourceBaseURL" : @"",
             @"ONGRedirectURL" : @"",
             };
}

@end