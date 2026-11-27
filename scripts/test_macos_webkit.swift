import Cocoa
import WebKit

final class Check: NSObject, WKNavigationDelegate {
    let script: String
    init(_ script: String) { self.script = script }
    func webView(_ view: WKWebView, didFinish navigation: WKNavigation!) {
        waitForOptions(view, remaining: 100)
    }
    func waitForOptions(_ view: WKWebView, remaining: Int) {
        view.evaluateJavaScript("!!document.querySelector('#platformSelect option[value=\"Filmin\"]')") { ready, _ in
            if ready as? Bool == true {
                view.evaluateJavaScript(self.script) { result, error in
                    if let error = error { print(error); exit(1) }
                    exit(result as? Bool == true ? 0 : 1)
                }
            } else if remaining > 0 {
                DispatchQueue.main.asyncAfter(deadline: .now() + 0.2) { self.waitForOptions(view, remaining: remaining-1) }
            } else { print("WKWebView: no aparecieron las plataformas"); exit(1) }
        }
    }
}
let app = NSApplication.shared
let width = CommandLine.arguments.count > 3 ? (Double(CommandLine.arguments[3]) ?? 1200) : 1200
let height = CommandLine.arguments.count > 4 ? (Double(CommandLine.arguments[4]) ?? 820) : 820
let window = NSWindow(contentRect:NSRect(x:0,y:0,width:width,height:height), styleMask:[.titled], backing:.buffered, defer:false)
let browser = WKWebView(frame:window.contentView!.bounds)
let check = Check(try String(contentsOfFile:CommandLine.arguments[2],encoding:.utf8))
browser.navigationDelegate = check
window.contentView = browser
window.orderFront(nil)
browser.load(URLRequest(url:URL(string:CommandLine.arguments[1])!))
DispatchQueue.main.asyncAfter(deadline:.now()+45){print("WKWebView: timeout");exit(1)}
app.run()
