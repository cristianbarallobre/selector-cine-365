'use strict';
// Regression: on a synced phone the library exists but posters/provider details
// are not included in OneDrive snapshots. Complete visible TMDB metadata without
// scanning the entire catalog or making any OneDrive calls.
const assert=require('node:assert/strict');
const fs=require('node:fs'),vm=require('node:vm');
const html=fs.readFileSync('android/app/src/main/assets/index.html','utf8');
const start=html.indexOf('const autoTmdbTried=new Set();');
const end=html.indexOf('function showPlatformModal(){',start);
assert.ok(start>=0&&end>start,'Visible-card on-demand downloader exists');
const body=html.slice(start,end);
const requested=[];
const results=new Map();
const parsed={entries:Array.from({length:100},(_,i)=>({
 title:i===0?'Parasite':'Film '+i,year:i===0?'2019':'2000',idx:i
}))};
const progress={textContent:''};
let renderCalls=0;
const ctx={
 window:{CineAndroid:{}},document:{hidden:false},platformConfigured:true,
 platformActive:false,filmPage:0,parsed,platformData:results,
 platformApi:async(path,details)=>{
  requested.push({path,details});
  return {
   state:'ok',metadata_version:3,poster_path:'/poster-'+details.title,
   tmdb_id:details.title==='Parasite'?496243:123,subscriptions:['Filmin']
  };
 },
 providerKey:f=>f.title.toLowerCase()+'|'+Number(f.year||0),
 registerIdentity:()=>{},
 scheduledPlatformRender:()=>{renderCalls++},
 updatePlatformProgress:()=>{progress.textContent='Disponibilidad de plataformas en España'},
 $:()=>progress,setTimeout,console,Promise,Set,Map,Math,Number
};
vm.createContext(ctx);
vm.runInContext(body+'\nthis.autoFill=fillVisibleTmdbGaps;',ctx);
(async()=>{
 await ctx.autoFill();
 assert.equal(requested.length,12,'At most 12 visible film details per visit');
 assert.ok(requested.every(x=>x.path==='/api/platforms/lookup' && x.details.force===false),'No cloud request or forced bypass of retry cooldown');
 assert.equal(requested[0].details.title,'Parasite');
 assert.equal(results.get('parasite|2019').poster_path,'/poster-Parasite');
 await ctx.autoFill();
 assert.equal(requested.length,24,'Next progressive batch fetches more missing visible entries');
 await ctx.autoFill();
 assert.equal(requested.length,24,'Only first 24 cards fetched on current page');
 ctx.filmPage=1;
 await ctx.autoFill();
 await ctx.autoFill();
 assert.equal(requested.length,48,'Max 48 per app session');
 await ctx.autoFill();
 assert.equal(requested.length,48,'No unbounded requests');
 assert.ok(renderCalls>0,'Cards rerender when data becomes available');
 ctx.platformConfigured=false;
 await ctx.autoFill();
 assert.equal(requested.length,48,'No requests without configured TMDB key');
 assert.match(html,/if\(window\.CineAndroid&&platformConfigured\)void fillVisibleTmdbGaps\(\)/,'Start on app opening only after TMDB local status');
 assert.match(html,/await initLibrary\(\);restartList\(\);\s*if\(window\.CineAndroid\)void fillVisibleTmdbGaps\(\)/,'Start after importing OneDrive list');
 assert.match(html,/filmPage\+\+;renderList\(!!window\.CineAndroid\);if\(window\.CineAndroid\)void fillVisibleTmdbGaps\(\)/,'Load newly visible page');
 assert.ok(!body.includes('/api/cloud/'),'TMDB completion does not sync to OneDrive');
 console.log('OK: TMDB poster lookup for visible Android movies, bounded 12/page and 48/session, key-gated, no OneDrive requests.');
})().catch(e=>{console.error(e);process.exitCode=1});
