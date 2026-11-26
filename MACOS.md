# Selector Cine 365 para macOS

La versión estable es 9.2.16. Requiere macOS 12 o posterior. En ** → Acerca de este Mac** comprueba si tu equipo indica un chip Apple o un procesador Intel.

- [Descargar para Apple Silicon](https://github.com/cristianbarallobre/selector-cine-365/releases/download/v9.2.16-desktop/Selector_Cine_365_Mac_Apple_Silicon_9.2.16.zip)
- [Descargar para Intel](https://github.com/cristianbarallobre/selector-cine-365/releases/download/v9.2.16-desktop/Selector_Cine_365_Mac_Intel_9.2.16.zip)
- [Comprobar las sumas SHA-256](https://github.com/cristianbarallobre/selector-cine-365/releases/download/v9.2.16-desktop/SHA256SUMS.txt)

## Instalación y actualización

1. Exporta la biblioteca como JSON desde **Mis listas** y cierra la app anterior.
2. Descomprime el ZIP y sustituye **Selector Cine 365.app** en **Aplicaciones**.
3. Abre la app. Si Gatekeeper bloquea la primera apertura, verifica la descarga y consulta **Ajustes del Sistema → Privacidad y seguridad → Abrir igualmente**. La app tiene firma ad hoc y carece de notarización de Apple.

La biblioteca permanece en **~/Library/Application Support/Selector_Cine_365/**. Los archivos regenerables de TMDB usan **~/Library/Caches/Selector_Cine_365/**. La autorización de OneDrive usa el llavero de macOS.

## Sincronizar con OneDrive

Windows y macOS guardan tus listas y marcas automáticamente en el equipo. Con una cuenta Microsoft vinculada, pulsa **Sincronizar con la nube** debajo del título para enviar y recibir cambios. La app comprueba si hay cambios remotos pendientes, pero solo los aplica cuando pulsas el botón. La ventana de sincronización bloquea la edición hasta que termina. Conecta o revisa la cuenta desde **Mis listas → OneDrive y copias de seguridad**.

## Funciones y comprobaciones

Mac y Windows comparten Biblioteca, Descubrir, estadísticas, Favoritos y sorteo. El motor de ventana de Mac es WKWebView y los diálogos de archivos son propios del sistema.

GitHub Actions genera paquetes separados para Intel y Apple Silicon y verifica su firma ad hoc y el contenido del ZIP. La prueba del flujo de sincronización manual se hizo en Windows con datos controlados. El recorrido completo con una cuenta de OneDrive en macOS queda pendiente.

[Manual completo](docs/MANUAL_COMPLETO_SELECTOR_CINE_365.md) · [Guía de sincronización](docs/GUIA_SINCRONIZACION_SEGURA.md) · [Pruebas de GitHub](https://github.com/cristianbarallobre/selector-cine-365/actions/workflows/desktop-release.yml)
