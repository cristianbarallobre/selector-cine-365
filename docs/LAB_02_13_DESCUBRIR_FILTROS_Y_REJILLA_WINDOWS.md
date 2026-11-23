# Selector Cine 365 — LAB-02.13 · Descubrir completo y filtros TMDB

**Fecha:** 23 de noviembre de 2026  
**Rama:** `experimental` · **Plataforma:** Windows LAB exclusivamente.  
La rama `main`, Android y macOS permanecen intactos.

## 1. Descubrir: rellenar las filas disponibles

**Fallo reproducido:** a pesar de que TMDB mostraba `1 / 500`, aparecían ocho tarjetas y quedaba un rectángulo vacío en una ventana de aproximadamente 1480 px.

**Causa:** el cálculo del número de tarjetas se hacía antes de que la rejilla CSS hubiese establecido sus columnas. Algunas mediciones iniciales de ancho eran insuficientes y el límite calculado no alcanzaba ni para dos filas, por lo que la página TMDB 1 detenía la carga.

**Corrección:**
- Medición del ancho real del grid, su panel y la vista, con un mínimo de **tres filas de tarjetas**.
- Si hay seis columnas, se buscan al menos 18 películas, solicitando páginas TMDB consecutivas cuando la primera devuelve ocho (8+8+8). Como TMDB pagina en bloques, pueden aparecer 24 títulos —cuatro filas completas—.
- Se conserva el máximo de ocho páginas TMDB consultadas por página visible, con respuestas temporalmente en memoria, cancelación de peticiones obsoletas al cambiar de categoría/filtro y navegación sin reutilizar rangos anteriores.
- El número de páginas mostrado distingue la página visible del rango de páginas TMDB utilizadas.
- Si TMDB no tiene suficientes resultados, puede quedar espacio vacío; no se inventan películas ni se crean tarjetas repetidas.

**Prueba real con Edge:** `TestLabDiscoverFillsVisibleRowsInEdge` utiliza un servidor local de pruebas que devuelve **ocho películas por llamada**. Comprueba que una ventana ancha muestra **al menos tres filas**, que se solicitan más páginas y que la paginación y los filtros conservan sus resultados.

## 2. Filtros de género y duración en Descubrir

Se añaden debajo de Tendencias, Populares, Mejor valoradas y Próximos estrenos:

- **Género:** acción, aventura, animación, comedia, crimen, documental, drama, familiar, fantasía, historia, terror, música, misterio, romance, ciencia ficción, película de TV, suspense, bélica y western.
- **Duración:** cualquier duración, hasta 90 minutos, de 91 a 120 y más de 120 minutos.

**Funcionamiento:** cuando se aplica algún filtro en una categoría, el backend Windows LAB consulta **`/discover/movie`** de TMDB con `with_genres`, `with_runtime.gte` y/o `with_runtime.lte`. Los valores se validan en el servidor para no admitir géneros arbitrarios ni rangos inesperados. Los filtros se combinan y actualizar uno devuelve la paginación a la primera página.

**Limitaciones transparentes:**
- Las listas sin filtros conservan los cuatro endpoints y su orden habitual.
- TMDB no permite filtrar directamente el ranking *Tendencias de la semana*. Con filtros activos, las películas se ordenan por **popularidad**; la interfaz lo indica.
- Los filtros de género y duración se aplican a las **categorías de descubrimiento**. Durante una búsqueda **por título**, se muestran desactivados y aparece una nota que explica que no se aplican: TMDB no permite combinar `/search/movie` con filtros `with_runtime` y `with_genres` de `/discover/movie`. No se ignoran de forma silenciosa.
- Próximos estrenos filtrados utiliza fecha de estreno a partir del día actual y ordenación cronológica; Mejor valoradas filtradas utiliza un mínimo de 200 votos para evitar puntuaciones sin respaldo.

## 3. Ayuda «?» alineada con las otras pestañas

Se ha reducido el tamaño de fuente y el grosor a los de los botones vecinos; además, un ajuste en Windows LAB mide la altura real de **Biblioteca** y aplica esa altura a «?», incluido tras redimensionar, para que no sobresalga visualmente.

## 4. Calidad y pruebas

- `scripts/test_lab_discover_filters_windows.cjs`: controles y valores válidos, validación del backend, distinción entre búsqueda y categorías, invalidación de caché al filtrar.
- `scripts/test_lab_discover_fill_windows.cjs`: rejilla de seis columnas simulando páginas TMDB de ocho resultados, tres filas mínimas, navegación sin repeticiones y errores.
- `lab_discover_viewport_windows_test.go`: navegador Edge/Chrome en Windows, pantalla ancha, 8 resultados/página y filtros de género/duración comprobados dentro de la aplicación HTML real.
- `lab_boot_smoke_windows_test.go`: arranque, películas, navegación, filtros de Biblioteca y altura de «?» igual a Biblioteca.
- Las pruebas previas de cache/pre-carga, seguridad, red, Biblioteca, sorteo y Descubrir siguen ejecutándose.

**Artefacto:** `Selector_Cine_365_LAB-02_13_Windows`, disponible como artifact de la acción LAB de `experimental`.

Antes de llevar estas mejoras a la rama estable, deben revisarse visualmente en Windows con datos TMDB reales; las pruebas de navegador utilizan respuestas simuladas, sin gastar solicitudes reales de la clave del usuario.
