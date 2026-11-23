# LAB-02.11 · Corrección de arranque y navegación (Windows)

**Rama:** `experimental`. **Sin cambios en:** `main`, Android, macOS.

## Error detectado en LAB-02.10

La pantalla mostraba la cabecera y el panel del sorteo, pero la lista quedaba vacía, los controles originales de ámbito/estado/plataforma seguían fuera de Filtros y algunas pestañas no respondían.

**Causas comprobadas mediante una prueba real con Microsoft Edge en Windows:**

1. La aplicación ejecutaba su inicialización original (`update()`, `renderList()` y `initLibrary()`) antes de llegar a las declaraciones `const` del nuevo panel de filtros. El primer renderizado accedía a `labFilterDecade` antes de inicializarla y detenía el script.
2. Al compilar/probar directamente en Windows, Git puede convertir los HTML a saltos de línea CRLF. La inyección de Descubrir buscaba una secuencia de cierre de script con LF exactos, de manera que parte de su JavaScript no se insertaba. Aparecían errores `labCreateDetailButton is not defined` y `labLeave is not defined`.
3. El arranque también incluye la llamada `labQueuePrefetchSoon()` añadida por Descubrir. El intento inicial de trasladar una línea exacta sin esa llamada no coincidía. Se corrigió para buscar el prefijo estable y trasladar la línea real íntegra.

## Corrección

- El inicio de Biblioteca y de la aplicación se ejecuta **solo después de inicializar Descubrir, Sorteo, diagnósticos y Filtros**.
- La inyección de Descubrir normaliza LF/CRLF para que ambas formas de compilación generen el mismo comportamiento.
- Se mantienen las novedades de LAB-02.10: filtro único en Biblioteca, ámbito/estado/plataforma dentro del panel, década/género/duración/ordenación, cinco indicadores de Estadísticas, botón «?» y relleno de Descubrir.
- No se modifican los datos locales, sus formatos ni el protocolo de OneDrive.

## Pruebas añadidas

`lab_boot_smoke_windows_test.go` realiza comprobaciones específicas:

- JavaScript **completo** generado por `makeHTML()` validado con `node --check`; no solo cada módulo aislado.
- Orden comprobado: todos los módulos y variables necesarios preceden a la línea real que inicia la interfaz.
- Inyección comprobada de `labLeave` y `labCreateDetailButton`.
- Prueba integral de interfaz en **Edge/Chrome headless en Windows**: aparecen películas en la biblioteca; los selectores anteriores se trasladan al panel y dejan de verse fuera; el panel abre/cierra; navegan Biblioteca, Estadísticas, Descubrir, Mis listas e Instrucciones; Estadísticas muestra cinco indicadores; el sorteo clásico elige una película y aparece el botón Detalles.
- Se conservan las pruebas Go y Node existentes, incluida la caché TMDB, el sorteo inteligente, las redes y el autollenado de Descubrir.

**Artefacto Windows:** `Selector_Cine_365_LAB-02_11_Windows` en la rama `experimental`.

Estas pruebas reducen el riesgo de repetir los errores de LAB-02.10, pero siguen siendo necesarias las comprobaciones visuales y la prueba del ejecutable en el PC del usuario antes de valorar una publicación estable.
