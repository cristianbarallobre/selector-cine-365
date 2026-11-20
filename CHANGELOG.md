# Historial de cambios — Selector Cine 365

**Versión coordinada v9 Beta 12 (20 de noviembre de 2026).** Los cambios recientes se basan en el código fuente v9, no en la rama antigua v8.2. Los paquetes Mac usan el mismo frontend y motor Go que Windows; no se afirma que se hayan validado manualmente en hardware real.

## v9 Beta 12 — Windows, macOS, Android (beta)

- Búsqueda global entre listas, sin duplicados por identidad TMDB o título/año; opción «Lista actual».
- Numeración dinámica de listas; fichas con «Listas: 1, 2…» y navegación/edición de pertenencia.
- Año y valoración **TMDB** (no IMDb) en las fichas. Únicamente el tick cambia el estado visto.
- Windows y Mac: paginador de estadísticas anclado al fondo; Mac conserva interfaz y lógica Go de Windows.
- Windows y Mac: filas de «Mis listas» visibles aunque la nube e historial ocupen el panel.
- «Actualizar historial» lee el historial local, informa de la consulta y no dibuja una caja extra al estar vacío.
- Android: mismo proyecto nativo de v9 con UI para móvil; mantiene cabecera fija, carga incremental y estadísticas sin paginar.
- No se modifica el motor de sincronización causal de Beta 10: se conservan revisiones, recuperación, copia local, historial y PKCE.

## v9 Beta 11 — Windows y Android (beta)

- Se introdujeron búsqueda global, listas numeradas y edición de pertenencias.
- Se corrigió la posición de los controles de paginación en Estadísticas de Windows.
- Android pasó al paquete permanente \`es.cine365.app\`, con \`versionCode=9011\`. La clave privada estable se conserva fuera de GitHub; versiones futuras necesitan la misma firma.

## v9 Beta 10 — Windows y Android (beta)

- Revisión de sincronización causal, protección contra la «resurrección» de marcas antiguas, historial y restauración.
- «Sincronizar ahora» y fecha de última comprobación; OneDrive Android mediante PKCE y Microsoft Graph Delta Query.
- Alineación de la sección Instrucciones y conservación de la interacción por tick.
- Android Beta 10 y betas anteriores usaban un identificador con sufijo distinto: requerían migración de datos al saltar al paquete permanente.

## v9 Betas 1–9 — Windows y Android (experimental)

- Primeras ediciones nativas de Android, caché TMDB, conexión OneDrive, sincronización automática, diseño móvil, controles de vista y documentación técnica.
- Beta 9 usaba \`es.cine365.app.v9beta9\`; no puede actualizarse directamente sobre el paquete nuevo conservando datos privados sin exportación/importación.

## Archivo v8 y anterior

Este registro recoge las principales funciones incorporadas a lo largo del desarrollo. Es un resumen; no pretende representar un historial de commits anteriores.

## v8.2 — Windows y macOS

- Primera versión de escritorio para Mac Apple Silicon y Mac Intel con Cocoa y WKWebView.
- Selector de carpetas para copias en macOS, compatible con rutas locales sincronizadas.
- Se mantienen las listas, sincronización de vistas, TMDB, estadísticas y la firma del autor.
- GitHub Actions compila tres ZIP de distribución: Windows, Mac Intel y Mac Apple Silicon.
- Los Mac se compilan y prueban de forma automatizada; resta comprobar su interfaz en un Mac de usuario.
- macOS: aplicación sin notarización ni certificado Apple Developer ID.

## v8.1.1 — Firma del autor (19 de noviembre de 2026)

- Firma discreta «Hecho por Cristian Barallobre Rúa» en la cabecera del programa.
- Metadatos de autoría en el HTML y en el servidor Go; constantes de autoría integradas en el ejecutable.
- Archivo `AUTORIA.md` para documentar la autoría y el repositorio canónico.
- Sin cambios en los datos personales, las colecciones, la conciliación de TMDB ni las copias de seguridad.

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