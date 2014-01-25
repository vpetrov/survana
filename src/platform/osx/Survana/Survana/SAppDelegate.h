//
//  SAppDelegate.h
//  Survana
//
//  Created by Victor Petrov on 1/16/14.
//  Copyright (c) 2014 The Neuroinformatics Research Group at Harvard University. All rights reserved.
//

#import <Cocoa/Cocoa.h>

@interface SAppDelegate : NSObject <NSApplicationDelegate> {
    IBOutlet NSMenu *statusMenu;
    
    NSStatusItem    *statusItem;
    NSImage         *statusImage;
    NSImage         *highlightedStatusImage;
}

- (IBAction)openDashboard:(id)sender;
- (IBAction)startServer:(id)sender;
- (IBAction)about:(id)sender;

//displays a warning window
- (BOOL)warning:(NSString*)message andTitle:(NSString*)title;
//displays an error window
- (BOOL)error:(NSString*)message andTitle:(NSString*)title;
//displays an informational window
- (BOOL)info:(NSString*)message andTitle:(NSString*)title;
//generic alert window with customizable style
- (BOOL)alert:(NSString*)message andTitle:(NSString*)title andStyle:(NSAlertStyle)style;
@end
