# Selector Cine 365 v9.0.0 — versión estable

**Publicado el 22 de noviembre de 2026.** Es la versión estable de referencia tras las pruebas del usuario en **Windows, Android y dispositivos macOS**.

## Qué se distribuye exactamente

Esta es una **promoción de las compilaciones ya probadas**, no una recompilación. El número de lanzamiento en GitHub es **v9.0.0**, mientras que los nombres internos de los programas conservan:

| Plataforma | Compilación incluida en v9.0.0 | Archivo |
| --- | --- | --- |
| Windows 10/11 x64 | Beta 14 | `Selector_Cine_365_Windows_Beta14.exe` |
| macOS Apple Silicon | Beta 14 | `Selector_Cine_365_Mac_AppleSilicon_Beta14.zip` |
| macOS Intel | Beta 14 | `Selector_Cine_365_Mac_Intel_Beta14.zip` |
| Android 8+ | Beta 13.13 Limpia · versionCode 9027 | `Selector_Cine_365_Android_Beta13_13_Limpia.apk` |

Los cuatro binarios y su **SHA-256 coinciden con la última Beta 14 coordinada**. No se han reemplazado ni vuelto a firmar. Windows es EXE portátil, no un MSI; las aplicaciones Mac no están firmadas/notarizadas por Apple; la APK Android sí está firmada con el certificado original.

## Mejoras incluidas en esta referencia

- Biblioteca con listas independientes, importación/exportación TXT y copia de seguridad JSON.
- Navegación Android adaptada a zonas seguras y pantalla con notch; búsqueda global y «Cargar más».
- Protección de sincronización: identifica posibles colisiones de ID entre listas diferentes y detiene uniones que parecen mezclar películas de distintas colecciones.
- Primera sincronización de un Android recién instalado: evita concatenar automáticamente el catálogo inicial de 365 películas a la biblioteca real procedente de OneDrive, cuando el catálogo inicial está intacto.
- Revisión de listas por nombre y recuento de películas antes de confirmar una combinación; resolución de conflictos de marcas.
- Windows conserva la sesión Microsoft cifrada con DPAPI; Android guarda credenciales localmente con Android Keystore. Mac todavía reautoriza tras reiniciar.
- Consulta TMDB y relleno gradual de carátulas visibles en Android cuando la clave está configurada, con descarga local independiente por dispositivo.

## Requisitos y actualización

- **Android:** identificador **`es.cine365.app.limpia`**, `versionCode=9027`. **Actualiza Beta 13.12 Limpia** sin desinstalar siempre que conserve la misma firma. **No actualiza** la antigua línea `es.cine365.app` ni APK de depuración con otra firma.
- **Windows:** sustituir el EXE sin borrar `%APPDATA%\Selector_Cine_365`.
- **Mac:** escoger ZIP Intel o Apple Silicon y conservar `~/Library/Application Support/Selector_Cine_365`.
- Los datos pueden usarse **sin OneDrive ni TMDB**. Antes de una actualización o revisión de sincronización exporta JSON y guárdalo fuera de la carpeta de OneDrive.

## Advertencia sobre datos anteriores

La versión estable evita nuevas fusiones de listas sospechosas, pero **no deshace mezclas ya guardadas por una beta anterior**. Si los recuentos por lista son incorrectos, cancela la sincronización y recupera una copia verificada. No borres copias ni snapshots de otros dispositivos para forzar una restauración.

Las carátulas y plataformas TMDB **no se copian de Windows a Android**: cada dispositivo consulta/cacha sus fichas, aunque los IDs de TMDB corregidos manualmente pueden formar parte de la biblioteca sincronizada. La clave TMDB no sale del dispositivo.

## Documentación incluida

- Manual completo actualizado en **PDF, Word y Markdown**.
- Guía de OneDrive y sincronización segura en **PDF y Word**.
- Notas de versión, historial de cambios y `SHA256SUMS.txt`.
- Código fuente de la etiqueta `v9.0.0` (ZIP y descargas GitHub automáticas).

## Estado de mantenimiento

`v9.0.0` es el **canal estable recomendado** por el autor tras sus pruebas. «Estable» aquí significa referencia para uso y soporte, **no garantía de ausencia total de errores**. Se conserva la beta anterior como referencia exacta; otras versiones históricas se archivan antes de retirar publicaciones redundantes, conservando sus etiquetas Git.

**Importante:** no se incluyen contraseñas, archivos de firma privada, tokens de Microsoft ni claves TMDB en la release.
