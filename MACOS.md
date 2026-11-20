# Selector Cine 365 — macOS v9 Beta 12

**Autor: Cristian Barallobre Rúa.** Aplicación para macOS 12 o posterior, en dos variantes: **Apple Silicon (arm64)** e **Intel (x86_64)**.

## Paridad con Windows

La edición Mac utiliza el **mismo código Go**, la **misma interfaz web de escritorio** (\`plantilla.html\`), \`collections.go\`, TMDB, revisiones causales, copias y estadísticas que Windows Beta 12. Únicamente cambia la ventana nativa (WKWebView en Mac frente a WebView2 en Windows), los diálogos del sistema y el empaquetado \`.app\`.

Incluye la búsqueda en todas las listas, listas numeradas, indicadores «Listas:», edición de pertenencia, valoración TMDB junto al año, historial y sincronización segura, además de las correcciones de Beta 12 para listas invisibles e historial vacío.

**Límite de verificación:** compilar en macOS y validar el contenido no equivale a haberlo probado manualmente en un Mac físico con OneDrive real. Las diferencias propias de WKWebView o de proveedores de carpetas deben evaluarse antes de llamar estable a la versión.

## Instalación y actualización

1. Descarga en [Releases](../../releases) el paquete de tu procesador:
   - \`Selector_Cine_365_Mac_Apple_Silicon_v9_BETA12.zip\`.
   - \`Selector_Cine_365_Mac_Intel_v9_BETA12.zip\`.
2. Guarda una copia JSON de tus listas antes de actualizar.
3. Cierra la app anterior, descomprime el ZIP y reemplaza \`Selector Cine 365.app\` en **Aplicaciones**.
4. Abre la aplicación. Si macOS bloquea el inicio de una app no notarizada, comprueba el origen y consulta **Ajustes del Sistema → Privacidad y seguridad → Abrir igualmente**.

Esta beta está firmada **ad hoc**; **NO utiliza un certificado Apple Developer ID y NO está notarizada**. Gatekeeper podría mostrar advertencias. No recomendamos desactivar Gatekeeper globalmente.

## Almacenamiento

Los datos se conservan en \`~/Library/Application Support/AtlasDelCine/\`. El nombre de la ruta es un identificador histórico para no perder listas existentes. La app puede funcionar localmente sin configurar nube.

Para usar nube en Mac, elige una **carpeta local que el proveedor de sincronización mantenga sincronizada y accesible**. Puedes seleccionar la carpeta del cliente OneDrive o un proveedor equivalente que ofrezca lectura y escritura reales. En el móvil Android, OneDrive se conecta a través de Microsoft Graph y debes seleccionar la misma carpeta dentro de tu OneDrive personal.

La sincronización **no es instantánea**. Para la primera conexión entre dispositivos pulsa **Revisar sincronización**, examina la información detectada y confirma. Los cambios posteriores sin conflicto se reciben periódicamente; las acciones concurrentes pueden requerir revisión.

## Compilar para macOS

Se necesitan Xcode Command Line Tools, Go 1.23+, \`sips\`, \`iconutil\`, \`codesign\` y \`plutil\`.

\`\`\`bash
CGO_ENABLED=1 go test ./...
MAC_ARCH=arm64 MAC_LABEL=Apple_Silicon bash scripts/build_macos_app.sh
# En otro trabajo o runner para Intel:
MAC_ARCH=amd64 MAC_LABEL=Intel bash scripts/build_macos_app.sh
\`\`\`

El script genera los dos ZIP nombrados con **v9 BETA12**, \`CFBundleShortVersionString = 9.0.0-beta.12\` y \`CFBundleVersion = 9012\`. Genera el icono \`.icns\` desde el \`icon.png\` actual, firma ad hoc y verifica la estructura de la app.

El binario para Apple Silicon se puede compilar desde un runner Intel cuando esté disponible el SDK correspondiente y un compilador C compatible; el flujo CI comprueba su arquitectura con \`file\`. Siempre que sea posible, ejecuta una prueba manual también en un Mac Apple Silicon.

## Ayuda y referencias

- [Manual Beta 12](docs/MANUAL_COMPLETO_SELECTOR_CINE_365.md).
- [Historial de cambios](CHANGELOG.md).
- [Guía OneDrive](docs/GUIA_AZURE_ENTRA_ONEDRIVE_SELECTOR_CINE_365.md).

No compartas ficheros locales \`colecciones.json\`, tokens ni claves de TMDB al informar de problemas públicamente.
