# Manual de Selector Cine 365

## Edición estable v9.0.0 — Windows/macOS Beta 14 y Android Beta 13.13

**Actualizado:** 22 de noviembre de 2026 · **Estado:** versión estable de distribución, con los binarios Beta 14/13.13 que ya fueron probados. Se conservan sus números de versión internos.

### 1. Qué incluye cada descarga

- **Windows:** EXE portátil de 64 bits para Windows 10/11 con WebView2. Cierra la app antes de reemplazar el EXE. Los datos locales están en `%APPDATA%\Selector_Cine_365`; sustituir el EXE no debería borrar listas.
- **Mac Apple Silicon / Intel:** ZIP con aplicación macOS Beta 14. Escoge arquitectura correcta; estas compilaciones de prueba no tienen notarización de Apple. Requieren macOS 12 o superior.
- **Android:** Beta 13.13 (versionCode `9027`), identificador **`es.cine365.app.limpia`**, Android 8 o superior. Actualiza únicamente la línea **Beta 13.12 Limpia** que usa la misma firma. La antigua app `es.cine365.app` es distinta; Android no permite actualizarla con esta APK.

Los ZIP/EXE/APK y sus sumas SHA-256 se encuentran juntos en [GitHub Releases](https://github.com/cristianbarallobre/selector-cine-365/releases/tag/v9.0.0). Descargar desde esa página evita depender de enlaces temporales de ChatGPT.

### 2. Primera apertura y biblioteca

Puedes usar Selector Cine 365 **sin TMDB ni OneDrive**. Elige la lista en Biblioteca, busca una película y marca su estado **solo con la casilla de visto**. El resto de la ficha no debe alterar la marca.

El selector de búsqueda distingue **Lista actual** y **Todas las listas**. El segundo busca globalmente por título; no significa fusionar físicamente las listas. Las listas tienen nombre, ID y texto propios. **Cargar más películas** en Android muestra más títulos de la lista filtrada; si muestra títulos ajenos, podría tratarse de una lista mezclada antes de este arreglo.

En **Listas / Mis listas**, crea, renombra y elimina colecciones, importa TXT como una lista nueva, exporta TXT con el nombre de la lista y exporta/recupera una copia JSON de toda la biblioteca. Si eliminas una lista, la sincronización conserva un registro de borrado. Las vistas de comparación son calculadas y no reemplazan las colecciones originales.

### 3. Guardado, copia JSON y recuperación

Todas las modificaciones se guardan primero **en el dispositivo**. No hace falta Microsoft. Antes de cambiar de versión o de intentar recuperar una nube problemática:

1. En Windows abre **Mis listas → Copias de seguridad / Exportar copia JSON**; conserva el JSON en una carpeta aparte.
2. En Android abre **Listas → Copia de seguridad**, exporta el JSON y guarda una copia en Descargas o fuera de la aplicación.
3. Comprueba las listas por **nombre y número de películas**. Una copia con una lista inexplicablemente aumentada puede estar contaminada.
4. No borres `colecciones.json`, los archivos de `sync_recuperacion` ni `sincronizacion/dispositivos/*.json` en OneDrive antes de decidir cuál es la copia buena.

**Importante:** instalar una APK sobre el mismo identificador y firma preserva normalmente los datos; antes de cualquier cambio realiza una copia JSON. **Desinstalar** elimina el almacenamiento privado del paquete. Android no restaura automáticamente desde copias del sistema en esta beta; OneDrive sigue siendo externo.

### 4. Conectar Microsoft y OneDrive sin sobrescribir datos

**Windows Beta 14:** entra en **Mis listas → Sincronización y copias de seguridad → Conectar Microsoft**. Autoriza en el navegador una vez. La autorización de renovación se conserva cifrada mediante **Windows DPAPI** en el perfil del usuario, incluso si usas un EXE nuevo. Conectar no importa listas inmediatamente. Selecciona las opciones de revisión cuando corresponda.

**Android Beta 13.13:** abre **Nube → Conectar cuenta Microsoft** y utiliza **la misma cuenta** que Windows. Android localiza o crea la carpeta `Selector_Cine_365`. Pulsa **↻ Revisar, enviar y recibir cambios**. Examina la vista previa. La versión actual muestra los **nombres y números de películas por lista**, además de las vistas y conflictos. Pulsa **Confirmar y sincronizar** solo si la combinación tiene sentido.

**Mac Beta 14:** comparte los controles de escritorio pero la sesión Microsoft sigue siendo temporal, no está guardada en Keychain; puede pedir autorización cada vez que se reinicia.

La nube utiliza una carpeta `Selector_Cine_365` dentro de OneDrive, con snapshots independientes en `Selector_Cine_365_Copias/sincronizacion/dispositivos/` (o la ubicación de compatibilidad del dispositivo previamente configurado). El código de revisión causal de las marcas vistas intenta resolver cambios entre dispositivos; conflictos simultáneos necesitan decisión. No hay transferencia automática de datos al simplemente conectar.

**Android funciona con sincronización manual:** las ediciones locales no se publican por cada cambio ni por abrir la aplicación. En escritorio puede existir sincronización automática tras vincular; cuando se investigue un problema de listas, selecciona **Solo cuando pulse sincronizar** en Windows y no aceptes nuevas vistas previas.

### 5. Protección contra mezcla de listas

Betas anteriores unían películas cuando dos listas de instalaciones diferentes compartían el ID `atlas`, incluso si eran listas distintas. Las versiones recomendadas detectan colisiones entre listas, cantidades desproporcionadas o snapshots contaminados y **detienen la sincronización sin publicar cambios**. En la instalación inicial, Android sustituye su biblioteca de ejemplo por la copia buena de OneDrive **solo si el ejemplo no fue modificado**.

**No reparan listas que ya fueron fusionadas.** Si ves películas de distintas listas juntas: cancela el envío, mantén intactos los JSON originales, exporta la biblioteca buena de Windows y recupera selectivamente después de revisar títulos y cantidades. No intentes arreglarlo borrando snapshots a ciegas.

### 6. TMDB y carátulas

TMDB ofrece carátulas, notas y disponibilidad de plataformas; requiere una clave configurada individualmente en cada dispositivo desde **Conectar TMDB**. La clave **no se sincroniza**. La caché de carátulas y proveedores tampoco se sincroniza; una película ya identificada en Windows puede aparecer inicialmente **Por consultar** en Android.

Los IDs que un usuario corrigió manualmente pueden transmitirse como `tmdb_overrides` con la biblioteca. El móvil realiza solicitudes limitadas para completar fichas visibles (hasta 12 por tanda, 48 por sesión, dos simultáneas) cuando TMDB está conectado. Para más películas, utiliza **Nube → Completar carátulas y plataformas pendientes**. Una ficha incorrecta se corrige con el icono de lápiz.

### 7. Sorteo y estadísticas

Puedes sortear entre las películas pendientes y, si hay datos de TMDB, restringir por plataforma disponible en España. El sorteo no altera las listas por sí solo. Las estadísticas de totales, vistas, décadas u otros apartados dependen del estado de la biblioteca y de las fichas TMDB que se hayan podido consultar.

### 8. Instalación y actualización de Android

**Paquete recomendado:** `es.cine365.app.limpia` · **Compilación Android incluida:** 13.13 · **versionCode:** 9027 · **certificado:** misma clave original que Beta 13.12 Limpia.

Para actualizar una Beta 13.12 Limpia ya instalada, descarga la APK de esta misma línea y **pulsa Instalar sin desinstalar**. Debe permanecer el mismo identificador de paquete, aumentar versionCode y mantenerse el certificado. Si Android muestra **Conflicto con un paquete**, no borres datos: verifica la identidad de la app instalada y que no estés intentando instalar sobre `es.cine365.app` o una APK de depuración con otra firma.

### 9. Dónde están los datos

| Plataforma | Ruta de biblioteca |
| --- | --- |
| Windows | `%APPDATA%\Selector_Cine_365\colecciones.json` |
| Mac | `~/Library/Application Support/Selector_Cine_365/colecciones.json` |
| Android | Directorio privado de `es.cine365.app.limpia`, carpeta `Selector_Cine_365` |
| OneDrive | Carpeta `Selector_Cine_365` de la cuenta Microsoft conectada |

Al cambiar de dispositivo, no copies contraseñas ni tokens cifrados. Exporta listas mediante JSON y utiliza el proceso explícito de revisión antes de unir bibliotecas.

### 10. Limitaciones pendientes

Esta es la **versión de distribución estable v9.0.0**, basada en las compilaciones que el usuario ya ha comprobado en sus dispositivos; eso no implica ausencia garantizada de errores. Falta validar todas las combinaciones reales de Microsoft y hardware; en Mac no hay persistencia segura de sesión con Keychain; los datos de TMDB se consultan individualmente; las fusiones históricas necesitan recuperar desde un JSON sano. Conserva siempre una copia externa antes de sincronizar.

[Notas completas de esta publicación](../RELEASE_NOTES_v9_BETA14_RECOMENDADA.md) · [Guía para recuperación segura](GUIA_SINCRONIZACION_SEGURA_BETA14_ANDROID13_13.md).

### 11. Instrucciones prácticas y errores habituales

**Películas de todas las listas dentro de una sola:** primero asegúrate de que no estás buscando en «Todas las listas» ni usando una vista de comparación. Si en la lista activa aparecen títulos ajenos incluso sin búsqueda global, **no sincronices**. Comprueba su recuento en la biblioteca de Windows y exporta la copia correcta de `colecciones.json`; las barreras de la versión actual previenen nuevas uniones sospechosas, pero no restauran los títulos si ya se mezclaron antes.

**Póster visible en Windows pero no en Android:** OneDrive comparte las listas, los vistos y ciertos IDs de TMDB, **no** la imagen ni la caché de proveedores descargada. Configura TMDB también en Android. Al navegar, el móvil consulta progresivamente una parte de las fichas visibles; para completar una lista utiliza «Nube → Completar carátulas y plataformas pendientes». No necesitas volver a iniciar sesión en Microsoft para una carátula.

**El APK indica «Conflicto con un paquete»:** identifica el paquete instalado. `es.cine365.app.limpia` puede actualizar la Beta 13.12 Limpia únicamente si conserva **el mismo certificado**; `es.cine365.app` y APK debug de prueba son aplicaciones o firmas diferentes. **No desinstales sin copia JSON**.

**Windows vuelve a mostrar la autorización de Microsoft:** el token de renovación se guarda cifrado con DPAPI bajo el usuario Windows. Se necesita iniciar sesión al menos una vez después de pasar de versiones que no persistían esa sesión; si Microsoft ha revocado los permisos o ha cambiado la cuenta, puede ser necesario repetir la autenticación. No publiques el archivo cifrado ni tus tokens en una incidencia.

**Mac pide conectar con Microsoft al reiniciar:** esta compilación no almacena todavía las credenciales en Keychain. Es una limitación conocida, no un fallo de la biblioteca local.

**No aparece una lista después de sincronizar:** revisa la vista previa, las eliminaciones sincronizadas y las copias de recuperación. Nunca elimines datos de otros dispositivos en OneDrive sin verificar antes sus snapshots.

**La descarga del chat no se inicia:** usa directamente la [página oficial de GitHub Releases](https://github.com/cristianbarallobre/selector-cine-365/releases/tag/v9.0.0), que mantiene enlaces permanentes de los instaladores.

### 12. Privacidad, código fuente y mantenimiento

La biblioteca es local por defecto. TMDB requiere tu clave propia para datos enriquecidos; OneDrive requiere autorización delegada a tu cuenta Microsoft para compartir listas. El programa no necesita contratar servidores para guardar tu biblioteca. Las claves privadas de firma, las contraseñas y los tokens **no forman parte del código fuente público**.

La etiqueta `v9.0.0` fija el código fuente de referencia de esta edición. Para futuras APK de Android, conservar **`es.cine365.app.limpia` y el certificado permanente**, aumentando el código de versión. No confundir la firma de la distribución con los paquetes debug temporales de CI.

Descarga el [código fuente](https://github.com/cristianbarallobre/selector-cine-365/tree/v9.0.0) o consulta [la política de mantenimiento](POLITICA_DE_RELEASES.md).

**Autor:** Cristian Barallobre Rúa.
