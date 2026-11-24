#!/bin/bash
set -euo pipefail
work=$(mktemp -d)
unzip -q "$1" -d "$work/package"
bundle="$work/package/Selector Cine 365.app"
codesign --verify --deep --strict "$bundle"
swiftc scripts/test_macos_webkit.swift -o "$work/probar-webkit"
# Start the packaged binary with an isolated profile and verify its own WKWebView heartbeat.
mkdir "$work/home"
HOME="$work/home" "$bundle/Contents/MacOS/SelectorCine365" >"$work/app.log" 2>&1 &
pid=$!
trap 'kill "$pid" 2>/dev/null || true; wait "$pid" 2>/dev/null || true' EXIT
log="$work/home/Library/Application Support/AtlasDelCine/diagnostico_macos.txt"
for i in {1..90}; do
  if grep -q 'WKWebView: interfaz iniciada' "$log" 2>/dev/null; then break; fi
  kill -0 "$pid" || { cat "$work/app.log"; exit 1; }
  sleep 1
done
grep -q 'WKWebView: interfaz iniciada' "$log" || { cat "$work/app.log" "$log"; exit 1; }
url=$(sed -n 's/^.*Servidor local activo en \(http:\/\/127\.0\.0\.1:[0-9]*\/\)$/\1/p' "$log" | tail -1)
test -n "$url"
"$work/probar-webkit" "$url" scripts/test_desktop_picker.js
