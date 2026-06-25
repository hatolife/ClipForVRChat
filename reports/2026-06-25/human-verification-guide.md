# Human Verification Guide

調査日: 2026-06-25

この文書は、セキュリティ監査後に人間が確認する必要がある作業を整理したものです。

自動テスト、静的解析、依存関係監査で確認できるものはCIまたはローカルコマンドで確認し、Windows実機、Discord実通信、GitHub Release運用、ユーザー向け文言の妥当性など、自動化だけでは判断しにくいものを人間確認の対象にします。

## 確認項目一覧

| No. | 項目 | 人間確認の要否 | 主な確認タイミング |
| --- | --- | --- | --- |
| 1 | `GO-2026-4550` 解消後の互換性 | 一部必要 | 依存関係更新後 |
| 2 | Windows実機でのACL確認 | 必要 | 診断データ作成後 |
| 3 | Windows実機でのGUI/CLI/Explorer確認 | 必要 | exeビルド後 |
| 4 | 診断データzipの中身確認 | issue対応後に必要 | 診断zip安全化対応後 |
| 5 | Webhook token保存方針 | 必要 | 設計判断時 |
| 6 | Discord実通信の動作確認 | 必要 | Discord投稿まわり変更後 |
| 7 | Release workflow権限とEnvironment protection確認 | 必要 | Release workflow変更時 |
| 8 | Release成果物とRelease本文の確認 | 必要 | RC/正式Release後 |
| 9 | ユーザー向け文言の確認 | 必要 | README/About/Release notes更新時 |

## 1. `GO-2026-4550` 解消後の互換性

### 人間がすべきか

依存関係更新と脆弱性解消そのものは、基本的に実装作業と自動検証で確認する項目です。

人間が確認すべきなのは、更新後に以下の実機挙動が壊れていないかです。

- 不具合報告用データ作成が成功する。
- `.zip.gpg` が復号できる。
- 復号後zipが破損していない。
- zipファイル引数暗号化が成功する。

### 自動検証手順

```sh
cd src
go get github.com/cloudflare/circl@latest
go mod tidy
go test ./...
go run golang.org/x/vuln/cmd/govulncheck@latest ./...
```

### 人間確認手順

1. Windows向けexeをビルドする。
2. アプリから不具合報告用データを作成する。
3. 生成された `.zip.gpg` を復号する。
4. 復号後zipを展開または `unzip -t` で検査する。
5. 任意のzipをexe引数へ渡して `.gpg` が生成されることを確認する。

```sh
gpg -d ClipForVRChat-diagnostics-YYYYMMDD-HHMMSS.zip.gpg > decrypted.zip
unzip -t decrypted.zip
```

## 2. Windows実機でのACL確認

### 目的

WindowsではUnixの `0600` と同じ意味でアクセス権が決まるわけではなく、実際の権限はDACLで決まります。そのため、`icacls` で実ファイルのACLを確認します。

### 確認手順

診断データ作成後の作業フォルダで以下を実行します。

```cmd
icacls config.json
icacls history.json
icacls logs
icacls diagnostics
```

### 現時点の確認結果

確認対象:

```text
C:\Users\user\Downloads\ClipForVRChat\ClipForVRChat-diagnostics-20260625-150432
```

`config.json`:

```text
config.json NT AUTHORITY\SYSTEM:(F)
            BUILTIN\Administrators:(F)
            9800X3D\user:(F)

1 個のファイルが正常に処理されました。0 個のファイルを処理できませんでした
```

`logs`:

```text
logs NT AUTHORITY\SYSTEM:(OI)(CI)(F)
     BUILTIN\Administrators:(OI)(CI)(F)
     9800X3D\user:(OI)(CI)(F)

1 個のファイルが正常に処理されました。0 個のファイルを処理できませんでした
```

`history.json`:

```text
history.json: 指定されたファイルが見つかりません。
0 個のファイルが正常に処理されました。1 個のファイルを処理できませんでした
```

`diagnostics`:

```text
diagnostics: 指定されたファイルが見つかりません。
0 個のファイルが正常に処理されました。1 個のファイルを処理できませんでした
```

### 現時点の判断

`config.json` と `logs` は、SYSTEM、Administrators、現在ユーザーに限定されており、確認結果としては妥当です。

`history.json` は、その時点で履歴が存在しなければファイルがないこと自体は異常ではありません。Discord投稿やQRコード検出など、履歴が作成される処理を行った後に再確認します。

`diagnostics` は、確認している場所がすでに診断データ作成先ディレクトリの内側であるため、同階層に `diagnostics` ディレクトリが存在しないのは自然です。ACL確認対象としては、診断データ作成元のアプリ実行フォルダ側にある `diagnostics/`、または作成された時刻付き作業フォルダそのものを確認します。

追加確認例:

```cmd
cd ..
icacls ClipForVRChat-diagnostics-20260625-150432
icacls ClipForVRChat-diagnostics-20260625-150432\config.json
icacls ClipForVRChat-diagnostics-20260625-150432\logs
```

## 3. Windows実機でのGUI/CLI/Explorer確認

### 確認手順

```pwsh
.\ClipForVRChat.exe --version
.\ClipForVRChat.exe --help
```

GUI側では以下を確認します。

- 通常起動できる。
- Explorerでファイル選択表示が動く。
- 不具合報告用データ生成ボタンが動く。
- 生成中に半透明オーバーレイと進捗表示が出る。
- 完了後に `.zip.gpg` がExplorerで選択表示される。
- 完了後に操作不能状態が解除される。

### 現時点の確認結果

```text
C:\Users\user\Downloads\ClipForVRChat\ClipForVRChat-diagnostics-20260625-150432>ClipForVRChat.exe --version

C:\Users\user\Downloads\ClipForVRChat\ClipForVRChat-diagnostics-20260625-150432>ClipForVRChat develop
ClipForVRChat.exe --help

C:\Users\user\Downloads\ClipForVRChat\ClipForVRChat-diagnostics-20260625-150432>ClipForVRChat
Usage: ClipForVRChat [--version]

Options:
  --version              バージョンを表示して終了します
  --help, -h             display this help and exit
```

### 現時点の判断

GUIサブシステムのexeからコンソールへ出力する処理は、PowerShell/cmd上でプロンプト表示と出力順が崩れることがあります。これは過去に確認したWindows GUI exeの親コンソール出力の制約に近い挙動です。

ただし、出力がユーザーに読めること、GUI起動時に不要なコンソールが表示されないことが優先です。CLIを主対象ユーザーが常用しない前提なら、現状は致命的ではありません。

確認時は、次の観点で判断します。

- `--version` の内容が表示される。
- `--help` の内容が表示される。
- 文字化けしない。
- 通常起動時に余計なコンソールが出ない。

## 4. 診断データzipの中身確認

### 確認タイミング

診断データ内のWebhook URL、Discord token、パス置換、zip/gpg内容の安全化に関するissue対応後に確認します。

### 確認手順

1. Webhook URLを設定する。
2. Discord投稿を1回実行し、履歴を作成する。
3. 不具合報告用データを生成する。
4. 確認用zipを展開する。
5. `.zip.gpg` を復号し、復号後zipも展開する。
6. 両方の中身を比較する。

確認コマンド例:

```sh
unzip -l ClipForVRChat-diagnostics-YYYYMMDD-HHMMSS.zip
gpg -d ClipForVRChat-diagnostics-YYYYMMDD-HHMMSS.zip.gpg > decrypted.zip
unzip -t decrypted.zip
unzip -l decrypted.zip
```

### 確認観点

- `config.json`、`history.json`、`logs/` が含まれる。
- 画像本体は含まれず、outputフォルダの内容はログや一覧情報として確認できる。
- `C:\Users\user\...` のような生ユーザーパスが、可能な範囲で `%USERPROFILE%` などに置換されている。
- Webhook URLの生値がない。
- `discordToken` の生値がない。
- 確認用zipと復号後zipで、含まれる情報に意図しない差がない。

## 5. Webhook token保存方針

### 現時点の方針

- Discord投稿を後から削除するため、アプリ内履歴にはtokenを残す。
- 診断データではtokenをダミー化する。
- 診断データ上では、削除用tokenが実際に設定されていたかどうかを判断できなくてもよい。

### DPAPIやWindows Credential Managerへ移す場合の対象

移す候補は、`history.json` に保存しているDiscord削除用tokenです。

Discord投稿を後から削除するには、投稿先Webhookのtokenが必要です。現在は履歴にtokenを保存することで、過去の投稿を履歴画面から削除できます。

DPAPIやWindows Credential Managerへ移す場合は、次のような設計になります。

- `history.json` にはmessage ID、URL、保存日時などの非秘密情報を残す。
- Discord削除用tokenはWindowsの保護ストレージに保存する。
- 履歴レコードには、保護ストレージ上のキーまたは参照IDだけを保存する。
- 削除時に保護ストレージからtokenを取り出してDiscord削除APIを呼ぶ。

### トレードオフ

利点:

- `history.json` だけを読まれてもDiscord tokenが漏れにくい。
- 診断データに履歴を含めやすくなる。

欠点:

- 実装と移行処理が複雑になる。
- Windows専用処理が増える。
- バックアップや別PC移行時に過去投稿の削除ができなくなる可能性がある。
- 保護ストレージの参照が壊れた場合、履歴上はDiscordにあるが削除できない状態が発生する。

現時点では、アプリ内履歴にはtokenを残し、診断データ生成時に必ずダミー化する方針が現実的です。

## 6. Discord実通信の動作確認

### 確認手順

1. Discord投稿ONで画像を処理する。
2. Discordへ投稿されることを確認する。
3. 履歴画面でDiscordありとして表示されることを確認する。
4. Discordから削除できることを確認する。
5. Discord投稿OFFで画像を処理する。
6. Discordへ投稿されないことを確認する。
7. 履歴画面でDiscordなしとして扱われることを確認する。

### 確認観点

- Discord投稿OFF時に外部送信されない。
- ローカル保存OFF時にoutputへ画像が残らない。
- QRコード読み取りだけONの場合、QRがある時だけ結果と履歴が出る。
- QRコードがない場合、結果なしの理由がユーザーに分かる。

## 7. Release workflow権限とEnvironment protection確認

### 確認手順

GitHubのRepository Settingsとworkflowファイルを確認します。

確認観点:

- build/test jobは `contents: read` になっている。
- Release作成jobだけ `contents: write` を持つ。
- Release用secretが不要なjobに渡っていない。
- `environment: release` などを使う場合、required reviewersが設定されている。
- tag pushだけで意図しないReleaseが公開されない。

### Environment protectionの意味

Environment protectionは、GitHub Actionsの特定jobが実行される前に、人間の承認やsecret公開範囲の制限を挟む仕組みです。

Release jobに設定すると、誤ったtag pushや想定外のworkflow変更が即座にRelease成果物作成へつながるリスクを下げられます。

## 8. Release成果物とRelease本文の確認

### 確認手順

RCまたは正式Release後にGitHub Release画面と成果物zipを確認します。

確認観点:

- Release本文が `RELEASE_NOTES.md` の対象バージョン由来である。
- Release assetsにzip、sha256、exe.ascがある。
- zip内に `ClipForVRChat.exe`、`README.md`、`LICENSE`、`Release-signing-public-key.url` がある。
- zip内に不要な `.asc` や公開鍵実体ファイルが混入していない。
- sha256が一致する。
- PGP署名検証が通る。

## 9. ユーザー向け文言の確認

### 確認観点

- 添付すべきファイルが `.zip.gpg` だと分かる。
- 平文zipは確認用だと分かる。
- 誤ってzipを添付すると中身が見えると分かる。
- ただし秘密情報は可能な範囲で除外またはダミー化されると分かる。
- TwitterとGitHub Issueへの案内が自然である。

## 優先順位

直近の優先度は以下です。

1. 診断zipのWebhook URL / Discord tokenダミー化対応後の中身確認。
2. `GO-2026-4550` 解消後の暗号化/復号実機確認。
3. Windows実機ACLの追加確認。
4. Release workflow権限とEnvironment protectionの確認。
5. Discord実通信の回帰確認。
