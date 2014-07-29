//
//  Alert.h
//  Survana
//
//  Created by Victor Petrov on 7/29/14.
//  Copyright (c) 2014 The Neuroinformatics Research Group at Harvard University. All rights reserved.
//

#import <Foundation/Foundation.h>

@interface Alert : NSObject

+ (BOOL)warning:(NSString*)message andTitle:(NSString*)title;
+ (BOOL)error:(NSString*)message andTitle:(NSString*)title;
+ (BOOL)info:(NSString*)message andTitle:(NSString*)title;
+ (BOOL)alert:(NSString*)message andTitle:(NSString*)title andStyle:(NSAlertStyle)style;

@end