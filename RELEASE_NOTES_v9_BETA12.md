# Selector Cine 365 — v9 Beta 12

**Versión preliminar para Windows, Mac Apple Silicon, Mac Intel y Android.** Autor: Cristian Barallobre Rúa. El proyecto continúa en fase beta, con sincronización experimental, y debe probarse antes de usarlo como único almacén de listas.

## Descargas

- **Windows 10/11 (64 bits):** \`Selector_Cine_365_Windows_v9_BETA12.zip\` — WebView2, sin ventana CMD.
- **Mac Apple Silicon (M1 y posteriores):** \`Selector_Cine_365_Mac_Apple_Silicon_v9_BETA12.zip\` — Cocoa + WKWebView.
- **Mac Intel:** \`Selector_Cine_365_Mac_Intel_v9_BETA12.zip\` — Cocoa + WKWebView.
- **Android 8+ (firmado):** \`Selector_Cine_365_Android_v9_BETA12.apk\` — paquete permanente \`es.cine365.app\`, versión interna \`9012\`.

La APK Android está **firmada con el mismo certificado utilizado por Beta 11**. La publicación vuelve a comprobar su contenido SHA-256 y firma v2/v3 y **no incorpora la clave privada**. Mac tiene firma *ad hoc* pero no Developer ID ni notarización. Antes de instalar en Mac puede ser necesario autorizar la aplicación en Privacidad y seguridad.

## Novedades

- Búsqueda global o dentro de la lista actual, con deduplicación de películas por TMDB o título/año.
- Listas numeradas y navegación desde las fichas, edición de pertenencia sin perder las marcas.
- Año junto a valoración **TMDB**, no IMDb.
- Solo el tick marca Vista/Pendiente; se mantiene sincronización de marcas entre listas.
- Estadísticas Windows/Mac con flechas fijas. Listas visibles incluso cuando Nube ocupa mucho espacio.
- Historial de cambios con botón que comunica la consulta, número de eventos y estado vacío discreto.
- Android mantiene la cabecera fija, carga incremental y navegación táctil.
- Mantiene cambios anteriores: modo local opcional, copias JSON, PKCE de Microsoft, sincronización causal, recepción periódica, revisión de conflictos y recuperación.

## Actualizar conservando datos

**Primero exporta y verifica una copia JSON.**

- **Windows:** cierra la app anterior y reemplaza el ejecutable; mantiene el directorio de datos actual.
- **Mac:** cierra la app y reemplaza \`Selector Cine 365.app\` en Aplicaciones. Mantiene \`~/Library/Application Support/AtlasDelCine/\`.
- **Android Beta 11 → Beta 12:** instala sobre la aplicación antigua, sin desinstalar. Ambas usan el mismo identificador y firma.
- **Android Beta 9 → Beta 12:** Beta 9 usaba otro identificador y firma. **No se puede actualizar sin migración.** Exporta JSON fuera de la aplicación, verifícalo, desinstala Beta 9 solo si ya lo tienes a salvo, instala Beta 12 e importa el JSON. Revisa los datos antes de activar OneDrive.

## Comprobaciones y límites

- CI de Go, JavaScript y Gradle, integridad de los ZIP, firmas y arquitecturas.
- Mac Intel y Apple Silicon empaquetados desde las **mismas fuentes v9 de Windows**; no se han validado mediante uso prolongado en Macs físicos del usuario.
- La sincronización OneDrive no es instantánea y los cambios concurrentes sin orden causal requieren revisión.
- Las betas no sustituyen una copia JSON independiente.

**Documentación:** [manual completo](../../blob/main/docs/MANUAL_COMPLETO_SELECTOR_CINE_365.md) · [Guía Mac](../../blob/main/MACOS.md) · [Guía Android](../../blob/main/android/LEEME_ANDROID_v9_BETA.md) · [Historial](../../blob/main/CHANGELOG.md).
