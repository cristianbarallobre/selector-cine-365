# Selector Cine 365 — Beta 14 (vista previa multiplataforma)

Estado: **desarrollo / pruebas**, no lanzamiento estable. No sustituye las releases ya publicadas.

## Objetivo

Conectar una biblioteca de Selector Cine 365 a **la misma carpeta lógica de OneDrive** en Android, Windows y macOS, sin seleccionar rutas manuales. Mantener copia local de cada dispositivo y conflictos por revisión causal.

## Funciones añadidas en escritorio (Windows/macOS)

- **Conectar cuenta Microsoft:** integración experimental directa con Microsoft Graph (permisos delegados `Files.ReadWrite`, `User.Read`, `offline_access`), sin instalar OneDrive.
- Se utiliza el flujo oficial **Authorization Code + PKCE (S256)**: la app abre el navegador del sistema para autorizar tu cuenta y recibe el retorno en una dirección local con estado de seguridad. No pide tu contraseña ni códigos manuales.
- La app reutiliza o crea la carpeta raíz `Selector_Cine_365` de OneDrive, compatible con Android Beta 13.5. Dentro respeta la estructura `Selector_Cine_365_Copias/sincronizacion/dispositivos` existente.
- **Primera sincronización:** vista previa obligatoria, conflictos con elección manual, copia local de recuperación antes de importar, y escritura solo del snapshot del dispositivo propio. No hay envío de la biblioteca solo por conectar la cuenta.
- **Después de confirmar la vinculación:** escritorio publica cambios propios en lotes y verifica novedades como máximo cada 90 segundos durante las comprobaciones activas, sin trasladar cada edición por separado.
- **Alternativa local:** botón «Detectar OneDrive» busca la ubicación sincronizada en Windows/macOS y evita seleccionar entre varias cuentas ambiguas. Las carpetas manuales quedan como opción avanzada.

## Compatibilidad y restricciones

- **Android:** por el momento sigue con Beta 13.5, usando la misma carpeta OneDrive desde Microsoft Graph. Su sincronización es manual para respetar las mejoras de batería.
- No se migra silenciosamente ninguna carpeta o configuración antigua. Si la ruta ya configurada es diferente, se debe exportar una copia JSON y revisar qué carpeta contiene las revisiones más recientes antes de cambiarla.
- **Windows Beta 14 Preview:** la sesión de Microsoft se conserva cifrada con **DPAPI (usuario actual de Windows)** en `%APPDATA%\\Selector_Cine_365\\onedrive_sesion_dpapi.bin`. Al reabrir o cambiar el EXE, recupera la autorización sin abrir el navegador. Los tokens de acceso siguen solo en RAM; el refresh token se renueva y se reescribe protegido. Desconectar elimina el archivo protegido. Es necesario autorizar **una última vez** al instalar esta corrección porque las versiones anteriores no conservaban la credencial.
- **Mac:** la sesión sigue siendo temporal y hay que volver a autorizar tras reiniciar; no se almacenan credenciales sin cifrar. Keychain está pendiente. Windows recupera únicamente sesiones del mismo usuario de Windows; si Microsoft revoca el permiso o cambia la cuenta, habrá que autenticarse otra vez.
- En Microsoft Entra → Registro de aplicaciones → Autenticación, hay que añadir **`http://localhost`** como URI de redirección para escritorio (plataforma aplicación móvil/escritorio). El puerto local cambia y Microsoft permite la coincidencia de `localhost` sin puerto. **No borres** la URI Android `msal696458fe-0672-47e5-9fdd-aa0332e9884b://auth`. La autenticación de escritorio aún no se ha probado con una cuenta Microsoft real.
- Las compilaciones de prueba de Windows y macOS no deben confundirse con instaladores finales firmados o notarizados.
- No se deben subir claves privadas, contraseñas, tokens Microsoft o bibliotecas personales a GitHub.

## Pruebas de regresión propuestas

1. Sin conectar la nube, crear/modificar una película y comprobar persistencia únicamente local.
2. En Windows o Mac, pulsar Conectar Microsoft, completar la autorización PKCE en el navegador y verificar que vuelve a la aplicación sin selector.
3. Antes de la primera sincronización, verificar que no se crean snapshots ni se reemplazan listas.
4. En Android Beta 13.5, pulsar ↻ y confirmar; en escritorio revisar y aceptar el mismo conjunto de listas.
5. Marcar y desmarcar en dos dispositivos desconectados; comprobar que aparece conflicto y que se mantiene la última decisión validada.
6. Simular OneDrive sin conexión, error 429, paginación sospechosa y snapshots corruptos; no deben sustituir datos locales.
7. En Windows, conectar una sola vez, cerrar la app y abrir otro EXE Beta 14; debe conservarse la misma cuenta y la misma vinculación sin navegador. Desconectar debe eliminar el archivo DPAPI e impedir restaurarlo. En Mac se sigue necesitando reconectar.

## Estado de publicación

Se publicará una release Beta 14 solo cuando la CI de Windows/macOS sea satisfactoria, el acceso real a OneDrive esté validado y la migración con Android se haya probado. Mientras tanto se trabaja en una rama y una PR separadas de `main`.

## Correcciones de interfaz / exportación para esta revisión de Beta 14

- El selector «SORTEAR ENTRE» se repuebla con las plataformas de la caché TMDB al abrir; no hace falta recargar TMDB para ver Netflix, Filmin u otros proveedores.
- Barra de scroll de «Mis listas» con carril y extremos dentro de las esquinas redondeadas del panel en WebView2.
- Exportar TXT usa el nombre real de la lista activa con caracteres válidos para ficheros; en vista de comparación, `Comparación de listas.txt`.
- Compilaciones CI de Windows y macOS Intel/Apple Silicon verificadas; Windows sigue siendo una vista previa experimental hasta validar la autenticación Microsoft Graph con cuenta real.

## Sesión Microsoft permanente en Windows (corrección de Beta 14)

- DPAPI protege la credencial de renovación con el usuario actual de Windows; **no** se copia ningún token sin cifrar al navegador, a la nube ni al repositorio.
- El archivo cifrado se guarda junto a las listas y al identificador del dispositivo, fuera de la carpeta de OneDrive. La sincronización de películas y la autorización de Microsoft son datos independientes.
- El EXE nuevo recupera automáticamente el refresh token. La primera petición Graph renueva el acceso de forma silenciosa; no se crea una cuenta nueva ni se vuelve a pedir confirmación de las listas previamente vinculadas.
- Si no hay conexión, se mantiene la credencial para el siguiente intento y se conservan todos los cambios locales. Solo una revocación explícita de Microsoft requiere volver a iniciar sesión.
- El almacenamiento seguro es exclusivo de Windows por ahora. El EXE nunca incluye ni distribuye las credenciales.
- Las pruebas incluyen fusión y versiones existentes, reinicios simulados, renovación del token, desconexión, cortes de red, revocación y DPAPI nativo en Windows CI.
