# Selector Cine 365 — v9 Beta 13 (sincronización y persistencia)

Actualización correctiva de **Windows y macOS**. Se ha realizado a partir del código v9 Beta 12.2, **no desde la antigua v8**.

## Cambios

1. **Mis listas de Windows/Mac:** se reservan siempre **cuatro posiciones por página**, aunque queden una, dos o tres listas. La nube y los demás controles ya no se desplazan al cambiar de página.
2. **Borrados sincronizados:** al eliminar una lista se guarda en la biblioteca compartida un registro de eliminación (\`deleted_lists\`). Windows y Mac lo respetan y no vuelven a crear esa lista a partir de copias antiguas. Las copias locales de recuperación se conservan.
3. **Identificación manual TMDB compartida:** al identificar con el lápiz, se almacena el ID confirmado (\`tmdb_overrides\`) en el JSON compartido. Cuando el otro dispositivo recibe la biblioteca, aplica esa identificación en su caché local sin exigir rehacerla con el lápiz. **La clave API personal de TMDB NO se sincroniza**; debe configurarse por separado en cada dispositivo.
4. **Reinicio de la app:** las claves TMDB permanecen en su fichero privado local y un fallo de comunicación de la interfaz al iniciar ya no se interpreta como pérdida de configuración. Hay reintentos y el motor conserva la clave guardada.
6. **Windows/Mac:** ejecutables actualizados desde el mismo código Go + HTML, conservando la sincronización causal de marcas Vista/Pendiente y los controles de estadísticas.

## Qué hacer antes de actualizar

**Importante: realiza una copia JSON de la biblioteca en ambos dispositivos.** Guárdala fuera de la carpeta OneDrive de sincronización. Recomendamos actualizar **todos los dispositivos vinculados** a Beta 13 antes de volver a probar borrados entre dispositivos.

- **Windows:** cierra la aplicación y sustituye el ejecutable; no abras dos versiones a la vez.
- **Mac:** sustituye la aplicación en Aplicaciones; la firma es ad hoc, no está notarizada por Apple.

**Para borrados efectuados ANTES de Beta 13:** las versiones antiguas no registraban las eliminaciones. Si una lista ya había reaparecido por una copia vieja, elimínala **una sola vez después de instalar Beta 13 en ambos dispositivos** y comprueba que desaparezca también del otro. No borres listas antes de asegurarte de que ambos dispositivos han recibido la actualización.

**TMDB:** si tu cuenta no tiene clave TMDB en uno de los dispositivos, los IDs manuales pueden recibirse con las listas, pero las consultas de póster/valoración/plataformas necesitan configurar tu clave personal localmente. No compartas esa clave ni capturas que la muestren.

## Comprobaciones

- Persistencia de configuración TMDB al reiniciar el motor Go.
- Compilaciones Windows y Mac basadas en el código de Beta 13.

**Limitación:** las pruebas automatizadas no sustituyen una campaña real de sincronización OneDrive con dos dispositivos conectados y cambios simultáneos. La fusión conserva una eliminación confirmada incluso si una versión vieja publica una copia antigua; al recuperar intencionadamente una lista borrada de una copia, crea una lista nueva con un ID nuevo. Las modificaciones simultáneas de IDs TMDB contradictorios no se resuelven automáticamente: revisa el caso a mano.

© 2026 Cristian Barallobre Rúa.
