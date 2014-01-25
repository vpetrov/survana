//
//  SAppDelegate.m
//  Survana
//
//  Created by Victor Petrov on 1/16/14.
//  Copyright (c) 2014 The Neuroinformatics Research Group at Harvard University. All rights reserved.
//

#import "SAppDelegate.h"
#include <sys/types.h>
#include <sys/socket.h>
#include <netinet/in.h>


@implementation SAppDelegate

- (void)applicationDidFinishLaunching :(NSNotification *)aNotification
{
    // Insert code here to initialize your application
    statusItem = [[NSStatusBar systemStatusBar] statusItemWithLength:NSSquareStatusItemLength];
    NSBundle *bundle = [NSBundle mainBundle];
    statusImage = [[NSImage alloc] initWithContentsOfFile:[bundle pathForResource:@"survana-logo-bw" ofType:@"png"]];
    highlightedStatusImage = [[NSImage alloc] initWithContentsOfFile:[bundle pathForResource:@"survana-logo-bw-inverted" ofType:@"png"]];
    //Sets the images in our NSStatusItem
    [statusItem setImage:statusImage];
    [statusItem setAlternateImage:highlightedStatusImage];
    [statusItem setMenu:statusMenu];
    [statusItem setToolTip:@"Survana"];
    [statusItem setHighlightMode:YES];
}

- (IBAction)openDashboard :(id)sender {
    NSURL *url = [NSURL URLWithString:@"https://localhost:4443/dashboard/"];
    
    [[NSWorkspace sharedWorkspace] openURL:url];
}

- (IBAction)startServer :(id)sender {
    NSBundle *bundle = [NSBundle mainBundle];
    NSString *serverDir = [bundle pathForResource:@"server" ofType:@""];
    NSString *serverBin = [bundle pathForResource:@"server/server" ofType:@""];
    NSTask *server = [[NSTask alloc] init];
    [server setCurrentDirectoryPath:serverDir];
    [server setLaunchPath:serverBin];
    
    NSLog(@"Launching %@", serverBin);
    [server launch];
}

- (IBAction)stopServer :(id)sender {
    NSLog(@"%@", [[NSWorkspace sharedWorkspace] runningApplications]);
}

//About menu action: displays version information and copyright notice
- (IBAction)about :(id)sender {
    [self info:@"Survana v.1.0\n\n(c) 2014 The Neuroinformatics Research Group at Harvard University" andTitle:@"About Survana"];
}

/* ALERTS */

//displays a warning window
- (BOOL)warning:(NSString*)message andTitle:(NSString*)title
{
    return [self alert:message andTitle:title andStyle:NSCriticalAlertStyle];
}

//displays an error window
- (BOOL)error:(NSString*)message andTitle:(NSString*)title
{
    return [self alert:message andTitle:title andStyle:NSCriticalAlertStyle];
}

//displays an informational window
- (BOOL)info:(NSString*)message andTitle:(NSString*)title
{
    return [self alert:message andTitle:title andStyle:NSInformationalAlertStyle];
}

//generic alert window with customizable style
- (BOOL)alert:(NSString*)message andTitle:(NSString*)title andStyle:(NSAlertStyle)style
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
