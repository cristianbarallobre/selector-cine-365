# Selector Cine 365 — v9 Beta 12.1 (corrección de Android)

**Actualización correctiva sobre v9 Beta 12.** Windows, macOS Intel y macOS Apple Silicon conservan **exactamente los mismos paquetes y funciones de Beta 12**. La única modificación del código de aplicación es Android.

## Qué se corrige

- **Android 15 y posteriores:** la app ya no debe situar el título bajo la cámara frontal/notch, la barra de estado ni la barra inferior de gestos. Un contenedor nativo respeta los insets de `systemBars` y `displayCutout` del sistema Android en los cuatro lados.
- **Android 11–14:** mismo tratamiento adaptable, también en orientación horizontal.
- **Android 9–10:** mantiene el ajuste tradicional y pide explícitamente no ocupar el recorte.
- **Sin márgenes fijos**: no depende del modelo, la marca ni del tamaño del notch.
- Se conserva el diseño con **cabecera fija**, desplazamiento interno, búsqueda global, TMDB y sincronización causal; sin cambios de backend.

## Descargas

- **Android Beta 12.1, firmada:** `Selector_Cine_365_Android_v9_BETA12_1.apk` — **versionCode 9013**; paquete `es.cine365.app`.
- Windows Beta 12: `Selector_Cine_365_Windows_v9_BETA12.zip` — mismo archivo sin recompilar.
- Mac Apple Silicon Beta 12: `Selector_Cine_365_Mac_Apple_Silicon_v9_BETA12.zip` — mismo archivo sin recompilar.
- Mac Intel Beta 12: `Selector_Cine_365_Mac_Intel_v9_BETA12.zip` — mismo archivo sin recompilar.

**Android Beta 11 o Beta 12 → Beta 12.1:** instalar encima, sin desinstalar. Mismo certificado SHA-256 `b8a17f39e4fa24d7c110d452161425c7524412c3a5a4d7459b1132ffd0825fb6`, firma v2/v3 válida; identificador de aplicación sin cambios, versión 9013.

**Android Beta 9 → Beta 12.1:** sigue requiriendo la migración inicial mediante copia JSON porque su paquete y firma eran diferentes. Nunca desinstales sin exportar y verificar tus datos primero.

## Verificación y límites

- Pruebas de regresión JavaScript, pruebas Java y compilación Gradle en GitHub Actions.
- APK firmada fuera del repositorio y verificada con `apksigner`. La clave privada no está en GitHub.
- Queda pendiente comprobar el notch y navegación por gestos en un dispositivo Redmi y otras pantallas con recortes reales. La beta sigue siendo experimental.
- La sincronización OneDrive no es instantánea; conserva copia JSON independiente.

**Documentación:** [manual](docs/MANUAL_COMPLETO_SELECTOR_CINE_365.md), [Android](android/LEEME_ANDROID_v9_BETA.md), [historial](CHANGELOG.md).
