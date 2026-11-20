#!/usr/bin/env python3
"""Reduce el espaciado excesivo del historial v9 en Word.

Pandoc hereda de la referencia v8.1.1 18 puntos entre párrafos de 'Body Text',
lo que produce una página casi vacía al final del historial v9 Beta 13.1.
Este script no cambia contenido ni fuente; solo espaciado. Compatible con futuras versiones.
"""
from pathlib import Path
import sys
from docx import Document
from docx.shared import Pt

if len(sys.argv)!=2:
    raise SystemExit('Uso: compactar_historial_v9.py historial.docx')
path=Path(sys.argv[1])
if not path.is_file():
    raise SystemExit('El historial Word no existe')
doc=Document(path)
for name in ('Body Text',):
    try:
        style=doc.styles[name]
        style.paragraph_format.space_before=Pt(4)
        style.paragraph_format.space_after=Pt(4)
    except KeyError:
        pass
for name in ('Heading 1','Heading 2','Heading 3'):
    try:
        style=doc.styles[name]
        style.paragraph_format.space_before=Pt(9 if name=='Heading 1' else 5)
        style.paragraph_format.space_after=Pt(3)
    except KeyError:
        pass
for par in doc.paragraphs:
    if par.style.name in ('Body Text','Compact'):
        par.paragraph_format.space_before=Pt(2)
        par.paragraph_format.space_after=Pt(3)
doc.save(path)
print('Espaciado y saltos de páginas normalizados:',path)
