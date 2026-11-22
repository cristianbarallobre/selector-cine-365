<!-- Documento histórico v8.1.1: no modificar la narración de esa edición como si describiera v9. -->

> **ARCHIVO HISTÓRICO — documentación cerrada en v8.1.1 (19 de noviembre de 2026).** No representa las funciones de la versión estable v9.0.0. Para Windows, Android y Mac consulta el [manual v9](MANUAL_COMPLETO_SELECTOR_CINE_365.md), el [historial reciente](../CHANGELOG.md) y la [página principal](../README.md).


<img src="../icon.png" style="width:0.80709in;height:0.80709in" />

**DOCUMENTACIÓN DEL PRODUCTO / REGISTRO DE CAMBIOS**

**Selector de Cine 365**

**Historial de versiones y documentación funcional**

<table>
<colgroup>
<col style="width: 100%" />
</colgroup>
<thead>
<tr class="header">
<th><p><strong>ÚLTIMA VERSIÓN DOCUMENTADA</strong></p>
<p><strong>v8.1.1 · Windows / WebView2</strong></p>
<p>Firma de autor: Cristian Barallobre Rúa · Versión publicada en GitHub con ejecutable de Windows y código fuente.</p></th>
</tr>
</thead>
<tbody>
</tbody>
</table>

**AUTOR Y CREADOR Cristian Barallobre Rúa**

## Acerca de este documento

Este documento reúne la evolución de Selector de Cine 365 desde su primer prototipo hasta la versión pública 8.1.1, junto con su funcionamiento, estructura de datos, requisitos, medidas de respaldo y notas de mantenimiento. La edición 8.1.1 añade la identificación visible del autor y el registro de procedencia en el código.

| **EDICIÓN DEL DOCUMENTO 1.2 · 19 de noviembre de 2026** |
|---------------------------------------------------------|

| **ALCANCE Aplicación de escritorio para Windows; el proyecto Android se documenta por separado** |
|--------------------------------------------------------------------------------------------------|

| **BASE DOCUMENTAL Paquetes de código fuente, archivos LEEME y decisiones registradas durante el desarrollo** |
|--------------------------------------------------------------------------------------------------------------|

<table>
<colgroup>
<col style="width: 100%" />
</colgroup>
<thead>
<tr class="header">
<th><p><strong>Nota sobre el estado de las versiones</strong></p>
<p>La v7.3 fue confirmada por el usuario funcionando con WebView2 y el icono correcto en la barra de tareas. La v8.1.1 está publicada en GitHub y su compilación y pruebas automatizadas han finalizado correctamente. Sigue pendiente una comprobación final del ejecutable en un equipo Windows con datos reales.</p></th>
</tr>
</thead>
<tbody>
</tbody>
</table>

**01 / FICHA DEL PRODUCTO**

# Descripción general

Selector de Cine 365 es una aplicación personal para organizar colecciones de películas, marcar qué títulos se han visto, descubrir qué ver mediante sorteo y consultar su disponibilidad de streaming en España.

## Funciones actuales

- **Biblioteca y seguimiento.** Múltiples colecciones de películas; estado Vista/Pendiente, búsqueda, filtros, carteles, puntuaciones y paginación.

- **Descubrimiento.** Sorteo de películas pendientes con filtro independiente por plataforma, además de acceso a información de los títulos.

- **Disponibilidad.** Proveedores en España y metadatos desde TMDB, con enlaces de consulta y corrección manual de identificaciones.

- **Análisis.** Estadísticas de progreso y distribución por décadas, géneros, países, directores y tiempo de visionado.

- **Gestión de listas.** Crear, renombrar, importar, combinar, exportar y eliminar; comparación virtual de coincidencias entre dos o más listas.

- **Persistencia y respaldo.** Guardado local automático en AppData y copia opcional en una carpeta sincronizada con la nube.

## Ficha técnica resumida

| **Elemento**               | **Situación en v8.1.1**                                                                                                                                          |
|----------------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| **Plataforma**             | Windows 10 y 11, arquitectura x64                                                                                                                                |
| **Tecnologías**            | Go para la aplicación y servidor local; HTML, CSS y JavaScript para la interfaz; Microsoft WebView2 para la ventana.                                             |
| **Conexión**               | Uso básico de listas sin internet. Carteles, valoraciones y plataformas necesitan acceso a TMDB y una clave configurada.                                         |
| **Persistencia principal** | %APPDATA%\AtlasDelCine\colecciones.json                                                                                                                          |
| **Estado de publicación**  | v8.1.1 publicada en GitHub; pruebas Go y compilación Windows correctas en GitHub Actions. Ejecutable sin firma digital; comprobación final en Windows pendiente. |

**02 / REGISTRO DE CAMBIOS**

# Evolución: prototipo y versiones 1 a 6.1

<table>
<colgroup>
<col style="width: 14%" />
<col style="width: 85%" />
</colgroup>
<thead>
<tr class="header">
<th><strong>Prototipo</strong></th>
<th><p><strong>Selector HTML · Origen</strong></p>
<p>Primera interfaz autocontenida en el navegador, con la lista de 365 películas y las marcas iniciales; sorteo de pendientes y edición o exportación de TXT.</p></th>
</tr>
</thead>
<tbody>
</tbody>
</table>

<table>
<colgroup>
<col style="width: 14%" />
<col style="width: 85%" />
</colgroup>
<thead>
<tr class="header">
<th><strong>v1</strong></th>
<th><p><strong>Primer ejecutable para Windows · Primera entrega</strong></p>
<p>Aplicación Go con interfaz básica de Windows. Sorteo de películas pendientes, marcado de vistas, importación/exportación TXT y guardado local.</p></th>
</tr>
</thead>
<tbody>
</tbody>
</table>

<table>
<colgroup>
<col style="width: 14%" />
<col style="width: 85%" />
</colgroup>
<thead>
<tr class="header">
<th><strong>v2</strong></th>
<th><p><strong>Corrección de estabilidad · Corrección</strong></p>
<p>Se corrigió el bloqueo «No responde» fijando la ventana y su bucle de mensajes al hilo de Windows. Se conservó el archivo de progreso anterior.</p></th>
</tr>
</thead>
<tbody>
</tbody>
</table>

<table>
<colgroup>
<col style="width: 14%" />
<col style="width: 85%" />
</colgroup>
<thead>
<tr class="header">
<th><strong>v3</strong></th>
<th><p><strong>Diseño visual oscuro y dorado · Rediseño</strong></p>
<p>Se reutilizó la estética del prototipo HTML en una ventana abierta con Edge o Chrome en modo aplicación; se mantuvieron el sorteo y el guardado.</p></th>
</tr>
</thead>
<tbody>
</tbody>
</table>

<table>
<colgroup>
<col style="width: 14%" />
<col style="width: 85%" />
</colgroup>
<thead>
<tr class="header">
<th><strong>v4</strong></th>
<th><p><strong>Biblioteca interactiva · Usabilidad</strong></p>
<p>Sustitución del bloque TXT visible por filas clicables de películas; buscador, filtros Todas/Pendientes/Vistas e icono de claqueta.</p></th>
</tr>
</thead>
<tbody>
</tbody>
</table>

<table>
<colgroup>
<col style="width: 14%" />
<col style="width: 85%" />
</colgroup>
<thead>
<tr class="header">
<th><strong>v5</strong></th>
<th><p><strong>Información de plataformas en España · Integración</strong></p>
<p>Integración con TMDB y datos de disponibilidad; separación entre suscripción y alquiler/compra, enlaces a JustWatch, actualización y caché de consultas.</p></th>
</tr>
</thead>
<tbody>
</tbody>
</table>

<table>
<colgroup>
<col style="width: 14%" />
<col style="width: 85%" />
</colgroup>
<thead>
<tr class="header">
<th><strong>v6</strong></th>
<th><p><strong>Biblioteca ampliada · Funciones</strong></p>
<p>Carteles pequeños y nota de TMDB; sorteo por plataforma; estadísticas, colecciones múltiples y copias de seguridad exportables.</p></th>
</tr>
</thead>
<tbody>
</tbody>
</table>

<table>
<colgroup>
<col style="width: 14%" />
<col style="width: 85%" />
</colgroup>
<thead>
<tr class="header">
<th><strong>v6.1</strong></th>
<th><p><strong>Filtro de plataformas en la lista · Filtros</strong></p>
<p>Filtro adicional para mostrar películas de una, varias o todas las plataformas; compatible con búsqueda y con Vistas/Pendientes, independiente del sorteo.</p></th>
</tr>
</thead>
<tbody>
</tbody>
</table>

**02 / REGISTRO DE CAMBIOS**

# Evolución: adopción de WebView2

<table>
<colgroup>
<col style="width: 14%" />
<col style="width: 85%" />
</colgroup>
<thead>
<tr class="header">
<th><strong>v7.0</strong></th>
<th><p><strong>Primera ventana nativa WebView2 · Experimental</strong></p>
<p>Intento de sustituir Chrome por una ventana independiente de Windows para conservar el icono de claqueta en la barra de tareas. La inicialización presentaba cierres y pantalla negra.</p></th>
</tr>
</thead>
<tbody>
</tbody>
</table>

<table>
<colgroup>
<col style="width: 14%" />
<col style="width: 85%" />
</colgroup>
<thead>
<tr class="header">
<th><strong>v7.1</strong></th>
<th><p><strong>Arranque aislado y alternativa compatible · Diagnóstico</strong></p>
<p>Se aisló WebView2 en un proceso y se incorporó un modo de respaldo mediante Chrome/Edge y un registro de diagnóstico.</p></th>
</tr>
</thead>
<tbody>
</tbody>
</table>

<table>
<colgroup>
<col style="width: 14%" />
<col style="width: 85%" />
</colgroup>
<thead>
<tr class="header">
<th><strong>v7.2</strong></th>
<th><p><strong>Trazas de navegación · Diagnóstico</strong></p>
<p>Se añadieron registros de eventos de WebView2, navegación y comprobaciones del servidor local para identificar por qué no cargaba la interfaz.</p></th>
</tr>
</thead>
<tbody>
</tbody>
</table>

<table>
<colgroup>
<col style="width: 14%" />
<col style="width: 85%" />
</colgroup>
<thead>
<tr class="header">
<th><strong>v7.3</strong></th>
<th><p><strong>Corrección funcional de WebView2 · Confirmada</strong></p>
<p>Se ajustó el ciclo de vida de referencias COM y se difirió la navegación. El usuario confirmó que cargaba correctamente y conservaba el icono propio al anclarla.</p></th>
</tr>
</thead>
<tbody>
</tbody>
</table>

<table>
<colgroup>
<col style="width: 14%" />
<col style="width: 85%" />
</colgroup>
<thead>
<tr class="header">
<th><strong>v7.3 Def.</strong></th>
<th><p><strong>Identificación de edición · Rotulación</strong></p>
<p>Se actualizó la etiqueta visible y el título de la ventana a «v7.3 · Definitiva», manteniendo el resto de la integración estable.</p></th>
</tr>
</thead>
<tbody>
</tbody>
</table>

<table>
<colgroup>
<col style="width: 14%" />
<col style="width: 85%" />
</colgroup>
<thead>
<tr class="header">
<th><strong>v7.4</strong></th>
<th><p><strong>Copia automática en carpeta de nube · Respaldo</strong></p>
<p>Segunda copia de colecciones en una carpeta local sincronizada (OneDrive, Drive, Dropbox u otra), con carpeta «actual» e histórico diario.</p></th>
</tr>
</thead>
<tbody>
</tbody>
</table>

<table>
<colgroup>
<col style="width: 14%" />
<col style="width: 85%" />
</colgroup>
<thead>
<tr class="header">
<th><strong>v7.5</strong></th>
<th><p><strong>Mejor identificación y hora de respaldo · Datos</strong></p>
<p>Búsqueda TMDB con variantes de título y año; corrección manual mediante ID; actualización visible de la hora al terminar el guardado en la nube.</p></th>
</tr>
</thead>
<tbody>
</tbody>
</table>

**02 / REGISTRO DE CAMBIOS**

# Evolución: versiones 7.6 a 8.0

<table>
<colgroup>
<col style="width: 14%" />
<col style="width: 85%" />
</colgroup>
<thead>
<tr class="header">
<th><strong>v7.6</strong></th>
<th><p><strong>Reorganización sin scroll general · Interfaz</strong></p>
<p>Las acciones de TXT y añadir películas pasaron a «Mis listas»; la biblioteca y las estadísticas se reorganizaron con paginación.</p></th>
</tr>
</thead>
<tbody>
</tbody>
</table>

<table>
<colgroup>
<col style="width: 14%" />
<col style="width: 85%" />
</colgroup>
<thead>
<tr class="header">
<th><strong>v7.7</strong></th>
<th><p><strong>Páginas numeradas y combinación de listas · Gestión</strong></p>
<p>Botones numéricos para pasar de página; función para incorporar por TXT o pegado varias películas a una colección existente sin reemplazarla.</p></th>
</tr>
</thead>
<tbody>
</tbody>
</table>

<table>
<colgroup>
<col style="width: 14%" />
<col style="width: 85%" />
</colgroup>
<thead>
<tr class="header">
<th><strong>v7.8</strong></th>
<th><p><strong>Paginación rápida · Navegación</strong></p>
<p>Hasta 11 páginas consecutivas visibles, saltos de ±10 páginas y acceso directo escribiendo el número de página.</p></th>
</tr>
</thead>
<tbody>
</tbody>
</table>

<table>
<colgroup>
<col style="width: 14%" />
<col style="width: 85%" />
</colgroup>
<thead>
<tr class="header">
<th><strong>v7.9</strong></th>
<th><p><strong>Eliminar una película concreta · Gestión</strong></p>
<p>Buscador de títulos dentro de la colección activa, selección y confirmación antes de eliminar; las demás listas no se modifican.</p></th>
</tr>
</thead>
<tbody>
</tbody>
</table>

<table>
<colgroup>
<col style="width: 14%" />
<col style="width: 85%" />
</colgroup>
<thead>
<tr class="header">
<th><strong>v8.0</strong></th>
<th><p><strong>Comparación y sincronización entre listas · Versión anterior</strong></p>
<p>Búsqueda de eliminación vacía al abrir y diálogo de tamaño constante; vista virtual de coincidencias entre listas seleccionadas; estado Vista/Pendiente sincronizado entre colecciones para el mismo título y año.</p></th>
</tr>
</thead>
<tbody>
</tbody>
</table>

## Detalles de la versión 8.0

- La comparación exige al menos dos listas y muestra la intersección: películas presentes en todas las colecciones seleccionadas.

- La colección «Películas coincidentes» se elige desde Biblioteca, es de solo lectura y se recalcula al cambiar los datos.

- El reconocimiento normaliza el título (acentos, signos y mayúsculas) y utiliza el año: títulos alternativos totalmente diferentes pueden no coincidir.

- Al iniciar la sincronización global, si una película figura vista en una lista y pendiente en otra, prevalece el estado Vista; conviene hacer copia antes de actualizar.

02 / REGISTRO DE CAMBIOS

# Evolución: versión 8.1

<table>
<colgroup>
<col style="width: 14%" />
<col style="width: 85%" />
</colgroup>
<thead>
<tr class="header">
<th><strong>v8.1</strong></th>
<th><p><strong>Identidad común y conciliación global · Versión anterior</strong></p>
<p>La coincidencia entre listas y la sincronización de películas vistas usan un criterio compartido: ID de TMDB cuando está disponible, con título normalizado y año como alternativa.</p></th>
</tr>
</thead>
<tbody>
</tbody>
</table>

## Cambios incorporados

- Coincidencias más completas: dos entradas con títulos diferentes pueden considerarse la misma película cuando comparten ID de TMDB; se evita relacionar películas solo por parecido superficial.

- Sincronización global: al abrir, importar o añadir listas y al obtener nuevas identificaciones, se reconcilian las marcas de todas las colecciones, incluidas las que ya estaban guardadas.

- Regla de estados: si cualquiera de las apariciones de una película está marcada Vista, la conciliación inicial mantiene Vista; cuando se desmarca manualmente, se desmarcan las demás apariciones reconocidas.

- Revisión de TMDB en Mis listas: se incorpora un indicador de progreso y la acción «Revisar TMDB»; continúan disponibles las correcciones manuales mediante el botón de lápiz.

- Identificaciones persistentes: se reutilizan la caché de plataformas y las correcciones manuales; los nuevos resultados se guardan en identidades_cine.json, incluido en el respaldo adicional en nube.

## Compatibilidad y cautelas

- Se mantiene la ventana nativa WebView2 confirmada desde la v7.3, las funciones de la v8.0 y el guardado principal en AppData.

- La identificación de títulos diferentes requiere conexión y clave TMDB; las coincidencias ya conocidas o de título y año pueden funcionar sin conexión. Las identificaciones ambiguas no se vinculan automáticamente.

02 / REGISTRO DE CAMBIOS

# Evolución: versión 8.1.1

<table>
<colgroup>
<col style="width: 14%" />
<col style="width: 85%" />
</colgroup>
<thead>
<tr class="header">
<th><strong>v8.1.1</strong></th>
<th><p><strong>Firma visual y procedencia del proyecto · Versión pública actual</strong></p>
<p>La aplicación incorpora el crédito visible «Hecho por Cristian Barallobre Rúa» y referencias de autoría en el código y en la documentación, sin cambiar la función principal de las colecciones.</p></th>
</tr>
</thead>
<tbody>
</tbody>
</table>

## Cambios incorporados

- Firma en la pantalla principal: crédito «Hecho por Cristian Barallobre Rúa» en la cabecera, integrado en la paleta oscura y dorada sin añadir desplazamiento general a la ventana.

- Marca de procedencia en el software: archivo autoria.go con nombre, copyright y repositorio canónico; metadatos de autor y atributos de procedencia en plantilla.html.

- Identificación adicional: el servidor local incluye cabeceras informativas con el creador y la dirección del proyecto. Estas marcas documentan la procedencia, pero no impiden modificar una copia de código abierto al público.

- Documentación: AUTORIA.md, README.md y CHANGELOG.md actualizados con el crédito del creador, la versión 8.1.1 y el enlace al repositorio.

- Distribución: el código se ha actualizado en GitHub; GitHub Actions ha ejecutado las pruebas Go, compilado el programa para Windows de 64 bits y publicado la Release v8.1.1 con un ZIP descargable.

## Publicación, pruebas y cautelas

- Repositorio oficial: https://github.com/cristianbarallobre/selector-cine-365

- Release oficial: https://github.com/cristianbarallobre/selector-cine-365/releases/tag/v8.1.1

- Confirmado en GitHub: ejecución satisfactoria de go test ./..., compilación cruzada para Windows x64 y creación del archivo ZIP de distribución. No debe confundirse con una prueba manual del ejecutable en Windows.

- Los datos existentes del usuario continúan almacenándose en AppData; la v8.1.1 no introduce una migración del formato de colecciones. Aun así, se recomienda mantener una copia de seguridad antes de actualizar.

- La firma de autor es identificativa, no constituye firma criptográfica ni certificación de Windows. El ejecutable publicado no está firmado digitalmente.

- Antes de abrir la v8.1 por primera vez, exportar una copia de seguridad: la conciliación puede convertir en vistas películas que antes figuraban pendientes en otras listas.

**03 / GUÍA DE USO**

# Uso habitual de la aplicación

## Biblioteca

- Elegir una colección en el selector superior. Si existe una comparación configurada, puede elegirse la vista virtual de coincidencias.

- Buscar por título, filtrar Vistas/Pendientes o limitar los resultados a una o varias plataformas.

- Navegar por páginas numéricas; usar los saltos de diez páginas o introducir un número de página.

- Sortear entre las películas pendientes de la selección actual o limitar el sorteo a un proveedor de streaming.

- En las colecciones reales, marcar o desmarcar una película; la v8.1 sincroniza el estado en las demás colecciones identificadas como la misma película. En la vista de coincidencias no se pueden editar títulos.

## Mis listas

- Crear, abrir, renombrar o eliminar una colección.

- Importar o pegar un TXT para reemplazar el contenido de la lista activa; «Añadir lista de películas» combina títulos sin sustituirla.

- «Eliminar película» permite buscar y borrar una entrada concreta con confirmación.

- «Comparar listas» configura las colecciones de la intersección virtual; los IDs de TMDB permiten reconocer títulos alternativos de una misma película.

- Exportar archivos o crear y restaurar copias de seguridad; configurar la carpeta de copia en nube y consultar «Revisar TMDB» para completar identificaciones pendientes.

## Estadísticas

El panel ofrece indicadores de avance y análisis por grupos temporales y categorías cinematográficas. Las cifras dependen de los metadatos disponibles y de las colecciones seleccionadas; no sustituyen una base externa de valoración o disponibilidad.

<table>
<colgroup>
<col style="width: 100%" />
</colgroup>
<thead>
<tr class="header">
<th><p><strong>Definición del estado de visionado</strong></p>
<p>Cada película se representa en su colección mediante una casilla: ☑ significa vista; ☐ significa pendiente. En la v8.1, los estados se reconcilian entre todas las apariciones de una misma película reconocidas por TMDB o, como respaldo, por título normalizado y año.</p></th>
</tr>
</thead>
<tbody>
</tbody>
</table>

**04 / DATOS Y CONTINUIDAD**

# Almacenamiento, copias y recuperación

La aplicación mantiene los datos en el equipo. La carpeta de nube es un respaldo adicional y no sustituye el guardado principal.

| **Archivo o ubicación**                     | **Contenido / función**                                                                                                           |
|---------------------------------------------|-----------------------------------------------------------------------------------------------------------------------------------|
| **%APPDATA%\AtlasDelCine\colecciones.json** | Fuente principal de colecciones, títulos, marcas ☑/☐, lista activa y selección de listas para comparación (comparison_ids).       |
| **mis_365_peliculas.txt**                   | Espejo en TXT de la colección inicial «Atlas del cine · 365», por compatibilidad con versiones antiguas.                          |
| **identificaciones_tmdb.json**              | Correcciones manuales de identificación de películas en TMDB.                                                                     |
| **identidades_cine.json**                   | Caché persistente de IDs TMDB utilizada para identificar títulos equivalentes y sincronizar listas; incluida en la copia de nube. |
| **plataformas_ajustes.json**                | Configuración local de TMDB, incluida la clave personal; no se copia al respaldo de nube.                                         |
| **plataformas_cache.json**                  | Datos consultados de plataformas / metadatos, con caducidad y actualización.                                                      |
| **copia_nube_ajustes.json**                 | Carpeta sincronizada elegida para el respaldo adicional.                                                                          |
| **copias_automaticas**                      | Instantáneas locales del JSON de colecciones; se conserva una versión previa por día con cambios.                                 |

## Esquema de respaldo en nube

Carpeta elegida → Selector_Cine_365_Copias → actual / historico / AAAA-MM-DD. En «actual» se guardan colecciones.json, mis_365_peliculas.txt, identificaciones_tmdb.json e identidades_cine.json cuando existen; «historico» conserva una instantánea diaria.

## Recuperar listas en otro ordenador

- Instalar la aplicación en el nuevo equipo y configurar el cliente de nube si se utiliza.

- Desde «Mis listas», elegir «Restaurar copia» y seleccionar colecciones.json de «actual» o del histórico.

- Comprobar la lista activa, el número de películas y las marcas; revisar también la recuperación de identidades_cine.json si se quieren mantener las equivalencias TMDB conocidas. Volver a configurar la clave TMDB cuando haga falta.

<table>
<colgroup>
<col style="width: 100%" />
</colgroup>
<thead>
<tr class="header">
<th><p><strong>Precaución ante conflictos</strong></p>
<p>La hora de «Copia en nube actualizada» refleja que la aplicación escribió en la carpeta sincronizada. La subida definitiva a OneDrive, Drive o Dropbox depende de su cliente. Comprueba que los archivos muestran el indicador de sincronización completada.</p></th>
</tr>
</thead>
<tbody>
</tbody>
</table>

**05 / OPERACIÓN Y MANTENIMIENTO**

# Requisitos, limitaciones y control de calidad

## Requisitos y dependencias

- Windows 10 u 11 de 64 bits; Microsoft Edge WebView2 Runtime para la ventana nativa.

- Go y recursos web integrados en el ejecutable. La app sirve la interfaz desde 127.0.0.1 en un puerto local dinámico.

- Conexión a internet y credencial de TMDB para consultar carteles, puntuaciones y plataformas de streaming.

- El ejecutable no está firmado digitalmente y Windows puede mostrar avisos de seguridad.

## Aspectos conocidos

- Disponibilidad de streaming: cambia por fecha y territorio; TMDB/JustWatch pueden contener retrasos o errores.

- Coincidencias entre listas: se prioriza el ID de TMDB; sin identificación se recurre a título normalizado y año. Una película sin ID puede quedar sin reconocer y un remake con otro año no debe mezclarse automáticamente.

- Identificación TMDB: la revisión automática necesita internet y clave configurada; si no hay coincidencia segura, puede fijarse un ID manualmente. Revisar posibles asociaciones erróneas antes de confirmar cambios masivos.

- Compatibilidad entre versiones: retroceder desde la 8.1 o la 8.1.1 puede dejar de aprovechar identidades_cine.json y la conciliación más reciente; exportar una copia de seguridad antes de cambiar de versión.

- Android: existe un proyecto de adaptación para Android Studio, pero no una APK compilada y probada incluida en esta línea de versiones de Windows.

## Verificación de versiones

| **Estado**                    | **Evidencia / significado**                                                                                                                                                                        |
|-------------------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| **Confirmado por el usuario** | v7.3: carga de WebView2, cierre normal e icono correcto al anclar a la barra de tareas.                                                                                                            |
| **Validado en preparación**   | v8.1: pruebas de Go, compilación y escenarios simulados de identidad TMDB. v8.1.1: pruebas Go y compilación Windows finalizadas con éxito en GitHub Actions; publicación de la Release verificada. |
| **Pendiente**                 | Prueba final de la v8.1.1 en Windows con listas y credencial TMDB reales; verificar el crédito en pantalla, el inicio WebView2, las coincidencias, las marcas y las copias de nube.                |

Repositorio del proyecto y versiones públicas: https://github.com/cristianbarallobre/selector-cine-365

## Convenciones para las próximas entregas

- Mostrar la versión en la interfaz y en el ejecutable.

- Registrar novedades, cambios, compatibilidad, migraciones y pruebas de cada versión.

- Publicar ZIP del programa y código, más instrucciones de instalación y recuperación.

**06 / NOTAS EDITORIALES**

# Alcance del historial y fuentes

Este historial reúne el desarrollo realizado entre el 17 y el 19 de noviembre de 2026. Las primeras subversiones no disponen de fechas de publicación formales y sistemáticas; se presentan en su orden de entrega. Desde la v8.1, los cambios y versiones publicadas se pueden contrastar en el repositorio oficial de GitHub.

## Fuentes empleadas

- Paquetes de código fuente de las versiones disponibles, especialmente los archivos LEEME, el código Go y la interfaz plantilla.html.

- Registros de arranque de WebView2 compartidos durante las correcciones de las versiones 7.1 a 7.3.

- Descripción de requisitos funcionales y comprobaciones comunicadas durante el desarrollo de Selector de Cine 365.

- Pruebas de implementación registradas en el material de las versiones 7.4 a 8.1.1; desde la 8.1.1, también se documentan la compilación Windows y las pruebas de Go ejecutadas mediante GitHub Actions.

## Créditos y servicios externos

La aplicación utiliza la API de The Movie Database (TMDB) para información cinematográfica y proveedores por territorio; la información de disponibilidad puede proceder de datos de JustWatch ofrecidos por TMDB. TMDB y JustWatch son servicios externos, sujetos a sus propias condiciones y a cambios de catálogo.

**TMDB:** https://www.themoviedb.org/

**JustWatch España:** https://www.justwatch.com/es

## Recomendación de publicación

<table>
<colgroup>
<col style="width: 100%" />
</colgroup>
<thead>
<tr class="header">
<th><p><strong>Mantener una fuente única de verdad</strong></p>
<p>El código, el changelog y la documentación publicada en docs/ deben corresponderse con la versión más reciente. Antes de actualizar, exportar colecciones.json y verificar la copia de identidades_cine.json.</p></th>
</tr>
</thead>
<tbody>
</tbody>
</table>

**FIN DEL DOCUMENTO · SELECTOR DE CINE 365 · v8.1.1**