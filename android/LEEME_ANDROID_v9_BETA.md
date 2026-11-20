# Selector Cine 365 — Android v9 Beta 12

**Autor:** Cristian Barallobre Rúa. **Estado:** versión beta experimental para Android 8.0+ (API 26). Conserva funcionamiento local sin cuenta Microsoft y sin clave TMDB.

## Instalación: elegir la situación correcta

**Si ya utilizas Beta 11 con el paquete \`es.cine365.app\`:** exporta por precaución un JSON, instala la APK **firmada** Beta 12 (versionCode \`9012\`) como actualización de \`9011\`. Debe conservar tus datos siempre que el certificado sea el **mismo que Beta 11**. No desinstales previamente.

**Si todavía tienes Beta 9 (\`es.cine365.app.v9beta9\`)**: su identificador y firma son diferentes, por lo que **no se puede actualizar directamente** conservando el almacenamiento privado. Antes de desinstalar, exporta un JSON fuera de la app, comprueba su existencia, desinstala Beta 9 si no quieres dos apps, instala la Beta 12 firmada, importa el JSON y revisa listas y marcas **sin conectar OneDrive aún**. Solo después reconecta y confirma una primera revisión de sincronización.

**Las APK de \`assembleRelease\` de GitHub Actions salen SIN FIRMA cuando la clave privada no se ha configurado allí. No las instales ni publiques como versión oficial.** La clave de firma debe guardarse aparte, fuera de GitHub, y reutilizarse para las siguientes betas. Mantener \`applicationId 'es.cine365.app'\` y aumentar \`versionCode\` en cada actualización.

## Lo que incluye Beta 12

- Búsqueda en lista actual y **todas las listas**, con coincidencias agrupadas por identidad TMDB o título/año.
- Listas numeradas en selector y gestión; cada ficha indica a cuáles pertenece y permite navegar o editar pertenencia.
- Año junto a nota **TMDB** (no IMDb). Estado Vista/Pendiente sincronizado entre listas; **solo el tick lo cambia**.
- Póster, disponibilidad en España, filtros por plataformas, sorteo y corrección manual de identificaciones.
- Cabecera y navegación propias para móvil, contenido desplazable, carga progresiva de películas y estadísticas desplazables sin paginación.
- Nube opcional, acceso a OneDrive personal mediante navegador Microsoft con **OAuth2 + PKCE**, guardado local, copias, revisiones causales, historial y recuperación.
- Correcciones Beta 10 mantenidas: protección contra «resurrección» de ticks antiguos, Sincronizar ahora, estado de última comprobación, y sin bloqueo general de la interfaz por una consulta automática.

## Configuración de Microsoft Entra

El cliente público registrado usa el ID \`696458fe-0672-47e5-9fdd-aa0332e9884b\` y la URI de retorno exacta:

\`\`\`text
msal696458fe-0672-47e5-9fdd-aa0332e9884b://auth
\`\`\`

Permisos delegados: \`Files.ReadWrite\` y \`User.Read\`, cuentas personales de Microsoft. No se necesita contraseña de cliente ni crear un registro nuevo en Azure para cada usuario. Si la cuenta se conectaba correctamente en betas anteriores, conserva esta configuración.

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

\`\`\`bash
cd android
gradle testDebugUnitTest assembleRelease
\`\`\`

Después firma y verifica el APK usando el certificado persistente del paquete \`es.cine365.app\`. No incluyas el keystore ni su contraseña en el repositorio, commits, releases ni logs. Una firma diferente impediría instalar como actualización aunque el \`versionCode\` aumente.

## Limitaciones de pruebas

Las pruebas de compilación, JavaScript y motor causal superadas no garantizan por sí solas el comportamiento en todos los modelos de teléfono ni la aceptación de Microsoft Entra en una sesión real. Recomendamos comprobar la actualización y la sincronización antes de considerar esta versión estable.

El programa utiliza datos de TMDB sin patrocinio ni certificación de TMDB. © 2026 Cristian Barallobre Rúa.
