# Rutas de datos de Selector Cine 365

| Sistema | Biblioteca | Caché y sesión |
| --- | --- | --- |
| Windows 9.2.19 | %APPDATA%\Selector_Cine_365\ | Cachés de la aplicación en su carpeta de datos de escritorio. Sesión OneDrive protegida por DPAPI |
| macOS 9.2.19 | ~/Library/Application Support/Selector_Cine_365/ | ~/Library/Caches/Selector_Cine_365/. Sesión OneDrive en el llavero |
| Android estable | Almacenamiento privado de es.cine365.app.limpia | Las otras ediciones usan almacenamiento independiente |
| OneDrive | Carpeta lógica Selector_Cine_365 cuando se usa Graph | Las copias y revisiones dependen de la configuración y la vinculación |

En escritorio la biblioteca se guarda en colecciones.json. Pueden existir copias en copias_automaticas y sync_recuperacion. Algunas instalaciones conservan un TXT de compatibilidad. La carpeta histórica AtlasDelCine puede contener diagnósticos y componentes de la ventana, como el cargador de WebView2 en Windows.

Los carteles y los datos de disponibilidad se pueden volver a descargar desde TMDB. En escritorio, la carpeta de caché **imagenes_descubrir** guarda los carteles grandes y retratos. La sincronización con OneDrive se realiza al pulsar el botón. Conserva la biblioteca, las copias y las identificaciones manuales al limpiar la caché. Para mover datos entre equipos o ediciones Android, utiliza la exportación JSON.

[Manual completo](MANUAL_COMPLETO_SELECTOR_CINE_365.md) · [Guía de recuperación](GUIA_SINCRONIZACION_SEGURA.md)

## Registros de diagnóstico (escritorio 9.2.19)

- Windows: `%APPDATA%\Selector_Cine_365\logs`.
- macOS: `~/Library/Application Support/Selector_Cine_365/logs`.
- Exportación: Mis listas → Mantenimiento → Exportar registro de diagnóstico.
