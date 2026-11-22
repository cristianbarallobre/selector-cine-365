# Ubicaciones de datos — Selector Cine 365 v9.0.0

**Vigente el 22 de noviembre de 2026.** Estos directorios pertenecen al programa, no a los ejecutables que se descargan.

| Plataforma | Datos locales |
| --- | --- |
| Windows | `%APPDATA%\Selector_Cine_365` |
| macOS Intel / Apple Silicon | `~/Library/Application Support/Selector_Cine_365` |
| Android estable | Almacenamiento privado del paquete **`es.cine365.app.limpia`**, carpeta `Selector_Cine_365` |
| OneDrive (opcional) | Carpeta lógica `Selector_Cine_365` de la cuenta Microsoft personal conectada; snapshots en `Selector_Cine_365_Copias/sincronizacion/dispositivos/` cuando se utilice ese esquema |

**Archivos habituales en escritorio:** `colecciones.json` (biblioteca), respaldo `copias_automaticas` y/o `sync_recuperacion` (si se generaron). En Windows la sesión Microsoft tiene un archivo DPAPI denominado `onedrive_sesion_dpapi.bin`, ligado al usuario que autorizó; no se debe copiar a otros equipos ni publicar.

**Android:** `es.cine365.app.limpia` y `es.cine365.app` son **dos aplicaciones diferentes**, con cachés, datos y preferencias privadas separadas. La APK v9.0.0 actualiza Beta 13.12 Limpia con la misma firma; no instala una actualización de la aplicación histórica. Desinstalar Android elimina normalmente los datos privados del paquete afectado, **pero no borra OneDrive**. Para recuperar listas exporta JSON y revisa las copias en la nube.

**Cachés TMDB:** se guardan localmente por dispositivo; la carátula descargada en Windows no se comparte a Android. La clave personal de TMDB tampoco viaja con OneDrive.

**Nombre antiguo:** `AtlasDelCine` no se migra automáticamente a la carpeta actual; conserva los archivos antiguos que necesites por separado.

Consulta [Manual completo](MANUAL_COMPLETO_SELECTOR_CINE_365.md), [Sincronización segura](GUIA_SINCRONIZACION_SEGURA.md) y [Release estable](https://github.com/cristianbarallobre/selector-cine-365/releases/tag/v9.0.0).
