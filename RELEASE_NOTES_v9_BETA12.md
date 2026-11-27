# Selector Cine 365 — v9 Beta 12

**Versión preliminar para Windows, Mac Apple Silicon y Mac Intel.** Autor: Cristian Barallobre Rúa. El proyecto continúa en fase beta, con sincronización experimental, y debe probarse antes de usarlo como único almacén de listas.

## Descargas

- **Windows 10/11 (64 bits):** `Selector_Cine_365_Windows_v9_BETA12.zip` — WebView2, sin ventana CMD.
- **Mac Apple Silicon (M1 y posteriores):** `Selector_Cine_365_Mac_Apple_Silicon_v9_BETA12.zip` — Cocoa + WKWebView.
- **Mac Intel:** `Selector_Cine_365_Mac_Intel_v9_BETA12.zip` — Cocoa + WKWebView.


## Novedades

- Búsqueda global o dentro de la lista actual, con deduplicación de películas por TMDB o título/año.
- Listas numeradas y navegación desde las fichas, edición de pertenencia sin perder las marcas.
- Año junto a valoración **TMDB**, no IMDb.
- Solo el tick marca Vista/Pendiente; se mantiene sincronización de marcas entre listas.
- Estadísticas Windows/Mac con flechas fijas. Listas visibles incluso cuando Nube ocupa mucho espacio.
- Historial de cambios con botón que comunica la consulta, número de eventos y estado vacío discreto.
- Mantiene cambios anteriores: modo local opcional, copias JSON, PKCE de Microsoft, sincronización causal, recepción periódica, revisión de conflictos y recuperación.

## Actualizar conservando datos

**Primero exporta y verifica una copia JSON.**

- **Windows:** cierra la app anterior y reemplaza el ejecutable; mantiene el directorio de datos actual.
- **Mac:** cierra la app y reemplaza `Selector Cine 365.app` en Aplicaciones. Mantiene `~/Library/Application Support/AtlasDelCine/`.

## Comprobaciones y límites

- Mac Intel y Apple Silicon empaquetados desde las **mismas fuentes v9 de Windows**; no se han validado mediante uso prolongado en Macs físicos del usuario.
- La sincronización OneDrive no es instantánea y los cambios concurrentes sin orden causal requieren revisión.
- Las betas no sustituyen una copia JSON independiente.

