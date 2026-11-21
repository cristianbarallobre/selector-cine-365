# Ubicaciones de datos — Selector Cine 365 (Beta 14 / Android 13.8)

A partir de estas compilaciones el nombre único de carpeta de datos es **`Selector_Cine_365`**.

| Plataforma | Carpeta |
| --- | --- |
| Windows | `%APPDATA%\Selector_Cine_365` |
| macOS | `~/Library/Application Support/Selector_Cine_365` |
| Android | `getFilesDir()/Selector_Cine_365` (almacenamiento privado de `es.cine365.app`) |
| OneDrive | `Selector_Cine_365` en la raíz de OneDrive (con el subdirectorio de sincronización existente) |

No hay migración ni lectura de la carpeta antigua de escritorio `AtlasDelCine`. Si existe, no se borra automáticamente: la nueva versión la ignora. El usuario que quiera una instalación completamente limpia puede borrar esa carpeta anterior por su cuenta.

**Android:** la aplicación conserva su identificador `es.cine365.app` para que las APK sigan actualizando el mismo paquete. La carpeta privada nueva se crea automáticamente; no se importa el contenido del directorio privado anterior. Las preferencias de autorización Microsoft se gestionan por separado mediante el almacenamiento de credenciales de la aplicación, no dentro de esta carpeta de biblioteca.

Al cambiar de nombre de carpeta en escritorio, el perfil de WebView2, las cachés de TMDB, la biblioteca y los datos de sincronización se generan nuevamente. **Los datos previos no son recuperados por la nueva versión salvo importación manual.**

La numeración visual y el arrastre de películas siguen vigentes; el orden compartido se almacena en la biblioteca dentro de la nueva ubicación.
