# Selector Cine 365 — v9 Beta 14 / Android Beta 13.13

**21 de noviembre de 2026 · Beta recomendada para pruebas, NO versión estable definitiva.**

Esta publicación reúne los cuatro instaladores verificados en una sola GitHub Release y el código fuente combinado. Es la referencia de pruebas que reemplaza a las betas anteriores para nuevas instalaciones; las anteriores siguen disponibles para recuperación y consulta histórica.

## Versiones incluidas

| Plataforma | Versión del programa | Paquete |
| --- | --- | --- |
| Windows x64 (10/11) | Beta 14 | EXE portátil (sin instalador MSI ni firma Authenticode) |
| macOS Intel (12+) | Beta 14 | ZIP con aplicación .app, sin notarización Apple |
| macOS Apple Silicon (12+) | Beta 14 | ZIP con aplicación .app, sin notarización Apple |
| Android (8+) | Beta 13.13, código 9027 | APK firmada con el certificado original, paquete `es.cine365.app.limpia` |

**Android:** esta APK **actualiza** Beta 13.12 **Limpia** con el mismo ID y certificado, sin desinstalar. **No actualiza** la antigua línea `es.cine365.app` ni APK firmadas temporalmente. Una instalación nueva de la línea Limpia tiene almacenamiento independiente.

## Correcciones estabilizadas

- **Barrera contra mezclas de listas:** Android y escritorio bloquean la fusión si dos listas reutilizan un ID pero son incompatibles. También bloquean ampliaciones desproporcionadas de copias antiguas ya mezcladas, sin escribir en la nube ante el error.
- **Primera sincronización Android:** una instalación nueva sin cambios adopta la biblioteca real de OneDrive como referencia; no añade el catálogo inicial de 365 películas a la lista de Windows.
- **Verificación por lista:** la vista previa Android enseña nombre y cantidad de películas de cada lista propuesta, con aviso de cancelar ante diferencias.
- **Guardado local y nube bajo control:** Android guarda localmente y solo publica después de revisar y confirmar. Windows permite sincronización Graph directa y copias por dispositivo.
- **Autorización Microsoft en Windows:** credencial de renovación cifrada con DPAPI y reutilizable al cambiar el EXE. En Mac sigue siendo temporal; Android usa Android Keystore.
- **TMDB Android:** relleno gradual de carátulas de las películas visibles, sin sincronizar claves. El ID de identificación manual puede viajar en `tmdb_overrides`, pero **las imágenes y cachés de plataformas no se comparten** entre equipos.
- **Interfaz:** zonas seguras de Android, ajustes Windows Beta 14, selección de plataformas y exportaciones TXT.

## Limitación de seguridad que debe conocer todo usuario

**Las protecciones no reparan automáticamente bibliotecas ya contaminadas por betas anteriores.** Si una lista tiene películas de otras listas, no confirms la sincronización y no borres los snapshots: primero exporta el JSON local bueno, compara sus listas y restaura desde una copia anterior verificada. **Una prueba automatizada satisfactoria no demuestra que una biblioteca personal de OneDrive esté intacta.**

En particular, no prometemos sincronización TMDB completa ni recuperación automática de carteles de Windows en el móvil. Las pruebas de Microsoft con cada cuenta real y la validación visual en todos los modelos de Mac y Android siguen abiertas.

## Instalación y conservación

1. Guarda un JSON independiente de las listas desde cada dispositivo antes de probar sincronización.
2. En Windows cierra el EXE anterior y abre el nuevo; la carpeta habitual es `%APPDATA%\Selector_Cine_365`. La sesión Microsoft DPAPI debería conservarse; si procede de una beta anterior a DPAPI, inicia sesión una vez.
3. En Android `es.cine365.app.limpia`, instala encima de Beta 13.12 Limpia **sin desinstalar**. Comprueba que el código sea 9027.
4. En Mac extrae el ZIP adecuado. Mac no tiene persistencia Keychain de la cuenta Microsoft en esta beta.
5. Antes de aceptar OneDrive comprueba que el total **de cada lista**, no solo el total de listas, coincide con lo esperado. Si aparece una mezcla, pulsa Cancelar.
6. Puedes usar la aplicación sin Microsoft ni TMDB; todas las listas se guardan localmente.

Documentación: [Manual recomendado](docs/MANUAL_BETA14_ANDROID13_13.md) · [Guía de sincronización segura](docs/GUIA_SINCRONIZACION_SEGURA_BETA14_ANDROID13_13.md) · [Historial](CHANGELOG.md) · [Política de versiones](docs/POLITICA_DE_RELEASES.md).

**Publicación:** prerelease. El código fuente disponible en la etiqueta de la release es una integración de las ramas de Android y Windows Beta 14; no se publican claves privadas ni datos personales.
