# Selector Cine 365 para macOS

Este manual corresponde a 9.2.26; consulta Releases/latest para confirmar la descarga disponible. Requiere macOS 12 o posterior. En ** → Acerca de este Mac** comprueba si tu equipo indica un chip Apple o un procesador Intel.

- [Descargar para Apple Silicon](https://github.com/cristianbarallobre/selector-cine-365/releases/download/v9.2.26-desktop/Selector_Cine_365_Mac_Apple_Silicon_9.2.26.zip)
- [Descargar para Intel](https://github.com/cristianbarallobre/selector-cine-365/releases/download/v9.2.26-desktop/Selector_Cine_365_Mac_Intel_9.2.26.zip)
- [Comprobar las sumas SHA-256](https://github.com/cristianbarallobre/selector-cine-365/releases/download/v9.2.26-desktop/SHA256SUMS.txt)

## Instalación y actualización

1. Exporta la biblioteca como JSON desde **Mis listas** y cierra la app anterior.
2. Descomprime el ZIP y sustituye **Selector Cine 365.app** en **Aplicaciones**.
3. Abre la app. Si Gatekeeper bloquea la primera apertura, verifica la descarga y consulta **Ajustes del Sistema → Privacidad y seguridad → Abrir igualmente**. La app tiene firma ad hoc y carece de notarización de Apple.

La biblioteca permanece en **~/Library/Application Support/Selector_Cine_365/**. Los archivos regenerables de TMDB usan **~/Library/Caches/Selector_Cine_365/**.

## Guardado y copia adicional

En **Mis listas → Guardado y copia adicional** puedes seleccionar una segunda ubicación. La biblioteca local sigue siendo la principal. Los servicios externos administran sus carpetas sin conexión directa desde la aplicación. Para trasladar los datos a otro equipo, exporta y restaura una copia JSON completa desde la aplicación. Incluye la biblioteca y el historial; no copies únicamente colecciones.json.

## Funciones y comprobaciones

Mac y Windows comparten Biblioteca, Descubrir, estadísticas, Favoritos y sorteo. El motor de ventana de Mac es WKWebView y los diálogos de archivos son propios del sistema.

GitHub Actions genera paquetes separados para Intel y Apple Silicon y verifica su firma ad hoc y el contenido del ZIP. Las pruebas de guardado y copia adicional utilizan carpetas temporales controladas. WebKit se comprueba a 960 × 680, 1200 × 820 y 1600 × 1000. El llavero protege la clave TMDB y la retira al desconectar.

[Manual completo](docs/MANUAL_COMPLETO_SELECTOR_CINE_365.md) · [Guía de recuperación](docs/GUIA_SINCRONIZACION_SEGURA.md) · [Pruebas de GitHub](https://github.com/cristianbarallobre/selector-cine-365/actions/workflows/desktop-release.yml)
