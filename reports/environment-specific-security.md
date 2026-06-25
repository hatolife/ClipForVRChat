# Environment Specific Security

調査日: 2026-06-25

## 対象環境

- Windows GUIデスクトップアプリ
- Wails/WebView2アプリ
- 補助CLI
- GitHub ActionsによるWindows Release build

## デスクトップアプリ

- 確認結果: サーバーとして待ち受ける機能は確認できない。
- 良い点: 管理者権限を前提にしない。起動ロックがあり、履歴・設定の多重更新リスクを下げている。
- 注意点: Wails frontendからGo methodを呼べるため、将来HTML injectionや外部コンテンツ読み込みを追加した場合の影響が大きい。

## GUIアプリ

- 確認結果: ドラッグ&ドロップ、ボタン操作、設定入力が主な入口。
- 良い点: Vueテンプレート中心で、`innerHTML` や `v-html` は確認できない。
- 注意点: `OpenURL` が任意URLを受け入れる。現在は定数URLとGitHub API由来のURLに限定されるが、許可リスト化が望ましい。

## CLIツール

- 確認結果: `--version`、`--help`、`.zip` 引数暗号化を持つ。
- 良い点: `.zip` 拡張子とzipとして読めることを確認してから暗号化する。
- 注意点: 任意ローカルzipを読み込み、同じ場所へ `.gpg` を作成する。明示的なユーザー操作なので許容範囲だが、上書き確認はない。

## ファイルシステム

- 確認結果: config/history/log/diagnostics/outputをアプリ実行フォルダ基準で扱う。
- 良い点: Go側で `0600` / `0700` 相当を指定。画像入力サイズとピクセル数上限がある。
- 注意点: Windows ACLとして期待通り制限されるか要確認。履歴内の絶対 `outputPath` は削除対象として解決される。

## パス処理

- 確認結果: 設定の出力先や監視フォルダは前後の空白と引用符を除去する。
- 良い点: ローカル保存ファイル名は `filepath.Base(sourcePath)` から作成され、入力パスのディレクトリ成分を保存名に使わない。
- 注意点: `history.json` 改ざん時の絶対パス削除リスクがある。診断ログでは作成時にパス置換するが、ログファイル自体は生パスを持ち得る。

## 一時ファイル

- 確認結果: テスト用の診断zip作成では `os.MkdirTemp` を使う。実アプリの診断データは `diagnostics/<timestamp>/` に残す。
- 注意点: 確認用zipを残す仕様により、秘密情報の平文残存リスクがある。

## シンボリックリンク・ショートカット

- 確認結果: 明示的なsymlink対策は確認できない。
- リスク: 出力先や履歴パスがsymlink経由の場合、ユーザー権限範囲で想定外の場所へ書き込み・削除する可能性がある。
- 推奨: 管理ディレクトリ配下判定と、削除対象の実体確認を追加する。

## 権限境界

- 確認結果: 管理者権限は不要。同一ユーザー権限で完結する。
- 注意点: 同一ユーザー権限で読めるファイルにはWebhook URLや履歴tokenが含まれる。

## ネットワーク接続

- Discord投稿・削除: timeoutあり。Webhook URLは `https`、Discordホスト、`/api/webhooks/{id}/{token}` 形式を検証。
- GitHub Release確認: timeoutあり。自動ダウンロード・自動更新はなし。
- Discord添付URL確認: trusted hostとpath確認あり。

## IPC

- Windows clipboardを使用。
- Wails runtime eventsでfrontendへ進捗通知。
- Windows Shell APIでExplorer選択表示。

## プロセス起動

- 明示的なshell command実行は確認できない。
- BrowserOpenURLとShell APIはOS機能経由で外部アプリを起動し得る。

## 動的ライブラリ読み込み

- `syscall.NewLazyDLL` でWindows標準DLLを参照。
- 任意DLLやプラグイン読み込みは確認できない。

## 自動更新

- 自動更新は対象外。
- 更新通知だけを行い、ユーザーがGitHub/BOOTHを開く。

## インストーラー・アンインストーラー

- 対象外。zip配布でありインストーラーは確認できない。

## 常駐プロセス

- アプリ起動中のみVRChat写真/Screenshotsフォルダを定期スキャンする。
- スキャン件数上限と1tick処理件数上限がある。

## Release成果物

- zip、sha256、exe detached PGP署名を作成。
- zip内に公開鍵URL、README、LICENSE、exeを含める。
- SBOMは未作成。
