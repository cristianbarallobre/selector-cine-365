# Configuración de OneDrive y Microsoft Entra

Documento histórico para versiones de escritorio anteriores a 9.2.24 y la edición de escritorio anterior. La edición de escritorio vigente no conecta cuentas Microsoft ni requiere configurar Entra.

## Conectar una cuenta personal


1. Pulsa **Conectar cuenta Microsoft**.
2. Autoriza el acceso en el navegador. Al terminar, vuelve a Selector Cine 365.
3. Usa la misma cuenta en todos los equipos que deban compartir la biblioteca.
4. Revisa la primera sincronización y confirma después de comprobar los recuentos.
5. En las versiones antiguas de escritorio, pulsa **Sincronizar con la nube** bajo el título cada vez que quieras enviar o recibir cambios. Una cuenta autorizada permite usar el botón. Los cambios se guardan localmente mientras tanto.

Windows usa DPAPI para proteger la autorización y macOS usa el llavero. Si Microsoft revoca la autorización, conecta la cuenta de nuevo.

## Registro utilizado por la app

- El ID público de aplicación es 696458fe-0672-47e5-9fdd-aa0332e9884b.
- El registro es un cliente público para cuentas personales de Microsoft.
- El flujo utiliza un código de autorización con PKCE S256.
- La redirección de escritorio es http://localhost.
- Los permisos delegados solicitados son Files.ReadWrite, User.Read y offline_access.

El ID identifica el registro. Mantén privados los tokens y las contraseñas. La app utiliza el navegador para pedir autorización.

## Revisar el registro

1. Accede al portal de Azure con la cuenta que administra el registro de Selector Cine 365.
2. En **Microsoft Entra ID → Registros de aplicaciones**, comprueba que el ID coincide.
4. En **Permisos de API**, verifica los permisos delegados de Microsoft Graph. Autoriza desde la app para probar el retorno.

[Portal de Azure](https://portal.azure.com/) · [Guardado y recuperación](https://github.com/cristianbarallobre/selector-cine-365/blob/main/docs/GUIA_SINCRONIZACION_SEGURA.md)

Autor: Cristian Barallobre Rúa.
