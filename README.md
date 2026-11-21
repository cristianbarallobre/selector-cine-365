# Selector Cine 365 🎬

Aplicación personal de cine multiplataforma para Windows, macOS y Android. Biblioteca con listas, películas vistas, estadísticas, TMDB opcional y sincronización opcional mediante OneDrive.

## Descarga recomendada para pruebas (21 de noviembre de 2026)

**[Selector Cine 365 — Release conjunta v9 Beta 14 / Android Beta 13.13](https://github.com/cristianbarallobre/selector-cine-365/releases/tag/v9.0.0-beta.14)**

| Plataforma | Archivo de esta release | Compatibilidad |
|---|---|---|
| Windows 10/11 x64 | `Selector_Cine_365_Windows_Beta14.exe` | EXE portátil, WebView2 |
| macOS Apple Silicon | `Selector_Cine_365_Mac_AppleSilicon_Beta14.zip` | macOS 12+, sin notarización |
| macOS Intel | `Selector_Cine_365_Mac_Intel_Beta14.zip` | macOS 12+, sin notarización |
| Android 8+ | `Selector_Cine_365_Android_Beta13_13_Limpia.apk` | `es.cine365.app.limpia`, versionCode 9027, firma original |

**Estado: beta recomendada para pruebas; no estable definitiva.** Se han pasado pruebas automatizadas, pero faltan comprobaciones de cuentas Microsoft y ciertos dispositivos. La seguridad de sincronización ha mejorado: se bloquean uniones de listas con ID coincidente pero bibliotecas diferentes. **No repara listas que ya se mezclaron:** conserva tus copias JSON y no sincronices si ves conteos incorrectos.

**Actualización Android:** instala la APK Beta 13.13 encima de la **Beta 13.12 Limpia**, sin desinstalar; solo son compatibles si conservan `es.cine365.app.limpia` y la misma firma. La línea antigua `es.cine365.app` tiene otro identificador y no puede actualizarse directamente con esta APK.

**TMDB:** carátulas, valoraciones y plataformas son cachés locales de cada dispositivo; las identificaciones manuales (ID TMDB) sí pueden viajar con las listas. Una carátula visible en Windows no tiene por qué estar descargada aún en Android.

### Documentación vigente

- [Manual actualizado de las versiones recomendadas](docs/MANUAL_BETA14_ANDROID13_13.md)
- [Guía de sincronización y recuperación segura](docs/GUIA_SINCRONIZACION_SEGURA_BETA14_ANDROID13_13.md)
- [Historial de versiones](CHANGELOG.md) y [política de publicaciones](docs/POLITICA_DE_RELEASES.md)
- [Notas completas](RELEASE_NOTES_v9_BETA14_RECOMENDADA.md)
- [Código fuente unificado](https://github.com/cristianbarallobre/selector-cine-365/tree/release/selector-cine-v9-beta14-beta13-13)

Las versiones anteriores se conservan en [GitHub Releases](https://github.com/cristianbarallobre/selector-cine-365/releases) y, cuando corresponda, en el archivo histórico privado.

---

## Archivo: cambios de betas anteriores

*Las notas que siguen documentan estados históricos, no versiones actuales recomendadas.*

## Android v9 Beta 13.5 — OneDrive sin selección obligatoria de carpetas

- Al conectar Microsoft, Android reutiliza la carpeta Graph previamente vinculada o localiza/crea automáticamente `Selector_Cine_365` en la raíz de OneDrive, **sin publicar ni importar datos en ese paso**.
- Las carpetas antiguas SAF no se migran a escondidas; el selector queda en **Opciones avanzadas** para una transición controlada.
- Para enviar y recibir cambios hay que pulsar ↻, revisar y confirmar; se mantienen los guardados locales.
- Compatible con la misma firma Android de Beta 13.3; **`versionCode 9019`**.
- **Límite:** Windows/macOS aún usan la carpeta local de OneDrive y deben apuntar a `Selector_Cine_365` para compartir biblioteca; no tienen inicio Microsoft Graph directo.
- [Notas técnicas y pasos de comprobación](RELEASE_NOTES_v9_BETA13_5_ANDROID.md). **La APK Beta 13.5 todavía no está publicada en GitHub Releases; no confundir el artefacto CI sin firma con una APK instalable.**

## Novedad Android Beta 13.3 — botón de sincronización manual

- Al marcar películas se guardan **en el móvil**; la app ya no sube datos a OneDrive por cada cambio o cambio de lista.
- Al abrir consulta la nube una vez en modo **solo lectura**, sin importar ni enviar automáticamente.
- El pequeño botón de dos flechas circulares junto a «Guardado» permite revisar **envío y recepción** de cambios, con confirmación previa y gestión de conflictos.
- No hay sincronización periódica ni al volver desde segundo plano. Se reducen escrituras de datos idénticos.
- Android `es.cine365.app`, versión interna **9017**, misma firma que Beta 13.2. Windows y macOS mantienen Beta 13.1.
- [Notas completas Beta 13.3](RELEASE_NOTES_v9_BETA13_3_ANDROID.md).

## Hotfix Android Beta 13.2 — batería y segundo plano

Android reduce comprobaciones automáticas de OneDrive de 30 a **60 segundos en primer plano**, deja de hacer el ping de 12 segundos, pausa la WebView al minimizar y elimina la inicialización doble del backend. Se conservan el guardado inmediato al minimizar con un cambio pendiente, las revisiones causales y la sincronización manual. **APK Android 9016**, mismo paquete y firma que Beta 13. [Descargar Beta 13.2](https://github.com/cristianbarallobre/selector-cine-365/releases/tag/v9.0.0-beta.13.2). Windows y Mac continúan en Beta 13.1. [Notas del hotfix](RELEASE_NOTES_v9_BETA13_2_ANDROID.md).

## Novedades v9 Beta 13 — listas y TMDB fiables

- **Windows y Mac:** siempre se reserva espacio para cuatro listas en cada página; no cambia la altura al ir a la última.
- **Borrar listas sí se sincroniza:** un registro persistente de eliminación impide que copias de otro dispositivo las vuelvan a crear.
- **Identificar con el lápiz:** el ID manual TMDB viaja con las listas entre Windows, Android y Mac. La clave API personal sigue siendo privada de cada dispositivo.
- **Al reabrir:** una lectura fallida temporal de TMDB no borra la configuración y la app vuelve a intentar conectar con su motor local.
- **Android:** mantiene el notch y el menú corregidos en Beta 12.2. Paquete `es.cine365.app`, versión interna 9015, firma de las Betas anteriores.

**Antes de probar OneDrive:** exporta JSON en ambos dispositivos y actualiza los dos a Beta 13. Si una lista ya había reaparecido antes de esta revisión, elimínala una vez con ambos actualizados. [Ver la release conjunta Beta 13.1](https://github.com/cristianbarallobre/selector-cine-365/releases/tag/v9.0.0-beta.13.1).

## Android v9 Beta 12.2 — notch y menú inferior sin margen doble

Beta 12.1 solucionó el título tapado por el notch, pero en algunos móviles Redmi el menú inferior quedaba elevado y aparecía una franja vacía. **Beta 12.2** corrige la doble reserva del espacio de navegación: Android aplica los márgenes seguros reales y el CSS ya no los suma otra vez. La APK firmada **9014** actualiza Beta 11, 12 y 12.1 sin desinstalar ni perder datos; conserva el mismo identificador y certificado. Desde Beta 9 sigue siendo necesaria una migración por JSON.

[Descargar Android Beta 12.2](https://github.com/cristianbarallobre/selector-cine-365/releases/tag/v9.0.0-beta.12.2). Windows y Mac utilizan los **mismos ZIP** de Beta 12.

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
- [Historial de versiones completo v1–v9 Beta 13.1](docs/HISTORIAL_V9_BETA13_1.md) — Markdown, Word y PDF.
- [Guía Mac: Intel y Apple Silicon](MACOS.md).
- [Guía Android: instalación, firma y OneDrive](android/LEEME_ANDROID_v9_BETA.md).
- [Azure / Entra ID / OneDrive](docs/GUIA_AZURE_ENTRA_ONEDRIVE_SELECTOR_CINE_365.md).
- [Plan técnico v9 y límites de sincronización](docs/PLAN_V9_MULTIPLATAFORMA_Y_SINCRONIZACION.md).
- [Cambios recientes](CHANGELOG.md) y [archivo histórico hasta v8.1.1](docs/HISTORIAL_VERSIONES.md).
- [Autoría](AUTORIA.md).

Los Word y PDF de `v8.1.1` documentan **esa versión histórica**. El nuevo historial v9 se encuentra en `docs/HISTORIAL_V9_BETA13_1.md` y sus ediciones Word y PDF; no sustituyen el manual de uso.

## Derechos y atribuciones

© 2026 Cristian Barallobre Rúa. Todos los derechos reservados. Este repositorio no incorpora una licencia pública de reutilización de código; no presupongas permisos que no se hayan otorgado.

**This product uses the TMDB API but is not endorsed or certified by TMDB.** Los datos de disponibilidad pueden proceder de TMDB / JustWatch y cambiar sin aviso. El proyecto no guarda las listas en TMDB.
