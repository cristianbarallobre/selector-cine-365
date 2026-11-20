# Selector Cine 365 · Android v9 beta 6 / Windows v9 beta 6

**Versión experimental.** Instalar Android beta 6 no sustituye a beta 4: paquete de prueba es.cine365.app.v9beta6. Exporta/importa tu copia JSON si quieres trasladar tus listas. La app de Windows beta 6 usa el directorio de datos habitual, así que guarda una copia y **no ejecutes dos versiones de Windows simultáneamente**.

## Qué cambia

- Ayuda clara para el usuario: **OneDrive es la opción recomendada para compartir Windows y Android**. El botón «Elegir carpeta local» no abre OneDrive en el móvil.
- Código de acceso directo a OneDrive PERSONAL en Android mediante API oficial **Microsoft Graph**, con selección de carpeta de primer nivel, creación opcional de Pruebas_Selector_Cine_365, desconexión y conservación del modo solo local.
- Inicio de sesión en Microsoft en navegador externo mediante flujo OAuth 2.0 de **código de dispositivo**: la app muestra un código de un solo uso. Los tokens se cifran con Android Keystore. **La app nunca solicita la contraseña de Microsoft.**
- La sincronización existente sigue siendo experimental y **requiere revisión manual al recibir cambios**. Cada dispositivo escribe su archivo propio en Selector_Cine_365_Copias/sincronizacion/dispositivos.
- Las listas y marcas se guardan localmente aunque OneDrive esté desconectado. No se suben tokens TMDB ni se vuelve a consultar TMDB solo por conectar OneDrive.
- Continúa la compatibilidad opcional con carpetas locales accesibles mediante Android Storage Access Framework; otras nubes dependen de la compatibilidad real del proveedor.

## Activación del acceso a OneDrive en Android Beta 6

**La Beta 6 incorpora el identificador público de Selector Cine 365 en Microsoft Entra.** La conexión OneDrive ya se puede iniciar desde la interfaz. La autorización real sigue dependiendo de la configuración del registro en Microsoft y de la prueba con una cuenta personal. La app no requiere contraseña Microsoft ni secretos.

1. La app ya está registrada en Microsoft Entra. Verificar que el tipo de cuenta sigue siendo **cuentas personales de Microsoft**.
2. En **Autenticación / Configuración avanzada**, habilitar **flujos de cliente público** para el flujo de código de dispositivo. No se necesita secreto de cliente ni URI de redirección para este flujo.
3. En **Permisos de API → Microsoft Graph → Delegados**, configurar Files.ReadWrite y User.Read. offline_access se solicita para renovar sesiones. Para la selección libre de carpetas se necesita Files.ReadWrite: la variante Files.ReadWrite.AppFolder solo serviría para la carpeta propia reservada de la app.
4. El **ID de aplicación (cliente)** ya está incorporado a la compilación oficial de Android Beta 6. Es público; **no se necesita ni se debe crear un client secret** para este flujo.
5. Instalar Beta 6 en Android real, conectar con Microsoft y verificar que se ven las carpetas de OneDrive. No se puede afirmar compatibilidad real hasta completar esta prueba.

**Advertencia:** el flujo de código de dispositivo requiere escribir un código temporal en la página oficial de Microsoft; no es todavía el inicio de sesión con una sola pulsación mediante MSAL que sería deseable para producción. En OneDrive Empresarial no está habilitado este flujo en esta beta, que usa el endpoint /consumers.

## Guía de pruebas (cuando esté registrado Microsoft)

1. En Windows abre tu carpeta local de OneDrive y crea Pruebas_Selector_Cine_365. Selecciona esa carpeta en Selector Cine 365 Windows Beta 6. No elijas la carpeta que contiene tus copias reales.
2. En Android abre **Nube → Conectar cuenta Microsoft**. Abre Microsoft y escribe el código temporal. Vuelve a la app.
3. En el desplegable de OneDrive selecciona Pruebas_Selector_Cine_365 o créala desde la propia app. Pulsa **Utilizar esta carpeta**.
4. En ambos equipos pulsa **Revisar sincronización**; inspecciona las cantidades y resuelve conflictos. **La selección sola no publica tus películas.**
5. Marca una película en Windows, revisa la sincronización en Android, y viceversa.
6. Para detener la nube usa **Desactivar** o **Desconectar cuenta Microsoft**. Los archivos remotos no se borran.
7. Haz copia JSON antes de hacer pruebas de desmarcado simultáneo o cambios de listas.

## Protecciones y límites

- No se publican datos al conectar la cuenta ni al elegir una carpeta: se exige primera revisión/confirmación.
- Antes de importar, se crea respaldo local. Cada dispositivo publica su JSON, sin sobrescribir el archivo independiente de otro.
- La sincronización bidireccional de entrada **aún no es automática**: hay que pulsar «Revisar sincronización».
- Renombrados/eliminaciones de listas simultáneos y recuperación ante cambios concurrentes siguen siendo casos experimentales.
- Se conserva caché persistente de TMDB sin actualizaciones automáticas.
- Microsoft Graph puede introducir límites, permisos y revisiones de consentimiento según la cuenta.

**Autor:** Cristian Barallobre Rúa. Este producto utiliza TMDB pero no está respaldado ni certificado por TMDB.
