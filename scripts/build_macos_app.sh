#!/bin/bash
set -euo pipefail
: "$MAC_ARCH"
: "$MAC_LABEL"
export CGO_ENABLED=1 GOOS=darwin GOARCH="$MAC_ARCH" MACOSX_DEPLOYMENT_TARGET=12.0
if [ "$MAC_ARCH" = "arm64" ] && [ "$(uname -m)" = "x86_64" ]; then
  # Compilación cruzada: el test se ejecuta en la CPU Intel anfitriona.
  GOARCH=amd64 CC=clang go test ./...
else
  go test ./...
fi
mkdir -p dist
APP='dist/Selector Cine 365.app'
mkdir -p "$APP/Contents/MacOS" "$APP/Contents/Resources"
go build -trimpath -ldflags='-s -w' -o "$APP/Contents/MacOS/SelectorCine365" .
file "$APP/Contents/MacOS/SelectorCine365"
cat > "$APP/Contents/Info.plist" <<'PLIST'
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0"><dict>
  <key>CFBundleExecutable</key><string>SelectorCine365</string>
  <key>CFBundleIdentifier</key><string>es.cristianbarallobre.selectorcine365</string>
  <key>CFBundleName</key><string>Selector Cine 365</string>
  <key>CFBundleDisplayName</key><string>Selector Cine 365</string>
  <key>CFBundleShortVersionString</key><string>8.2</string>
  <key>CFBundleVersion</key><string>820</string>
  <key>CFBundlePackageType</key><string>APPL</string>
  <key>CFBundleIconFile</key><string>AppIcon</string>
  <key>LSMinimumSystemVersion</key><string>12.0</string>
  <key>NSHighResolutionCapable</key><true/>
  <key>NSAppTransportSecurity</key><dict><key>NSAllowsLocalNetworking</key><true/></dict>
  <key>NSHumanReadableCopyright</key><string>© 2026 Cristian Barallobre Rúa. Todos los derechos reservados.</string>
</dict></plist>
PLIST
mkdir -p dist/icon.iconset
for size in 16 32 128 256 512; do
  sips -s format png -z "$size" "$size" icon.png --out "dist/icon.iconset/icon_"$size"x"$size".png" >/dev/null
  double=$((size*2))
  sips -s format png -z "$double" "$double" icon.png --out "dist/icon.iconset/icon_"$size"x"$size"@2x.png" >/dev/null
done
iconutil -c icns dist/icon.iconset -o "$APP/Contents/Resources/AppIcon.icns"
rm -rf dist/icon.iconset
plutil -lint "$APP/Contents/Info.plist"
codesign --force --deep --sign - "$APP"
codesign --verify --deep --strict --verbose=2 "$APP"
cat > dist/LEEME_MAC.txt <<'NOTE'
Selector Cine 365 v9.0 beta
Autor: Cristian Barallobre Rúa
https://github.com/cristianbarallobre/selector-cine-365

Mueve «Selector Cine 365.app» a Aplicaciones.
La app no está notarizada por Apple ni firmada con Developer ID.
Si macOS la bloquea, verifica la procedencia y usa Ajustes del Sistema
→ Privacidad y seguridad → Abrir igualmente.
Datos: ~/Library/Application Support/AtlasDelCine/
NOTE
(
  cd dist
  zip -qry "Selector_Cine_365_v9_beta_Mac_$MAC_LABEL.zip" 'Selector Cine 365.app' LEEME_MAC.txt
  unzip -t "Selector_Cine_365_v9_beta_Mac_$MAC_LABEL.zip" | tail -1
)
