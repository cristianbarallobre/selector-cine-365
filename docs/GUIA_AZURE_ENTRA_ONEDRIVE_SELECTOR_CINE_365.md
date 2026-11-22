# Selector Cine 365 — Configuración de Microsoft Entra y OneDrive

**Guía vigente para v9.0.0 estable (22 de noviembre de 2026).** Esta distribución conserva Windows/macOS Beta 14 y Android Beta 13.13 como nombres internos. Este documento sustituye la guía Beta 13.1; la edición anterior puede recuperarse desde el historial Git.

## 1. ¿Tengo que configurar OneDrive?

**No.** Selector Cine 365 funciona localmente en todos los dispositivos sin configurar Microsoft ni TMDB. OneDrive es opcional y sirve para **compartir bibliotecas y marcas de vistas** entre Windows, Mac y Android.

Si quieres sincronizar, usa **la misma cuenta Microsoft personal** en todos los dispositivos. Las películas, las listas y el estado visto se almacenan primero en la memoria local del programa. Las copias de OneDrive son adicionales.

## 2. Registro Microsoft Entra que utiliza el proyecto

| Campo | Valor comprobado en este proyecto |
| --- | --- |
| Nombre de aplicación | Selector Cine 365 |
| Application (client) ID público | `696458fe-0672-47e5-9fdd-aa0332e9884b` |
| Tipos de cuenta | Cuentas personales de Microsoft |
| Tipo de aplicación | Cliente público, móvil y escritorio |
| Flujo OAuth | Authorization Code + PKCE (S256) |
| Redirección Android | `msal696458fe-0672-47e5-9fdd-aa0332e9884b://auth` |
| Redirección Windows/macOS | `http://localhost` |
| Permisos Graph | `Files.ReadWrite`, `User.Read`, `offline_access` |

El ID del cliente es **público**, no una contraseña. **No crees ni copies secretos de cliente** dentro del proyecto. No compartas contraseñas Microsoft, tokens, claves de TMDB ni certificados privados de firma Android.

### Si el registro ya existe y funciona

**No crees otro.** Solo comprueba en Microsoft Entra → **Registros de aplicaciones → Selector Cine 365 → Autenticación** que estén las dos redirecciones apropiadas; no borres la URI Android al configurar la de escritorio. El registro debe permitir el flujo de cliente público.

### Si necesitas revisar o reproducir la configuración

1. Accede a [Microsoft Azure Portal](https://portal.azure.com/) con la cuenta autorizada para administrar el registro.
2. Abre **Microsoft Entra ID → Registros de aplicaciones** y localiza **Selector Cine 365**.
3. En **Información general**, comprueba el Application (client) ID anterior.
4. En **Autenticación**, revisa que se permita iniciar sesión con cuentas Microsoft personales y que el cliente tenga las redirecciones Android y `http://localhost` como aplicación móvil/escritorio.
5. En **Permisos de API**, verifica los permisos delegados requeridos de Microsoft Graph y la posibilidad de acceso offline para renovar autorizaciones.
6. Guarda los cambios y autoriza desde las aplicaciones, nunca introduciendo contraseñas Microsoft dentro de un formulario de la app.

No es necesario contratar servidores, crear una máquina virtual o un registro distinto por cada dispositivo.

## 3. Conectar Windows v9.0.0 (compilación Beta 14)

1. Abre Selector Cine 365 y entra en **Mis listas → Sincronización y copias de seguridad**.
2. Selecciona **Conectar Microsoft**. La app abre el navegador oficial para completar PKCE.
3. La app utiliza un puerto local efímero y recibe la autorización; no copies códigos de acceso en chats.
4. Se localiza/crea la carpeta `Selector_Cine_365` de la cuenta y se prepara la comparación.
5. **Antes de sincronizar**, revisa las listas del PC y las copias remotas. No confirmes si alguna lista incluye un número inesperado de películas.

La credencial de renovación se almacena cifrada con **DPAPI del usuario de Windows** en `%APPDATA%\Selector_Cine_365\onedrive_sesion_dpapi.bin`. Esto permite continuar la sesión al cerrar la aplicación o sustituir el EXE. Una autorización anterior a DPAPI requiere iniciar sesión una vez. Si Microsoft revoca el permiso o cambias de cuenta, será necesario autorizar de nuevo.

Para comprobar datos dudosos selecciona **Solo cuando pulse sincronizar** y no permitas que una copia posiblemente mezclada se publique en otros equipos.

## 4. Conectar Android v9.0.0 (compilación Beta 13.13 Limpia)

1. Comprueba que tienes instalada `es.cine365.app.limpia` (versionCode 9027).
2. En **Nube**, elige **Conectar cuenta Microsoft** e inicia sesión con la misma cuenta que en Windows.
3. La app localiza o crea `Selector_Cine_365` en OneDrive; **conectar no debe importar ni publicar automáticamente las listas**.
4. Pulsa el botón de sincronización manual **↻**, espera la vista previa y revisa **el nombre y la cantidad de películas de cada lista**.
5. Solo si el resultado es correcto resuelve posibles conflictos de vistas y elige **Confirmar y sincronizar**.

Android **no sincroniza las ediciones automáticamente en segundo plano**. Marcar películas se guarda en el móvil; enviar/recibir requiere la acción manual con confirmación.

**Importante para nuevos dispositivos:** la lista inicial de muestra con 365 películas se sustituye por la biblioteca real de la nube únicamente si está intacta. Esto evita añadir todas las películas de muestra a una lista sincronizada.

## 5. Conectar Mac (Intel y Apple Silicon)

Mac comparte la lógica de interfaz y las copias de escritorio de Beta 14. Sin embargo, la versión actual **no guarda todavía la sesión Microsoft con Keychain**, por lo que puede solicitar una nueva autorización tras reiniciar. No significa que se haya borrado tu biblioteca: los datos locales se conservan en `~/Library/Application Support/Selector_Cine_365`.

Las versiones para Mac son archivos ZIP que incluyen `.app`; están distribuidas sin firma/notarización de Apple.

## 6. Qué datos se envían (y cuáles no)

| Se comparte mediante OneDrive | Permanece local |
| --- | --- |
| Listas y contenido | Clave personal de TMDB |
| Marcas vistas/pendientes | Carátulas descargadas de TMDB |
| IDs TMDB corregidos manualmente cuando están en `tmdb_overrides` | Caché de plataformas y valoraciones |
| Registros de cambios, eliminaciones y snapshots de sincronización | Contraseñas, sesiones Microsoft protegidas y claves de firma |
| Identificadores de listas y dispositivos | Archivos de recuperación no publicados localmente |

OneDrive mantiene los archivos bajo la carpeta lógica `Selector_Cine_365`, con snapshots en la estructura `Selector_Cine_365_Copias/sincronizacion/dispositivos/` cuando corresponde al modo de sincronización.

**Una carátula visible en Windows no aparece instantáneamente en Android**: cada dispositivo la consulta por separado a TMDB. Este comportamiento no se arregla volviendo a conectar la cuenta Microsoft.

## 7. Recuperación y prevención de conflictos

La versión estable bloquea mezclas sospechosas de listas que comparten un ID pero contienen películas incompatibles. **No puede separar automáticamente una lista que ya fue mezclada en versiones anteriores**.

1. Exporta un JSON de la biblioteca de Windows si allí las listas siguen correctamente separadas.
2. Guarda ese JSON **fuera** de la carpeta de OneDrive.
3. Compara nombres y recuentos de cada lista con la vista previa de Android.
4. Si aparece contenido inesperado, **Cancelar**; no borres snapshots ni pulses sincronizar.
5. Conserva la copia local de cada dispositivo y los respaldos `sync_recuperacion`. La recuperación debe basarse en una biblioteca buena verificada.

## 8. Solución de problemas frecuentes

- **No aparece «Conectar Microsoft» o falla el retorno:** comprueba la versión, `http://localhost` para escritorio y la URI `msal...` para Android. Revisa que el navegador no esté bloqueando la devolución a la app.
- **Se vuelve a solicitar la autorización en Windows:** verifica que usas el mismo usuario de Windows y que Microsoft no haya revocado los permisos. DPAPI protege la credencial pero no impide una revocación.
- **Mac solicita entrar cada vez:** limitación conocida hasta implementar almacenamiento en Keychain.
- **Android muestra «Conflicto con un paquete»:** el identificador estable es `es.cine365.app.limpia`, diferente al histórico `es.cine365.app`. No desinstales sin una copia JSON.
- **Las listas se mezclan o los conteos cambian mucho:** cancela la sincronización y recupera una biblioteca intacta antes de intentar otra vez.
- **No aparecen carátulas:** configura TMDB en cada dispositivo; OneDrive no lleva esas imágenes.

[Manual completo v9.0.0](MANUAL_COMPLETO_SELECTOR_CINE_365.md) · [Guía de sincronización segura](GUIA_SINCRONIZACION_SEGURA.md) · [Release estable](https://github.com/cristianbarallobre/selector-cine-365/releases/tag/v9.0.0).
