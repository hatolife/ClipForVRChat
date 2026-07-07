# Dependency and Supply Chain

## Supply Chain

- Go moduleは `src/go.mod` / `src/go.sum` で管理されている。
- frontendは `src/frontend/package-lock.json` でlockされている。
- Spout2は `tools/spout-capture/CMakeLists.txt` でcommit revisionとarchive SHA256を固定している。
- AvatarBeacon unitypackageはRelease workflowでpackageと署名をダウンロードし、固定fingerprintのOpenPGP署名で検証している。

## CI/CD

- CIは `contents: read`。
- Release workflowはbuild/sign/package/releaseを分離し、Release作成jobだけ `contents: write`。
- 署名secretはsign jobだけで参照される。
- artifact digestをGitHub APIで検証してからsign/packageへ進む。

## GitHub Actions

- 問題あり: 外部Actionは `actions/checkout@v7` などmajor tag参照で、commit SHA pinningではない。
- 改善済み: 2026-06-25監査時に指摘されていたworkflow全体 `contents: write` は、現行Release workflowではjob単位へ縮小されている。

## Release Workflow

- Release tag形式は `vX.Y.Z` / `vX.Y.Z-aN` / `vX.Y.Z-bN` / `vX.Y.Z-rcN` に制限されている。
- GitHub Release本文は `RELEASE_NOTES.md` から抽出される。
- 公開Assetは通常zip、exe asc、separated zip、AvatarBeacon unitypackageに制限される。
- sha256やbuild metadataはartifactには含めるが公開Assetからは除外される設計。

## Dependency Pinning

- Go: module version固定あり。
- npm: lock fileあり。
- Wails/go/node/govulncheck versionはworkflow envで固定。
- Spout2: revisionとSHA256固定あり。
- GitHub Actions: major tag参照のため完全固定ではない。

## Lock File

- `src/go.sum` あり。
- `src/frontend/package-lock.json` あり。

## Reproducible Build

完全再現可能ビルドではない。Release時刻、Git revision、Windows runner、Wails/Go/Node toolchain、GitHub Actions環境に依存する。ただしversion metadataとbuild metadataを生成して追跡可能性を高めている。

## PGP署名

Release workflowでexe detached signatureを作成する。通常zipには `.exe.asc` が同梱される。AvatarBeacon packageも固定fingerprintで署名検証される。

## SHA256

Release workflowは内部検証用にexe/zip/helper/DLL/packageのSHA256を算出する。ただし公開Assetとしてsha256ファイルは添付しない運用。

## SBOM

SBOM生成は確認できない。Info扱いで将来導入を推奨する。

## 依存関係の脆弱性

ローカルでは外部DB問い合わせを避けたため、`npm audit` / `govulncheck` は未実行。CI/Release workflowには組み込まれている。`src/go.mod` では過去finding対象の `github.com/cloudflare/circl` が `v1.6.4` に更新済みであることを確認した。

## 配布物の改ざん対策

- PGP detached signature。
- Artifact digest検証。
- Release asset allowlist。
- AvatarBeacon package署名検証。
- Spout2 archive SHA256固定。

## ビルド環境の再現性

- Go/Node/Wails/govulncheck versionはworkflow envで固定。
- GitHub-hosted runnerの詳細やWindows SDK/MSVCはrunnerに依存する。
- Spout helperのWindows実ビルドはローカル未確認。
