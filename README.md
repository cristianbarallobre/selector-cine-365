# Selector de Cine 365 🎬

Aplicación de escritorio para **Windows 10/11 de 64 bits** que permite llevar un registro personal de películas, sortear el siguiente visionado y comprobar plataformas de streaming en España.

**Última versión del proyecto:** `v8.1` · **Tecnologías:** Go, HTML/CSS/JavaScript, Microsoft WebView2.

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

Si solo quieres utilizar el programa, **no necesitas Go ni descargar el código**. Descarga el ZIP para Windows de la sección [Releases](../../releases), extráelo y ejecuta el archivo `.exe`.

En esta publicación, la lista inicial contiene las 365 películas **sin marcas de vista personales**. Si ya usabas la aplicación, tus propias marcas seguirán almacenadas en AppData.

La aplicación utiliza el motor **Microsoft Edge WebView2 Runtime**. Si no está disponible, consulta la [página oficial de WebView2](https://developer.microsoft.com/microsoft-edge/webview2/).

La consulta de carteles, notas, identificaciones y plataformas necesita una clave personal de TMDB; se introduce desde la aplicación. La gestión local de las listas puede utilizarse sin ella.

## Código y compilación

Se necesita **Go 1.23 o posterior** para compilar desde el código fuente.

```bash
go test ./...
GOOS=windows GOARCH=amd64 go build -o Selector_Cine_365.exe .
```

En Windows puedes ejecutar `go build -o Selector_Cine_365.exe .` directamente. El recurso `rsrc_amd64.syso` aporta el icono de Windows a la compilación AMD64.

La interfaz web está en `plantilla.html`, la lista inicial en `lista_inicial.txt`, el servidor local en `main.go`, las colecciones en `collections.go`, la integración con TMDB en `providers.go` e `identity.go`, y la ventana WebView2 en `window_windows.go`.

## Datos personales y copias de seguridad

Los datos del usuario se guardan por defecto en:

```text
%APPDATA%\AtlasDelCine\
```

Entre otros archivos se utilizan `colecciones.json`, `mis_365_peliculas.txt`, `identidades_cine.json` y `plataformas_ajustes.json`. **No subas estos archivos ni tus claves de API a GitHub.** El contenido de la carpeta AppData pertenece a cada instalación local.

Antes de actualizar a la v8.1 es recomendable exportar una copia de seguridad: la conciliación de identificaciones puede modificar las marcas de visto de varias listas.

## Atribución

Este proyecto utiliza la API de [The Movie Database (TMDB)](https://www.themoviedb.org/), pero no está respaldado ni certificado por TMDB.

**This product uses the TMDB API but is not endorsed or certified by TMDB.**

La información de disponibilidad de plataformas procede de JustWatch a través de TMDB. La disponibilidad varía según el país y la fecha.

## Licencia

Por el momento el proyecto **no incluye una licencia de reutilización de código**. Si el repositorio se hace público y se desea permitir modificaciones o redistribución, debe añadirse una licencia adecuada y revisar las obligaciones relativas a imágenes, datos y marcas de terceros.

## Historial

Consulta [CHANGELOG.md](CHANGELOG.md) para ver los cambios de cada versión.