//
//  Alert.m
//  Survana
//
//  Created by Victor Petrov on 7/29/14.
//  Copyright (c) 2014 The Neuroinformatics Research Group at Harvard University. All rights reserved.
//

#import "Alert.h"

@implementation Alert

/* ALERTS */

//displays a warning window
+ (BOOL)warning:(NSString*)message andTitle:(NSString*)title
{
    return [self alert:message andTitle:title andStyle:NSCriticalAlertStyle];
}

//displays an error window
+ (BOOL)error:(NSString*)message andTitle:(NSString*)title
{
    return [self alert:message andTitle:title andStyle:NSCriticalAlertStyle];
}

//displays an informational window
+ (BOOL)info:(NSString*)message andTitle:(NSString*)title
{
    return [self alert:message andTitle:title andStyle:NSInformationalAlertStyle];
}

//generic alert window with customizable style
+ (BOOL)alert:(NSString*)message andTitle:(NSString*)title andStyle:(NSAlertStyle)style
{
    NSAlert *alert = [[NSAlert alloc] init];
    [alert addButtonWithTitle:@"OK"];
    [alert setMessageText:title];
    [alert setInformativeText:message];
    [alert setAlertStyle:style];
    [alert runModal];
    
    return YES;
}

@end
