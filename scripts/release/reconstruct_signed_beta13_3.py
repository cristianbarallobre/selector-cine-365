#!/usr/bin/env python3
"""Reconstruye EXACTAMENTE el APK de Beta 13.3 ya firmada a partir de su artefacto CI.

El delta solo contiene bytes PUBLICOS del APK de distribución, no material
privado de firma. No se puede reutilizar para firmar otras versiones.
Si la compilación de origen difiere en un solo byte, aborta.
"""
from pathlib import Path
import base64, hashlib, struct, sys, zlib

SOURCE_SHA256 = "5ed201cd6a92278d1b9cb589bee43c79e124287849a3f7226234517a2a350774"
SIGNED_SHA256 = "e8db284ed38f33c2e91a0fcdb6b6cade684f18d8168cdd935ba4755f3ac70cf0"
MARKER = b"CINEDELTA1"

def digest(data):
    return hashlib.sha256(data).hexdigest()

def restore(unsigned_file, delta_file, output_file):
    src=Path(unsigned_file).read_bytes()
    if digest(src)!=SOURCE_SHA256:
        raise RuntimeError("APK sin firma inesperada: firma y build NO son compatibles; publicación cancelada")
    data=zlib.decompress(base64.b64decode(Path(delta_file).read_text().strip(),validate=True))
    if not data.startswith(MARKER):
        raise RuntimeError("Delta de distribución v9 Beta 13.3 inválido")
    pos=len(MARKER)
    out=bytearray()
    while pos<len(data):
        op=data[pos];pos+=1
        if op==0:
            off,size=struct.unpack_from(">II",data,pos);pos+=8
            if size==0 or off+size>len(src):
                raise RuntimeError("Copia fuera de rango")
            out+=src[off:off+size]
        elif op==1:
            size=struct.unpack_from(">I",data,pos)[0];pos+=4
            if size==0 or pos+size>len(data):
                raise RuntimeError("Literal fuera de rango")
            out+=data[pos:pos+size];pos+=size
        else:
            raise RuntimeError("Operación inválida")
    result=bytes(out)
    if digest(result)!=SIGNED_SHA256:
        raise RuntimeError("La huella SHA256 de APK firmada NO coincide. No publicar.")
    Path(output_file).parent.mkdir(parents=True,exist_ok=True)
    Path(output_file).write_bytes(result)
    print("APK firmada restaurada y comprobada, SHA-256:",SIGNED_SHA256)

if __name__=="__main__":
    if len(sys.argv)!=4:
        raise SystemExit("Uso: reconstruct_signed_beta13_3.py UNSIGNED.apk signed_delta.b64 OUTPUT.apk")
    restore(*sys.argv[1:])
