# Selector Cine 365

Aplicación de uso personal para organizar colecciones de cine. Puedes marcar películas como vistas, buscarlas en tus listas y sortear qué ver. TMDB añade fichas, carteles y disponibilidad en España. En escritorio el guardado principal es local y puedes elegir una carpeta para conservar una copia adicional.

## Descargas

| Equipo | Versión estable | Requisitos |
| --- | --- | --- |
| Windows | [9.2.25 para Windows](https://github.com/cristianbarallobre/selector-cine-365/releases/download/v9.2.25-desktop/Selector_Cine_365_Windows_9.2.25.exe) | Windows 10 u 11 de 64 bits y Microsoft Edge WebView2 Runtime |
| Mac con chip de Apple | [9.2.25 para Apple Silicon](https://github.com/cristianbarallobre/selector-cine-365/releases/download/v9.2.25-desktop/Selector_Cine_365_Mac_Apple_Silicon_9.2.25.zip) | macOS 12 o posterior |
| Mac con procesador Intel | [9.2.25 para Intel](https://github.com/cristianbarallobre/selector-cine-365/releases/download/v9.2.25-desktop/Selector_Cine_365_Mac_Intel_9.2.25.zip) | macOS 12 o posterior |

Windows y Mac comparten la interfaz y las funciones de escritorio.

## Instalar o actualizar

- **Windows.** Cierra la app anterior y abre el nuevo EXE portátil. Si falta WebView2, instala el runtime de Microsoft. SmartScreen puede mostrar un aviso al abrir la descarga.
- **Mac.** Descomprime el ZIP de tu procesador y mueve la aplicación a **Aplicaciones**. Los paquetes tienen firma ad hoc y Apple no los ha notarizado. Verifica su origen antes de autorizar la apertura en **Privacidad y seguridad**.

[Última versión de escritorio](https://github.com/cristianbarallobre/selector-cine-365/releases/latest) · [Sumas SHA-256](https://github.com/cristianbarallobre/selector-cine-365/releases/download/v9.2.25-desktop/SHA256SUMS.txt) · [Todas las publicaciones](https://github.com/cristianbarallobre/selector-cine-365/releases)

## Guardado y copia adicional

En Windows y macOS los datos se guardan automáticamente en el equipo. En **Mis listas → Guardado y copia adicional** puedes elegir otra carpeta. La copia contiene la biblioteca personal, no las claves ni la caché. Si la carpeta pertenece a un servicio de almacenamiento, ese servicio gestiona sus propias copias. La aplicación no conecta cuentas ni combina bibliotecas. Para trasladar los datos a otro equipo, exporta e importa el JSON manualmente.

El historial de actividad conserva hasta 10.000 movimientos confirmados de la biblioteca, con búsqueda, filtros y paginación.

## Manuales

- [Manual completo](docs/MANUAL_COMPLETO_SELECTOR_CINE_365.md) · [PDF](docs/MANUAL_COMPLETO_SELECTOR_CINE_365.pdf) · [Word](docs/MANUAL_COMPLETO_SELECTOR_CINE_365.docx)
- [Guía de Mac](MACOS.md)
- [Guardado y recuperación](docs/GUIA_SINCRONIZACION_SEGURA.md) · [PDF](docs/GUIA_SINCRONIZACION_SEGURA.pdf) · [Word](docs/GUIA_SINCRONIZACION_SEGURA.docx)
- [Rutas de datos](docs/UBICACIONES_DATOS_SELECTOR_CINE_365.md) · [Índice de documentación](docs/README.md) · [Historial](CHANGELOG.md)

## Código y versiones anteriores

El tag de escritorio identifica el código que se compiló para Windows y macOS 9.2.25. La rama principal integra el código de escritorio y la documentación; las descargas estables se identifican mediante su tag. El desarrollo de la siguiente versión se valida antes de publicarse. [Fuente 9.2.25](https://github.com/cristianbarallobre/selector-cine-365/tree/v9.2.25-desktop)

Autor: Cristian Barallobre Rúa.
