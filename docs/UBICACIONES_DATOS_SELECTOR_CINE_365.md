# Ubicaciones de datos — Selector Cine 365 (21/11/2026)

| Plataforma | Datos privados del programa |
|---|---|
| Windows Beta 14 | `%APPDATA%\Selector_Cine_365` |
| macOS Beta 14 | `~/Library/Application Support/Selector_Cine_365` |
| Android Beta 13.13 Limpia | `getFilesDir()/Selector_Cine_365`, almacenamiento del paquete **`es.cine365.app.limpia`** |
| OneDrive | Carpeta `Selector_Cine_365` de la cuenta Microsoft, normalmente con `Selector_Cine_365_Copias/sincronizacion/dispositivos/*.json` |

Windows guarda la autorización Microsoft de renovación **cifrada mediante DPAPI** en `onedrive_sesion_dpapi.bin`, dentro de su carpeta de datos, nunca en la nube. Android protege sus credenciales mediante Android Keystore. Mac no persiste todavía su credencial Microsoft.

**Android histórico:** `es.cine365.app` es otro paquete; no comparte automáticamente datos privados ni se actualiza directamente con `es.cine365.app.limpia`. Ambas líneas pueden coexistir en un teléfono. Si desinstalas cualquiera de ellas, se elimina normalmente el almacenamiento privado de esa aplicación. OneDrive sigue intacto y la recuperación requiere revisión.

**Biblioteca:** `colecciones.json`; copias en `copias_automaticas` y `sync_recuperacion`, si existen. Las películas vistas y las revisiones causales se sincronizan con snapshots, no con credenciales. No borres una carpeta por encontrarla vacía antes de verificar la cuenta OneDrive activa y su contenido real.

**Migración:** no hay importación automática garantizada desde el nombre antiguo de escritorio `AtlasDelCine`; conserva los datos anteriores por separado y usa un JSON de exportación si los necesitas.

[Guía de recuperación](GUIA_SINCRONIZACION_SEGURA_BETA14_ANDROID13_13.md).
