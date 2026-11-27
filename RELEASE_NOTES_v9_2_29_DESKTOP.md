# Selector Cine 365 9.2.29 — escritorio

Versión centrada en seguridad, recuperación, coherencia de los datos y estabilidad. Incluye todos los cambios de las candidatas 9.2.26–9.2.28.

- Servidor local: validación estricta de host y origen, bloqueo de marcos y conexiones externas desde la interfaz, credenciales en cabecera para las operaciones de escritura.
- Copia adicional: no reescribe los datos por cambiar de lista activa; muestra la fecha real del archivo; reintenta si la carpeta no está disponible. Conserva una copia dañada en cuarentena antes de reconstruirla desde el guardado principal.
- Historial: si falta la instantánea activa y el respaldo local, busca el historial exacto en copias completas o de recuperación y verifica su SHA-256. No sustituye una biblioteca válida ni descarta movimientos si no encuentra el historial exacto.
- TMDB: Biblioteca y Descubrir aprovechan la disponibilidad válida más reciente de una misma película por ID; las fichas antiguas abiertas pueden actualizar sus plataformas desde la caché local.
- Nuevas pruebas de origen HTTP, archivos dañados o vacíos, carpeta adicional inaccesible, cambios de navegación, recuperación exacta y disponibilidad compartida.
- Mantiene las cuatro señales de disponibilidad, la ayuda revisada, el guardado principal local, los filtros, el historial y la copia adicional opcional.
- Manuales de escritorio actualizados en Markdown, PDF y Word.

Antes de actualizar, exporta una copia completa, cierra la aplicación y sustituye el ejecutable o la app de tu arquitectura. No requiere firma comercial ni notarización de Apple.
