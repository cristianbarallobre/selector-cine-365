# Selector Cine 365 — LAB-02.10 · Biblioteca, Estadísticas y Descubrir (Windows)

**Rama:** `experimental` · **Plataforma:** Windows LAB exclusivamente. La versión estable (`main`) y macOS **no se modifican**.

## Biblioteca: botón único «Filtros»

En la parte superior de Biblioteca quedan visibles el buscador, el botón **☷ Filtros** con el número de filtros activos y el contador de películas mostradas. Los controles «Lista actual / Todas las listas» y «Todas / Pendientes / Vistas» dejan de ocupar filas y se trasladan al panel de filtros. También se mueve allí el filtro múltiple de plataformas, sin perder sus opciones ni su funcionamiento.

El botón abre una ventana de filtros con:

| Categoría | Opciones |
|---|---|
| Buscar en | Lista actual / Todas las listas |
| Estado | Todas / Pendientes / Vistas |
| Plataformas | Todas las plataformas o una o varias plataformas de streaming; se reutiliza el selector anterior |
| Década | Décadas realmente presentes en las películas del ámbito elegido |
| Género | Géneros disponibles en los metadatos TMDB de esas películas |
| Duración | Cualquiera, hasta 90, de 91 a 120 y más de 120 minutos |
| Ordenar por | Orden de la lista, título A–Z/Z–A, estreno reciente/antiguo, nota TMDB, duración ascendente/descendente, vistas primero y pendientes primero |

- **Todas las listas ahora funciona aunque no escribas nada en el buscador**, y permite combinar estado, género, década, duración y plataformas.
- La ordenación **no edita el orden de los archivos TXT**: se deshabilita arrastrar películas cuando está activa una ordenación o un filtro que impida reordenar con seguridad.
- El modo sin filtros conserva el comportamiento anterior. Al cambiar de lista se restablecen los filtros avanzados, para no hacer parecer vacía otra lista por una selección olvidada.
- Los filtros de duración y género se basan en metadatos TMDB locales. Si no se conoce el dato, la película no puede cumplir un filtro específico de duración o género. Sin tales filtros, puede mostrarse con normalidad.
- **Restablecer filtros** devuelve el ámbito a la lista actual, el estado a todas, las plataformas a todas y los controles adicionales a sus valores iniciales; no modifica qué películas están vistas.
- Se conserva la paginación Biblioteca, los botones Detalles y el sorteo inteligente.

## Estadísticas: cinco indicadores útiles

Se retiran las tarjetas **Pendientes** y **Fichas TMDB**. Quedan:

1. **Películas vistas**.
2. **Horas registradas** (suma de duraciones conocidas).
3. **Nota media TMDB** de las películas vistas que tengan valoración disponible (no es una nota personal).
4. **Década más vista** entre los años conocidos.
5. **Género más visto** entre los géneros disponibles en TMDB.

Se conservan los gráficos ya existentes por década, géneros, países de producción y dirección. Si faltan metadatos, se muestra «—»; no se inventan valores.

## Navegación: ayuda más discreta

La pestaña **Instrucciones** se representa como un botón pequeño **«?»**, con su descripción accesible y su contenido original intactos.

## Descubrir: relleno automático de filas

Si una página de TMDB devuelve menos títulos de los necesarios para ocupar las filas visibles, Descubrir continúa con páginas **consecutivas**, hasta completar el espacio cuando hay suficientes títulos. El cálculo considera las columnas y la altura disponibles.

- La siguiente página visible empieza **después de la última página TMDB utilizada**, sin repetir películas de la vista anterior.
- La etiqueta de navegación indica la página visible y el rango de páginas TMDB usado.
- Cuando se cambia de categoría o búsqueda se inicia una secuencia independiente.
- Al volver hacia atrás se reutilizan las respuestas de la sesión cuando siguen disponibles, evitando solicitudes redundantes.
- Se limita a **ocho páginas TMDB por pantalla** para evitar una ráfaga ilimitada de consultas y cumplir las restricciones del servicio.
- Al redimensionar y aumentar la capacidad se vuelven a calcular las tarjetas necesarias sin mezclar la nueva página con otras ya mostradas.
- Si se agotan las páginas disponibles o hay un fallo de red, se muestran los resultados efectivamente recuperados y un mensaje explicativo; no se inventan títulos.

## Seguridad y alcance

El nuevo panel usa solamente el estado que ya existe en las listas y los metadatos de la caché TMDB de LAB-02.4. **No añade grabaciones en disco ni modifica el protocolo de OneDrive**. Continúa el refuerzo de conexiones de LAB-02.9: reintentos solo para consultas seguras, sin repetición automática de escrituras.

Los cambios de interfaz se inyectan desde los archivos `lab_filters_windows.html`, `lab_filters_inject_windows.go` y `lab_discover_windows.html` marcados para Windows. `plantilla.html` sigue intacta.

## Pruebas

- `node scripts/test_lab_filters_windows.cjs`: combinación de década, género, duración, ordenación y estadísticas.
- `node scripts/test_lab_discover_fill_windows.cjs`: completar una página corta con otra, navegar sin duplicados, volver con caché y manejar errores.
- `go test -tags lab -run TestLab -count=1 -v .` en Windows: verificar inyección real, coexistencia con Descubrir, sorteo y las funciones anteriores.
- Pruebas de regresión existentes de Windows, caché, precarga y OneDrive.

**Artefacto:** `Selector_Cine_365_LAB-02_10_Windows` en [GitHub Actions de experimental](https://github.com/cristianbarallobre/selector-cine-365/actions/workflows/ci-lab.yml?query=branch%3Aexperimental).

La revisión visual final en Windows sigue pendiente de probar el EXE real en el PC.
