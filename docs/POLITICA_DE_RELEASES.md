# Política de publicaciones y recuperación

**Selector Cine 365 · Actualizada el 20 de noviembre de 2026**

## Releases públicas y versiones históricas

| Publicación | Estado | Activos |
|---|---|---|
| [v9 Beta 13.1](https://github.com/cristianbarallobre/selector-cine-365/releases/tag/v9.0.0-beta.13.1) | **Release coordinada actual** | Windows/Mac Beta 13.1, Android Beta 13 firmada, manuales PDF e historial Word/PDF |
| [v9 Beta 12.2](https://github.com/cristianbarallobre/selector-cine-365/releases/tag/v9.0.0-beta.12.2) | **Beta más reciente**, Android y escritorio sin cambios | APK Android 9014, Windows/Mac Beta 12 y documentación |
| [v9 Beta 12.1](https://github.com/cristianbarallobre/selector-cine-365/releases/tag/v9.0.0-beta.12.1) | Histórica: primer arreglo notch | APK Android 9013 y escritorio original |
| [v9 Beta 12](https://github.com/cristianbarallobre/selector-cine-365/releases/tag/v9.0.0-beta.12) | Beta para pruebas | Windows, Android, Mac Apple Silicon, Mac Intel y dos manuales PDF |
| [v9 Beta 9](https://github.com/cristianbarallobre/selector-cine-365/releases/tag/v9.0.0-beta.9) | Histórica | APK y ZIP de Windows originales |
| [v8.2](https://github.com/cristianbarallobre/selector-cine-365/releases/tag/v8.2) | Histórica | Windows y dos variantes de Mac |
| [v8.1.1](https://github.com/cristianbarallobre/selector-cine-365/releases/tag/v8.1.1) | Histórica | Paquete de la edición 8.1.1 |

Las **17 releases redundantes** se retiraron solamente después de verificar que los ejecutables y las APK originales estaban respaldados.

## Archivo histórico privado

En **GitHub → Releases → Drafts**, los mantenedores pueden acceder al borrador **Archivo privado de betas históricas v9**, identificado como **archivo-historico-beta-v9-2026-11**.

Contiene un archivo llamado **Archivo_Instaladores_Betas_Anteriores_v9.zip**, de unos 24 MB, que incluye:

- Todos los APK y ZIP originales de las 17 releases retiradas, organizados por el nombre de su tag original.
- Un **MANIFIESTO.json** con el tamaño, la huella SHA-256 y la referencia de commit.
- Un archivo **LEEME.txt** con instrucciones.

El proceso comprobó la integridad de los binarios originales, generó el archivo, lo subió a GitHub y **lo descargó de nuevo para confirmar su SHA-256 antes de borrar las releases públicas redundantes**.

Este borrador es visible solamente para quienes tienen permisos de mantenimiento. **No lo borres sin crear otra copia independiente**; es recomendable guardarlo también fuera de GitHub.

## Código fuente y etiquetas conservados

**No se han eliminado ramas, etiquetas Git ni el historial de commits.** Se comprobaron las referencias de los 17 tags históricos después de la limpieza.

El código fuente puede recuperarse desde las ramas y los tags. No obstante, **volver a compilar un APK no garantiza recuperar la misma firma original**. Por eso se conservaron también los binarios históricos en el borrador.

La nueva release combinada de **Beta 9** incluye exactamente los APK/ZIP originales de sus dos releases anteriores, verificados con SHA-256.

## Cómo recuperar un instalador antiguo

1. Abre GitHub con permisos de mantenimiento y entra en **Releases → Drafts → Archivo privado de betas históricas v9**.
2. Descarga el archivo ZIP de instaladores históricos.
3. Dentro del ZIP, localiza la carpeta **originales / nombre-del-tag / nombre-del-archivo**.
4. Comprueba el SHA-256 del instalador con **MANIFIESTO.json**.
5. Si necesitas republicarlo, utiliza la etiqueta Git original y el commit registrado en el manifiesto.

Para Android, **Beta 9 y Beta 12 no comparten identidad ni firma**: requiere migración JSON. **Beta 11 y Beta 12** usan el mismo identificador estable y requieren el mismo certificado privado de firma.

## Regla para nuevas versiones

1. Publicar **una única release por versión** con Windows, Android y ambas arquitecturas de macOS, si todas superan las verificaciones.
2. Incluir notas de cambios, manuales PDF y comprobación de integridad de las descargas.
3. En Android conservar el identificador **es.cine365.app**, aumentar el código de versión y firmar siempre con el mismo certificado. Nunca subir la clave privada ni sus contraseñas al repositorio.
4. No sustituir datos de usuario ni afirmar que OneDrive sincroniza instantáneamente; conservar el sistema de revisiones causales y resolución de conflictos.
5. Antes de retirar otras releases, realizar y verificar una copia de sus instaladores y conservar los tags.
6. Mantener versiones anteriores con valor histórico o de recuperación, así como los manuales editables en la carpeta docs.

## Registro de mantenimiento

- **20-11-2026:** consolidadas 20 releases originales en cuatro visibles: Beta 12, Beta 9 combinada, v8.2 y v8.1.1.
- Archivados los 17 conjuntos de activos antiguos en un borrador de GitHub con hashes SHA-256.
- Los PDF actualizados de Manual completo y Guía Microsoft Entra/OneDrive están disponibles en el repositorio y en la release Beta 12.
- El script utilizado se conserva en la rama **maintenance/organizar-releases-v9-2026-11**, dentro de scripts/maintenance.

## Consulta documental

Los cambios de software de las etapas v1 a v9 Beta 12 están reunidos en [Historial completo v9 Beta 12](HISTORIAL_V9_BETA12.md), disponible también en Word y PDF. Se conserva por separado la edición histórica v8.1.1 sin reescribirla.

## Publicación coordinada 13.1 (20-11-2026)

Se publica una única release `v9.0.0-beta.13.1` con **Windows Beta 13.1**, **Mac Intel y Apple Silicon Beta 13.1** y **la APK Android Beta 13 original (código 9015)**. Android no incrementa versión porque la reparación de TMDB era exclusiva del servidor Go de escritorio. Se conservan los tags y las publicaciones previas; ninguna de las antiguas APK ni sus firmas se sobrescribe.
