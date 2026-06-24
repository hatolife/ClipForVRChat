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

VERSION="$(git -C "${ROOT_DIR}" tag -l 'v*' | sort -rV | head -n1)"
if [[ -z "${VERSION}" ]]; then
  VERSION="develop"
fi
REVISION="$(git -C "${ROOT_DIR}" rev-parse --short=7 HEAD)"
RELEASE_TIME="$(date -u +%Y-%m-%dT%H:%M:%SZ)"
APP_VERSION="${VERSION}.${REVISION}"
EXE_PATH="${SRC_DIR}/build/bin/${APP_NAME}.exe"
ZIP_PATH="${DIST_DIR}/${APP_NAME}-${APP_VERSION}-windows-amd64.zip"
WAILS_JSON="${SRC_DIR}/wails.json"
WINDOWS_INFO_JSON="${SRC_DIR}/build/windows/info.json"
WAILS_JSON_BACKUP="$(mktemp)"
WINDOWS_INFO_JSON_BACKUP="$(mktemp)"
WINDOWS_INFO_JSON_EXISTED=0

cp "${WAILS_JSON}" "${WAILS_JSON_BACKUP}"
if [[ -f "${WINDOWS_INFO_JSON}" ]]; then
  cp "${WINDOWS_INFO_JSON}" "${WINDOWS_INFO_JSON_BACKUP}"
  WINDOWS_INFO_JSON_EXISTED=1
fi
cleanup() {
  cp "${WAILS_JSON_BACKUP}" "${WAILS_JSON}"
  if [[ "${WINDOWS_INFO_JSON_EXISTED}" == "1" ]]; then
    cp "${WINDOWS_INFO_JSON_BACKUP}" "${WINDOWS_INFO_JSON}"
  else
    rm -f "${WINDOWS_INFO_JSON}"
  fi
  rm -f "${WAILS_JSON_BACKUP}" "${WINDOWS_INFO_JSON_BACKUP}"
}
trap cleanup EXIT

echo "Building ${APP_NAME} for windows/amd64"
echo "Version: ${VERSION}"
echo "Revision: ${REVISION}"
echo "Release time: ${RELEASE_TIME}"

node "${ROOT_DIR}/scripts/write-wails-info.mjs" "${VERSION}" "${WAILS_JSON}"
node "${ROOT_DIR}/scripts/write-wails-windows-info.mjs" "${VERSION}" "${REVISION}" "${WINDOWS_INFO_JSON}" >/dev/null

(
  cd "${SRC_DIR}"
  wails build -platform windows/amd64 -ldflags "-X main.version=${VERSION} -X main.revision=${REVISION} -X main.releaseTime=${RELEASE_TIME}"
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
