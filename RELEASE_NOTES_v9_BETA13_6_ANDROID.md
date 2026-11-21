# Selector Cine 365 — Android v9 Beta 13.6

**`es.cine365.app` · `versionCode 9020` · `9.0.0-beta.13.6`**

## Correcciones

- **Exportar TXT** utiliza el nombre de la lista activa. Ejemplo: `Terror de los 80.txt`. Para la vista de comparación, `Comparación de listas.txt`. El sistema limpia caracteres inválidos en Windows y Android sin modificar el título dentro de Selector Cine 365.
- **Completar fichas TMDB pendientes**: en Ajustes aparece `Completar carátulas y plataformas pendientes`. La consulta se limita a la lista abierta y a películas sin datos, con error o sin cartel identificado; evita volver a descargar todas las fichas que ya funcionan.
- Los fallos transitorios TMDB no se conservan para siempre; dejan de bloquear una nueva consulta después de 30 minutos. La recuperación manual puede forzar el reintento.
- Restauración del selector de plataformas desde la caché local en Android.

## Qué sincroniza OneDrive

OneDrive comparte listas, estados vista/pendiente, revisiones por dispositivo e identificaciones TMDB confirmadas cuando se guardan en la biblioteca. **No incluye carátulas descargadas, datos detallados de TMDB, catálogos de plataformas ni la clave privada TMDB**: cada dispositivo mantiene su propia caché para proteger credenciales y evitar archivos grandes y consultas periódicas.

Si en Windows se ven todos los carteles pero en Android faltan algunos, abre en Android la lista afectada, entra en Ajustes y pulsa **Completar carátulas y plataformas pendientes** con conexión a Internet y clave TMDB configurada. Confirma después que las fichas aparecen. La caché de plataformas de cada dispositivo puede variar temporalmente por fecha de consulta.

## Seguridad y compatibilidad

- Se mantiene la conexión automática a la carpeta `Selector_Cine_365` de OneDrive de Beta 13.5.
- Android sigue guardando localmente y enviando/recibiendo solo mediante botón ↻ y confirmación; **sin sincronización periódica en segundo plano**.
- APK compilada con pruebas automatizadas, firmada localmente con la clave privada original; se comprobaron firmas v2/v3 y huella SHA-256 del certificado previo. La clave y contraseña no se suben a GitHub.
- Instalar sobre Beta 13.5 sin desinstalar; hacer copia JSON antes y comprobar en el Redmi.

© 2026.
