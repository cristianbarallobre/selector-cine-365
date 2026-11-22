# Selector Cine 365 🎬

**Versión estable recomendada: v9.0.0 — 22 de noviembre de 2026.**

Selector Cine 365 es una aplicación de organización de películas para **Windows**, **macOS (Intel y Apple Silicon)** y **Android**. Permite administrar colecciones independientes, marcar películas vistas, buscar, sortear películas pendientes, consultar TMDB y, opcionalmente, sincronizar mediante OneDrive. También funciona completamente **sin cuenta Microsoft y sin clave TMDB**.

La versión estable reúne **exactamente las compilaciones que se probaron en los dispositivos**, sin modificar sus binarios:

- Windows y macOS: **Beta 14** (nombre de versión interno conservado).
- Android: **Beta 13.13 Limpia**, `versionCode 9027`, `es.cine365.app.limpia`, misma firma privada original de Beta 13.12 Limpia.

El **canal de distribución es estable v9.0.0**; los nombres «Beta 14» y «Beta 13.13» seguirán apareciendo dentro de la aplicación porque no se ha recompilado ninguna de ellas. Esto permite preservar exactamente lo que ya fue probado y comprobar sus huellas SHA-256.

## Descargar la versión estable

**[Abrir GitHub Releases — Selector Cine 365 v9.0.0 ESTABLE](https://github.com/cristianbarallobre/selector-cine-365/releases/tag/v9.0.0)**

| Plataforma | Descarga directa desde GitHub | Requisitos |
| --- | --- | --- |
| Windows 10/11 x64 | [Selector_Cine_365_Windows_Beta14.exe](https://github.com/cristianbarallobre/selector-cine-365/releases/download/v9.0.0/Selector_Cine_365_Windows_Beta14.exe) | WebView2; EXE portátil, sin instalador MSI |
| Android 8+ | [Selector_Cine_365_Android_Beta13_13_Limpia.apk](https://github.com/cristianbarallobre/selector-cine-365/releases/download/v9.0.0/Selector_Cine_365_Android_Beta13_13_Limpia.apk) | Paquete `es.cine365.app.limpia`, APK con firma original |
| Mac Apple Silicon | [Selector_Cine_365_Mac_AppleSilicon_Beta14.zip](https://github.com/cristianbarallobre/selector-cine-365/releases/download/v9.0.0/Selector_Cine_365_Mac_AppleSilicon_Beta14.zip) | macOS 12+, ZIP con `.app`, sin notarización |
| Mac Intel | [Selector_Cine_365_Mac_Intel_Beta14.zip](https://github.com/cristianbarallobre/selector-cine-365/releases/download/v9.0.0/Selector_Cine_365_Mac_Intel_Beta14.zip) | macOS 12+, ZIP con `.app`, sin notarización |

**La descarga es desde GitHub**, no desde un enlace temporal del chat. La publicación incluye `SHA256SUMS.txt`, los manuales PDF/DOCX y el código fuente ZIP.

### Actualizar sin perder datos

En **Windows**, cierra la aplicación y utiliza el nuevo EXE. Los datos residen en `%APPDATA%\Selector_Cine_365` y no en el ejecutable.

En **Android**, el paquete estable es **`es.cine365.app.limpia`**. Puede instalarse encima de Beta 13.12 Limpia, con la misma firma, sin desinstalar. **No** es una actualización de la antigua aplicación `es.cine365.app`, que tiene datos privados independientes. Nunca desinstales sin exportar previamente una copia JSON de las listas.

En **Mac**, elige el ZIP de tu arquitectura y reemplaza la aplicación conservando tu carpeta de datos. Estas compilaciones todavía no tienen firma/notarización Apple.

### Sincronización y carátulas: límites conocidos

- Cada lista debe conservar **sus propias películas**. Las versiones estables bloquean fusiones sospechosas de listas con IDs coincidentes. **No reparan automáticamente datos mezclados por versiones anteriores**: ante conteos inesperados, cancela la sincronización y recupera una copia JSON verificada.
- Android **sincroniza bajo confirmación manual**; Windows dispone de integración Microsoft Graph y puede sincronizar según las opciones seleccionadas. El modo local nunca exige OneDrive.
- Las identificaciones manuales de TMDB pueden viajar como IDs; las **carátulas y plataformas descargadas se guardan localmente** en cada dispositivo y pueden tardar en completarse. Las claves TMDB y los tokens de Microsoft no se comparten.
- Windows conserva la sesión Microsoft mediante DPAPI; Android usa su almacenamiento seguro y Mac todavía puede solicitar iniciar sesión tras reiniciar.

## Documentación vigente

- [Manual completo v9 estable (Markdown)](docs/MANUAL_COMPLETO_SELECTOR_CINE_365.md)
- [Manual completo en PDF](docs/MANUAL_COMPLETO_SELECTOR_CINE_365.pdf)
- [Guía de sincronización segura](docs/GUIA_SINCRONIZACION_SEGURA.md)
- [Configuración de Microsoft Entra y OneDrive](docs/GUIA_AZURE_ENTRA_ONEDRIVE_SELECTOR_CINE_365.md)
- [Historial de versiones](CHANGELOG.md) · [Política de releases](docs/POLITICA_DE_RELEASES.md)
- [Rutas de datos por plataforma](docs/UBICACIONES_DATOS_SELECTOR_CINE_365.md)
- [Notas de lanzamiento estable](RELEASE_NOTES_v9_ESTABLE.md)
- [Código fuente v9.0.0](https://github.com/cristianbarallobre/selector-cine-365/tree/v9.0.0)

**Histórico:** [Beta 14, versión de referencia previa](https://github.com/cristianbarallobre/selector-cine-365/releases/tag/v9.0.0-beta.14), [otras versiones conservadas](https://github.com/cristianbarallobre/selector-cine-365/releases) y [documentación histórica](docs/POLITICA_DE_RELEASES.md). El historial y los tags de Git se conservan aunque se retiren algunos instaladores públicos redundantes.

## Desarrollo y seguridad

El código de **Windows/macOS (Go + interfaz web)** y **Android (Java/Gradle + WebView)** está en este repositorio. El tag `v9.0.0` fija la fuente asociada al lanzamiento; los archivos de distribución conservan las huellas SHA-256 de las betas que fueron probadas.

**Nunca publiques** contraseñas, archivos de firma Android (`.p12`, `.jks`), claves TMDB ni tokens de OneDrive en el código o los informes de incidencias. La firma permanente Android requiere secretos protegidos de GitHub Actions para futuras compilaciones; la APK publicada utiliza el certificado original, sin incluir el material privado.

**Autor:** Cristian Barallobre Rúa.
