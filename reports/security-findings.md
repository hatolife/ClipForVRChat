# Security Findings

調査日: 2026-06-25

対象コミット: `8e46ffdefaff79eb86fd7ef02606c45ff9df4ec7`

## SEC-001: 暗号関連依存関係に到達可能な既知脆弱性がある

- 重大度: Medium
- 分類: Vulnerable Dependency / Cryptography
- CWE: CWE-327, CWE-682
- 該当ファイル: `src/go.mod`, `src/diagnostic_package.go`
- 該当箇所: `github.com/cloudflare/circl@v1.6.2`, `encryptDiagnosticZip()`
- 確信度: High

### 問題

`govulncheck ./...` が `GO-2026-4550` を検出した。`github.com/cloudflare/circl@v1.6.2` のsecp384r1 CombinedMult計算に関する脆弱性で、修正版は `v1.6.3` 以降である。

検出経路は `src/diagnostic_package.go:156` の `openpgp.Encrypt` と、同ファイルの公開鍵読み込み処理である。現在の埋め込み公開鍵はcv25519系だが、脆弱な暗号モジュールがアプリの暗号処理依存として到達可能である点はRelease前に解消すべきである。

### 攻撃者視点

診断データ暗号化処理やzip暗号化CLIを利用するユーザーに対し、暗号処理の実装不備を攻撃チェーンへ組み込む余地がある。現時点で、このアプリの埋め込み鍵だけで秘密情報を直接復号できるとまでは断定しない。

### 影響

診断データの秘匿性または暗号処理の信頼性低下。CI/Release workflowの `govulncheck` 失敗によるRelease停止。

### 再現または確認方法

```sh
cd src
go run golang.org/x/vuln/cmd/govulncheck@latest ./...
```

### 推奨修正

`github.com/cloudflare/circl` を `v1.6.3` 以上へ更新する。直接requireがない場合でも、`go get github.com/cloudflare/circl@v1.6.4` で `go.mod` / `go.sum` を更新できるか確認する。必要に応じて `github.com/ProtonMail/go-crypto` も更新候補を確認する。

### 修正例

```sh
cd src
go get github.com/cloudflare/circl@v1.6.4
go mod tidy
go test ./...
go run golang.org/x/vuln/cmd/govulncheck@latest ./...
```

### 修正後の確認方法

`govulncheck ./...` が0件で終了すること。診断パッケージ作成とzip復号テストが通ること。

## SEC-002: 診断データの確認用zipにWebhook URLやDiscord tokenが平文で残る

- 重大度: Medium
- 分類: Sensitive Data Exposure
- CWE: CWE-312, CWE-532
- 該当ファイル: `src/diagnostic_package.go`, `src/internal/appcore/history.go`, `src/internal/appcore/config.go`
- 該当箇所: `prepareDiagnosticDataDirectory()`, `addRedactedTextFileToDirectory()`
- 確信度: High

### 問題

不具合報告用データは `diagnostics/<timestamp>/` に確認用zipと暗号化済み `.zip.gpg` を作成する。zip内のテキストはパスを `%USERPROFILE%` などへ置換するが、Webhook URLや履歴内の `discordToken` は置換対象ではない。

そのため、ユーザーが誤って暗号化前zipをGitHub Issueなどへ添付すると、Discord Webhook URLや削除用token、履歴URL、ログ内容が公開され得る。

### 攻撃者視点

公開IssueやSNS投稿から誤添付されたzipを入手し、Webhook URLを使ってDiscordチャンネルへ投稿したり、履歴内tokenとmessage IDで投稿済みメッセージ削除を試みる可能性がある。

### 影響

Discord Webhookの不正利用、投稿済み画像の削除、ユーザーのローカルパス・履歴・利用状況の漏えい。

### 再現または確認方法

Webhook URLを設定した状態で不具合報告用データを生成し、確認用zip内の `config.json` と `history.json` にWebhook URL/tokenが含まれるか確認する。

### 推奨修正

確認用zipは「ユーザー確認用の安全なzip」として扱い、Webhook URL、Discord token、message ID、必要に応じてQR URLやファイル名をマスクする。暗号化済み `.zip.gpg` だけに完全な情報を入れる場合でも、平文zipには明確に安全化済みデータだけを残す設計にする。

### 修正例

- `sanitizeDiagnosticJSON()` を追加し、`config.json` の `discord.webhookUrl`、`autoPhoto.webhookUrl`、`screenshotAutoPost.webhookUrl` を `"<redacted>"` にする。
- `history.json` の `discordToken` を `"<redacted>"` にする。
- `manifest.json` に「平文zipは秘密情報を除外済み」と記録する。

### 修正後の確認方法

診断zip内に `discord.com/api/webhooks`、`discordapp.com/api/webhooks`、`discordToken` の生値が含まれないことをテストで確認する。

## SEC-003: 履歴ファイルにDiscord削除用tokenを長期保存している

- 重大度: Medium
- 分類: Secret Storage
- CWE: CWE-522, CWE-922
- 該当ファイル: `src/internal/appcore/history.go`, `src/internal/appcore/discord.go`
- 該当箇所: `HistoryEntry.DiscordToken`, `AddResultsToHistory()`, `DeleteDiscordMessage()`
- 確信度: High

### 問題

Discord投稿後、削除機能のためにWebhook ID、message ID、tokenを `history.json` へ保存している。ファイルは `0600` 相当で保存されるが、Windows ACLとして期待通り最小化されるかは要確認であり、診断zipにも含まれ得る。

### 攻撃者視点

同一ユーザー権限で `history.json` を読めるマルウェアや別ユーザーが、Webhook tokenを流用してDiscord投稿や削除を行う可能性がある。

### 影響

Discord Webhookの不正利用、履歴内投稿の削除、Webhook URLの再利用。

### 再現または確認方法

Discord投稿後の `history.json` を確認し、`discordToken` が保存されることを確認する。

### 推奨修正

token保存を最小化する。選択肢として、削除機能を維持する場合はWindows Credential ManagerやDPAPIで暗号化する、または削除可能期間を短くしてtokenを破棄する。診断データには既定で含めない。

### 修正例

- `history.json` にはmessage IDとURLだけを保存し、削除時に現在の設定Webhook URLからtokenを再取得する。
- 取得不能な場合は削除不可としてUI表示する。
- token保存が必要な場合はOS保護ストレージへ移す。

### 修正後の確認方法

`history.json` に生tokenが保存されないこと。既存履歴の移行処理でtokenが削除または保護されること。

## SEC-004: Release workflowの権限とAction pinningに改善余地がある

- 重大度: Medium
- 分類: Supply Chain / CI/CD Hardening
- CWE: CWE-829
- 該当ファイル: `.github/workflows/release.yml`, `.github/workflows/ci.yml`
- 該当箇所: `permissions: contents: write`, `uses: actions/...@v*`, `uses: softprops/action-gh-release@v3`
- 確信度: Medium

### 問題

Release workflowはworkflow全体で `contents: write` を持つ。Release作成には必要だが、ビルド・テスト・署名などの前段stepにも同じ権限が付与される。また、外部Actionはメジャーバージョンタグで参照されており、コミットSHA pinningではない。

### 攻撃者視点

外部Action、依存するAction配布経路、またはworkflow変更権限が侵害された場合、Release作成権限や成果物改ざんへつながる可能性がある。

### 影響

Release成果物改ざん、意図しないRelease作成、GitHub tokenの悪用。

### 再現または確認方法

`.github/workflows/release.yml` の権限と `uses:` 参照を確認する。

### 推奨修正

job単位で最小権限を設定し、可能ならRelease作成jobとビルドjobを分ける。外部ActionはSHA pinningを検討し、Dependabot/Renovateで更新運用する。Release environment protectionも検討する。

### 修正例

- build/test job: `contents: read`
- release upload job: `contents: write`
- `uses: softprops/action-gh-release@<commit-sha>`

### 修正後の確認方法

workflow権限がjob単位になっていること。Release workflowが従来通り成果物を作成し、署名・zip検査が通ること。

## SEC-005: WailsブリッジのURLオープン処理が任意URLを受け入れる

- 重大度: Low
- 分類: URL Handling / UX Security
- CWE: CWE-939
- 該当ファイル: `src/app.go`, `src/frontend/src/main.js`
- 該当箇所: `OpenURL(url string)`
- 確信度: Medium

### 問題

`OpenURL` は受け取ったURLをそのまま `runtime.BrowserOpenURL` へ渡す。現状の呼び出し元はアプリ内定数やGitHub Releases API由来のURLであり、直接の任意入力は限定的である。ただし将来、QR URLや履歴URLを開く導線が追加されると、任意スキームや意図しない外部URLを開くリスクが増える。

### 攻撃者視点

ユーザーに細工したURLを開かせる、または将来のUI変更でQR URLを直接開く導線が入った場合に、フィッシングや外部アプリ起動へ誘導する。

### 影響

フィッシング、外部アプリ起動、ユーザー誤操作。

### 再現または確認方法

`src/app.go:72` がURL検証なしでブラウザ起動していることを確認する。

### 推奨修正

`OpenURL` を用途別メソッドに分け、GitHub、BOOTH、Discordヘルプ、Twitterなど許可ホストのみ開く。任意URLを開く必要がある場合は `https` のみに制限し、確認ダイアログを表示する。

### 修正例

```go
func (a *App) OpenTrustedURL(raw string) error {
    parsed, err := url.Parse(raw)
    if err != nil || parsed.Scheme != "https" {
        return fmt.Errorf("開けないURLです")
    }
    switch strings.ToLower(parsed.Hostname()) {
    case "github.com", "hatolife.booth.pm", "support.discord.com", "x.com":
    default:
        return fmt.Errorf("許可されていないURLです")
    }
    runtime.BrowserOpenURL(a.ctx, parsed.String())
    return nil
}
```

### 修正後の確認方法

許可URLは開け、`file:`, `javascript:`, `ms-*:`、未知ホストが拒否されるテストを追加する。

## SEC-006: ローカル削除処理が履歴内の絶対パスを信頼する

- 重大度: Low
- 分類: Path Handling
- CWE: CWE-22, CWE-73
- 該当ファイル: `src/internal/appcore/history.go`
- 該当箇所: `ResolveHistoryOutputPath()`, `removeHistoryOutputFile()`
- 確信度: Medium

### 問題

履歴からローカル保存ファイルを削除する際、`OutputPath` が絶対パスの場合はそのまま削除対象になる。通常はアプリが作成した履歴であり、同一ユーザー権限内の操作に限られるが、`history.json` を改ざんできるローカル攻撃者や誤編集により、アプリ外のファイル削除を誘導できる。

### 攻撃者視点

ユーザーの `history.json` を改ざんし、履歴画面で「ローカルから削除」を押させることで任意のユーザーファイル削除を狙う。

### 影響

同一ユーザー権限内の任意ファイル削除補助。権限昇格には直結しない。

### 再現または確認方法

`history.json` の `outputPath` に既存ファイルの絶対パスを設定し、履歴画面の削除操作を確認する。

### 推奨修正

アプリが管理するoutputディレクトリ配下のファイルだけを削除対象にする。履歴保存時に `OutputPath` と管理ディレクトリを正規化し、削除前に `filepath.Rel` で配下判定を行う。

### 修正例

- `ResolveHistoryOutputPath` とは別に `ResolveManagedOutputPath` を作る。
- `config.Image.OutputDirectory` 配下以外は削除不可にする。

### 修正後の確認方法

output配下は削除でき、output外の絶対パスは削除拒否されるテストを追加する。

## SEC-007: 診断ログに詳細な操作情報とローカル情報が蓄積される

- 重大度: Low
- 分類: Logging / Privacy
- CWE: CWE-532
- 該当ファイル: `src/internal/appcore/diagnostic.go`, `src/app_diagnostic.go`, `src/app.go`, `src/internal/appcore/processor.go`
- 該当箇所: `AppendDiagnosticLog()`, `logStartupLocked()`, `logProcessingStartLocked()`
- 確信度: High

### 問題

診断ログには起動exeパス、SHA-256、設定サマリ、画面遷移、ボタンクリック、処理元、ファイル名、出力パス、QR検出結果などが保存される。Webhook URLは設定サマリで生値を避けているが、ログ全体として利用状況やローカルパスが含まれ得る。

### 攻撃者視点

ローカルのログファイルまたは誤添付された診断zipから、ユーザー名、画像名、利用タイミング、QR URLなどの情報を収集する。

### 影響

プライバシー情報の漏えい。直接のコード実行や権限昇格にはつながらない。

### 再現または確認方法

アプリを起動して画像処理後、`logs/YYYY-MM-DD.log` を確認する。

### 推奨修正

ログ出力時点でパスを環境変数表記へ置換する。QR URLやファイル名は必要性に応じてハッシュ化または件数のみ記録するモードを用意する。診断データ作成時だけでなく、ログファイル自体の内容を最小化する。

### 修正例

- `AppendDiagnosticLog` の呼び出し前に `DiagnosticLogRedactor` を通す。
- `process result` は件数中心にし、詳細はDebugモード時だけ出す。

### 修正後の確認方法

ログに `%USERPROFILE%` などの置換が適用されること。Webhook URLやDiscord tokenがログへ出ないこと。

## SEC-008: Windows上のファイル権限がGoのモード指定だけに依存している

- 重大度: Info
- 分類: Platform Hardening
- CWE: CWE-732
- 該当ファイル: `src/internal/appcore/config.go`, `src/internal/appcore/diagnostic.go`
- 該当箇所: `WritePrivateFile()`, `AppendDiagnosticLog()`
- 確信度: Medium

### 問題

設定、履歴、ログ、診断zipは `0600` 相当で保存されるが、Windows上でGoのfile mode指定が期待するDACL最小化になっているかは実機確認が必要である。

### 攻撃者視点

同一PC上の別ユーザーやバックアップ/同期ソフトが設定・履歴・ログを読める場合、Webhook URLや利用履歴を取得できる。

### 影響

ローカル情報漏えい。

### 再現または確認方法

Windows実機で `icacls config.json history.json logs diagnostics` を確認する。

### 推奨修正

Windows専用実装でDACLを明示する、または保存先をユーザープロファイル配下のアプリデータディレクトリへ移すことを検討する。

### 修正例

Windowsでは `golang.org/x/sys/windows` を使い、現在ユーザーとAdministrators/Systemのみ許可するDACLを設定する。

### 修正後の確認方法

Windows実機で `icacls` の出力をテスト・手順化する。

## SEC-009: SBOMと再現可能ビルドの整備が未完了

- 重大度: Info
- 分類: Supply Chain Transparency
- CWE: CWE-1104
- 該当ファイル: `.github/workflows/release.yml`
- 該当箇所: Release artifact generation
- 確信度: High

### 問題

Releaseではzip、sha256、PGP署名が作成されるが、SBOM、依存関係一覧、ビルド環境メタデータ、再現可能ビルド手順は成果物として添付されない。

### 攻撃者視点

利用者や第三者監査者が、配布物に含まれる依存関係とソースの対応を追跡しづらい。

### 影響

脆弱性対応時の影響範囲確認や配布物検証が遅れる。

### 再現または確認方法

Release workflowの成果物一覧を確認する。

### 推奨修正

CycloneDXやSPDX形式のSBOMをReleaseに添付する。Go/npmの依存関係一覧、ビルド時Go/Node/Wailsバージョン、commit SHAを成果物メタデータとして残す。

### 修正例

- `cyclonedx-gomod` と `cyclonedx-npm` の導入。
- `dist/build-metadata.json` の生成。

### 修正後の確認方法

Release assetsにSBOMとbuild metadataが添付されること。
