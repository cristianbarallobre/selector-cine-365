# Selector Cine 365 — Android v9 Beta 13.5

**Android `es.cine365.app` · `versionCode 9019` · Beta experimental.**

## OneDrive sin selector de carpetas

- Después de conectar la cuenta Microsoft, Android reutiliza la carpeta OneDrive ya vinculada si existe.
- En una instalación sin carpeta vinculada, busca o crea **`Selector_Cine_365`** en la raíz de esa cuenta y la vincula, sin publicar ni importar datos por el mero hecho de conectarla.
- **No cambia ni migra automáticamente una configuración anterior de proveedor SAF** (carpeta local/OneDrive en explorador): se muestra aviso y se mantiene el estado anterior. La selección manual de una carpeta diferente queda como opción avanzada y requiere exportar JSON previamente.
- **Windows y Mac deben usar esa misma carpeta de OneDrive**: el cliente de escritorio conserva su mecanismo de carpeta local sincronizada y puede requerir vincularla una sola vez. No se ha añadido autenticación Graph nativa a escritorio.
- La subida y recepción de películas siguen siendo **manuales**: pulsar ↻, revisar conflictos y confirmar. La app guarda inmediatamente las modificaciones en local sin subirlas automáticamente.
- El código de Beta 13.4 para evitar subidas idénticas y no aceptar listados remotos incompletos se conserva.
- Antes de instalar, exportar copia JSON local, no desinstalar la versión anterior, e instalar la APK **firmada con la clave original**. La APK de CI sin firmar no sirve como actualización.

## Pruebas imprescindibles en móvil

1. Instala encima de Beta 13.3/13.4, confirma que las listas y marcas persisten.
2. En Ajustes → OneDrive verifica que aparece cuenta conectada y carpeta preparada, sin selector obligatorio.
3. Sincroniza mediante ↻, revisa diferencias y acepta solo si esperadas.
4. Marca una película sin pulsar ↻: comprueba que sigue guardada localmente y no se sube hasta confirmación.
5. Comprueba en Windows/Mac con la **misma carpeta** y confirma sincronización sin duplicados.
6. Si utilizabas una carpeta distinta, **no migres sin respaldo JSON**: comprueba ruta anterior y utiliza opciones avanzadas para seleccionar manualmente la ruta común.

## Alcance

Esta beta simplifica el vínculo OneDrive en **Android**; no instala OneDrive en escritorio ni modifica el motor de Windows/Mac. No hay servidor propio ni migraciones destructivas de bibliotecas.

© 2026 Cristian Barallobre Rúa.
