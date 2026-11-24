'use strict';
const assert=require('node:assert/strict');
const fs=require('node:fs');

const html=fs.readFileSync('plantilla.html','utf8');
const scripts=[...html.matchAll(/<script>([\s\S]*?)<\/script>/g)];
assert.ok(scripts.length,'Desktop HTML must have JavaScript');
for(const [,code] of scripts)assert.doesNotThrow(()=>new Function(code),'Desktop JavaScript must parse');
const js=scripts.at(-1)[1];

assert.ok(/\.library-platform-menu\{[^}]*width:100%;max-width:100%/.test(html),'Library dropdown fills the trigger width');
assert.ok(!html.includes('width:min(385px'),'Library dropdown must not use the old 385px cap');
for(const id of ['listSwitcher','sortPlatform']){
  assert.ok(html.includes('id="'+id+'Picker"'),'Custom picker exists: '+id);
  assert.ok(html.includes('id="'+id+'Button"'),'Custom button exists: '+id);
  assert.ok(html.includes('id="'+id+'Options"'),'Custom option menu exists: '+id);
  assert.ok(html.includes('id="'+id+'Menu" hidden'),'Custom menu starts closed: '+id);
}
assert.ok(html.includes('.single-select-button{')&&html.includes('cursor:pointer'),'Buttons must show a hand cursor');
assert.ok(html.includes('.single-select-option{width:100%;text-align:left;cursor:pointer'),'Dropdown options must show a hand cursor');
assert.ok(html.includes('id="platformSelect" class="single-select-native"'),'Keep platform select as source of truth');
assert.ok(html.includes('id="listSwitcher" class="single-select-native"'),'Keep list select as source of truth');
assert.ok(js.includes("if(old!==select.value)select.dispatchEvent(new Event('change'"),'Native change events remain connected');
assert.ok(js.includes("closeSingleChoiceMenu('listSwitcher')")&&js.includes("closeSingleChoiceMenu('sortPlatform')"),'Menus close on outside clicks / Escape');

class FakeElement {
  constructor() {
    this.hidden=true;this.children=[];this.dataset={};this.attrs={};this.listeners={};
    this.className='';this.textContent='';this.disabled=false;
    this.classList={
      toggle:(name,on)=>{this.attrs['class:'+name]=!!on;},
      remove:name=>{this.attrs['class:'+name]=false;}
    };
  }
  setAttribute(key,val){this.attrs[key]=String(val);}
  getAttribute(key){return this.attrs[key]??null;}
  addEventListener(type,fn){(this.listeners[type]??=[]).push(fn);}
  dispatchEvent(event){for(const fn of this.listeners[event.type]||[])fn(event);return true;}
  append(...children){this.children.push(...children);}
  appendChild(child){this.children.push(child);}
  replaceChildren(){this.children.length=0;}
  querySelectorAll(query){return query==='.single-select-option'?this.children.filter(x=>x.className.includes('single-select-option')):[];}
  getBoundingClientRect(){return {top:100,bottom:200,height:100};}
  focus(){document.activeElement=this;}
  click(){if(!this.disabled)this.dispatchEvent({type:'click'});}
}
class FakeSelect extends FakeElement {
  constructor(choices){
    super();this.options=choices.map(([value,textContent])=>({value,textContent}));this.value=this.options[0].value;
  }
  get value(){return this.selectedOptions[0]?.value||'';}
  set value(value){for(const o of this.options)o.selected=o.value===value;}
  get selectedOptions(){return this.options.filter(o=>o.selected);}
  get selectedIndex(){return this.options.findIndex(o=>o.value===this.value);}
}
const elements={};
for(const prefix of ['listSwitcher','sortPlatform']){
 for(const suffix of ['Button','ButtonLabel','Menu','Options','Picker'])elements[prefix+suffix]=new FakeElement();
}
for(const id of ['sortPlatformReset','listPlatformMenu','listPlatformButton'])elements[id]=new FakeElement();
elements.listSwitcher=new FakeSelect([['a','1. Mi primera lista'],['b','2. Mi segunda lista']]);
elements.platformSelect=new FakeSelect([['all','Todas las películas pendientes'],['any','Cualquier plataforma (streaming)'],['Netflix','Netflix'],['Filmin','Filmin']]);
elements.platformSelect.multiple=true;
const document={activeElement:null,createElement:()=>new FakeElement()};
const $=id=>{assert.ok(elements[id],'Unknown DOM id '+id);return elements[id];};
const start=js.indexOf('function closeSingleChoiceMenu('),end=js.indexOf('function eligibleFilms()',start);
assert.ok(start>=0&&end>start,'Custom picker helpers exist');
const controls=new Function('document','$','window','closeListPlatformMenu',
  js.slice(start,end)+'\nreturn {syncSingleChoiceMenu,closeSingleChoiceMenu};'
)(document,$,{innerHeight:850},()=>{});
let sortChanges=0,listChanges=0;
elements.platformSelect.addEventListener('change',()=>sortChanges++);
elements.listSwitcher.addEventListener('change',()=>listChanges++);

controls.syncSingleChoiceMenu('sortPlatform','platformSelect');
controls.syncSingleChoiceMenu('listSwitcher','listSwitcher');
assert.equal(elements.sortPlatformButtonLabel.textContent,'Todas las películas pendientes');
assert.equal(elements.sortPlatformOptions.children.length,4);
assert.equal(elements.sortPlatformOptions.children[0].getAttribute('aria-selected'),'true');
assert.equal(elements.sortPlatformOptions.children[1].getAttribute('aria-selected'),'false');

elements.sortPlatformButton.click();
assert.equal(elements.sortPlatformMenu.hidden,false,'Sort dropdown opens');
assert.equal(elements.sortPlatformButton.getAttribute('aria-expanded'),'true');
elements.sortPlatformOptions.children[2].click();
assert.equal(elements.platformSelect.value,'Netflix','Choose one streaming platform');
assert.equal(sortChanges,1,'One change event');
assert.equal(elements.sortPlatformMenu.hidden,false,'Permanece abierto para añadir otra plataforma');
assert.equal(elements.sortPlatformButtonLabel.textContent,'Netflix','Trigger shows current choice');

elements.sortPlatformOptions.children[3].click();
assert.deepEqual(elements.platformSelect.selectedOptions.map(o=>o.value),['Netflix','Filmin']);
assert.equal(elements.sortPlatformButtonLabel.textContent,'Netflix, Filmin');
elements.sortPlatformOptions.children[2].click();
assert.deepEqual(elements.platformSelect.selectedOptions.map(o=>o.value),['Filmin']);
elements.sortPlatformOptions.children[1].click();
assert.deepEqual(elements.platformSelect.selectedOptions.map(o=>o.value),['any']);
elements.sortPlatformReset.click();
assert.equal(elements.platformSelect.value,'all','Reset returns to all pending movies');
assert.equal(sortChanges,5,'Reset triggers change event');
assert.equal(elements.sortPlatformButtonLabel.textContent,'Todas las películas pendientes');

elements.listSwitcherButton.click();
assert.equal(elements.listSwitcherMenu.hidden,false,'List dropdown opens');
elements.listSwitcherOptions.children[1].click();
assert.equal(elements.listSwitcher.value,'b','Chosen collection changes');
assert.equal(listChanges,1,'Collection change event preserved');
assert.equal(elements.listSwitcherButtonLabel.textContent,'2. Mi segunda lista');
assert.equal(elements.listSwitcherMenu.hidden,true);

elements.listSwitcherButton.click();
elements.listSwitcherOptions.dispatchEvent({type:'keydown',key:'Escape',preventDefault(){},stopPropagation(){}});
assert.equal(elements.listSwitcherMenu.hidden,true,'Escape closes dropdown');
assert.equal(document.activeElement,elements.listSwitcherButton,'Escape returns keyboard focus');

console.log('OK: Beta 14 desktop custom dropdown CSS, accessibility, sort/list selection, reset and keyboard Escape.');
