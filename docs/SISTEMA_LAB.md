# Sistema de desarrollo LAB — Selector Cine 365

**Inicio: 22 de noviembre de 2026 · primera serie: LAB-01.**

## Política obligatoria

1. `main` es **producción**. Se conserva como rama predeterminada y no se escribe en ella durante las pruebas. La referencia pública estable sigue siendo el tag y la release [v9.0.0](https://github.com/cristianbarallobre/selector-cine-365/releases/tag/v9.0.0), incluidos sus ejecutables originales.
2. `experimental` es el canal de ensayos. Cualquier nueva función, corrección experimental o refactorización se desarrolla exclusivamente aquí o en ramas derivadas (`lab/<tema>`).
3. Se numeran los ensayos **LAB-01, LAB-02, ...**, sin mezclar el concepto de LAB con la numeración de releases estables.
4. No se crean automáticamente tags ni releases desde LAB. El workflow `.github/workflows/ci-lab.yml` tiene `contents: read`, solo prueba y sube artefactos temporales de GitHub Actions.
5. Solo con aprobación expresa se incorpora a producción **cada cambio aceptado**. Preferir rama de publicación nueva a partir de `main` y cherry-pick/PR con los commits aprobados; no fusionar a ciegas todo `experimental`.
6. Antes de aprobar: pruebas automáticas y manuales, regresión de listas y sincronización, y copia JSON local; conservar la versión oficial y sus hashes sin alteraciones.
7. No copiar secretos, datos reales, credenciales de OneDrive ni listas del usuario de la versión estable a LAB sin decisión explícita.

### Situación inicial congelada

- Rama principal: `main`, commit inicial `1fbbe89361f03dc548332611775aa794bb949009`.
- Release pública que no se modifica: `v9.0.0`.
- Rama de ensayo: `experimental`, nacida de ese commit.
- Los archivos y la release estables **no se han modificado** al implantar LAB.

> **Protección de GitHub pendiente de activar en ajustes:** actualmente `main` no tiene protección configurada. El propietario debe crear un ruleset/branch protection para `main`: bloquear force-push y borrado, exigir pull request y comprobaciones, y no permitir merges automáticos. Esto no lo puede imponer el código de la rama experimental. Mientras tanto, respetar la prohibición de escribir en `main`.

## Aislamiento de ejecución

| Componente | Estable | LAB |
| --- | --- | --- |
| Fuente Git | `main` | `experimental` / `lab/*` |
| Windows datos | `%APPDATA%\\Selector_Cine_365` | `%APPDATA%\\Selector_Cine_365_LAB` |
| macOS datos | `~/Library/Application Support/Selector_Cine_365` | `~/Library/Application Support/Selector_Cine_365_LAB` |
| OneDrive directo escritorio | `Selector_Cine_365` | `Selector_Cine_365_LAB` |
| Copias de nube | `Selector_Cine_365_Copias` | `Selector_Cine_365_LAB_Copias` |
| Android paquete | `es.cine365.app.limpia` | `es.cine365.app.lab` |
| Android datos | Almacenamiento privado estable | Almacenamiento privado LAB |
| macOS app y bundle ID | Selector Cine 365 | Selector Cine 365 LAB (`es.cristianbarallobre.selectorcine365.lab`) |

**Android LAB:** el acceso Microsoft/OneDrive permanece **desactivado intencionadamente**. No se puede reutilizar el callback Microsoft del APK estable sin posibles conflictos al tener ambas apps instaladas. Para habilitar OneDrive en Android LAB más adelante debe configurarse y probarse un registro OAuth LAB independiente. No usar un APK LAB para probar sincronización cruzada Android-escritorio hasta entonces. La sincronización mediante carpetas elegidas manualmente tampoco debe apuntar nunca a la carpeta estable.

**Windows/macOS LAB:** el cliente Microsoft Graph utiliza la carpeta remota `Selector_Cine_365_LAB` y sesiones locales separadas. Cuando se use copia a carpeta, el subdirectorio de LAB queda diferenciado; **no seleccionar ni reutilizar la carpeta de datos estable como destino**.

## Compilar el modo LAB

Es imprescindible compilar con la opción LAB: **clonar la rama por sí solo no cambia las rutas** de los builds normales.

```bash
git fetch origin
git switch experimental

# Windows (desde Windows o mediante compilación cruzada)
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -tags lab -ldflags="-H windowsgui -s -w" -o Selector_Cine_365_LAB.exe .

# Android, desde android/
gradle assembleDebug -PcineLab

# macOS, en un equipo Mac con Xcode CLI Tools
MAC_ARCH=arm64 MAC_LABEL=Apple_Silicon CC="clang -arch arm64" bash scripts/build_macos_lab.sh
```

**GitHub Actions:** `Actions → LAB - pruebas y compilaciones independientes` ejecuta pruebas y compila Windows, Android y macOS (Intel y Apple Silicon) tras cada push a `experimental`. Los artefactos están en la ejecución de Actions, no en Releases oficiales. Los binarios Windows no tienen Authenticode, los de Mac van firmados ad hoc (sin notarización) y Android usa **firma de depuración** generada por el runner: para cambiar a otro APK LAB puede ser necesario desinstalar **solo LAB** y volverlo a instalar; jamás desinstalar la app estable para una prueba LAB. No se ha afirmado que la compilación esté validada hasta que el CI termine correctamente.

## Checklist antes de aceptar una funcionalidad

- [ ] `main` y `v9.0.0` continúan inalterados.
- [ ] LAB se abre con identidad visual LAB y sus datos locales no aparecen en producción.
- [ ] Android LAB puede convivir con `es.cine365.app.limpia` sin reemplazarla.
- [ ] Windows/macOS LAB no usan `Selector_Cine_365` en AppData ni OneDrive.
- [ ] No se reutilizan cuentas/sesiones/archivos de pruebas en la biblioteca oficial.
- [ ] Pruebas automáticas verdes en las plataformas afectadas.
- [ ] Prueba manual y conformidad expresa antes de integrar commits seleccionados.

**Nunca** publicar una build LAB como `latest`, ni sustituir descargas de [la release v9.0.0](https://github.com/cristianbarallobre/selector-cine-365/releases/tag/v9.0.0).
