#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SRC_DIR="${ROOT_DIR}/src"
DIST_DIR="${ROOT_DIR}/dist/local-windows"
APP_NAME="ClipForVRChat"
ZIP=0
VERSION_PROVIDED=0

while [[ $# -gt 0 ]]; do
  case "$1" in
    --zip)
      ZIP=1
      shift
      ;;
    --version)
      if [[ $# -lt 2 ]]; then
        echo "--version requires a value." >&2
        exit 2
      fi
      VERSION="$2"
      VERSION_PROVIDED=1
      shift 2
      ;;
    *)
      echo "usage: $0 [--version vX.Y.Z] [--zip]" >&2
      exit 2
      ;;
  esac
done

if ! command -v git >/dev/null 2>&1; then
  echo "git is required." >&2
  exit 1
fi

if ! command -v wails >/dev/null 2>&1; then
  echo "wails is required. Install Wails CLI in WSL before running this script." >&2
  exit 1
fi

if [[ "${VERSION_PROVIDED}" == "0" ]]; then
  VERSION="${VERSION:-v0.1.7}"
fi
REVISION="$(git -C "${ROOT_DIR}" rev-parse --short=7 HEAD)"
RELEASE_TIME="$(date -u +%Y-%m-%dT%H:%M:%SZ)"
BUILD_CHANNEL="release"
if [[ "${VERSION_PROVIDED}" == "0" ]]; then
  BUILD_CHANNEL="develop"
fi
APP_VERSION="${VERSION}.${REVISION}"
DISPLAY_REVISION="${REVISION}"
if [[ "${BUILD_CHANNEL}" == "develop" ]]; then
  APP_VERSION="${APP_VERSION}.develop"
  DISPLAY_REVISION="${REVISION}.develop"
fi
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
echo "Build channel: ${BUILD_CHANNEL}"
echo "Release time: ${RELEASE_TIME}"

node "${ROOT_DIR}/scripts/write-wails-info.mjs" "${VERSION}" "${WAILS_JSON}"
node "${ROOT_DIR}/scripts/write-wails-windows-info.mjs" "${VERSION}" "${DISPLAY_REVISION}" "${WINDOWS_INFO_JSON}" >/dev/null

(
  cd "${SRC_DIR}"
  wails build -platform windows/amd64 -ldflags "-X main.version=${VERSION} -X main.revision=${REVISION} -X main.releaseTime=${RELEASE_TIME} -X main.buildChannel=${BUILD_CHANNEL}"
)

if [[ ! -f "${EXE_PATH}" ]]; then
  echo "Build failed: ${EXE_PATH} was not created." >&2
  exit 1
fi

if [[ "${ZIP}" == "1" ]]; then
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
