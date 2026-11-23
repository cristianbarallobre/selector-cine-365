'use strict';
const fs=require('node:fs');
const assert=require('node:assert/strict');
const lab=fs.readFileSync('lab_discover_windows.html','utf8');
const common=fs.readFileSync('plantilla.html','utf8');
const source=fs.readFileSync('lab_discover_windows.go','utf8');
const workflow=fs.readFileSync('.github/workflows/ci-lab.yml','utf8');

assert.match(source,/\/\/go:build windows/,'Cambios únicamente en Windows');
assert.match(source,/LAB-02\.12 WINDOWS/,'Cabecera de la última versión');
assert.match(workflow,/Selector_Cine_365_LAB-02_12_Windows/,'Artefacto nuevo');
assert.match(common,/class="page-btn" id="filmPrev">‹ Anterior/,'Botón de la biblioteca de referencia');
assert.match(common,/class="page-btn" id="filmNext">Siguiente ›/,'Botón siguiente de referencia');
assert.match(lab,/class="page-btn" id="labPrevPage" type="button">‹ Anterior<\/button>/,
 'Anterior en Descubrir reutiliza la clase y texto de la biblioteca');
assert.match(lab,/class="page-btn" id="labNextPage" type="button">Siguiente ›<\/button>/,
 'Siguiente en Descubrir reutiliza la clase y texto de la biblioteca');
assert.match(lab,/class="page-label" id="labPageLabel"/,'Etiqueta de páginas con mismo estilo');
assert.match(lab,/class="page-btn" id="labDetailClose" type="button" aria-label="Cerrar ficha" title="Cerrar ficha">✕<\/button>/,
 'Cierre como cruz, con etiqueta accesible y el estilo de paginación');
assert.ok(!lab.includes('✕ Cerrar'),'No debe quedar el texto Cerrar al lado de la X');
assert.match(lab,/\.lab-detail-top \.page-btn\{[^}]*width:34px/,'La cruz mantiene botón compacto');
assert.match(lab,/\.header-side \.header-tag\{order:1\}/,'Versión arriba');
assert.match(lab,/\.header-side \.creator-credit\{order:2\}/,'Hecho por justo bajo versión');
assert.match(lab,/\.header-side \.saved\{order:3\}/,'Guardado debajo del autor');
assert.match(common,/class="saved" id="saveState"/,'Los estados de guardado originales se conservan');
assert.match(common,/class="creator-credit"/,'El crédito del autor original se conserva');
assert.match(lab,/\.lab-detail-actions select\{[^}]*cursor:pointer/,'Cursor de mano para elegir lista en Detalles');
assert.match(lab,/#labDiscoverView select:not\(:disabled\),#labDetailOverlay select:not\(:disabled\)/,
 'Cursor de mano para selects interactivos');
assert.match(lab,/#labDiscoverView button:not\(:disabled\),#labDetailOverlay button:not\(:disabled\)/,
 'Cursor de mano para botones interactivos');
assert.match(lab,/#labDiscoverView button:disabled,#labDetailOverlay button:disabled\{cursor:not-allowed\}/,
 'No dar cursor de mano a botones deshabilitados');
console.log('OK: LAB-02.12 orden de cabecera, paginación unificada, X y hover Windows');
