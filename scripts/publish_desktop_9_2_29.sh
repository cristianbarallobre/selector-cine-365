#!/usr/bin/env bash
set -euo pipefail
version=$(sed -n 's/^const appVersion = "\(.*\)"$/\1/p' version.go)
test "$version" = "9.2.29"
tag="v$version-desktop"
if gh release view "$tag" --repo "$GITHUB_REPOSITORY" >/dev/null 2>&1; then
  echo "Ya existe la publicación $tag; no se sustituirá automáticamente."
  exit 1
fi
root="$ARTIFACT_ROOT"
out="$RUNNER_TEMP/cine-release"
mkdir -p "$out"
win="Selector_Cine_365_Windows_$version.exe"
intel="Selector_Cine_365_Mac_Intel_$version.zip"
arm="Selector_Cine_365_Mac_Apple_Silicon_$version.zip"
cp "$root/Windows/$win" "$out/"
cp "$root/Mac_Intel/$intel" "$out/"
cp "$root/Mac_Apple_Silicon/$arm" "$out/"
for file in "$win" "$intel" "$arm"; do test -s "$out/$file"; done
unzip -tq "$out/$intel"
unzip -tq "$out/$arm"
(cd "$out"; sha256sum "$win" "$intel" "$arm" > SHA256SUMS.txt; sha256sum -c SHA256SUMS.txt)
for name in MANUAL_COMPLETO_SELECTOR_CINE_365 GUIA_SINCRONIZACION_SEGURA; do
  for ext in pdf docx; do test -s "$root/Manuals/$name.$ext"; done
done
# Only merge the exact tested head. Never force-push or overwrite a release.
if [ "$GITHUB_REF" = "refs/heads/integration/desktop-9.2.29-main" ]; then
  info=$(gh pr view "$PR_NUMBER" --repo "$GITHUB_REPOSITORY" --json headRefOid,headRefName,baseRefName,state,isDraft)
  jq -e --arg sha "$GITHUB_SHA" '.headRefOid==$sha and .headRefName=="integration/desktop-9.2.29-main" and .baseRefName=="main" and .state=="OPEN" and .isDraft==false' <<< "$info"
  merge_sha=$(gh api --method PUT "repos/$GITHUB_REPOSITORY/pulls/$PR_NUMBER/merge" \
    -f sha="$GITHUB_SHA" -f merge_method=merge -f commit_title="Integrar escritorio 9.2.29 validado" --jq .sha)
else
  test "$GITHUB_REF" = "refs/heads/main"
  merge_sha="$GITHUB_SHA"
fi
git fetch --no-tags origin main
test "$(git rev-parse origin/main)" = "$merge_sha"
git diff --exit-code "$GITHUB_SHA" "$merge_sha" --
git switch --detach "$merge_sha"
for name in MANUAL_COMPLETO_SELECTOR_CINE_365 GUIA_SINCRONIZACION_SEGURA; do
  for ext in pdf docx; do cp "$root/Manuals/$name.$ext" "docs/$name.$ext"; done
done
git config user.name "github-actions[bot]"
git config user.email "41898282+github-actions[bot]@users.noreply.github.com"
git add docs/MANUAL_COMPLETO_SELECTOR_CINE_365.{pdf,docx} docs/GUIA_SINCRONIZACION_SEGURA.{pdf,docx}
if ! git diff --cached --quiet; then git commit -m "docs: generar manuales de escritorio 9.2.29 en PDF y Word"; fi
git push origin HEAD:main
target=$(git rev-parse HEAD)
# Tag the actual merged source plus the generated manuals; only docs differ from
# the tested source. Upload to a draft before making the release public.
gh release create "$tag" "$out/$win" "$out/$intel" "$out/$arm" "$out/SHA256SUMS.txt" \
  --repo "$GITHUB_REPOSITORY" --target "$target" --title "Selector Cine 365 $version · escritorio" \
  --notes-file RELEASE_NOTES_v9_2_29_DESKTOP.md --draft
gh release view "$tag" --repo "$GITHUB_REPOSITORY" --json assets --jq '.assets[].name' | sort > "$out/actual-assets.txt"
printf '%s\n' "$win" "$intel" "$arm" SHA256SUMS.txt | sort > "$out/expected-assets.txt"
diff -u "$out/expected-assets.txt" "$out/actual-assets.txt"
gh release edit "$tag" --repo "$GITHUB_REPOSITORY" --draft=false --latest
gh release view "$tag" --repo "$GITHUB_REPOSITORY" --json url,isDraft,tagName
