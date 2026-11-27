//go:build darwin

// AppKit + WebKit: ventana nativa macOS, sin navegador externo y sin dependencias
// de terceros. Autor: Cristian Barallobre Rúa.
#import <Cocoa/Cocoa.h>
#import <WebKit/WebKit.h>
#include <stdlib.h>
extern void cine_log_ui_event(char *event);
extern char *cine_store_draft(char *message);

@interface CineMacDelegate : NSObject <NSApplicationDelegate, NSWindowDelegate, WKNavigationDelegate, WKUIDelegate, WKDownloadDelegate, WKScriptMessageHandler>
@property(retain) NSURL *appURL;
@property(retain) NSMapTable *downloads;
@property(retain) WKWebView *appView;
@property BOOL preparingClose;
@property BOOL closeApproved;
@end

@implementation CineMacDelegate
- (void)userContentController:(WKUserContentController *)controller didReceiveScriptMessage:(WKScriptMessage *)message {
    WKSecurityOrigin *source = message.frameInfo.securityOrigin;
    if (!message.frameInfo.isMainFrame ||
        ![source.protocol isEqualToString:self.appURL.scheme] || ![source.host isEqualToString:self.appURL.host] ||
        source.port != self.appURL.port.integerValue || ![message.body isKindOfClass:[NSString class]]) return;
    NSString *event = message.body;
    if ([message.name isEqualToString:@"cineDiagnostics"]) {
        if (event.length <= 128) cine_log_ui_event((char *)[event UTF8String]);
    } else if ([message.name isEqualToString:@"cineDraft"] && [event lengthOfBytesUsingEncoding:NSUTF8StringEncoding] <= 64 * 1024 * 1024 + 256) {
        char *ack = cine_store_draft((char *)[event UTF8String]);
        if (ack) {
            NSString *reply = [NSString stringWithUTF8String:ack];
            free(ack);
            // Reply consists solely of validated ASCII protocol fields.
            NSString *script = [NSString stringWithFormat:@"window.cineDraftAck?.('%@')", reply];
            [message.webView evaluateJavaScript:script completionHandler:nil];
        }
    }
}
- (void)webViewWebContentProcessDidTerminate:(WKWebView *)webView {
    cine_log_ui_event("renderer_terminated");
}
- (BOOL)applicationShouldTerminateAfterLastWindowClosed:(NSApplication *)sender {
    return YES;
}
- (void)windowWillClose:(NSNotification *)notification {
    [NSApp stop:nil];
}
- (void)prepareClose:(void (^)(BOOL))reply {
    self.preparingClose = YES;
    __block BOOL finished = NO;
    void (^finish)(BOOL, NSString *) = ^(BOOL saved, NSString *detail) {
        if (finished) return;
        finished = YES;
        if (saved) {
            cine_log_ui_event("close_ready");
            self.preparingClose = NO;
            reply(YES);
            return;
        }
        cine_log_ui_event("close_failed");
        NSAlert *alert = [[[NSAlert alloc] init] autorelease];
        alert.alertStyle = NSAlertStyleWarning;
        alert.messageText = @"No se pudo confirmar el guardado local.";
        NSString *reason = detail.length ? detail : @"La interfaz no respondió; no se puede identificar qué cambios quedaron pendientes.";
        alert.informativeText = [reason stringByAppendingString:@"\n\nSi cierras ahora, puedes perder cambios sin guardar. Puedes permanecer en la aplicación para reintentar o exportar una copia de emergencia."];
        [alert addButtonWithTitle:@"Cancelar"];
        [alert addButtonWithTitle:@"Cerrar de todas formas"];
        NSModalResponse choice = [alert runModal];
        cine_log_ui_event(choice == NSAlertSecondButtonReturn ? "close_forced" : "close_cancelled");
        self.preparingClose = NO;
        reply(choice == NSAlertSecondButtonReturn);
    };
    dispatch_after(dispatch_time(DISPATCH_TIME_NOW, 20 * NSEC_PER_SEC), dispatch_get_main_queue(), ^{
        finish(NO, nil);
    });
    [self.appView callAsyncJavaScript:
        @"if (typeof window.cinePrepareClose !== 'function') return true; return await window.cinePrepareClose() ? true : (window.cineCloseFailureDetail || 'Guardado no confirmado.');"
        arguments:@{} inFrame:nil inContentWorld:[WKContentWorld pageWorld]
        completionHandler:^(id result, NSError *error) {
            BOOL saved = !error && [result isKindOfClass:[NSNumber class]] && [result boolValue];
            NSString *detail = !error && [result isKindOfClass:[NSString class]] ? result : nil;
            finish(saved, detail);
        }];
}
- (BOOL)windowShouldClose:(NSWindow *)window {
    if (self.closeApproved) return YES;
    if (!self.preparingClose) {
        [self prepareClose:^(BOOL saved) {
            self.closeApproved = saved;
            if (saved) [window close];
        }];
    }
    return NO;
}
- (NSApplicationTerminateReply)applicationShouldTerminate:(NSApplication *)sender {
    if (self.closeApproved) return NSTerminateNow;
    if (self.preparingClose) return NSTerminateCancel;
    [self prepareClose:^(BOOL saved) {
        self.closeApproved = saved;
        [sender replyToApplicationShouldTerminate:saved];
    }];
    return NSTerminateLater;
}
- (void)webView:(WKWebView *)webView decidePolicyForNavigationAction:(WKNavigationAction *)action
 decisionHandler:(void (^)(WKNavigationActionPolicy))decisionHandler {
    NSURL *url = action.request.URL;
    WKSecurityOrigin *source = action.sourceFrame.securityOrigin;
    if (action.shouldPerformDownload && [url.scheme isEqualToString:@"blob"] &&
        [source.protocol isEqualToString:self.appURL.scheme] && [source.host isEqualToString:self.appURL.host] &&
        source.port == self.appURL.port.integerValue) {
        decisionHandler(WKNavigationActionPolicyDownload);
        return;
    }
    if (([url.host isEqualToString:self.appURL.host] && [url.scheme isEqualToString:self.appURL.scheme] &&
         [url.port isEqual:self.appURL.port]) || [url.absoluteString isEqualToString:@"about:blank"]) {
        decisionHandler(WKNavigationActionPolicyAllow);
        return;
    }
    // TMDB y las plataformas se abren en el navegador, sin reemplazar la app.
    if ([url.scheme isEqualToString:@"https"] || [url.scheme isEqualToString:@"http"]) {
        [[NSWorkspace sharedWorkspace] openURL:url];
    }
    decisionHandler(WKNavigationActionPolicyCancel);
}
- (void)webView:(WKWebView *)view runOpenPanelWithParameters:(WKOpenPanelParameters *)options
 initiatedByFrame:(WKFrameInfo *)frame completionHandler:(void (^)(NSArray<NSURL *> *))reply {
    NSOpenPanel *panel = [NSOpenPanel openPanel];
    panel.allowsMultipleSelection = options.allowsMultipleSelection;
    panel.canChooseDirectories = NO;
    [panel beginSheetModalForWindow:view.window completionHandler:^(NSModalResponse result){
        reply(result == NSModalResponseOK ? panel.URLs : nil);
    }];
}
- (void)webView:(WKWebView *)view runJavaScriptConfirmPanelWithMessage:(NSString *)message
 initiatedByFrame:(WKFrameInfo *)frame completionHandler:(void (^)(BOOL))reply {
    NSAlert *alert = [[[NSAlert alloc] init] autorelease];
    alert.messageText = message;
    [alert addButtonWithTitle:@"Aceptar"];
    [alert addButtonWithTitle:@"Cancelar"];
    [alert beginSheetModalForWindow:view.window completionHandler:^(NSModalResponse result){reply(result == NSAlertFirstButtonReturn);}];
}
- (void)webView:(WKWebView *)view runJavaScriptAlertPanelWithMessage:(NSString *)message
 initiatedByFrame:(WKFrameInfo *)frame completionHandler:(void (^)(void))reply {
    NSAlert *alert = [[[NSAlert alloc] init] autorelease];
    alert.messageText = message;
    [alert addButtonWithTitle:@"Aceptar"];
    [alert beginSheetModalForWindow:view.window completionHandler:^(NSModalResponse result){reply();}];
}
- (void)webView:(WKWebView *)view navigationAction:(WKNavigationAction *)action didBecomeDownload:(WKDownload *)download {
    download.delegate = self;
}
- (void)download:(WKDownload *)download decideDestinationUsingResponse:(NSURLResponse *)response
 suggestedFilename:(NSString *)filename completionHandler:(void (^)(NSURL *))reply {
    NSSavePanel *panel = [NSSavePanel savePanel];
    panel.nameFieldStringValue = filename.lastPathComponent;
    [panel beginSheetModalForWindow:download.webView.window completionHandler:^(NSModalResponse result){
        if (result != NSModalResponseOK) { reply(nil); return; }
        NSError *error = nil;
        NSURL *dir = [[NSFileManager defaultManager] URLForDirectory:NSItemReplacementDirectory
            inDomain:NSUserDomainMask appropriateForURL:panel.URL create:YES error:&error];
        if (!dir) { reply(nil); [NSApp presentError:error]; return; }
        NSURL *temp = [dir URLByAppendingPathComponent:@"descarga"];
        [self.downloads setObject:@[temp, panel.URL] forKey:download];
        reply(temp);
    }];
}
- (void)downloadDidFinish:(WKDownload *)download {
    NSArray *urls = [self.downloads objectForKey:download];
    if (!urls) return;
    NSURL *temp = urls[0], *target = urls[1];
    if (rename(temp.fileSystemRepresentation, target.fileSystemRepresentation) != 0)
        [NSApp presentError:[NSError errorWithDomain:NSPOSIXErrorDomain code:errno userInfo:nil]];
    [[NSFileManager defaultManager] removeItemAtURL:temp.URLByDeletingLastPathComponent error:nil];
    [self.downloads removeObjectForKey:download];
}
- (void)download:(WKDownload *)download didFailWithError:(NSError *)error resumeData:(NSData *)data {
    NSArray *urls = [self.downloads objectForKey:download];
    if (urls) [[NSFileManager defaultManager] removeItemAtURL:[urls[0] URLByDeletingLastPathComponent] error:nil];
    [self.downloads removeObjectForKey:download];
    if (error.code != NSURLErrorCancelled) [NSApp presentError:error];
}
- (WKWebView *)webView:(WKWebView *)webView createWebViewWithConfiguration:(WKWebViewConfiguration *)configuration
 forNavigationAction:(WKNavigationAction *)navigationAction windowFeatures:(WKWindowFeatures *)windowFeatures {
    NSURL *url = navigationAction.request.URL;
    if ([url.scheme isEqualToString:@"https"] || [url.scheme isEqualToString:@"http"]) {
        [[NSWorkspace sharedWorkspace] openURL:url];
    }
    return nil;
}
@end

int cine_macos_run(const char *targetURL) {
    if (!targetURL) return 1;
    @autoreleasepool {
        NSString *address = [NSString stringWithUTF8String:targetURL];
        NSURL *url = [NSURL URLWithString:address];
        if (!url) return 2;

        NSApplication *app = [NSApplication sharedApplication];
        [app setActivationPolicy:NSApplicationActivationPolicyRegular];

        NSMenu *menuBar = [[NSMenu alloc] initWithTitle:@"MenuBar"];
        NSMenuItem *appMenuItem = [[NSMenuItem alloc] initWithTitle:@"Selector de Cine 365" action:NULL keyEquivalent:@""];
        NSMenu *appMenu = [[NSMenu alloc] initWithTitle:@"Selector de Cine 365"];
        NSMenuItem *quitItem = [[NSMenuItem alloc] initWithTitle:@"Salir de Selector de Cine 365" action:@selector(terminate:) keyEquivalent:@"q"];
        [appMenu addItem:quitItem];
        [appMenuItem setSubmenu:appMenu];
        [menuBar addItem:appMenuItem];
        [app setMainMenu:menuBar];

        NSUInteger style = NSWindowStyleMaskTitled | NSWindowStyleMaskClosable |
                           NSWindowStyleMaskMiniaturizable | NSWindowStyleMaskResizable;
        NSWindow *window = [[NSWindow alloc] initWithContentRect:NSMakeRect(0, 0, 1210, 830)
                  styleMask:style backing:NSBackingStoreBuffered defer:NO];
        if (!window) return 3;
        [window setTitle:@"Selector de Cine 365 · Cristian Barallobre Rúa"];
        [window setMinSize:NSMakeSize(950, 670)];
        [window center];

        CineMacDelegate *delegate = [[CineMacDelegate alloc] init];
        delegate.appURL = url;
        delegate.downloads = [NSMapTable strongToStrongObjectsMapTable];
        [app setDelegate:delegate];
        [window setDelegate:delegate];

        WKWebViewConfiguration *configuration = [[WKWebViewConfiguration alloc] init];
        configuration.websiteDataStore = [WKWebsiteDataStore nonPersistentDataStore];
        [configuration.userContentController addScriptMessageHandler:delegate name:@"cineDiagnostics"];
        [configuration.userContentController addScriptMessageHandler:delegate name:@"cineDraft"];
        WKWebView *webView = [[WKWebView alloc] initWithFrame:[[window contentView] bounds]
                                                 configuration:configuration];
        delegate.appView = webView;
        webView.autoresizingMask = NSViewWidthSizable | NSViewHeightSizable;
        [webView setNavigationDelegate:delegate];
        [webView setUIDelegate:delegate];
        [[window contentView] addSubview:webView];
        [webView loadRequest:[NSURLRequest requestWithURL:url]];
        [window makeKeyAndOrderFront:nil];
        [app activateIgnoringOtherApps:YES];
        [app run];
        [configuration.userContentController removeScriptMessageHandlerForName:@"cineDiagnostics"];
        [configuration.userContentController removeScriptMessageHandlerForName:@"cineDraft"];
        return 0;
    }
}
