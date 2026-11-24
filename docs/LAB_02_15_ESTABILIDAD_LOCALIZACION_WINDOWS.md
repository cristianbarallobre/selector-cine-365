# Selector Cine 365 · Windows LAB-02.15

**Rama:** `experimental`. No se publican instaladores ni releases estables. Android y macOS no cambian.

## 1. Biografías de intérpretes y fechas

- Las fichas de TMDB solicitan primero `es-ES`. Si la biografía no está disponible, se intenta `es-MX`.
- Si tampoco existe una biografía en español, se muestra un aviso en castellano y un botón opcional **«Ver biografía original en inglés»**. El texto inglés no se muestra como si fuese español ni se traduce artificialmente.
- Las fechas ISO se conservan en datos y ordenaciones pero **se muestran** como `11 de octubre de 1962` (o `11/10/1962` en vista breve). Las fechas de sincronización ya empleaban `es-ES`.
- Se traduce «Acting» como «Interpretación» y otros nombres de departamento. Nombres propios y lugares que devuelve TMDB no se traducen de forma automática.
- No se incorpora un traductor de terceros ni se envían biografías a servicios externos no configurados.

## 2. Integridad y copias

Los archivos permanentes continúan en `%APPDATA%\Selector_Cine_365_LAB\` y no se mezclan con los de la versión estable.

- **Copias diarias:** `copias_automaticas\`, máximo 60 archivos diarios con nomenclatura de la aplicación (las copias personales con otros nombres nunca se eliminan).
- **Copias rápidas:** `copias_rapidas_lab\`, máximo 24 snapshots, separadas al menos 15 minutos. Solo se respaldan estados válidos de la biblioteca, antes de sobrescribirlos.
- **Recuperación de OneDrive y sincronización:** `sync_recuperacion\`, mantiene el mecanismo anterior, no se borra con la limpieza de caché.
- **Archivo dañado en el arranque:** Windows LAB busca la última copia válida en los historiales locales. Si existe, conserva el original con sufijo `.corrupt_FECHA.json` y recupera el estado válido.
- **Sin copia válida:** se interrumpe el inicio con un error claro, sin sobrescribir el archivo original ni sustituirlo silenciosamente por la biblioteca inicial.
- **Daño detectado durante una sesión:** la escritura rechaza sobrescribir ese archivo y solicita reiniciar para recuperación.

Una recuperación puede devolver la biblioteca a una fecha anterior. Comprueba las listas antes de volver a sincronizar. No es una garantía de recuperación si el equipo pierde también las copias.

## 3. Historial y vista previa

El panel de sincronización ya dispone de **Actividad reciente**, **Copias locales de recuperación** y confirmación antes de restaurar. LAB incorpora las copias rápidas al selector y muestra una vista previa validada, con cantidad de listas, películas vistas y Favoritos. Ninguna copia se restaura por seleccionarla.

La restauración guarda otra copia local del estado previo e interviene en la sincronización normal. Antes de restaurar deben terminarse las operaciones de guardado.

## 4. Mantenimiento

En el panel de sincronización aparece **Mantenimiento de Windows LAB**, con:

- Tamaño de biblioteca, caché de fichas TMDB, plataformas y carteles; número de copias locales.
- Estado de conexión TMDB, OneDrive o carpeta de nube, y aviso de conflictos pendientes.
- Rutas de datos permanentes y caché temporal.
- Revisión de imágenes antiguas y botón confirmado para **limpiar únicamente carteles descargados**, que después pueden volver a descargarse.
- Acceso a las releases estables de GitHub, sin intentar actualizar automáticamente una versión de laboratorio.

**No se eliminan** listas, puntuaciones, corazones, identificaciones manuales, claves API, sesión de OneDrive ni el perfil WebView2. La caché WebView2 no está incluida en el presupuesto de carteles, al estar gestionada por el navegador.

## 5. Pruebas y límites

Se ejecutan las regresiones previas de navegación, filtros, búsqueda (40 listas y 10.000 películas), sincronización, Favoritos y descubrimiento; también:

- Biografías es-ES/es-MX e inglés solo bajo acción explícita; fechas en formato español.
- Arranque con JSON dañado sin sobrescribirlo; recuperación desde snapshot válido y aislamiento del original.
- Retención acotada de copias; no borrar copias personales.
- Validación de copias antes de mostrar una restauración.
- Autenticación del panel de mantenimiento y limpieza únicamente de imágenes autorizadas.

**Pendiente de una fase de distribución:** firma del EXE y creación de instalador. No se simula una firma sin certificado. **Pendiente de pruebas reales:** cierres forzados, red inestable, semanas de uso y revisión de la recuperación antes de enviar cambios a otros dispositivos. Primero probar Windows LAB; después portar los cambios aprobados a Android y Mac.
