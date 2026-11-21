'use strict';
const assert=require('node:assert/strict');
const fs=require('node:fs');

const read=name=>fs.readFileSync(name,'utf8');
const manifest=read('android/app/src/main/AndroidManifest.xml');
const rulesOld=read('android/app/src/main/res/xml/backup_rules.xml');
const rulesNew=read('android/app/src/main/res/xml/data_extraction_rules.xml');
const build=read('android/app/build.gradle');
const html=read('android/app/src/main/assets/index.html');
const backend=read('android/app/src/main/java/es/cine365/app/Backend.java');
const cloud=read('android/app/src/main/java/es/cine365/app/CloudBackup.java');
const drive=read('android/app/src/main/java/es/cine365/app/OneDriveCloud.java');

assert.match(manifest,/android:allowBackup="false"/,'No automatic restore on a fresh install');
assert.match(manifest,/android:fullBackupContent="@xml\/backup_rules"/,'Android 8-11 backup policy');
assert.match(manifest,/android:dataExtractionRules="@xml\/data_extraction_rules"/,'Android 12+ transfer policy');
for(const xml of [rulesOld,rulesNew]){
  for(const type of ['root','device_root','external']){
    assert.ok(xml.includes('<exclude domain="'+type+'" path="." />'),'Private files excluded from automatic restore: '+type);
  }
}
assert.match(rulesNew,/<cloud-backup>[\s\S]*?<\/cloud-backup>/);
assert.match(rulesNew,/<device-transfer>[\s\S]*?<\/device-transfer>/);
assert.match(build,/applicationId 'es\.cine365\.app'/);
assert.match(build,/versionCode 9025/);
assert.match(build,/versionName '9\.0\.0-beta\.13\.11'/);
assert.match(backend,/new File\(activity\.getFilesDir\(\),"Selector_Cine_365"\)/,'Android local files in private storage');
assert.match(drive,/AndroidKeyStore/,'Microsoft token protected by Android Keystore');
assert.match(drive,/getSharedPreferences\("onedrive"/,'Microsoft refresh token stored securely');
assert.match(html,/1 · Conecta tu cuenta Microsoft/);
assert.match(html,/2 · Revisa las listas que ya tienes en OneDrive/);
assert.match(html,/3 · Sincroniza cuando quieras/);
assert.match(html,/copias de un Android anterior/);
assert.match(html,/esta versión no solicita restaurarlo desde copias automáticas de Android/);
assert.match(html,/no envía cambios automáticamente al abrir ni en segundo plano/);
assert.match(html,/No se sustituye ninguna copia de Windows/);
assert.match(html,/«Confirmar y sincronizar»/);
assert.ok(!html.includes('Después recibirás los cambios seguros automáticamente.'));
assert.ok(!html.includes('En Windows abre OneDrive y elige la carpeta'));
assert.ok(html.includes('Selector_Cine_365_Copias'));
assert.match(html,/id="syncConfirmButton"/);
assert.match(html,/id="syncPreviewPanel" hidden/);
for(const [,code] of html.matchAll(/<script>([\s\S]*?)<\/script>/g))assert.doesNotThrow(()=>new Function(code));
const planStart=cloud.indexOf('private Plan buildPlan('), previewStart=cloud.indexOf('synchronized JSONObject preview(',planStart);
const confirmStart=cloud.indexOf('synchronized JSONObject synchronize(',previewStart);
assert.ok(planStart>=0&&previewStart>planStart&&confirmStart>previewStart);
assert.ok(!cloud.slice(planStart,confirmStart).includes('putFile('),'Preview and plan may not publish or overwrite OneDrive');
assert.match(cloud,/if\(folder\.isEmpty\(\)\|\|!activated\)return status\(\);/,'No upload until initial pairing confirmed');
assert.match(cloud,/if\(previewToken\.isEmpty\(\)\|\|!previewToken\.equals\(expected\)\)/,'Confirmation requires fresh preview token');
assert.match(cloud,/rescueLocal\(previousLocal\)/,'Save local recovery before first import');
assert.match(cloud,/String snapshotName=sync\.deviceId\(\)\+"\.json"/,'Uploads own device snapshot only');
assert.match(backend,/\/\/ Only save locally\. Upload\/download requires explicit user action\./);
assert.ok(!html.includes('setInterval(()=>void syncDevices()'),'Never sync periodically on mobile');
console.log('OK: Beta 13.11 fresh install, Android auto-restore disabled, three-step OneDrive preview, no silent uploads.');
