# Android Beta 13.13 — aislamiento de listas y recuperación segura

- Las versiones anteriores fusionaban automáticamente películas de listas con el mismo ID, especialmente `atlas` después de reinstalar Android. Una lista podía terminar con películas de otras listas dentro de su propio texto, aunque el selector mostrase una sola lista.
- Beta 13.13 detiene cualquier unión sospechosa, incluidas bibliotecas ya contaminadas o listas con el mismo ID pero nombres distintos. No envía ni importa datos cuando detecta ese problema.
- En instalación nueva y sin modificaciones locales, la lista de ejemplo de 365 películas **se sustituye por la biblioteca real de Windows** encontrada en OneDrive; no se anexa al final. Se priorizan snapshots de escritorio y se omiten únicamente ejemplos Android intactos sin cambios.
- En la vista previa se muestran los nombres y números de películas por lista. Cancela si algún número no cuadra; no pulses Confirmar ni Sincronizar.
- **Si las listas ya están mezcladas:** no se pueden separar automáticamente con garantías. Cambia Windows a sincronización manual y conserva `%APPDATA%\Selector_Cine_365\colecciones.json` y las copias `sync_recuperacion`. Guarda también una copia JSON desde Android antes de intentar restaurar. Solo restaura desde una biblioteca comprobada, por ejemplo la de Windows; no borres snapshots de OneDrive sin conservarlos.
- Esta versión mantiene el identificador de la variante limpia `es.cine365.app.limpia`, `versionCode 9027`, `9.0.0-beta.13.13` y firma el APK de prueba entregado con el mismo certificado original Beta 13.12 (`b8a17f39e4fa24d7c110d452161425c7524412c3a5a4d7459b1132ffd0825fb6`). La clave privada no se ha añadido a GitHub.
- Tests Android (`ListIdentityIsolationTest`) y Windows (`sync_list_identity_beta14_test.go`) protegen los casos detectados. La sincronización real de la instalación del usuario debe comprobarse antes de publicar una versión estable.

---
# Android Beta 13.12 — cargar carteles y plataformas pendientes en el móvil

- Al abrir la Biblioteca con TMDB ya conectado, Android completa **de forma limitada** las fichas visibles que aún no tiene guardadas. También lo hace tras confirmar una sincronización de OneDrive, seleccionar otra lista o cargar más películas.
- Solo consulta hasta 12 fichas por tanda, con dos peticiones concurrentes, y como máximo 48 por sesión; respeta la caché y el enfriamiento de errores. No recorre las 365 películas ni solicita nada a OneDrive por ello.
- Las listas y las marcas de vistas se sincronizan con OneDrive; las carátulas, valoraciones y plataformas son fichas TMDB locales. Para completar el resto de una lista sigue existiendo **Nube → Completar carátulas y plataformas pendientes**.
- Sin una clave TMDB configurada en el móvil no se hacen consultas; la clave personal sigue sin sincronizarse ni compartirse con Windows.
- Nueva prueba JavaScript de consultas limitadas y sin subida a la nube. Android `versionCode 9026`, `versionName 9.0.0-beta.13.12`.
- **Certificado de pruebas:** la APK debug de GitHub Actions está firmada con una clave temporal diferente de Beta 13.11. No es una actualización directa: instalarla sobre una beta firmada con otra clave provoca un error de Android. Utilizar el certificado original en las siguientes versiones sigue siendo necesario.

---
# Selector Cine 365 — Android v9 Beta 13.11 (instalación limpia de pruebas)

**Estado:** Beta 13.11 (versionCode `9025`) para probar Android recién instalado con Windows Beta 14 y OneDrive ya configurado. No es una release estable.

## Instalar desde cero, sin datos antiguos

1. Si desinstalaste `es.cine365.app`, los archivos privados del Android anterior normalmente desaparecieron. Una copia antigua de Android/Google o la transferencia entre teléfonos podía restaurar configuraciones sin avisar; desde esta Beta se **desactiva Auto Backup** y se excluyen explícitamente archivos privados, credenciales y datos de dispositivo tanto de las copias del sistema como de la transferencia de Android (reglas para Android 8–11 y Android 12+). Esto no borra copias anteriores ya existentes fuera del móvil.
2. **OneDrive es independiente del móvil:** la carpeta `Selector_Cine_365` y los archivos de Windows y de posibles Android anteriores siguen existiendo. Nunca borres esos archivos solo por haber desinstalado la app.
3. Abre la Beta 13.11 y confirma que puedes usarla en local. En **Ajustes y sincronización → 1. Conecta tu cuenta Microsoft**, utiliza **la misma cuenta de OneDrive** que Windows. Android localizará o creará la carpeta raíz correcta; conectar no importa ni sube películas.
4. En **2. Revisa las listas**, pulsa **↻ Revisar, enviar y recibir cambios**. Comprueba el número de listas y de películas vistas; puede haber copias de Windows y de un Android anterior. Si hay conflictos, decide explícitamente; si ves algo raro, **Cancelar** no cambia ni publica datos.
5. Pulsa **Confirmar y sincronizar** solo después de revisar. Se guarda una recuperación local antes de combinar y se envía **únicamente el JSON de este móvil**. Las copias de Windows permanecen intactas.
6. En **3. Sincroniza cuando quieras**, vuelve a usar ↻ para enviar o recibir cambios. Android **no sincroniza por intervalos ni automáticamente en segundo plano**; al abrir consulta la nube solo en modo lectura si ya está vinculado.

## Certificado de firma y distribución de prueba

- Las compilaciones GitHub Actions de esta rama generan una `app-debug.apk` **firmada con la clave de depuración temporal del runner** para instalaciones totalmente limpias, y un `app-release-unsigned.apk` que no sirve para instalar directamente.
- **Esta APK de depuración NO utiliza el certificado original de las betas publicadas.** Solo úsala si no tienes instalada una APK con otra firma; para actualizaciones posteriores sin reinstalar habrá que distribuir una versión firmada con **el certificado original de Selector Cine 365**, cuya clave privada no se guarda en GitHub.
- La aplicación sigue usando `applicationId es.cine365.app`; Android 8.0+ (API 26). Los cambios de interfaz y reglas de restauración no alteran los datos compartidos de Windows.

## Comprobaciones

- Prueba JVM de primera vinculación: se fusiona una biblioteca Windows ya vista con el estado inicial del móvil, se respeta la revisión causal de Windows frente a una copia Android antigua y las decisiones realmente concurrentes se dejan para revisión.
- Pruebas de prohibición de restauración automática, revisión antes de subir y compilación debug/release, más verificación de la firma y versión de la APK.
- Falta validar la instalación y la sincronización real en Redmi. Conserva una copia JSON de Windows antes de aceptar cualquier importación.

---
# Selector Cine 365 — Android v9 Beta 13.5

**Autor:** Cristian Barallobre Rúa. **Estado:** versión beta experimental para Android 8.0+ (API 26). Conserva funcionamiento local sin cuenta Microsoft y sin clave TMDB.

## Android Beta 13.5 — conectar OneDrive sin elegir carpeta

- **Conexión directa Microsoft:** si ya existe una carpeta de OneDrive enlazada desde Android, la reutiliza. Si no existe, busca o crea `Selector_Cine_365` directamente en OneDrive y la conecta sin necesidad de elegirla de una lista.
- No se importa ni publica ninguna biblioteca al conectar: sigue siendo obligatorio pulsar el botón **↻**, revisar conflictos y confirmar.
- Si Android tenía configurada una carpeta local con SAF, mantiene esa configuración y pide revisión antes de migrar. No se cambia ni borra ninguna carpeta en silencio.
- La selección manual de carpetas (incluida la local) pasa a «Opciones avanzadas» para quien la necesite.
- La versión sigue funcionando sin cuenta Microsoft, sin conexión a Internet y sin servidor propio.
- Windows y macOS conservan la conexión mediante carpeta local OneDrive y deben utilizar la **misma carpeta `Selector_Cine_365`** para compartir biblioteca; la autenticación Microsoft Graph nativa de escritorio no está implementada.
- Código de versión **9019** y nombre **9.0.0-beta.13.5**, mismo ID `es.cine365.app`, firmado con certificado anterior en la APK distribuida para pruebas. Copia JSON antes de instalar y nunca desinstales para actualizar.
- [Notas completas Beta 13.5](../RELEASE_NOTES_v9_BETA13_5_ANDROID.md).

## Android Beta 13.3 — guardado local y sincronización manual

- **Una lista no es una subida:** al cambiar de lista no se muestra «Guardando…»; si la selección cambia, se conserva en la memoria local sin subir nada.
- Todas las modificaciones de películas o listas se guardan automáticamente **solo en el móvil**, sin consultar ni publicar en OneDrive.
- **OneDrive solo se consulta al abrir la app**, una vez y de forma **solo lectura** cuando ya se ha vinculado. No se importan ni envían datos automáticamente.
- Nuevo botón compacto de dos flechas circulares en la parte superior, junto al estado de guardado. Al pulsarlo se abre Nube, se muestra un resumen comparado y debes **confirmar** para enviar o recuperar cambios; ante conflictos eliges qué conservar.
- No quedan comprobaciones periódicas ni al regresar a la app desde otra pantalla, y no se consulta TMDB al pasar de una lista a otra.
- Sigue disponible el botón específico «Enviar copia» en Nube para subir únicamente por voluntad del usuario, pero para recuperar y combinar cambios se recomienda el de las dos flechas.
- **Android versión 9017** `es.cine365.app` y el mismo certificado de todas las Beta 11–13.2: instalar la nueva APK encima sin borrar datos. Exporta un JSON antes de actualizar.
- **Windows y Mac Beta 13.1 no cambian**. Pueden seguir sincronizando automáticamente según su configuración, así que en Android habrá que pulsar ↻ y confirmar después de los cambios del PC.
- Sin red ni OneDrive conectado, la app sigue guardando tus listas localmente como siempre.

[Notas de versión Beta 13.3](../RELEASE_NOTES_v9_BETA13_3_ANDROID.md).

## Hotfix Android Beta 13.2 — consumo y actividad en segundo plano

- Solo hay una inicialización del motor Java al arrancar, no dos.
- Al minimizar, la WebView pausa sus temporizadores JavaScript; al volver, se reanudan.
- Se quitan los pings internos cada 12 segundos, innecesarios en Android.
- OneDrive se comprueba automáticamente cada **60 segundos mientras la app está visible**, además de al retomarla. El backend no inicia comprobaciones automáticas cuando está en segundo plano.
- Si acabas de marcar una película y minimizas antes del guardado de 240 ms, se escribe localmente y se conserva su revisión de sincronización; la copia de nube se mantiene.
- `versionCode 9016`, `versionName 9.0.0-beta.13.2`, mismo paquete y certificado que Beta 13. **Se instala encima sin desinstalar.** Las claves y ajustes permanecen.
- Las mejoras de consumo deben medirse con varias horas de uso en un móvil físico; no hay una cifra de ahorro prometida.

[Notas completas de Android Beta 13.2](../RELEASE_NOTES_v9_BETA13_2_ANDROID.md).

## Compatibilidad con la release conjunta Beta 13.1

La release `v9.0.0-beta.13.1` incluye **esta misma APK Android Beta 13**, `es.cine365.app`, versión interna **9015**, firmada con el mismo certificado persistente que Beta 11–12.2. No se ha cambiado Android para este hotfix: Beta 13.1 repara exclusivamente una ruta de caché TMDB del motor Go de Windows y Mac. **No desinstales Android si ya tienes Beta 13.** Para sincronizar listas eliminadas e identidades manuales utiliza Windows/Mac Beta 13.1 y Android Beta 13.

## Actualización Beta 13 — sincronización de listas e identidades

- Ahora se guardan en la nube los ID de las listas eliminadas para que un snapshot antiguo no pueda recuperarlas involuntariamente.
- Los ID corregidos manualmente en TMDB con el lápiz se incluyen en la biblioteca y se aplican en el otro dispositivo cuando se sincroniza.
- La clave de TMDB continúa siendo local y privada; no se sincroniza.
- La app reintenta consultar el estado TMDB al inicio sin confundir un error de conexión con la pérdida de la clave.
- Se mantienen las zonas seguras para notch y barra de gestos de Beta 12.2.

**Instalación:** Beta 13 usa `es.cine365.app`, `versionCode 9015` y el mismo certificado de Beta 11/12/12.1/12.2. Instala encima sin desinstalar. Exporta JSON antes y actualiza también Windows/Mac antes de borrar listas; la Beta 9 antigua sigue necesitando importación de JSON.

## Actualización Beta 12.2: notch y barra inferior adaptables

Beta 12.1 corrigió el título tapado por la cámara frontal/notch utilizando los `WindowInsets` nativos de `systemBars` y `displayCutout`. En algunos Redmi, la barra inferior quedaba demasiado arriba porque el HTML sumaba de nuevo `env(safe-area-inset-bottom)` a un espacio que Android ya había reservado. **Beta 12.2 elimina esa segunda reserva**, manteniendo el ajuste nativo en los cuatro lados; no se usa un margen fijo por modelo. Debe comprobarse visualmente en el Redmi con gestos, botones clásicos y rotación antes de considerarlo completamente validado.

**Archivo histórico Beta 12.2:** desde Beta 11, Beta 12 o Beta 12.1 se podía instalar aquella APK encima, **sin desinstalar**. Conserva `es.cine365.app` y la misma firma; el código de versión pasa a `9014`. Recomendamos una copia JSON antes de actualizar.

## Instalación: elegir la situación correcta

**Si ya utilizas Beta 11 con el paquete `es.cine365.app`:** exporta por precaución un JSON, instala la APK **firmada** Beta 12.2 (versionCode `9014`) como actualización de `9011`, `9012` o `9013`. Debe conservar tus datos siempre que el certificado sea el **mismo que Beta 11**. No desinstales previamente.

**Si todavía tienes Beta 9 (`es.cine365.app.v9beta9`)**: su identificador y firma son diferentes, por lo que **no se puede actualizar directamente** conservando el almacenamiento privado. Antes de desinstalar, exporta un JSON fuera de la app, comprueba su existencia, desinstala Beta 9 si no quieres dos apps, instala la Beta 12.2 firmada, importa el JSON y revisa listas y marcas **sin conectar OneDrive aún**. Solo después reconecta y confirma una primera revisión de sincronización.

**Las APK de `assembleRelease` de GitHub Actions salen SIN FIRMA cuando la clave privada no se ha configurado allí. No las instales ni publiques como versión oficial.** La clave de firma debe guardarse aparte, fuera de GitHub, y reutilizarse para las siguientes betas. Mantener `applicationId 'es.cine365.app'` y aumentar `versionCode` en cada actualización.

## Lo que incluye Beta 12.2

- Búsqueda en lista actual y **todas las listas**, con coincidencias agrupadas por identidad TMDB o título/año.
- Listas numeradas en selector y gestión; cada ficha indica a cuáles pertenece y permite navegar o editar pertenencia.
- Año junto a nota **TMDB** (no IMDb). Estado Vista/Pendiente sincronizado entre listas; **solo el tick lo cambia**.
- Póster, disponibilidad en España, filtros por plataformas, sorteo y corrección manual de identificaciones.
- Cabecera y navegación propias para móvil, contenido desplazable, carga progresiva de películas y estadísticas desplazables sin paginación.
- Nube opcional, acceso a OneDrive personal mediante navegador Microsoft con **OAuth2 + PKCE**, guardado local, copias, revisiones causales, historial y recuperación.
- Correcciones Beta 10 mantenidas: protección contra «resurrección» de ticks antiguos, Sincronizar ahora, estado de última comprobación, y sin bloqueo general de la interfaz por una consulta automática.

## Configuración de Microsoft Entra

El cliente público registrado usa el ID `696458fe-0672-47e5-9fdd-aa0332e9884b` y la URI de retorno exacta:

```text
msal696458fe-0672-47e5-9fdd-aa0332e9884b://auth
```

Permisos delegados: `Files.ReadWrite` y `User.Read`, cuentas personales de Microsoft. No se necesita contraseña de cliente ni crear un registro nuevo en Azure para cada usuario. Si la cuenta se conectaba correctamente en betas anteriores, conserva esta configuración.

Para usar OneDrive: entra en **Nube → Conectar cuenta Microsoft**, autoriza en el navegador externo, escoge la **misma carpeta de OneDrive** que utiliza Windows/Mac y pulsa **Revisar sincronización** antes de aceptar cambios. Seleccionar una carpeta no publica datos por sí solo.

## Cómo se sincronizan las marcas

- Las acciones sobre una película llevan identificador de dispositivo y revisiones causales. Un cambio nuevo que conoce al anterior prevalece, aunque los relojes locales no coincidan.
- Si dos dispositivos modifican el mismo estado **sin conocer mutuamente el cambio**, la app pide intervención. No se elige silenciosamente el estado del dispositivo con el reloj más adelantado.
- La app comprueba OneDrive periódicamente mientras se utiliza y al recuperar el primer plano. Los cambios pueden tardar en propagarse; no es un servicio en segundo plano permanente.
- Tras confirmar la primera asociación, el sistema no debe tratar copias históricas sin revisión como autoridad sobre marcas más recientes. Guarda copias JSON antes de restauraciones.
- La simultaneidad de borrados/renombrados de listas y la sincronización con archivos a medio subir siguen siendo áreas experimentales.

## Historial y restauración

**Actualizar historial** vuelve a leer los eventos **registrados localmente** por las versiones que incluyen este registro. No importa películas, no sincroniza y no reconstruye acciones realizadas en versiones antiguas. Para comprobar cambios entre dispositivos utiliza **Sincronizar ahora** o **Revisar sincronización**.

Si restableces una copia, confirma que el contenido sea correcto antes de dejar que vuelva a compartir cambios con la nube.

## Compilación

Con Java 17, Gradle 8.7 y Android SDK 35:

```bash
cd android
gradle testDebugUnitTest assembleRelease
```

Después firma y verifica el APK usando el certificado persistente del paquete `es.cine365.app`. No incluyas el keystore ni su contraseña en el repositorio, commits, releases ni logs. Una firma diferente impediría instalar como actualización aunque el `versionCode` aumente.

## Limitaciones de pruebas

Las pruebas de compilación, JavaScript y motor causal superadas no garantizan por sí solas el comportamiento en todos los modelos de teléfono ni la aceptación de Microsoft Entra en una sesión real. Recomendamos comprobar la actualización y la sincronización antes de considerar esta versión estable.

El programa utiliza datos de TMDB sin patrocinio ni certificación de TMDB. © 2026 Cristian Barallobre Rúa.
