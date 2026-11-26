# Manual de Selector Cine 365

Windows y macOS 9.2.19. Android estable de la publicación v9.0.0.

## Instalar y proteger la biblioteca

1. Elige la descarga de tu sistema desde el README de GitHub. En Mac comprueba si tu procesador es Intel o Apple Silicon.
2. Si ya usas la app, exporta una copia JSON y cierra la versión anterior.
3. En Windows abre el nuevo EXE portátil. Si el equipo lo necesita, instala WebView2 Runtime. En Mac descomprime el ZIP, lleva la app a Aplicaciones y reemplaza la anterior.
4. Si Windows o Mac avisan sobre el editor, verifica primero que el archivo procede de la publicación oficial y comprueba su SHA-256. Los paquetes Mac llevan firma ad hoc y carecen de notarización.
5. En Android instala la APK Limpia con el mismo identificador y la misma firma que tu instalación compatible. Para cambiar de edición, importa una copia JSON en la nueva.

[Descargas de escritorio](https://github.com/cristianbarallobre/selector-cine-365/releases/tag/v9.2.19-desktop) · [APK oficial de Android](https://github.com/cristianbarallobre/selector-cine-365/releases/tag/v9.0.0)

## Biblioteca y listas

- Abre el selector superior para elegir una lista. Puedes cambiar de lista desde Estadísticas y mantener abierta la sección.
- Usa Buscar y elige **Lista actual** o **Todas las listas**. Los filtros de escritorio permiten combinar el estado vista o pendiente con plataformas, década, género y duración, y elegir el orden.
- Usa la paginación o el campo **Página** para saltar. Al sincronizar manualmente con OneDrive, la app conserva la sección, la búsqueda y la página si la lista activa sigue disponible.
- Marca el estado desde la casilla. El indicador **Listas** muestra las colecciones que contienen la película. Pulsa un número para ir a esa lista o el lápiz para editar la pertenencia.
- Si quitas todas las listas numeradas y la película permanece en Favoritos, puedes guardar directamente. Si va a desaparecer de todas las listas, la app pide confirmación.
- Desde **Mis listas** puedes importar o añadir películas, gestionar colecciones y configurar la comparación de listas. En escritorio puedes reordenar arrastrando el tirador cuando estás viendo el orden original de la lista.

## Favoritos y valoraciones en escritorio

Pulsa el corazón para añadir una película. Aparece la lista **♥ Favoritos** en el selector y el corazón pequeño junto a los números te lleva a ella. La lista se oculta al quitar el último favorito. Si el mismo ID de TMDB identifica varios títulos de una película, el favorito y la valoración se comparten.

Al añadir la película, se marca como vista donde ya figure. Si después la importas desde un TXT, puede heredar esa marca desde Favoritos. En Favoritos se muestra como vista y la casilla queda bloqueada. Puedes ordenar Favoritos con el tirador y el orden se guarda. Las películas que solo están en Favoritos también pueden cargar su ficha. La valoración personal, de uno a cinco, es independiente de la nota pública de TMDB.

## Sorteo y estadísticas

El sorteo normal selecciona películas pendientes de la lista activa. En **Sortear entre** puedes marcar varias plataformas y entran las películas disponibles en cualquiera. **Cualquier plataforma** limita el sorteo a títulos con streaming registrado. El modo 🧠 **Sorteo inteligente** permite escoger vistas, pendientes o ambas y filtrar por género y duración. Las estadísticas reflejan la lista activa y se completan con los datos de TMDB disponibles.

## TMDB y Descubrir

Configura una clave de TMDB en la app si quieres carteles, disponibilidad en España y fichas. **Descubrir** permite buscar películas y personas, explorar el catálogo y abrir las fichas. Al pasar el puntero por un actor o director puedes ver su retrato. Pulsa su nombre para abrir la biografía y las películas en las que actuó o que dirigió. Los carteles grandes y retratos tienen una caché local de hasta 32 MiB o 640 imágenes. La disponibilidad puede cambiar. Pulsa **Actualizar** en Biblioteca para volver a consultar las fichas y plataformas de la lista activa; los datos en caché no caducan automáticamente. Conviene comprobar la disponibilidad en el servicio.

Al añadir una película desde Descubrir, Biblioteca reutiliza la ficha y el cartel de TMDB guardados localmente. Al borrar una búsqueda de Descubrir, se recuperan la categoría, la página y el filtro de género anteriores.

Si la app ha elegido otra película, abre **Identificar película en TMDB** y pega su ID numérico o su enlace. Si los años difieren, verifica el título y pulsa **Guardar de todos modos** si la ficha es correcta. El año de tu lista se conserva. **Volver a búsqueda automática** retira la identificación manual.

## Copias y OneDrive

La biblioteca se guarda automáticamente en el equipo. Desde **Mis listas → Copia de seguridad** en escritorio o **Ajustes** en Android, exporta un JSON para conservar todas las listas. Una exportación TXT corresponde a la lista seleccionada.

1. En escritorio abre **Mis listas → OneDrive y copias de seguridad → Conectar cuenta Microsoft**. Revisa los recuentos y confirma la primera vinculación.
2. Pulsa **Sincronizar con la nube**, debajo del título, para enviar y recibir cambios. La ventana de progreso bloquea la edición hasta terminar.
3. **Subida pendiente** indica que quedan cambios locales por enviar. **Recepción pendiente** indica que la app encontró cambios remotos por aplicar. **Recepción sin comprobar** indica que todavía no pudo verificar si hay cambios remotos.
4. Si surge un conflicto, revisa los cambios y elige qué estado conservar. Con una carpeta local de sincronización, usa **Mis listas → OneDrive y copias de seguridad → Sincronizar ahora**.

Si una edición no se guardó por un conflicto local, carga la biblioteca guardada y repite solo el cambio que falte. Las copias de recuperación están en **Mis listas → OneDrive y copias de seguridad → Actividad reciente, historial y copias de recuperación**. Restaurar una copia sustituye la biblioteca completa. Las claves TMDB y la caché permanecen en cada dispositivo.

[Guía de sincronización](https://github.com/cristianbarallobre/selector-cine-365/blob/main/docs/GUIA_SINCRONIZACION_SEGURA.md) · [Rutas de datos](https://github.com/cristianbarallobre/selector-cine-365/blob/main/docs/UBICACIONES_DATOS_SELECTOR_CINE_365.md)

## Android

La edición pública se identifica como **Beta 13.13 Limpia**. Tiene biblioteca, navegación táctil, estadísticas, sorteo, copias y sincronización mediante revisión manual. La edición Android 9.2.5 es una versión de pruebas para Android 8 o posterior. Se instala en paralelo con la estable. Una instalación de pruebas con firma distinta requiere reinstalar esa copia. [Descarga de Android 9.2.5](https://github.com/cristianbarallobre/selector-cine-365/releases/tag/v9.2.5-android).

## Diagnóstico y cierre seguro (escritorio 9.2.18–9.2.19)

Si el guardado no se confirma, revisa los cambios pendientes, reintenta o exporta una copia JSON de emergencia. Es posible cerrar voluntariamente con advertencia de pérdida de cambios. En Mis listas → Mantenimiento → Exportar registro de diagnóstico puedes guardar los registros técnicos locales para revisar incidencias. No se envían automáticamente ni incluyen claves de acceso o contenido de listas.

Los archivos originales se conservan en `logs` dentro de la carpeta de datos: `%APPDATA%/Selector_Cine_365/logs` en Windows y `~/Library/Application Support/Selector_Cine_365/logs` en macOS. El registro del motor y el de la ventana nativa de Windows rotan al alcanzar 2 MiB; se mantienen hasta dos copias anteriores por archivo.

Autor: Cristian Barallobre Rúa.
