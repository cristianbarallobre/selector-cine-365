'use strict';
const assert=require('node:assert/strict');
const fs=require('node:fs');
const vm=require('node:vm');
const html=fs.readFileSync('lab_discover_windows.html','utf8');
const js=html.split('<!-- LAB_DISCOVERY_SCRIPT -->')[1];
const go=fs.readFileSync('lab_discover_windows.go','utf8');
new vm.Script('(()=>{\n'+js+'\n})();');
assert.match(html,/value="featured">Destacadas/);
assert.match(html,/value="rated">Mejor valoradas/);
assert.match(html,/value="popular">Más populares/);
assert.match(html,/value="recent">Recientes/);
assert.match(go,/VoteAverage float64/);
assert.match(go,/VoteCount int/);
assert.match(go,/CastOrder int/);
assert.ok(!go.includes('if len(films)>16{films=films[:16]}'),'No truncar filmografía antes de ordenar');

class Fake {
 constructor(tag){this.tagName=tag;this.children=[];this.events={};this.value='featured';this.textContent='';this.className='';}
 replaceChildren(...children){this.children=[...children];}
 appendChild(node){this.children.push(node);}
 addEventListener(event,fn){this.events[event]=fn;}
 click(){this.events.click?.();}
}
const elements={labPersonFilms:new Fake('div'),labPersonFilmSort:new Fake('select'),labPersonFilmsCount:new Fake('small')};
let filmID=0;
const ctx={
 $:id=>{assert.ok(elements[id],id);return elements[id];},
 labDetailCurrent:()=>({kind:'person'}),
 document:{createElement:tag=>new Fake(tag)},
 labShowDetail:id=>{filmID=id;},
 Math,Number,String,Array
};
vm.createContext(ctx);
const start=js.indexOf('let labPersonMovies=[];');
const end=js.indexOf('async function labOpenPerson(id,record=true){',start);
assert.ok(start>0&&end>start);
vm.runInContext(js.slice(start,end),ctx);
const movies=[
 {id:1,title:'Taquillazo secundario',vote_average:5,vote_count:500,popularity:400,order:45,release_date:'2024-01-01'},
 {id:2,title:'Papel protagonista',vote_average:8.5,vote_count:1000,popularity:55,order:0,release_date:'2010-01-01'},
 {id:3,title:'Peli votada una vez',vote_average:10,vote_count:1,popularity:5,order:1,release_date:'2015-01-01'},
 {id:4,title:'Muy reciente',vote_average:6.5,vote_count:100,popularity:4,order:2,release_date:'2026-12-01'}
];
const top=mode=>ctx.labSortedPersonMovies(movies,mode).map(x=>x.id);
assert.equal(top('featured')[0],2,'Priorizar buena crítica y papel relevante antes que cameo famoso');
assert.equal(top('popular')[0],1,'Popularidad TMDB original');
assert.equal(top('recent')[0],4,'Ordenar por estreno más reciente');
assert.equal(top('rated')[0],2,'Usar rating ajustado, no película con un solo voto');
assert.ok(ctx.labMovieCredibleRating(movies[2])<ctx.labMovieCredibleRating(movies[1]));
assert.equal(ctx.labSortedPersonMovies([{id:0,title:'X'},...movies, {id:3,title:'Repetida'}].filter(x=>x.id!==3),'recent').length,3);

const all=Array.from({length:35},(_,i)=>({
 id:100+i,title:'Obra '+(i+1),vote_average:7,vote_count:300,
 popularity:i,order:i,release_date:'2023-01-01',character:'Intérprete'
}));
vm.runInContext('labPersonMovies = '+JSON.stringify(all),ctx);
ctx.labRenderPersonFilmography();
assert.equal(elements.labPersonFilms.children.length,16,'Mostrar únicamente 16 tarjetas a la vez');
assert.match(elements.labPersonFilmsCount.textContent,/16 de 35/);
elements.labPersonFilmSort.value='popular';
elements.labPersonFilmSort.events.change();
assert.match(elements.labPersonFilms.children[0].children[0].textContent,/Obra 35/);
elements.labPersonFilms.children[0].click();
assert.equal(filmID,134,'Abrir la ficha de la película de la selección');
console.log('OK: 4 criterios TMDB independientes, destacado ponderado y filmografía de 16 sobre catálogo completo');
