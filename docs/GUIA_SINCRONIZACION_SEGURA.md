# Guardado y recuperación de Selector Cine 365

Windows y macOS 9.2.30.

## Guardado principal y copia opcional

La biblioteca se guarda automáticamente en el equipo, sin conectar cuentas. En **Mis listas → Guardado y copia adicional**, elige una carpeta y pulsa **Guardar ubicación**. Cada instalación conserva su propia copia en `Selector_Cine_365_Copias`, con `colecciones.json`, `colecciones_anterior.json` y hasta 14 copias diarias sujetas a un límite de espacio.

La copia contiene listas, estados de visionado, favoritos, valoraciones, identificaciones e historial. No incluye claves, sesiones ni cachés. **Copiar ahora** solicita una copia. Si la carpeta no está disponible, el guardado principal continúa y la app reintenta. Desactivar la copia conserva los archivos existentes.

## Trasladar una biblioteca

Exporta una copia JSON completa desde **Mis listas → Copia de seguridad**. En el equipo de destino, exporta primero su biblioteca actual y usa **Restaurar copia**. No se fusionan bibliotecas ni se aplican cambios automáticamente. Los servicios externos administran sus carpetas, no la aplicación.

## Resolver fallos de guardado

Si un guardado no se confirma, reintenta o exporta un JSON de emergencia. Se intenta conservar un borrador local recuperable. En **Mis listas → Mantenimiento → Exportar registro de diagnóstico** puedes guardar el registro técnico.

## Historial de actividad

Consulta hasta 10.000 movimientos confirmados en **Mis listas → Guardado y copia adicional → Actividad reciente, historial y copias de recuperación**. Hay búsqueda, filtros y páginas de 20. El historial se guarda por separado con huella de integridad y espejo de recuperación. Se incluye en las exportaciones completas y la copia adicional.

## Volver a una copia anterior

1. Exporta la biblioteca actual antes de sustituirla.
2. Elige un JSON completo o revisa las copias locales en el apartado anterior.
3. Revisa las cantidades y confirma. Se guarda una copia previa. Si el JSON contiene historial, se restaura; si no, se conserva el actual. Comprueba tus listas.

[Manual completo](https://github.com/cristianbarallobre/selector-cine-365/blob/main/docs/MANUAL_COMPLETO_SELECTOR_CINE_365.md) · [Rutas de datos](https://github.com/cristianbarallobre/selector-cine-365/blob/main/docs/UBICACIONES_DATOS_SELECTOR_CINE_365.md)

Autor: Cristian Barallobre Rúa.

## Salvaguardas de la copia adicional y del historial


Cambiar solo de lista activa conserva la selección en el equipo, pero no reescribe la copia adicional cuando su contenido e historial son iguales. La fecha de la copia corresponde al archivo guardado. Si la copia actual está dañada, se conserva su contenido en un archivo `colecciones_danada_<huella>.json` antes de reconstruirla desde la biblioteca principal. Si no se puede conservar el original, no se sustituye. Los archivos dañados conservados quedan disponibles para revisión manual.


Si la biblioteca es válida pero su historial no está disponible, se buscan copias completas que contengan exactamente el historial comprometido, verificado por su huella SHA-256. Si no se encuentra ninguna, no se descarta historial ni se sustituye automáticamente la biblioteca; conserva los archivos y revisa una copia completa antes de intervenir.
