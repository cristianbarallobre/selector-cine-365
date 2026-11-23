# LAB-02.9 — Conexiones más resistentes y unificación visual (Windows experimental)

**Fecha:** 23 de noviembre de 2026  
**Rama:** `experimental`  
**Plataforma:** Windows LAB exclusivamente. La versión estable (`main`), Android y macOS no reciben estos cambios.

## 1. Recuperación de fallos de red

Se refuerza la comunicación de la interfaz WebView2 con el servidor local de Selector Cine 365:

- Las consultas **GET/HEAD** realizadas a las rutas locales `/api/` se reintentan como máximo **dos veces** después de la primera, con intervalos de 300 y 900 ms.
- Solo se reintentan **fallos de conexión** y errores HTTP **408, 502, 503 y 504**. Los 400, 401, 403, 404, 429 y respuestas válidas no se reintentan, para no insistir ante entradas erróneas, credenciales incorrectas o límites de TMDB.
- Las peticiones externas a otros dominios, las comprobaciones `/api/ping` y las solicitudes canceladas no se reintentan.
- **POST, PUT, PATCH y DELETE no se reintentan nunca automáticamente.** Así evitamos duplicar escrituras, cambios de estado, sincronizaciones o confirmaciones de OneDrive si la respuesta se perdió después de que el servidor hubiera completado la operación.
- Si falla un guardado, se muestra **«Guardado no confirmado · comprueba la lista antes de cerrar»**, sin asegurar que el servidor no haya guardado el cambio.
- Los errores de OneDrive conservan un mensaje reconocible incluso si la respuesta HTTP del servidor local no contiene JSON.

### Diagnóstico local sin información sensible

El endpoint interno `POST /api/lab/network-log` requiere el token de sesión y solo recibe metadatos de error validados: ruta **sin parámetros ni consultas**, método, tipo de fallo, código HTTP e intentos. No recibe nombres de películas, contenido de listas, claves ni cabeceras.

Los errores finales se anotan en:

```text
%APPDATA%\Selector_Cine_365_LAB\diagnostico_red_lab.txt
```

El archivo reside en la carpeta de usuario de LAB y hereda los permisos de acceso (ACL) de esa carpeta en Windows; la ruta del registro requiere autorización para recibir eventos. Rota al alcanzar 256 KiB y conserva como máximo una copia anterior. El diagnóstico está separado de las listas y de OneDrive.

**Importante:** los reintentos reducen errores *transitorios*, pero no solucionan por sí solos caídas prolongadas de Internet, errores del servidor ni problemas permanentes de autorización.

## 2. Mejoras visuales

- El encabezado derecho muestra, en este orden: **Versión → Hecho por → Guardado**.
- Los botones **Anterior** y **Siguiente** de Descubrir usan la clase, dimensiones y estilo de los de paginación de Biblioteca.
- El botón de cierre de la ficha Detalles muestra **solo «✕»**, con etiqueta accesible «Cerrar ficha».
- Los controles interactivos de Descubrir y Detalles, incluidos desplegables como **«Añadir a la lista»**, muestran la **mano** al pasar el cursor; los botones desactivados no la muestran.
- Se conservan la selección de películas pendientes, vistas o ambas, los filtros de duración y género, Detalles y el resto de LAB-02.7.

## 3. Pruebas

La automatización de Windows ejecuta los tests existentes de biblioteca, caché, Descubrir y sorteo, más:

- `node scripts/test_lab_network_windows.cjs`: reintentos, códigos HTTP, escrituras no duplicadas, abortos y privacidad del diagnóstico.
- `node scripts/test_lab_visual_polish_windows.cjs`: navegación, cierre, orden de encabezado y cursor.
- `go test -tags lab -run TestLab -count=1 -v .` en Windows: rutas de diagnóstico con autenticación, validación, rotación y resultado HTML real, incluyendo los estilos de Descubrir.

El ejecutable se entrega como artefacto GitHub Actions **`Selector_Cine_365_LAB-02_9_Windows`**. No se mezcla con las versiones públicas de la rama estable.

## 4. Pruebas manuales recomendadas

1. Abrir LAB-02.9, comprobar **Versión → Hecho por → Guardado**.
2. Entrar en Descubrir, comprobar Anterior y Siguiente en páginas 1 y 2, con el estilo de Biblioteca.
3. Abrir Detalles y comprobar que el botón de cierre sea únicamente la cruz; revisar el cursor sobre **Añadir a la lista**.
4. Probar la biblioteca, la ficha de una película, el sorteo inteligente y el guardado de una marca.
5. Ante un error repetido de red, consultar `diagnostico_red_lab.txt` para obtener la ruta y el tipo de fallo, sin compartir archivos con datos personales.

**La validación visual final requiere probar el EXE en Windows real.**
