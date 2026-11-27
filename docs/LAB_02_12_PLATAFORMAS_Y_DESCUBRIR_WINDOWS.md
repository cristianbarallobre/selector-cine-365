# LAB-02.12 — Plataformas visibles, filtros accesibles y Descubrir sin huecos

**Fecha:** 23 de noviembre de 2026  
**Plataforma:** Windows LAB · Rama **`experimental`** exclusivamente.  
No se han modificado la versión estable `main` ni macOS.

## 1. Plataforma sin desplegable en el filtro de Biblioteca

El cuadro **Filtrar biblioteca** presenta directamente:

- Encabezado **Seleccionar plataformas**.
- Una opción **«Todas las plataformas»**, marcada de forma visible al abrir el panel con el filtro vacío. Significa que se muestra cualquier película, con independencia de si tiene streaming registrado.
- Debajo, todas las plataformas disponibles como casillas seleccionables, con desplazamiento dentro de la lista cuando hay muchas. El usuario **no tiene que abrir otro menú** dentro del diálogo para verlas.
- Seleccionar Netflix, Filmin u otra plataforma desmarca «Todas las plataformas» y muestra solo películas con ofertas de streaming en alguna de las plataformas elegidas. Marcar de nuevo «Todas las plataformas» elimina las restricciones específicas de plataforma.
- El control anterior **«VER EN LA LISTA»**, el botón con desplegable y el pequeño botón «Todas» dejan de mostrarse. Internamente se reutilizan los mismos checkboxes y el estado de filtro, sin duplicar sus eventos.

## 2. Limpiar filtros fuera del diálogo

En la barra superior de Biblioteca, junto a **☷ Filtros**, aparece **↺ Limpiar filtros**. Está desactivado si no hay filtros activos.

Restablece el alcance a lista actual, el estado a Todas, plataformas a Todas, década y género sin filtro, duración cualquiera y ordenación original. No altera películas vistas ni el contenido de las listas TXT.

## 3. Descubrir aprovecha la rejilla completa

Causa concreta corregida: al empezar con el grid CSS `auto-fill` vacío, `getComputedStyle(...).gridTemplateColumns` podía devolver `none`. El cálculo anterior interpretaba que había una única columna y detenía la carga tras una página TMDB, aunque el ancho de la ventana permitiera seis columnas y varias filas.

Ahora el número de columnas se calcula a partir del **ancho real del grid vacío**, el espacio entre tarjetas y el ancho mínimo de cada tarjeta (180 px). Si una página TMDB devuelve solo ocho resultados y caben doce, se solicitan páginas consecutivas hasta rellenar las filas o agotar los resultados disponibles. Se respetan el límite previo de ocho consultas por página visual y las respuestas en caché; las páginas siguientes no vuelven a repetir páginas TMDB que ya se utilizaron.

Los tests simulan específicamente que `gridTemplateColumns === "none"` antes de dibujar las primeras películas y que se completan dos páginas de ocho resultados sin duplicados.

## 4. Información de películas añadidas desde Descubrir

La comprobación del usuario confirma que la película **sí acababa recibiendo sus datos en Biblioteca**, solo con retraso. No había indicios de que se perdiera el título o su identidad en las listas.

Ahora, al pulsar **Añadir pendiente** en una ficha de Descubrir:

1. Se conserva el guardado de la lista y del identificador TMDB como hasta ahora.
2. Se transfieren **en el momento**, desde la ficha ya cargada, el cartel, la valoración, los votos, la duración, los géneros y las plataformas disponibles al estado de Biblioteca.
3. Se solicita en segundo plano una actualización del registro nativo de plataformas, sin retrasar el renderizado ni repetir el guardado.
4. Si esa actualización falla, identifica otra película o devuelve datos parciales, **no sobrescribe ni borra la ficha que ya se estaba mostrando**.

Esto beneficia también al caso de añadir una película a una lista nueva y abrir dicha lista después. Tras reiniciar la aplicación, los datos se siguen recuperando mediante la caché y el sistema nativo de plataformas existentes; la primera aparición depende de los datos disponibles localmente.

## 5. Pruebas

- `scripts/test_lab_filters_windows.cjs`: «Todas las plataformas» seleccionada por defecto, casillas visibles sin desplegable, reset externo, filtros avanzados y métricas.
- `scripts/test_lab_discover_fill_windows.cjs`: grid vacío con `gridTemplateColumns=none`, relleno 8+8, paginación sin títulos repetidos y fallo de red.
- `scripts/test_lab_discover_add_windows.cjs`: aparición inmediata de carátula, datos, plataformas y protección contra actualizaciones incorrectas/incompletas.
- `lab_boot_smoke_windows_test.go`: prueba real en navegador Edge/Chrome con comprobación de opciones visibles, «Todas» seleccionada, selección de plataforma, limpieza desde fuera, navegación, Estadísticas y sorteo.
- Se mantienen las pruebas de regresión Go, caché TMDB y demás pantallas.

**Artefacto:** `Selector_Cine_365_LAB-02_12_Windows` en [GitHub Actions de la rama experimental](https://github.com/cristianbarallobre/selector-cine-365/actions/workflows/ci-lab.yml?query=branch%3Aexperimental).

Sigue siendo necesario comprobar el comportamiento visual con una lista real y las consultas reales a TMDB en el PC antes de considerar esta LAB candidata a estable.
