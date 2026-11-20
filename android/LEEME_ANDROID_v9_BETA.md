# Selector Cine 365 — Android v9 Beta 13.2

**Autor:** Cristian Barallobre Rúa. **Estado:** versión beta experimental para Android 8.0+ (API 26). Conserva funcionamiento local sin cuenta Microsoft y sin clave TMDB.

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
