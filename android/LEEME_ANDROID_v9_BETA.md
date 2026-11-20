# Selector Cine 365 · Android v9 Beta 10

**Beta experimental para OneDrive PERSONAL.** Se mantiene el inicio de sesión en el navegador Microsoft mediante **OAuth 2.0 Authorization Code con PKCE**, que ya funcionaba en Beta 7. La novedad de Beta 10 es la recepción automática de cambios tras confirmar la primera conexión, manteniendo la resolución manual de conflictos.

## Configuración de Microsoft Entra (sin cambios respecto a Beta 7)

**Si ya conectaste OneDrive con Beta 7, no necesitas cambiar nada en Azure.** La aplicación mantiene el mismo identificador público y la misma dirección de retorno:

```text
msal696458fe-0672-47e5-9fdd-aa0332e9884b://auth
```

Si necesitas verificar la configuración en otro móvil, el registro de Microsoft Entra se encuentra en **Registros de aplicaciones → Selector Cine 365 → Autenticación**.
Debe estar seleccionada la plataforma **Aplicaciones móviles y de escritorio**.
La dirección anterior debe estar guardada como URI de redirección.
Mantén **Tipos de cuenta compatibles → Solo cuentas personales**, **Permitir flujos de cliente público → Sí** y en **Permisos de API** los permisos delegados `Files.ReadWrite` y `User.Read`.
No añadas ninguna contraseña ni secreto de cliente. No se necesita habilitar servicios de Azure de pago.

**No vuelvas a registrar la aplicación** si ya funcionaba la conexión en Beta 7.

## Instalación y prueba en Android

- APK Beta 10: paquete `es.cine365.app.v9beta10`, separado de la Beta 7, por lo que **no sobrescribe datos anteriores**.
- **Antes de las pruebas**, exporta un archivo JSON desde tu instalación actual y restáuralo en Beta 10 si quieres trasladar tus películas. **No desinstales la Beta 7** hasta confirmar las listas.
- **Atención con el navegador:** Beta 7 y Beta 10 usan la misma dirección de retorno de Microsoft. Si Android te pregunta qué aplicación debe abrir, elige **Selector Cine 365 · Beta 10**. Si regresa automáticamente a Beta 7, borra los valores predeterminados de apertura de Beta 7 en Ajustes de Android antes de repetir el inicio de sesión; no borres sus datos sin conservar antes un JSON.
- Pulsa **Nube → Conectar cuenta Microsoft**. La app abre la página oficial de Microsoft automáticamente en el navegador del móvil.
- Inicia sesión y acepta los permisos de Microsoft Graph. Al finalizar, el navegador debería **volver directamente a Selector Cine 365** sin solicitar ningún código manual.
- Selecciona o crea una carpeta de pruebas (por ejemplo `Selector_Cine_365`) en OneDrive y pulsa **Utilizar esta carpeta**. Eso **no publica** datos por sí solo.
- En Windows Beta 10, selecciona esa misma carpeta local sincronizada dentro de OneDrive.
- Pulsa **Revisar sincronización** en los dos dispositivos; inspecciona los datos y confirma. Después de la primera revisión, los cambios compatibles entrantes se reciben automáticamente. Solo se revisan los conflictos.

**La autenticación completa y el retorno del navegador aún deben validarse en un móvil real.** El éxito de compilación no demuestra que Azure haya aceptado la URI o que Microsoft permita el consentimiento.

## Beta 10: sincronización causal, actividad y recuperación

- **No basta con comparar fechas:** cada marca Vista/Pendiente contiene un identificador de dispositivo, número de secuencia y reloj causal. Una acción que conocía la anterior tiene prioridad, aunque sus equipos tengan relojes distintos. Cuando dos acciones opuestas son concurrentes, el sistema solicita revisión manual.
- Tras vincular OneDrive por primera vez, **las copias históricas sin revisión ya no pueden volver a marcar películas por sí solas**. El motor también deja de considerar el archivo antiguo `actual/colecciones.json` como autoridad sobre las marcas una vez vinculado.
- Las listas con películas repetidas respetan las revisiones causales al reconciliar las marcas (se evita que aparezca de nuevo un tick anterior).
- **Sincronizar ahora** fuerza una comprobación sin acceder a TMDB y muestra un aviso de progreso mientras opera.
- En Android la detección incremental de archivos usa Microsoft Graph Delta Query cuando está disponible; si falla por incompatibilidad del servicio se realiza una revisión completa segura. En primer plano se solicita una comprobación aproximadamente cada 30 segundos, con un mínimo de 25 entre comprobaciones. La propagación real de OneDrive puede tardar más.
- En Windows se comprueba la carpeta local sincronizada aproximadamente cada 25 segundos. No se promete una latencia exacta ni se mantiene actividad cuando las aplicaciones están cerradas.
- **Historial** muestra las decisiones que registra esa instalación y los cambios recibidos desde otros dispositivos. Los horarios se muestran como referencia; el orden causal NO depende exclusivamente de los relojes.
- **Copias de recuperación:** se pueden explorar copias locales, confirmar la restauración y rescatar primero el estado previo. Restaurar un estado anterior implica nuevas modificaciones locales que podrían compartirse si la nube está activada.
- Las credenciales de Microsoft y la configuración del registro de Azure **no cambian**; solo se reinstala la nueva Beta y se selecciona la misma carpeta de nube.
- **Precaución Beta10:** APK diferente de Beta9 (`es.cine365.app.v9beta10`). Antes de probar exporta una copia JSON, instálala sin desinstalar la anterior e importa los datos. Android puede preguntar qué Beta abrir al volver del navegador de Microsoft; escoge Beta10.
- **No hay verificación física todavía:** los tests cubren la secuencia Windows visto → Android recibe → Android pendiente → Windows recibe, pero debe confirmarse en el móvil Xiaomi y el PC reales.

## Sincronización automática después de la primera revisión y bloqueos de Nube

- **Primera vez:** confirma la sincronización en Windows y Android con «Revisar sincronización», usando la misma carpeta de OneDrive.
- **Estabilidad Android Beta 10:** ya no se desactiva toda la interfaz durante una consulta automática a OneDrive; solo se protegen los controles de Nube mientras una operación manual está ejecutándose, con un aviso visible.\n- **Siguientes veces:** en ambos dispositivos los cambios seguros se reciben automáticamente **al abrir, volver a la app y mientras se mantiene abierta**. La comprobación es periódica, no instantánea: OneDrive necesita tiempo para subir los archivos.
- **Conflictos:** si se ha marcado y desmarcado la misma película en dos dispositivos sin haber recibido el cambio del otro, **NO se decide automáticamente**. Pulsa «Revisar sincronización» para escoger qué conservar.
- **Sin nube:** todos los cambios se guardan en cada dispositivo de forma local, sin requerir OneDrive ni TMDB.
- **Uso de datos:** no se consulta TMDB para sincronizar. La aplicación revisa los archivos de OneDrive con intervalos limitados; no funciona como un servicio permanente en segundo plano.
- **Límites todavía experimentales:** renombrados o eliminaciones de listas hechos simultáneamente en dos dispositivos, o archivos que no hayan terminado de subirse a OneDrive. Conservar copias JSON de recuperación.

## Windows Beta 10

Compilada específicamente con el subsistema **Windows GUI**, con comprobación automática de la cabecera PE. **No abre una ventana CMD**, a diferencia del viejo ejecutable Beta 5. Si en el título de una ventana aparece BETA5, cierra esa versión y abre exactamente el EXE `Selector_Cine_365_Windows_v9_BETA8.exe`. No ejecutes dos versiones simultáneamente: ambas comparten datos locales.

## Seguridad y funcionamiento

- PKCE S256, estado aleatorio verificado y válido una sola vez, redirect comprobado en la actividad Android, login en navegador externo; no se envían contraseñas ni tokens a la WebView.
- Tokens cifrados con **Android Keystore** y almacenados en memoria interna de la app; «Desconectar» borra la sesión local.
- Los datos se guardan localmente incluso sin nube. Las solicitudes a TMDB solo ocurren para películas nuevas no guardadas o cuando se pide actualizar; se mantiene la caché.
- Cada dispositivo conserva su propio snapshot de sincronización. Los cambios entrantes no son todavía automáticos y hay límites experimentales con renombrados/eliminaciones de listas.
- «Elegir carpeta local» utiliza Android SAF y **no significa conectar directamente OneDrive**. Otras nubes son opcionales si ofrecen carpetas de lectura y escritura reales.

**Autor:** Cristian Barallobre Rúa. Selector Cine 365 utiliza datos de TMDB sin patrocinio ni certificación de TMDB.
