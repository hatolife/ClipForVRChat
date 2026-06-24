# ClipForVRChat プロダクト問題点チェック報告書

調査日: 2026-06-24

対象コミット: `8e93ac1 docs(cli): update v0.1.6 behavior`

対象範囲: Go/Wails バックエンド、Vue フロントエンド、設定/履歴保存、Discord連携、自動投稿、CLI、CI/Release、README/SPEC/Release Notes、依存関係

## 概要

現行コードは、過去のセキュリティレビューで高リスクだった Webhook URL の任意送信、履歴URLの任意アクセス、画像サイズ上限、ファイル権限、Go標準ライブラリ脆弱性については大きく改善済みです。`go test`、フロントエンドビルド、npm監査、`govulncheck` も通っています。

一方で、v0.1.6 リリース前に優先して確認すべき問題が残っています。特に、Release workflow が作る成果物と README / アプリ内About / SPEC の説明が矛盾しており、PGP検証手順がユーザーに誤案内される状態です。

## 実施した確認

| 確認 | 結果 |
| --- | --- |
| `go test ./...` | 成功 |
| `go vet ./...` | 成功。サンドボックス内ではWSL一時ディレクトリエラーが出たため通常環境で再実行 |
| `npm run build` | 成功 |
| `npm audit --omit=dev` | 既知脆弱性 0 件 |
| `go run golang.org/x/vuln/cmd/govulncheck@latest ./...` | 到達可能な脆弱性 0 件。require module上の未到達脆弱性 1 件 |
| `go run github.com/securego/gosec/v2/cmd/gosec@latest ./...` | `G304` 1 件。診断ログ出力先の可変パス |
| `node scripts/extract-release-notes.mjs v0.1.6 ...` | 成功。v0.1.6 の本文を抽出可能 |
| `src/build/bin/ClipForVRChat.exe --version` | WSL/UNC実行環境の制約でこちらでは安定確認不可。ユーザー側の実機確認結果を優先 |

## 優先度付き Findings

### 1. Release成果物とドキュメントのPGP説明が矛盾している

優先度: 高

Release workflow は zip 内に `.asc` を入れず、公開鍵 `.asc` も添付しない設計です。実際に `dist/ClipForVRChat` へ入れるのは exe、README、LICENSE、`Release-signing-public-key.url` で、Release asset は zip、zip.sha256、exe.asc のみです。

根拠:

- `.github/workflows/release.yml:87-93`
- `.github/workflows/release.yml:120-123`
- `.github/workflows/release.yml:128-150`
- `.github/workflows/release.yml:165-168`

しかし README、アプリ内About、SPEC は、zip内に `ClipForVRChat.exe.asc` と公開鍵 `.asc` があり、Releaseにも公開鍵 `.asc` が添付されると説明しています。

根拠:

- `README.md:31-33`
- `src/frontend/src/main.js:697-707`
- `src/SPEC.md:658-670`

影響:

- ユーザーがREADMEやアプリ内案内どおりに検証しようとしても、必要と書かれているファイルが見つからない。
- Release workflow の検証条件と仕様書が反対のことを言っているため、今後の修正で意図しない成果物混入または説明不足が再発しやすい。

推奨対応:

- 現行workflowに合わせるなら、README / About / SPEC を「zip内は exe、README、LICENSE、公開鍵URLショートカット。署名ファイルはRelease assetの `.exe.asc`」へ統一する。
- 逆にzip内に署名ファイルを同梱する方針なら、Release workflow の検証条件と添付物一覧を仕様に合わせて戻す。

### 2. Discord削除が途中で失敗すると、成功済み削除の履歴が保存されない

優先度: 中

`DeleteDiscordHistoryEntries` は、複数選択された履歴を順番に削除し、全件完了した後で一度だけ `SaveHistory` します。途中の1件で `DeleteDiscordMessage` が失敗すると、その前に成功して `DiscordDeleted = true` にした履歴も保存されません。

根拠:

- `src/app.go:171-196`

影響:

- Discord上では削除済みなのに、アプリ履歴では未削除のまま残る。
- 次回削除時に 404 や削除済み状態とぶつかり、ユーザーが状況を判断しづらい。

推奨対応:

- 1件削除が成功するたびに保存する。
- または、全件処理を続行して成功分は保存し、失敗分だけエラーとしてUIに返す。

### 3. 不正な config.json を開くと、読み込み失敗後もアクティブな設定パスだけ切り替わる

優先度: 中

GUIの `OpenSettings(path)` は、渡されたパスを先に `a.configPath` へ代入してから `LoadConfig` しています。読み込みに失敗しても `a.configPath` は戻らないため、以後の保存、履歴、診断ログが失敗した設定ファイルの隣へ向く可能性があります。起動引数でJSONを渡した場合も、読み込み失敗時に `runUI(configPath, state)` へ失敗したパスが渡ります。

根拠:

- `src/app.go:249-258`
- `src/main.go:71-82`
- `src/frontend/src/main.js:321-324`

影響:

- 破損したconfigをドロップしただけで、アプリ内部の保存先が意図せず変わる。
- その後の保存や履歴整理で、ユーザーが想定していない場所に `config.json` / `history.json` / `diagnostic.log` が作られる可能性がある。

推奨対応:

- 候補パスをローカル変数で読み込んで成功した場合だけ `a.configPath` と `state.ConfigPath` を更新する。
- 起動引数のJSON読み込み失敗時も、デフォルトconfigPathのままエラー表示するか、明示的に「このconfigは開けない」として保存系操作を無効化する。

### 4. URL自動コピーの失敗を握りつぶし、成功メッセージを出す可能性がある

優先度: 中

`CopySingleURLIfNeeded` の戻り値が複数箇所で `_ =` により無視されています。一方で `resultMessage` はURLがあれば無条件で「画像URLをクリップボードにコピーしました。」を返します。

根拠:

- `src/app.go:281`
- `src/app.go:320`
- `src/app.go:347`
- `src/app.go:373-376`
- `src/main.go:120`

影響:

- クリップボード書き込みが失敗しても、ユーザーはコピー済みだと思ってVRChatへ貼り付けようとする。
- GUIサブシステムや権限、リモートデスクトップ、特殊なクリップボード状態で発生すると原因が分かりづらい。

推奨対応:

- コピー失敗を `state.Message` または `Result.Error` に反映する。
- URL取得自体は成功しているため、エラー扱いではなく「URL取得成功、コピー失敗。サムネイルをクリックして再コピーしてください」のように分ける。

### 5. OSSライセンス表示が依存関係全体を網羅していない

優先度: 中

アプリ内のOSSライセンス表示は手書きの固定リストです。現在は direct dependency だけ見ても `github.com/gofrs/flock`、`github.com/makiuchi-d/gozxing`、`github.com/skip2/go-qrcode` が表示に含まれていません。

根拠:

- `src/app.go:75-85`
- `src/go.mod:5-14`
- `src/frontend/package.json:6-10`

影響:

- 配布アプリのライセンス表示として不完全になる可能性がある。
- 依存追加時に手作業で漏れやすい。

推奨対応:

- 少なくとも direct dependency を全て追加する。
- 中長期的には `go.mod` / `package-lock.json` からライセンス一覧を生成するか、Release成果物へライセンス一覧ファイルを同梱する。

### 6. 自動投稿の監視フォルダ異常やスキャン上限到達がユーザーへ見えない

優先度: 低

自動投稿のスキャンでは、監視フォルダが空、存在しない、アクセス不能、または 5000 件上限に到達した場合でも、UIに明示的な警告が出ません。`filepath.WalkDir` のエラーも無視されます。

根拠:

- `src/internal/appcore/autophoto.go:120-137`

影響:

- ユーザーは自動投稿をONにしたつもりでも、フォルダ設定ミスや権限エラーで投稿されない。
- 画像が多すぎるフォルダを指定した場合、どこまで監視されているか分からない。

推奨対応:

- 保存時または watcher 起動時にフォルダ存在・アクセス可否を検証する。
- スキャン上限到達やWalkDirエラーを診断ログとUI通知に出す。

### 7. 過去のセキュリティレビューを最新報告書として誤読しやすい

優先度: 低

`reports/2026-06-21-security-review.md` は 2026-06-21 時点の報告書として有用ですが、現行コードでは修正済みの内容も高リスクとして残っています。例として、Goバージョン差分の指摘は現在の `go.mod` / CI / Release が `1.26.4` へ揃っており、`govulncheck` でも到達可能な脆弱性は0件です。

根拠:

- `reports/2026-06-21-security-review.md:96-116`
- `src/go.mod:3`
- `.github/workflows/ci.yml:44-48`
- `.github/workflows/release.yml:21-25`

影響:

- 今後の調査で、修正済みの問題と未対応の問題を取り違えやすい。

推奨対応:

- 監査報告書は `reports/` 配下に日付付きファイル名で保存し、作成時点のスナップショットであることを明確にする。
- 現行状態のセキュリティレビューが必要な場合は、新しい日付の報告書を追加する。

### 8. ローカルWindowsビルドのバージョンが最新ローカルタグ固定で、リリース候補確認に向かない

優先度: 低

`scripts/build-windows-from-wsl.sh` はローカルの最新 `v*` タグを自動採用します。今回のように `v0.1.6` タグを削除した状態では、v0.1.6向け修正を確認していても exe の表示バージョンは `v0.1.5.<commit>` になります。

根拠:

- `scripts/build-windows-from-wsl.sh:19-25`
- ローカルタグ一覧は `v0.1.5` まで

影響:

- リリース候補の実機確認で、ユーザーが「直したexeなのか」「旧リリースなのか」を判断しづらい。

推奨対応:

- `VERSION=v0.1.6 ./scripts/build-windows-from-wsl.sh` のような上書き手段を用意する。
- または `--version v0.1.6` 引数を追加する。

### 9. 診断ログの書き込み先が可変パスとして `gosec` に検出される

優先度: 低

`gosec` は `AppendDiagnosticLog` の `os.OpenFile(path, ...)` を `G304` として検出しました。現在の設計では `DiagnosticLogPath(configPath)` によりconfigの隣へ出力されるため、任意コマンド実行のような問題ではありません。ただし、不正なconfig読み込み時に `configPath` が切り替わる問題と組み合わさると、ユーザーが意図しない場所へ `diagnostic.log` を作る可能性があります。

根拠:

- `src/internal/appcore/diagnostic.go:11-27`
- `src/app.go:249-258`

推奨対応:

- finding 3 を先に直す。
- 必要なら診断ログの出力先をアプリ既定ディレクトリ配下へ固定する。

## 修正済みと判断した過去リスク

- Discord Webhook URL は `discord.com` / `discordapp.com` の `/api/webhooks/{id}/{token}` に制限済み。
- 履歴URLの表示・可用性確認はDiscord CDN系URLに制限済み。
- 画像入力はファイルサイズとピクセル数の上限が入っている。
- `config.json` / `history.json` / 出力画像 / 診断ログは `0600` 相当で保存される。
- CI/Release の Go バージョンは `1.26.4` へ揃っている。
- npm監査と `govulncheck` では、現時点で到達可能な既知脆弱性は見つからなかった。

## 推奨する次の対応順

1. Release成果物と README / About / SPEC のPGP説明を統一する。
2. Discord削除の部分成功を履歴に保存する。
3. config読み込み失敗時にアクティブ設定パスを変えない。
4. URL自動コピー失敗をUIに出す。
5. OSSライセンス一覧を依存関係に合わせて更新する。
6. 自動投稿の監視フォルダ異常をUIまたは診断ログに出す。
7. 監査報告書を日付付きファイル名で管理し、最新状態と誤読されないようにする。
