# Historial de cambios — Selector de Cine 365

Este registro recoge las principales funciones incorporadas a lo largo del desarrollo. Es un resumen; no pretende representar un historial de commits anteriores.

## v8.1 — Identificación y sincronización entre listas

- Reconocimiento prioritario de películas equivalentes mediante el identificador de TMDB.
- Reconciliación de marcas de «Vista» y «Pendiente» entre colecciones existentes e importadas.
- Revisión de identificaciones pendiente y conservación de resultados en `identidades_cine.json`.
- Mejoras en la colección virtual de coincidencias.

## v8.0 — Coincidencias y estados compartidos

- Colección virtual de películas presentes en varias listas seleccionadas.
- Sincronización inicial de marcas entre listas, por título normalizado y año.
- Ventana de eliminación de tamaño fijo y búsqueda vacía hasta introducir texto.

## v7.9 — Eliminación individual

- Buscador en «Mis listas» para seleccionar y eliminar una película.
- Confirmación antes del borrado.

## v7.8 — Navegación rápida

- Más números de página, saltos por bloques e introducción directa del número de página.

## v7.7 — Importación acumulativa

- Añadir varias películas de un TXT a una colección ya existente, sin sustituirla.
- Detección de duplicados por título y año.

## v7.6 — Organización de interfaz

- Traslado de las herramientas de TXT y de añadir películas a «Mis listas».
- Reorganización sin desplazamiento general y paginación de biblioteca.

## v7.5 — Mejoras en TMDB y nube

- Búsqueda por variantes de título y corrección manual de ID TMDB.
- Actualización visual del estado de la copia en nube.

## v7.4 — Copias adicionales

- Copias automáticas en una carpeta local sincronizada (OneDrive, Dropbox, etc.).
- Historial diario de seguridad, independiente del guardado principal en AppData.

## v7.3 — WebView2 operativo

- Ventana nativa de Windows con WebView2, icono propio en la barra de tareas y diseño web conservado.
- Versiones intermedias de WebView2 fueron experimentales y presentaron errores de inicialización.

## v7.0–v7.2 — Pruebas de WebView2

- Primeros intentos de ventana independiente de Chrome y mejoras de diagnóstico de arranque.
- No deben considerarse versiones estables.

## v6.1 — Filtros en la biblioteca

- Filtro por una o varias plataformas, independiente del selector usado para el sorteo.

## v6.0 — Biblioteca ampliada

- Carteles y puntuaciones TMDB, estadísticas, sorteo por plataforma, múltiples listas y copias de seguridad.

## v5.0 — Plataformas

- Consulta de disponibilidad en España usando TMDB / JustWatch, con diferenciación entre suscripción y alquiler o compra.

## v4.0 — Interfaz de colección

- Películas en filas clicables, buscador, filtros por estado y nuevo icono de claqueta.

## v3.0 — Diseño oscuro

- Interfaz oscura y dorada utilizando un navegador en modo aplicación.

## v2.0 — Estabilidad

- Corrección de bloqueos de la primera versión del ejecutable.

## v1.0 — Primer selector

- Sorteo de películas pendientes y marcado de películas vistas.