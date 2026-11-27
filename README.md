# Selector Cine 365

Aplicación de uso personal para organizar colecciones de cine. Puedes marcar películas como vistas, buscarlas en tus listas y sortear qué ver. TMDB añade fichas, carteles y disponibilidad en España. OneDrive es opcional para compartir la biblioteca entre equipos.

## Descargas

| Equipo | Versión estable | Requisitos |
| --- | --- | --- |
| Windows | [9.2.21 para Windows](https://github.com/cristianbarallobre/selector-cine-365/releases/download/v9.2.21-desktop/Selector_Cine_365_Windows_9.2.21.exe) | Windows 10 u 11 de 64 bits y Microsoft Edge WebView2 Runtime |
| Mac con chip de Apple | [9.2.21 para Apple Silicon](https://github.com/cristianbarallobre/selector-cine-365/releases/download/v9.2.21-desktop/Selector_Cine_365_Mac_Apple_Silicon_9.2.21.zip) | macOS 12 o posterior |
| Mac con procesador Intel | [9.2.21 para Intel](https://github.com/cristianbarallobre/selector-cine-365/releases/download/v9.2.21-desktop/Selector_Cine_365_Mac_Intel_9.2.21.zip) | macOS 12 o posterior |
| Android | [APK estable de la publicación v9.0.0](https://github.com/cristianbarallobre/selector-cine-365/releases/tag/v9.0.0) | Android 8 o posterior. Compilación Beta 13.13 Limpia |

Windows y Mac comparten la interfaz y las funciones de escritorio. Android tiene interfaz táctil propia. La [compilación de prueba de Android 9.2.5](https://github.com/cristianbarallobre/selector-cine-365/releases/tag/v9.2.5-android) se instala en paralelo con la versión estable.

## Instalar o actualizar

- **Windows.** Cierra la app anterior y abre el nuevo EXE portátil. Si falta WebView2, instala el runtime de Microsoft. SmartScreen puede mostrar un aviso al abrir la descarga.
- **Mac.** Descomprime el ZIP de tu procesador y mueve la aplicación a **Aplicaciones**. Los paquetes tienen firma ad hoc y Apple no los ha notarizado. Verifica su origen antes de autorizar la apertura en **Privacidad y seguridad**.
- **Android.** La APK pública usa el identificador es.cine365.app.limpia. Actualiza la instalación Limpia cuando coinciden el identificador y el certificado. Exporta un JSON antes de sustituir cualquier edición.

[Última versión de escritorio](https://github.com/cristianbarallobre/selector-cine-365/releases/latest) · [Sumas SHA-256](https://github.com/cristianbarallobre/selector-cine-365/releases/download/v9.2.21-desktop/SHA256SUMS.txt) · [Todas las publicaciones](https://github.com/cristianbarallobre/selector-cine-365/releases)

## Sincronizar con OneDrive

Windows y macOS guardan tus listas y marcas automáticamente en el equipo. Con una cuenta Microsoft vinculada, pulsa **Sincronizar con la nube** debajo del título para enviar y recibir cambios. La app comprueba si hay cambios remotos pendientes, pero solo los aplica cuando pulsas el botón. La ventana de sincronización bloquea la edición hasta que termina. Antes de cambiar de equipo, sincroniza en el que acabas de usar y luego en el siguiente.

## Manuales

- [Manual completo](docs/MANUAL_COMPLETO_SELECTOR_CINE_365.md) · [PDF](docs/MANUAL_COMPLETO_SELECTOR_CINE_365.pdf) · [Word](docs/MANUAL_COMPLETO_SELECTOR_CINE_365.docx)
- [Guía de Mac](MACOS.md) · [Guía de Android](android/LEEME_ANDROID_v9_BETA.md)
- [Sincronización y recuperación](docs/GUIA_SINCRONIZACION_SEGURA.md) · [Configuración de Entra](docs/GUIA_AZURE_ENTRA_ONEDRIVE_SELECTOR_CINE_365.md)
- [Rutas de datos](docs/UBICACIONES_DATOS_SELECTOR_CINE_365.md) · [Índice de documentación](docs/README.md) · [Historial](CHANGELOG.md)

## Código y versiones anteriores

Cambiar de lista activa no marca una subida a OneDrive por sí solo.

El tag de escritorio identifica el código que se compiló para Windows y macOS 9.2.21. La rama principal centraliza la documentación. La APK pública de Android tiene su propia publicación. [Fuente 9.2.21](https://github.com/cristianbarallobre/selector-cine-365/tree/v9.2.21-desktop) · [Publicación Android](https://github.com/cristianbarallobre/selector-cine-365/releases/tag/v9.0.0)

Autor: Cristian Barallobre Rúa.
