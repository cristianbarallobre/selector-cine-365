# Selector Cine 365 9.2.26 — escritorio

- Historial de actividad separado del archivo principal, con referencia de integridad, respaldo local y recuperación protegida.
- Exportación y restauración conjunta de biblioteca e historial; copia adicional actual, anterior y hasta 14 diarias. Las copias anteriores sin historial siguen siendo compatibles.
- Retirada del motor y de los controles de conexión a cuentas externas: solo guardado local y carpeta adicional opcional.
- Disponibilidad TMDB con actualización en segundo plano, caducidad de 24 horas, conservación del último dato válido y control de saturación de consultas.
- Credenciales TMDB protegidas mediante DPAPI en Windows y llavero en macOS. Separación de caché regenerable en Windows.
- Mantenimiento: uso de disco, memoria del motor, historial y copias; limpieza opcional de archivos de caché antiguos con duplicado idéntico.
- Pruebas automatizadas de guardado, restauración, caché, interfaz en diferentes tamaños y compilación para Windows, macOS Intel y macOS Apple Silicon.
- Uso personal: sin firma comercial ni notarización de Apple; los paquetes Mac conservan la firma ad hoc técnica.

**Antes de actualizar:** exporta una copia completa en la aplicación, cierra la anterior y sustituye el ejecutable. En macOS selecciona tu arquitectura. Comprueba SHA256SUMS.txt. No se importan bibliotecas de otras instalaciones automáticamente.
