# Manual de Selector Cine 365

Windows y macOS 9.2.30.

## Instalar y proteger la biblioteca

1. Elige la descarga de tu sistema desde el README de GitHub. En Mac comprueba si tu procesador es Intel o Apple Silicon.
2. Si ya usas la app, exporta una copia JSON y cierra la versión anterior.
3. En Windows abre el nuevo EXE portátil. Si el equipo lo necesita, instala WebView2 Runtime. En Mac descomprime el ZIP, lleva la app a Aplicaciones y reemplaza la anterior.
4. Si Windows o Mac avisan sobre el editor, verifica primero que el archivo procede de la publicación oficial y comprueba su SHA-256. Los paquetes Mac llevan firma ad hoc y carecen de notarización.

[Descargas de escritorio](https://github.com/cristianbarallobre/selector-cine-365/releases/tag/v9.2.30-desktop)

## Biblioteca y listas

- Abre el selector superior para elegir una lista. Puedes cambiar de lista desde Estadísticas y mantener abierta la sección.
- Usa Buscar y elige **Lista actual** o **Todas las listas**. Los filtros de escritorio permiten combinar el estado vista o pendiente con plataformas, década, género y duración, y elegir el orden.
- Usa la paginación o el campo **Página** para saltar.
- Marca el estado desde la casilla. El indicador **Listas** muestra las colecciones que contienen la película. Pulsa un número para ir a esa lista o el lápiz para editar la pertenencia.
- Si quitas todas las listas numeradas y la película permanece en Favoritos, puedes guardar directamente. Si va a desaparecer de todas las listas, la app pide confirmación.
- Desde **Mis listas** puedes importar o añadir películas, gestionar colecciones y configurar la comparación de listas. En escritorio puedes reordenar arrastrando el tirador cuando estás viendo el orden original de la lista.

## Favoritos y valoraciones en escritorio

Pulsa el corazón para añadir una película. Aparece la lista **♥ Favoritos** en el selector y el corazón pequeño junto a los números te lleva a ella. La lista se oculta al quitar el último favorito. Si el mismo ID de TMDB identifica varios títulos de una película, el favorito y la valoración se comparten.

Al añadir la película, se marca como vista donde ya figure. Si después la importas desde un TXT, puede heredar esa marca desde Favoritos. En Favoritos se muestra como vista y la casilla queda bloqueada. Puedes ordenar Favoritos con el tirador y el orden se guarda. Las películas que solo están en Favoritos también pueden cargar su ficha. La valoración personal, de uno a cinco, es independiente de la nota pública de TMDB.

## Sorteo y estadísticas

El sorteo normal selecciona películas pendientes de la lista activa. En **Sortear entre** puedes marcar varias plataformas y entran las películas disponibles en cualquiera. **Cualquier plataforma** limita el sorteo a títulos con streaming registrado. El modo 🧠 **Sorteo inteligente** permite escoger vistas, pendientes o ambas y filtrar por género y duración. Las estadísticas reflejan la lista activa y se completan con los datos de TMDB disponibles.

## TMDB y Descubrir

Configura una clave de TMDB en la app si quieres carteles, disponibilidad en España y fichas. **Descubrir** permite buscar películas y personas, explorar el catálogo y abrir las fichas. Al pasar el puntero por un actor o director puedes ver su retrato. Pulsa su nombre para abrir la biografía y las películas en las que actuó o que dirigió. Los carteles grandes y retratos tienen una caché local de hasta 32 MiB o 640 imágenes. La disponibilidad puede cambiar. Pulsa **Actualizar** en Biblioteca para volver a consultar las fichas y plataformas de la lista activa; la disponibilidad guardada caduca a las 24 horas y se actualiza en segundo plano cuando procede. Si falla la red, se conserva la última información útil; pasa el ratón sobre el icono para ver la fecha de comprobación. En Biblioteca el círculo verde indica una consulta correcta de menos de 24 horas, el reloj amarillo datos pendientes de actualizar, el círculo rojo con exclamación una consulta fallida y el interrogante gris ausencia de una comprobación válida. Conviene comprobar la disponibilidad en el servicio.

Al añadir una película desde Descubrir, Biblioteca reutiliza la ficha y el cartel de TMDB guardados localmente. Al borrar una búsqueda de Descubrir, se recuperan la categoría, la página y el filtro de género anteriores.

Si la app ha elegido otra película, abre **Identificar película en TMDB** y pega su ID numérico o su enlace. Si los años difieren, verifica el título y pulsa **Guardar de todos modos** si la ficha es correcta. El año de tu lista se conserva. **Volver a búsqueda automática** retira la identificación manual.

## Guardado y copia adicional

La biblioteca se guarda automáticamente en el equipo. Desde **Mis listas → Copia de seguridad** exporta un JSON completo para conservar todas las listas, favoritos, valoraciones, identificaciones e historial. Una exportación TXT corresponde a la lista seleccionada.

1. Abre **Mis listas → Guardado y copia adicional** y pulsa **Elegir carpeta**.
2. Pulsa **Guardar ubicación**. Cada instalación conserva su propio archivo `colecciones.json` y una versión anterior dentro de `Selector_Cine_365_Copias`. La interfaz muestra la ruta completa. Se guardan hasta 14 copias diarias adicionales, limitadas también por espacio.
3. La copia se actualiza tras el guardado principal y puedes pulsar **Copiar ahora**. Si la ubicación no está disponible, los datos principales se conservan y la app reintenta la copia. Desactivar no borra los archivos.
4. Si necesitas recuperar una copia, utiliza **Restaurar copia** y elige su JSON. Restaurar sustituye la biblioteca completa; exporta antes los datos actuales.

No se conectan cuentas ni se combinan bibliotecas automáticamente. Una carpeta gestionada por un servicio externo queda bajo responsabilidad de ese servicio. Las claves TMDB y la caché permanecen en cada dispositivo. La clave se protege con DPAPI o el llavero de macOS. En Mantenimiento puedes consultar el espacio utilizado y retirar solo duplicados antiguos verificados. Las copias locales de recuperación siguen en **Mis listas → Guardado y copia adicional → Actividad reciente, historial y copias de recuperación**.

[Guía de guardado y recuperación](https://github.com/cristianbarallobre/selector-cine-365/blob/main/docs/GUIA_SINCRONIZACION_SEGURA.md) · [Rutas de datos](https://github.com/cristianbarallobre/selector-cine-365/blob/main/docs/UBICACIONES_DATOS_SELECTOR_CINE_365.md)

## Historial de actividad

En **Mis listas → Guardado y copia adicional → Actividad reciente, historial y copias de recuperación** consulta las últimas modificaciones confirmadas. Puedes buscar por película, lista o acción, filtrar por tipo y recorrer páginas de 20 movimientos. Se conservan hasta 10.000 registros locales. Los movimientos antiguos de visionado se conservan cuando están disponibles; no se reconstruyen acciones que nunca fueron registradas. El historial se guarda físicamente en instantáneas verificadas dentro de historial_local, con un espejo de recuperación. Las exportaciones completas y las copias adicionales incluyen los movimientos. Al restaurar una copia que contiene historial, se recupera también su historial; con una copia antigua sin historial se conserva el actual.

## Diagnóstico y cierre seguro (escritorio 9.2.18–9.2.30)

Desde 9.2.20 se conserva un borrador local de los cambios pendientes cuando puede escribirse; se ofrece recuperar copias válidas si falta o se daña la biblioteca principal. Si el guardado no se confirma, revisa los cambios pendientes, reintenta o exporta una copia JSON de emergencia. Es posible cerrar voluntariamente con advertencia de pérdida de cambios. En Mis listas → Mantenimiento → Exportar registro de diagnóstico puedes guardar los registros técnicos locales para revisar incidencias. No se envían automáticamente ni incluyen claves de acceso o contenido de listas.

Los archivos originales se conservan en `logs` dentro de la carpeta de datos: `%APPDATA%/Selector_Cine_365/logs` en Windows y `~/Library/Application Support/Selector_Cine_365/logs` en macOS. El registro del motor y el de la ventana nativa de Windows rotan al alcanzar 2 MiB; se mantienen hasta dos copias anteriores por archivo.

Autor: Cristian Barallobre Rúa.

## Salvaguardas de la copia adicional y del historial


Cambiar solo de lista activa conserva la selección en el equipo, pero no reescribe la copia adicional cuando su contenido e historial son iguales. La fecha de la copia corresponde al archivo guardado. Si la copia actual está dañada, se conserva su contenido en un archivo `colecciones_danada_<huella>.json` antes de reconstruirla desde la biblioteca principal. Si no se puede conservar el original, no se sustituye. Los archivos dañados conservados quedan disponibles para revisión manual.


Si la biblioteca es válida pero su historial no está disponible, se buscan copias completas que contengan exactamente el historial comprometido, verificado por su huella SHA-256. Si no se encuentra ninguna, no se descarta historial ni se sustituye automáticamente la biblioteca; conserva los archivos y revisa una copia completa antes de intervenir.
