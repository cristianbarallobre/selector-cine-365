# LAB-03 — Apariencia personalizable (solo Windows)

**Rama:** `experimental` · **Canal:** Windows LAB · **No se publica en `main`**.

## Novedades

La nueva pestaña **🎨 Apariencia** aparece junto a Instrucciones y permite personalizar el fondo completo de la ventana sin modificar la biblioteca:

- **Oscuro original**: recupera la apariencia habitual, sin fondo añadido ni transparencias.
- **Cinematográfico**: luces cálidas sobre tonos oscuros.
- **Espacial**: luces y estrellas sobre tonos azulados y violetas.
- **Clásico**: efecto de cortinaje rojo oscuro.
- **Tu propia imagen**: seleccionar una fotografía **JPG o PNG** del PC (máximo **12 MB**; dimensiones de 100×100 a 8192×8192 píxeles). Se valida realmente su formato antes de guardarla. Se puede volver a elegir la imagen guardada después de probar otros temas.

Los fondos predeterminados son **gradientes integrados**, por lo que no se necesitan descargas ni servicios externos. La imagen propia se visualiza como fondo de ventana; los paneles centrales tienen un fondo oscuro translúcido para mantener la lectura y el fondo se aprecia mejor en los laterales.

Controles en Apariencia:

- **Opacidad de los paneles:** de 60 % a 100 %. A 100 % los paneles quedan sólidos.
- **Desenfoque de la imagen:** de 0 a 16 píxeles.
- **Vista previa:** refleja el fondo y la transparencia elegidos.
- **Restaurar apariencia original:** vuelve a la configuración predeterminada; no elimina la imagen que hayas subido.
- Los cambios se aplican inmediatamente y se guardan para las siguientes aperturas.

## Seguridad y almacenamiento

Se guardan dos archivos privados en la carpeta local de LAB (no en los archivos de listas):

```text
%APPDATA%\Selector_Cine_365_LAB\apariencia_lab.json
%APPDATA%\Selector_Cine_365_LAB\fondo_personalizado_lab.bin
```

- La app solo guarda estos archivos localmente; **la función no los sube a TMDB, Microsoft ni OneDrive**.
- Las rutas de la API local `/api/lab/appearance` y `/api/lab/appearance/wallpaper` requieren el token de la sesión.
- Se limita el tamaño de imagen, se verifica que sea JPEG o PNG y se restringen sus dimensiones antes de guardar.
- No se modifica `plantilla.html`, ni las listas, ni los datos de sincronización.
- El código solo se activa en **Windows compilado con `-tags lab`**. Existe una implementación inactiva para Mac/Linux. Android no contiene este menú.
- La app Windows oficial y la release estable siguen conservando su fondo original.

## Cómo probar

1. Descarga el artefacto **`Selector_Cine_365_LAB-03_Windows`** desde [GitHub Actions de `experimental`](https://github.com/cristianbarallobre/selector-cine-365/actions/workflows/ci-lab.yml?query=branch%3Aexperimental).
2. Cierra LAB-02.4 y abre el EXE nuevo, sin borrar ni mover tus carpetas LAB.
3. Entra en **🎨 Apariencia**; prueba los fondos Cinematográfico, Espacial y Clásico.
4. Ajusta opacidad y desenfoque. Cambia a Biblioteca y Descubrir para comprobar la lectura y el fondo de los laterales.
5. Carga una foto JPG/PNG, vuelve al fondo Oscuro original y después selecciona **Usar mi imagen guardada**.
6. Cierra el programa y vuelve a abrirlo: deben recuperarse fondo, opacidad y desenfoque. Prueba también **Restaurar apariencia original**.

## Validación

- `node scripts/test_lab_appearance_windows.cjs`: sintaxis JS, navegación, aislamiento y controles.
- `go test -tags lab -run TestLab -count=1 -v .` en Windows: API, permisos, imagen JPG/PNG, preferencias guardadas e inyección real de HTML.
- Compilación cruzada Windows LAB y tests existentes de Descubrir, biblioteca y sincronización.
- **La revisión visual real en tu Windows queda pendiente hasta que pruebes el ejecutable.** No integrar en `main` sin tu aprobación.

**Importante:** el fondo personalizado se conserva en LAB y no se migra a la app estable. Si en el futuro se aprueba portar Apariencia a otras plataformas, se planificará de manera separada.
