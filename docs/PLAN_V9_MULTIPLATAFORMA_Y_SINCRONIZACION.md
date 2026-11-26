# Estado de Selector Cine 365

Windows y macOS 9.2.19 comparten el código Go y la interfaz web de escritorio. Las ventanas nativas usan WebView2 y WKWebView respectivamente. Android mantiene una interfaz táctil y un backend Java propios.

- **Escritorio vigente.** Biblioteca, búsqueda global, favoritos con orden manual, filtros, Descubrir, selección múltiple de plataformas, estadísticas y sincronización manual opcional con OneDrive. [Publicación 9.2.19](https://github.com/cristianbarallobre/selector-cine-365/releases/tag/v9.2.19-desktop)
- **Android vigente.** APK Limpia de la publicación v9.0.0. [Descarga](https://github.com/cristianbarallobre/selector-cine-365/releases/tag/v9.0.0)
- **Android de pruebas.** La candidata 9.2.4 usa un paquete separado. Se comprobó con pruebas Java y recorridos táctiles en Chromium. Sigue pendiente comprobar esa revisión en Android físico.

La sincronización manual se probó en Windows con una carpeta de prueba y cambios en ambos lados. GitHub Actions compila Windows x64 y macOS Intel y Apple Silicon. Sigue pendiente probar el flujo completo en macOS con una cuenta de OneDrive.

[Manuales y archivo](README.md) · [Pruebas en GitHub Actions](https://github.com/cristianbarallobre/selector-cine-365/actions/workflows/desktop-release.yml)
