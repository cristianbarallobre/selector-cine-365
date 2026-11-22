> **ARCHIVO HISTÓRICO:** esta documentación conserva el estado de una versión anterior. La referencia de instalación y soporte actual es [Selector Cine 365 v9.0.0 estable](https://github.com/cristianbarallobre/selector-cine-365/releases/tag/v9.0.0) y su [manual completo](MANUAL_COMPLETO_SELECTOR_CINE_365.md). No seguir instrucciones de sincronización antiguas sin revisar las nuevas.

# Sincronización segura — Beta 14 Windows/Mac y Beta 13.13 Android

**21/11/2026 · Para pruebas con dos dispositivos y OneDrive.**

## Qué se sincroniza

Se comparten `colecciones.json` (listas, marcas, orden cuando procede, `tmdb_overrides`, eliminaciones) y snapshots individuales en `sincronizacion/dispositivos/*.json`. **No** viajan claves de TMDB, tokens Microsoft ni archivos de caché de carátulas/plataformas. Windows conserva su sesión mediante DPAPI; Android guarda su sesión con Android Keystore; Mac todavía necesita reautorizar.

## Reglas antes de pulsar confirmar

1. Exporta la biblioteca buena de Windows como JSON y guárdala **fuera** de OneDrive.
2. Comprueba cuántas películas tiene **cada lista** en Windows, no solo el total de vistas.
3. Si Android está recién instalado, abre Nube, conecta la misma cuenta y realiza una **vista previa**. El contenido inicial de 365 películas es solo un ejemplo local: la sincronización nueva lo sustituye por la biblioteca real **si no lo has modificado**.
4. En la vista previa, verifica **nombre + cantidad de películas de cada lista**; si una lista tiene inexplicablemente películas de varias, pulsa **Cancelar**.
5. Después de revisar la biblioteca, resuelve únicamente conflictos reales de películas vistas. Pulsa **Confirmar y sincronizar** para recibir/enviar una vez.
6. En Android vuelve a sincronizar manualmente cuando quieras. Si habilitas automático en Windows, ten en cuenta que pueden llegar snapshots nuevos sin abrir Android.

## Si las listas ya aparecen mezcladas

**No sincronices ningún dispositivo.** Una beta anterior añadía filas de otra lista cuando ambos registros compartían `id: atlas`. El código nuevo detiene nuevas fusiones sospechosas, pero **no separa por sí mismo las películas ya añadidas**.

- Conserva una copia de `colecciones.json` del Windows que tenga las listas correctas. Si no están correctas, busca un archivo anterior bajo `%APPDATA%\Selector_Cine_365\copias_automaticas` o `sync_recuperacion` y comprueba sus listas antes de restaurar.
- Conserva los snapshots antiguos de OneDrive. **No elimines las carpetas de sincronización** para forzar una recuperación.
- Evita confirmar una vista previa donde algún conteo sea inesperado; la solución requiere una recuperación selectiva desde una copia buena.

## Carpetas y ajustes

En Windows: `%APPDATA%\Selector_Cine_365`. En Android: almacenamiento privado `es.cine365.app.limpia` dentro de `Selector_Cine_365`. En nube: `Selector_Cine_365` de OneDrive; según las rutas históricas, los JSON pueden estar en `Selector_Cine_365_Copias/sincronizacion/dispositivos`.

Una carpeta vacía de Windows no equivale a una carpeta vacía de OneDrive. Puedes usar la aplicación **sin nube**; el guardado en disco sigue operativo. No es necesario abrir Microsoft al iniciar Windows si ya se guardó una autorización DPAPI válida.

[Manual de usuario](MANUAL_BETA14_ANDROID13_13.md) · [Descargas de esta beta](https://github.com/cristianbarallobre/selector-cine-365/releases/tag/v9.0.0-beta.14).
