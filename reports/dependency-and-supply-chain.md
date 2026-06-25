# Dependency And Supply Chain

調査日: 2026-06-25

## 直接依存

### Go

- `github.com/ProtonMail/go-crypto v1.4.1`
- `github.com/alexflint/go-arg v1.6.1`
- `github.com/disintegration/imaging v1.6.2`
- `github.com/gofrs/flock v0.13.0`
- `github.com/makiuchi-d/gozxing v0.1.1`
- `github.com/skip2/go-qrcode v0.0.0-20200617195104-da1b6568686e`
- `github.com/wailsapp/wails/v2 v2.12.0`
- `golang.design/x/clipboard v0.8.0`
- `golang.org/x/image v0.43.0`

### npm

- `vue 3.5.38`
- `vite 6.4.3`
- `@vitejs/plugin-vue 5.2.4`

## 間接依存

GoはWailsとgo-crypto由来の依存が多い。`go list -m all` では `github.com/cloudflare/circl v1.6.2`、`golang.org/x/crypto v0.41.0`、`golang.org/x/net v0.42.0` などを確認した。

npmはVite/Vue由来でRollup、esbuild、PostCSS、Vue compiler/runtimeを含む。`npm ls --all --omit=dev` の optional dependency 不足はプラットフォーム別optional依存であり、Linux上では通常の表示である。

## 依存関係のリスク

- `govulncheck` が `github.com/cloudflare/circl@v1.6.2` の `GO-2026-4550` を到達可能として検出。
- `go list -m -u all` で複数の更新候補がある。特に `cloudflare/circl` は `v1.6.4` が候補。
- `npm outdated --long` で `vite 8.1.0`、`@vitejs/plugin-vue 6.0.7` のメジャー更新候補がある。ただし `npm audit --omit=dev` は0件。

## 既知脆弱性の確認方法

```sh
cd src
go run golang.org/x/vuln/cmd/govulncheck@latest ./...
```

```sh
cd src/frontend
npm audit --omit=dev
```

```sh
cd src
go run github.com/securego/gosec/v2/cmd/gosec@latest ./...
```

## パッケージマネージャのロック状況

- Go: `src/go.sum` あり。
- npm: `src/frontend/package-lock.json` あり。
- CI/Release: `npm ci` を使用。
- Go toolはmodule cache利用。Wails CLIはバージョン固定してinstall/cache。

## ビルドスクリプトのリスク

- `scripts/build-windows-from-wsl.sh` はローカル開発用。`wails` と `git` に依存し、引数は `--version` と `--zip` のみ。
- Release workflowはtag名を成果物名・Release本文抽出に使う。現状のtag運用ルールとGitHubのref制約に依存する。

## CI/CDのリスク

- CIは `contents: read` で良い。
- Releaseは `contents: write` が必要だがworkflow全体に設定されている。
- 外部ActionはSHA pinningではない。
- GPG秘密鍵はRelease job内で取り込む。`GNUPGHOME` は一時ディレクトリだが、cleanupは明示されていない。

## リリース成果物作成時のリスク

- 成果物にzip、sha256、exe detached signatureを添付。
- zip内に署名公開鍵URLを入れる。
- 公開鍵ファイル本体の誤混入検査はある。
- SBOMやbuild metadataはない。

## 改ざん耐性

- sha256: zip破損・改ざん確認に有効。ただしsha256自体もReleaseページから取得するため、Releaseページ侵害時にはPGP署名の方が重要。
- PGP署名: exeのdetached signatureあり。
- GitHub Release: tagとRelease workflowに依存。

## SBOM

未作成。ReleaseにCycloneDX/SPDXを添付すると、依存関係監査とユーザー説明がしやすくなる。

## Reproducible Build

未整備。Wails/WebView2/Go/Node/npmのバージョンはworkflowに固定されているが、完全な再現可能ビルド手順やbit-for-bit再現性は未確認。

## SLSA観点

- Source: GitHub repository。
- Build: GitHub-hosted runner。
- Provenance: 未生成。
- Dependency pinning: Go/npm lockあり。GitHub ActionsはSHA pinningなし。
- Signing: exe PGP署名あり。tag署名の必須化は未確認。

## リリース権限

- Release workflow: `contents: write`
- Release secrets: `CI_RELEASE_GPG_PRIVATE_KEY`, `CI_RELEASE_GPG_PASSPHRASE`, `CI_RELEASE_GPG_FINGERPRINT`
- 推奨: GitHub Environmentsの承認、tag作成権限制限、Release job分離。

## シークレット管理

- GitHub Actions secretsを使用。
- アプリ側はWebhook URLをconfig/historyへ保存する。
- 診断データにシークレットが含まれ得るため、平文zipの扱いを修正すべき。

## 外部アクション

- `actions/checkout@v7`
- `actions/setup-go@v6`
- `actions/setup-node@v6`
- `actions/cache@v6`
- `softprops/action-gh-release@v3`

SHA pinningを検討する。
