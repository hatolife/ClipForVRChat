# ClipForVRChat セキュリティチェック報告書

調査日: 2026-06-21  
対象コミット: `2c2a2a3 Add VRChat photo auto posting`  
対象範囲: Go/Wails バックエンド、Vue フロントエンド、設定/履歴保存、Discord 連携、VRChat 写真フォルダ監視、CI/Release、依存関係

## 概要

重大なハードコード済み秘密情報や、フロントエンドの `v-html` / `innerHTML` / `eval` のような直接的な XSS 起点は見つかりませんでした。npm 依存関係の既知脆弱性も 0 件でした。

一方で、画像と Discord Webhook を扱うアプリとしては、次のリスクを優先して見直すべきです。

| 優先度 | 主な問題 | 影響 |
| --- | --- | --- |
| 高 | Discord Webhook URL が Discord ドメインに制限されていない | 悪意ある設定ファイル経由で、画像を任意サーバーへ送信される可能性 |
| 高 | Webhook URL / Discord 削除トークン / 履歴が平文保存される | PC内の別ユーザーやバックアップ流出時に投稿・削除権限や画像履歴が漏れる |
| 高 | Go 標準ライブラリに到達可能な既知脆弱性が検出された | HTTPS/HTTP 通信や証明書検証周辺で DoS 等の影響を受ける可能性 |
| 中 | 画像ファイルのサイズ・ピクセル数上限がない | 巨大画像や細工画像でメモリ/CPUを消費し、アプリが固まる可能性 |
| 中 | 履歴内 URL の検証がなく、任意 URL へアクセスする処理がある | ローカルで改ざんされた履歴から内部ネットワークへのリクエストや追跡が可能 |
| 中 | 自動投稿のフォルダ監視が再帰走査かつ対象確認が弱い | 誤設定や大量ファイルで意図しない投稿・負荷増大が起きる |
| 中 | 自動投稿と手動処理の履歴保存が排他制御なし | 履歴・削除フラグが失われる不具合につながる |

## 実施した確認

| 確認 | 結果 |
| --- | --- |
| `go test ./...` | 成功。ただしテストファイルなし |
| `go vet ./...` | 指摘なし |
| `go run golang.org/x/vuln/cmd/govulncheck@latest ./...` | 到達可能な Go 標準ライブラリ脆弱性 8 件 |
| `go run github.com/securego/gosec/v2/cmd/gosec@latest -fmt=json ./...` | 15 件検出 |
| `npm audit --json` | 脆弱性 0 件 |
| `npm run build` | 成功 |
| 秘密情報パターン検索 | 追跡ファイル内に Webhook 実値、秘密鍵、GitHub token らしき文字列は見つからず |

## 詳細 findings

### 1. Webhook URL が任意送信先として使える

優先度: 高

`UploadDiscord` は `webhookURL` をそのまま `http.NewRequest` に渡しています。URL の scheme、host、path が Discord Webhook かどうかを検証していません。

根拠:

- `src/internal/appcore/discord.go:30-60`
- `src/internal/appcore/processor.go:91-95`
- `src/internal/appcore/autophoto.go:67-73`

悪用シナリオ:

1. 攻撃者が `config.json` を配布する。
2. ユーザーがその JSON をアプリへドロップする。
3. 設定内の `webhookUrl` が `https://attacker.example/upload` のような任意 URL でも保存・使用される。
4. 以後、画像処理または VRChat 写真自動投稿で、画像が Discord ではなく攻撃者のサーバーへ送られる。

JSON ドロップで設定ファイルを開ける経路もあります。

- `src/frontend/src/main.js:192-201`
- `src/app.go:197-201`
- `src/main.go:45-57`

推奨対応:

- Discord Webhook URL は `https://discord.com/api/webhooks/{id}/{token}` と `https://discordapp.com/api/webhooks/{id}/{token}` 程度に限定する。
- 設定ファイルから読み込んだ Webhook URL が既存設定と異なる場合は、保存前に明示確認する。
- 自動投稿を有効にする設定ファイルを読み込む場合は、初回保存時に警告を出す。

### 2. Webhook URL と削除用トークンが平文保存される

優先度: 高

`config.json` には Webhook URL が平文で保存されます。さらに `history.json` には Discord 投稿の削除に使える `discordToken` が保存されます。どちらも `0644` で書き込まれるため、Unix 系環境では同一マシン上の他ユーザーから読める可能性があります。Windows でも、バックアップ、同期、共有フォルダ、ログ添付などで流出すると、投稿先チャンネルへの投稿や過去投稿の削除に使われます。

根拠:

- 設定構造: `src/internal/appcore/config.go:35-42`
- 設定保存: `src/internal/appcore/config.go:94-103`
- 履歴構造: `src/internal/appcore/history.go:13-28`
- 履歴保存: `src/internal/appcore/history.go:52-60`
- Discord token の履歴保存: `src/internal/appcore/history.go:73-83`
- `gosec`: `G306` `config.go:103`, `history.go:60`

推奨対応:

- Webhook URL は OS の安全な資格情報ストアに保存する。Windows なら DPAPI または Windows Credential Manager を検討する。
- 最低限、設定ファイルと履歴ファイルは `0600` 相当で保存する。
- `history.json` に削除用 token を保存し続ける必要があるか再検討する。必要なら暗号化する。
- README に「設定ファイルと履歴ファイルを共有しない」注意を明記する。

### 3. Go 標準ライブラリに到達可能な既知脆弱性がある

優先度: 高

`govulncheck` により、ローカルの Go 1.26.1 標準ライブラリに対して到達可能な脆弱性が 8 件検出されました。主に `net/http`、`crypto/tls`、`crypto/x509`、`net` 周辺です。Discord 投稿・削除・履歴 URL チェックで HTTPS/HTTP 通信を使うため、影響経路があります。

検出例:

| ID | 概要 | 修正版 |
| --- | --- | --- |
| GO-2026-5039 | `net/textproto` のエラー文字列エスケープ問題 | Go 1.26.4 |
| GO-2026-5037 | `crypto/x509` の非効率な hostname 解析 | Go 1.26.4 |
| GO-2026-4971 | Windows の `net` で NUL byte による panic | Go 1.26.3 |
| GO-2026-4918 | HTTP/2 transport の無限ループ | Go 1.26.3 |
| GO-2026-4870 | TLS 1.3 KeyUpdate による接続保持 DoS | Go 1.26.2 |

CI/Release は `go-version: "1.24.x"` を使用していますが、`go.mod` は `go 1.25.0` です。この差分により、ローカル確認と配布バイナリの実際の標準ライブラリが揃いません。

根拠:

- `src/go.mod:3`
- `.github/workflows/ci.yml:44-48`
- `.github/workflows/release.yml:21-25`

推奨対応:

- CI/Release の Go を現在の安全な patch version に固定する。
- `go.mod` の Go バージョンと CI の Go バージョンを合わせる。
- `govulncheck ./...` を CI に追加し、リリース前に失敗させる。

### 4. 画像処理にファイルサイズ・ピクセル数上限がない

優先度: 中

画像ファイルは `os.ReadFile` で全体をメモリに読み込み、その後 `image.Decode` します。ファイルサイズやデコード後ピクセル数の上限がないため、巨大画像や圧縮率の高い細工画像でメモリ・CPUを消費する可能性があります。

根拠:

- `src/internal/appcore/image.go:43-48`
- `src/internal/appcore/image.go:51-72`
- `src/internal/appcore/processor.go:61-70`
- `gosec`: `G304` `image.go:44`

影響:

- ドラッグ&ドロップされた画像でアプリが固まる。
- 自動投稿フォルダに巨大画像を置かれると、起動中に繰り返し重くなる。
- 複数画像処理時に UI 応答性が下がる。

推奨対応:

- 読み込み前に `os.Stat` でファイルサイズ上限を設ける。
- `image.DecodeConfig` でピクセル数上限を確認してから `Decode` する。
- 自動投稿ではファイルサイズ・拡張子・MIME を検証し、上限超過を履歴/通知に残す。

### 5. 履歴 URL チェックが任意 URL へリクエストする

優先度: 中

履歴整理機能は、履歴にある URL に対して `HEAD` または `Range: bytes=0-0` の `GET` を送ります。`history.json` はローカルファイルなので、改ざんされた場合に任意 URL へアクセスできます。

根拠:

- `src/internal/appcore/history.go:115-160`
- `src/app.go:163-169`
- UI 実行経路: `src/frontend/src/main.js:315-325`

リスク:

- `http://127.0.0.1:...` や LAN 内 URL へアクセスさせるローカル SSRF 的な挙動。
- 履歴画面で `<img :src="item.thumbnail || item.url">` により、外部 URL への画像取得が発生する。
- 外部サーバーにアプリ起動や履歴画面閲覧のタイミングが伝わる。

根拠:

- `src/frontend/src/main.js:487-491`

推奨対応:

- 履歴に保存・表示・検査する URL を Discord CDN または Discord attachments の HTTPS URL に限定する。
- `http://`、localhost、private IP、file scheme、data URL の扱いを明確に拒否する。
- 履歴ファイルの改ざんを想定し、読み込み時に URL validation を行う。

### 6. VRChat 写真自動投稿の誤設定リスク

優先度: 中

自動投稿は指定フォルダを 2 秒ごとに再帰走査し、新しい `.png`, `.jpg`, `.jpeg`, `.webp` を投稿します。対象フォルダが広すぎる場合、意図しない画像投稿や大量ファイル走査が発生します。

根拠:

- `src/internal/appcore/autophoto.go:24-42`
- `src/internal/appcore/autophoto.go:83-109`
- `src/internal/appcore/autophoto.go:111-122`
- `src/frontend/src/main.js:581-598`

推奨対応:

- 自動投稿を有効にする際、監視対象フォルダを表示して確認させる。
- デフォルトの VRChat 写真フォルダ以外を選んだ場合は警告する。
- 走査対象ファイル数や 1 tick あたり処理件数に上限を設ける。
- 初回起動時の既存ファイル無視は良い設計なので維持する。

### 7. 履歴保存と UI 状態更新に排他制御がない

優先度: 中

手動処理、履歴削除、自動投稿 watcher が同じ `a.state` と `history.json` を更新しますが、mutex やファイルロックがありません。自動投稿中に手動処理や履歴削除を行うと、履歴エントリや `cleared` / `discordDeleted` フラグが失われる可能性があります。

根拠:

- 手動処理の履歴保存: `src/app.go:242-264`
- 自動投稿の履歴保存と state 更新: `src/app.go:349-374`
- 履歴の load-modify-save: `src/internal/appcore/history.go:63-88`, `src/internal/appcore/history.go:91-112`

影響:

- Discord から削除済みなのに履歴に反映されない。
- 履歴整理中に新規エントリが消える。
- UI 上の結果一覧と `history.json` が食い違う。

推奨対応:

- `App` に mutex を持たせ、`state` と履歴ファイル更新を直列化する。
- 履歴保存は「読み込み、変更、保存」を単一関数で排他実行する。
- Discord 削除は 1 件ずつ結果を記録し、部分成功でも履歴に保存する。

### 8. 設定ファイルの読み込み元と保存先が任意パスになり得る

優先度: 中

JSON を 1 つだけ渡す/ドロップすると、その JSON を設定ファイルとして開き、`a.configPath` がそのパスへ変更されます。その後の保存や履歴ファイルは、その JSON と同じディレクトリに作られます。ローカルアプリとしては自然な仕様ですが、信頼できない `config.json` を開いた場合の影響が大きいです。

根拠:

- `src/main.go:45-57`
- `src/app.go:197-201`
- `src/internal/appcore/history.go:30-35`
- `src/internal/appcore/config.go:78-103`

リスク:

- 意図しない場所へ `config.json` / `history.json` を作成・上書きする。
- 攻撃者の Webhook URL、出力フォルダ、自動投稿設定を取り込む。

推奨対応:

- 外部設定ファイルを開いた場合は「この設定ファイルを信頼するか」を明示する。
- 保存前に Webhook URL、出力先、自動投稿 ON/OFF の差分を確認表示する。
- 履歴保存先は exe 横またはユーザーデータディレクトリに固定することを検討する。

### 9. エラーメッセージに外部レスポンス本文を含めている

優先度: 低から中

Discord 投稿・削除に失敗した場合、レスポンス本文をそのままエラーに含めます。通常の Discord API レスポンスなら大きな問題になりにくいですが、Webhook URL が任意 URL にできるため、攻撃者サーバーから長大または紛らわしい本文を返されると UI 表示やログ添付時の情報混乱につながります。

根拠:

- `src/internal/appcore/discord.go:66-72`
- `src/internal/appcore/discord.go:109-110`

推奨対応:

- エラー本文は長さ制限を設ける。
- Webhook URL を Discord に限定する。
- ユーザー向けメッセージと詳細ログを分ける。

### 10. `OpenURL` に scheme/host 制限がない

優先度: 低

現在 UI から渡される URL はアプリ内の固定値または OSS ライセンスの固定値です。ただし Wails に公開されている `OpenURL` は任意文字列をブラウザに渡します。フロントエンドが将来ユーザー入力 URL を渡すようになると、危険な scheme や意図しないアプリ起動につながる可能性があります。

根拠:

- `src/app.go:59-60`
- `src/frontend/src/main.js:336-338`

推奨対応:

- `https://` のみ許可する。
- host allow-list を持つ。
- 将来、履歴 URL やユーザー入力を `OpenURL` に渡さない設計を維持する。

### 11. 配布物の信頼性確認が弱い

優先度: 低から中

Release workflow は Windows exe を zip 化して GitHub Release に添付しますが、署名、チェックサム、SBOM、依存脆弱性チェック、SLSA/provenance はありません。個人配布アプリとしてはよくある状態ですが、Windows exe を配布するならユーザーが改ざん有無を確認しづらいです。

根拠:

- `.github/workflows/release.yml:49-70`

推奨対応:

- zip の SHA256 チェックサムを release artifact に追加する。
- 可能ならコード署名を導入する。
- SBOM を生成して添付する。
- `govulncheck` と `npm audit` を CI/Release に入れる。

### 12. `gosec` のその他検出

優先度: 低

`src/tools/make_icon.go` は配布アプリ本体ではなくアイコン生成ツールですが、`gosec` は次を検出しました。

- `G115`: `uint32 -> uint8` 変換。`src/tools/make_icon.go:155-159`
- `G304`: 変数パスへの `os.Create`。`src/tools/make_icon.go:166-168`
- `G301`: `build/icons` の `0755`。`src/tools/make_icon.go:17-20`

固定入力の開発用ツールなので実運用リスクは低いです。ただし CI や配布前生成に組み込む場合は、警告を抑制するか修正しておくとよいです。

## 問題が見つからなかった点

- 追跡ファイル内に Discord Webhook 実値、秘密鍵、GitHub token らしき文字列は見つかりませんでした。
- フロントエンドに `v-html`、`innerHTML`、`eval`、任意コマンド実行は見つかりませんでした。
- Wails の file drop は有効ですが、WebView drop は無効化されています。`src/main.go:140-143`
- `npm audit` では脆弱性 0 件でした。
- `go vet` は指摘なしでした。

## 推奨対応順

1. Webhook URL を Discord の HTTPS Webhook URL に限定する。
2. `config.json` / `history.json` の保存権限を絞り、Webhook URL と削除 token の暗号化または OS 資格情報ストア保存を検討する。
3. Go のビルドバージョンを安全な patch version に揃え、`govulncheck` を CI/Release に追加する。
4. 画像のファイルサイズ・ピクセル数・処理件数に上限を設ける。
5. 履歴 URL の validation を追加し、Discord CDN 以外への自動アクセスを避ける。
6. 自動投稿 ON と外部設定ファイル読み込みに確認 UI を追加する。
7. 履歴ファイルと `App.state` 更新に排他制御を入れる。
8. Release にチェックサム、可能ならコード署名/SBOMを追加する。

## 参考: 実行コマンド

```bash
go test ./...
go vet ./...
go run golang.org/x/vuln/cmd/govulncheck@latest ./...
go run github.com/securego/gosec/v2/cmd/gosec@latest -fmt=json ./...
npm audit --json
npm run build
```
