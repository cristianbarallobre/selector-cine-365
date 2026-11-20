#!/usr/bin/env python3
"""Genera docs/MANUAL_COMPLETO_SELECTOR_CINE_365.pdf desde Markdown."""
from pathlib import Path
import mistune
from weasyprint import HTML, CSS

root = Path(__file__).resolve().parents[1]
source = root / "docs" / "MANUAL_COMPLETO_SELECTOR_CINE_365.md"
target = root / "docs" / "MANUAL_COMPLETO_SELECTOR_CINE_365.pdf"

md = mistune.create_markdown(plugins=["table", "url", "strikethrough"])
body = md(source.read_text(encoding="utf-8"))
css = r"""
@page {
  size:A4;
  margin:19mm 18mm 18mm 18mm;
  @bottom-center {
    content:"Selector Cine 365 · Manual completo | Página " counter(page) " de " counter(pages);
    font-family:"DejaVu Sans",Arial,sans-serif;font-size:8pt;color:#6c7b91;
  }
}
*{box-sizing:border-box}
html{font-family:"DejaVu Sans",Arial,sans-serif;color:#142435;font-size:9.4pt}
body{line-height:1.55}
h1{font-size:29pt;margin:0 0 9mm;padding:0 0 5mm;border-bottom:4px solid #1578bc;color:#0b213b;line-height:1.1}
h2{font-size:14pt;color:#0a568a;margin:9mm 0 3mm;page-break-after:avoid;line-height:1.25}
h3{font-size:11pt;color:#28465d;margin:5mm 0 2mm;page-break-after:avoid}
p{margin:0 0 3mm;text-align:left}
strong{font-weight:700;color:#132e4a}
a{color:#0a68a8;text-decoration:none;overflow-wrap:anywhere}
blockquote{background:#eff5fa;padding:10px 12px;margin:7mm 0;border-left:4px solid #1977b7;color:#34536c}
blockquote p{margin:0}
ul,ol{margin:1mm 0 4mm;padding-left:6mm}
li{margin:0 0 1.5mm}
code{font-family:"DejaVu Sans Mono",monospace;background:#edf2f6;border-radius:3px;padding:1px 2px;overflow-wrap:anywhere;word-break:break-word;font-size:8.8pt}
pre{white-space:pre-wrap;background:#f0f5f8;border:1px solid #d5e2ed;padding:11px 13px;border-radius:7px;line-height:1.7;page-break-inside:avoid}
pre code{padding:0;background:transparent}
table{width:100%;border-collapse:collapse;margin:4mm 0 5mm;table-layout:fixed;font-size:8pt}
thead{display:table-header-group}
tr{page-break-inside:avoid}
th{background:#153b5f;color:white;text-align:left}
th,td{border:1px solid #d7e2ec;padding:7px 8px;vertical-align:top;overflow-wrap:anywhere}
tbody tr:nth-child(even) td{background:#f3f7fb}
hr{border:none;border-top:1px solid #b7c9d9;margin:8mm 0}
"""
html = f'<!doctype html><html lang="es"><head><meta charset="utf-8"></head><body>{body}</body></html>'
HTML(string=html).write_pdf(target, stylesheets=[CSS(string=css)])
print(f"Generado {target} ({target.stat().st_size} bytes)")
