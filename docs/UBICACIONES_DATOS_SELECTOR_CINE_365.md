# Rutas de datos de Selector Cine 365

| Sistema | Biblioteca | Caché y sesión |
| --- | --- | --- |
| Windows 9.2.26 | %APPDATA%\Selector_Cine_365\ | %LOCALAPPDATA%\Selector_Cine_365\cache\ (caché regenerable). |
| macOS 9.2.26 | ~/Library/Application Support/Selector_Cine_365/ | ~/Library/Caches/Selector_Cine_365/ |
| Copia adicional | Carpeta elegida / Selector_Cine_365_Copias / identificador de instalación | colecciones.json, colecciones_anterior.json y hasta 14 diarias |

En escritorio la biblioteca se guarda en colecciones.json. Pueden existir copias en copias_automaticas y sync_recuperacion. Algunas instalaciones conservan un TXT de compatibilidad. La carpeta histórica AtlasDelCine puede contener diagnósticos y componentes de la ventana, como el cargador de WebView2 en Windows.


[Manual completo](MANUAL_COMPLETO_SELECTOR_CINE_365.md) · [Guía de recuperación](GUIA_SINCRONIZACION_SEGURA.md)

## Registros de diagnóstico (escritorio 9.2.26)

- Windows: `%APPDATA%\Selector_Cine_365\logs`.
- macOS: `~/Library/Application Support/Selector_Cine_365/logs`.
- Exportación: Mis listas → Mantenimiento → Exportar registro de diagnóstico.

## Recuperación de cambios pendientes

Los borradores de guardado están en `sync_recuperacion`.

## Historial de actividad

El historial de hasta 10.000 movimientos se guarda en `historial_local/`, en instantáneas identificadas por SHA-256. `colecciones.json` apunta a la instantánea confirmada y `sync_dispositivo_v2.json` conserva el espejo de recuperación. Las copias locales y las exportaciones completas incluyen el historial. No copies únicamente `colecciones.json` para trasladar la biblioteca: usa Exportar copia desde la app.

La clave TMDB queda protegida mediante DPAPI en Windows o el llavero de macOS. Mantenimiento muestra los tamaños de las carpetas y permite eliminar duplicados de caché verificados.
