# Selector Cine 365 · Android v9 beta 4

Esta versión es **experimental**. Conserva tu aplicación anterior y exporta una copia JSON de las listas antes de probar. La APK de esta beta tiene un identificador distinto, por lo que no accede automáticamente a la biblioteca ni a la caché de TMDB de otras betas.

## Novedades

- Estadísticas completas: sin flechas, todas las filas visibles y scroll vertical.
- Sexta pestaña «Ayuda» con instrucciones de TMDB y funciones básicas.
- Datos de TMDB guardados sin caducidad automática. Al abrir se lee la caché local: solo películas nuevas sin ficha o actualizaciones manuales generan solicitudes. Los carteles también se guardan en el móvil.
- Sincronización experimental con la versión **Windows v9 beta 4**, con vista previa, conflictos y archivos independientes por dispositivo.
- Sincronización **opcional**: la biblioteca se guarda en el móvil incluso sin nube ni clave API.

## Instalación y seguridad

- Instalar `Selector_Cine_365_Android_v9_BETA4.apk`. Nombre de paquete de prueba: `es.cine365.app.v9beta4`.
- La Beta 3 no se reemplaza; para conservar películas vistas, exporta «Copia de seguridad» JSON desde la anterior y restáurala en esta.
- El token TMDB se configura en cada instalación. No se sube a la nube.
- La caché persiste entre aperturas **de esta misma instalación**, pero no se comparte automáticamente entre aplicaciones con identificadores distintos.
- No desinstales la versión anterior hasta comprobar todas tus listas.

## Nube: no depende de OneDrive

**Windows** permite seleccionar cualquier carpeta local sincronizada por un servicio externo (OneDrive, Google Drive, Dropbox u otro). **Android** usa el selector de carpetas del sistema (Storage Access Framework, SAF), que solo funciona si el proveedor elegido permite seleccionar carpetas y ofrece **lectura y escritura reales**. Que una nube aparezca en «Archivos» no garantiza esa capacidad.

1. Elige una carpeta compartida del proveedor en Windows y asegúrate de que se sincronice en el disco local.
2. Android → **Nube → Elegir carpeta**: selecciona **la misma carpeta compartida** o directamente `Selector_Cine_365_Copias`.
3. Android → **Activar copia**: el dispositivo aún **no publica** datos compartidos hasta que revises la primera combinación.
4. **Revisar sincronización**: muestra listas, vistas y conflictos; confirma solamente después de comprobarlos.
5. En Windows v9 beta 4 usa **Revisar sincronización con otros dispositivos**, resuelve posibles conflictos y confirma. A partir de entonces cada equipo publica su propio archivo de dispositivo al guardar.
6. Para recibir cambios nuevos del otro dispositivo en esta beta, vuelve a usar **Revisar sincronización**. **No es todavía sincronización entrante automática**.

Si tu aplicación de nube en Android no permite seleccionar una carpeta con lectura/escritura, **no es compatible con este acceso directo por SAF**. Como alternativa, una aplicación de sincronización de carpetas de terceros puede mantener una carpeta local del móvil sincronizada con esa nube y puedes seleccionarla desde Android. Debes comprobar que admita sincronización bidireccional y que no borre conflictos sin avisarte. También puedes exportar/importar JSON manualmente.

### Protección de datos

- `actual/` y `historico/` siguen siendo copias de recuperación de Windows.
- `sincronizacion/dispositivos/<id>.json` contiene la copia específica de cada instalación: **cada dispositivo escribe solo su archivo**.
- Se hace una copia local de recuperación antes de importar; no se aplican cambios contradictorios sin elección.
- No se sincronizan tokens TMDB, cachés de imágenes ni contraseñas.
- Las eliminaciones y renombrados simultáneos de listas, y la sincronización automática de entrada, **todavía requieren más pruebas y protección**. No uses esta beta como único respaldo.

## Comprobaciones recomendadas

1. Importar tu copia JSON de la beta anterior; marcar una película en Android, revisar en Windows e importar allí.
2. Desmarcarla en Windows, revisar de nuevo en Android y comprobar que cambió.
3. Cambiar de estado la misma película con ambos dispositivos desconectados y comprobar el aviso de conflicto.
4. Desactivar la nube y confirmar que los datos se siguen guardando localmente.
5. Cerrar y reabrir la aplicación: comprobar que se muestran carteles y fichas guardados sin volver a descargar todo.

**Autor:** Cristian Barallobre Rúa. Este producto utiliza TMDB, pero no está respaldado ni certificado por TMDB.
