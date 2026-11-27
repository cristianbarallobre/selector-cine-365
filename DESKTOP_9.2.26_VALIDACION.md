# Validación de la candidata 9.2.26

La descarga estable sigue identificada en Releases. Este documento describe el código en validación.

- Guardado: la biblioteca confirma una referencia SHA-256 a una instantánea del historial escrita previamente. Las copias locales son autocontenidas; se conserva un espejo de recuperación compatible con versiones antiguas.
- Copias adicionales: biblioteca e historial, actual y anterior, hasta 14 diarias. Sin claves, caché ni sesiones. Las copias automáticas locales tienen límites de espacio; las pendientes no se retiran por antigüedad.
- Se ha eliminado la implementación del motor de sincronización remota y la autenticación asociada. Se mantienen los campos y utilidades de revisión local necesarios para leer formatos anteriores. No se conectan cuentas.
- TMDB: cuatro consultas simultáneas como máximo, pausa ante 429, reutilización de datos guardados y aviso de fecha. DPAPI en Windows y llavero en macOS, con pruebas de migración, desconexión y reconexión.
- Mantenimiento: tamaños de datos, caché, copias, historial y registros; limpieza voluntaria de duplicados idénticos. No se elimina una carpeta personal.
- Pruebas: Windows y macOS Intel/Apple Silicon, guardado y recuperación, benchmarks de actividad, disco y caché. WebKit a 960×680, 1200×820 y 1600×1000. No sustituye una revisión visual manual en equipos físicos.
- No se cambia automáticamente la release publicada. Para un traslado o una vuelta a otra versión, exportar una copia completa desde la aplicación.
