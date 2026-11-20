# Selector Cine 365 · Android v9 Beta 8

**Beta experimental para OneDrive PERSONAL.** Se ha sustituido el inicio de sesión con código temporal por el inicio interactivo desde navegador Microsoft mediante **OAuth 2.0 Authorization Code con PKCE**. La Beta 7 anterior no debe usarse para probar el login: su método de código fallaba en las pruebas reales.

## Requisito de Microsoft Entra: configurar una sola dirección de retorno

La aplicación ya tiene registrado el ID público de cliente de Microsoft. Antes del primer inicio interactivo, hay que **añadir** en el mismo registro de Microsoft Entra de Selector Cine 365 esta dirección de retorno exacta:

```text
msal696458fe-0672-47e5-9fdd-aa0332e9884b://auth
```

1. En **Azure → Microsoft Entra ID → Registros de aplicaciones → Selector Cine 365 → Autenticación** (o Authentication Preview).
2. En **Configuración de URI de redirección**, pulsa **Agregar una plataforma** y elige **Aplicaciones móviles y de escritorio** (*Mobile and desktop applications*).
3. Si permite URI personalizados, pega exactamente **msal696458fe-0672-47e5-9fdd-aa0332e9884b://auth** y pulsa Guardar. Microsoft documenta este formato `msal{client_id}://auth` para aplicaciones móviles/escritorio.
4. Mantén **Tipos de cuenta compatibles → Solo cuentas personales**, **Permitir flujos de cliente público → Sí** y en **Permisos de API** los permisos delegados `Files.ReadWrite` y `User.Read`.
5. No añadas ninguna contraseña ni secreto de cliente. No se necesita habilitar servicios de Azure de pago.

**Si la página no permite registrar esa URI o muestra un error, no improvises otra.** Copia el mensaje de Azure y revisaremos la plataforma que utiliza el portal antes de distribuir otra APK.

## Instalación y prueba en Android

- APK Beta 8: paquete `es.cine365.app.v9beta8`, separado de la Beta 7, por lo que **no sobrescribe datos anteriores**.
- **Antes de las pruebas**, exporta un archivo JSON desde tu instalación actual y restáuralo en Beta 8 si quieres trasladar tus películas. **No desinstales la Beta 7** hasta confirmar las listas.
- Pulsa **Nube → Conectar cuenta Microsoft**. La app abre la página oficial de Microsoft automáticamente en el navegador del móvil.
- Inicia sesión y acepta los permisos de Microsoft Graph. Al finalizar, el navegador debería **volver directamente a Selector Cine 365** sin solicitar ningún código manual.
- Selecciona o crea una carpeta de pruebas (por ejemplo `Pruebas_Selector_Cine_365`) en OneDrive y pulsa **Utilizar esta carpeta**. Eso **no publica** datos por sí solo.
- En Windows Beta 8, selecciona esa misma carpeta local sincronizada dentro de OneDrive.
- Pulsa **Revisar sincronización** en los dos dispositivos; inspecciona los datos y confirma. Para recibir nuevos cambios del otro dispositivo, **hay que revisar de nuevo**.

**La autenticación completa y el retorno del navegador aún deben validarse en un móvil real.** El éxito de compilación no demuestra que Azure haya aceptado la URI o que Microsoft permita el consentimiento.

## Sincronización automática después de la primera revisión

- **Primera vez:** confirma la sincronización en Windows y Android con «Revisar sincronización», usando la misma carpeta de OneDrive.
- **Siguientes veces:** en ambos dispositivos los cambios seguros se reciben automáticamente **al abrir, volver a la app y mientras se mantiene abierta**. La comprobación es periódica, no instantánea: OneDrive necesita tiempo para subir los archivos.
- **Conflictos:** si se ha marcado y desmarcado la misma película en dos dispositivos sin haber recibido el cambio del otro, **NO se decide automáticamente**. Pulsa «Revisar sincronización» para escoger qué conservar.
- **Sin nube:** todos los cambios se guardan en cada dispositivo de forma local, sin requerir OneDrive ni TMDB.
- **Uso de datos:** no se consulta TMDB para sincronizar. La aplicación revisa los archivos de OneDrive con intervalos limitados; no funciona como un servicio permanente en segundo plano.
- **Límites todavía experimentales:** renombrados o eliminaciones de listas hechos simultáneamente en dos dispositivos, o archivos que no hayan terminado de subirse a OneDrive. Conservar copias JSON de recuperación.

## Windows Beta 8

Compilada específicamente con el subsistema **Windows GUI**, con comprobación automática de la cabecera PE. **No abre una ventana CMD**, a diferencia del viejo ejecutable Beta 5. Si en el título de una ventana aparece BETA5, cierra esa versión y abre exactamente el EXE `Selector_Cine_365_Windows_v9_BETA8.exe`. No ejecutes dos versiones simultáneamente: ambas comparten datos locales.

## Seguridad y funcionamiento

- PKCE S256, estado aleatorio verificado y válido una sola vez, redirect comprobado en la actividad Android, login en navegador externo; no se envían contraseñas ni tokens a la WebView.
- Tokens cifrados con **Android Keystore** y almacenados en memoria interna de la app; «Desconectar» borra la sesión local.
- Los datos se guardan localmente incluso sin nube. Las solicitudes a TMDB solo ocurren para películas nuevas no guardadas o cuando se pide actualizar; se mantiene la caché.
- Cada dispositivo conserva su propio snapshot de sincronización. Los cambios entrantes no son todavía automáticos y hay límites experimentales con renombrados/eliminaciones de listas.
- «Elegir carpeta local» utiliza Android SAF y **no significa conectar directamente OneDrive**. Otras nubes son opcionales si ofrecen carpetas de lectura y escritura reales.

**Autor:** Cristian Barallobre Rúa. Selector Cine 365 utiliza datos de TMDB sin patrocinio ni certificación de TMDB.
