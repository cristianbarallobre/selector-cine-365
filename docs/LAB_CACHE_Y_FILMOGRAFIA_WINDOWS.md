# Windows LAB — Filmografía y almacenamiento de caché

## Filmografía de intérpretes

La ficha de intérprete conserva los créditos de películas de TMDB (se evita
duplicar identificadores). Se presentan las primeras 16 tras ordenar **toda**
la filmografía recibida y se puede alternar sin repetir peticiones:

- **Destacadas** (por defecto): 50 % de valoración bayesiana
  `(nota × votos + 6 × 120)/(votos + 120)` normalizada a 0–1; 35 %
  de prominencia aproximada a partir de `order` de reparto; 15 %
  de popularidad TMDB comprimida con logaritmo.
- **Mejor valoradas**: misma nota bayesiana, para que una puntuación con un
  único voto no domine sobre obras con muchos votos.
- **Más populares**: popularidad que devuelve TMDB.
- **Recientes**: fecha de estreno descendente.

En los empates, se utiliza la fecha de estreno y después el ID TMDB.
El orden de reparto no es una medida exacta de protagonismo, y no se presenta
la clasificación «Destacadas» como oficial de TMDB.

## Ubicaciones

- **Datos del usuario (permanentes)**:
  `%APPDATA%\Selector_Cine_365_LAB\` (Roaming). Se mantienen allí las listas,
  el TXT, la sincronización, las identificaciones manuales y los tokens de
  OneDrive cifrados con DPAPI. Nunca se borran al limpiar una caché.
- **Windows LAB, caché regenerable**:
  `%LOCALAPPDATA%\Selector_Cine_365_LAB\cache_lab\`.
  Aquí se guardan `tmdb_fichas_lab.json`, `plataformas_cache.json`,
  `carteles\`, diagnósticos LAB y el perfil WebView2 de LAB.
- **Versión estable**: no se cambia su ubicación de archivos ni se borran
  carpetas utilizadas por ella.

Cuando existe una caché anterior en Roaming, LAB la copia a LocalAppData
solo si no existe allí una copia más nueva. No borra el original, porque
la versión estable podría estar utilizándolo. Los carteles antiguos se
reutilizan puntualmente sin necesidad de duplicar todos al iniciar.

## Límites y conservación

| Tipo | Política Windows LAB |
|---|---|
| Fichas TMDB | 2.800 fichas y aproximadamente 14 MiB máximo de archivo; al crecer, expulsar menos usadas hasta objetivo ~11 MiB |
| Índice de títulos resueltos, regenerable | Máximo 5.500 entradas; no afecta identificaciones corregidas |
| Resultados de plataformas | Máximo 3.200 películas o 10 MiB; se descartan registros más antiguos |
| Carteles pequeños | Hasta 900 archivos o 48 MiB, y caducidad de 120 días |
| Diagnósticos de red | Rotación de archivo de ~256 KiB, más un histórico del mismo tamaño |
| WebView2 | Carpeta LocalAppData; caché de navegador gestionada por Chromium/WebView2, sin límite duro añadido por nuestra aplicación |

Las fichas usan caducidades existentes: metadatos 30 días, valoración/vídeos
7 días, disponibilidad 24 horas. La caducidad pide refrescar cuando es
necesario; la expulsión por límite libera espacio. Una ficha expulsada se
descarga otra vez al consultarla (si hay clave TMDB y conexión).

**Importante:** los presupuestos de las cachés propias no constituyen un
límite máximo del espacio total de `cache_lab`, porque WebView2 administra
su propio perfil. Si este perfil crece anormalmente, habrá que gestionarlo
con las APIs de limpieza de WebView2 evitando borrar cookies o credenciales
en uso. Los archivos antiguos de Roaming se mantienen por compatibilidad
y pueden ocupar espacio adicional hasta una retirada controlada.

## Pruebas

- `cache_store_test.go`: copia de datos sin sobreescritura,
  presupuesto de carteles y aislamiento de configuraciones.
- `lab_cache_budget_windows_test.go`: ruta LocalAppData, expulsión LRU y
  persistencia incluso cuando la caché alcance el máximo.
- `scripts/test_lab_person_sort_windows.cjs`: cuatro ordenaciones y
  navegación desde películas destacadas.
