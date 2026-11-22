# Política de versiones, descargas y recuperación

**Revisión vigente: 22 de noviembre de 2026.**

## Versión recomendada

**[Selector Cine 365 v9.0.0 ESTABLE](https://github.com/cristianbarallobre/selector-cine-365/releases/tag/v9.0.0)** es la única versión marcada como **Latest** para nuevos usuarios. Esta release usa **sin recompilar ni modificar** los cuatro ejecutables comprobados en Windows Beta 14, Mac Beta 14 (Intel/Apple Silicon) y Android Beta 13.13 Limpia.

Aunque la etiqueta general es `v9.0.0`, los binarios conservan «Beta 14» / «Beta 13.13» como versión interna. Los cuatro SHA-256 son los de la [release fuente de referencia](https://github.com/cristianbarallobre/selector-cine-365/releases/tag/v9.0.0-beta.14). El archivo `SHA256SUMS.txt` documenta las huellas. No renombrar binarios como si hubieran sido recompilados.

**Actualización Android:** paquete estable `es.cine365.app.limpia`, código 9027, firma original. Las futuras actualizaciones deben conservar **el mismo ID de paquete + firma de certificado**, aumentando versionCode. La línea histórica `es.cine365.app` permanece incompatible por identidad y utiliza almacenamiento privado propio. No distribuir firmados con certificados temporales como actualización sobre la línea estable.

## Política de retención

| Publicación | Función | Estado recomendado |
| --- | --- | --- |
| `v9.0.0` | Versión estable multiplataforma | **Última versión (Latest)** |
| `v9.0.0-beta.14` | Referencia exacta de los binarios ahora estables | Conservar como prerelease histórica |
| `v9.0.0-beta.13.3` | Última APK histórica de la línea `es.cine365.app` | Conservar para recuperación del paquete antiguo |
| `v9.0.0-beta.13.1` | Última combinación histórica completa anterior a Beta 14 | Conservar para recuperación |
| `v9.0.0-beta.9` | Histórico de la primera generación Android/Windows | Conservar por incompatibilidades de identidad |
| `v8.2` | Histórica pre-v9 (Windows y Mac) | Conservar |
| `v9.0.0-beta.12`, `beta.12.1`, `beta.12.2`, `beta.13.2`, `v8.1.1` | Betas/versión anterior redundantes | Archivar instaladores, manifiesto y notas **antes** de retirar sus releases públicas |

Las publicaciones históricas ya archivadas en el borrador privado **`archivo-historico-beta-v9-2026-11`** no deben borrarse sin crear previamente una copia adicional. Para la limpieza de noviembre se conserva un segundo respaldo verificable con los archivos de las releases retiradas. Las etiquetas Git y el historial de commits **no se eliminan**, aunque desaparezca la página pública de una release redundante.

## Seguridad de la publicación

1. No recompilar binarios de una versión que se declara idéntica a la ya probada. Copiar sus activos y verificar SHA-256.
2. Antes de publicar, verificar el ZIP del código fuente y comprobar la firma Android original (v2/v3, huella del certificado).
3. Publicar solo una release marcada **Latest**; las betas anteriores permanecen como prerelease o archivos históricos.
4. Generar y conservar manual PDF/Word, notas, historial, manifiesto `SHA256SUMS.txt` y ZIP del código fuente.
5. No publicar jamás `.p12`, `.jks`, contraseñas, tokens Microsoft, claves de TMDB ni bibliotecas personales.
6. Mantener copia JSON fuera de OneDrive antes de cambios de versión o sincronización; no prometer recuperación automática de copias mezcladas.
7. Los EXE Windows no están firmados con Authenticode y los ZIP Mac no están notarizados. Documentar estas limitaciones con claridad, sin llamarlos «instaladores firmados».

## Descargas e historial

- [Última versión estable](https://github.com/cristianbarallobre/selector-cine-365/releases/latest)
- [Lista de releases](https://github.com/cristianbarallobre/selector-cine-365/releases)
- [Código fuente estable](https://github.com/cristianbarallobre/selector-cine-365/tree/v9.0.0)
- [Historial de cambios](../CHANGELOG.md) y [guía de sincronización](GUIA_SINCRONIZACION_SEGURA.md)

**Recuperación:** si un archivo histórico ya no aparece en Releases, consultar los borradores de archivo privados con permisos de mantenimiento y verificar la huella SHA-256 antes de reutilizarlo. Las etiquetas originales siguen accesibles en Git.

## Tratamiento de datos de versiones anteriores

La versión estable incluye barreras contra la mezcla de listas con el mismo ID. **No reparan automáticamente una colección que ya estuviera mezclada**. Evitar acciones destructivas de sincronización hasta confirmar que cada lista contiene los títulos correctos. Conservar `colecciones.json` y los snapshots anteriores, especialmente cuando los recuentos discrepan.
