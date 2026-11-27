# Estado de Selector Cine 365

Windows y macOS 9.2.26 comparten el código Go y la interfaz web de escritorio. Las ventanas nativas usan WebView2 y WKWebView respectivamente.

- **Escritorio vigente.** Biblioteca, búsqueda global, favoritos con orden manual, filtros, Descubrir, selección múltiple de plataformas, estadísticas y guardado local, copia adicional opcional e historial de actividad. [Publicación 9.2.26](https://github.com/cristianbarallobre/selector-cine-365/releases/tag/v9.2.26-desktop)

GitHub Actions compila Windows x64 y macOS Intel y Apple Silicon. El guardado local y la copia adicional se prueban con carpetas temporales. El historial es independiente de la biblioteca; no existe conexión a cuentas externas. Se validan también las credenciales TMDB, caducidad de caché, rendimiento e interfaz.

[Manuales y archivo](README.md) · [Pruebas en GitHub Actions](https://github.com/cristianbarallobre/selector-cine-365/actions/workflows/desktop-release.yml)
