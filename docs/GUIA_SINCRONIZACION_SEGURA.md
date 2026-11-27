# Guardado y recuperación de Selector Cine 365

Windows y macOS 9.2.26.

## Guardado principal y copia opcional

La biblioteca de escritorio se guarda automáticamente en el equipo. No necesita conectar una cuenta ni consultar archivos remotos. En **Mis listas → Guardado y copia adicional**, elige una carpeta y pulsa **Guardar ubicación**. Puedes elegir un disco externo o una carpeta gestionada por otro programa. La aplicación crea una carpeta propia por instalación, con `colecciones.json` y, tras una modificación, `colecciones_anterior.json`.

El archivo contiene listas, estados de visionado, favoritos, valoraciones, identificaciones de películas e historial de actividad. Además se conservan hasta 14 copias diarias sujetas a un límite de espacio. No incluye sesiones, claves ni cachés. **Copiar ahora** solicita una copia del estado principal. Si la carpeta no está disponible, el guardado local sigue funcionando y la app reintenta. Desactivar la copia conserva los archivos existentes.

## Trasladar una biblioteca

Exporta una copia JSON completa desde **Mis listas → Copia de seguridad**. En el equipo de destino, exporta primero su biblioteca actual y usa **Restaurar copia** para seleccionar el JSON de origen. No se fusionan bibliotecas ni se aplican cambios automáticamente. La copia de un equipo no sustituye la copia de otro.

## Resolver fallos de guardado

Si el guardado local no está confirmado, reintenta o exporta un JSON de emergencia. Desde 9.2.20 se intenta conservar un borrador recuperable de las ediciones antes del guardado principal. En Mis listas → Mantenimiento → Exportar registro de diagnóstico puedes exportar el registro técnico local.

## Historial de actividad

Consulta las modificaciones confirmadas desde **Mis listas → Guardado y copia adicional → Actividad reciente, historial y copias de recuperación**. Hay búsqueda, filtros y páginas de 20 movimientos, con un máximo de 10.000 registros locales. El historial se guarda por separado con una huella de integridad y un espejo de recuperación. Las exportaciones completas y la copia adicional incluyen sus movimientos; no incluyen claves, cachés ni registros de diagnóstico.

## Volver a una copia anterior

1. Exporta la biblioteca actual antes de sustituirla.
2. Elige un JSON verificado o consulta **Mis listas → Guardado y copia adicional → Actividad reciente, historial y copias de recuperación**.
3. Revisa las cantidades y confirma la restauración. Se conserva una copia local anterior con historial. Si el JSON seleccionado contiene historial, lo restaura; si es una copia antigua sin historial, conserva el actual. Comprueba tus listas y movimientos.

La aplicación no conecta cuentas ni borra archivos de carpetas externas. Utiliza la exportación manual para trasladar bibliotecas entre equipos.

[Manual completo](https://github.com/cristianbarallobre/selector-cine-365/blob/main/docs/MANUAL_COMPLETO_SELECTOR_CINE_365.md) · [Rutas de datos](https://github.com/cristianbarallobre/selector-cine-365/blob/main/docs/UBICACIONES_DATOS_SELECTOR_CINE_365.md)

Autor: Cristian Barallobre Rúa.
