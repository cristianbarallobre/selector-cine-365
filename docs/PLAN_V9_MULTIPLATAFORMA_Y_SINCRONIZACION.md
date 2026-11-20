# Selector Cine 365 — Plan funcional y técnico v9.0

**Estado (20-11-2026):** **v9 Beta 13 Android y v9 Beta 13.1 Windows/macOS; implementación experimental en Windows y Android; la versión para Mac Intel y Apple Silicon comparte el código de Windows y está publicada como beta en GitHub Releases.** No está aprobada todavía como versión estable: faltan pruebas prolongadas de sincronización en dispositivos reales.  
**Autor:** Cristian Barallobre Rúa.  
**Plataformas previstas:** Windows, Android y macOS (Apple Silicon / Intel).  
**Base actual:** código v9 Beta 13.1 de escritorio Go/HTML/WKWebView-WebView2 y código Android Java/Gradle dentro de este mismo repositorio. El plan recoge objetivos originales y puede incluir funciones todavía pendientes.

## Seguimiento del estado v9 Beta 13.1 (Windows/Mac) y Beta 13 (Android)

**Ya implementado en la beta (sujeto a pruebas reales):** modo local sin nube; Windows, Android y versiones Mac Intel/Apple Silicon publicadas; base compartida para Windows y Mac; búsqueda global/listas numeradas; estado Vista/Pendiente desde el tick; PKCE/OneDrive personal en Android; revisiones causales, detección periódica, historial y copias de recuperación; interfaz Android adaptada.

**Implementado en Beta 13:** borrados de listas con tombstones persistentes, ID de TMDB corregidos manualmente compartidos sin compartir API key, cuatro huecos estables en paginación de listas. **Beta 13.1 de escritorio:** ruta local `/api/platforms/cache` reparada tras detectar HTTP 404 al abrir.

**Pendiente o con límites importantes:** pruebas de uso real en Mac Intel y Apple Silicon (la compilación CI ya pasó); notarización y firma Apple Developer ID; automatización segura de firma Android para futuras publicaciones; comprobación física prolongada de propagación y conflictos de OneDrive; manejo de borrados/renombrados simultáneos; revisión exhaustiva de accesibilidad móvil.

**Compatibilidad:** Android Beta 9 y Beta 13 no comparten identificador ni certificado, por lo que requieren migración única mediante JSON. Android Beta 11, 12 y 13 comparten identificador permanente y deben firmarse con el mismo certificado. La versión estable v9.0 todavía **no** se ha publicado.

## 1. Objetivos y principios innegociables

1. Las aplicaciones tendrán el nombre visible **Selector Cine 365** y el título destacado **¿Qué película veo hoy?**.
2. La firma de autor **Hecho por Cristian Barallobre Rúa** seguirá apareciendo de forma discreta.
3. El icono definitivo será **original**, creado para este proyecto a partir de una claqueta de cine con marco azul y acentos dorados. No se afirmará que el icono actual sea de creación original sin verificar sus fuentes.
4. Android utilizará controles táctiles y su propia composición; escritorio conservará su interfaz optimizada para ventanas.
5. El almacenamiento local seguirá funcionando sin conexión y **sin clave TMDB**. Carteles, datos enriquecidos e identificación de TMDB podrán requerir una clave personal.
6. El primer dispositivo **no se elegirá por sistema operativo**. Cualquier instalación podrá ser la primera o incorporarse a una biblioteca existente.
7. Conectar un dispositivo nuevo **nunca publicará automáticamente una biblioteca vacía o inicial** encima de otra con cambios de usuario.
8. La sincronización será reversible, tendrá copias de recuperación y nunca tendrá como única fuente de verdad el último archivo `colecciones.json` que llegue.
9. No se incluirán credenciales de OneDrive ni claves TMDB en los archivos compartidos.
10. v8.2 y las versiones previas se mantendrán como recuperación hasta superar pruebas reales de v9.0.

## 1.1 Funcionamiento local sin nube — requisito obligatorio

**El guardado local es el modo predeterminado y completo en Windows, Android y macOS.** Ningún usuario necesita configurar OneDrive, crear una cuenta Microsoft ni activar sincronización para usar la biblioteca.

- Al primer arranque, la aplicación crea sus archivos de datos privados en el propio dispositivo y guarda inmediatamente las ediciones posteriores.
- Sin carpeta vinculada, no se intenta publicar, sustituir ni leer ninguna biblioteca remota; la pantalla Nube mostrará «No configurada» e informará de que el guardado local funciona.
- La sincronización solo se inicia por decisión expresa del usuario, con revisión de la biblioteca remota y confirmación antes de publicar cambios.
- «Desconectar nube» **no elimina** las listas ni las películas guardadas localmente. La aplicación sigue funcionando sin conexión.
- Si más adelante se conecta la nube, las modificaciones locales previas se conservan y se incorporan mediante fusión y resolución de conflictos; jamás se reemplaza automáticamente la biblioteca local con un catálogo vacío o una copia remota por defecto.
- El servicio TMDB es independiente de la nube: una clave personal de TMDB sirve para fichas, carteles e identificación, pero no es necesaria para gestionar y guardar las películas.
- Android debe permitir exportar una copia de seguridad local antes de desinstalar o borrar datos, ya que el almacenamiento privado de la aplicación puede borrarse al desinstalarla.
- Las pruebas de aceptación incluirán arrancar desde cero, crear y editar listas, cerrar y reabrir, trabajar sin Internet, desactivar nube y reconectarla después sin pérdida.

## 2. Identidad y marca

- Nombre en icono, barra de título, pantallas, diálogos y documentación: **Selector Cine 365**.
- Encabezado principal: **¿Qué película veo hoy?**
- Crédito: **Hecho por Cristian Barallobre Rúa**.
- Nuevos recursos del mismo original: `icon.png`, `icon.ico`, recurso Windows `.syso`, `.icns` en macOS e iconos adaptativos de Android.
- El logotipo se elaborará como una ilustración original, no como una mera copia pixelada del icono existente.
- **Compatibilidad:** no renombrar de golpe `%APPDATA%\\AtlasDelCine`, `~/Library/Application Support/AtlasDelCine`, `atlas` ni las claves de JSON: son identificadores internos históricos. La marca visible cambiará; los datos antiguos deberán seguir cargando. Un traslado de rutas solo podrá hacerse con migración, copia y reversión comprobadas.

## 3. Android: pantallas y comportamiento exactos

### 3.1 Estructura de pantalla y desplazamiento

**Zona superior fija**, por encima de la línea señalada en las capturas:
- Zona segura del sistema operativo (barra de estado).
- Icono oficial, encabezado, indicador `Guardado / Sin guardar / Sincronizando` y crédito.
- Selector de colección activa (si procede a la vista).
- Separador visual inferior de la cabecera.

**Zona central desplazable (vertical)**:
- Solo el contenido de la pestaña seleccionada.
- Debe conservar la posición de desplazamiento por pestaña cuando el usuario cambia de sección.
- Los menús contextuales y la edición no quedarán ocultos detrás del teclado o la barra inferior.

**Zona inferior fija**:
- Cinco destinos: **Biblioteca | Sorteo | Estadísticas | Listas | Nube**.
- Etiquetas visibles, iconos con contraste y superficies táctiles de ~48 dp como mínimo.
- Soporte de teléfonos estrechos, tamaños de letra del sistema, navegación con gestos y teclado en pantalla.

### 3.2 Biblioteca

- Selector de colección; si es una colección virtual «Coincidentes», se identifica como solo lectura.
- Buscador con búsqueda por título original, traducido y año.
- Filtros grandes **Todas / Pendientes / Vistas**.
- Botón **Plataformas** que abre un panel inferior con filtros en vez de ocupar permanentemente la pantalla.
- Tarjetas verticales de película: póster, título, año, estado visto, etiquetas de disponibilidad, nota TMDB y botón de ficha.
- Cambio de visto con control táctil grande y confirmación visual de guardado.
- Lista vertical de desplazamiento continuo y carga por lotes, en vez de paginador minúsculo.
- Marcado «Vista» propagado a otras listas solo después de identificar con fiabilidad la misma película.

### 3.3 Sorteo

- Pantalla propia con elección de colección y filtro de plataformas.
- Botón principal grande **Elegir película**.
- Resultado con cartel, título, año, descripción, disponibilidad y estado de visto.
- Opciones «Volver a sortear» y «Marcar vista» fáciles de pulsar.
- No incluir por defecto películas ya vistas.

### 3.4 Estadísticas

- Contadores: vistas, pendientes, porcentaje, tiempo cuando existan duraciones fiables.
- Décadas, géneros, países y directores en paneles legibles; no presentar datos TMDB inexistentes como cero real.
- Cambiar de colección sin volver a Biblioteca.
- Gráficos/páginas adaptados a desplazamiento táctil.

### 3.5 Listas

- Una tarjeta por colección: nombre, número de películas, vistas y estado activo.
- Al tocar una tarjeta: **Abrir / Añadir películas / Importar TXT / Exportar / Renombrar / Eliminar**.
- Comparación de al menos dos listas, con colección virtual no editable.
- Importación de líneas con o sin `☐`/`☑`; el año es opcional, pero ayuda a identificar versiones.
- La eliminación debe solicitar confirmación. Las eliminaciones sincronizadas necesitan evento específico, no se deducen simplemente por «ausencia».

### 3.6 Nube

- Estado de conexión: **sin configurar / OneDrive conectado / última sincronización / errores**.
- Botón **Conectar OneDrive** o **Elegir carpeta compatible**.
- Primera sincronización: **Examinar biblioteca encontrada → Comparar → Importar/Combinar → Activar sincronización**.
- «Sincronizar ahora», «Ver historial», «Revisar conflictos», «Exportar copia de seguridad» y «Desconectar».
- Mostrar cantidad de listas, películas y cambios que entrarán/saldrán **antes** de la primera publicación.
- Si no hay conexión, guardar cambios localmente y mostrar «Cambios pendientes de sincronizar».

## 4. Sincronización sin orden obligatorio

### 4.1 Compatibilidad con la carpeta de OneDrive ya utilizada

Mantener y **no sobrescribir** estos archivos existentes:

```text
OneDrive/
└── Selector_Cine_365_Copias/
    ├── actual/
    │   └── colecciones.json          # copia Windows heredada
    ├── historico/
    │   └── AAAA-MM-DD/               # puntos de restauración heredados
    └── sincronizacion/
        └── v2/
            ├── manifest.json         # identificador y esquema de biblioteca
            ├── dispositivos/
            │   ├── <id-dispositivo-1>/
            │   │   └── eventos/
            │   └── <id-dispositivo-2>/
            │       └── eventos/
            ├── instantaneas/
            └── recuperacion/
```

Los archivos de `sincronizacion/v2` son una **estructura prevista**, no la estructura de la beta 8.3. Mantener soporte de lectura del formato antiguo `sincronizacion/dispositivos` durante una migración explícita.

### 4.2 Primera conexión: decisiones según el contenido, no según el equipo

| Estado local | Estado compartido | Acción |
|---|---|---|
| Instalación nueva con catálogo de fábrica | Biblioteca compartida existente | Importar primero. Nunca publicar los estados por defecto. |
| Biblioteca personal existente | Carpeta de nube nueva | Crear biblioteca compartida tras resumen y confirmación. |
| Biblioteca personal existente | Biblioteca compartida existente | Previsualizar y fusionar sin borrar; si hay conflictos, preguntar. |
| Dos instalaciones nuevas conectadas a la vez | Carpeta nueva | Identificador único, creación condicional y relectura; no pisarse. |
| Carpeta antigua con `actual/colecciones.json` | Sin protocolo nuevo | Importar como «origen heredado», conservar copia y generar respaldo antes de migrar. |
| Carpeta inaccesible o sin permisos | Desconocido | No crear una falsa copia local ni indicar éxito; avisar y permitir reintentar. |

La aplicación solo considerará «cambio del usuario» una acción real y registrada, no un `☐` que venía preinstalado.

### 4.3 Modelo de datos compartido

- `library_id`: UUID común a la biblioteca sincronizada.
- `device_id`: UUID generado por instalación, persistente, independiente de «Windows», «Android» o «Mac».
- Listas con `list_id` estable y nombre editable; título y año no bastan para decidir que dos listas distintas son la misma.
- Película lógica: **ID de TMDB confirmado**; sin ID, título normalizado + año cuando sea inequívoco. Nunca fusionar automáticamente remakes diferentes.
- Eventos inmutables con `event_id`, `device_id`, secuencia local, `list_id`, identidad de película, operación, valor y versiones conocidas de otros dispositivos.
- Eventos: `mark_seen`, `mark_pending`, `add_movie`, `remove_movie`, `create_list`, `rename_list`, `delete_list`, `override_tmdb_id`.
- Manifest e instantáneas con versión de esquema y comprobaciones de integridad; un archivo de índice no debe ser el único lugar donde se registran cambios.

### 4.4 Resolución de cambios

- **Sin conflicto:** importar cambios nuevos una sola vez; los eventos se deduplican por `event_id`.
- **Cambio causal posterior:** prevalece el posterior cuando un dispositivo ya había visto el cambio anterior.
- **Cambios simultáneos desconectados sobre la misma película con decisiones opuestas:** generar conflicto visible; no decidir únicamente por la hora del reloj.
- **Alta de una película frente a eliminación concurrente:** preservar los datos y ofrecer resolución.
- **Eliminaciones:** conservar registros de eliminación (tombstones) y permitir restaurar desde historial.
- **Sin identidad TMDB confiable:** no propagar cambios entre películas ambiguas.
- Cada fusión produce un resumen de `añadidas / marcadas vistas / marcadas pendientes / conflictos` y **se guarda localmente antes de subir nuevos cambios**.
- Las escrituras serán idempotentes y comprobadas; cuando el proveedor admita escritura condicional (ETag), aprovecharla para evitar sustituciones accidentales.

### 4.5 Copias de seguridad y permisos

- Copia local antes de migraciones y antes de fusiones masivas.
- Conservar `actual/` e `historico/` existentes como respaldo independiente; no utilizarlos como el canal de escritura simultánea entre dispositivos.
- No compartir API Key de TMDB, tokens de Microsoft, rutas locales de usuario ni cachés de carteles.
- No considerar una sincronización «realizada» hasta confirmar lectura, escritura y relectura de la información esperada.
- Modo desconectado: cambios pendientes, reintento seguro y estado comprensible.

## 5. OneDrive en Android: acceso real

Windows y macOS pueden trabajar con la carpeta local sincronizada de OneDrive. En Android:

1. Intentar selección de una carpeta con **Storage Access Framework** y comprobar **lectura y escritura reales** del proveedor seleccionado.
2. Si OneDrive no expone una carpeta seleccionable y escribible, **no simular que está sincronizado**: ofrecer acceso oficial con **inicio de sesión Microsoft (MSAL) y Microsoft Graph**.
3. La autenticación de Android utilizará un cliente público con permisos delegados y autorización explícita. Nunca llevar secreto de aplicación incrustado.
4. No pedir al usuario copiar manualmente archivos como única vía de sincronización permanente.
5. Mantener importación/exportación manual como alternativa cuando el acceso de OneDrive no esté disponible.

**Requisito pendiente de implantación:** registrar la aplicación Microsoft necesaria para MSAL, configurar sus URI de redirección/permisos y probar cuentas personales de OneDrive. Ninguna beta actual ha demostrado todavía la sincronización continua real en un Xiaomi Mi 9T.

## 6. Archivos: mapa concreto de implementación

### Go compartido (Windows y macOS)

| Archivo | Trabajo previsto |
|---|---|
| `main.go` | Registrar endpoints de sincronización y estado, conservando arranque local actual. |
| `collections.go` | Persistencia segura y migraciones de colecciones, sin cambiar estructuras antiguas a ciegas. |
| `cloud.go` | Mantener copias `actual` e `historico`; conectar con adaptador del nuevo protocolo. |
| `identity.go` y `providers.go` | Reutilizar IDs TMDB confirmados; mantener clave personal local. |
| `sync_protocol.go` (nuevo) | Esquemas JSON versionados, validación, UUID, eventos. |
| `sync_merge.go` (nuevo) | Fusión causal, detección de conflictos y deduplicación. |
| `sync_migrate.go` (nuevo) | Importación de la carpeta OneDrive heredada. |
| `sync_store.go` (nuevo) | Intercambio seguro con carpeta local y registro de cambios. |
| `sync_*_test.go` (nuevos) | Pruebas de concurrencia, importación y no pérdida de datos. |
| `window_windows.go`, `window_darwin.go`, `macos_view_darwin.m` | Conservar motores nativos de ventana ya funcionales. |
| `plantilla.html` | Marca visible, estado de sincronización, asistente de conexión y conflictos. |

### Android (proyecto `Cine365_Android_Studio`)

| Archivo | Trabajo previsto |
|---|---|
| `app/src/main/assets/index.html` | Cabecera fija + contenido desplazable + barra inferior fija; vistas táctiles. |
| `MainActivity.java` | Control WebView, selectores de documentos, estado de red y flujo de autenticación. |
| `Backend.java` | API local, persistencia, importación y operaciones de listas. |
| `SyncStore.java` | Sustituir la fusión beta basada solo en instantáneas por eventos causales versionados. |
| `CloudBackup.java` | Adaptador SAF, comprobación de permisos y distinción copia/sincronización. |
| `OneDriveAuth.java` (nuevo) | Autenticación MSAL sin secretos incrustados. |
| `OneDriveRepository.java` (nuevo) | Lectura/escritura de archivos y reintentos con Microsoft Graph. |
| `SyncEngine.java` (nuevo) | Fusiones, deduplicación, conflictos y cambios sin conexión. |
| `SyncMigration.java` (nuevo) | Detectar biblioteca de fábrica, archivos heredados y primer dispositivo. |
| `Tmdb.java` | Mantener clave local y correlación de identidades. |
| `AndroidManifest.xml`, `app/build.gradle` | Permisos mínimos, configuración MSAL, versiones y firma de distribución. |
| `res/` | Icono adaptativo oficial y recursos Android. |
| `app/src/test/` (nuevo) | Pruebas de reglas de sincronización y migraciones. |

### Empaquetado y documentación

- `icon.png`, `icon.ico`, recursos `.syso`, icono macOS `.icns` y recursos Android generados desde **un único original maestro**.
- `README.md`, `CHANGELOG.md`, `AUTORIA.md`, `MACOS.md`: actualizar instrucciones y atribuciones.
- `.github/workflows/`: pruebas y construcción por plataforma; **beta separada de las Releases estables** hasta validación.
- `docs/`: manual de conexión y recuperación, riesgos y migración de v8.1.1/v8.2/v8.3 beta.

## 7. Pruebas de aceptación obligatorias

1. Instalar Android nuevo con el catálogo sin marcar y enlazar OneDrive con 87 vistos del PC: aparecen 87 vistos en Android y Windows conserva 87.
2. Empezar **solo** en Android y añadir 10 vistas; instalar Windows después: Windows recibe las 10, sin cambiarlas a pendientes.
3. Marcar una película como vista en Windows; importar en Android; desmarcarla Android; recibir en Windows: el estado final se conserva en ambos.
4. Dos dispositivos cambian una misma película **sin conexión** de manera contraria: se señala el conflicto, no se pierden datos.
5. OneDrive desconectado, permisos caducados, archivos incompletos o sincronización retrasada: no se declara éxito ni se borra la biblioteca.
6. Múltiples listas, duplicados con títulos traducidos y IDs TMDB: estados coherentes en todos los dispositivos.
7. Dos remakes de título idéntico y años distintos: no se mezclan.
8. Importación de `Selector_Cine_365_Copias/actual/colecciones.json` antigua: conservada íntegramente y con copia restaurable.
9. Una instalación Android con datos anteriores a la conexión: se conserva su biblioteca y se propone combinación, no reemplazo.
10. Cambios de nombre y eliminación de listas, restauración desde histórico y compatibilidad con backup legacy.
11. En Xiaomi Mi 9T: apertura, navegación táctil, teclado, selección real de OneDrive, modo sin conexión, vuelta de pantalla y permisos persistentes.
12. En Windows y Mac: app nativa, copia y sincronización, sin pérdida en actualizaciones.
13. El cambio de marca visible no interrumpe el acceso a archivos guardados históricamente bajo `AtlasDelCine`.
14. Compilaciones Android APK, Windows EXE y Mac .app verificadas; nunca publicar una build que no se haya probado suficientemente como estable.

## 8. Entregas por etapas

**v9.0-beta.1 — Identidad y Android:** icono original, nueva marca, cabecera fija, pestañas y controles táctiles, migración local sin nube.

**v9.0-beta.2 — Sincronización segura:** protocolo v2, migración de backup Windows, detección de biblioteca existente, vista previa, registro de eventos y conflictos; sin publicación estable.

**v9.0-beta.3 — OneDrive real:** acceso SAF cuando sea posible, MSAL/Graph de respaldo, pruebas con Xiaomi y Windows, tareas de reconexión.

**v9.0 candidata — Multiplataforma:** integración macOS, paquetes firmados en la medida disponible, ensayo de restauración, documentación y pruebas cruzadas.

**v9.0 estable:** solo tras verificar los escenarios anteriores y conservar copia y posibilidad de volver a la versión estable anterior.

## 9. Estado de implementación y referencia histórica

**Situación actual — 20 de noviembre de 2026:** la última beta coordinada publicada es **v9 Beta 12** para Windows, Android y macOS Intel/Apple Silicon. El motor Go y la interfaz de escritorio son compartidos entre Windows y Mac; Android mantiene su proyecto Java/Gradle y una interfaz táctil diferenciada.

**Implementado en Beta 12, pendiente de validación prolongada en hardware real:** funcionamiento local, TMDB opcional, listas numeradas, búsqueda global, manejo de pertenencia, sincronización causal y copias de seguridad, historial, OneDrive personal con PKCE en Android y aplicaciones Mac empaquetadas.

**Limitaciones actuales:** sincronización OneDrive no instantánea, posibles conflictos concurrentes que requieren revisión, pruebas reales de Mac y Android pendientes, Mac sin notarización Developer ID, clave privada Android conservada fuera del repositorio y necesidad de migrar JSON desde Android Beta 9.

**Referencia histórica:** cuando se redactó el primer plan, **v8.2** era la última edición estable de escritorio y se ensayaban prototipos v8.3. Esos datos ya no describen la implementación actual; se conservan como contexto de diseño, no como instrucciones de instalación o diagnóstico.

La próxima **v9 estable** debe aprobarse solo después de superar los casos de prueba funcionales, sincronización, restauración y actualización sin pérdida de datos descritos en esta especificación. Consulta el [historial reciente](../CHANGELOG.md) y el [historial detallado de v9](HISTORIAL_V9_BETA12.md).
