# LAB-02.5 · Sorteo inteligente (Windows experimental)

## Qué incorpora

Nuevo interruptor **🎯 Sorteo inteligente** dentro de **Biblioteca → Tu próximo visionado**. Está **desactivado por defecto** y convive con el sorteo clásico sin modificarlo.

Cuando se activa, permite elegir:

- **Qué te apetece:** Equilibrado (valoración y variedad de géneros de la colección vista), Algo corto (favorece títulos de menor duración), Una buena valoración (favorece notas TMDB con suficientes votos) y Sorpréndeme (favorece títulos menos populares).
- **Duración:** Cualquiera, hasta 90 minutos, 91–120 minutos y más de 120 minutos.
- **Género:** Todos o cualquier género que aparezca en los metadatos TMDB de las películas pendientes.
- **Plataforma:** se sigue usando el selector existente, sin cambios.
- **No repetir:** cuando haya alternativas, evita las últimas ocho elecciones inteligentes durante la sesión.

El sorteo inteligente muestra una explicación breve junto al resultado. No marca películas como vistas por sí solo. Los botones **Copiar título** y **Marcar como vista** conservan su comportamiento, también en el modo de coincidencias de solo lectura.

## Selección y criterios

1. Solo entran películas pendientes de la lista actual o de la vista activa existente. Se respetan los filtros de streaming del sorteo, incluido «Todas».
2. Los filtros de género y duración son **estrictos**. Una ficha sin información suficiente no supera un filtro específico de ese dato; sin filtros, sigue participando.
3. Se excluyen de la bolsa las últimas ocho elegidas, salvo que con ello no queden alternativas. En ese caso, vuelven a participar (si es posible se evita la última elección).
4. Los criterios «Equilibrado», «Algo corto», «Una buena valoración» y «Sorpréndeme» cambian **probabilidades**, no prometen siempre la película mejor puntuada o más corta.
5. Se usa `crypto.getRandomValues` cuando el navegador lo permite, con alternativa aleatoria si no está disponible.
6. Las películas sin valoración fiable siguen siendo elegibles. La nota TMDB no es IMDb.

## Datos y privacidad

El modo inteligente utiliza `platformData` y la caché de fichas de LAB-02.4 **ya disponible en el ordenador**. La nueva ruta `GET /api/discover/smart-cache` devuelve solamente un resumen de duración, géneros y votos de películas en caché; no descarga datos nuevos de Internet ni escribe archivos. Es una ruta local con token de sesión y límite de 2500 fichas devueltas. La precarga TMDB existente sigue funcionando por separado.

La lista de últimas ocho elecciones se mantiene **únicamente en memoria durante la sesión** y no se incluye en JSON, OneDrive ni sincronización.

## Cómo probar

1. Ejecutar el artefacto **`Selector_Cine_365_LAB-02_5_Windows`** de [GitHub Actions en la rama experimental](https://github.com/cristianbarallobre/selector-cine-365/actions/workflows/ci-lab.yml?query=branch%3Aexperimental).
2. En Biblioteca, abrir el selector de plataforma de siempre y comprobar un sorteo clásico con el interruptor apagado.
3. Activar **Sorteo inteligente**, cambiar la preferencia a **Algo corto** y sortear varias veces.
4. Probar un género y duración específica; observar el contador de candidatas.
5. Seleccionar una plataforma concreta y comprobar que solo salgan títulos disponibles en ella.
6. Probar **Sorpréndeme**, **Una buena valoración** y **Equilibrado**, viendo la explicación.
7. Volver al sorteo clásico desactivando el interruptor.
8. Probar cuando falten datos TMDB: los filtros específicos deben mostrar cero candidatas si no hay ninguna ficha que los cumpla, sin inventar géneros ni duraciones.

## Aislamiento y pruebas

- Solo la rama `experimental`; la rama `main` no se modifica.
- El menú de Apariencia de LAB-03 **no vuelve a introducirse**.
- Solo Windows compilado con `-tags lab`, sin cambios en macOS.
- `plantilla.html` permanece intacta y el código se inyecta con `//go:build windows`.
- Pruebas de regresión del sorteo y del JavaScript, y pruebas Windows Go para autorización de caché, resumen de metadatos e integración HTML.
- El resultado visual y el comportamiento real en el PC requieren revisión manual antes de considerar la función estable.
