#!/usr/bin/env bash

set -euo pipefail

OUT_DIR="${1:-dist}"

ROOT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")/.." && pwd)"
DIST_DIR="${ROOT_DIR}/${OUT_DIR}"
CACHE_DIR="${ROOT_DIR}/.cache/go-build"

mkdir -p "${CACHE_DIR}"
export GOCACHE="${CACHE_DIR}"

rm -rf "${DIST_DIR}"
mkdir -p "${DIST_DIR}/assets"

echo "Building WebAssembly binary in ${DIST_DIR}"
GOOS=js GOARCH=wasm go build -o "${DIST_DIR}/main.wasm" "${ROOT_DIR}"

WASM_EXEC="$(go env GOROOT)/misc/wasm/wasm_exec.js"
if [[ ! -f "${WASM_EXEC}" ]]; then
  ALT_WASM_EXEC="$(go env GOROOT)/lib/wasm/wasm_exec.js"
  if [[ -f "${ALT_WASM_EXEC}" ]]; then
    WASM_EXEC="${ALT_WASM_EXEC}"
  else
    echo "wasm_exec.js not found in GOROOT" >&2
    exit 1
  fi
fi

cp "${WASM_EXEC}" "${DIST_DIR}/wasm_exec.js"
cp "${ROOT_DIR}/web/index.html" "${DIST_DIR}/index.html"
cp -R "${ROOT_DIR}/assets/static" "${DIST_DIR}/assets/"

echo "Artifacts written to ${DIST_DIR}"
