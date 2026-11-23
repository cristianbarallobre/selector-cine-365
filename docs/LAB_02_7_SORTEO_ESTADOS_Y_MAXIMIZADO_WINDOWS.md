# LAB-02.7 — Sorteo por estado de visionado y arranque maximizado (Windows)

**Canal:** `experimental` — Windows LAB. No modifica `main`, Android ni macOS. No incorpora Apariencia.

## Pantalla de la aplicación

La ventana nativa de **Selector Cine 365 LAB para Windows** se abre **maximizada** de forma predeterminada, ocupando el escritorio disponible pero conservando la barra de título, la barra de tareas y la opción de restaurar el tamaño normal de la ventana. No se utiliza un modo quiosco que impida salir.

En caso de que falle WebView2 y sea necesario arrancar el modo compatible con Edge/Chrome, también se solicita maximizar la ventana de respaldo.

La versión estable no cambia su tamaño de apertura.

## Sorteo inteligente

Al activar **Sorteo inteligente**, los controles son:

| Control | Opciones |
| --- | --- |
| Películas | **No vistas (pendientes)** · **Vistas** · **Vistas y no vistas** |
| Duración | Cualquier duración · Hasta 90 minutos · De 91 a 120 minutos · Más de 120 minutos |
| Género | Todos los géneros o un género real disponible de TMDB |
| Plataforma | El selector superior ya existente: todas las películas o las que cumplan un servicio de streaming seleccionado |

Se ha eliminado **«¿Qué te apetece?»**, junto con las opciones Equilibrado, Algo corto, Una buena valoración y Sorpréndeme. También se ha retirado el párrafo de aclaraciones situado debajo de Género; los filtros son suficientemente explícitos.

### Comportamiento

- El estado inicial del sorteo inteligente es **No vistas (pendientes)**, para conservar el comportamiento que ya tenía. Puede cambiarse a **Vistas** o **Vistas y no vistas** para volver a ver una película.
- Los filtros de duración y género se aplican **a cualquiera de los tres estados**. Si se exige duración concreta, las películas cuya duración TMDB no se conoce se excluyen, en lugar de asumir un tiempo.
- Se respeta siempre la selección superior de plataforma. **«Todas las plataformas»** conserva la semántica anterior: admite películas sin plataforma registrada; **«Cualquier plataforma (streaming)»** exige una oferta registrada.
- Se evitan las ocho últimas elecciones inteligentes cuando hay alternativas; no se guarda el historial en la nube.
- Se mantiene una ponderación suave de valoraciones y variedad de géneros a partir de la caché TMDB; no se introducen solicitudes externas nuevas por elegir películas.
- El contador de candidatas indica si está contabilizando pendientes, vistas o películas de ambos estados.
- **Detalles** sigue disponible en el resultado. Si sale una película ya vista, el botón de acción inferior ofrece **«Marcar como pendiente»**, en vez de «Marcar como vista»; si es pendiente, ofrece **«Marcar como vista»**.
- **El sorteo clásico no cambia:** solo considera películas pendientes y respeta la plataforma seleccionada.
- Los cambios de selección no alteran las marcas de visionado; solo se escriben en las listas cuando el usuario pulsa expresamente una acción para cambiar su estado.

## Pruebas

Incluidas pruebas automáticas de:
- Estados pendientes, vistas y ambos; combinaciones con género, plataforma y duraciones de 80, 85, 90, 91 y 238 minutos.
- Botón de marcado coherente con el estado de la película elegida.
- Integración HTML, ausencia de aclaraciones y del antiguo selector «¿Qué te apetece?».
- Arranque maximizado exclusivo de Windows LAB y respaldo Edge/Chrome.
- Regresiones Go y pruebas nativas de la caché y precarga TMDB.

Artefacto: **`Selector_Cine_365_LAB-02_7_Windows`** en [GitHub Actions](https://github.com/cristianbarallobre/selector-cine-365/actions/workflows/ci-lab.yml?query=branch%3Aexperimental).

La revisión visual final corresponde a las pruebas en Windows del usuario.
