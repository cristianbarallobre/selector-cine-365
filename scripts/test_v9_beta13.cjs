'use strict';
const assert=require('node:assert/strict'),fs=require('node:fs');
function source(p){
 const s=fs.readFileSync(p,'utf8');
 const scripts=[...s.matchAll(/<script>([\s\S]*?)<\/script>/g)].map(m=>m[1]);
 assert.ok(scripts.length,'JS existente en '+p);
 for(const code of scripts)assert.doesNotThrow(()=>new Function(code),'Cada script JS debe compilar '+p);
 const script=scripts[scripts.length-1];
 for(const needle of ['cloudAutoReceive','cineRecentIntents','cineLoadRevisions','film-check-action','togg.append(tick,poster,main)','data-search-scope="all"','filmMembershipIndex','applyFilmMembershipLists','membershipOverlay','Conectar TMDB']){
  assert.ok(s.includes(needle),'Falta en '+p+': '+needle);
 }
 return {html:s,script};
}
const desktop=source('plantilla.html'),android=source('android/app/src/main/assets/index.html');
assert.ok(desktop.html.includes('data-view="instructions"'),'Windows debe conservar Instrucciones');
assert.ok(android.html.includes('data-view="help"'),'Android debe conservar Ayuda');
for(const {html,script} of [desktop,android]){
 assert.ok(!html.includes('main.inert=true'),'No bloquear toda la interfaz');
 assert.ok(script.includes('const readonly=virtualMode&&!global'),'Las coincidencias son solo lectura, la búsqueda global editable');
 assert.ok(script.includes('if(f._global)tick.dataset.globalGroup=f._group'),'Marcar visto con tick global');
 assert.ok(script.includes('else tick.dataset.idx=f.idx'),'Marcar visto con tick local');
 assert.ok(script.includes('if(!global&&selectedListPlatforms.size)'),'Filtro streaming no oculta búsqueda global');
 assert.ok(script.includes("num.addEventListener('click',()=>openFilmInList(m))"),'Números navegables');
}
assert.ok(desktop.script.includes('body.appendChild(line)'),'Estadísticas con cuerpo independiente');
assert.ok(desktop.html.includes('#view-statistics .stat-pager{flex:0 0 32px'),'Pie fijo de estadísticas Windows');
assert.ok(android.html.includes('Android: estadísticas completas sin paginación'),'Mantener estadísticas Android sin paginar');
assert.ok(android.html.includes('window.CineAndroid?Math.min(2000,'),'Mantener carga incremental Android');
const gradle=fs.readFileSync('android/app/build.gradle','utf8');
assert.ok(gradle.includes("applicationId 'es.cine365.app'"),'Identificador estable Android');
assert.ok(!gradle.includes('applicationIdSuffix'),'No generar un paquete distinto en cada Beta');
assert.ok(gradle.includes('versionCode 9015'),'Versión de Android incrementada');
const script=desktop.script,start=script.indexOf('function applyFilmMembershipLists('),end=script.indexOf('function toggleSearchedFilm(',start);
assert.ok(start>=0&&end>start);
function parse(t){const lines=t.split('\n'),entries=[];
 lines.forEach((l,idx)=>{const m=l.match(/^([☐☑])\s+(.*?)(?:\s+\((\d{4})\))?$/);
  if(m)entries.push({idx,checked:m[1]==='☑',title:m[2],name:m[2],year:m[3]||''});
 });return {lines,entries};
}
const key=f=>f.title+'|'+f.year;
const apply=new Function('parse','groupKey',script.slice(start,end)+';return applyFilmMembershipLists;')(parse,key);
const lists=[
 {id:'1',name:'Primera',text:'☑ A (2024)\n☐ B (2025)'},
 {id:'2',name:'Segunda',text:'☐ A (2024)'},
 {id:'3',name:'Tercera',text:'☐ C (2018)'}
];
const ids=lists.map(x=>x.id);
assert.equal(apply(lists,{title:'A',year:'2024',checked:true},new Set(['2','3']),new Map()),2);
assert.equal(lists[0].text,'☐ B (2025)','No modifica otras películas');
assert.equal(lists[1].text,'☐ A (2024)','No duplica');
assert.equal(lists[2].text,'☐ C (2018)\n☑ A (2024)','Conserva estado visto');
assert.deepEqual(lists.map(x=>x.id),ids,'No modifica IDs');
const a=script.indexOf('function filmMembershipIndex('),b=script.indexOf('function numberedListName(',a);
const index=new Function('libraryState','parse','buildIdentityGroups','groupKey','let membershipCache=null;'+script.slice(a,b)+';return filmMembershipIndex;')(
 {lists:[{id:'a',name:'A',text:'☐ Viaje (2020)'},{id:'b',name:'B',text:'☑ Journey (2020)'}]},parse,()=>new Map(),()=>314
)();
assert.equal(index.rows.length,1,'Una ficha por identidad TMDB');
assert.deepEqual(index.rows[0]._lists.map(x=>x.n),[1,2],'Listas numeradas');
assert.equal(index.rows[0].checked,true,'Estado global visto');
assert.deepEqual(index.rows[0]._searchNames,['Viaje','Journey'],'Alias en búsqueda global');
console.log('OK: sintaxis Windows+Android, sincronización preservada, búsqueda y listas y estadísticas');

// Regresiones visuales de Beta 12 (captura real comprobada con navegador).
assert.ok(desktop.html.includes('#view-collections .collection-rows{display:grid;grid-template-columns:minmax(0,1fr)'),'Las listas siempre han de tener altura visible');
assert.ok(desktop.script.includes('const visible=window.innerWidth<=760?3:4;'),'Paginación de listas independiente del panel de nube');
assert.ok(desktop.html.includes('id="cloudHistoryChecked"'),'La consulta de historial informa de su resultado');
assert.ok(desktop.script.includes("button.textContent='Consultando…'"),'Botón de historial indica operación en curso');
assert.ok(desktop.script.includes("target.classList.toggle('is-empty',items.length===0)"),'El estilo vacío de historial no añade una caja');
console.log('OK: regresiones de altura de listas y consulta de historial Beta 12');


// Nuevo: comprobar el recorte por las áreas seguras nativas de Android 15+
const mainActivity=fs.readFileSync('android/app/src/main/java/es/cine365/app/MainActivity.java','utf8');
const manifest=fs.readFileSync('android/app/src/main/AndroidManifest.xml','utf8');
assert.ok(mainActivity.includes('new FrameLayout(this)'),'WebView debe vivir en un contenedor seguro');
assert.ok(mainActivity.includes('setDecorFitsSystemWindows(false)'),'Android 11+ debe controlar edge-to-edge');
assert.ok(mainActivity.includes('WindowInsets.Type.systemBars() | WindowInsets.Type.displayCutout()'),'Combinar barras y recorte real');
assert.ok(mainActivity.includes('view.setPadding(bounds.left, bounds.top, bounds.right, bounds.bottom)'),'El viewport respeta los cuatro bordes');
assert.ok(mainActivity.includes('return insets;'),'No consumir insets y romper otros componentes');
assert.ok(mainActivity.includes('LAYOUT_IN_DISPLAY_CUTOUT_MODE_NEVER'),'Proteger recorte con APIs 28–29');
assert.ok(mainActivity.includes('Build.VERSION.SDK_INT >= 30'),'No utilizar getInsets Type en API antiguos');
assert.ok(!mainActivity.includes('setPadding(0, 35,'),'No introducir margen fijo por modelo');
assert.ok(android.html.includes('Android v9 Beta 13'));
assert.ok(gradle.includes("applicationId 'es.cine365.app'"),'Mantener package');
assert.ok(gradle.includes("versionName '9.0.0-beta.13'"),'Versión');
assert.ok(manifest.includes('Selector Cine 365 · Beta 12.2'));
console.log('OK: safe area dinámica para notch, navegación gestual y compatibilidad Android antiguos');


// Un único cálculo de las zonas seguras: el contenedor nativo ajusta el WebView.
// No añadir un segundo inset inferior en el HTML (problema visto en Redmi).
const html=android.html;
assert.ok(!html.includes('env(safe-area-inset-bottom'), 'El HTML NO debe duplicar el inset inferior ya aplicado por MainActivity');
assert.ok(html.includes('padding:6px 3px 7px!important'), 'Barra inferior con padding normal, sin notch duplicado');
assert.ok(html.includes('padding:12px 13px 76px!important'), 'Espacio bajo contenido calculado para el menú sin un segundo inset');
assert.ok(html.includes('.mobile-bottom-nav{padding-bottom:7px!important;bottom:0!important}'), 'La barra acaba en el borde seguro inferior de WebView');
assert.ok(mainActivity.includes('view.setPadding(bounds.left, bounds.top, bounds.right, bounds.bottom)'), 'Mantener ajuste nativo del recorte y navegación');
console.log('OK: Beta 12.2 elimina la doble aplicación de insets sin perder protección del notch');


// Beta 13: cuatro posiciones estables y persistencia de TMDB/listas borradas.
for(const page of [desktop,android]){
 assert.ok(page.script.includes('libraryState.deleted_lists[id]=true'),'Los borrados deben dejar tombstones');
 assert.ok(page.script.includes('libraryState.tmdb_overrides[key]=id'),'Las correcciones manuales TMDB se guardan en biblioteca');
 assert.ok(page.script.includes('reconcileTMDBLibrary()'),'Los cambios de TMDB recibidos se aplican localmente');
 assert.ok(page.script.includes("'/api/platforms/overrides/import'"),'Importación de identidades a motor local');
 assert.ok(page.script.includes("initPlatforms.retries"),'Reintentos de configuración tras reinicio');
}
assert.ok(desktop.script.includes('cine-list-placeholder'),'La página corta conserva posiciones vacías');
assert.ok(desktop.html.includes('grid-auto-rows:70px!important'),'Altura constante para las cuatro filas');
const src=fs.readFileSync('sync_desktop.go','utf8');
assert.ok(src.includes('result.DeletedLists[other.ID]'),'Windows no reintroduce listas eliminadas');
assert.ok(src.includes('result.TMDBOverrides[key]=id'),'Windows combina correcciones TMDB');
const mobile=fs.readFileSync('android/app/src/main/java/es/cine365/app/SyncStore.java','utf8');
assert.ok(mobile.includes('deleted.optBoolean(other.optString("id"))'),'Android no reintroduce listas borradas');
assert.ok(mobile.includes('result.put("tmdb_overrides",tmdb)'),'Android combina identificaciones TMDB');
console.log('OK: 4 filas, borrados en ambos sentidos, TMDB manual y persistencia en relanzamiento.');
