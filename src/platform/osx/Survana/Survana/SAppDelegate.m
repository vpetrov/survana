//
//  SAppDelegate.m
//  Survana
//
//  Created by Victor Petrov on 1/16/14.
//  Copyright (c) 2014 The Neuroinformatics Research Group at Harvard University. All rights reserved.
//

#import "SAppDelegate.h"
#import "SSettings.h"
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
    
    //initialize the settings window
    settingsWindow = [[SSettings alloc] initWithWindowNibName:@"SSettings"];
    //get the path to the server configuration file
    NSString *serverConfig = [[NSBundle mainBundle] pathForResource:@"server/survana" ofType:@"json"];
    //set the path to the configuration file for the settings window
    [settingsWindow setFilePath:serverConfig];
}

- (IBAction)openDashboard :(id)sender {
    NSURL *url = [NSURL URLWithString:@"https://localhost:4443/dashboard/"];
    
    [[NSWorkspace sharedWorkspace] openURL:url];
}

- (IBAction)start:(id)sender {
    [self startMongoDB];
    [self startServer];
}

- (IBAction)settings:(id)sender {
    [NSApp activateIgnoringOtherApps:YES];
    [settingsWindow showWindow:nil];
}

- (void)startServer {
    NSBundle *bundle = [NSBundle mainBundle];
    NSString *serverDir = [bundle pathForResource:@"server" ofType:@""];
    NSString *serverBin = [bundle pathForResource:@"server/server" ofType:@""];
    NSTask *server = [[NSTask alloc] init];
    NSLog(@"ServerDir %@", serverDir);
    NSLog(@"ServerBin %@", serverBin);
    [server setCurrentDirectoryPath:serverDir];
    [server setLaunchPath:serverBin];
    
    NSLog(@"Launching %@", serverBin);
    [server launch];
}

- (void)startMongoDB {
    NSBundle *bundle = [NSBundle mainBundle];
    NSString *serverDir = [bundle pathForResource:@"mongodb" ofType:@""];
    NSString *serverBin = [bundle pathForResource:@"mongodb/bin/mongod" ofType:@""];
    NSString *serverData = [bundle pathForResource:@"db" ofType:@""];
    NSArray *args = [[NSArray alloc] initWithObjects:@"--dbpath", serverData, nil];
    NSTask *server = [[NSTask alloc] init];
    NSLog(@"ServerDir %@", serverDir);
    NSLog(@"ServerBin %@", serverBin);
    NSLog(@"ServerBin %@", serverData);
    [server setCurrentDirectoryPath:serverDir];
    [server setLaunchPath:serverBin];
    [server setArguments:args];
    
    
    NSLog(@"Launching %@", serverBin);
    [server launch];
}

- (IBAction)stopServer :(id)sender {
    NSLog(@"Not implemented");
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
