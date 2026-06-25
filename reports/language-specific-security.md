# Language Specific Security

調査日: 2026-06-25

## Go

### 確認した観点

- 入力検証: 画像サイズ、ピクセル数、Webhook URL、Discord添付URL、zip形式。
- エラー処理: 多くのファイル・HTTP処理はエラーを返す。診断ログ書き込みは失敗を握りつぶす設計。
- パス正規化: 保存名は `filepath.Base` を使用。履歴削除では絶対パスを許す。
- HTTP timeout: Discord投稿60秒、Discord削除30秒、GitHub更新確認6秒、URL availability15秒。
- 暗号: ProtonMail/go-crypto OpenPGPを使用。独自暗号実装は確認できない。
- unsafe: Windows clipboardで `unsafe.Slice` を使用。Win32 API境界として局所化されている。
- 並行処理: `App` はmutexで状態を保護。自動処理watcherはcontext cancelあり。
- 乱数: アプリコードで直接の乱数利用は確認できない。暗号ライブラリ内部に依存。
- ログ: ユーザー操作と処理結果を詳細に記録。秘密情報・パスの扱いに注意。

### 良い点

- `DecodeImageWithLimit()` で入力バイト数とピクセル数を確認してからdecodeしている。
- Discord Webhook URLは `https`、ホスト、パス形式を検証している。
- Discord添付URLは信頼ホストと `/attachments/` パスを確認している。
- HTTP clientにtimeoutが設定されている。
- `gosec` は0件。

### 問題・注意点

- `govulncheck` が `github.com/cloudflare/circl@v1.6.2` の到達可能脆弱性を検出。
- 診断用の平文zipにWebhook URLや履歴tokenが残る。
- `history.json` の `OutputPath` に絶対パスが入ると削除対象になる。
- Windowsファイル権限は実機ACL確認が必要。

## JavaScript / Vue

### 確認した観点

- DOM XSS: `innerHTML`、`insertAdjacentHTML`、`v-html` は確認できない。
- URL handling: `openURL()` はGo側へURLを渡す。
- 秘密情報表示: Webhook URL入力は `type="password"`。
- 依存関係: `npm audit --omit=dev` は0件。

### 良い点

- Vueテンプレートのテキストバインディング中心。
- 外部URLは主にアプリ内定数。
- 診断データ作成中のUIロックがある。

### 問題・注意点

- `openURL()` の許可ホスト制限はGo側にもJS側にもない。
- GitHub Release API由来URLを表示・オープンする導線がある。GitHub APIを信頼しているため現状大きなリスクではないが、許可リスト化が望ましい。
- `npm outdated` でVite系のメジャー更新候補がある。脆弱性ではないが保守計画が必要。

## YAML / GitHub Actions

### 確認した観点

- 権限: CIは `contents: read`、Releaseは `contents: write`。
- 依存関係監査: npm audit、govulncheck。
- 署名: GPG秘密鍵をRepository Secretsから取り込み、exe detached signatureを作成。
- 成果物検査: zip/署名/sha256の存在確認、不要asc混入検査。

### 良い点

- Release前にgo test、npm audit、govulncheckを実行。
- Release成果物にsha256とPGP署名がある。
- Wails CLI cacheでCI高速化済み。

### 問題・注意点

- Release workflow全体が `contents: write`。
- 外部ActionがSHA pinningではなくメジャーバージョン指定。
- SBOM生成はない。

## Bash / PowerShell / Node scripts

### 確認した観点

- `scripts/build-windows-from-wsl.sh`: `set -euo pipefail`、引数解析あり。
- Release PowerShell: tag名をバージョン判定に使用し、zip名などに使用。
- Node scripts: Release notes抽出とWails metadata生成。

### 良い点

- shell scriptは未知引数を拒否。
- Release notes抽出は見出しを正規表現escapeしている。

### 注意点

- Release workflowのtag名はファイル名に使われる。GitHub tag名の制約に依存しているが、workflow上でも安全な文字種へ正規化するとより堅い。
