#!/usr/bin/env python3
"""Exporta las fuentes EXCLUSIVAMENTE desde el commit de la release Beta 13.1.
La documentación Markdown se actualiza desde la rama principal para incluir
las correcciones de manuales posteriores a la etiqueta sin tocar el código.
Nunca incluye claves privadas de firma ni archivos de usuario.
"""
from pathlib import Path
import subprocess, hashlib, json, zipfile, re

RELEASE_TAG="v9.0.0-beta.13.1"
COMMIT="dc56af0c6243da2baf425c1914dc6885aae38238"
SOURCE=Path("source")
LATEST=Path("latest-docs")
OUTPUT=Path("dist")
OUTPUT.mkdir(exist_ok=True)

def git(folder,*args):
    return subprocess.check_output(["git","-C",str(folder),*args],stderr=subprocess.STDOUT).decode().strip()

head=git(SOURCE,"rev-parse","HEAD")
assert head==COMMIT,f"No se empaquetará otro commit: {head}"
docs_head=git(LATEST,"rev-parse","HEAD")
tracked=[p.decode() for p in subprocess.check_output(["git","-C",str(SOURCE),"ls-files","-z"]).split(b"\0") if p]
assert len(tracked)>=100, f"checkout incompleto: {len(tracked)} archivos"
assert len(tracked)==len(set(tracked))
assert "versionCode 9015" in (SOURCE/"android/app/build.gradle").read_text()
assert "versionName '9.0.0-beta.13'" in (SOURCE/"android/app/build.gradle").read_text()
assert "applicationId 'es.cine365.app'" in (SOURCE/"android/app/build.gradle").read_text()
assert "Beta 13.1" in (SOURCE/"plantilla.html").read_text()
assert "9.0.0-beta.13.1" in (SOURCE/"scripts/build_macos_app.sh").read_text()
assert 'case "/api/platforms/status", "/api/platforms/cache"' in (SOURCE/"main.go").read_text()
assert (SOURCE/"android/app/src/main/java/es/cine365/app/MainActivity.java").exists()

SENSITIVE=re.compile(r"(^|/)(\.env($|\.)|[^/]*\.(jks|keystore|p12|pfx|pem|key)$|[^/]*(?:private_key|secret_key)[^/]*)", re.I)
assert not any(SENSITIVE.search(p) for p in tracked),"Hay un archivo sospechoso en el repositorio"
def doc_or_notes(p):
    return (p.startswith("docs/") and p.endswith(".md")) or p in {
        "README.md", "CHANGELOG.md", "MACOS.md",
        "android/LEEME_ANDROID_v9_BETA.md",
        "RELEASE_NOTES_v9_BETA13_1.md"
    }
latest_docs={p:(LATEST/p).read_bytes() for p in tracked if doc_or_notes(p) and (LATEST/p).is_file()}
assert len(latest_docs)>7, "Falta documentación principal"
root_sources={
    p for p in tracked if "/" not in p and (
        p.endswith((".go",".mod",".sum",".m",".syso",".ico",".png",".html",".txt",".md"))
        or p==".gitignore"
    )
}
common_docs={p for p in tracked if (p.startswith("docs/") and p.endswith(".md"))
             or p in ("README.md","CHANGELOG.md","AUTORIA.md","MACOS.md",
                       "RELEASE_NOTES_v9_BETA13_1.md")}
common_scripts={p for p in tracked if p.startswith("scripts/") and
                (p.endswith(".cjs") or p.endswith(".py") and not p.startswith("scripts/release/"))}
android_files={p for p in tracked if p.startswith("android/")}
assert len(android_files)>15 and len(root_sources)>25

targets={
  "Windows": {
    "label":"Windows — v9 Beta 13.1",
    "files":root_sources | common_docs | common_scripts,
    "build":("Go 1.23+, Windows 10/11 x64 y WebView2.\n"
             "go test ./...\n"
             "go build -trimpath -ldflags=\"-H windowsgui -s -w\" -o Selector_Cine_365.exe .\n"),
  },
  "Android": {
    "label":"Android — v9 Beta 13 (APK publicada en la release 13.1)",
    "files":android_files | common_docs | common_scripts |
            {p for p in tracked if p.startswith("scripts/release/") and "beta13" in p},
    "build":("Java 17, Android SDK 35, Gradle 8.7.\n"
             "cd android && gradle testDebugUnitTest assembleRelease\n"
             "AVISO: assembleRelease sin la clave privada crea APK sin firmar. "
             "Para instalar encima de Beta 12.x se necesita EXACTAMENTE el "
             "mismo certificado, que NO se distribuye en estos archivos.\n"),
  },
  "Mac_Apple_Silicon": {
    "label":"macOS Apple Silicon — v9 Beta 13.1 (ARM64)",
    "files":root_sources | common_docs | common_scripts | {"scripts/build_macos_app.sh"},
    "build":("Compilar en macOS 12+ con Go 1.23+ y Xcode Command Line Tools.\n"
             "export MAC_ARCH=arm64 MAC_LABEL=Apple_Silicon\n"
             "bash scripts/build_macos_app.sh\n"
             "La firma del script es ad hoc, no Apple Developer ID ni notarización.\n"),
  },
  "Mac_Intel": {
    "label":"macOS Intel — v9 Beta 13.1 (x86_64)",
    "files":root_sources | common_docs | common_scripts | {"scripts/build_macos_app.sh"},
    "build":("Compilar en macOS 12+ con Go 1.23+ y Xcode Command Line Tools.\n"
             "export MAC_ARCH=amd64 MAC_LABEL=Intel\n"
             "bash scripts/build_macos_app.sh\n"
             "La firma del script es ad hoc, no Apple Developer ID ni notarización.\n"),
  },
  "Proyecto_Completo": {
    "label":"Proyecto COMPLETO — Windows 13.1, Android 13, Mac ARM/Intel 13.1",
    "files":set(tracked),
    "build":("Todos los archivos versionados de la release (incluidas "
             "pruebas, configuraciones, documentación, CI y recursos). "
             "No contiene la clave privada Android.\n"),
  },
}

def archive_name(target):
    if target=="Android":
        return "Selector_Cine_365_Android_v9_BETA13_CODIGO_FUENTE.zip"
    if target=="Windows":
        return "Selector_Cine_365_Windows_v9_BETA13_1_CODIGO_FUENTE.zip"
    if target=="Mac_Intel":
        return "Selector_Cine_365_Mac_Intel_v9_BETA13_1_CODIGO_FUENTE.zip"
    if target=="Mac_Apple_Silicon":
        return "Selector_Cine_365_Mac_Apple_Silicon_v9_BETA13_1_CODIGO_FUENTE.zip"
    return "Selector_Cine_365_Proyecto_Completo_v9_BETA13_1_CODIGO_FUENTE.zip"

summary=[]
for target,details in targets.items():
    chosen=sorted(details["files"])
    assert len(chosen)>10 and set(chosen)<=set(tracked)
    prefix="Selector_Cine_365_"+target+"_Fuentes_v9_BETA13_1"
    notes=(f"SELECTOR CINE 365 — CÓDIGO FUENTE\n{details['label']}\n"
           f"Release oficial: https://github.com/cristianbarallobre/selector-cine-365/releases/tag/{RELEASE_TAG}\n"
           f"FUENTES de código exactas de commit Git: {COMMIT}\n"
           f"Documentación Markdown actualizada desde main: {docs_head}\n"
           "Comparten el mismo frontend y backend Go Windows y Mac; "
           "los dos Mac son compilaciones de distinta arquitectura, no proyectos independientes.\n"
           "Archivos de código originales y recursos versionados; no incluye APK/EXE/\n"
           "credenciales privadas de firma ni datos del usuario.\n\n"
           "INSTRUCCIONES DE COMPILACIÓN:\n"+details["build"]+
           "\nVer MANIFIESTO_SHA256.json para las huellas de cada archivo y su procedencia.\n")
    manifest={"tag":RELEASE_TAG,"source_commit":COMMIT,"latest_docs_commit":docs_head,
              "platform":target,"files":[]}
    out=OUTPUT/archive_name(target)
    with zipfile.ZipFile(out,"w",compression=zipfile.ZIP_DEFLATED,compresslevel=8) as z:
        z.writestr(prefix+"/LEEME_CODIGO_FUENTE.txt",notes)
        for p in chosen:
            docs_updated=p in latest_docs
            data=latest_docs[p] if docs_updated else (SOURCE/p).read_bytes()
            if SENSITIVE.search(p):raise RuntimeError("Archivo privado: "+p)
            z.writestr(prefix+"/"+p,data)
            manifest["files"].append({
                "path":p, "sha256":hashlib.sha256(data).hexdigest(),
                "origin":("main_documentation" if docs_updated else "release_tag")
            })
        z.writestr(prefix+"/MANIFIESTO_SHA256.json",
                   json.dumps(manifest,ensure_ascii=False,indent=2)+"\n")
    with zipfile.ZipFile(out) as z:
        assert z.testzip() is None
        assert len(z.namelist())==len(chosen)+2
    summary.append({"platform":target,"name":out.name,"files":len(chosen),
                    "bytes":out.stat().st_size,
                    "sha256":hashlib.sha256(out.read_bytes()).hexdigest()})
    print(target,":",len(chosen),"archivos,",out.stat().st_size,"bytes")
(OUTPUT/"RESUMEN_ZIPS_FUENTES.json").write_text(json.dumps({
    "tag":RELEASE_TAG,"source_commit":COMMIT,"latest_docs_commit":docs_head,
    "archives":summary},indent=2,ensure_ascii=False)+"\n")
assert len(summary)==5
print("OK: CINCO ZIP COMPROBADOS Y VERIFICADOS")
