# Selector Cine 365 — v9 Beta 13.1 (Windows y macOS)


## Corrección principal de Beta 13.1: TMDB en escritorio

- Windows/Mac Beta 13 podían quedarse en «Conexión local TMDB no disponible. Reintentando…» después de cerrar y volver a abrir la aplicación.
- **Causa comprobada:** el frontend solicitaba `GET /api/platforms/cache`, que el módulo de TMDB admitía, pero que no estaba registrado en el router HTTP principal. Respondía **404** aun estando la clave API guardada en AppData.
- **Solución:** ruta registrada, lectura de estado separada de la caché y botón Actualizar capaz de recuperar una consulta fallida. Si la caché falla no se borra la clave ni se indica erróneamente que no está configurada.
- Se añadieron pruebas del servidor (estado/caché responden 200 con token; 403 sin token). El ejecutable Windows se compiló como aplicación GUI, sin consola. Los Mac se reconstruyeron desde el mismo código corregido.

## Funciones de Beta 13 conservadas (todas las plataformas)

- «Mis listas» ocupa siempre el espacio de cuatro filas en Windows/Mac aunque la última página tenga menos listas.
- Las eliminaciones llevan registros `deleted_lists` compartidos, para que copias antiguas de otros dispositivos **no resuciten listas borradas**.
- Las identificaciones manuales con el lápiz guardan IDs de TMDB en `tmdb_overrides`, que pueden sincronizarse. **La clave API personal jamás se incluye en la biblioteca compartida**.
- Persistencia local de listas, marcas de visto y preferencias, funcionamiento sin nube y sincronización causal de las marcas.

## Descargas de esta release

| Plataforma | Archivo |
|---|---|
| Windows 10/11 x64 | `Selector_Cine_365_Windows_v9_BETA13_1.zip` |
| Mac Apple Silicon | `Selector_Cine_365_Mac_Apple_Silicon_v9_BETA13_1.zip` |
| Mac Intel | `Selector_Cine_365_Mac_Intel_v9_BETA13_1.zip` |

Incluye además **Manual completo PDF**, **Guía Microsoft Entra/OneDrive PDF** e **Historial completo v9 Beta 13.1 en PDF y Word**.

## Instalación y datos


**Sincronización:** actualiza todos los dispositivos vinculados antes de volver a borrar listas. Los borrados realizados antes de Beta 13 no incluían marcadores: si una lista histórica reaparece, bórrala una vez con ambos dispositivos ya actualizados. La sincronización por OneDrive no es instantánea; pueden surgir conflictos concurrentes que requieran revisión.


**Mac:** firma ad hoc, **sin notarización Apple Developer ID**. Es posible que Gatekeeper solicite autorizar su apertura. Compilación CI verificada, pero no sustituye pruebas manuales en hardware físico.

## Pruebas y trazabilidad

- [Windows Beta 13.1: compilación y pruebas de caché HTTP](https://github.com/cristianbarallobre/selector-cine-365/actions/runs/44829805619).
- [macOS Beta 13.1 Intel/Apple Silicon](https://github.com/cristianbarallobre/selector-cine-365/actions/runs/44830855769).


