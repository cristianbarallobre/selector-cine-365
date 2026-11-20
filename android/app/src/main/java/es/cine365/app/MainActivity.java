package es.cine365.app;

import android.app.Activity;
import android.content.ClipData;
import android.content.ClipboardManager;
import android.content.Context;
import android.content.Intent;
import android.content.res.Configuration;
import android.graphics.Color;
import android.net.Uri;
import android.os.Bundle;
import android.provider.DocumentsContract;
import android.webkit.JavascriptInterface;
import android.webkit.ValueCallback;
import android.webkit.WebChromeClient;
import android.webkit.WebResourceRequest;
import android.webkit.WebResourceResponse;
import android.webkit.WebView;
import android.webkit.WebViewClient;
import android.webkit.WebSettings;
import android.widget.Toast;
import org.json.JSONObject;
import java.io.ByteArrayInputStream;
import java.io.ByteArrayOutputStream;
import java.io.InputStream;
import java.nio.charset.StandardCharsets;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.concurrent.CountDownLatch;
import java.util.concurrent.TimeUnit;

public class MainActivity extends Activity {
    private static final int OPEN_FILE = 401;
    private static final int EXPORT_FILE = 402;
    private static final int PICK_FOLDER = 403;
    private static final String APP_ORIGIN = "https://appassets.androidplatform.net";
    private static final String MICROSOFT_LOGIN_HOST = "login.microsoftonline.com";
    private static final String MICROSOFT_CALLBACK_SCHEME = "msal696458fe-0672-47e5-9fdd-aa0332e9884b";
    private WebView webView;
    private Backend backend;
    private final ExecutorService workers = Executors.newFixedThreadPool(4);
    private ValueCallback<Uri[]> fileCallback;
    private String exportText, exportName, exportMime;
    private CountDownLatch folderLatch;
    private String pickedFolder = "";
    private String folderError = "";

    @Override public void onCreate(Bundle state) {
        super.onCreate(state);
        getWindow().setStatusBarColor(Color.rgb(13,19,30));
        getWindow().setNavigationBarColor(Color.rgb(13,19,30));
        backend = new Backend(this);
        webView = new WebView(this);
        webView.setBackgroundColor(Color.rgb(13,19,30));
        setContentView(webView);
        WebSettings settings = webView.getSettings();
        settings.setJavaScriptEnabled(true);
        settings.setDomStorageEnabled(true);
        settings.setAllowFileAccess(false);
        settings.setAllowContentAccess(true);
        settings.setMixedContentMode(WebSettings.MIXED_CONTENT_NEVER_ALLOW);
        settings.setBuiltInZoomControls(false);
        settings.setDisplayZoomControls(false);
        settings.setSupportZoom(false);
        webView.addJavascriptInterface(new NativeBridge(), "CineAndroid");
        webView.setWebViewClient(new WebViewClient() {
            @Override public WebResourceResponse shouldInterceptRequest(WebView view, WebResourceRequest request) {
                Uri uri = request.getUrl();
                if (!APP_ORIGIN.equals(uri.getScheme()+"://"+uri.getHost())) return null;
                // La WebView no dispone de un servidor HTTP: resolvemos aquí los carteles de TMDB.
                if ("/api/poster".equals(uri.getPath())) {
                    try {
                        String posterPath=uri.getQueryParameter("path");
                        byte[] image=backend.tmdb.poster(posterPath);
                        String mime=posterPath.endsWith(".png")?"image/png":posterPath.endsWith(".webp")?"image/webp":"image/jpeg";
                        return new WebResourceResponse(mime,null,new ByteArrayInputStream(image));
                    } catch(Exception e) { return errorResponse(); }
                }
                if (!uri.getPath().startsWith("/assets/")) return null;
                String file = uri.getPath().substring(8);
                if (!file.equals("index.html") && !file.equals("icon.png")) return errorResponse();
                try {
                    InputStream input = getAssets().open(file);
                    String mime = file.endsWith(".png") ? "image/png" : "text/html";
                    return new WebResourceResponse(mime, file.endsWith(".png") ? null : "utf-8", input);
                } catch(Exception e) { return errorResponse(); }
            }
            @Override public boolean shouldOverrideUrlLoading(WebView view, WebResourceRequest request) {
                Uri uri=request.getUrl();
                if (APP_ORIGIN.equals(uri.getScheme()+"://"+uri.getHost())) return false;
                if (uri.getScheme()!=null && (uri.getScheme().equals("https") || uri.getScheme().equals("http"))) {
                    try { startActivity(new Intent(Intent.ACTION_VIEW,uri)); } catch(Exception ignored){}
                    return true;
                }
                return true;
            }
        });
        webView.setWebChromeClient(new WebChromeClient() {
            @Override public boolean onShowFileChooser(WebView webView, ValueCallback<Uri[]> callback, FileChooserParams params) {
                if(fileCallback!=null)fileCallback.onReceiveValue(null);
                fileCallback=callback;
                Intent i = new Intent(Intent.ACTION_GET_CONTENT);
                i.setType("*/*"); i.addCategory(Intent.CATEGORY_OPENABLE);
                try { startActivityForResult(Intent.createChooser(i,"Seleccionar lista o copia JSON"), OPEN_FILE); }
                catch(Exception e){fileCallback.onReceiveValue(null);fileCallback=null;return false;}
                return true;
            }
        });
        webView.loadUrl(APP_ORIGIN + "/assets/index.html");
        // Android may create a new activity after the browser returns.
        handleMicrosoftRedirect(getIntent());
    }
    @Override protected void onNewIntent(Intent intent){
        super.onNewIntent(intent);
        setIntent(intent);
        handleMicrosoftRedirect(intent);
    }
    private void handleMicrosoftRedirect(Intent intent){
        if(intent==null||!Intent.ACTION_VIEW.equals(intent.getAction()))return;
        Uri callback=intent.getData();
        if(callback==null||!MICROSOFT_CALLBACK_SCHEME.equalsIgnoreCase(callback.getScheme())
            ||!"auth".equalsIgnoreCase(callback.getHost()))return;
        workers.submit(()->{
            boolean ok=false;
            String message;
            try{
                backend.onedrive.finish(callback);
                ok=true;
                message="Cuenta Microsoft conectada. Ya puedes elegir tu carpeta de OneDrive.";
            }catch(Exception e){
                message=e.getMessage()==null?"Microsoft no pudo completar el inicio de sesión.":e.getMessage();
            }
            final boolean success=ok;
            final String notice=message;
            runOnUiThread(()->{
                if(webView!=null){
                    webView.evaluateJavascript("if(window.cineMicrosoftAuthFinished)window.cineMicrosoftAuthFinished("
                        +success+","+JSONObject.quote(notice)+")",null);
                }
                toast(notice);
            });
        });
    }
    private WebResourceResponse errorResponse() {return new WebResourceResponse("text/plain","utf-8",new ByteArrayInputStream(new byte[0]));}
    private class NativeBridge {
        @JavascriptInterface public void request(int id, String path, String method, String body) {
            workers.submit(()->{
                Backend.ApiResponse answer;
                try {answer=backend.handle(path,method,body);}
                catch(Exception e){answer=new Backend.ApiResponse(500,"{\"error\":"+JSONObject.quote(e.getMessage()==null?"Error desconocido":e.getMessage())+"}");}
                Backend.ApiResponse result=answer;
                String script="window.__cineNativeResolve("+id+","+result.status+","+JSONObject.quote(result.body)+","+JSONObject.quote(result.cloudHeader)+")";
                runOnUiThread(()->{if(webView!=null)webView.evaluateJavascript(script,null);});
            });
        }
        @JavascriptInterface public void exportFile(String name,String text,String mime) {
            runOnUiThread(()->{
                exportName=name;exportText=text;exportMime=mime;
                Intent intent=new Intent(Intent.ACTION_CREATE_DOCUMENT);
                intent.addCategory(Intent.CATEGORY_OPENABLE);
                intent.setType(mime==null?"text/plain":mime.split(";")[0]);
                intent.putExtra(Intent.EXTRA_TITLE,name);
                try{startActivityForResult(intent,EXPORT_FILE);}catch(Exception e){toast("No se puede exportar el archivo");}
            });
        }
        @JavascriptInterface public void saveOnExit(String state){
            try{JSONObject obj=new JSONObject(state);if(Backend.validLibrary(obj))backend.storeState(obj);}catch(Exception ignored){}
        }
        @JavascriptInterface public void closeApp(){runOnUiThread(()->finish());}
        @JavascriptInterface public void openExternal(String href){runOnUiThread(()->{
            try{Uri uri=Uri.parse(href);if("https".equals(uri.getScheme()))startActivity(new Intent(Intent.ACTION_VIEW,uri));}
            catch(Exception ignored){}
        });}
        @JavascriptInterface public void openMicrosoftLogin(String authorizationUrl){runOnUiThread(()->{
            try{
                Uri uri=Uri.parse(authorizationUrl);
                if(!"https".equals(uri.getScheme())||!MICROSOFT_LOGIN_HOST.equals(uri.getHost())
                    ||!"/consumers/oauth2/v2.0/authorize".equals(uri.getPath())
                    ||!BuildConfig.MICROSOFT_CLIENT_ID.equals(uri.getQueryParameter("client_id")))
                    throw new Exception("Dirección de Microsoft no reconocida.");
                startActivity(new Intent(Intent.ACTION_VIEW,uri));
            }catch(Exception e){toast("No se pudo iniciar sesión en Microsoft: "+e.getMessage());}
        });}
        @JavascriptInterface public void copy(String text){runOnUiThread(()->{
            ClipboardManager c=(ClipboardManager)getSystemService(Context.CLIPBOARD_SERVICE);
            c.setPrimaryClip(ClipData.newPlainText("película",text));
        });}
    }
    String selectFolder() throws Exception {
        folderLatch=new CountDownLatch(1);
        pickedFolder="";folderError="";
        runOnUiThread(()->{
            Intent i=new Intent(Intent.ACTION_OPEN_DOCUMENT_TREE);
            i.addFlags(Intent.FLAG_GRANT_READ_URI_PERMISSION|Intent.FLAG_GRANT_WRITE_URI_PERMISSION|Intent.FLAG_GRANT_PERSISTABLE_URI_PERMISSION);
            try {startActivityForResult(i,PICK_FOLDER);}catch(Exception e){folderError="No se pudo abrir el selector de carpetas";folderLatch.countDown();}
        });
        if(!folderLatch.await(120,TimeUnit.SECONDS)) throw new Exception("Tiempo agotado al elegir carpeta");
        if(!folderError.isEmpty())throw new Exception(folderError);
        return pickedFolder;
    }
    @Override protected void onActivityResult(int requestCode,int resultCode,Intent data){
        super.onActivityResult(requestCode,resultCode,data);
        if(requestCode==OPEN_FILE){
            if(fileCallback!=null){fileCallback.onReceiveValue(resultCode==RESULT_OK && data!=null && data.getData()!=null?new Uri[]{data.getData()}:null);fileCallback=null;}
        } else if(requestCode==EXPORT_FILE){
            if(resultCode==RESULT_OK && data!=null && data.getData()!=null){
                try(java.io.OutputStream os=getContentResolver().openOutputStream(data.getData(),"wt")){
                    if(os!=null)os.write(exportText.getBytes(StandardCharsets.UTF_8));
                    toast("Archivo guardado");
                }catch(Exception e){toast("Error al guardar: "+e.getMessage());}
            }
            exportText=null;
        } else if(requestCode==PICK_FOLDER){
            if(resultCode==RESULT_OK && data!=null && data.getData()!=null){
                Uri uri=data.getData();
                try{
                    getContentResolver().takePersistableUriPermission(uri,Intent.FLAG_GRANT_READ_URI_PERMISSION|Intent.FLAG_GRANT_WRITE_URI_PERMISSION);
                    pickedFolder=uri.toString();
                }catch(Exception e){folderError="No se han concedido permisos de escritura";}
            }
            if(folderLatch!=null)folderLatch.countDown();
        }
    }
    private void toast(String message){Toast.makeText(this,message,Toast.LENGTH_LONG).show();}
    @Override protected void onResume(){
        super.onResume();
        // Check OneDrive when the user returns from another app, without any TMDB refresh.
        if(webView!=null)webView.evaluateJavascript(
            "if(window.cineSyncResume)window.cineSyncResume()",null);
    }
    @Override public void onConfigurationChanged(Configuration newConfig){super.onConfigurationChanged(newConfig);}
    @Override public void onBackPressed(){
        // Do not discard an in-progress list edit when pressing Android Back.
        if(webView!=null)webView.evaluateJavascript("(function(){"+
                "var ov=document.querySelector('.overlay:not([hidden])');"+
                "if(ov){var close=ov.querySelector('[id^=close], [id^=cancel]');if(close)close.click();else ov.hidden=true;}"+
                "else if(document.querySelector('.nav-btn.active')?.dataset.view!=='library'){document.querySelector('.nav-btn[data-view=library]').click();}"+
                "else{window.CineAndroid.closeApp();}"+
                "})()",null);
    }
    @Override protected void onDestroy(){
        workers.shutdownNow();
        if(webView!=null){webView.destroy();webView=null;}
        super.onDestroy();
    }
}
