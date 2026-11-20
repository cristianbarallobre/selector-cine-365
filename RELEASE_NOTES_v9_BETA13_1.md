# Selector Cine 365 — v9 Beta 13.1 (Windows, Android y macOS)

**20 de noviembre de 2026 · Versión preliminar para pruebas.** Esta publicación reúne las versiones actuales para las cuatro plataformas. Son **Windows v9 Beta 13.1**, **Mac Intel/Apple Silicon v9 Beta 13.1** y la **APK Android v9 Beta 13 original**, con firma verificada y `versionCode 9015`. No existe una recompilación Android «Beta 13.1»: el arreglo de esta revisión afecta exclusivamente al backend Go de escritorio.

## Corrección principal de Beta 13.1: TMDB en escritorio

- Windows/Mac Beta 13 podían quedarse en «Conexión local TMDB no disponible. Reintentando…» después de cerrar y volver a abrir la aplicación.
- **Causa comprobada:** el frontend solicitaba `GET /api/platforms/cache`, que el módulo de TMDB admitía, pero que no estaba registrado en el router HTTP principal. Respondía **404** aun estando la clave API guardada en AppData.
- **Solución:** ruta registrada, lectura de estado separada de la caché y botón Actualizar capaz de recuperar una consulta fallida. Si la caché falla no se borra la clave ni se indica erróneamente que no está configurada.
- Se añadieron pruebas del servidor (estado/caché responden 200 con token; 403 sin token). El ejecutable Windows se compiló como aplicación GUI, sin consola. Los Mac se reconstruyeron desde el mismo código corregido.

## Funciones de Beta 13 conservadas (todas las plataformas)

- «Mis listas» ocupa siempre el espacio de cuatro filas en Windows/Mac aunque la última página tenga menos listas.
- Las eliminaciones llevan registros `deleted_lists` compartidos, para que copias antiguas de otros dispositivos **no resuciten listas borradas**.
- Las identificaciones manuales con el lápiz guardan IDs de TMDB en `tmdb_overrides`, que pueden sincronizarse. **La clave API personal jamás se incluye en la biblioteca compartida**.
- Persistencia local de listas, marcas de visto y preferencias, funcionamiento sin nube y sincronización causal de las marcas.
- Android conserva la adaptación a notch y navegación por gestos de Beta 12.2; sin cambio de identificador ni certificado de firma.

## Descargas de esta release

| Plataforma | Archivo |
|---|---|
| Windows 10/11 x64 | `Selector_Cine_365_Windows_v9_BETA13_1.zip` |
| Mac Apple Silicon | `Selector_Cine_365_Mac_Apple_Silicon_v9_BETA13_1.zip` |
| Mac Intel | `Selector_Cine_365_Mac_Intel_v9_BETA13_1.zip` |
| Android 8.0+ | `Selector_Cine_365_Android_v9_BETA13.apk` (APK firmada, versión interna 9015) |

Incluye además **Manual completo PDF**, **Guía Microsoft Entra/OneDrive PDF** e **Historial completo v9 Beta 13.1 en PDF y Word**.

## Instalación y datos

**Antes de actualizar, exporta un JSON en cada dispositivo y guarda una copia fuera de la carpeta OneDrive compartida.** Cierra Windows/Mac y sustituye el ejecutable o el paquete `.app`. Android Beta 11/12/12.1/12.2 puede actualizarse directamente instalando **Beta 13 encima, sin desinstalar**, pues conserva `es.cine365.app` y el certificado persistente. La Beta 9 Android tenía otro identificador y firma: requiere migración de datos mediante copia JSON.

**Sincronización:** actualiza todos los dispositivos vinculados antes de volver a borrar listas. Los borrados realizados antes de Beta 13 no incluían marcadores: si una lista histórica reaparece, bórrala una vez con ambos dispositivos ya actualizados. La sincronización por OneDrive no es instantánea; pueden surgir conflictos concurrentes que requieran revisión.

**TMDB:** tu clave se guarda localmente en Windows/Mac/Android, no se incorpora en el instalador. No la subas a GitHub. Los IDs corregidos manualmente pueden viajar con la biblioteca; para descargar carteles y datos cada dispositivo necesita su propia clave válida.

**Mac:** firma ad hoc, **sin notarización Apple Developer ID**. Es posible que Gatekeeper solicite autorizar su apertura. Compilación CI verificada, pero no sustituye pruebas manuales en hardware físico.

## Pruebas y trazabilidad

- [Windows Beta 13.1: compilación y pruebas de caché HTTP](https://github.com/cristianbarallobre/selector-cine-365/actions/runs/44829805619).
- [macOS Beta 13.1 Intel/Apple Silicon](https://github.com/cristianbarallobre/selector-cine-365/actions/runs/44830855769).
- [Android Beta 13 y sincronización causal](https://github.com/cristianbarallobre/selector-cine-365/actions/runs/44826692958).
- Huella SHA-256 de APK Android firmada: `0cff436f96b15c324169e4626746fd5a905d0c26ba5d9ddc769b81ed51f5eabe`.
- Certificado Android SHA-256: `b8a17f39e4fa24d7c110d452161425c7524412c3a5a4d7459b1132ffd0825fb6`.

**Limitación:** falta completar una prueba de uso prolongado de OneDrive entre Windows/Mac y Android reales, incluidos borrados simultáneos. Esta versión sigue etiquetada como *pre-release*.

[Manual](docs/MANUAL_COMPLETO_SELECTOR_CINE_365.md) · [Historial actualizado](docs/HISTORIAL_V9_BETA13_1.md) · [Guía Android](android/LEEME_ANDROID_v9_BETA.md) · [GitHub Releases](https://github.com/cristianbarallobre/selector-cine-365/releases)
