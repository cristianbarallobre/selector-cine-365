# Selector Cine 365 — v9 Beta 12
Basada en el código completo de v9 Beta 10, rama v9-android-beta, no en main/v8.2.

## Cambios
- Buscador: Lista actual / Todas las listas. Una sola ficha por identidad TMDB o título y año.
- Listas numeradas dinámicamente, con identificadores internos que no cambian.
- En fichas: año, valoración TMDB (no IMDb), indicador Listas: 1, 2, 3 con acceso directo y botón para añadir o quitar películas de otras listas.
- Solo el tick cambia Vista/Pendiente; el tick de búsqueda global respeta el historial causal.
- Estadísticas Windows: filas en contenedor desplazable y flechas fijas abajo. Android conserva la vista con scroll, sin paginación.
- Mantiene los cambios previos de v9: OneDrive/PKCE, historial y revisiones causales, sincronización y recuperación, «Sincronizar ahora», TMDB, instrucciones anchas, uso local, barra fija Android y arranque Windows sin CMD.

## Android: transición de Beta 9 a aplicación de ID estable
La Beta 9 usaba el paquete es.cine365.app.v9beta9 y una firma de depuración generada en GitHub Actions. No se dispone de su clave privada, por lo que no puede actualizarse directamente manteniendo datos.
Beta 12 usará es.cine365.app, sin sufijos de versión; las posteriores actualizaciones mantendrán ese mismo paquete y necesitarán la MISMA clave privada de firma.
Procedimiento seguro:
1. En Beta 9 exportar copia JSON. Guardarla en una carpeta fuera del almacenamiento privado de la app.
2. Confirmar que la copia JSON existe y contiene tus listas. Mantener Beta 9 hasta comprobarlo.
3. Desinstalar Beta 9 solo si deseas evitar dos aplicaciones simultáneas. No eliminar la copia JSON.
4. Instalar el APK Beta 12 FIRMADO, importar la copia y verificar las listas localmente antes de conectar nube.
5. Si usas OneDrive, reconectar, escoger la carpeta existente y revisar el primer intercambio de datos antes de confirmar.

La clave privada se guarda fuera del GitHub público y debe conservarse con su contraseña para firmar todas las versiones futuras del paquete es.cine365.app. No distribuir ni instalar APK sin firmar.

## Windows
Extraer el nuevo ZIP y ejecutar su EXE. Cerrar versiones anteriores y realizar una copia JSON antes. Mantiene su carpeta local AppData.

## Pruebas y precauciones
Pruebas de regresión JavaScript Windows/Android, pruebas Go del motor de sincronización y pruebas Android de revisiones.
Queda pendiente comprobar en dispositivos reales el comportamiento visual y la propagación por OneDrive, cuya latencia no es instantánea. Los conflictos concurrentes requieren revisión.

## Correcciones Beta 12 en Windows
- Mis listas permanece visible con 3–4 filas por página, incluso cuando se activa la sección OneDrive e historial. Las filas ya no se encogen a altura cero ni quedan recortadas.
- «Actualizar historial» consulta realmente /api/sync/history (historial local) e indica lectura en curso, hora de consulta y número de movimientos encontrados.
- Un historial vacío aparece sin contenedor destacado; explica que solo existen eventos desde las versiones que registran actividad.

## Compatibilidad Android
- Mismo identificador `es.cine365.app` y mismo certificado de firma que la APK Beta 11 generada en este proyecto. El código de versión 9012 actualiza 9011 sin cambiar la clave. Una instalación Beta 9 anterior con identificador distinto requiere la migración única del JSON descrita arriba.
