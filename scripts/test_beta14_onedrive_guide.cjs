'use strict';
// Regresiones de la guía de sincronización de escritorio Beta 14.
const assert=require('node:assert/strict');
const fs=require('node:fs');

const html=fs.readFileSync('plantilla.html','utf8');
const allScripts=[...html.matchAll(/<script>([\s\S]*?)<\/script>/g)].map(x=>x[1]);
assert.ok(allScripts.length);
for(const script of allScripts)assert.doesNotThrow(()=>new Function(script),'JS debe compilar');
const script=allScripts.at(-1);

const labels=[
 'Conecta OneDrive','Revisa y confirma','Sincroniza y comprueba',
 'Selector_Cine_365/sincronizacion/dispositivos',
 'Confirmar y sincronizar',
 'Otras formas de conexión y opciones avanzadas',
 'Actividad reciente, historial y copias de recuperación'
];
for(const label of labels)assert.ok(html.includes(label),'Falta texto o sección: '+label);
for(const id of [
 'cloudWizardStatus','cloudStepConnect','cloudStepReview','cloudStepSync',
 'cloudGraphConnect','cloudGraphStatus','cloudGraphDisconnect',
 'cloudSyncReview','cloudSyncNowFast','cloudSyncMode','cloudGraphLogin',
 'cloudAutoFolder','cloudFolder','cloudBrowse','cloudActivate','cloudBackupNow','cloudDisable',
 'cloudStatus','cloudAutoStatus','cloudLastCheck','cloudLastReceive','cloudSyncManualResult',
 'cloudHistory','cloudHistoryRefresh','cloudHistoryChecked','cloudRecoveryList','cloudRecoveryRestore',
 'cloudRecoveryStatus','cloudSyncPanel','cloudSyncSummary','cloudSyncConflicts','cloudSyncConfirm'
]){
 const count=[...html.matchAll(new RegExp('id="'+id+'"','g'))].length;
 assert.equal(count,1,'ID único para preservar funciones y eventos: '+id);
}
assert.ok(html.includes('id="cloudSyncReview" disabled'),'Revisión bloqueada antes de conectar');
assert.ok(html.includes('id="cloudSyncNowFast" disabled'),'Sync bloqueado antes de vincular');
assert.ok(html.includes('.cloud-guide-steps{display:grid;grid-template-columns:repeat(3,minmax(0,1fr))'),'Tres pasos en escritorio');
assert.ok(html.includes('@media(max-width:1000px){.cloud-guide-heading{flex-direction:column}'),'Pasos adaptables al ancho');

class Fake {
 constructor(){this.dataset={};this.textContent='';this.disabled=false;this.value='';this.attrs={};this.classList={toggle:(k,v)=>{this.attrs[k]=v;}};}
}
const ids=[
 'cloudWizardStatus','cloudStepConnect','cloudStepReview','cloudStepSync',
 'cloudSyncReview','cloudSyncNowFast','cloudReviewHint','cloudSyncHint',
 'cloudGraphConnect','cloudStatus','cloudFolder','cloudLastCheck','cloudLastReceive'
];
const elements=Object.fromEntries(ids.map(id=>[id,new Fake()]));
const $=id=>{assert.ok(elements[id],'Missing UI mock: '+id);return elements[id];};
const document={activeElement:null};
const from=script.indexOf('let cineGraphConnected=false;');
const to=script.indexOf('async function cloudRefresh()',from);
assert.ok(from>=0&&to>from);
const handlers=new Function('$','document',script.slice(from,to)+
 '\nreturn {cloudShow,updateCloudGuide};')($,document);

handlers.cloudShow({provider:'',enabled:false,auto_sync:true});
assert.match(elements.cloudWizardStatus.textContent,/Paso 1 de 3/);
assert.equal(elements.cloudSyncReview.disabled,true);
assert.equal(elements.cloudSyncNowFast.disabled,true);
assert.equal(elements.cloudStepConnect.dataset.state,'ready');
assert.match(elements.cloudStatus.textContent,/guardado local/);

handlers.cloudShow({provider:'onedrive-direct',enabled:true,sync_paired:false,needs_review:true,last_check_ms:0,last_received_ms:0});
assert.match(elements.cloudWizardStatus.textContent,/Paso 2 de 3/);
assert.equal(elements.cloudSyncReview.disabled,false);
assert.equal(elements.cloudSyncNowFast.disabled,true);
assert.equal(elements.cloudStepConnect.dataset.state,'done');
assert.equal(elements.cloudStepReview.dataset.state,'ready');
assert.match(elements.cloudStatus.textContent,/vacía/);

handlers.cloudShow({provider:'onedrive-direct',enabled:true,sync_paired:true,needs_review:false,last_check_ms:0,last_received_ms:0});
assert.match(elements.cloudWizardStatus.textContent,/Sincronización disponible/);
assert.equal(elements.cloudSyncReview.disabled,false);
assert.equal(elements.cloudSyncNowFast.disabled,false);
assert.equal(elements.cloudStepReview.dataset.state,'done');
assert.equal(elements.cloudStepSync.dataset.state,'ready');
assert.match(elements.cloudStatus.textContent,/dispositivos/);

handlers.cloudShow({provider:'onedrive-direct',enabled:true,sync_paired:true,needs_review:true,last_check_ms:0,last_received_ms:0});
assert.equal(elements.cloudSyncNowFast.disabled,true,'Conflictos bloquean un envío ciego');
assert.match(elements.cloudSyncHint.textContent,/conflictos/);

handlers.cloudShow({provider:'folder',enabled:true,sync_paired:false,needs_review:true,folder:'C:\\Users\\Una\\OneDrive\\Selector_Cine_365',last_check_ms:0});
assert.equal(elements.cloudGraphConnect.disabled,true,'No mezclar carpeta local y cuenta Graph');

const confirmation=script.slice(script.indexOf("$('cloudSyncConfirm').addEventListener("),script.indexOf('// Se activa solo tras la primera revisión manual.'));
assert.ok(!confirmation.includes('location.reload()'),'Sin recargar después de una sincronización confirmada');
assert.ok(confirmation.includes('await initLibrary();restartList();'),'Refrescar biblioteca sin terminar la sesión');
assert.ok(confirmation.includes('await cloudRefresh()'),'Actualizar estado tras confirmar');
assert.ok(confirmation.includes('La copia podría haberse enviado'),'Timeout no se presenta como pérdida definitiva');
assert.ok(script.includes('Se agotó el tiempo de espera'));
assert.ok(script.includes('La carpeta puede seguir vacía'));

console.log('OK: guía de 3 pasos, estados desconectado/pendiente/vinculado/conflicto, controles, historial y actualización sin recarga.');
