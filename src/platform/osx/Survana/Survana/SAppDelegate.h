//
//  SAppDelegate.h
//  Survana
//
//  Created by Victor Petrov on 1/16/14.
//  Copyright (c) 2014 The Neuroinformatics Research Group at Harvard University. All rights reserved.
//

#import <Cocoa/Cocoa.h>
#import "SSettings.h"
#import "Alert.h"

@interface SAppDelegate : NSObject <NSApplicationDelegate> {
    IBOutlet NSMenu *statusMenu;
    IBOutlet NSMenuItem *dashboardMenu;
    IBOutlet NSMenuItem *startMenu;
    IBOutlet NSMenuItem *stopMenu;
    
    NSString        *environmentPath;
    NSString        *serverConfig;
    NSString        *dbDataPath;
    
    NSString        *servicesPath;
    NSStatusItem    *statusItem;
    NSImage         *statusImage;
    NSImage         *highlightedStatusImage;
    SSettings       *settingsWindow;
    
    NSNumber        *pidServer;
    NSNumber        *pidMongoDB;
}

- (IBAction)openDashboard:(id)sender;
- (IBAction)start:(id)sender;
- (IBAction)stop:(id)sender;
- (IBAction)about:(id)sender;
- (IBAction)settings:(id)sender;
- (IBAction)quit:(id)sender;

@end
