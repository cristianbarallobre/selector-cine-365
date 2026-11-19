//go:build darwin

// AppKit + WebKit: ventana nativa macOS, sin navegador externo y sin dependencias
// de terceros. Autor: Cristian Barallobre Rúa.
#import <Cocoa/Cocoa.h>
#import <WebKit/WebKit.h>

@interface CineMacDelegate : NSObject <NSApplicationDelegate, NSWindowDelegate, WKNavigationDelegate, WKUIDelegate>
@end

@implementation CineMacDelegate
- (BOOL)applicationShouldTerminateAfterLastWindowClosed:(NSApplication *)sender {
    return YES;
}
- (void)windowWillClose:(NSNotification *)notification {
    [NSApp stop:nil];
}
- (void)webView:(WKWebView *)webView decidePolicyForNavigationAction:(WKNavigationAction *)action
 decisionHandler:(void (^)(WKNavigationActionPolicy))decisionHandler {
    NSURL *url = action.request.URL;
    NSString *host = url.host ?: @"";
    if ([host isEqualToString:@"127.0.0.1"] || [host isEqualToString:@"localhost"] ||
        [url.scheme isEqualToString:@"about"]) {
        decisionHandler(WKNavigationActionPolicyAllow);
        return;
    }
    // TMDB y las plataformas se abren en el navegador, sin reemplazar la app.
    if ([url.scheme isEqualToString:@"https"] || [url.scheme isEqualToString:@"http"]) {
        [[NSWorkspace sharedWorkspace] openURL:url];
    }
    decisionHandler(WKNavigationActionPolicyCancel);
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
        [app setDelegate:delegate];
        [window setDelegate:delegate];

        WKWebViewConfiguration *configuration = [[WKWebViewConfiguration alloc] init];
        configuration.websiteDataStore = [WKWebsiteDataStore nonPersistentDataStore];
        WKWebView *webView = [[WKWebView alloc] initWithFrame:[[window contentView] bounds]
                                                 configuration:configuration];
        webView.autoresizingMask = NSViewWidthSizable | NSViewHeightSizable;
        [webView setNavigationDelegate:delegate];
        [webView setUIDelegate:delegate];
        [[window contentView] addSubview:webView];
        [webView loadRequest:[NSURLRequest requestWithURL:url]];
        [window makeKeyAndOrderFront:nil];
        [app activateIgnoringOtherApps:YES];
        [app run];
        return 0;
    }
}
