> DOCUMENTO ARCHIVADO: esta edición corresponde a Windows/Mac Beta 13.1 y Android Beta 13; la vigente es [Manual completo v9.0.0 estable](../MANUAL_COMPLETO_SELECTOR_CINE_365.md). Se conserva íntegra para recuperar el contexto técnico de las versiones anteriores.

> **MANUAL ACTUALIZADO:** para las versiones recomendadas del 21/11/2026 (**Windows/Mac Beta 14 y Android Beta 13.13**) usa [Manual Beta 14 / Android 13.13](MANUAL_BETA14_ANDROID13_13.md) y [Guía de sincronización segura](GUIA_SINCRONIZACION_SEGURA_BETA14_ANDROID13_13.md). La documentación que continúa a continuación es un **archivo de Beta 13.1/13**, conservado sin borrar su historia. No apliques las instrucciones antiguas de instalación o sincronización a las versiones nuevas.

# Selector Cine 365

## Manual completo de uso: Windows, macOS y Android

**Edición:** v9 Beta 13.1 Windows/Mac y v9 Beta 13 Android (guía de pruebas, no versión estable)  
**Fecha:** 20 de noviembre de 2026  
**Autor del programa:** Cristian Barallobre Rúa  
**Objeto de este manual:** biblioteca, películas vistas, listas, búsqueda global, sorteo, estadísticas, TMDB, copias y sincronización opcional con OneDrive.

> Este manual describe la release coordinada v9 Beta 13.1 para Windows y Mac y la APK Android Beta 13. La interfaz de escritorio se comparte entre Windows (WebView2) y Mac (WKWebView); Android tiene navegación propia. No se ha validado toda la funcionalidad en hardware Mac real. Antes de probar nuevas betas, exporta una copia JSON y conserva una versión anterior. La nube es opcional: la biblioteca también funciona sin Internet ni clave TMDB.

## Android Beta 13.3 — guardar localmente y sincronizar cuando quieras

En Android, **cambiar de lista deja de enviar un snapshot a OneDrive**: cambiar la lista activa se guarda de forma silenciosa solo en el móvil. Marcar películas o editar listas se guarda automáticamente en el teléfono, sin necesidad de pulsar un botón. Al abrir la aplicación, la nube se consulta una vez para preparar una **vista previa de solo lectura**; no se importan ni publican cambios por abrir.

El botón compacto de **dos flechas circulares**, situado en la cabecera junto al estado de guardado, abre Nube y compara el estado local y remoto. Revisa el resumen, resuelve conflictos si aparecen y pulsa **Confirmar y sincronizar** para recibir y enviar los cambios. Sin confirmación no se escribe ni se importa nada. En segundo plano no hay comprobaciones periódicas ni al regresar desde otra app.

La versión de Android es **9017** (`9.0.0-beta.13.3`), con el mismo paquete `es.cine365.app` y el mismo certificado que las Betas anteriores desde Beta 11. Instala encima sin borrar datos, pero realiza una copia JSON antes. Windows/Mac siguen en **Beta 13.1** y conservan su sincronización propia; para recibir desde el PC en Android debes pulsar las dos flechas.

**Importante:** el modo manual evita tráfico innecesario, pero los cambios locales no estarán copiados en OneDrive hasta que los sincronices expresamente. Para no perder nada al cambiar de teléfono, pulsa y confirma el icono de sincronización o exporta JSON.

## Corrección Beta 13.1 — TMDB al abrir Windows o Mac

Al reabrir Beta 13, podía aparecer permanentemente «Conexión local TMDB no disponible. Reintentando…», incluso con una clave válida guardada. El problema no estaba en tu clave ni en OneDrive: el servidor local tenía la ruta de caché `/api/platforms/cache` implementada, pero no conectada al enrutador principal, por lo que respondía HTTP 404. **Beta 13.1** corrige esa ruta y distingue el estado de credenciales del estado de la caché. Si la caché falla, la clave local no se borra ni se marca erróneamente como ausente; el botón Actualizar permite volver a consultarla.

El archivo de Windows es `Selector_Cine_365_Windows_v9_BETA13_1.zip` y los de Mac terminan en `v9_BETA13_1.zip`. **Android sigue en Beta 13 (9015)**; su backend es diferente y no necesita este ajuste. El programa puede usarse localmente sin TMDB ni nube.

## Correcciones de Beta 13 — listas, TMDB y sincronización

En Windows/Mac, la sección Mis listas mantiene cuatro filas por página para que la parte inferior no se mueva en páginas cortas. Al borrar una lista, se registra su ID en `deleted_lists` y se comparte con los otros dispositivos; una copia anterior no debe recuperarla. Las identificaciones manuales TMDB se guardan como IDs en `tmdb_overrides` y se reciben por OneDrive, pero la clave API no se comparte. Al abrir de nuevo, un fallo temporal de la interfaz no debe hacer que la clave configurada parezca perdida.

**Antes de probar:** exporta un JSON en ambos dispositivos y actualiza todos a Beta 13. Los borrados antiguos no llevaban marcador: si una lista reapareció antes de esta beta, elimínala de nuevo cuando ambos estén actualizados. Una copia de recuperación conserva el contenido para restaurar intencionadamente en una nueva lista.

## Correcciones de pantalla de Android Beta 12.1 y 12.2

En móviles con notch o cámara frontal recortada, **Beta 12.1** adaptó el contenido a las zonas seguras de Android mediante `WindowInsets`, `systemBars` y `displayCutout`. Algunos Redmi dejaron una franja vacía bajo el menú porque el CSS sumaba nuevamente el margen inferior. **Beta 12.2 (versionCode 9014)** elimina esa doble reserva para que las pestañas queden encima de la navegación del sistema, conservando la protección del notch. Se puede instalar sobre Beta 11, 12 y 12.1 con la misma firma, sin desinstalar y preservando los datos locales. Desde Beta 9 se necesita exportar e importar JSON. Windows/Mac no requieren cambios. Falta validación visual definitiva en el móvil real.

## 1. Primeros pasos y conceptos básicos

Selector Cine 365 permite guardar colecciones de películas, identificar cuáles has visto, filtrar las pendientes, elegir una película al azar y consultar estadísticas. Cada lista tiene sus propios títulos. Las marcas de visto pueden reconciliarse entre listas cuando se reconoce la misma película.

En Windows y Mac las secciones principales son **Biblioteca**, **Estadísticas**, **Mis listas** e **Instrucciones**. En Android se utilizan las pestañas **Biblioteca**, **Sorteo**, **Estadísticas**, **Listas**, **Nube** y **Ayuda**. La navegación puede variar por el tamaño de la pantalla.

**No es obligatorio configurar todo al empezar.** Puedes crear listas y marcar películas como vistas sin conectar OneDrive y sin añadir una clave de TMDB. Si más adelante activas la nube, revisa primero los cambios antes de unir dispositivos.

## 2. Biblioteca y películas vistas

1. Elige arriba la lista que quieras consultar.
2. Usa **Buscar película** con **Lista actual** o **Todas las listas**. En la búsqueda global solo se muestran fichas cuando se escribe un título; se agrupan entradas equivalentes por TMDB o título/año, sin duplicados. Los filtros **Todas**, **Pendientes** y **Vistas** limitan lo que aparece, sin borrar películas.
3. Para marcar como vista una película, pulsa **únicamente la casilla cuadrada con el tick**, situada en su ficha. Pulsa de nuevo esa casilla para volver a pendiente.
4. Tocar el cartel, el título o el resto de la ficha **no debe cambiar** su estado. Los otros botones pueden abrir información o permitir corregir la identificación.
5. En Windows y Mac puedes pasar entre páginas con números, flechas y, cuando se muestre, el campo para ir a una página concreta. En Android la biblioteca puede mostrar **Cargar más** para ampliar el listado sin perder la posición.
6. Cada ficha muestra **año y valoración TMDB** cuando existen datos (no es IMDb), además de **Listas: 1, 2, …**. Pulsa un número para ir a esa lista o el icono de edición para añadir o quitar esa película de otras listas. La numeración se actualiza cuando cambia el orden de las listas; sus IDs internos permanecen iguales.

Si una marca cambia también en otra lista, comprueba que se trate realmente de la misma película, especialmente cuando existan títulos iguales y años distintos. Antes de importar muchas películas, guarda una copia JSON.

## 3. Añadir, importar y administrar listas

En **Mis listas / Listas** puedes crear una colección, abrirla, cambiarle el nombre, eliminarla y, según la pantalla, comparar varias. Cada lista mantiene su título y contenido; cambiar de lista no implica borrarla.

Para añadir películas puedes escribirlas individualmente o importar una relación de títulos desde un TXT, cuando esté disponible. El formato de texto más sencillo es:

```
☐ Ciudadano Kane (1941)
☑ Casablanca (1942)
☐ 2001: Una odisea del espacio (1968)
```

La casilla vacía significa pendiente; la marcada, vista. Se admite numeración en muchas listas. Revisa títulos y años antes de confirmar, especialmente si la app detecta duplicados. Si dos listas incluyen una película identificada como la misma, la marca puede aplicarse a ambas.

En Windows y Mac, **Mis listas** incorpora números de página para moverte por todas las colecciones, además de **Anterior**, **Siguiente** y el salto directo por número. Si comparas listas, las coincidencias se muestran como consulta en la biblioteca; esa vista puede ser de solo lectura.

**Antes de eliminar una lista**, exporta una copia JSON. La sincronización experimental puede conservar datos de otros dispositivos cuando se realizan eliminaciones o cambios de nombre simultáneos.

## 4. Sorteo de películas

En **Sorteo** (o en el panel derecho de la biblioteca Windows) puedes elegir al azar entre las películas pendientes. Es posible limitar el sorteo a una plataforma de streaming si sus datos están disponibles.

1. Elige la colección y, si procede, las plataformas.
2. Pulsa **Elegir una película al azar**.
3. Consulta el título y la información ofrecida.
4. Si finalmente la ves, utiliza **Marcar como vista**.

El sorteo no elimina la película de tu lista. Si filtras por plataforma, recuerda que la disponibilidad puede cambiar y depende de las consultas realizadas a TMDB; no equivale a una garantía de reproducción gratuita.

## 5. Estadísticas

La sección **Estadísticas** resume títulos totales, vistas y pendientes. Según los datos de tus listas y de TMDB, también puede mostrar décadas, años, géneros, países, directores, progreso y otros apartados.

En **Windows y Mac**, los paneles de estadísticas con paginación mantienen los controles en la parte inferior aunque la última página tenga menos filas. En **Android** se recorren mediante desplazamiento vertical; no hace falta ir abriendo paneles con flechas. Una estadística puede faltar cuando la película no tiene suficientes datos identificados. No es necesario volver a consultar TMDB cada vez que se abre la aplicación.

## 6. Conectar TMDB: carteles, valoraciones y plataformas

TMDB es un servicio externo de información sobre películas. **Es opcional**: no almacena tus listas de Selector Cine 365 ni sustituye a OneDrive. En Windows el botón de configuración se llama **Conectar TMDB**; en Android está disponible desde los ajustes o la ayuda.

1. Entra en <https://www.themoviedb.org/> y crea una cuenta si no la tienes.
2. En **Ajustes → API**, sigue las condiciones de TMDB para obtener una credencial.
3. Copia la clave admitida por la aplicación y pégala en **Conectar TMDB / Configurar TMDB**.
4. Guarda la clave y realiza la consulta inicial cuando lo necesites.
5. Usa **Actualizar** o **Revisar TMDB** cuando quieras solicitar datos de nuevo o corregir identificaciones. Las consultas pueden gastar datos móviles.

La aplicación mantiene una **caché local persistente** de información ya consultada. La finalidad es no recargar en cada apertura los carteles y plataformas que ya conoce. Nuevas películas o actualizaciones solicitadas pueden requerir Internet. La disponibilidad en España se obtiene mediante datos de proveedores proporcionados por TMDB y puede variar.

**Seguridad:** no compartas tu clave de TMDB. Si instalas una beta Android como aplicación independiente, puede ser necesario configurar la clave otra vez; importar una copia de listas no siempre traslada la caché o los ajustes.

## 7. Guardado local y copias de seguridad

Las listas se guardan localmente en cada dispositivo. El guardado no depende de tener una carpeta de nube configurada. La opción de copia JSON permite respaldar y transferir tu biblioteca antes de cambiar de versión.

**Consejo de seguridad:** exporta una copia antes de importar muchas películas, eliminar listas, conectar un dispositivo por primera vez o instalar una beta nueva. Conserva el archivo en un lugar distinto de la carpeta de sincronización. Comprueba que puedes leerlo antes de borrar una instalación anterior.

**Windows:** las betas comparten el mismo almacén local con versiones previas. No abras dos ejecutables simultáneamente.  
**macOS:** reemplaza la aplicación anterior en Aplicaciones; conserva las listas de `~/Library/Application Support/AtlasDelCine/`. Mac usa el mismo frontend que Windows y firma ad hoc no notarizada.  
**Android:** Beta 11 y 12 usan el paquete `es.cine365.app` y requieren la misma clave de firma, por lo que Beta 12 actualiza Beta 11. Beta 9 (`es.cine365.app.v9beta9`) necesita **una migración mediante exportar/importar JSON**, porque su firma e identificador son distintos. No desinstales Beta 9 antes de verificar el JSON fuera de la app.

Si desactivas la nube, las películas siguen guardadas localmente. Desconectar una cuenta de OneDrive no significa borrar sus archivos remotos.

## 8. Sincronización con OneDrive: conexión inicial

**Recomendación para una instalación nueva:** utiliza una carpeta de OneDrive llamada `Selector_Cine_365`. Ese nombre es solo una sugerencia; si ya sincronizas con otra carpeta, **no la renombres ni la muevas** sin realizar antes una copia JSON y volver a configurar los dos dispositivos.

### Paso 1. Prepara Windows o macOS

Instala el cliente de OneDrive en el equipo de escritorio e inicia sesión. Dentro de OneDrive crea o elige la carpeta `Selector_Cine_365`. En Selector Cine 365 Windows o Mac abre la configuración de nube, usa **Elegir carpeta**, selecciona la ruta local sincronizada y activa la copia. Espera a que OneDrive complete su sincronización.

### Paso 2. Conecta Android

En el móvil entra en **Nube → Conectar cuenta Microsoft**. Se abrirá la página oficial de Microsoft en el navegador. Inicia sesión, acepta los permisos de lectura y escritura de archivos y permite que Microsoft te devuelva a Selector Cine 365. No necesitas copiar códigos temporales ni introducir tu contraseña dentro de la app.

### Paso 3. Elige la misma carpeta

En Android actualiza o consulta la lista de carpetas de OneDrive. Selecciona `Selector_Cine_365` o crea esa carpeta si todavía no existe. Pulsa **Utilizar esta carpeta**. **La carpeta debe ser exactamente la misma que usa Windows o Mac**, aunque Android la vea mediante Microsoft Graph y Windows como una ruta local.

### Paso 4. Primera revisión manual

En **los dispositivos vinculados**, pulsa **Revisar sincronización**. Lee el resumen de listas, películas vistas y cambios detectados; resuelve cualquier conflicto y confirma. Esta revisión inicial es obligatoria para evitar que una biblioteca antigua sobrescriba otra.

### Paso 5. Uso habitual

Tras la primera conexión confirmada, los cambios compatibles se reciben automáticamente al abrir o reanudar la app y de forma periódica mientras se utiliza. No es instantáneo: OneDrive puede tardar en subir y descargar los archivos. Los conflictos se resuelven desde **Revisar sincronización**.

La carpeta compartida crea subcarpetas de copias y sincronización, entre ellas `Selector_Cine_365_Copias`, `actual` y `sincronizacion`, con archivos separados por dispositivo. **No borres ni modifiques a mano esos archivos.**

## 9. Qué sucede cuando ambos dispositivos cambian cosas

Ejemplo sin conflicto: marcas una película como vista en Windows. El cambio se guarda y se publica en OneDrive. Cuando el móvil consulte la nube, podrá recibir esa marca sin pulsar el botón de revisión. También funciona en sentido contrario.

Ejemplo con conflicto: mientras los dos dispositivos están sin recibir cambios, uno marca una película como vista y el otro como pendiente. Para evitar sobrescribir silenciosamente una decisión, la app solicita que abras **Revisar sincronización** y elijas qué estado conservar.

La sincronización de listas renombradas o eliminadas simultáneamente sigue siendo una función **experimental**. Si trabajas sin conexión durante mucho tiempo en ambos dispositivos, exporta copias JSON antes de unir los datos. La sincronización automática no funciona como servicio permanente con la app cerrada; las comprobaciones se realizan al abrir, al volver y mientras se usa.

## 10. La sección Nube de Android, por orden

1. **Conectar Microsoft:** autoriza el acceso a OneDrive en el navegador.
2. **Elegir carpeta:** selecciona o crea la carpeta `Selector_Cine_365`.
3. **Preparar Windows:** utiliza esa misma carpeta local del cliente OneDrive.
4. **Revisar y confirmar:** realiza la primera combinación o resuelve conflictos.
5. **Sincronización automática:** se comprobarán cambios compatibles durante el uso.

En operaciones manuales de Nube aparece un **aviso de progreso**. Hasta que termine, los controles de esa sección se mantienen bloqueados para impedir pulsaciones duplicadas. Esto no debe bloquear la biblioteca ni la navegación general durante una comprobación automática.

Si aparece un error de red, espera a recuperar la conexión y vuelve a intentarlo. No pulses repetidamente **Conectar** ni **Revisar** mientras se esté mostrando el progreso. Si OneDrive está tardando, la app puede avisar de la espera sin asumir que el proceso ha terminado.

## 11. Microsoft Entra: configuración del propietario de la app

Para conectar OneDrive en Android, el desarrollador registró **Selector Cine 365** en Microsoft Entra como aplicación pública compatible con **cuentas personales de Microsoft**. La configuración incluye un ID público de aplicación, permisos delegados `Files.ReadWrite` y `User.Read`, flujos de cliente público y una dirección de retorno para Android. La conexión utiliza OAuth con **PKCE**, sin secreto de cliente.

**Un usuario normal NO necesita crear una cuenta Azure ni registrar una aplicación nueva.** Solo tiene que iniciar sesión en su OneDrive personal desde la pantalla Microsoft de Android.

Para rehacer la configuración del desarrollador o resolver incidencias de autorización, consulta el documento específico del repositorio: **Guía de Azure, Microsoft Entra y OneDrive de Selector Cine 365**. No publiques contraseñas, secretos, claves de sesión ni tokens en GitHub.

## 12. Historial local, sincronizar ahora y recuperar copias

**Actualizar historial** no pide nuevos datos a OneDrive. Consulta el registro local de cambios realizados o recibidos por las versiones que implementan esa función, indica la hora de comprobación y cuántos movimientos muestra. Si está vacío, un mensaje sencillo explica que no existen eventos registrados: **no** implica que tus listas estén vacías ni que haya ocurrido un fallo. Las versiones antiguas no pueden aportar retrospectivamente cambios que nunca registraron.

**Sincronizar ahora** fuerza la comprobación de cambios de la nube cuando ya está configurada. **Revisar sincronización** permite inspeccionar los cambios y resolver conflictos antes de confirmar. Son funciones diferentes.

Las **copias de recuperación** permiten restaurar estados anteriores con confirmación. Antes de restaurar, exporta un JSON adicional: restaurar puede generar revisiones nuevas que se compartirán si la nube está activada.

## 13. Problemas habituales y soluciones

| Síntoma | Qué comprobar |
|---|---|
| Android no responde al tocar o deslizar | Comprueba qué versión beta estás usando. La v9 conserva la corrección que elimina el bloqueo general de interfaz durante la recepción automática. No toques controles de Nube mientras se muestra un progreso; si se repite, anota los pasos y el tiempo de bloqueo. |
| No aparece la carpeta de OneDrive en Android | Conecta primero tu cuenta Microsoft dentro de Nube. El selector local de Android no equivale a OneDrive. Pulsa Actualizar carpetas. |
| Windows y Android no ven los mismos datos | Verifica que ambos apuntan a la misma carpeta y que OneDrive ha terminado de copiarla. Comprueba que la primera sincronización se confirmó. |
| La app pide revisar un conflicto | Entra en Revisar sincronización, examina las películas indicadas y elige un estado. No borres copias para forzar la operación. |
| No aparecen carteles o plataformas | Comprueba la clave de TMDB, la conexión y la disponibilidad de esos títulos. Revisa TMDB manualmente; no borres las listas. |
| Android solicita iniciar sesión de nuevo | Vuelve a conectar Microsoft en la app. Los tokens de acceso pueden caducar o revocarse. |
| Se abre CMD al ejecutar Windows | Comprueba que arrancas el EXE Beta 12 con subsistema Windows GUI. Algunas betas antiguas abrían consola. |
| No se reconocen películas tras importar | Revisa el formato de texto, los títulos y los años. Importa desde una copia JSON válida cuando quieras recuperar toda la biblioteca. |

## 14. Privacidad, límites y buenas prácticas

Selector Cine 365 guarda tus colecciones localmente. Si activas OneDrive, publicará archivos de sincronización en la carpeta autorizada de tu cuenta Microsoft para compartirlos entre tus dispositivos. TMDB se utiliza para obtener información cinematográfica y plataformas, no para guardar tu biblioteca.

No compartas la carpeta de sincronización con terceros, no publiques tu biblioteca personal en repositorios y evita copiar a capturas públicas claves, tokens o datos de acceso. Mantén una copia JSON fuera de la nube si la biblioteca es importante para ti.

Las betas son versiones de prueba y pueden tener errores. La sincronización segura intenta preservar las marcas y avisar en conflictos, pero ninguna aplicación sustituye las copias de seguridad independientes. Si se presenta un fallo, conserva los archivos, anota el dispositivo y la versión, y consulta el historial del repositorio antes de actualizar.

## 15. Enlaces útiles

- Proyecto y versiones: <https://github.com/cristianbarallobre/selector-cine-365>
- Historial v9 Beta 12 en Markdown, Word y PDF: <https://github.com/cristianbarallobre/selector-cine-365/blob/main/docs/HISTORIAL_V9_BETA13_1.md>
- Manual completo (esta guía): <https://github.com/cristianbarallobre/selector-cine-365/blob/main/docs/MANUAL_COMPLETO_SELECTOR_CINE_365.pdf>
- Guía de Microsoft Azure, Entra y OneDrive: <https://github.com/cristianbarallobre/selector-cine-365/blob/main/docs/GUIA_AZURE_ENTRA_ONEDRIVE_SELECTOR_CINE_365.md>
- TMDB: <https://www.themoviedb.org/>
- Portal Microsoft Azure: <https://portal.azure.com/>
- Instrucciones Mac: <https://github.com/cristianbarallobre/selector-cine-365/blob/main/MACOS.md>
- Guía Android: <https://github.com/cristianbarallobre/selector-cine-365/blob/main/android/LEEME_ANDROID_v9_BETA.md>
