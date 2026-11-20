#!/usr/bin/env python3
"""Reconstruye EXACTAMENTE el APK de Beta 13.2 ya firmada a partir de su artefacto CI.

El delta solo contiene bytes PUBLICOS del APK de distribución, no material
privado de firma. No se puede reutilizar para firmar otras versiones.
Si la compilación de origen difiere en un solo byte, aborta.
"""
from pathlib import Path
import base64, hashlib, struct, sys, zlib

SOURCE_SHA256 = "ecaacec9a1a9364945c0c0d3209198c179deebb633c9456247cd5bfea8d71e92"
SIGNED_SHA256 = "16698e3dcb5b12ea74eaff749b48af2045c8b90b50b9128c32bb61d9d6d76853"
MARKER = b"CINEDELTA1"

def digest(data):
    return hashlib.sha256(data).hexdigest()

def restore(unsigned_file, delta_file, output_file):
    src=Path(unsigned_file).read_bytes()
    if digest(src)!=SOURCE_SHA256:
        raise RuntimeError("APK sin firma inesperada: firma y build NO son compatibles; publicación cancelada")
    data=zlib.decompress(base64.b64decode(Path(delta_file).read_text().strip(),validate=True))
    if not data.startswith(MARKER):
        raise RuntimeError("Delta de distribución v9 Beta 13.2 inválido")
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
        raise SystemExit("Uso: reconstruct_signed_beta13_2.py UNSIGNED.apk signed_delta.b64 OUTPUT.apk")
    restore(*sys.argv[1:])
