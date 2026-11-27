# Rutas de datos de Selector Cine 365

| Sistema | Biblioteca | Caché y sesión |
| --- | --- | --- |
| Windows 9.2.25 | %APPDATA%\Selector_Cine_365\ | Cachés de la aplicación en su carpeta de datos de escritorio. |
| macOS 9.2.25 | ~/Library/Application Support/Selector_Cine_365/ | ~/Library/Caches/Selector_Cine_365/ |
| Copia adicional | Carpeta elegida / Selector_Cine_365_Copias / identificador de instalación | colecciones.json y colecciones_anterior.json |

En escritorio la biblioteca se guarda en colecciones.json. Pueden existir copias en copias_automaticas y sync_recuperacion. Algunas instalaciones conservan un TXT de compatibilidad. La carpeta histórica AtlasDelCine puede contener diagnósticos y componentes de la ventana, como el cargador de WebView2 en Windows.


[Manual completo](MANUAL_COMPLETO_SELECTOR_CINE_365.md) · [Guía de recuperación](GUIA_SINCRONIZACION_SEGURA.md)

## Registros de diagnóstico (escritorio 9.2.25)

- Windows: `%APPDATA%\Selector_Cine_365\logs`.
- macOS: `~/Library/Application Support/Selector_Cine_365/logs`.
- Exportación: Mis listas → Mantenimiento → Exportar registro de diagnóstico.

## Recuperación de cambios pendientes

Los borradores de guardado están en `sync_recuperacion`.

## Historial de actividad

El historial de hasta 10.000 movimientos se guarda junto al estado local de la biblioteca y su espejo de metadatos. No se incluye en las exportaciones JSON ordinarias ni en la copia adicional. Las copias locales de recuperación pueden contener versiones del historial.
