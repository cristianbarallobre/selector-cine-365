# Estado de Selector Cine 365

Windows y macOS 9.2.7 comparten el código Go y la interfaz web de escritorio. Las ventanas nativas usan WebView2 y WKWebView respectivamente. Android mantiene una interfaz táctil y un backend Java propios.

- **Escritorio vigente.** Biblioteca, búsqueda global, favoritos con orden manual, filtros, Descubrir, selección múltiple de plataformas, estadísticas y sincronización opcional. [Publicación 9.2.7](https://github.com/cristianbarallobre/selector-cine-365/releases/tag/v9.2.7-desktop)
- **Android vigente.** APK Limpia de la publicación v9.0.0. [Descarga](https://github.com/cristianbarallobre/selector-cine-365/releases/tag/v9.0.0)
- **Android de pruebas.** La candidata 9.2.4 usa un paquete separado. Se comprobó con pruebas Java y recorridos táctiles en Chromium. Sigue pendiente comprobar esa revisión en Android físico.

La publicación de escritorio verificó la actualización de Windows desde 9.2.0 y el arranque en Mac Intel. Apple Silicon se compiló para ARM64 con verificación de firma, arquitectura y contenido del ZIP. Sigue pendiente el arranque de esta versión en un Mac con chip Apple.

[Manuales y archivo](README.md) · [Pruebas en GitHub Actions](https://github.com/cristianbarallobre/selector-cine-365/actions/runs/45642426229)
