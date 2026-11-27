// LAB-02.8: instalada antes de las primeras peticiones, solo Windows LAB.
// Se reintentan únicamente GET locales. Nunca repetir POST/PUT/PATCH/DELETE:
// una respuesta perdida no demuestra que la escritura no se haya aplicado.
const labOriginalFetch=window.fetch.bind(window);
const labRetryHTTP=new Set([408,502,503,504]);
const labLastErrorLog=new Map();
function labNetworkInfo(input,options){
 const method=String(options?.method||input?.method||'GET').toUpperCase();
 const url=typeof input==='string'?input:input?.url;
 try{
  const parsed=new URL(url,window.location.href);
  const local=parsed.origin===window.location.origin&&parsed.pathname.startsWith('/api/');
  return {method,path:local?parsed.pathname:'',local};
 }catch(e){return {method,path:'',local:false};}
}
function labNetworkSource(path){
 if(path.startsWith('/api/discover/')||path.startsWith('/api/platforms/')||path.startsWith('/api/identities'))return 'TMDB';
 if(path.startsWith('/api/extra-backup/'))return 'Copia adicional';
 return 'la aplicación';
}
function labNetworkReport(info,kind,status,attempts){
 // No se envían URL completas, querystrings, títulos, contraseñas ni cuerpos.
 if(!info.local||!info.path||info.path==='/api/ping'||info.path==='/api/lab/network-log')return;
 if(!/^[a-z0-9/_-]+$/.test(info.path.slice(1))||info.path.length>78)return;
 const key=info.method+'|'+info.path+'|'+kind+'|'+status;
 const now=Date.now();
 if(now-(labLastErrorLog.get(key)||0)<5000)return;
 labLastErrorLog.set(key,now);
 if(labLastErrorLog.size>75)labLastErrorLog.clear();
 const body=JSON.stringify({route:info.path,method:info.method,kind,status,attempts});
 void labOriginalFetch('/api/lab/network-log',{
  method:'POST',headers:{'Content-Type':'application/json','X-Atlas-Token':token},
  body,cache:'no-store'
 }).catch(()=>{});
}
const labNetworkPause=ms=>new Promise(resolve=>setTimeout(resolve,ms));
async function labResilientFetch(input,options){
 const info=labNetworkInfo(input,options);
 // No alterar contenidos externos, peticiones de control, ni el registro mismo.
 if(!info.local||info.path==='/api/ping'||info.path==='/api/lab/network-log')return labOriginalFetch(input,options);
 const safeRead=info.method==='GET'||info.method==='HEAD';
 const signal=options?.signal||input?.signal;
 const tries=safeRead?3:1;
 for(let attempt=1;attempt<=tries;attempt++){
  try{
   const response=await labOriginalFetch(input,options);
   if(safeRead&&labRetryHTTP.has(response.status)&&attempt<tries&&!signal?.aborted){
    // No dejar cuerpos abiertos antes de reintentar una consulta.
    try{await response.body?.cancel()}catch(e){}
    await labNetworkPause(attempt===1?300:900);
    continue;
   }
   if(response.status===408||response.status===429||response.status>=500){
    labNetworkReport(info,'http',response.status,attempt);
   }
   return response;
  }catch(err){
   if(signal?.aborted||err?.name==='AbortError')throw err;
   const networkError=err?.name==='TypeError'||err?.name==='NetworkError';
   if(safeRead&&networkError&&attempt<tries){
    await labNetworkPause(attempt===1?300:900);
    continue;
   }
   if(networkError){
    labNetworkReport(info,'network',0,attempt);
    const service=labNetworkSource(info.path);
    if(safeRead)throw new Error('No se pudo consultar '+service+' a través del servidor local después de '+attempt+' intentos. Comprueba la conexión y vuelve a probar.');
    throw new Error('Conexión interrumpida al guardar o cambiar datos de '+service+'. No se ha repetido la operación: comprueba si se aplicó antes de intentarlo de nuevo.');
   }
   throw err;
  }
 }
}
window.fetch=labResilientFetch;
