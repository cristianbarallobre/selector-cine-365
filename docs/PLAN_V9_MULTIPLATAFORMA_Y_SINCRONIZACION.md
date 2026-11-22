# Selector Cine 365 — Estado técnico y mantenimiento de v9.0.0

**Actualizado el 22 de noviembre de 2026.** El documento original del plan de las betas 13/13.1 se conserva en [Archivo histórico](archivo/PLAN_V9_BETA13_1_ORIGINAL.md).

## Estado actual

**Versión estable recomendada:** [v9.0.0](https://github.com/cristianbarallobre/selector-cine-365/releases/tag/v9.0.0), integrada con los ejecutables comprobados en Windows Beta 14, macOS Beta 14 (Intel/Apple Silicon) y Android Beta 13.13 Limpia. Promoción **sin cambios funcionales ni recompilación** respecto a la beta conjunta publicada el 21 de noviembre.

| Funcionalidad | Estado en v9.0.0 |
| --- | --- |
| Guardado local sin cuentas ni TMDB | Implementado en las tres plataformas |
| Listas múltiples, búsqueda global, vistas, sorteo, estadísticas | Implementado |
| Exportar e importar copias JSON / listas TXT | Implementado según la interfaz |
| OneDrive personal con OAuth PKCE | Implementado; requiere cuenta Microsoft y revisión de los datos |
| Android: sincronización manual con confirmación | Implementado |
| Windows: sesión Microsoft persistente DPAPI | Implementado y probado en CI |
| Mac: sesión Microsoft persistente Keychain | **Pendiente**, la sesión puede expirar al reiniciar |
| Mezcla de listas con IDs antiguos repetidos | Nuevas fusiones sospechosas bloqueadas; **recuperación histórica manual** |
| Android: carátulas y proveedores TMDB | Caché local y recuperación gradual de fichas visibles |
| Transferencia de caché TMDB de Windows a Android | **No implementado**; solo pueden compartirse los IDs manuales |
| Firma APK Android | APK estable firmada con certificado original; futuras compilaciones deben reutilizarlo |
| Firma Windows Authenticode y notarización macOS | **No implementadas** |

## Identidades y compatibilidad

**Android estable:** `es.cine365.app.limpia`, versionCode 9027. No sustituirlo por la línea histórica `es.cine365.app`; cada paquete tiene sus datos privados. Una futura actualización dentro de la línea Limpia requiere versionCode superior y **la misma clave de firma privada original**, conservada fuera del repositorio.

**Datos:** `Selector_Cine_365` en el perfil de cada plataforma y en OneDrive. La biblioteca en cada equipo conserva sus listas locales, que no requieren internet. Ver [rutas](UBICACIONES_DATOS_SELECTOR_CINE_365.md).

## Estrategia de mantenimiento

1. Mantener `v9.0.0` como referencia reproducible de código y binarios hasta que exista un problema confirmado.
2. Los arreglos futuros deben comenzar en una rama desde `main`, incluir pruebas de regresión y no cambiar el identificador/firma Android.
3. Antes de publicar, realizar pruebas en Android y Windows/macOS; comparar hashes de binarios y revisar la sincronización con una biblioteca de pruebas, nunca con datos reales irremplazables.
4. Mantener un archivo privado verificado de releases redundantes y conservar **tags históricos**. Solo una release debe estar marcada `Latest`.
5. Regenerar el manual Markdown/PDF/DOCX y la guía de Microsoft de la misma revisión fuente. No mezclar instrucciones Beta 13 con Beta 14.
6. Ante copias ya mezcladas, recuperar primero desde `colecciones.json`/JSON sanos; ninguna migración silenciosa se considera segura.
7. No incluir credenciales privadas ni keystore en GitHub.

## Mejoras posibles en una versión posterior

- Persistencia de autorización Microsoft en Keychain de macOS.
- Rutina explícita, revisable y no destructiva para restaurar listas previamente mezcladas.
- Mejor reutilización multiplataforma de identificaciones TMDB y manejo de errores/cupos, sin exponer API keys.
- Pruebas de actualización Android con paquete estable y certificado permanente en CI protegido.
- Firma y notarización de Windows/Mac para facilitar instalación.
- Verificaciones de accesibilidad y consumo de batería en más dispositivos.

**No se introducen estas mejoras en v9.0.0**: la decisión actual es conservar las compilaciones probadas.

[Manual completo](MANUAL_COMPLETO_SELECTOR_CINE_365.md) · [Historial](../CHANGELOG.md) · [Política de releases](POLITICA_DE_RELEASES.md).
