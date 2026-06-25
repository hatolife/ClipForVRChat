# Project Profile

調査日: 2026-06-25

対象ブランチ: `fix/discord-upload-disabled-posts`

対象コミット: `8e46ffdefaff79eb86fd7ef02606c45ff9df4ec7`

備考: 監査用 issue、todo、報告書追加前の製品コードを対象に確認した。

## 構成

- プロジェクト種別: Windows向けGUIデスクトップアプリ。補助的に `--version`、`--help`、zip暗号化CLIを持つ。
- 主要言語: Go
- 補助言語: JavaScript、CSS、HTML、PowerShell、Bash、GitHub Actions YAML
- ビルドシステム: Wails v2、Vite
- パッケージマネージャ: Go Modules、npm
- フレームワーク: Wails v2、Vue 3、Vite
- 対象実行環境: Windows desktop。開発・ローカルビルドはWSLからのWindowsビルド補助あり。
- 対象OS: 配布対象はWindows。CIはWindows/Ubuntu。
- 配布形態: GitHub ReleaseのWindows amd64 zip。zip、sha256、exe detached PGP署名を添付。
- エントリポイント: `src/main.go`
- GUIバインド: `src/app.go` の `App` methods
- 設定ファイル: exe同階層の `config.json`
- 履歴ファイル: exe同階層の `history.json`
- ログ: exe同階層の `logs/YYYY-MM-DD.log`
- 診断データ: exe同階層の `diagnostics/<timestamp>/`

## 入力経路

- GUIへの画像ファイルドラッグ&ドロップ
- CLI引数の画像ファイル、`.json` 設定ファイル、`.zip` ファイル
- クリップボード画像
- `config.json`
- `history.json`
- 設定画面のWebhook URL、出力先、監視フォルダ
- VRChat写真フォルダ、Screenshotsフォルダの自動スキャン
- Discord Webhook APIレスポンス
- GitHub Releases APIレスポンス
- QRコード内URL

## 出力経路

- ローカル縮小画像ファイル
- `config.json`
- `history.json`
- `logs/YYYY-MM-DD.log`
- 診断用zipと暗号化済み `.zip.gpg`
- Discord Webhookへの画像投稿
- Discord Webhookメッセージ削除リクエスト
- クリップボードへのURLコピー
- ブラウザでのGitHub/BOOTH/Discordヘルプ/Twitter URL表示
- Explorerでのファイル選択表示

## ネットワーク

- Discord Webhook投稿: `https://discord.com` / `https://discordapp.com`
- Discord添付URL確認: `https://cdn.discordapp.com` / `https://media.discordapp.net`
- Discord Webhookメッセージ削除
- GitHub Releases API: `https://api.github.com/repos/hatolife/ClipForVRChat/releases/latest`
- ブラウザ起動による外部サイト表示

## ファイルシステム

- 画像読み込み、画像保存、履歴・設定・ログ書き込み
- 診断データ作成時にconfig/history/log/exeを収集
- 履歴画面からローカル保存ファイルを削除
- 自動処理フォルダの再帰スキャン
- Release workflowでzip作成、署名、検査

## プロセス・IPC・動的読み込み

- 外部プロセス起動: アプリ本体から直接のシェル実行は確認できない。ブラウザ起動とExplorer表示はWails runtime / Windows Shell API経由。
- IPC: WailsのGo/Frontendブリッジ、Windowsクリップボード、Windows Shell API。
- 動的ライブラリ読み込み: Windows API呼び出しで `user32.dll`、`kernel32.dll`、`shell32.dll`、`ole32.dll` を `syscall.NewLazyDLL` で参照。
- サーバー機能: 配布アプリとして外部から待ち受けるHTTPサーバーは確認できない。Wails開発時のVite dev serverは開発用途。

## 認証・認可

- アプリ独自の認証・認可はない。
- Discord Webhook URL自体が投稿・削除権限を持つ秘密情報。
- GitHub Actions Release署名鍵はRepository Secretsで管理。

## 暗号・署名・ハッシュ

- 診断zip暗号化: ProtonMail/go-crypto OpenPGP
- 診断zip暗号化公開鍵: `poppo@hato.life` の公開鍵をソースへ埋め込み
- 起動ログ: exe SHA-256を記録
- Release: zip SHA-256、exe detached PGP署名
- 自動更新: なし。更新通知のみで、ダウンロードやインストールは行わない。

## 依存関係

- Go direct dependency: Wails、go-arg、go-crypto、imaging、flock、gozxing、go-qrcode、clipboard、x/image
- JS direct dependency: Vue、Vite、@vitejs/plugin-vue
- lock file: `src/go.sum`、`src/frontend/package-lock.json`

## CI/CD

- CI: `.github/workflows/ci.yml`
  - docs check
  - npm ci / npm audit
  - frontend build
  - go test
  - govulncheck
  - Wails build
- Release: `.github/workflows/release.yml`
  - tag `v*` と手動実行
  - npm audit、govulncheck、test、build
  - PGP署名
  - zip/sha256作成
  - Release本文抽出
  - GitHub Release作成

## セキュリティ上重要なファイル

- `src/main.go`: CLI、起動引数、起動ロック、初期状態
- `src/app.go`: Wailsブリッジ、履歴削除、設定保存、URL表示
- `src/diagnostic_package.go`: 診断データ収集、zip、OpenPGP暗号化
- `src/app_diagnostic.go`: 起動ログ、設定ログ、exeハッシュ
- `src/internal/appcore/config.go`: 設定保存権限、デフォルト値
- `src/internal/appcore/history.go`: 履歴保存、Discord削除状態、ローカルファイル削除
- `src/internal/appcore/discord.go`: Webhook URL検証、Discord投稿・削除
- `src/internal/appcore/image.go`: 画像入力サイズ制限、保存先処理
- `src/internal/appcore/autophoto.go`: 自動スキャン制限
- `src/internal/appcore/update.go`: GitHub Release確認
- `src/internal/appcore/clipboard_native_windows.go`: Windows clipboard unsafe境界
- `src/reveal_windows.go`: Windows Shell API
- `.github/workflows/ci.yml`: CI security checks
- `.github/workflows/release.yml`: Release権限、署名、成果物作成
- `src/frontend/src/main.js`: GUI入力、Wails API呼び出し
