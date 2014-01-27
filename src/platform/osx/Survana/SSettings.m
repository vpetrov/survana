//
//  SSettings.m
//  Survana
//
//  Created by Victor Petrov on 1/25/14.
//  Copyright (c) 2014 The Neuroinformatics Research Group at Harvard University. All rights reserved.
//

#import "SSettings.h"

@interface SSettings ()

@end

@implementation SSettings

//Dashboard Tab
@synthesize cbAuthentication;
@synthesize txtDashUsername;
@synthesize txtDashPassword;

//Web Server Tab
@synthesize txtIP;
@synthesize txtPort;
@synthesize txtSSLCertificate;
@synthesize txtSSLKey;
@synthesize txtWWW;

//Database Tab
@synthesize cbDatabase;
@synthesize txtDBHost;
@synthesize txtDBUsername;
@synthesize txtDBPassword;

- (id)initWithWindow:(NSWindow *)window
{
    self = [super initWithWindow:window];
    if (self) {
        // Initialization code here.
    }
    return self;
}

- (void)windowDidLoad
{
    [super windowDidLoad];
    
    dialog = [NSOpenPanel openPanel];
    [dialog setAllowsMultipleSelection:NO];
    
    //Select MongoDB
    [cbDatabase selectItemAtIndex:0];
    
    //load settings
    [self loadConfiguration:filename];
    
    [[self window] orderFront:nil];
}

- (void) loadConfiguration:(NSString*)file {
    NSLog(@"Loading data from %@", file);
    
    NSData *data = [NSData dataWithContentsOfFile:file];
    
    if (data == nil) {
        NSLog(@"Failed to read file %@", file);
        return;
    }
    
    NSError *error;
    configuration = [NSJSONSerialization JSONObjectWithData:data options:NSJSONReadingMutableLeaves|NSJSONReadingMutableContainers error:&error];
    if (configuration == nil ) {
        NSLog(@"Failed to decode JSON data: %@", error);
        return;
    }
    
    /* Update UI fields */
    
    //Dashboard tab
    [cbAuthentication selectItemAtIndex:0]; //only built-in is supported ATM
    [self updateStringField:txtDashUsername for:@"admin"];
    [self updateStringField:txtDashPassword for:@"password"];
    
    //Web Server tab
    [self updateStringField:txtIP for:@"ip"];
    if (configuration[@"port"] != nil) {
        //[txtPort setStringValue:configuration[@"port"]];
    }

    [self updateStringField:txtWWW for:@"www"];
    [self updateStringField:txtSSLCertificate for:@"sslcert"];
    [self updateStringField:txtSSLKey for:@"sslkey"];
    
    //Database tab
    [cbAuthentication selectItemAtIndex:0]; //only MongoDB is supported ATM
    
    NSURL *url = [NSURL URLWithString:configuration[@"db"]];

    if (url.path != nil) {
        [txtDBHost setStringValue:url.path];
    }

    if (url.user != nil) {
        [txtDBUsername setStringValue:url.user];
    }

    if (url.password != nil) {
        [txtDBPassword setStringValue:url.password];
    }
   
    NSLog(@"Loaded JSON: %@", configuration);
}


- (void) saveConfiguration:(NSString*)file {
    
    if (configuration != nil) {
        NSLog(@"Saving data to %@: %@", file, configuration);
        
        NSError *error;
        NSData *data = [NSJSONSerialization dataWithJSONObject:configuration options:NSJSONWritingPrettyPrinted error:&error];
        if (data == nil) {
            NSLog(@"Failed to serialize to JSON: %@", configuration);
            return;
        }
        
        [data writeToFile:file options:kNilOptions error:&error];
        if (error != nil ) {
            NSLog(@"Failed to write configuration file: %@", error);
            return;
        }
        
        NSLog(@"Wrote JSON to %@: %@", file, data);
    } else {
        NSLog(@"No configuration data to save");
    }
    
    //close Settings window
    [self close];
}

- (void)setFilePath:(NSString *)path {
    filename = path;
}

- (void)updateStringField:(NSTextField*)field for:(NSString*)name {
    NSString *value = configuration[name];
    
    if (value == nil) {
        return;
    }
    
    [field setStringValue:value];
}

-(IBAction)saveSettings:(id)sender {
    NSLog(@"Saving settings");
    //Dashboard tab
    configuration[@"authentication"] = [[cbAuthentication stringValue] lowercaseString];
    configuration[@"admin"] = [txtDashUsername stringValue];
    configuration[@"password"] = [txtDashPassword stringValue];
    
    //Web Server tab
    configuration[@"ip"] = [txtIP stringValue];
    configuration[@"port"]= [NSNumber numberWithInteger:[txtPort integerValue]];
    configuration[@"www"] = [txtWWW stringValue];
    configuration[@"sslcert"] = [txtSSLCertificate stringValue];
    configuration[@"sslkey"] = [txtSSLKey stringValue];
    
    //Database tab
    NSMutableString *dburl = [NSMutableString stringWithString:@"mongodb://"];
    NSString *dbUsername = [txtDBUsername stringValue];
    NSString *dbPassword = [txtDBPassword stringValue];
    
    if ([dbUsername length] > 0) {
        [dburl appendString:dbUsername];
        [dburl appendString:@":"];
        [dburl appendString:dbPassword];
    }

    [dburl appendString:[txtDBHost stringValue]];
    configuration[@"db"] = dburl;
    
    //save configuration
    [self saveConfiguration:filename];
}

-(IBAction)browseForSSLCertificate:(id)sender {
    NSLog(@"Browsing for SSL Certificate");
    [self browseForFile:txtSSLCertificate];
}

-(IBAction)browseForSSLKey:(id)sender {
    NSLog(@"Browsing for SSL Key");
    [self browseForFile:txtSSLKey];
}

-(IBAction)browseForWWW:(id)sender {
    NSLog(@"Browsing for WWW");
    [self browseForFolder:txtWWW];
}

-(void)browseForFile:(NSTextField*)field {
    [dialog setCanChooseDirectories:NO];
    [dialog setCanChooseFiles:YES];

    if ([dialog runModal] == NSOKButton) {
        NSString *path = [[dialog URL] path];
        [field setStringValue:path];
    }
}

-(void)browseForFolder:(NSTextField*)field {
    [dialog setCanChooseDirectories:YES];
    [dialog setCanChooseFiles:NO];

    if ([dialog runModal] == NSOKButton) {
        NSString *path = [[dialog URL] path];
        NSLog(@"Dir: %@", path);
        [field setStringValue:path];
    }
}

@end
