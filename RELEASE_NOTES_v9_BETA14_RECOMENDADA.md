
# Selector Cine 365 — v9 Beta 14

**21 de noviembre de 2026 · Beta recomendada para pruebas, NO versión estable definitiva.**

Esta publicación reúne los cuatro instaladores verificados en una sola GitHub Release y el código fuente combinado. Es la referencia de pruebas que reemplaza a las betas anteriores para nuevas instalaciones; las anteriores siguen disponibles para recuperación y consulta histórica.

## Versiones incluidas

| Plataforma | Versión del programa | Paquete |
| --- | --- | --- |
| Windows x64 (10/11) | Beta 14 | EXE portátil (sin instalador MSI ni firma Authenticode) |
| macOS Intel (12+) | Beta 14 | ZIP con aplicación .app, sin notarización Apple |
| macOS Apple Silicon (12+) | Beta 14 | ZIP con aplicación .app, sin notarización Apple |


## Correcciones estabilizadas


## Limitación de seguridad que debe conocer todo usuario

**Las protecciones no reparan automáticamente bibliotecas ya contaminadas por betas anteriores.** Si una lista tiene películas de otras listas, no confirms la sincronización y no borres los snapshots: primero exporta el JSON local bueno, compara sus listas y restaura desde una copia anterior verificada. **Una prueba automatizada satisfactoria no demuestra que una biblioteca personal de OneDrive esté intacta.**


## Instalación y conservación

1. Guarda un JSON independiente de las listas desde cada dispositivo antes de probar sincronización.
2. En Windows cierra el EXE anterior y abre el nuevo; la carpeta habitual es `%APPDATA%\Selector_Cine_365`. La sesión Microsoft DPAPI debería conservarse; si procede de una beta anterior a DPAPI, inicia sesión una vez.
4. En Mac extrae el ZIP adecuado. Mac no tiene persistencia Keychain de la cuenta Microsoft en esta beta.
5. Antes de aceptar OneDrive comprueba que el total **de cada lista**, no solo el total de listas, coincide con lo esperado. Si aparece una mezcla, pulsa Cancelar.
6. Puedes usar la aplicación sin Microsoft ni TMDB; todas las listas se guardan localmente.


