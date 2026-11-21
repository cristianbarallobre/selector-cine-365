# Selector Cine 365 — Android v9 Beta 13.3

**Actualización Android (versionCode 9017).** Windows y macOS conservan la última Beta 13.1 y sus propios ajustes de sincronización.

## Nuevo funcionamiento

- **Guardar en el móvil ≠ sincronizar con OneDrive.** Marcar Vista/Pendiente, añadir o eliminar películas, identificar con TMDB o cambiar de lista guarda el estado exclusivamente en Android. No se sube ninguna copia a OneDrive por esas acciones.
- **Al cambiar de lista** no aparece «Guardando…» ni se realiza ninguna consulta a OneDrive o TMDB. Si cambia la lista activa, su selección se conserva silenciosamente en el móvil.
- **Solo se consulta OneDrive al abrir la aplicación (una vez).** Si tienes una carpeta conectada y previamente vinculada, se realiza una *vista previa de solo lectura*: puede señalar posibles conflictos, pero **no aplica cambios ni publica archivos sin autorización**. En modo local, no intenta conectarse a OneDrive.
- **Nuevo botón compacto con dos flechas circulares** junto al indicador de guardado, visible desde la cabecera. Abre el apartado Nube y prepara una comparación bidireccional de los datos locales y remotos. **Es necesario revisar y confirmar** el resumen antes de recibir y enviar cambios; los conflictos exigen elección expresa.
- **Sin comprobaciones periódicas o al volver de segundo plano.** Se eliminan las consultas automáticas por temporizador, al cambiar de foco y al reanudar.
- Los cambios locales siguen protegidos con guardado automático en el almacenamiento del teléfono. Un cambio pendiente al minimizar se guarda, pero nunca se envía sin pulsar el botón manual.
- Se evita reescribir el TXT de las 365 películas cuando no ha cambiado, y no se escribe una biblioteca idéntica de nuevo.
- Los botones específicos de la sección Nube para enviar copia o revisar sincronización siguen siendo acciones **manuales**.
- Android conserva la protección de notch, estadísticas, marcadores causales, lista eliminada mediante tombstone, identificación TMDB manual y clave API privada.

## Instalación

- `applicationId es.cine365.app`, versión `9.0.0-beta.13.3`, `versionCode 9017`.
- **Mismo certificado privado que Android Beta 11–13.2.** Instala la APK firmada encima de Beta 13.2 (o anteriores compatibles) **sin desinstalar** ni borrar los datos.
- Realiza una copia JSON de seguridad desde el móvil y, si sincronizas, también desde Windows antes de actualizar.
- Windows/Mac Beta 13.1 siguen comprobando su carpeta OneDrive según su propia configuración. Android ya no importa cambios automáticamente: después de editar en el PC, pulsa ↻ en el móvil para recibirlos, y viceversa. La carpeta elegida y las credenciales privadas se conservan.

## Comprobaciones

La compilación y las pruebas cubren el guardado local sin tráfico de red, ausencia de sincronización por intervalo o reanudación, comprobación de nube de solo lectura al abrir, botón con revisión y confirmación, sincronización causal, guardado al minimizar y compatibilidad con la versión anterior.

**Limitación:** no se garantiza un porcentaje fijo de ahorro de batería. Hace falta medir varias horas de uso real. No usar el icono manual no crea automáticamente una copia nueva en OneDrive: antes de cambiar de móvil, perder datos o reinstalar, recuerda pulsarlo y confirmar o exportar JSON.

© 2026 Cristian Barallobre Rúa.
