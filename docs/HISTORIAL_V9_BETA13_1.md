> **ARCHIVO HISTÓRICO:** esta documentación conserva el estado de una versión anterior. La referencia de instalación y soporte actual es [Selector Cine 365 v9.0.0 estable](https://github.com/cristianbarallobre/selector-cine-365/releases/tag/v9.0.0) y su [manual completo](MANUAL_COMPLETO_SELECTOR_CINE_365.md). No seguir instrucciones de sincronización antiguas sin revisar las nuevas.

# Selector Cine 365 — Historial completo de versiones

**Autor del software:** Cristian Barallobre Rúa  
**Proyecto:** [cristianbarallobre/selector-cine-365](https://github.com/cristianbarallobre/selector-cine-365)

> Documento de evolución funcional y técnica. La edición v9 Beta 13.1 es una **versión de pruebas**, no una versión estable. Esta fuente en Markdown genera las ediciones Word y PDF con el mismo contenido. El historial clásico v8.1.1 se conserva intacto por motivos históricos.

## 1. Situación actual y distribución

| Plataforma | Edición publicada | Tecnología | Estado |
|---|---|---|---|
| Windows 10/11, 64 bits | v9 Beta 13.1 | Go y WebView2; interfaz HTML compartida | Beta publicada |
| macOS, Apple Silicon | v9 Beta 13.1 | Go y WKWebView; interfaz HTML compartida | Beta publicada, sin notarización |
| macOS, Intel | v9 Beta 13.1 | Go y WKWebView; interfaz HTML compartida | Beta publicada, sin notarización |


## 2. Evolución de la serie v9

### v9 Beta 13.1 — 20 de noviembre de 2026 (Windows y macOS)

- Tras abrir Windows Beta 13 se podía mostrar indefinidamente «Conexión local TMDB no disponible. Reintentando…», aunque la clave TMDB personal estuviera guardada. La ruta `GET /api/platforms/cache` estaba implementada en el módulo TMDB, pero **no figuraba en el enrutador principal**; devolvía HTTP 404.
- Se conecta esa ruta al servidor local y se separa la comprobación de credenciales de la carga de caché. Un fallo de caché no cambia la configuración y el botón Actualizar puede recuperar el estado.
- Se incorporan pruebas de las rutas HTTP: estado y caché autenticados responden 200, las peticiones sin token obtienen 403.
- Los binarios Mac Intel y Apple Silicon se reconstruyen a partir del mismo código Go y HTML de Windows Beta 13.1, no desde Beta 13.
- La publicación `v9.0.0-beta.13.1` reúne cuatro instaladores y la documentación actualizada. Es una beta: la sincronización real prolongada aún debe validarse.

### v9 Beta 13 — 20 de noviembre de 2026

- Estabilidad visual: cuatro posiciones en Mis listas de escritorio, sin saltos entre páginas.
- Sincronización: las listas eliminadas registran una marca persistente que impide su reintroducción desde snapshots antiguos de otros dispositivos.
- Identificaciones TMDB corregidas manualmente: se guardan como ID en la biblioteca compartida; la clave API personal nunca se publica ni se sincroniza.
- Se reintenta leer el estado de TMDB al reabrir, en vez de confundir errores temporales con pérdida de configuración.
- Regresiones y compilaciones verificadas para Windows; falta validar sincronización real con ambos dispositivos simultáneamente.

### v9 Beta 12.2 — 20 de noviembre de 2026

- Se mantienen todas las funciones anteriores; Windows y Mac continúan en Beta 12.
- Compilación y pruebas automáticas correctas; aún falta comprobación visual física de la zona inferior.

### v9 Beta 12.1 — 20 de noviembre de 2026

- La altura extra del menú inferior introducida por la doble reserva se corrigió después en Beta 12.2.

### v9 Beta 12 — 20 de noviembre de 2026

**Windows y macOS (mismo código de escritorio):**

- Se corrigió el espacio de «Mis listas», cuyos elementos podían quedar ocultos por la zona de nube e historial.
- «Actualizar historial» comunica que está consultando y después muestra la hora de lectura y la cantidad de eventos encontrados.
- El mensaje de historial sin movimientos usa un formato discreto, sin recuadro discordante.
- Los controles de paginación en Estadísticas permanecen al pie de los paneles.
- Las compilaciones macOS Intel y Apple Silicon recuperan las mismas funciones de escritorio que Windows y se distribuyen en una release conjunta.


- Búsqueda dentro de la lista actual o entre todas las listas; resultados globales agrupados por identidad TMDB o título y año.
- Numeración automática de listas; cada ficha indica «Listas: 1, 2, …» y permite abrirlas por número.
- Edición de pertenencia de una película a otras listas manteniendo el estado de Vista/Pendiente.
- Año y valoración **TMDB** juntos cuando se dispone de datos; esta valoración **no es IMDb**.
- Únicamente el tick cambia el estado de Vista/Pendiente; pulsar el título o el cartel no debe hacerlo.
- Continúan el modo de almacenamiento local, las copias JSON y el mecanismo de sincronización y revisión causal de Beta 10.


### v9 Beta 11 — 20 de noviembre de 2026

- Se introdujeron el buscador global, listas numeradas y modificación de la pertenencia de películas desde sus fichas.
- Windows incorporó correcciones de distribución de Estadísticas y de las fichas.
- Se conservaron el historial, recuperación, sincronización causal, Conectar TMDB y el modo local de Beta 10.
- En esta etapa la firma y empaquetado se verificaron mediante pruebas automáticas y navegadores; estas verificaciones no sustituyen una prueba real de uso en cada dispositivo.

### v9 Beta 10 — 20 de noviembre de 2026

- Se reforzó la conciliación de modificaciones entre Windows mediante revisiones causales.
- Se corrigieron casos en los que reaparecían marcas antiguas tras desmarcar una película desde otro dispositivo.
- Se añadieron historial de acciones locales/recibidas, información de última comprobación, botón «Sincronizar ahora» y vías de recuperación de copias.
- Se respetó la interacción por tick y la sección de Instrucciones de escritorio se alineó con el ancho de los demás paneles.

### v9 Beta 9 — 20 de noviembre de 2026

- Edición de referencia anterior a la búsqueda global y a la consolidación de Beta 12.
- Se consolidaron el comportamiento de cabecera y desplazamiento móvil, las estadísticas y la estabilidad de sincronización.
- Fue publicada por separado en Windows; sus dos instaladores originales se conservan juntos en la [release histórica Beta 9](https://github.com/cristianbarallobre/selector-cine-365/releases/tag/v9.0.0-beta.9).

### v9 Beta 8 y Beta 7 — 20 de noviembre de 2026

- Beta 8 amplió la comprobación automática de OneDrive, manteniendo las primeras conexiones bajo revisión de usuario.
- Estas betas continuaron la integración móvil, el guardado local y la recuperación de marcas desde listas.
- Los paquetes originales se encuentran en el archivo histórico de instaladores, no en las cuatro releases públicas actuales.

### v9 Betas 1 a 6 — 19 y 20 de noviembre de 2026

- Incorporación progresiva de la integración TMDB, la conexión OneDrive y el trabajo de compatibilidad entre escritorio y teléfono.
- Mejoras del icono, presentación, estados compartidos y herramientas de copia y restauración.
- La serie era experimental; estos paquetes no se consideran intercambiables ni garantizan compatibilidad directa de firma con Beta 12.

## 3. Versiones anteriores a v9

| Versión | Evolución principal |
|---|---|
| **v8.2** | Última edición publicada como estable antes de la serie v9; Windows y macOS. |
| **v8.1.1** | Autoría, firma visual y metadatos en las aplicaciones. |
| **v8.1** | Identificación y conciliación de una misma película entre listas. |
| **v8.0** | Intersección/coincidencias entre listas y estados compartidos. |
| **v7.9** | Eliminación individual de películas. |
| **v7.8** | Navegación rápida entre páginas y listas. |
| **v7.7** | Importación acumulativa de listas. |
| **v7.6** | Reorganización de la interfaz. |
| **v7.5** | Mejoras de TMDB y gestión de la nube. |
| **v7.4** | Nuevas opciones de copias de seguridad. |
| **v7.3** | Entorno de escritorio WebView2 operativo. |
| **v7.0–v7.2** | Pruebas de integración de WebView2. |
| **v6.1** | Filtros en la biblioteca. |
| **v6.0** | Biblioteca ampliada. |
| **v5.0** | Integración inicial de plataformas de visualización. |
| **v4.0** | Interfaz de organización de películas. |
| **v3.0** | Diseño oscuro. |
| **v2.0** | Mejoras de estabilidad. |
| **v1.0** | Primer selector con sorteo y marca de vista. |

Las descripciones de las series v1 a v8 proceden del registro conservado en **CHANGELOG.md**. Los documentos de v8.1.1 no se sustituyen: forman parte de la historia del proyecto.

## 4. Sincronización: capacidades y límites


El mecanismo v9 guarda revisiones y relaciones causales para decidir qué modificación conoce a cuál. No se limita a comparar la hora del reloj: dos cambios independientes, realizados sin haberse recibido mutuamente, pueden necesitar revisión antes de aceptar uno.

**Funciones disponibles:** sincronización y comprobaciones periódicas mientras las apps están activas, «Sincronizar ahora», «Revisar sincronización», historial local de eventos, exportación JSON y recuperación de copias.

**Limitaciones:** la propagación por OneDrive no es instantánea; desconexiones, modificaciones simultáneas, borrados o renombrados concurrentes y restauraciones requieren precaución. No se ha concluido una campaña prolongada de pruebas en todos los dispositivos físicos. Estas versiones son betas.

«Actualizar historial» solo vuelve a consultar el registro **local** de eventos almacenados por las versiones compatibles; no sustituye a «Sincronizar ahora» ni recrea eventos anteriores a la función.

## 5. Actualización y protección de datos

Antes de actualizar, guarda una copia **JSON fuera del almacenamiento privado de la app** y verifica que contiene tus listas.

- **Windows:** cierra el programa anterior y sustituye el ejecutable; conserva la carpeta de datos.
- **macOS:** cierra la aplicación y sustituye su paquete en Aplicaciones; Mac conserva la ruta histórica de almacenamiento y no está notarizada por Apple Developer ID.


## 6. Publicación, archivo y recuperación histórica

En la limpieza inicial se conservaron cuatro publicaciones públicas. Posteriormente se publicaron las betas correctivas; la [release v9 Beta 13.1](https://github.com/cristianbarallobre/selector-cine-365/releases/tag/v9.0.0-beta.13.1) es la edición coordinada actual. Los instaladores históricos y las releases anteriores siguen disponibles en GitHub.

Los instaladores de 17 releases redundantes se guardaron byte por byte en una **release borrador** accesible para los mantenedores, con manifiesto SHA-256, y se verificó la copia remota antes de retirar los enlaces públicos. Las ramas, etiquetas Git y commits originales permanecen conservados.


## 7. Documentación de referencia

- [README — versión actual y descargas](../README.md).
- [CHANGELOG — evolución compacta](../CHANGELOG.md).
- [Manual completo v9 Beta 13.1](MANUAL_COMPLETO_SELECTOR_CINE_365.md).
- [Guía macOS](../MACOS.md).
- [Guía Microsoft Entra y OneDrive](GUIA_AZURE_ENTRA_ONEDRIVE_SELECTOR_CINE_365.md).
- [Plan técnico v9](PLAN_V9_MULTIPLATAFORMA_Y_SINCRONIZACION.md).
- [Historial histórico hasta v8.1.1](HISTORIAL_VERSIONES.md).
- [Política de releases](POLITICA_DE_RELEASES.md).

**Autor y titular indicado en el proyecto:** Cristian Barallobre Rúa. © 2026. TMDB y Microsoft OneDrive son servicios externos; el proyecto no está certificado ni patrocinado por TMDB. 

*Fin del historial completo · v9 Beta 13.1.*
