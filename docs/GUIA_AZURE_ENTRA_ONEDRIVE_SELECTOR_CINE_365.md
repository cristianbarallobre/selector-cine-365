# Configuración de OneDrive y Microsoft Entra

Guía para Windows y macOS 9.2.7 y la APK estable de Android. La app incluye una configuración pública de Microsoft. Puedes conectar tu cuenta desde la app y consultar esta guía si administras el registro o tienes un problema de autorización.

## Conectar una cuenta personal

En Windows y Mac abre **Mis listas** y busca el apartado de sincronización con OneDrive. En Android entra en **Ajustes**.

1. Pulsa **Conectar cuenta Microsoft**.
2. Autoriza el acceso en el navegador. Al terminar, vuelve a Selector Cine 365.
3. Usa la misma cuenta en todos los equipos que deban compartir la biblioteca.
4. Revisa la primera sincronización y confirma después de comprobar los recuentos.

Windows usa DPAPI para proteger la autorización, macOS usa el llavero y Android su almacenamiento protegido. Si Microsoft revoca la autorización, conecta la cuenta de nuevo.

## Registro utilizado por la app

- El ID público de aplicación es 696458fe-0672-47e5-9fdd-aa0332e9884b.
- El registro es un cliente público para cuentas personales de Microsoft.
- El flujo utiliza un código de autorización con PKCE S256.
- La redirección de escritorio es http://localhost.
- La redirección de Android es msal696458fe-0672-47e5-9fdd-aa0332e9884b://auth.
- Los permisos delegados solicitados son Files.ReadWrite, User.Read y offline_access.

El ID identifica el registro. Mantén privados los tokens y las contraseñas. La app utiliza el navegador para pedir autorización.

## Revisar el registro

1. Accede al portal de Azure con la cuenta que administra el registro de Selector Cine 365.
2. En **Microsoft Entra ID → Registros de aplicaciones**, comprueba que el ID coincide.
3. En **Autenticación**, revisa el tipo de cuentas, el cliente público y las dos redirecciones. Conserva ambas para escritorio y Android.
4. En **Permisos de API**, verifica los permisos delegados de Microsoft Graph. Autoriza desde la app para probar el retorno.
5. Si el retorno falla, comprueba qué edición de Android abrió el navegador o vuelve a iniciar desde la app de escritorio para obtener un nuevo enlace.

[Portal de Azure](https://portal.azure.com/) · [Sincronización segura](https://github.com/cristianbarallobre/selector-cine-365/blob/main/docs/GUIA_SINCRONIZACION_SEGURA.md)

Autor: Cristian Barallobre Rúa.
