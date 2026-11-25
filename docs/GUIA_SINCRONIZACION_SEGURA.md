# Sincronización segura de Selector Cine 365

Windows y macOS 9.2.7. Android estable de la publicación v9.0.0.

## Preparar la primera conexión

1. Abre la app en el equipo que tiene la biblioteca más completa. Anota los recuentos de las listas importantes.
2. Exporta una copia JSON en **Mis listas** en escritorio o **Ajustes** en Android. Guarda una copia fuera de la carpeta sincronizada.
3. En la sección de OneDrive, pulsa **Conectar cuenta Microsoft**. Autoriza en el navegador y vuelve a la app.
4. En el otro equipo, autoriza la misma cuenta. Revisa la sincronización y comprueba que aparecen las listas esperadas.

La carpeta habitual se llama **Selector_Cine_365**. Si antes usabas una carpeta distinta, revisa la configuración avanzada y las copias existentes antes de cambiarla.

## Revisar y confirmar

- Comprueba el nombre y los recuentos de **cada lista** en la vista previa.
- Si dos equipos tomaron decisiones opuestas sobre una película, elige cuál conservar antes de confirmar.
- Cancela si aparecen películas de muestra mezcladas con una lista real o si el número de películas resulta inesperado. Conserva las copias de ambos equipos.
- Tras vincular la biblioteca, Windows y Mac pueden recibir cambios automáticamente. Android estable guarda localmente y utiliza la revisión manual para intercambiar cambios.

## Qué viaja entre equipos

Se comparten los datos de biblioteca compatibles con cada versión. El escritorio 9.2.7 guarda marcas, favoritos, valoraciones y orden de Favoritos. Las identificaciones corregidas viajan como IDs de TMDB. La APK estable de Android pertenece a una etapa anterior.

Las claves de TMDB, los carteles, la caché y las sesiones Microsoft se guardan en cada dispositivo. Windows protege la autorización con DPAPI y Mac con el llavero. También puedes usar la app sin conectar OneDrive.

## Volver a una copia anterior

1. Si ves una biblioteca equivocada, evita confirmar nuevos cambios en OneDrive.
2. Exporta la biblioteca actual para conservar el estado antes de recuperar.
3. Elige un JSON verificado o usa las opciones de recuperación de la app. Restaurar un JSON sustituye la biblioteca local completa.
4. Abre las listas importantes y comprueba los recuentos. Revisa la vista previa antes de volver a sincronizar.

[Manual completo](https://github.com/cristianbarallobre/selector-cine-365/blob/main/docs/MANUAL_COMPLETO_SELECTOR_CINE_365.md) · [Rutas de datos](https://github.com/cristianbarallobre/selector-cine-365/blob/main/docs/UBICACIONES_DATOS_SELECTOR_CINE_365.md)

Autor: Cristian Barallobre Rúa.
