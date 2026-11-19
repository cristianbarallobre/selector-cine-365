# Selector Cine 365 — Android v9.0 beta 2

**Autor:** Cristian Barallobre Rúa. Repositorio oficial: https://github.com/cristianbarallobre/selector-cine-365

## ¿Qué incluye esta prueba?

- Identidad visual renovada, icono de claqueta generado para el proyecto y cabecera fija.
- Cinco pantallas: Biblioteca, Sorteo, Estadísticas, Listas y Nube.
- Scroll solo dentro del contenido de la sección seleccionada y controles táctiles grandes.
- Compatible con el formato original de colecciones (ID interno `atlas`), conservando el nombre y carpeta de almacenamiento históricos cuando corresponda.
- Primera conexión a carpeta compartida en modo **solo lectura**. Vista previa con número de listas y películas vistas; hay que confirmar para importar y empezar a publicar el snapshot de Android.
- Las copias de PC `Selector_Cine_365_Copias/actual/colecciones.json` y `historico/` nunca son modificadas por Android. Android utiliza su propio fichero en `sincronizacion/dispositivos/`.
- Revisión de cambios hechos sin conexión por dispositivo mediante relojes causales. Los estados contradictorios concurrentes se muestran para resolverlos.
- La importación crea una copia local de recuperación antes de guardar.

## Instalar en Xiaomi Mi 9T

Descarga el APK de la prepublicación v9.0-beta.2 y ábrelo desde el gestor de archivos. **Se instala junto a tu antigua APK** con el identificador `es.cine365.app.v9beta2`, sin eliminar los datos de la beta anterior. Se instala aparte porque la firma de desarrollo puede cambiar entre compilaciones. Para trasladar tus listas, expórtalas desde la beta 1 e impórtalas en la beta 2. Android puede advertir sobre la firma de desarrollo; no desactives Play Protect de forma general.

## Mejoras de la beta 2

- Las secciones vuelven arriba cuando cambias de pestaña.
- Los cinco iconos inferiores son monocromos, minimalistas y del mismo tamaño.
- «Cargar más» incorpora 24 películas sin reconstruir la lista ni desplazar el scroll.
- Los carteles se cargan desde TMDB mediante la ruta nativa de Android, con caché local.
- Las tarjetas crecen según el título y las plataformas, con controles en una fila independiente.

## Sincronización con OneDrive

1. En Windows utiliza **Mis listas > Copia en nube** para dejar una copia de seguridad actualizada en la carpeta habitual. No elimines nada.
2. Abre Android > Nube > Elegir carpeta. Selecciona la carpeta que contiene `Selector_Cine_365_Copias`, o esa misma carpeta.
3. Pulsa Activar copia (solo registra el permiso) y **Revisar sincronización**.
4. Lee los números de listas y vistos del móvil y los previstos tras importar. Resuelve conflictos si aparecen y confirma.
5. Solo entonces el móvil publica su archivo, **no** la copia `actual` del PC.

La sincronización a través de la app OneDrive para Android **depende de que el proveedor permita elegir y escribir carpetas mediante el Storage Access Framework**. OneDrive puede no aparecer en el selector de carpetas o no dar permisos persistentes. No podemos saltarnos esa restricción. Como alternativa temporal, exporta/importa JSON manualmente. **El inicio de sesión directo con Microsoft Graph todavía no está implementado** y requiere registrar una aplicación de Microsoft con sus permisos. No incluyas contraseñas en código ni archivos compartidos.

Windows v8.2 / v8.1.1 sigue haciendo respaldos, **pero todavía no recibe los cambios de Android**. No confundas esta beta con sincronización bidireccional completa de escritorio. El protocolo causal del móvil es una base para completar las versiones Windows y Mac.

La beta requiere pruebas físicas y debe mantenerse una copia independiente de las colecciones existentes.
