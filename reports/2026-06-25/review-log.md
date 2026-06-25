# Review Log

調査日: 2026-06-25

対象コミット: `8e46ffdefaff79eb86fd7ef02606c45ff9df4ec7`

## 実施した確認作業

- `reports/2026-06-25/security-audit-prompt.md` の元になった監査プロンプト全体を確認。
- リポジトリ構成、主要ファイル、CI/Release workflowを確認。
- Goコードの入力経路、ファイル処理、ネットワーク処理、暗号処理、ログ処理を確認。
- FrontendコードのURL処理、Wails API呼び出し、HTML挿入リスクを確認。
- Go/npm依存関係の一覧、更新候補、脆弱性チェックを実行。
- テスト、frontend build、gosecを実行。

## 読んだ主要ファイル

- `README.md`
- `RELEASE_NOTES.md`
- `src/SPEC.md`
- `src/SETTINGS_SPEC.md`
- `src/main.go`
- `src/app.go`
- `src/app_diagnostic.go`
- `src/diagnostic_package.go`
- `src/internal/appcore/config.go`
- `src/internal/appcore/diagnostic.go`
- `src/internal/appcore/discord.go`
- `src/internal/appcore/history.go`
- `src/internal/appcore/image.go`
- `src/internal/appcore/processor.go`
- `src/internal/appcore/autophoto.go`
- `src/internal/appcore/update.go`
- `src/internal/appcore/qrcode.go`
- `src/internal/appcore/clipboard.go`
- `src/internal/appcore/clipboard_native_windows.go`
- `src/reveal_windows.go`
- `src/cli_console_windows.go`
- `src/frontend/src/main.js`
- `src/frontend/index.html`
- `src/go.mod`
- `src/frontend/package.json`
- `.github/workflows/ci.yml`
- `.github/workflows/release.yml`
- `scripts/build-windows-from-wsl.sh`
- `scripts/extract-release-notes.mjs`
- `scripts/write-wails-info.mjs`
- `scripts/write-wails-windows-info.mjs`

## 実行したコマンド

```sh
git status --short --branch
sed -n '1,900p' reports/2026-06-25/security-audit-prompt.md
find . -path './src/frontend/node_modules' -prune -o -path './src/build/bin' -prune -o -path './.git' -prune -o -maxdepth 3 -type f -print
rg --glob '!src/frontend/node_modules/**' ...
go test ./...
npm audit --omit=dev
go run golang.org/x/vuln/cmd/govulncheck@latest ./...
go run github.com/securego/gosec/v2/cmd/gosec@latest ./...
go list -m all
npm ls --all --omit=dev
npm run build
go list -m -u all
npm outdated --long
git ls-files
```

## 実行結果の要約

- `go test ./...`: 成功。
- `npm audit --omit=dev`: 0 vulnerabilities。
- `go run github.com/securego/gosec/v2/cmd/gosec@latest ./...`: 0 issues。
- `npm run build`: 成功。
- `go run golang.org/x/vuln/cmd/govulncheck@latest ./...`: 失敗。`GO-2026-4550` を検出。
- `go list -m -u all`: `github.com/cloudflare/circl v1.6.2 [v1.6.4]` など更新候補を確認。
- `npm outdated --long`: `vite` と `@vitejs/plugin-vue` のメジャー更新候補を確認。脆弱性としては検出されていない。

## 実行できなかったコマンド

- Windows実機での `icacls` 確認は未実行。現在の作業環境がWSL/Linuxであるため。
- WailsのWindows exeビルドは今回の監査では未実行。既存CI/Release workflowとローカルbuild scriptを静的確認した。
- Discord/GitHubへの攻撃的・能動的な外部検証は実施していない。プロンプトの制約に従い、第三者サービスへのスキャンや攻撃は行っていない。

## 環境上の制約

- Windows ACL、Explorer、WebView2の実機挙動はWSL上では確認できない。
- `govulncheck` の結果は2026-06-25時点の脆弱性DBに依存する。
- 実Discord Webhookを使った投稿・削除は行っていない。

## 未確認事項

- Windows上でのconfig/history/log/diagnosticsの実ACL。
- Release workflowを実際に走らせた場合の成果物とGitHub Release本文。
- 診断データ生成UIの実機操作。
- WebView2/Wails runtimeの詳細なsandbox設定。
- GPG Release署名鍵の運用手順とアクセス権限。

## 追加で人間が確認すべき事項

- `GO-2026-4550` の実影響と更新後の互換性。
- 診断データに確認用zipを残す前提で、Webhook URLやDiscord tokenなどをダミー化する範囲。
- Webhook tokenを履歴へ保存し続ける必要性。
- Release workflowのenvironment protectionとtag作成権限。
- ユーザー向け文言で、zipではなく `.zip.gpg` を添付する注意が十分か。
