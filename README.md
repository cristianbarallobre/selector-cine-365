# Selector Cine 365 🎬

Aplicación personal de cine para **Windows 10/11**, **macOS Apple Silicon e Intel** y **Android**, creada por **Cristian Barallobre Rúa**. Sirve para organizar listas, registrar películas vistas, decidir qué ver, consultar datos de TMDB y sincronizar bibliotecas entre dispositivos.

**Versión coordinada más reciente:** **v9 Beta 12 (20 de noviembre de 2026)**. Es una *beta*, no una edición estable o validada en todos los dispositivos. Las versiones históricas principales permanecen visibles; los instaladores de las betas intermedias se han respaldado en un borrador privado, conservando las etiquetas Git y las ramas.

## Descargar

**[Releases y descargas de Selector Cine 365](https://github.com/cristianbarallobre/selector-cine-365/releases)**

Las **cuatro releases públicas** se organizan así: [Beta 12](https://github.com/cristianbarallobre/selector-cine-365/releases/tag/v9.0.0-beta.12) (actual, cuatro plataformas), [Beta 9](https://github.com/cristianbarallobre/selector-cine-365/releases/tag/v9.0.0-beta.9) (histórica, Windows y Android), [v8.2](https://github.com/cristianbarallobre/selector-cine-365/releases/tag/v8.2) y [v8.1.1](https://github.com/cristianbarallobre/selector-cine-365/releases/tag/v8.1.1). Las releases redundantes están archivadas; **no se han borrado sus fuentes ni sus tags**. Consulta la [política de publicaciones y recuperación](docs/POLITICA_DE_RELEASES.md).

| Plataforma | Archivo de la Beta 12 | Requisitos |
|---|---|---|
| Windows | `Selector_Cine_365_Windows_v9_BETA12.zip` | Windows 10/11 de 64 bits y Microsoft WebView2 Runtime |
| Mac Apple Silicon | `Selector_Cine_365_Mac_Apple_Silicon_v9_BETA12.zip` | macOS 12+; M1/M2/M3/M4 y posteriores |
| Mac Intel | `Selector_Cine_365_Mac_Intel_v9_BETA12.zip` | macOS 12+; procesador Intel |
| Android | `Selector_Cine_365_Android_v9_BETA12.apk` | Android 8.0+ (API 26); instalar únicamente la APK firmada publicada |

**Importante:** si una plataforma aún no figura entre los archivos de la release, no descargues una APK sin firma ni sustituyas otra versión como si fuera equivalente. Consulta las notas de la publicación.

## Novedades v9 Beta 12

- **Buscar en la lista actual o en todas:** resultado global agrupado por identidad TMDB o título y año, sin repetir fichas de la misma película.
- **Listas numeradas y navegables:** «Listas: 1, 2, 3» en cada ficha; pulsar el número abre la lista. Desde el botón de edición puedes añadir o retirar una película de otras listas.
- **Marcar Vista/Pendiente:** solo el tick cambia el estado; los carteles y títulos no deben hacerlo. Si la misma película está identificada en varias listas, el estado se comparte.
- **Carteles, año y valoración TMDB**, disponibilidad en España y enlaces a proveedores cuando haya datos. La nota mostrada es de **TMDB, no IMDb**.
- **Windows y Mac:** misma interfaz web incrustada en aplicaciones nativas; Estadísticas con paginación fija en la parte inferior.
- **Windows Beta 12:** corregidas las filas invisibles en «Mis listas» cuando Nube ocupa espacio; «Actualizar historial» comunica el resultado de la consulta y muestra discretamente el estado vacío.
- **Android:** cabecera y navegación propias, desplazamiento en las secciones, biblioteca con carga progresiva y estadísticas sin paginación.
- **Nube opcional:** almacenamiento local sin OneDrive y sin clave TMDB. Si conectas la nube, hay revisión inicial, sincronización periódica, revisiones causales, historial y recuperación de copias.

**La sincronización es experimental:** depende de OneDrive, no es instantánea y pide revisión manual ante modificaciones concurrentes no ordenables. No prometemos que siempre gane la hora del reloj: el historial causal y las revisiones determinan la precedencia.

## Actualizar sin perder datos

**Antes de instalar cualquier beta, exporta un JSON de la biblioteca y comprueba que puedes abrirlo.** Conserva una segunda copia fuera de la carpeta de sincronización.

- **Windows:** cierra la versión anterior, sustituye el ejecutable y no abras dos versiones a la vez; el almacén local no cambia.
- **Mac:** cierra la aplicación anterior y reemplaza la `.app` en Aplicaciones. Conserva los datos de `~/Library/Application Support/AtlasDelCine/`. La app Mac tiene firma *ad hoc* y **no está notarizada**, por lo que macOS puede solicitar autorización manual.
- **Android Beta 11 → Beta 12:** ambas usan el identificador `es.cine365.app` y requieren el **mismo certificado privado de firma**. Con la APK publicada firmada correctamente, la Beta 12 reemplaza a Beta 11 y conserva sus datos.
- **Android Beta 9 → Beta 12:** Beta 9 usaba `es.cine365.app.v9beta9` y una firma distinta. No existe actualización directa con conservación automática de datos. Exporta JSON desde Beta 9, verifica el archivo, desinstala la versión antigua solo después, instala Beta 12 e importa el JSON. Comprueba todo antes de conectar OneDrive.

**Nunca subas claves privadas Android (keystore/P12), contraseñas, tokens OAuth, claves TMDB ni datos personales al repositorio o a sus logs.** En versiones futuras de Android conserva el identificador y el certificado, sin sufijos de paquete por beta.

## Desarrollo

Fuentes compartidas para escritorio: **Go 1.23+** y `plantilla.html`; **Windows** utiliza WebView2 y **macOS** utiliza WKWebView. **Android** dispone de proyecto Gradle y Java propio, con la interfaz adaptada en `android/app/src/main/assets/index.html`.

### Windows

```bash
go test ./...
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags="-H windowsgui -s -w" -o Selector_Cine_365_Windows.exe .
```

### macOS

Requiere un Mac con Xcode Command Line Tools, Go y las herramientas de firma. Para el empaquetado de Intel/ARM se utiliza `scripts/build_macos_app.sh` desde un runner macOS (consulta `MACOS.md`). Ambas variantes compilan **el mismo código v9 de Windows**, salvo la envoltura nativa.

### Android

```bash
cd android
gradle testDebugUnitTest assembleRelease
```

La APK release de Gradle se genera sin firma si no existe configuración externa: **no es instalable ni apta para distribución**. Debe firmarse mediante un proceso seguro con una clave privada persistente y verificarse antes de publicar. No incluir esa clave en Git.

## Almacenamiento y nube

| Sistema | Carpeta privada habitual |
|---|---|
| Windows | `%APPDATA%\\AtlasDelCine\\` |
| macOS | `~/Library/Application Support/AtlasDelCine/` |
| Android | Memoria interna privada de la aplicación |

Los nombres históricos `AtlasDelCine` y `atlas` permanecen **solo como identificadores de compatibilidad**, aunque la marca visible sea Selector Cine 365. No los renombres ni muevas manualmente sin un procedimiento de migración comprobado.

En Windows y Mac puedes seleccionar una carpeta local sincronizada por el proveedor admitido; Android dispone de conexión a **OneDrive personal** con Microsoft OAuth2 + PKCE mediante el navegador. La primera sincronización debe revisarse y confirmarse en cada dispositivo. Al usar una carpeta con datos reales, no publiques una biblioteca vacía sobre otra.

## Documentación

- [Manual completo v9 Beta 12](docs/MANUAL_COMPLETO_SELECTOR_CINE_365.md) (también PDF).
- [Guía Mac: Intel y Apple Silicon](MACOS.md).
- [Guía Android: instalación, firma y OneDrive](android/LEEME_ANDROID_v9_BETA.md).
- [Azure / Entra ID / OneDrive](docs/GUIA_AZURE_ENTRA_ONEDRIVE_SELECTOR_CINE_365.md).
- [Plan técnico v9 y límites de sincronización](docs/PLAN_V9_MULTIPLATAFORMA_Y_SINCRONIZACION.md).
- [Cambios recientes](CHANGELOG.md) y [archivo histórico hasta v8.1.1](docs/HISTORIAL_VERSIONES.md).
- [Autoría](AUTORIA.md).

Los PDF y Word de `v8.1.1` documentan **esa versión histórica**, no constituyen el manual de v9.

## Derechos y atribuciones

© 2026 Cristian Barallobre Rúa. Todos los derechos reservados. Este repositorio no incorpora una licencia pública de reutilización de código; no presupongas permisos que no se hayan otorgado.

**This product uses the TMDB API but is not endorsed or certified by TMDB.** Los datos de disponibilidad pueden proceder de TMDB / JustWatch y cambiar sin aviso. El proyecto no guarda las listas en TMDB.
