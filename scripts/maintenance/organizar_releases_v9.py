#!/usr/bin/env python3
"""Consolidar las releases v9 sin perder binarios, tags ni código fuente."""
import os, subprocess, json, hashlib, tempfile, zipfile
from pathlib import Path

REPO="cristianbarallobre/selector-cine-365"
BRANCH="maintenance/organizar-releases-v9-2026-11"
OLD=[
"v9.0.0-windows-beta.9","v9.0.0-windows-beta.8","v9.0.0-windows-beta.10",
"v9.0.0-android-beta.9","v9.0.0-android-beta.8","v9.0.0-android-beta.7",
"v9.0.0-android-beta.6","v9.0.0-android-beta.10",
"manual-azure-entra-onedrive-2026-11",
"v9.0.0-windows-beta.5.1","v9.0.0-windows-beta.5","v9.0.0-windows-beta.4",
"v9.0.0-android-beta.5","v9.0.0-android-beta.4","v9.0.0-android-beta.3",
"v9.0.0-android-beta.2","v9.0.0-android-beta.1"
]
STABLE={"v9.0.0-beta.12","v8.2","v8.1.1"}
FINAL=STABLE|{"v9.0.0-beta.9"}
ARCHIVE="archivo-historico-beta-v9-2026-11"
BETA9_COMMIT="764b72bac3580b89a04c5f753fbb61e223fa5966"
BETA12_FILES={
"Selector_Cine_365_Android_v9_BETA12.apk",
"Selector_Cine_365_Mac_Apple_Silicon_v9_BETA12.zip",
"Selector_Cine_365_Mac_Intel_v9_BETA12.zip",
"Selector_Cine_365_Windows_v9_BETA12.zip"
}
MANUALS=[
"docs/GUIA_AZURE_ENTRA_ONEDRIVE_SELECTOR_CINE_365.pdf",
"docs/MANUAL_COMPLETO_SELECTOR_CINE_365.pdf"
]
def cmd(*args):
    p=subprocess.run(args,stdout=subprocess.PIPE,stderr=subprocess.PIPE,text=True)
    if p.returncode:raise RuntimeError("Falló "+repr(args)+": "+p.stderr[-2000:])
    return p.stdout.strip()
def api(endpoint):
    return json.loads(cmd("gh","api","-X","GET","repos/"+REPO+"/"+endpoint))
def rel(tag):return api("releases/tags/"+tag)
def listing():return api("releases?per_page=100")
def names(r):return {x["name"] for x in r["assets"]}
def hashfile(p):return hashlib.sha256(Path(p).read_bytes()).hexdigest()
def tags():
    return {tag:api("git/ref/tags/"+tag)["object"]["sha"] for tag in OLD}
def checktags(saved):
    now=tags()
    assert now==saved,"Alguno de los tags históricos ha desaparecido o cambiado"
def download(tag,name,dest):
    dest=Path(dest);dest.mkdir(parents=True,exist_ok=True)
    cmd("gh","release","download",tag,"-R",REPO,"--pattern",name,"--dir",str(dest))
    p=dest/name
    assert p.is_file(),f"No se descargó {tag}/{name}"
    return p
def check_public(expected):
    pub=[x["tag_name"] for x in listing() if not x["draft"]]
    assert len(pub)==len(expected) and set(pub)==expected,str(pub)

assert os.getenv("GITHUB_REPOSITORY")==REPO
assert os.getenv("GITHUB_REF_NAME")==BRANCH
assert os.getenv("GITHUB_ACTIONS")=="true" and os.getenv("GH_TOKEN")
old=set(OLD)
assert len(old)==17
published={x["tag_name"] for x in listing() if not x["draft"]}
if published==FINAL:
    assert BETA12_FILES <= names(rel("v9.0.0-beta.12"))
    print("Limpieza previamente terminada: 4 releases públicas")
    raise SystemExit(0)
assert published==STABLE|old,"Releases modificadas desde la última comprobación: "+repr(published)
assert BETA12_FILES <= names(rel("v9.0.0-beta.12"))
assert len(rel("v8.2")["assets"])==3 and len(rel("v8.1.1")["assets"])==1
saved_tags=tags()
work=Path(tempfile.mkdtemp(prefix="cine_releases_"))
manifest=[]
archive=work/"Archivo_Instaladores_Betas_Anteriores_v9.zip"
with zipfile.ZipFile(archive,"w",compression=zipfile.ZIP_DEFLATED,compresslevel=6) as z:
    for tag in OLD:
        obj=rel(tag)
        assert obj["assets"] and not obj["draft"]
        for a in obj["assets"]:
            path=download(tag,a["name"],work/"originales"/tag)
            data=path.read_bytes()
            assert len(data)==a["size"],"Tamaño incorrecto "+tag
            record={"tag_original":tag,"commit_original":obj["target_commitish"],
                    "asset":a["name"],"bytes":len(data),"sha256":hashfile(path)}
            manifest.append(record)
            z.writestr("originales/"+tag+"/"+a["name"],data)
            print("RESPALDADO",tag,a["name"])
    z.writestr("MANIFIESTO.json",json.dumps(manifest,indent=2,ensure_ascii=False))
    z.writestr("LEEME.txt",
      "Archivo histórico interno de Selector Cine 365. Contiene los ZIP y APK ORIGINALES "
      "de las 17 releases antiguas, con SHA256 en MANIFIESTO.json. "
      "Los tags Git y fuentes originales se conservan. "
      "No borrar la release DRAFT sin realizar otra copia independiente.")
with zipfile.ZipFile(archive) as z:
    assert z.testzip() is None
    for row in manifest:
        assert hashlib.sha256(z.read("originales/"+row["tag_original"]+"/"+row["asset"])).hexdigest()==row["sha256"]
assert len(manifest)>=17 and archive.stat().st_size>500000
print("Archivo histórico listo:",len(manifest),"archivos; bytes:",archive.stat().st_size)
notes=work/"ARCHIVO_NOTAS.txt"
notes.write_text("BORRADOR privado que conserva las APK y ZIP originales de betas antiguas; "
                 "no es una release para instalar. No eliminar sin copia adicional.",encoding="utf-8")
cmd("gh","release","create",ARCHIVE,str(archive),"-R",REPO,"--draft",
    "--target","main","--title","Archivo privado de betas históricas v9","--notes-file",str(notes))
obj=next((x for x in listing() if x["tag_name"]==ARCHIVE),None)
assert obj and obj["draft"] and any(
  x["name"]==archive.name and x["size"]==archive.stat().st_size for x in obj["assets"])
copy=download(ARCHIVE,archive.name,work/"comprobado")
assert hashfile(copy)==hashfile(archive),"El archivo de GitHub no coincide con el original"
print("VERIFICADA la copia remota de GitHub, byte por byte")

windows=work/"originales"/"v9.0.0-windows-beta.9"/"Selector_Cine_365_Windows_v9_BETA9.zip"
android=work/"originales"/"v9.0.0-android-beta.9"/"Selector_Cine_365_Android_v9_BETA9.apk"
assert windows.exists() and android.exists()
notes9=work/"NOTAS_BETA9.txt"
notes9.write_text(
    "HISTÓRICO - Beta 9 para Windows y Android. Ambas descargas son idénticas "
    "a las anteriores. Si estás instalando por primera vez utiliza Beta 12. "
    "Android Beta 9 utiliza otro ID y otra firma; para migrar a Beta 12 es "
    "necesario exportar JSON, conservarlo fuera de la aplicación e importarlo.",
    encoding="utf-8")
cmd("gh","release","create","v9.0.0-beta.9",str(windows),str(android),
    "-R",REPO,"--target",BETA9_COMMIT,"--prerelease",
    "--title","Selector Cine 365 v9 Beta 9 (histórica: Windows y Android)",
    "--notes-file",str(notes9))
assert names(rel("v9.0.0-beta.9"))=={windows.name,android.name}
for f in (windows,android):
    confirmed=download("v9.0.0-beta.9",f.name,work/"comprobado"/"beta9")
    assert hashfile(confirmed)==hashfile(f)
print("Beta 9 consolidada y verificada")

for p in MANUALS:
    f=Path(p)
    assert f.is_file() and f.stat().st_size>10000 and f.read_bytes().startswith(b"%PDF-")
cmd("gh","release","upload","v9.0.0-beta.12","-R",REPO,*MANUALS)
assert BETA12_FILES|{Path(x).name for x in MANUALS}<=names(rel("v9.0.0-beta.12"))
for p in MANUALS:
    local=Path(p)
    confirmed=download("v9.0.0-beta.12",local.name,work/"comprobado"/"manuales")
    assert hashfile(confirmed)==hashfile(local)
print("Manuales actualizados de Beta12 publicados y verificados")
checktags(saved_tags)
assert next((x for x in listing() if x["tag_name"]==ARCHIVE),{}).get("draft")
print("COPIA VERIFICADA, BETA9 CONSOLIDADA Y MANUALES PUBLICADOS; retirando releases redundantes")
for tag in OLD:
    cmd("gh","release","delete",tag,"-R",REPO,"--yes")
    print("RETIRADA:",tag)
check_public(FINAL)
checktags(saved_tags)
assert BETA12_FILES<=names(rel("v9.0.0-beta.12"))
assert next((x for x in listing() if x["tag_name"]==ARCHIVE),{}).get("draft")
print("LISTO: cuatro releases visibles; borrador con archivo intacto; 17 tags conservados")
