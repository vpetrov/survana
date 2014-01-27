//
//  SSettings.h
//  Survana
//
//  Created by Victor Petrov on 1/25/14.
//  Copyright (c) 2014 The Neuroinformatics Research Group at Harvard University. All rights reserved.
//

#import <Cocoa/Cocoa.h>

@interface SSettings : NSWindowController {
    NSOpenPanel* dialog;
    NSString* filename;
    NSMutableDictionary* configuration;
}

//Dashboard Tab
@property (strong, nonatomic) IBOutlet NSComboBox *cbAuthentication;
@property (strong, nonatomic) IBOutlet NSTextField *txtDashUsername;
@property (strong, nonatomic) IBOutlet NSTextField *txtDashPassword;


//Web Server Tab
@property (strong, nonatomic) IBOutlet NSTextField *txtIP;
@property (strong, nonatomic) IBOutlet NSTextField *txtPort;
@property (strong, nonatomic) IBOutlet NSTextField *txtSSLCertificate;
@property (strong, nonatomic) IBOutlet NSTextField *txtSSLKey;
@property (strong, nonatomic) IBOutlet NSTextField *txtWWW;

//Database Tab
@property (strong, nonatomic) IBOutlet NSComboBox *cbDatabase;
@property (strong, nonatomic) IBOutlet NSTextField *txtDBHost;
@property (strong, nonatomic) IBOutlet NSTextField *txtDBUsername;
@property (strong, nonatomic) IBOutlet NSTextField *txtDBPassword;

-(void)browseForFile:(NSTextField*)field;
-(void)browseForFolder:(NSTextField*)field;

-(void)setFilePath:(NSString*)path;
-(void)loadConfiguration:(NSString*)file;
-(void)updateStringField:(NSTextField*)field for:(NSString*)name;

-(IBAction)saveSettings:(id)sender;
-(IBAction)browseForSSLCertificate:(id)sender;
-(IBAction)browseForSSLKey:(id)sender;

@end
