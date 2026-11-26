# Sincronización segura de Selector Cine 365

Windows y macOS 9.2.15. Android estable de la publicación v9.0.0.

## Preparar la primera conexión

1. Abre la app en el equipo que tiene la biblioteca más completa. Anota los recuentos de las listas importantes.
2. Exporta una copia JSON en **Mis listas** en escritorio o **Ajustes** en Android. Guarda una copia fuera de la carpeta sincronizada.
3. En escritorio, abre **Mis listas → OneDrive y copias de seguridad** y pulsa **Conectar cuenta Microsoft**. Autoriza en el navegador y vuelve a la app.
4. En el otro equipo, autoriza la misma cuenta. Revisa la sincronización y comprueba que aparecen las listas esperadas.

La carpeta habitual se llama **Selector_Cine_365**. Si antes usabas una carpeta distinta, revisa la configuración avanzada y las copias existentes antes de cambiarla.

## Revisar y confirmar

- Comprueba el nombre y los recuentos de **cada lista** en la vista previa.
- Si dos equipos tomaron decisiones opuestas sobre una película, elige cuál conservar antes de confirmar.
- Cancela si aparecen películas de muestra mezcladas con una lista real o si el número de películas resulta inesperado. Conserva las copias de ambos equipos.
## Sincronizar después de vincular

Windows y Mac guardan tus ediciones automáticamente en el equipo. Para intercambiar los cambios, pulsa **Sincronizar con la nube** bajo el título una vez conectada la cuenta Microsoft. La aplicación termina el guardado local, abre una ventana de progreso y bloquea la edición mientras envía y recibe. Si los dispositivos tomaron decisiones incompatibles, se abre una revisión antes de aplicar el cambio.

- **Subida pendiente** significa que hay cambios guardados en este equipo que todavía no se enviaron.
- **Recepción pendiente** significa que la comprobación encontró cambios de otro dispositivo que todavía no se aplicaron. **Recepción sin comprobar** indica que no se pudo verificar el estado remoto.
- **Sin cambios por subir o por recibir** muestra el resultado de la última comprobación. Antes de pasar a otro equipo, sincroniza en el primero y luego en el segundo.
- Con una carpeta local configurada, usa **Sincronizar ahora** dentro de **Mis listas → OneDrive y copias de seguridad**. La app escribe la copia en esa carpeta y el cliente de nube se encarga de enviarla.
- En Android estable, usa el control de sincronización y la revisión de su propia interfaz.

Si la operación falla, los cambios que ya guardaste siguen en el equipo. Revisa el aviso y prueba de nuevo cuando recuperes la conexión.

## Qué viaja entre equipos

Se comparten los datos de biblioteca compatibles con cada versión. El escritorio 9.2.15 guarda marcas, favoritos, valoraciones y orden de Favoritos. Las identificaciones corregidas viajan como IDs de TMDB. La APK estable de Android pertenece a una etapa anterior.

Las claves de TMDB, los carteles, la caché y las sesiones Microsoft se guardan en cada dispositivo. Windows protege la autorización con DPAPI y Mac con el llavero. También puedes usar la app sin conectar OneDrive.

## Volver a una copia anterior

1. Si ves una biblioteca equivocada, evita confirmar nuevos cambios en OneDrive.
2. Exporta la biblioteca actual para conservar el estado antes de recuperar.
3. Elige un JSON verificado o consulta **Mis listas → OneDrive y copias de seguridad → Actividad reciente, historial y copias de recuperación**. Si solo falta una edición, repítela en la biblioteca actual. Recuperar una copia sustituye toda la biblioteca local.
4. Abre las listas importantes y comprueba los recuentos. Revisa la vista previa antes de volver a sincronizar.

[Manual completo](https://github.com/cristianbarallobre/selector-cine-365/blob/main/docs/MANUAL_COMPLETO_SELECTOR_CINE_365.md) · [Rutas de datos](https://github.com/cristianbarallobre/selector-cine-365/blob/main/docs/UBICACIONES_DATOS_SELECTOR_CINE_365.md)

Autor: Cristian Barallobre Rúa.
