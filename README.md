# Selector de Cine 365 🎬

Aplicación de escritorio para **Windows 10/11 y macOS (Apple Silicon e Intel)** que permite llevar un registro personal de películas, sortear el siguiente visionado y comprobar plataformas de streaming en España.

**Última versión del proyecto:** `v8.2` · **Tecnologías:** Go, HTML/CSS/JavaScript, Microsoft WebView2 (Windows) y WKWebView (macOS).

**Creado por Cristian Barallobre Rúa.** © 2026. Todos los derechos reservados.

## Funciones principales

- Crear e importar varias listas de películas y marcar títulos como vistos o pendientes.
- Buscar, filtrar y navegar por páginas en cada colección.
- Sortear películas pendientes, con filtros por servicios de streaming.
- Mostrar carteles, notas y disponibilidad en España a través de TMDB y datos de plataformas procedentes de JustWatch.
- Comparar listas mediante una colección virtual de películas coincidentes.
- Reconocer títulos equivalentes utilizando IDs de TMDB y sincronizar las marcas de visto entre listas.
- Consultar estadísticas de visionado, corregir identificaciones de TMDB y eliminar películas con confirmación.
- Guardar los datos localmente en AppData y hacer una segunda copia en una carpeta sincronizada que configure el usuario.

## Descargar y usar

Si solo quieres utilizar el programa, **no necesitas Go ni descargar el código**. Desde [Releases](../../releases/latest) descarga el ZIP correspondiente a **Windows**, **Mac Apple Silicon** o **Mac Intel**. En Windows ejecuta el `.exe`; en Mac, mueve la `.app` a Aplicaciones. [Instrucciones para macOS](MACOS.md).

En esta publicación, la lista inicial contiene las 365 películas **sin marcas de vista personales**. Si ya usabas la aplicación, tus propias marcas seguirán almacenadas en AppData.

Windows utiliza **Microsoft Edge WebView2 Runtime**; macOS utiliza **WKWebView**, incluido en el sistema. Si no está disponible, consulta la [página oficial de WebView2](https://developer.microsoft.com/microsoft-edge/webview2/).

La consulta de carteles, notas, identificaciones y plataformas necesita una clave personal de TMDB; se introduce desde la aplicación. La gestión local de las listas puede utilizarse sin ella.

## Código y compilación

Se necesita **Go 1.23 o posterior** para compilar desde el código fuente.

```bash
go test ./...
GOOS=windows GOARCH=amd64 go build -o Selector_Cine_365.exe .
```

En Windows puedes ejecutar `go build -o Selector_Cine_365.exe .` directamente. El recurso `rsrc_windows_amd64.syso` aporta el icono de Windows a la compilación AMD64.

La ventana nativa de macOS está en `window_darwin.go` y `macos_view_darwin.m`, y el selector de carpeta en `cloud_picker_darwin.go`. La interfaz web está en `plantilla.html`, la lista inicial en `lista_inicial.txt`, el servidor local en `main.go`, las colecciones en `collections.go`, la integración con TMDB en `providers.go` e `identity.go`, y la ventana WebView2 en `window_windows.go`.

## Datos personales y copias de seguridad

Los datos del usuario se guardan por defecto en:

```text
Windows: %APPDATA%\AtlasDelCine\
macOS:   ~/Library/Application Support/AtlasDelCine/
```

Entre otros archivos se utilizan `colecciones.json`, `mis_365_peliculas.txt`, `identidades_cine.json` y `plataformas_ajustes.json`. **No subas estos archivos ni tus claves de API a GitHub.** El contenido de la carpeta AppData pertenece a cada instalación local.

Antes de actualizar desde una versión anterior a la v8.1 es recomendable exportar una copia de seguridad: la conciliación de identificaciones puede modificar las marcas de visto de varias listas.

## Autoría del proyecto

**Autor:** Cristian Barallobre Rúa. La app muestra la firma en la cabecera. El repositorio incorpora también metadatos de autoría en `plantilla.html`, `autoria.go` y en las respuestas HTTP del servidor local. Estos datos no son una medida anticopia ni equivalen a una prueba jurídica concluyente; GitHub conserva el historial de cambios.

Para más información, consulta [AUTORIA.md](AUTORIA.md).

## Atribución

Este proyecto utiliza la API de [The Movie Database (TMDB)](https://www.themoviedb.org/), pero no está respaldado ni certificado por TMDB.

**This product uses the TMDB API but is not endorsed or certified by TMDB.**

La información de disponibilidad de plataformas procede de JustWatch a través de TMDB. La disponibilidad varía según el país y la fecha.

## Licencia

Por el momento el proyecto **no incluye una licencia de reutilización de código**. Si el repositorio se hace público y se desea permitir modificaciones o redistribución, debe añadirse una licencia adecuada y revisar las obligaciones relativas a imágenes, datos y marcas de terceros.

## Historial y documentación oficial

El historial completo, actualizado hasta la **v8.1.1** y con el crédito de **Cristian Barallobre Rúa**, está disponible en los siguientes formatos:

- [Leer el historial en GitHub](docs/HISTORIAL_VERSIONES.md).
- [Descargar historial en PDF](docs/Selector_de_Cine_365_Historial_de_versiones_v8.1.1.pdf).
- [Descargar historial editable en Word](docs/Selector_de_Cine_365_Historial_de_versiones_v8.1.1.docx).

Los documentos Word y PDF publicados en `docs/` se generan automáticamente a partir del historial en Markdown mediante GitHub Actions. Consulta también [CHANGELOG.md](CHANGELOG.md) para ver el resumen de los cambios de cada versión.

**Descargas (Windows, Mac Apple Silicon y Mac Intel):** [Selector de Cine 365 v8.2](https://github.com/cristianbarallobre/selector-cine-365/releases/tag/v8.2).
## macOS

Consulta [instrucciones de instalación, almacenamiento y seguridad para Mac](MACOS.md). Las apps de Mac tienen firma ad hoc, pero no están notarizadas por Apple. Se requiere una prueba práctica de interfaz en un Mac real.
