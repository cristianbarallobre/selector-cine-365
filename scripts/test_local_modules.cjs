'use strict';
const fs=require('node:fs'),vm=require('node:vm'),assert=require('node:assert/strict');
const page=fs.readFileSync('plantilla.html','utf8');
for(const marker of ['__CINE_LOCAL_BACKUP_MARKUP__','__CINE_LOCAL_BACKUP_SCRIPT__','__CINE_LOCAL_ACTIVITY_SCRIPT__']){
 assert.equal(page.split(marker).length,2,'Módulo ausente o duplicado: '+marker);
}
assert.doesNotMatch(page,/login\.microsoftonline\.com|\/api\/onedrive-desktop\/|\/api\/cloud\/(?:auto|config|sync|browse)|id="cloudGraphConnect"/);
const backup=fs.readFileSync('ui/local_backup.js','utf8'),activity=fs.readFileSync('ui/activity.js','utf8');
new vm.Script(backup);new vm.Script(activity);
assert.match(backup,/\/api\/extra-backup\//);
assert.match(activity,/\/api\/sync\/history/);
console.log('OK: interfaz local, módulos separados y sin controles de nube');
