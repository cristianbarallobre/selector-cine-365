# LAB-02.6 · Ajustes del sorteo inteligente y ficha del resultado (Windows)

**Rama:** `experimental`. **Plataforma:** exclusivamente Windows LAB. No afecta a `main`, Android ni macOS. LAB-03 Apariencia sigue eliminada.

## Correcciones y novedades

### 1. Algo corto = máximo 90 minutos

- La opción **«Algo corto · hasta 90 min»** ahora descarta de manera estricta las películas de duración mayor a 90 minutos, incluidas películas de 91 y 238 minutos.
- **Al elegir «Algo corto» se pone automáticamente «Hasta 90 minutos» en Duración**, evitando que permanezca visible «Cualquier duración».
- Si se cambia el selector Duración a «Cualquier duración», «De 91 a 120» o «Más de 120» mientras está seleccionado «Algo corto», el modo cambia a **Equilibrado** para no presentar criterios contradictorios.
- Las películas sin duración TMDB conocida no cumplen el filtro corto. Si no hay coincidencias, se desactiva el botón de sorteo y aparece un mensaje explicativo; no se selecciona una larga por defecto.
- El filtro sigue respetando género, las películas pendientes y el selector de plataforma.

### 2. Selector de plataforma

- El valor visible por defecto pasa de «Todas las películas pendientes» a **«Todas las plataformas»** tanto en el botón como en el menú y su reconstrucción tras cargar TMDB.
- Es **un cambio de etiqueta exclusivamente en Windows LAB**. El valor interno sigue siendo `all`, que considera todas las películas pendientes (también aquellas sin streaming registrado). Para limitar a películas con streaming, se conserva «Cualquier plataforma (streaming)».
- No se cambia la lógica de disponibilidad ni los nombres de plataformas concretas.

### 3. Detalles de la película sorteada

- Cuando se elige una película, aparece un botón **«Detalles»** en la esquina superior derecha de la tarjeta de resultado. Antes del sorteo permanece oculto.
- Funciona tanto con **«🎲 Elegir una película al azar»** como con el **Sorteo inteligente**.
- Abre la ficha de película **existente en LAB-02**, con sinopsis, géneros, reparto, dirección, duración, tráiler y plataformas en España cuando TMDB dispone de esos datos.
- Si no se puede identificar una película de forma unívoca, la acción usa la búsqueda de Descubrir para seleccionar la ficha correcta. No asigna a ciegas un ID TMDB.
- El botón no cambia el estado de vista/pendiente, ni altera Copiar título ni Marcar como vista. El resto de la tarjeta sigue sin ser clicable.

## Pruebas

Las pruebas automáticas comprueban:
- Filtrado estricto de 80, 90, 91 y 238 minutos; duración desconocida y combinación con género.
- Sincronización entre «Qué te apetece» y «Duración».
- Aparición/ocultación del botón Detalles y apertura de la película sorteada.
- Interfaz inyectada en Windows LAB, etiqueta «Todas las plataformas» en los tres puntos y ausencia de cambios en la plantilla compartida.
- Regresiones anteriores de Descubrir, biblioteca y caché TMDB.

Artefacto de Windows: **`Selector_Cine_365_LAB-02_6_Windows`** en [GitHub Actions experimental](https://github.com/cristianbarallobre/selector-cine-365/actions/workflows/ci-lab.yml?query=branch%3Aexperimental).

La comprobación visual en Windows (espaciado y tamaño de la tarjeta) queda pendiente de las pruebas del usuario.
