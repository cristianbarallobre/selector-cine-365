# LAB-02 — Descubrir y fichas completas (solo Windows)

**Canal:** rama `experimental`. **No está en `main` ni en la release estable `v9.0.0`.**

## Alcance


## Qué se ha añadido

- Nueva sección **✦ Descubrir** en la navegación de Windows LAB.
- **Tendencias**, **populares**, **mejor valoradas** y **próximos estrenos** consultados a TMDB en español. Región España para categorías que admiten región.
- Búsqueda por título, con resultados paginados (hasta cinco páginas por consulta en esta primera prueba).
- Tarjetas con carátula, título, año y puntuación TMDB, claramente diferenciada de una nota personal.
- Ficha de película con sinopsis, reparto (hasta diez intérpretes), dirección, géneros, duración, valoración, enlace al tráiler y enlace a TMDB.
- Disponibilidad de España separada por suscripción, gratis, anuncios, alquiler y compra; enlace a la plataforma cuando TMDB lo facilite.
- Añadir a **cualquiera de las listas actuales** como película pendiente. Evita duplicados por ID conocido de TMDB o por título y año. No importa ni reemplaza una lista entera.
- **Botón «Detalles» explícito** junto al título en Biblioteca y dentro de cada tarjeta de Descubrir. Las fichas ya no se abren al pulsar cualquier punto de la fila o carátula. Cuando la identificación es ambigua, el botón lleva a resultados para elegir manualmente.
- **Corrección visual:** el recuadro de reserva se oculta al cargar correctamente la carátula en la ficha.

## Requisitos y seguridad

- Hay que conectar la clave personal de TMDB desde **Biblioteca → Conectar TMDB** en la propia instalación LAB. La clave no se copia automáticamente desde la instalación estable.
- La búsqueda es optativa y necesita internet; las listas locales siguen funcionando desconectadas.
- El servidor local exige el token de la sesión para las rutas `/api/discover/*`. No entrega la clave TMDB al navegador; utiliza el conector TMDB existente.
- Las carátulas usan una ruta proxy limitada a rutas de imagen válidas de TMDB. Los enlaces externos se abren con `noopener noreferrer`.
- El modo LAB utiliza `Selector_Cine_365_LAB` para los datos locales y otra carpeta en OneDrive; no se deben conectar las copias de producción como destino manual.
- Los cambios se guardan mediante el procedimiento habitual de la biblioteca y sus precauciones de sincronización; no se hace una importación masiva automática al pulsar una película.
- Para comprobar las fichas de una película añadida sin conexión, se conserva el título en la biblioteca pero la información TMDB completa no estará disponible hasta volver a conectarse.

## Descargar y probar

Abrir [Actions → LAB - pruebas y compilaciones independientes](https://github.com/cristianbarallobre/selector-cine-365/actions/workflows/ci-lab.yml?query=branch%3Aexperimental) y descargar el artefacto **`Selector_Cine_365_LAB-02_Windows`** de una ejecución correcta.

Se genera un EXE portable de Windows x64 sin firma Authenticode, con datos separados. El botón Descubrir no aparece en la aplicación oficial ni en el canal de laboratorio de Mac. No se publica como release oficial.


## Revisión LAB-02.4 — altura de filas y centrado (solo Windows)

- **Biblioteca Windows LAB**: se optimiza la ocupación vertical de las fichas. En la plantilla compartida existían reglas antiguas de 85 px y una posterior de 96 px (altura efectiva); se deja intacta la plantilla.
- Las filas de LAB ahora se ajustan automáticamente desde **96 px (si falta espacio) hasta 110 px (cuando sobra)**. Se usa la altura real disponible en cada renderizado y se conserva el número de películas por página que cabía antes, para evitar que hacer las fichas algo más altas provoque todavía más hueco vacío.
- Carátulas ligeramente más grandes (42 × 60 px), centradas junto al texto y los botones. El aspecto, la altura de otras vistas y la navegación de la aplicación estable siguen inalterados.
- **Pruebas**: validación estática de las reglas de estilo e inyección solo en LAB; casos de distintas alturas de ventana para evitar recortes y pérdidas de capacidad; compilación Windows y pruebas Windows nativas.
- **Prueba manual recomendada**: abrir una lista con varias películas, mirar cómo se distribuyen sus filas; agrandar y reducir la ventana; pasar de página y volver de Descubrir para comprobar que todo sigue alineado.
- Descarga: `Selector_Cine_365_LAB-02_4_Windows` de una ejecución correcta de [Actions LAB](https://github.com/cristianbarallobre/selector-cine-365/actions/workflows/ci-lab.yml?query=branch%3Aexperimental).

## Revisión LAB-02.3 — Correcciones de interacción y paginación (solo Windows)

- En **Descubrir**, se elimina por completo la **paginación superior**. Se mantiene únicamente la inferior, debajo de las tarjetas.
- Solo el botón **Detalles** abre cada ficha; las fotos, los títulos y el fondo de la tarjeta ya no son clicables ni activables con el teclado. El botón permite abrir la ficha, ver sus enlaces y añadir la película a una lista.
- En **Biblioteca**, el orden sigue **Pendiente/Vista → Detalles → Lápiz**, pero ahora hay **13 px de separación** entre Detalles y el lápiz, además del margen anterior.
- **Corrección del fallo de una sola película:** el navegador ejecutaba el manejador normal de Biblioteca antes de ocultar Descubrir. Como ambas vistas tienen distribución flex, la lista calculaba cuántas filas cabían mientras compartían temporalmente el espacio, obteniendo a veces 1 sola fila. Ahora se oculta Descubrir durante la fase de captura del clic, **antes** de mostrar/renderizar Biblioteca. Se añade observador de tamaño para rehacer la paginación si una pestaña o ventana altera su altura. No se modifica el contenido de las listas.
- Se conservan la caché, la precarga de fichas y la paginación inferior completa de TMDB.

Descarga **`Selector_Cine_365_LAB-02_3_Windows`** de la última ejecución satisfactoria en [Actions LAB](https://github.com/cristianbarallobre/selector-cine-365/actions/workflows/ci-lab.yml?query=branch%3Aexperimental). Para probar: abre Descubrir, vuelve a Biblioteca, comprueba que se muestren varias películas; abre y cierra alguna ficha y repite el cambio de pestaña. Las compilaciones anteriores LAB-02.2 no incluyen estos cambios.

## Revisión LAB-02.2 — Descubrir y orden de botones (solo Windows)

- La sección **Descubrir** utiliza un scroll vertical propio. Sus paneles ajustan la altura al contenido y ya no recortan los nombres, botones «Detalles», el resto de las carátulas ni la navegación.
- En **Descubrir**, pulsar cualquier zona de la tarjeta (imagen, título o fondo), o su botón «Detalles», abre la misma ficha con la opción de añadir a la lista. Teclado: Tab hasta la tarjeta, Intro o Espacio para abrirla.
- La búsqueda y las categorías paginan los resultados que publica TMDB: 20 títulos por página es el tamaño de una página, **no el total de películas**. Se muestra el total de resultados devuelto por TMDB y controles Anterior/Siguiente tanto **arriba** como **debajo** de las tarjetas, hasta 500 páginas (límite del API de TMDB).
- En **Biblioteca**, los controles a la derecha de cada fila se ordenan: **Pendiente/Vista → Detalles → Lápiz**. La flecha externa sigue eliminada; al pulsar el resto de la fila no se abre la ficha.
- Si TMDB aún no ha anunciado la fecha de estreno, la película se puede añadir como pendiente sin año, conservando su ID para futuras consultas.
- **Por validar manualmente** en Windows: vista a distintos tamaños, scroll de 20 tarjetas, cambio a página 2 y posteriores, apertura desde foto/título/Detalles y guardado en una lista existente.

Descarga el artefacto `Selector_Cine_365_LAB-02_2_Windows` de la ejecución correcta más reciente de [Actions LAB](https://github.com/cristianbarallobre/selector-cine-365/actions/workflows/ci-lab.yml?query=branch%3Aexperimental). No utilices el EXE LAB-02.1 para probar estas correcciones.

## Revisión LAB-02.1 — caché y botones (solo Windows)

- La biblioteca mantiene el botón **Detalles** al lado del título, la pastilla Vista/Pendiente y el lápiz para corregir TMDB. Se retira el botón externo **↗** de las filas; sus enlaces siguen dentro de la ficha.
- Descubrir también abre las fichas exclusivamente desde el botón **Detalles**. Ninguna tarjeta completa tiene que ser clicable.
- Desaparece el recuadro de carátula de reserva cuando ya se cargó la imagen.
- La primera apertura con TMDB conectado lanza **precarga de fichas de todas las listas**, priorizando la lista activa, sin bloquear la biblioteca. Una misma película se descarga solo una vez si comparte ID o título y año.
- La cola trabaja **secuencialmente**, espera al menos 900 ms entre películas y respeta las respuestas de límite de TMDB. Hay progreso **Fichas guardadas: N / total** y botones **Pausar / Reanudar**. Si se cierra la app, el progreso de los datos guardados se conserva y las fichas pendientes se reprograman al iniciar.
- Los detalles se guardan en `%APPDATA%\\Selector_Cine_365_LAB\\tmdb_fichas_lab.json`, por separado de `colecciones.json`. Los identificadores resueltos también se guardan, de modo que una película ya preparada abre «Detalles» sin nueva búsqueda.
- Actualizaciones de caché: **30 días** para información general, **7 días** para nota/tráiler y **24 horas** para disponibilidad de plataformas. Si la ficha ya existe se muestra al instante, aunque esté obsoleta, y se refresca en segundo plano.
- Fichas guardadas disponibles **sin internet** (texto y enlaces conservados; las carátulas dependen de la caché del navegador). Las plataformas pueden estar desactualizadas y se indican como tales al refrescar.
- La caché no copia claves TMDB ni tokens; utiliza la cuenta TMDB ya configurada localmente en LAB. Ni las listas ni los datos compartidos de OneDrive se transforman para esta precarga.
- **Pendiente de validación manual en Windows:** comprobar primeras 365 fichas, rendimiento real de TMDB, conexión lenta, duplicados, sinopsis y disponibilidad de España. La compilación CI por sí sola no prueba el comportamiento en tu PC.


## Validación

- Pruebas Go del proyecto y del canal LAB.
- Comprobación de sintaxis y aislamiento del JavaScript inyectado: `node scripts/test_lab_discover_windows.cjs`.
- Compilación cruzada de un ejecutable Windows con etiqueta LAB.
- Pendiente de **prueba manual en Windows con una clave TMDB real** antes de considerar aprobada LAB-02.

**El paso a producción requiere aprobación explícita.** No fusionar `experimental` automáticamente con `main`.
