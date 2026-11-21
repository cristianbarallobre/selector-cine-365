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
- **Los tokens de escritorio solo se guardan en memoria durante la sesión.** Hay que volver a autorizar al reiniciar la aplicación; no se guardan credenciales en texto plano. El almacenamiento seguro de sesión con Keychain/credenciales de Windows queda pendiente.
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
7. Cerrar sesión, reiniciar, comprobar que la copia local permanece y que hace falta reconectar Microsoft.

## Estado de publicación

Se publicará una release Beta 14 solo cuando la CI de Windows/macOS sea satisfactoria, el acceso real a OneDrive esté validado y la migración con Android se haya probado. Mientras tanto se trabaja en una rama y una PR separadas de `main`.

## Correcciones de interfaz / exportación para esta revisión de Beta 14

- El selector «SORTEAR ENTRE» se repuebla con las plataformas de la caché TMDB al abrir; no hace falta recargar TMDB para ver Netflix, Filmin u otros proveedores.
- Barra de scroll de «Mis listas» con carril y extremos dentro de las esquinas redondeadas del panel en WebView2.
- Exportar TXT usa el nombre real de la lista activa con caracteres válidos para ficheros; en vista de comparación, `Comparación de listas.txt`.
- Compilaciones CI de Windows y macOS Intel/Apple Silicon verificadas; Windows sigue siendo una vista previa experimental hasta validar la autenticación Microsoft Graph con cuenta real.
