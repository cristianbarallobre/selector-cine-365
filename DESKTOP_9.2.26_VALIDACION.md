# Candidata 9.2.26 de escritorio
- Caducidad de disponibilidad TMDB: 24 h; se conservan resultados útiles si la red falla. Hasta 24 fichas antiguas por ciclo de fondo en Biblioteca.
- Copia completa con historial (máx. 10.000 eventos), sin claves ni caché. La importación sustituye biblioteca e historial; crea antes un punto de rescate.
- Copia adicional: actual, anterior, hasta 14 copias diarias; oportunidad de terminar al cerrar limitada a 2 segundos.
- Caché TMDB en LocalAppData en Windows y Library/Caches en macOS. El origen antiguo se conserva y se aprovechan cachés compatibles; no mover ni borrar la biblioteca.
- Credencial TMDB mediante DPAPI o Llavero de macOS. Si falla la migración no se borra la configuración anterior.
- Escrituras protegidas de fichas; copia transaccional de metadatos sin ciclo JSON; diagnóstico de memoria del motor y benchmark sintético.
- Pendiente de validar mediante CI y pruebas nativas. La release estable no cambia automáticamente.