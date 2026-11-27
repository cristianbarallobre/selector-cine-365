# Selector Cine 365 v9.0.0 — versión estable


## Qué se distribuye exactamente

Esta es una **promoción de las compilaciones ya probadas**, no una recompilación. El número de lanzamiento en GitHub es **v9.0.0**, mientras que los nombres internos de los programas conservan:

| Plataforma | Compilación incluida en v9.0.0 | Archivo |
| --- | --- | --- |
| Windows 10/11 x64 | Beta 14 | `Selector_Cine_365_Windows_Beta14.exe` |
| macOS Apple Silicon | Beta 14 | `Selector_Cine_365_Mac_AppleSilicon_Beta14.zip` |
| macOS Intel | Beta 14 | `Selector_Cine_365_Mac_Intel_Beta14.zip` |


## Mejoras incluidas en esta referencia

- Biblioteca con listas independientes, importación/exportación TXT y copia de seguridad JSON.
- Protección de sincronización: identifica posibles colisiones de ID entre listas diferentes y detiene uniones que parecen mezclar películas de distintas colecciones.
- Revisión de listas por nombre y recuento de películas antes de confirmar una combinación; resolución de conflictos de marcas.

## Requisitos y actualización

- **Windows:** sustituir el EXE sin borrar `%APPDATA%\Selector_Cine_365`.
- **Mac:** escoger ZIP Intel o Apple Silicon y conservar `~/Library/Application Support/Selector_Cine_365`.
- Los datos pueden usarse **sin OneDrive ni TMDB**. Antes de una actualización o revisión de sincronización exporta JSON y guárdalo fuera de la carpeta de OneDrive.

## Advertencia sobre datos anteriores

La versión estable evita nuevas fusiones de listas sospechosas, pero **no deshace mezclas ya guardadas por una beta anterior**. Si los recuentos por lista son incorrectos, cancela la sincronización y recupera una copia verificada. No borres copias ni snapshots de otros dispositivos para forzar una restauración.


## Documentación incluida

- Manual completo actualizado en **PDF, Word y Markdown**.
- Guía de OneDrive y sincronización segura en **PDF y Word**.
- Notas de versión, historial de cambios y `SHA256SUMS.txt`.
- Código fuente de la etiqueta `v9.0.0` (ZIP y descargas GitHub automáticas).

## Estado de mantenimiento

`v9.0.0` es el **canal estable recomendado** por el autor tras sus pruebas. «Estable» aquí significa referencia para uso y soporte, **no garantía de ausencia total de errores**. Se conserva la beta anterior como referencia exacta; otras versiones históricas se archivan antes de retirar publicaciones redundantes, conservando sus etiquetas Git.

**Importante:** no se incluyen contraseñas, archivos de firma privada, tokens de Microsoft ni claves TMDB en la release.
