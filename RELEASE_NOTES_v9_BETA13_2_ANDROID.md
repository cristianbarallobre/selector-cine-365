# Selector Cine 365 — Android v9 Beta 13.2 (batería y tareas en segundo plano)

**Actualización específica de Android.** Windows y macOS conservan su última versión **Beta 13.1**. No se cambian el diseño, las listas, la biblioteca compartida ni la clave de TMDB.

## Qué se ha optimizado

- **Menos trabajo al arrancar:** se ha corregido una doble construcción del motor `Backend` que inicializaba archivos y servicios dos veces sin necesidad.
- **WebView pausada al minimizar:** se suspende su ejecución y los temporizadores de JavaScript cuando Android coloca la aplicación en segundo plano, y se reanudan al volver.
- **Ediciones protegidas al salir rápidamente:** si hay una película marcada o una lista editada pendiente de los 240 ms de guardado, se escribe inmediatamente en almacenamiento local sin perder su revisión causal de sincronización; el guardado normal en nube no se elimina.
- **Sin ping innecesario:** se elimina la petición interna `/api/ping` repetida cada 12 segundos. Android no tiene el servidor HTTP de Windows: usa su propio puente Java.
- **Sincronización optimizada:** comprobación automática de OneDrive cada **60 segundos en primer plano** (en lugar de 30), más la comprobación al volver a la app. El backend descarta comprobaciones automáticas que arranquen estando en segundo plano, sin impedir las sincronizaciones manuales.
- **Sin pérdida de funciones:** se mantiene la sincronización causal Vista/Pendiente, eliminaciones de listas, identificaciones TMDB compartidas, las copias y el funcionamiento exclusivamente local.

## Actualización sin borrar datos

- App Android: `es.cine365.app`.
- Versión: `9.0.0-beta.13.2` — `versionCode 9016` (Beta 13 = 9015).
- Se distribuye con el **mismo certificado de firma** de Beta 11, 12, 12.1, 12.2 y 13. **Instalar encima sin desinstalar.**
- Las credenciales personales de TMDB y Microsoft permanecen guardadas localmente; nunca se distribuyen con la APK.
- Exporta un JSON de la biblioteca antes de actualizar, por precaución.

## Pruebas y límites

Se prueban la compilación Java/Gradle, la lógica de sincronización JSON, la protección en segundo plano, la frecuencia de consulta, las comprobaciones de la UI y la compatibilidad de versión/firma. **La mejora real de batería debe comprobarse en el móvil con varias horas de uso**, porque Android calcula los porcentajes relativos y el consumo depende del dispositivo, pantalla, red y sincronización. Una cuota del 50% del gasto de las apps durante solo ocho minutos no es el 50% de la batería total.

Estas mejoras afectan al uso **en segundo plano** y a las comprobaciones repetidas; no desactivan el guardado ni garantizan que el consumo sea cero mientras se esté usando la app.

© 2026 Cristian Barallobre Rúa.
