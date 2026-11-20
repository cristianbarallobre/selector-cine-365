# Selector Cine 365 — v9 Beta 12.2 (Android, notch y barra inferior)

**Actualización correctiva Android del 20 de noviembre de 2026.** Windows, macOS Intel y macOS Apple Silicon conservan **exactamente los mismos ejecutables de Beta 12**. No se ha cambiado el motor de sincronización, OneDrive, TMDB ni la gestión de listas.

## Corrección

En algunos móviles con notch o recorte de cámara, como determinados Redmi, **Beta 12.1** impidió que el título quedara tapado por la cámara frontal, pero la barra de seis pestañas inferiores quedaba demasiado arriba y aparecía una franja vacía. El contenedor nativo de Android aplicaba los márgenes dinámicos de `WindowInsets` y el CSS añadía otra vez `env(safe-area-inset-bottom)`.

**Beta 12.2 retira únicamente el margen duplicado** del menú, del cuerpo y del área principal. Conserva el contenedor seguro nativo para el notch, barras de estado, gestos y pantalla en horizontal. No utiliza un margen fijo según la marca o el modelo.

## Descargas

- **Android 8.0+ (APK firmada):** `Selector_Cine_365_Android_v9_BETA12_2.apk`, `es.cine365.app`, `versionCode=9014`.
- **Windows 10/11 x64:** `Selector_Cine_365_Windows_v9_BETA12.zip` (sin cambios respecto a Beta 12).
- **Mac Apple Silicon:** `Selector_Cine_365_Mac_Apple_Silicon_v9_BETA12.zip` (sin cambios).
- **Mac Intel:** `Selector_Cine_365_Mac_Intel_v9_BETA12.zip` (sin cambios).
- **Documentación:** manual completo y guía OneDrive PDF, historial v9 en PDF y Word.

## Instalación sin pérdida de datos

**Si utilizas Beta 11, 12 o 12.1 de Android:** guarda una copia JSON por precaución y **instala esta APK encima**. No desinstales antes. El identificador de aplicación y el certificado coinciden, y el código 9014 es mayor que 9011, 9012 y 9013. Tus datos locales permanecen en la aplicación.

**Si vienes de Android Beta 9:** aquella APK usa otro identificador y certificado; es necesaria una única exportación/importación JSON. Comprueba que la copia existe fuera de la app antes de desinstalar la versión antigua.

**Windows y Mac:** no necesitan descargar nuevamente sus ejecutables si ya tienen Beta 12.

## Verificaciones y límites

- Compilación Gradle y pruebas de sincronización y JavaScript superadas en la revisión [44806971301](https://github.com/cristianbarallobre/selector-cine-365/actions/runs/44806971301).
- APK firmada con esquemas Android v2 y v3; mismo certificado SHA-256 `b8a17f39e4fa24d7c110d452161425c7524412c3a5a4d7459b1132ffd0825fb6`.
- Huella SHA-256 de la APK exacta: `dcbcae51dbe47cfdf9ffc0273ceb08077bca95d3391a0b98989e97f60e955554`.
- La clave privada **no** está en GitHub; la release se verifica byte por byte antes de publicarse.
- **Pendiente:** comprobación visual en el Redmi real con navegación por gestos/botones y rotación. Es una beta experimental y OneDrive no sincroniza instantáneamente.

## Documentación y fuentes

- [Manual actualizado](https://github.com/cristianbarallobre/selector-cine-365/blob/main/docs/MANUAL_COMPLETO_SELECTOR_CINE_365.md)
- [Guía Android](https://github.com/cristianbarallobre/selector-cine-365/blob/main/android/LEEME_ANDROID_v9_BETA.md)
- [Historial de cambios](https://github.com/cristianbarallobre/selector-cine-365/blob/main/CHANGELOG.md)
- [Historial completo Word/PDF](https://github.com/cristianbarallobre/selector-cine-365/blob/main/docs/HISTORIAL_V9_BETA12.md)
