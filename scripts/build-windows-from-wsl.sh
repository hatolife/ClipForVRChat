#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SRC_DIR="${ROOT_DIR}/src"
DIST_DIR="${ROOT_DIR}/dist/local-windows"
APP_NAME="ClipForVRChat"

if ! command -v git >/dev/null 2>&1; then
  echo "git is required." >&2
  exit 1
fi

if ! command -v wails >/dev/null 2>&1; then
  echo "wails is required. Install Wails CLI in WSL before running this script." >&2
  exit 1
fi

VERSION="$(git -C "${ROOT_DIR}" rev-parse --short=12 HEAD)"
EXE_PATH="${SRC_DIR}/build/bin/${APP_NAME}.exe"
ZIP_PATH="${DIST_DIR}/${APP_NAME}-${VERSION}-windows-amd64.zip"

echo "Building ${APP_NAME} for windows/amd64"
echo "Version: ${VERSION}"

(
  cd "${SRC_DIR}"
  wails build -platform windows/amd64 -ldflags "-X main.version=${VERSION}"
)

if [[ ! -f "${EXE_PATH}" ]]; then
  echo "Build failed: ${EXE_PATH} was not created." >&2
  exit 1
fi

if [[ "${1:-}" == "--zip" ]]; then
  rm -rf "${DIST_DIR}/${APP_NAME}"
  mkdir -p "${DIST_DIR}/${APP_NAME}"
  cp "${EXE_PATH}" "${DIST_DIR}/${APP_NAME}/${APP_NAME}.exe"
  (
    cd "${DIST_DIR}"
    rm -f "${ZIP_PATH}" "${ZIP_PATH}.sha256"
    zip -qr "$(basename "${ZIP_PATH}")" "${APP_NAME}"
    sha256sum "$(basename "${ZIP_PATH}")" >"$(basename "${ZIP_PATH}").sha256"
  )
  echo "Created ${ZIP_PATH}"
  echo "Created ${ZIP_PATH}.sha256"
else
  sha256sum "${EXE_PATH}"
fi
