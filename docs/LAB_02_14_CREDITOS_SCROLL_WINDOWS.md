# LAB-02.14 — Actores y personajes separados; desplazamiento de Mis listas (Windows)

**Fecha:** 23 de noviembre de 2026  
**Rama:** `experimental` · **Plataforma:** Windows LAB solamente.  
**Sin cambios en:** `main` ni macOS.

## Ficha de película: reparto más claro

En la sección **Reparto y dirección** de las fichas de Descubrir (también cuando se abre una ficha desde Biblioteca), ya no se concatenan actores y personajes con un punto separador ambiguo.

Ahora se muestran:

- **Dirección** en una línea identificada con el nombre del director destacado.
- **Reparto y personajes**, en una cuadrícula de tarjetas individuales.
- En cada tarjeta, el **nombre del actor** aparece en negrita y blanco; debajo, en color secundario, el texto **Personaje: [nombre]**.
- Si TMDB conoce el actor pero no su papel, solo se muestra el nombre del actor; no se inventa un personaje.
- El contenido se crea con `textContent`, evitando interpretar como HTML los nombres procedentes de TMDB.

Ejemplo de formato: **Andrew Koji** y, debajo, **Personaje: Ryu**; **Noah Centineo** y, debajo, **Personaje: Ken Masters**.

Los datos ya estaban disponibles en el JSON de la ficha (`cast[].name` y `cast[].character`). No se realizan nuevas llamadas a TMDB, no se pierden datos de caché ni se modifica el formato de las listas.

## Mis listas: desplazamiento igual a Descubrir

La captura mostraba la barra de desplazamiento **dentro del panel redondeado** de Mis listas, mientras Descubrir la situaba en el lateral de la **vista completa**.

En LAB-02.14:

- `#view-collections` es el contenedor que hace scroll, como `#labDiscoverView`.
- El panel interior deja de tener scrollbar propia y permite que todo su contenido fluya verticalmente.
- El scroll comparte grosor, color, comportamiento de contención y separación de borde con Descubrir.
- Cabecera, selector de lista y navegación superior permanecen fuera del área desplazable.
- Se preserva la paginación de listas, todos los botones de gestión y la sección OneDrive; las funcionalidades de sincronización no cambian.
- El cambio está aislado al CSS de Windows LAB: `plantilla.html` y otras plataformas se mantienen exactamente igual.

## Pruebas

- `scripts/test_lab_cast_collections_scroll_windows.cjs` verifica renderización estructurada de actores, roles y director, el tratamiento de papeles vacíos, el uso de datos TMDB y la ubicación del scroll.
- `lab_boot_smoke_windows_test.go` verifica con **Edge/Chrome headless en Windows** que el desplazamiento está en `#view-collections`, que el panel no tiene scroll propio y que el encabezado y las pestañas no se desplazan al mover la lista.
- `lab_discover_viewport_windows_test.go` abre una ficha TMDB simulada en el navegador y comprueba que **Andrew Koji** y **Ryu** aparecen separados, con el rótulo «Personaje», y que **Kitao Sakurai** aparece como director.
- Se conservan todas las pruebas anteriores de Descubrir, filtros, Biblioteca, nube y caché.

**Artefacto:** `Selector_Cine_365_LAB-02_14_Windows` en GitHub Actions de la rama `experimental`.

Como en versiones LAB anteriores, será necesaria una comprobación visual con el EXE en Windows del usuario para confirmar el aspecto de las tarjetas y el desplazamiento con sus listas reales.
