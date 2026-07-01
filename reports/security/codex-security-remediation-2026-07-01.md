# Codex Security findings remediation 2026-07-01

## 概要

Codex Securityが過去commitに対して作成したfindingを、現在HEAD基準で再検証する。
古いcommitへ戻すのではなく、現在の `release/v0.1.8` の実装を基準に、未修正または部分修正の問題だけを修正する。

## 入力

- ユーザー指定CSV `reports/security/codex-security-findings-2026-06-29.csv` は現在の作業ツリーに存在しなかった。
- ユーザー指定ディレクトリ内の `reports/security/2026-07-01T04-48-55.763Z/codex-security-findings-2026-07-01T04-48-55.763Z.csv` を一次情報として使用した。
- 外部finding URLは参照していない。

## 作業開始時点

- ブランチ: `release/v0.1.8`
- HEAD: `e47a169594f087fd4988474be26e5408740fb476`
- remoteとの差分: `origin/release/v0.1.8...HEAD` は `ahead 2, behind 2`
- 未コミット変更: 作業開始時点の `git status --short --branch` では未コミット変更なし
- 分岐メモ: local側に `e47a169 docs(reports): add codex security reports`、remote側に `3fb9fe8 add codex security reports` があり、同種のreport追加コミットで分岐している。既存変更は破棄しない。

## 分類

- auto-photo / screenshot auto-post / imported config consent
  - Imported configs can start silent auto-photo exfiltration
- settings UI tab visibility / hidden sensitive config
  - Tabbed settings can hide auto-post exfiltration config
- GitHub Actions tag/ref name shell injection
  - Tag-name shell injection exposes release signing secrets
  - GitHub Actions tag name command injection
  - Release workflow injects untrusted tag names into PowerShell
- release workflow cached executable trust
  - Release build trusts cached Wails executable
- unpinned tools/actions supply-chain risk
  - Unpinned govulncheck runs in release workflow
  - Release workflow uses unpinned write-capable actions
- PGP verification trust anchor
  - PGP verification trusts a release-bundled public key

## Findingごとの現在HEAD判定

| Title | Severity | Original commit | Relevant paths | 現在HEAD状態 | 判定理由 | 修正方針 |
| --- | --- | --- | --- | --- | --- | --- |
| Tabbed settings can hide auto-post exfiltration config | high | `dc532f209103` | `src/frontend/src/main.js`, `src/app.go`, `src/internal/appcore/autophoto.go` | fixed | デフォルトタブは現在 `feature` で自動処理ON/OFFは見える。今回、保存前に自動処理の監視フォルダとDiscord投稿先を表示する確認モーダルを追加した。 | 対応済み。 |
| Release build trusts cached Wails executable | high | `22ca489c3d9d` | `.github/workflows/release.yml`, `.github/workflows/ci.yml` | fixed | Release/CI workflowから `actions/cache` による `~/go/bin/wails.exe` 復元を削除し、固定バージョン `v2.12.0` を毎回 `go install` する。 | 対応済み。 |
| Tag-name shell injection exposes release signing secrets | high | `0fb2dfe1e5f6` | `.github/workflows/release.yml` | fixed | `Prepare version` をCheckout直後に移動し、`vX.Y.Z` / `vX.Y.Z-rcN` だけを許可する。Release jobはbuild job outputの検証済みtagだけを使う。 | 対応済み。 |
| PGP verification trusts a release-bundled public key | high | `749bd2efbab9` | `.github/workflows/release.yml`, `README.md`, `src/frontend/src/main.js`, `issues/038-about-official-distribution-and-pgp.md` | fixed | README/About/SPECでRelease同梱URLや公開鍵だけでは真正性の根拠にならないことを明記し、`release-signing@hato.life` のfingerprintを照合する案内へ変更した。 | 対応済み。 |
| GitHub Actions tag name command injection | high | `2a3fa25208bf` | `.github/workflows/release.yml` | fixed | tag検証を厳格化し、artifact名、release title、upload pathは検証済みtagから生成する。 | 対応済み。 |
| Unpinned govulncheck runs in release workflow | high | `5c935b18630e` | `.github/workflows/release.yml`, `.github/workflows/ci.yml` | already fixed | 現在は `GOVULNCHECK_VERSION: v1.1.4` を使い、`@latest` は残っていない。 | 追加修正なし。再検索で `@latest` がないことを確認する。 |
| Imported configs can start silent auto-photo exfiltration | high | `7efdf8d5720a` | `src/main.go`, `src/app.go`, `src/internal/appcore/autophoto.go`, `src/internal/appcore/config.go` | fixed | `App.startup` は `ModeResults` の時だけwatcherを起動するため、起動引数jsonを開いただけではwatcherは開始しない。今回、保存時に自動処理内容の確認を必須にした。 | 対応済み。 |
| Release workflow injects untrusted tag names into PowerShell | high | `62630835b5b6` | `.github/workflows/release.yml` | fixed | PowerShell runブロックへ未検証tagを直接埋め込まず、env経由で受けた値をCheckout直後に厳格検証してからjob outputとして使う。 | 対応済み。 |
| Release workflow uses unpinned write-capable actions | high | `9a816fec98b8` | `.github/workflows/release.yml` | already fixed / residual low | 現在は `softprops/action-gh-release` を使わず、write権限jobでは `gh` CLIでrelease作成/アップロードしている。`actions/download-artifact@v7` はwrite権限jobで使われるがGitHub公式action。 | third-party release actionは復活させない。write権限jobの操作をtag検証と`gh`に限定する。 |

## 実施した修正

- `.github/workflows/release.yml`
  - `workflow_dispatch` に `tag_name` 入力を追加した。
  - Checkout直後に `Prepare version` を置き、`^v\d+\.\d+\.\d+(-rc\d+)?$` 以外を拒否するようにした。
  - release jobは `github.ref_name` ではなくbuild job outputの検証済み `tag_name`、`draft`、`prerelease` を使うようにした。
  - `actions/cache` で復元した `wails.exe` を信頼する経路を削除し、Wails CLIは固定バージョンを毎回installするようにした。
- `.github/workflows/ci.yml`
  - CI側の `wails.exe` cacheも削除し、固定バージョンinstallに統一した。
- `src/frontend/src/main.js`, `src/frontend/src/style.css`
  - autoPhoto / screenshotAutoPost が有効な設定を保存する前に、監視フォルダとDiscord投稿先を表示する確認モーダルを追加した。
  - Webhook URLはIDと短縮tokenだけを表示し、生tokenを画面に出さない。
- `README.md`, `src/frontend/src/main.js`, `src/SPEC.md`, `RELEASE_NOTES.md`, `issues/038-about-official-distribution-and-pgp.md`
  - `release-signing@hato.life` のfingerprint `BE40 AA8D 082F 493F 613B C072 21DC 3486 1B40 E77D` を明記した。
  - Release同梱URLや公開鍵だけでは公開鍵自体の真正性を確認できないことを明記した。
  - `Good signature` の過剰表現を避け、信頼済みfingerprintの公開鍵で検証した場合に限る表現へ変更した。

## 未対応または意図的に対応しない項目

- `reports/security/codex-security-findings-2026-06-29.csv` は存在しないため、指定ディレクトリ内CSVを使用した。
- `softprops/action-gh-release` は現在HEADに存在しないため、commit SHA pinではなく「使わない」方針を維持した。

## 実行したテスト

- `npm run build` (`src/frontend`): 成功
- `go test ./...` (`src`): 成功
- `node scripts/check-wails-api-surface.mjs`: 成功
- `python3` + PyYAMLで `.github/workflows/ci.yml` / `.github/workflows/release.yml` 読み込み: 成功
- `npm audit --omit=dev` (`src/frontend`): 成功、0 vulnerabilities
- `go run golang.org/x/vuln/cmd/govulncheck@v1.1.4 ./...` (`src`): 成功、called vulnerability 0件
- `git diff --check`: 成功
- `git show --name-status --oneline <commit_hash> -- <relevant_paths>`: 9 findingすべてについて実行し、過去commitの変更対象を確認した
- 自己レビュー検索: `github.ref_name`, `@latest`, `actions/cache`, `wails.exe`, `softprops`, `action-gh-release`, `Good signature`, `Release-signing-public-key.asc` を `.github`、README、frontend、SPEC、Release Notes、issues、報告書で検索した。workflow上の `github.ref_name` はenv入力1箇所のみで、runブロック直埋めではない。`@latest`、`actions/cache`、`wails.exe`、`softprops` はworkflowには残っていない。

## 残リスク

- Release workflowのWindows/Release成果物生成はローカルLinuxでは完全実行できないため、最終確認はGitHub Actionsで行う。
- `workflow_dispatch` は `tag_name` 入力必須にしたため、従来の「ブランチから手動実行してdraft Releaseを作る」用途は使えない。今回のsecurity要件に合わせ、保守的にRC/正式tagだけを許可する。
- 自動処理の保存前確認はフロントエンドUIの確認であり、backend APIを直接呼ぶローカル攻撃者には設定保存自体を止めない。対象はユーザー操作時の同意/可視化であり、既存config互換性を維持するためbackend schemaは変更していない。

## レビュー時に見るべき差分

- `.github/workflows/release.yml`
- `.github/workflows/ci.yml`
- `src/frontend/src/main.js`
- `README.md`
- `src/SPEC.md`
- `issues/038-about-official-distribution-and-pgp.md`
