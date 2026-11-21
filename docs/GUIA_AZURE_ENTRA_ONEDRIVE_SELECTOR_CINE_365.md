# Selector Cine 365
## Guía de configuración de Microsoft Azure, Entra ID y OneDrive

**Edición:** actualizada para la release v9 Beta 13.1 Windows/Mac y APK Android Beta 13, 20 de noviembre de 2026  
**Alcance:** conexión de OneDrive personal con Selector Cine 365 para Android y escritorio Windows/macOS.  
**Estado comprobado:** el inicio de sesión de Microsoft y la sincronización entre dispositivos llegaron a funcionar tras el cambio al inicio interactivo de Android Beta 7. La sincronización automática y las revisiones causales se conservan en Beta 13, pero deben seguir probándose con dispositivos reales.

> **Idea principal:** no necesitas crear una aplicación de Azure por cada teléfono, ni usar códigos de dispositivo, ni contratar servidores. El registro de Microsoft Entra sirve para autorizar el acceso de Android a OneDrive. En Windows y macOS se utiliza la carpeta de OneDrive sincronizada por el cliente instalado en el ordenador.

## 1. Lo que finalmente funcionó

- **Nombre del registro en Microsoft Entra:** `Selector Cine 365`.
- **Tipo de cuentas admitidas:** `Solo cuentas personales` de Microsoft (OneDrive personal, Outlook, Hotmail, Live).
- **Identificador de aplicación (cliente):** `696458fe-0672-47e5-9fdd-aa0332e9884b`.
- **Plataforma para el retorno de inicio de sesión:** `Aplicaciones móviles y de escritorio`.
- **URI de redirección exacta:** `msal696458fe-0672-47e5-9fdd-aa0332e9884b://auth`.
- **Permitir flujos de cliente público:** `Sí`.
- **Permisos delegados de Microsoft Graph:** `Files.ReadWrite` y `User.Read`.
- **Permiso que solicita la app para mantener la sesión:** `offline_access`.
- **Método de inicio de sesión final:** navegador de Microsoft, **OAuth 2.0 Authorization Code + PKCE**. Microsoft devuelve automáticamente el control a Android.
- **Secreto de cliente:** **NO crear ni introducir** uno. El identificador de cliente es público; los secretos y contraseñas no lo son.

Este resumen sirve para **revisar una instalación que ya está creada**. Los pasos siguientes explican cómo reproducir la configuración desde cero.

## 2. Acceder a Azure y Microsoft Entra

1. Entra en **https://portal.azure.com/** e inicia sesión con la cuenta Microsoft que administra el registro.
2. En el buscador del portal escribe **Microsoft Entra ID** y abre el servicio.
3. En el menú selecciona **Registros de aplicaciones** y busca **Selector Cine 365**.
4. Si el registro **ya existe y funciona, no lo vuelvas a crear**. Entra en él para comprobar o actualizar su configuración.
5. Si no existe y necesitas reproducirlo, pulsa **Nuevo registro**.

**Sobre «Probar Azure gratis»:** durante nuestra configuración se utilizó el acceso a Azure para disponer de un directorio y poder registrar la aplicación. No se necesita crear una máquina virtual, contratar almacenamiento de Azure ni pasar automáticamente a «Pago por uso» para que Selector Cine 365 acceda a OneDrive mediante un registro estándar de Entra. Revisa siempre las condiciones de tu cuenta de Azure por separado.

## 3. Crear el registro correcto

1. En **Nombre** escribe `Selector Cine 365`.
2. En **Tipos de cuenta compatibles** abre el desplegable y elige **Solo cuentas personales**. **No** elijas «Solo mi organización» ni «Solo inquilino único».
3. En el formulario inicial puedes dejar **URI de redirección** vacía: se añadirá después en **Autenticación**.
4. Pulsa **Registrar**.
5. En **Información general** localiza **Id. de aplicación (cliente)**. Es el valor que debe usar la APK Android. No lo confundas con el **Identificador de objeto** ni con el **Id. del directorio (inquilino)**.

**Si vuelves a registrar la aplicación desde cero**, Microsoft generará un identificador nuevo. En ese caso habrá que incorporarlo al código/compilación de Selector Cine 365 y actualizar la URI de retorno; las APK existentes no funcionarán automáticamente con otro identificador.

## 4. Corregir el error que apareció al elegir cuentas personales

Durante la configuración encontramos este aviso al intentar guardar «Solo cuentas personales»:

`Property api.requestedAccessTokenVersion is invalid`

La solución que funcionó fue:

1. En el menú del registro de **Selector Cine 365**, abre **Manifiesto**.
2. Busca con `Ctrl + F` la propiedad `requestedAccessTokenVersion` (dentro de `api`).
3. Establece su valor numérico en **2**, no en `null` ni en `1`:

```json
"requestedAccessTokenVersion": 2
```

4. Guarda el manifiesto.
5. Vuelve a **Autenticación (Preview) → Supported accounts / Tipos de cuenta admitidos**.
6. Elige **Solo cuentas personales** y pulsa **Guardar**. En nuestra configuración esta vez se guardó **sin error**.

Esta corrección responde al error específico que apareció en nuestro registro. Si Microsoft cambia la interfaz o el manifiesto, modifica solo esa propiedad y comprueba el mensaje mostrado por el portal.

## 5. Activar la autenticación de una aplicación móvil

1. En el registro, entra en **Autenticación** o **Authentication (Preview)**.
2. Busca la pestaña **Configuración** o el apartado de **Configuración avanzada**.
3. Ajusta **Permitir flujos de cliente público** a **Sí** y guarda.
4. Abre **Configuración de URI de redirección**.
5. Pulsa **Agregar plataforma** y elige **Aplicaciones móviles y de escritorio** (*Mobile and desktop applications*).
6. Añade como **URI de redirección personalizada** la dirección exacta:

```text
msal696458fe-0672-47e5-9fdd-aa0332e9884b://auth
```

7. Guarda y vuelve a abrir el apartado para comprobar que quedó registrada.

**Es importante conservar todos los caracteres:** `msal` + identificador de aplicación + `://auth`. La app Android espera recibir esa misma dirección cuando termina la autorización.

## 6. Configurar los permisos de Microsoft Graph

1. En el menú del registro abre **Permisos de API**.
2. Pulsa **Agregar un permiso → Microsoft Graph → Permisos delegados**.
3. Añade **`Files.ReadWrite`**. La aplicación lo utiliza para leer y escribir los archivos necesarios en OneDrive.
4. Comprueba que también aparece **`User.Read`**.
5. Guarda. Al iniciar sesión en Android, Microsoft puede pedir al usuario que acepte estos permisos.

La aplicación solicita además el ámbito `offline_access` para poder renovar la sesión, cuando Microsoft lo permita. **No hace falta crear «Certificados y secretos».**

**Alcance del permiso:** `Files.ReadWrite` no está limitado por Microsoft a la carpeta de pruebas; el permiso delegado permite leer y modificar archivos accesibles para la cuenta. **El diseño de Selector Cine 365 utiliza la carpeta elegida** y guarda allí sus propios archivos. Por prudencia, selecciona siempre una carpeta específica para la app y conserva tus copias importantes fuera de las pruebas.

## 7. Conectar OneDrive desde Android (método que funcionó)

1. Instala una versión de Selector Cine 365 Android con inicio de sesión **PKCE** (introducido en **Beta 7** y mantenido hasta **Beta 12**).
2. Abre la aplicación y entra en **Nube → Conectar cuenta Microsoft**.
3. Se abre **automáticamente el navegador del móvil** con la página oficial de Microsoft. **No copies ni pegues códigos temporales**.
4. Inicia sesión con la cuenta personal de Microsoft que contiene el OneDrive que vas a usar.
5. Acepta los permisos solicitados si estás de acuerdo con ellos.
6. Al completar la autorización, el navegador te devuelve a **Selector Cine 365**.
7. En la aplicación selecciona o crea la carpeta **`Selector_Cine_365`** y pulsa **Utilizar esta carpeta**.
8. Seleccionar la carpeta **no fusiona todavía las bibliotecas**. La primera vez pulsa **Revisar sincronización**, comprueba los datos y confirma.

Si conservas una beta antigua instalada en paralelo, Android puede preguntar qué aplicación abre el retorno de Microsoft. **Beta 11 y 12** comparten `es.cine365.app`, por lo que Beta 12 firmada con el mismo certificado actualiza Beta 11 sin duplicarla. **Beta 9** usa otro paquete y no puede reemplazarse conservando sus datos privados: exporta e importa JSON antes de retirarla.

## 8. Conectar Windows o macOS a la misma carpeta

1. Instala o abre el cliente de **OneDrive para Windows** y comprueba que sincroniza tu cuenta personal.
2. Dentro de tu OneDrive crea, si aún no existe, **`Selector_Cine_365`**. Mantenla separada de tus carpetas de trabajo y copias habituales.
3. Abre **Selector Cine 365 para Windows → Nube / Copias y sincronización**.
4. Pulsa **Elegir carpeta** y selecciona en el disco la carpeta de OneDrive **`Selector_Cine_365`**.
5. Comprueba que en Android has seleccionado **esa misma carpeta de OneDrive**, no una carpeta local cualquiera del teléfono.
6. En la **primera vinculación** pulsa **Revisar sincronización** en ambos dispositivos, verifica las listas y los estados, y confirma.

**Diferencia importante:** Windows accede a una carpeta que sincroniza el cliente de OneDrive. Android accede a OneDrive directamente mediante Microsoft Graph. **El registro de Microsoft Entra es para el acceso autorizado de Android**, no para instalar el cliente de OneDrive en Windows.

## 9. Cómo comprobar que realmente funciona

1. Haz primero una **copia JSON** de tus listas en ambos dispositivos.
2. Después de vincularlos, en Windows marca una película de prueba como **vista**.
3. Deja que OneDrive termine de subir los cambios. En Android comprueba que la marca se recibe; si estás usando la versión de sincronización manual, pulsa **Revisar sincronización**.
4. Repite la prueba al revés: marca otra película distinta en Android y comprueba el cambio en Windows.
5. En **Beta 8** se añadió recepción automática después de la primera revisión: al abrir, volver a la app y periódicamente mientras se usa. No es inmediata y **el comportamiento completo todavía debe validarse en dispositivos reales**.
6. Si ambos dispositivos realizan decisiones contrarias sobre la misma película sin haber recibido los cambios del otro, usa **Revisar sincronización** para resolver el conflicto; no debe decidirse silenciosamente.

**No desinstales la última versión que te funciona ni borres las carpetas de OneDrive durante las pruebas.** Los archivos de sincronización se guardan bajo `Selector_Cine_365_Copias/sincronizacion/dispositivos` dentro de la carpeta elegida.

## 10. Errores que encontramos y cómo se resolvieron

### A. `AADSTS16000` al entrar en Microsoft Entra

Ocurrió al intentar usar una cuenta personal en un directorio al que no tenía acceso. Se pudo continuar tras disponer de acceso a Azure y al directorio propio. **No significa que tu OneDrive estuviera dañado.** Si reaparece, comprueba qué directorio y qué cuenta están seleccionados en el portal.

### B. `unauthorized_client` en Android

El registro mostraba **«Solo mi organización»**, incompatible con el OneDrive personal que queríamos utilizar. Cambiamos a **«Solo cuentas personales»** y comprobamos la opción de cliente público. Si el cambio falla con el mensaje de versión de token, aplica primero el apartado 4.

### C. `Property api.requestedAccessTokenVersion is invalid`

Se solucionó estableciendo `"requestedAccessTokenVersion": 2` en **Manifiesto**, guardando y volviendo a elegir **Solo cuentas personales**.

### D. «El código no funcionó» en la página de Microsoft

El primer sistema de Android utilizaba **códigos de dispositivo**. Lo repetimos varias veces sin éxito. La solución definitiva **no fue volver a escribir el código**: se cambió el método de autenticación de la app a **Authorization Code + PKCE**, con navegador y retorno automático. No tiene sentido volver al procedimiento de `microsoft.com/devicelogin` para esta app.

### E. Microsoft no vuelve a Android tras iniciar sesión

Comprueba que está registrada la URI `msal696458fe-0672-47e5-9fdd-aa0332e9884b://auth`, que la app instalada admite PKCE y que Android no está abriendo por defecto una beta anterior.

### F. Parece que no se sincronizan los cambios

Comprueba que Windows y Android apuntan a **la misma carpeta**, que el cliente OneDrive de Windows ha terminado de sincronizarla y que la **primera revisión** fue confirmada en ambos. En betas antiguas la recepción era manual; en Beta 8 hay recepción automática limitada y conflictos que requieren revisión.

## 11. Mantenimiento y seguridad

- El **registro de aplicación** de Microsoft Entra no tiene una fecha de caducidad anual fija. Se conserva mientras el registro y el directorio permanezcan activos y Microsoft permita su uso.
- La **prueba gratuita o el crédito promocional de Azure** y el **registro de aplicación en Entra ID** son cosas distintas. No conviertas la cuenta a pago por uso solo por creer que Selector Cine 365 lo necesita.
- La sesión de OneDrive en Android puede solicitar una **nueva autorización** si caduca, es revocada o Microsoft lo requiere; eso no implica registrar otra aplicación.
- Conserva acceso a la cuenta administradora del registro y atiende los correos oficiales de Microsoft sobre el directorio o la cuenta.
- **No compartas** contraseñas, tokens, claves privadas o secretos de cliente. El **Id. de aplicación (cliente)** y la URI de retorno que figuran aquí son datos públicos de configuración.
- Haz copias JSON antes de experimentar con otra cuenta, otra carpeta, otra beta o bibliotecas distintas. El **modo sin nube** sigue siendo válido.
- La autenticación **no descarga TMDB** por sí misma; la app conserva su caché de películas para evitar consultas innecesarias.

## 12. Enlaces oficiales y archivos del proyecto

- [Portal de Microsoft Azure](https://portal.azure.com/)
- [Portal Microsoft Entra](https://entra.microsoft.com/)
- [Registrar una aplicación](https://learn.microsoft.com/en-us/graph/auth-register-app-v2)
- [Cambiar tipos de cuenta](https://learn.microsoft.com/en-us/entra/identity-platform/howto-modify-supported-accounts)
- [Versión 2 de tokens y manifiesto](https://learn.microsoft.com/en-us/entra/identity-platform/reference-microsoft-graph-app-manifest)
- [Registrar una URI de redirección](https://learn.microsoft.com/en-us/entra/identity-platform/how-to-add-redirect-uri)
- [Flujos de cliente público](https://learn.microsoft.com/en-us/entra/identity-platform/scenario-mobile-app-configuration)
- [OAuth 2.0 con PKCE](https://learn.microsoft.com/en-us/entra/identity-platform/v2-oauth2-auth-code-flow)
- [Permisos Microsoft Graph](https://learn.microsoft.com/en-us/graph/permissions-reference)
- [Repositorio de Selector Cine 365](https://github.com/cristianbarallobre/selector-cine-365)

---

**Nota de conservación:** esta guía documenta la configuración que funcionó para Selector Cine 365 y los problemas reales encontrados en el proceso. La interfaz de Azure/Entra puede cambiar de nombre o ubicación; los valores técnicos (tipo de cuenta, permisos, Id. de cliente y URI de retorno) son la referencia principal mientras la app siga usando este registro.
## Nota sobre Beta 13.1 de escritorio y Beta 13 Android

- El registro público de Microsoft Entra, los permisos delegados y la URI de redirección no han cambiado al avanzar de Beta 12 a Beta 12.2. La corrección de notch y menú inferior solo afecta al diseño Android y no modifica la sincronización.
- No se necesitan credenciales de Azure para utilizar la biblioteca local. La sincronización es opcional.
- Android almacena y protege sus tokens en su espacio de aplicación; **no** los exportes en capturas ni los publiques en GitHub.
- La app solo revisa la carpeta OneDrive seleccionada para su biblioteca; conserva una copia JSON independiente antes de conectar otro dispositivo.
- Esta guía documenta la configuración, no garantiza instantaneidad ni ausencia de conflictos en la sincronización real.

## Cambio en Android Beta 13.3: sincronización manual, sin tocar Entra

La cuenta Microsoft, permisos y carpeta de OneDrive **no requieren ninguna modificación**. Solo cambia cuándo Android intercambia archivos: se realiza una consulta **de solo lectura al abrir la app**, y después ninguna comprobación programada ni al regresar de segundo plano. Para **enviar o recibir** cambios, pulsa el nuevo icono de dos flechas circulares de la cabecera y confirma el resumen de comparación. La marca «Guardado» indica almacenamiento **en el móvil**, no subida a OneDrive. Se mantienen todas las copias de recuperación y la resolución de conflictos. Windows/Mac Beta 13.1 mantienen sus propios mecanismos automáticos.

## Nota de compatibilidad con Beta 13.1

En esta publicación Windows y Mac incluyen la corrección de caché TMDB del servidor local Go. Android permanece en Beta 13 (versión interna 9015) con su firma persistente. **La corrección TMDB no requiere modificar Microsoft Entra ni volver a iniciar sesión en OneDrive**. Los IDs de películas identificadas manualmente y los marcadores de listas eliminadas se comparten en la biblioteca, pero **la clave privada API TMDB nunca se sincroniza**. Antes de compartir borrados entre dispositivos exporta una copia JSON independiente y actualiza Windows/Mac y Android a estas versiones.
