# Estado de Selector Cine 365

Windows y macOS 9.2.24 comparten el código Go y la interfaz web de escritorio. Las ventanas nativas usan WebView2 y WKWebView respectivamente. Android mantiene una interfaz táctil y un backend Java propios.

- **Escritorio vigente.** Biblioteca, búsqueda global, favoritos con orden manual, filtros, Descubrir, selección múltiple de plataformas, estadísticas y guardado local y copia adicional opcional. [Publicación 9.2.24](https://github.com/cristianbarallobre/selector-cine-365/releases/tag/v9.2.24-desktop)
- **Android vigente.** APK Limpia de la publicación v9.0.0. [Descarga](https://github.com/cristianbarallobre/selector-cine-365/releases/tag/v9.0.0)
- **Android de pruebas.** La edición de pruebas 9.2.5 usa un paquete separado y requiere Android 8 o posterior. [Publicación de pruebas](https://github.com/cristianbarallobre/selector-cine-365/releases/tag/v9.2.5-android). No sustituye a la APK estable.

GitHub Actions compila Windows x64 y macOS Intel y Apple Silicon. El guardado local y la copia adicional se prueban con carpetas temporales. Android mantiene su backend anterior.

[Manuales y archivo](README.md) · [Pruebas en GitHub Actions](https://github.com/cristianbarallobre/selector-cine-365/actions/workflows/desktop-release.yml)
