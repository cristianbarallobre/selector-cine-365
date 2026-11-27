#!/usr/bin/env bash
set -euo pipefail
out="$1"
mkdir -p "$out/previews"
version=$(sed -n 's/^const appVersion = "\(.*\)"$/\1/p' version.go)
test -n "$version"
for name in MANUAL_COMPLETO_SELECTOR_CINE_365 GUIA_SINCRONIZACION_SEGURA; do
  pandoc "docs/$name.md" --from=gfm --to=docx --metadata lang=es-ES \
    --reference-doc="docs/$name.docx" -o "$out/$name.docx"
  unzip -tq "$out/$name.docx"
  libreoffice -env:UserInstallation="file://$out/lo-profile" --headless \
    --convert-to pdf --outdir "$out" "$out/$name.docx"
  test -s "$out/$name.pdf"
  pdfinfo "$out/$name.pdf"
  pdftotext -layout "$out/$name.pdf" "$out/$name.txt"
  grep -F "$version" "$out/$name.txt"
  pdftoppm -scale-to 1000 -png "$out/$name.pdf" "$out/previews/$name"
done
