'use strict';
const fs=require('node:fs');
const assert=require('node:assert/strict');
const cloud=fs.readFileSync('android/app/src/main/java/es/cine365/app/CloudBackup.java','utf8');
const graph=fs.readFileSync('android/app/src/main/java/es/cine365/app/OneDriveCloud.java','utf8');

// Source-level guards complement the Java tests and catch accidental regressions.
assert.match(cloud,/Uri own=findChild\(devices,snapshotName\);/);
assert.match(cloud,/if\(own==null\|\|!java\.util\.Arrays\.equals\(read\(own\),data\)\)\s*\{\s*putFile\(devices,snapshotName,data,true\);/);
assert.match(cloud,/if\(own==null\|\|!java\.util\.Arrays\.equals\(read\(own\),data\)\)\s*throw new Exception\("El proveedor no confirmó la escritura/);
assert.match(graph,/if\(list==null\)throw new Exception\("OneDrive devolvió un listado de archivos incompleto"\)/);
assert.match(graph,/if\(next\.isEmpty\(\)\)return all;/);
assert.match(graph,/throw new Exception\("OneDrive devolvió demasiadas páginas de archivos; sincronización cancelada/);
assert.doesNotMatch(graph,/JSONArray list=result\.optJSONArray\("value"\);if\(list==null\)break/);
console.log('OK: identical snapshots are not uploaded; incomplete OneDrive listings fail closed');
