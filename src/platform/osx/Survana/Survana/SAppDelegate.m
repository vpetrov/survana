//
//  SAppDelegate.m
//  Survana
//
//  Created by Victor Petrov on 1/16/14.
//  Copyright (c) 2014 The Neuroinformatics Research Group at Harvard University. All rights reserved.
//

#import "SAppDelegate.h"
#import "SSettings.h"
#import <sys/sysctl.h>
#import <signal.h>
#import <unistd.h>

@implementation SAppDelegate

- (void)applicationDidFinishLaunching :(NSNotification *)aNotification
{
    // Insert code here to initialize your application
    statusItem = [[NSStatusBar systemStatusBar] statusItemWithLength:NSSquareStatusItemLength];
    NSBundle *bundle = [NSBundle mainBundle];
    servicesPath = [bundle pathForResource:@"services" ofType:@""];

    statusImage = [[NSImage alloc] initWithContentsOfFile:[bundle pathForResource:@"survana-logo-bw" ofType:@"png"]];
    highlightedStatusImage = [[NSImage alloc] initWithContentsOfFile:[bundle pathForResource:@"survana-logo-bw-inverted" ofType:@"png"]];
    //Sets the images in our NSStatusItem
    [statusItem setImage:statusImage];
    [statusItem setAlternateImage:highlightedStatusImage];
    [statusItem setMenu:statusMenu];
    [statusItem setToolTip:@"Survana"];
    [statusItem setHighlightMode:YES];
    
    //create all files and directories
    if (![self initEnvironment]) {
        [self quit:nil];
    }

    //initialize the settings window
    settingsWindow = [[SSettings alloc] initWithWindowNibName:@"SSettings"];
    
    //set the path to the configuration file for the settings window
    [settingsWindow setFilePath:serverConfig];
    
    [self updateServiceStatus];
}

- (IBAction)openDashboard :(id)sender {
    NSURL *url = [NSURL URLWithString:@"https://localhost:4443/dashboard/"];
    
    [[NSWorkspace sharedWorkspace] openURL:url];
}

- (IBAction)start:(id)sender {
    [self startMongoDB];
    [self startServer];

    //wait for 0.5 seconds, then update the menu
    [self performSelector:@selector(updateServiceStatus) withObject:nil afterDelay:0.5];
} 

- (IBAction)settings:(id)sender {
    [NSApp activateIgnoringOtherApps:YES];
    [settingsWindow showWindow:nil];
}

//creates all required folders and files
- (BOOL)initEnvironment {
    NSBundle *bundle = [NSBundle mainBundle];
    
    //create the environment path
    environmentPath = [NSString stringWithFormat:@"%@/%@", NSHomeDirectory(), @".survana"];
    if (![self createFolder:environmentPath withPermissions:0711]) {
        return NO;
    }
    
    //get the path to the bundled configuration file
    NSString *bundledServerConfig = [bundle pathForResource:@"survana" ofType:@"json"];

    //the actual configuration should be in the env folder
    serverConfig = [NSString stringWithFormat:@"%@/%@", environmentPath, @"config.json"];
    
    if (![self createFile:serverConfig from:bundledServerConfig withPermissions:0600]) {
        return NO;
    }
    
    NSLog(@"server config path: %@", serverConfig);
    
    //create the database folder
    dbDataPath = [NSString stringWithFormat:@"%@/%@", environmentPath, @"mongodb"];
    
    if (![self createFolder:dbDataPath withPermissions:0711]) {
        return NO;
    }

    return YES;
}

- (BOOL)createFolder:(NSString*)folderPath {
    NSFileManager *fs = [NSFileManager defaultManager];
    BOOL isDir;
    if (! [fs fileExistsAtPath:folderPath isDirectory:&isDir]) {
        NSError *error;
        if (! [fs createDirectoryAtPath:folderPath withIntermediateDirectories:YES attributes:nil error:&error]) {
            NSString *message = [NSString stringWithFormat:@"Failed to create directory: %@", error];
            NSLog(@"%@", message);
            [self error:message andTitle:@"Error"];
            return NO;
        }
    }
    
    return YES;
}

- (BOOL)createFolder:(NSString*)folderPath withPermissions:(int)octal {
    if (![self createFolder:folderPath]) {
        return NO;
    }
    
    if (![self setPermissions:folderPath to:octal]) {
        return NO;
    }
    
    return YES;
}

- (BOOL) createFile:(NSString*)destPath from:(NSString*)srcPath {
    NSError *error;
    NSFileManager *fs = [NSFileManager defaultManager];
    //attempt to copy the file if it doesn't exist
    if (! [fs fileExistsAtPath:destPath]) {
        NSLog(@"Copying %@ to %@", srcPath, destPath);
        if (![[NSFileManager defaultManager] copyItemAtPath:srcPath toPath:destPath error:&error]) {
            NSString *message = [NSString stringWithFormat:@"Failed to copy file %@ from %@: %@", destPath, srcPath, error];
            NSLog(@"%@", message);
            [self error:message andTitle:@"Error"];
            return NO;
        }
    }

    return YES;
}

- (BOOL) createFile:(NSString*)destPath from:(NSString*)srcPath withPermissions:(int)octal {
    
    if (![self createFile:destPath from:srcPath]) {
        return NO;
    }
    
    if (![self setPermissions:destPath to:octal]) {
        return NO;
    }
    
    return YES;
}

- (BOOL) setPermissions:(NSString*)destPath to:(int)octal {
    NSFileManager *fs = [NSFileManager defaultManager];
    NSMutableDictionary *perms = [[NSMutableDictionary alloc] init];
    [perms setObject:[NSNumber numberWithInt:octal] forKey:NSFilePosixPermissions];
    
    NSError *error;
    
    if (![fs setAttributes:perms ofItemAtPath:destPath error:&error]) {
        NSString *message = [NSString stringWithFormat:@"Failed to set permissions of %@ to %d: %@", destPath, octal, error];
        NSLog(@"%@", message);
        [self error:message andTitle:@"Error"];
        return NO;
    };

    return YES;
}

- (void)startServer {
    NSString *serverDir = [NSString stringWithFormat:@"%@/%@", servicesPath, @"server"];
    NSString *serverBin = [NSString stringWithFormat:@"%@/%@", serverDir, @"server"];
    
    NSTask *server = [[NSTask alloc] init];
    NSLog(@"ServerDir %@", serverDir);
    NSLog(@"ServerBin %@", serverBin);
    [server setCurrentDirectoryPath:serverDir];
    [server setLaunchPath:serverBin];
    [server setArguments:[NSArray arrayWithObjects:@"--config", serverConfig, nil]];
    
    NSLog(@"Launching %@ %@", serverBin, server.arguments);
    [server launch];
}

- (void)startMongoDB {
    
    NSString *serverDir = [NSString stringWithFormat:@"%@/%@", servicesPath, @"mongodb"];
    NSString *serverBin = [NSString stringWithFormat:@"%@/%@", serverDir, @"bin/mongod"];

    NSArray *args = [[NSArray alloc] initWithObjects:@"--dbpath", dbDataPath, nil];
    NSTask *server = [[NSTask alloc] init];
    NSLog(@"ServerDir %@", serverDir);
    NSLog(@"ServerBin %@", serverBin);
    NSLog(@"ServerData %@", dbDataPath);
    [server setCurrentDirectoryPath:serverDir];
    [server setLaunchPath:serverBin];
    [server setArguments:args];
    
    NSLog(@"Launching %@", serverBin);
    [server launch];
    
    //wait for 5 seconds, then insert forms.json into the database
    [self performSelector:@selector(importForms) withObject:nil afterDelay:5];
}

- (void) importForms {
    NSString *serverDir = [NSString stringWithFormat:@"%@/%@", servicesPath, @"mongodb"];
    NSString *mongoImportBin = [NSString stringWithFormat:@"%@/%@", serverDir, @"bin/mongoimport"];
    NSString *formsJSON = [[NSBundle mainBundle] pathForResource:@"forms" ofType:@"json"];
    
    NSTask *mongoImport = [[NSTask alloc] init];
    NSLog(@"ImportDir %@", serverDir);
    NSLog(@"mongoIMportBin %@", mongoImportBin);
    
    [mongoImport setCurrentDirectoryPath:serverDir];
    [mongoImport setLaunchPath:mongoImportBin];
    [mongoImport setArguments:[NSArray arrayWithObjects:@"-d", @"dashboard_test", @"-c", @"forms", formsJSON, nil]];
    
    NSLog(@"Launching %@ %@", mongoImportBin, mongoImport.arguments);
    [mongoImport launch];
}

- (IBAction)stop:(id)sender {
    if ([pidServer intValue] > 0) {
        NSLog(@"Shutting down Survana server");
        [self killProcess:pidServer];
    }
    
    if ([pidMongoDB intValue] > 0) {
        NSLog(@"Shutting down MongoDB");
        [self mongoShutdown];
    }
    
    //wait for 0.5 seconds, then update the menu
    [self performSelector:@selector(updateServiceStatus) withObject:nil afterDelay:0.5];
}

//launches 'mongo shutdown.js', because mongod doesn't support --shutdown on OSX, and sending the kill signal
//doesn't make the process go away
- (void)mongoShutdown {
    NSString *mongoDir = [NSString stringWithFormat:@"%@/%@", servicesPath, @"mongodb"];
    NSString *mongoBin = [NSString stringWithFormat:@"%@/%@", mongoDir, @"bin/mongo"];
    NSString *shutdownJS = [[NSBundle mainBundle] pathForResource:@"shutdown" ofType:@"js"];
    NSArray *args = [NSArray arrayWithObjects:shutdownJS, nil];
    NSTask *mongo = [[NSTask alloc] init];

    [mongo setCurrentDirectoryPath:mongoDir];
    [mongo setLaunchPath:mongoBin];
    [mongo setArguments:args];
    
    NSLog(@"Launching %@ %@", mongoBin, args);
    [mongo launch];
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

- (BOOL)killProcess:(NSNumber*)pid_number {
    pid_t pid = [pid_number intValue];
    pid_t gpid = getpgid(pid);
    
    if (gpid < 1) {
        char *err = strerror(errno);
        NSLog(@"Failed to get process group id for pid %@: %s\n", pid_number, err);
        return NO;
    }
    
    //kill process group
    if (killpg(gpid, SIGTERM) == -1) {
        char *err = strerror(errno);
        NSLog(@"Failed to send kill signal to process group %d for pid %@: %s\n", gpid, pid_number, err);
        return NO;
    }

    return YES;
}

- (void)updateServiceStatus {
    [self getPIDs];
    [self updateServiceMenu];
}

- (void)getPIDs {
    NSArray *procs = [self runningProcesses];
    NSUInteger nprocs = [procs count];
    NSDictionary *proc;
    NSString *proc_name;
    
    //set both pids to -1
    pidServer = [NSNumber numberWithInt:-1];
    pidMongoDB = [NSNumber numberWithInt:-1];
    
    for (int i = 0; i < nprocs; i++) {
        proc = procs[i];
        proc_name = proc[@"name"];
        
        if ([proc_name isEqualToString:@"server"]) {
            pidServer = proc[@"pid"];
        } else if ([proc_name isEqualToString:@"mongod"]) {
            pidMongoDB = proc[@"pid"];
        }
    }
}

- (void)updateServiceMenu {
    //if any services are enabled, allow the user to stop them
    if (([pidServer intValue] >= 0) || ([pidMongoDB intValue] > 0)) {
        [startMenu setEnabled:NO];
        [stopMenu setEnabled:YES];
    } else {
        [startMenu setEnabled:YES];
        [stopMenu setEnabled:NO];
    }
    
    if (([pidServer intValue] >= 0) && ([pidMongoDB intValue] >= 0)) {
        [dashboardMenu setEnabled:YES];
    } else {
        [dashboardMenu setEnabled:NO];
    }
}

- (NSArray *)runningProcesses {
    
    int mib[4] = {CTL_KERN, KERN_PROC, KERN_PROC_ALL, 0};
    size_t size;
    int st = sysctl(mib, sizeof(mib) / sizeof(CTL_KERN), NULL, &size, NULL, 0);
    
    struct kinfo_proc * process = NULL;
    struct kinfo_proc * newprocess = NULL;
    
    do {
        size += size / 10;
        newprocess = realloc(process, size);
        
        if (!newprocess){
            
            if (process){
                free(process);
            }
            
            return nil;
        }
        
        process = newprocess;
        st = sysctl(mib, sizeof(mib) / sizeof(CTL_KERN), process, &size, NULL, 0);
        
    } while (st == -1 && errno == ENOMEM);
    
    if (st == 0) {
        
        if (size % sizeof(struct kinfo_proc) == 0){
            unsigned long nprocess = size / sizeof(struct kinfo_proc);
            
            if (nprocess){
                
                NSMutableArray * array = [[NSMutableArray alloc] init];
                
                for (unsigned long i = 0; i < nprocess; i++){
                    
                    NSNumber *processID = [NSNumber numberWithInt:process[i].kp_proc.p_pid];
                    NSString *processName = [[NSString alloc] initWithFormat:@"%s", process[i].kp_proc.p_comm];
                    
                    NSDictionary * dict = [[NSDictionary alloc] initWithObjects:[NSArray arrayWithObjects:processID, processName, nil]
                                                                        forKeys:[NSArray arrayWithObjects:@"pid", @"name", nil]];
                    [array addObject:dict];
                }
                
                free(process);
                return array;
            }
        }
    }
    
    return nil;
}

- (IBAction)quit:(id)sender {
    //stop all services
    [self stop:nil];
    //exit after 0.5 seconds
    [NSApp performSelector:@selector(terminate:) withObject:nil afterDelay:0.5];
}

@end
