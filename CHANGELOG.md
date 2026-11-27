# Historial de cambios — Selector Cine 365

## 9.2.30 de escritorio

- Retirado el enlace a FilmAffinity de las fichas y su generador de URL. Se conservan TMDB y tráileres.
- Prueba de regresión para verificar la ausencia del enlace.
- Sin cambios en los datos ni en el guardado. No se ha confirmado la causa de los bloqueos intermitentes.

[Publicación 9.2.30](https://github.com/cristianbarallobre/selector-cine-365/releases/tag/v9.2.30-desktop)

## 9.2.29 de escritorio

- Protección del origen HTTP local y bloqueo de conexiones de la interfaz a orígenes ajenos.
- Reparación conservadora de copias adicionales, recuperación del historial exacto desde copias verificadas y ausencia de escrituras adicionales por navegación.
- Disponibilidad TMDB coherente por ID entre Biblioteca y Descubrir; refresco de fichas abiertas desde caché.
- Pruebas de fallos de almacenamiento y manuales actualizados.

[Publicación 9.2.29](https://github.com/cristianbarallobre/selector-cine-365/releases/tag/v9.2.29-desktop)

## 9.2.28 (candidata de escritorio)

- Indicadores minimalistas de disponibilidad en Biblioteca y antigüedad visible en Sorteo y fichas.
- Ayuda del interrogante revisada y manuales actualizados.
- Conserva las comprobaciones automáticas y manuales de plataformas.

Cambios incorporados en 9.2.29; la candidata pública no se publicó como estable.

## 9.2.27 (pruebas)

- Primera versión del indicador de última comprobación, con fecha exacta y distinción entre datos antiguos, errores y datos sin consultar.

## 9.2.26 de escritorio

- Historial de actividad separado del archivo principal, con referencia de integridad, respaldo local y recuperación protegida.
- Exportación y restauración conjunta de biblioteca e historial; copia adicional actual, anterior y hasta 14 diarias. Las copias anteriores sin historial siguen siendo compatibles.
- Retirada del motor y de los controles de conexión a cuentas externas: solo guardado local y carpeta adicional opcional.
- Disponibilidad TMDB con actualización en segundo plano, caducidad de 24 horas, conservación del último dato válido y control de saturación de consultas.
- Credenciales TMDB protegidas mediante DPAPI en Windows y llavero en macOS. Separación de caché regenerable en Windows.
- Mantenimiento: uso de disco, memoria del motor, historial y copias; limpieza opcional de archivos de caché antiguos con duplicado idéntico.
- Pruebas automatizadas de guardado, restauración, caché, interfaz en diferentes tamaños y compilación para Windows, macOS Intel y macOS Apple Silicon.
- Uso personal: sin firma comercial ni notarización de Apple; los paquetes Mac conservan la firma ad hoc técnica.

Cambios de la candidata 9.2.26 incorporados a la 9.2.28.

## 9.2.25 de escritorio

- Historial de actividad local con hasta 10.000 movimientos: altas y bajas de películas, cambios de estado, listas, favoritos, valoraciones, orden, identificaciones y restauración.
- Búsqueda por película, lista o acción, filtros por tipo y páginas de 20 movimientos.
- Se conservan los movimientos antiguos disponibles. Solo se registran cambios confirmados en el guardado principal.
- Edición exclusiva para Windows y macOS.

[Publicación 9.2.25](https://github.com/cristianbarallobre/selector-cine-365/releases/tag/v9.2.25-desktop)

## 9.2.24 de escritorio

- Se retira la conexión directa a OneDrive y la sincronización entre dispositivos en escritorio.
- Guardado principal local protegido y segunda copia opcional en una carpeta elegida por el usuario, con copia anterior, verificación y reintento.
- Restauración manual sin fusionar bibliotecas. Las copias anteriores se conservan.

[Publicación 9.2.24](https://github.com/cristianbarallobre/selector-cine-365/releases/tag/v9.2.24-desktop)

## 9.2.22 de escritorio

- Cambiar de lista conserva la selección en el equipo, sin activar una subida a OneDrive si no cambian los datos compartidos.
- Indicadores amarillos con reloj para operaciones pendientes y flecha circular para sincronización en curso.

[Publicación 9.2.22](https://github.com/cristianbarallobre/selector-cine-365/releases/tag/v9.2.22-desktop)

## 9.2.20 de escritorio

- Borrador persistente de la interfaz, tiempos límite para guardados y consultas, recuperación cuando falta la biblioteca principal y avisos de fallos de copias.
- Biblioteca y revisiones guardadas juntas, huella de archivo para detectar cambios externos y retirada de borradores confirmados.
- Conservación del historial de listas eliminadas sin el límite arbitrario de 1.000 registros.

[Publicación 9.2.20](https://github.com/cristianbarallobre/selector-cine-365/releases/tag/v9.2.20-desktop)

## 9.2.19 de escritorio

- Registro de diagnóstico local del motor y la ventana, con rotación y exportación desde Mantenimiento.
- Diagnóstico de guardado, recuperación, sincronización, TMDB y comunicación con la interfaz, sin registrar contenido de listas ni credenciales.

[Publicación 9.2.19](https://github.com/cristianbarallobre/selector-cine-365/releases/tag/v9.2.19-desktop)

## 9.2.18 de escritorio

- El motor local de Windows permanece activo mientras esté abierta la ventana nativa.
- Avisos de cambios sin confirmar, reintento de guardado, copia JSON de emergencia y cierre voluntario con advertencia.

[Publicación 9.2.18](https://github.com/cristianbarallobre/selector-cine-365/releases/tag/v9.2.18-desktop)

## 9.2.17 de escritorio

- Al añadir una película desde Descubrir, Biblioteca reutiliza el cartel y la ficha de TMDB guardados localmente.
- Al borrar una búsqueda, Descubrir recupera la categoría, la página y el filtro de género anteriores.
- Se conservan los carteles disponibles en caché al actualizar las fichas de Biblioteca.
- Compilaciones de Windows x64 y macOS Intel y Apple Silicon verificadas mediante GitHub Actions.

[Código fuente 9.2.17 (release archivada)](https://github.com/cristianbarallobre/selector-cine-365/tree/v9.2.17-desktop)

## 9.2.16 de escritorio

- La sincronización conserva las retiradas de películas aunque una copia anterior todavía las incluya. La película sigue en sus otras listas y puede añadirse de nuevo después.
- Si una versión anterior ya volvió a añadir la película, quítala de nuevo en la 9.2.16 y sincroniza.
- El estado de la nube aparece junto al botón de sincronizar.

[Código fuente 9.2.16 (release archivada)](https://github.com/cristianbarallobre/selector-cine-365/tree/v9.2.16-desktop)

## 9.2.15 de escritorio

- Sincronización manual con OneDrive en Windows y macOS. El botón de la cabecera muestra subida y recepción pendientes cuando hay una cuenta Microsoft conectada.
- La ventana de sincronización bloquea la edición mientras se aplican los cambios. El guardado local sigue siendo automático.
- Se reaplican las retiradas de películas compatibles con una revisión concurrente y se actualizan los números de lista al recargar.
- Los avisos de guardado explican qué hacer. Las copias anteriores se consultan en Mantenimiento.

[Código fuente 9.2.15 (release archivada)](https://github.com/cristianbarallobre/selector-cine-365/tree/v9.2.15-desktop) · [Guía de sincronización](docs/GUIA_SINCRONIZACION_SEGURA.md)

## 9.2.7 de escritorio

- Windows y macOS incorporan las correcciones de escritorio de la serie 9.2.
- Favoritos guarda su orden manual, carga fichas aunque el título solo figure allí y reconoce títulos equivalentes con el mismo ID de TMDB.
- La recepción automática conserva la lista, la búsqueda y la página. El selector se alinea con el panel de sorteo.
- Quitar una película de las listas numeradas conserva Favoritos. La identificación manual permite aceptar una diferencia de año.
- Windows utiliza el icono de la app. Mac guarda la caché aparte de la biblioteca y protege la sesión de OneDrive en el llavero.

[Código fuente 9.2.7 (release archivada)](https://github.com/cristianbarallobre/selector-cine-365/tree/v9.2.7-desktop) · [Manual completo](docs/MANUAL_COMPLETO_SELECTOR_CINE_365.md)

---

## v9.0.0 — estable (22 de noviembre de 2026)


- **Limitaciones:** las bibliotecas históricas ya contaminadas requieren recuperación manual desde copias verificadas; en Mac Microsoft aún no persiste mediante Keychain; la caché TMDB es local por dispositivo.
- **Mantenimiento:** documentación consolidada, archivos históricos etiquetados, releases redundantes respaldadas antes de retirada. No se han cambiado funcionalidades del ejecutable al promocionarlo a estable.

[Descargar estable](https://github.com/cristianbarallobre/selector-cine-365/releases/tag/v9.0.0) · [Manual completo](docs/MANUAL_COMPLETO_SELECTOR_CINE_365.md) · [Notas](RELEASE_NOTES_v9_ESTABLE.md).

---

## Etapa anterior: Beta 14 (21 de noviembre de 2026)

## v9 Beta 14 — candidato recomendado (21/11/2026)

- **Sincronización segura:** fusión por ID ya no añade a ciegas películas de listas distintas con el mismo `atlas`. Se bloquean listas sospechosas, incluyendo copias ya contaminadas con un gran excedente de películas.
- **Windows Beta 14:** sesión de OneDrive cifrada con DPAPI, vista previa y corrección de scroll/selector, exportación TXT y bloqueo de conflictos entre listas.
- **Mac Intel/Apple Silicon Beta 14:** la misma protección de listas de Go en las compilaciones; sesión Microsoft temporal, todavía sin Keychain.
- **Limitación:** los datos históricamente fusionados no se separan solos. Las betas no están certificadas con todas las cuentas y dispositivos reales. Se recomienda respaldo JSON antes de sincronizar.


---

## Historial anterior (archivo)


- Comprobación única de nube al abrir (lectura sin importar ni enviar). Se eliminan los intervalos, comprobaciones al reanudar y subida de restauraciones locales.
- Nuevo botón pequeño de **dos flechas circulares** junto al indicador de guardado; inicia comparación bidireccional con **confirmación previa obligatoria** y resolución de conflictos. Los guardados en el móvil siguen siendo automáticos.
- Se suprimen operaciones innecesarias de disco para datos idénticos y consultas TMDB automáticas al cambiar de lista.
- Pruebas de guardado local sin tráfico de nube, única consulta al inicio, no sincronización en segundo plano y conservación de la fusión causal.


- WebView suspendida cuando la actividad se minimiza; se reanuda al volver. Guardado local inmediato para ediciones pendientes y preservación de revisiones.
- Eliminado ping de 12 segundos; sincronización automática OneDrive cada 60 segundos mientras visible, además de al regresar, sin nuevas comprobaciones automáticas en segundo plano.
- La mejora real de batería requiere comprobarse en el móvil durante más tiempo; las capturas de ocho minutos no bastan para cuantificarla.

## v9 Beta 13.1 — Windows y Mac (20 de noviembre de 2026)

- Reparada la ruta local `GET /api/platforms/cache`, implementada en TMDB pero omitida del enrutador principal. Esto provocaba HTTP 404 y dejaba el mensaje «Conexión local TMDB no disponible. Reintentando…» al reabrir, aunque la clave estuviera guardada.
- La lectura de estado y caché se realiza por separado. Si falla la caché, no se interpreta que falte la clave ni se sobrescribe la configuración; el botón Actualizar reintenta.
- Añadidas pruebas automáticas de HTTP 200 para estado y caché, y HTTP 403 sin token; se conservaron las pruebas de listas y sincronización.
- Las cuatro plataformas y la documentación actualizada están incluidas en la release conjunta `v9.0.0-beta.13.1`.

## v9 Beta 13 — Windows y Mac (20 de noviembre de 2026)

- «Mis listas» en Windows/Mac reserva cuatro filas de altura constante por página, incluso si hay menos listas.
- Sincronización de los IDs manuales de TMDB (`tmdb_overrides`) dentro de la biblioteca. La clave API no se exporta y permanece privada en cada dispositivo.
- Reintento al leer configuración TMDB al abrir; una caída temporal no cambia el estado a «Sin configurar». Se comprueba persistencia de la clave en disco tras reiniciar.
- Pruebas causales de marcas, borrados en ambos sentidos, restauración de identificación e integridad de biblioteca; copias locales antes de fusiones.


- El menú vuelve al límite seguro inferior y se conserva la protección del título respecto a la cámara frontal.
- Mismo paquete `es.cine365.app` y certificado que Beta 11, 12 y 12.1; código de versión **9014**, instalable directamente encima sin desinstalar.


- Windows y Mac no cambian; se reutilizan sus paquetes de Beta 12.

## v9 Beta 12 — Windows y macOS (beta)

- Búsqueda global entre listas, sin duplicados por identidad TMDB o título/año; opción «Lista actual».
- Numeración dinámica de listas; fichas con «Listas: 1, 2…» y navegación/edición de pertenencia.
- Año y valoración **TMDB** (no IMDb) en las fichas. Únicamente el tick cambia el estado visto.
- Windows y Mac: paginador de estadísticas anclado al fondo; Mac conserva interfaz y lógica Go de Windows.
- Windows y Mac: filas de «Mis listas» visibles aunque la nube e historial ocupen el panel.
- «Actualizar historial» lee el historial local, informa de la consulta y no dibuja una caja extra al estar vacío.
- No se modifica el motor de sincronización causal de Beta 10: se conservan revisiones, recuperación, copia local, historial y PKCE.

## v9 Beta 11 — Windows (beta)

- Se introdujeron búsqueda global, listas numeradas y edición de pertenencias.
- Se corrigió la posición de los controles de paginación en Estadísticas de Windows.

## v9 Beta 10 — Windows (beta)

- Revisión de sincronización causal, protección contra la «resurrección» de marcas antiguas, historial y restauración.
- Alineación de la sección Instrucciones y conservación de la interacción por tick.

## v9 Betas 1–9 — Windows (experimental)

- Beta 9 usaba `es.cine365.app.v9beta9`; no puede actualizarse directamente sobre el paquete nuevo conservando datos privados sin exportación/importación.

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
