#!/usr/bin/env python3
"""Selector Cine 365: genera PDF desde Markdown, con ReportLab y DejaVu.
Uso: python generar_guia_azure_entra_onedrive.py ENTRADA.md SALIDA.pdf
"""
import html, re, sys
from pathlib import Path
from reportlab.lib import colors
from reportlab.lib.pagesizes import A4
from reportlab.lib.styles import ParagraphStyle
from reportlab.pdfbase import pdfmetrics
from reportlab.pdfbase.ttfonts import TTFont
from reportlab.platypus import HRFlowable, Paragraph, SimpleDocTemplate, Spacer, Table, TableStyle

NAVY=colors.HexColor('#17293E'); SLATE=colors.HexColor('#55697F')
GOLD=colors.HexColor('#D3A04E'); PALE=colors.HexColor('#EFF4F9')
BORDER=colors.HexColor('#D6E2EA'); INK=colors.HexColor('#243241')


def load_fonts():
    folder=Path('/usr/share/fonts/truetype/dejavu')
    if not folder.exists(): raise RuntimeError('Install fonts-dejavu-core')
    for name,file in [('CineSans','DejaVuSans.ttf'),('CineBold','DejaVuSans-Bold.ttf'),('CineMono','DejaVuSansMono.ttf')]:
        pdfmetrics.registerFont(TTFont(name,str(folder/file)))
    pdfmetrics.registerFontFamily('CineSans',normal='CineSans',bold='CineBold')


def style(name,size,leading,color=INK,**kw):
    kw.setdefault('fontName','CineSans')
    return ParagraphStyle(name,fontSize=size,leading=leading,textColor=color,**kw)


def inline(s):
    pattern=r'(\[[^\]]+\]\(https?://[^)]+\)|\*\*[^*]+\*\*|`[^`]+`|https?://[^\s]+)'
    tokens=re.split(pattern,s)
    result=[]
    for t in tokens:
        if t.startswith('**') and t.endswith('**'):
            result.append('<b>'+html.escape(t[2:-2])+'</b>')
        elif t.startswith('`') and t.endswith('`'):
            result.append('<font name="CineMono" size="7.8" color="#18527A">'+html.escape(t[1:-1])+'</font>')
        elif t.startswith('['):
            label,url=t[1:].split('](',1);url=url[:-1]
            result.append('<link href="'+html.escape(url,quote=True)+'" color="#225E91">'+html.escape(label)+'</link>')
        elif t.startswith('https://') or t.startswith('http://'):
            url=t.rstrip('.,;)')
            result.append('<link href="'+html.escape(url,quote=True)+'" color="#225E91">'+html.escape(url)+'</link>'+html.escape(t[len(url):]))
        else: result.append(html.escape(t))
    return ''.join(result).strip()


def decorate(canvas,doc):
    w,h=A4;canvas.saveState()
    canvas.setFillColor(NAVY);canvas.rect(0,h-10,w,10,fill=1,stroke=0)
    canvas.setFillColor(GOLD);canvas.rect(0,h-13,w,3,fill=1,stroke=0)
    if doc.page>1:
        canvas.setFont('CineBold',7.6);canvas.setFillColor(SLATE)
        canvas.drawString(43,h-36,'SELECTOR CINE 365 / MICROSOFT ENTRA Y ONEDRIVE')
        canvas.setStrokeColor(BORDER);canvas.line(43,h-44,w-43,h-44)
    canvas.setStrokeColor(BORDER);canvas.line(43,42,w-43,42)
    canvas.setFont('CineSans',7.1);canvas.setFillColor(SLATE)
    canvas.drawString(43,30,'Guía de configuración · Noviembre de 2026')
    canvas.drawRightString(w-43,30,'Página '+str(doc.page))
    canvas.restoreState()


def box(lines,fmt,width,code=False):
    ps=[Paragraph(html.escape(x) if code else inline(x),fmt) for x in lines]
    tbl=Table([[ps]],colWidths=[width],hAlign='LEFT')
    rules=[('BACKGROUND',(0,0),(-1,-1),PALE),('LEFTPADDING',(0,0),(-1,-1),14),
           ('RIGHTPADDING',(0,0),(-1,-1),12),('TOPPADDING',(0,0),(-1,-1),10),
           ('BOTTOMPADDING',(0,0),(-1,-1),10)]
    rules.append(('BOX',(0,0),(-1,-1),.6,BORDER) if code else ('LINEBEFORE',(0,0),(0,-1),3.5,GOLD))
    tbl.setStyle(TableStyle(rules));return tbl

def build(source,target):
    load_fonts(); W=A4[0]-86
    st={
        'eyebrow':style('eyebrow',8.4,13,GOLD,fontName='CineBold',spaceAfter=10),
        'title':style('title',24,31,NAVY,fontName='CineBold',spaceAfter=9),
        'subtitle':style('subtitle',13.4,20,SLATE,spaceAfter=13),
        'head':style('head',12.1,18,NAVY,fontName='CineBold',spaceBefore=18,spaceAfter=7,keepWithNext=True),
        'subhead':style('subhead',9.8,15,colors.HexColor('#245B85'),fontName='CineBold',spaceBefore=11,spaceAfter=5,keepWithNext=True),
        'body':style('body',8.6,13.35,INK,spaceAfter=8,splitLongWords=1),
        'bullet':style('bullet',8.55,13.1,INK,leftIndent=15,firstLineIndent=-12,spaceAfter=5,splitLongWords=1),
        'num':style('num',8.55,13.1,INK,leftIndent=20,firstLineIndent=-17,spaceAfter=5,splitLongWords=1),
        'quote':style('quote',8.5,13.4,NAVY),
        'code':style('code',8,13.0,NAVY,fontName='CineMono',wordWrap='CJK'),
    }
    doc=SimpleDocTemplate(str(target),pagesize=A4,leftMargin=43,rightMargin=43,
        topMargin=60,bottomMargin=57,
        title='Selector Cine 365 - Guía Microsoft Azure Entra y OneDrive',
        author='Documentación de Selector Cine 365')
    flow=[];buf=[];quotes=[];codes=[];in_code=False;is_cover=True
    def flush():
        if buf:
            flow.append(Paragraph(inline(' '.join(buf)),st['body']));buf.clear()
    def flush_quotes():
        if quotes:
            flow.append(box([' '.join(quotes)],st['quote'],W));flow.append(Spacer(1,11));quotes.clear()
    for line in Path(source).read_text(encoding='utf-8').splitlines():
        line=line.rstrip()
        if line.startswith('```'):
            flush();flush_quotes()
            if in_code:
                flow.append(box(codes,st['code'],W,code=True));flow.append(Spacer(1,10));codes=[]
            in_code=not in_code;continue
        if in_code:codes.append(line);continue
        if line.startswith('> '):flush();quotes.append(line[2:]);continue
        flush_quotes()
        if not line:flush();continue
        if line=='---':
            flush();flow.append(HRFlowable(width='100%',color=BORDER,thickness=.7,spaceBefore=9,spaceAfter=12));continue
        if line.startswith('# '):
            flush();flow.append(Spacer(1,19));flow.append(Paragraph('GUÍA / DOCUMENTO DE REFERENCIA',st['eyebrow']))
            flow.append(Paragraph(inline(line[2:]),st['title']));continue
        if line.startswith('## '):
            flush()
            if is_cover:
                flow.append(Paragraph(inline(line[3:]),st['subtitle']))
                flow.append(HRFlowable(width='100%',color=GOLD,thickness=2,spaceBefore=0,spaceAfter=16))
                is_cover=False
            else:flow.append(Paragraph(inline(line[3:]),st['head']))
            continue
        if line.startswith('### '):flush();flow.append(Paragraph(inline(line[4:]),st['subhead']));continue
        if line.startswith('- '):flush();flow.append(Paragraph('• '+inline(line[2:]),st['bullet']));continue
        found=re.match(r'^(\d+)\.\s+(.*)',line)
        if found:
            flush();flow.append(Paragraph('<b>'+found.group(1)+'.</b> '+inline(found.group(2)),st['num']));continue
        buf.append(line.strip())
    flush();flush_quotes()
    if codes:flow.append(box(codes,st['code'],W,code=True))
    doc.build(flow,onFirstPage=decorate,onLaterPages=decorate)
    print('Built PDF:',target,Path(target).stat().st_size,'bytes')

if __name__=='__main__':
    if len(sys.argv)!=3:raise SystemExit('Usage: generate.py input.md output.pdf')
    build(sys.argv[1],sys.argv[2])