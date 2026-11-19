# Selector de Cine 365 en macOS

**Autor:** Cristian Barallobre Rúa.

Se ofrece para **Apple Silicon (M1 y posteriores)** y **Mac Intel**. Ambos paquetes llevan
una app nativa para macOS con ventana de WKWebView, sin dependencia de Microsoft WebView2.

## Instalación

1. Descarga el ZIP de tu procesador desde **Releases** en GitHub.
2. Descomprímelo y mueve **Selector de Cine 365.app** a **Aplicaciones**.
3. Ábrelo desde Aplicaciones. La primera vez macOS podría mostrar una advertencia
   porque no está firmado ni notarizado con un certificado Apple Developer ID.
4. Si macOS lo bloquea, ve a **Ajustes del Sistema → Privacidad y seguridad**
   y utiliza **Abrir igualmente** para esta aplicación, solo si confías en su origen.

**AVISO:** El programa se firma de forma ad hoc para comprobar la integridad de
los binarios generados, pero **no** cuenta con una firma Developer ID ni con
notarización de Apple. Se recomienda comprobar el repositorio y la procedencia.

Las listas se guardan en `~/Library/Application Support/AtlasDelCine/`. No se
sobrescriben datos de Windows: la instalación de macOS guarda su propia biblioteca.
El botón de copias en la nube puede seleccionar una carpeta local sincronizada,
por ejemplo iCloud Drive, Dropbox o Google Drive.

## Desarrollo

Se requiere un Mac con Xcode Command Line Tools y Go 1.23+.

```bash
CGO_ENABLED=1 go test ./...
CGO_ENABLED=1 go build -o SelectorCine365 .
```

`window_darwin.go` y `macos_view_darwin.m` implementan la ventana Cocoa/WebKit.
`cloud_picker_darwin.go` implementa el selector de carpetas de macOS.

GitHub Actions compila dos variantes de `.app` en runners macOS Intel y ARM64
reales. Se verifica compilación y empaquetado, pero el flujo de publicación no
hace una prueba interactiva de la interfaz. Agradecemos informes de uso real.
