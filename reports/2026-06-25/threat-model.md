# Threat Model

調査日: 2026-06-25

## 想定ユーザー

- VRChatで外部画像URLを使いたい一般Windowsユーザー。
- Discord Webhookを作成できるユーザー。
- GitHub ReleaseまたはBOOTHからアプリを入手するユーザー。

## 想定攻撃者

- ユーザーに細工した画像やzipを処理させる第三者。
- 公開IssueやSNSに添付された診断データを収集する第三者。
- 同一PC上でユーザー権限のファイルを読めるローカル攻撃者。
- Discord Webhook URLを入手した第三者。
- Release workflowや依存関係を狙うサプライチェーン攻撃者。

## 保護すべき資産

- Discord Webhook URL
- Discord message ID、Webhook ID/token
- `config.json`
- `history.json`
- `logs/YYYY-MM-DD.log`
- 診断データzipと `.zip.gpg`
- ユーザーの画像、画像名、保存パス
- QRコードから読み取ったURL
- Release署名鍵とGitHub Release権限
- 配布exeの完全性

## 信頼境界

- ユーザーが選択・ドロップしたファイルとアプリ内部処理の境界
- `config.json` / `history.json` と実行中アプリの境界
- Wails frontendとGo backendの境界
- Discord/GitHub APIレスポンスとアプリ内部状態の境界
- ユーザー確認用の安全化済みzipと暗号化済み診断データの境界
- GitHub Actions runnerとRelease成果物の境界

## 入力経路

- CLI引数
- GUI操作
- ドラッグ&ドロップ
- クリップボード
- 設定ファイル
- 履歴ファイル
- 画像ファイル
- 自動監視フォルダ
- Discord APIレスポンス
- GitHub Releases APIレスポンス

## 出力経路

- Discord Webhook投稿
- Discord削除API
- ローカル画像保存
- config/history/log/diagnostics書き込み
- クリップボード
- ブラウザ起動
- Explorer表示
- Release assets

## 権限境界

- アプリは通常ユーザー権限で動作する。
- 管理者権限は不要。
- ローカルファイル操作は同一ユーザー権限の範囲。
- Discord Webhook URLは外部サービス上の投稿・削除権限として扱う必要がある。
- GitHub Actions Release jobは `contents: write` を持つ。

## 攻撃面

- 画像デコーダへの不正画像入力
- 監視フォルダへの大量ファイル投入
- Webhook URLの漏えい
- 診断zipの誤添付
- 履歴ファイル改ざんによるローカル削除誘導
- 任意URLオープン導線の将来拡張
- Release workflowのAction/依存関係
- OpenPGP暗号処理依存関係

## STRIDE

### Spoofing

- Discord Webhook URLが漏れると第三者が投稿者のように投稿できる。
- GitHub Release以外の配布物を公式と誤認する可能性がある。PGP署名手順で軽減。

### Tampering

- `config.json` や `history.json` の改ざんで保存先、監視先、削除対象を変更できる。
- Release workflowや依存Actionの改ざんは配布物改ざんにつながる。

### Repudiation

- ログはユーザー操作を記録するが、改ざん耐性や署名はない。
- Discord削除操作の監査証跡はローカル履歴に依存する。

### Information Disclosure

- config/history/log/diagnosticsにWebhook URL、token、パス、QR URL、利用状況が含まれ得る。
- 確認用zipに秘密情報が残っている場合、その誤添付が最大の情報漏えい経路になる。確認用zip自体はユーザー確認用として必要なため、Webhook URLやtokenを入れない設計にする。

### Denial of Service

- 大量画像や巨大画像による処理負荷。ファイル数・サイズ・ピクセル数制限で一部軽減。
- Discord/GitHub API障害による機能低下。

### Elevation of Privilege

- 現状、管理者権限へ昇格する経路は確認できない。
- Windows DLL search orderや外部プロセス起動は限定的だが、Wails/WebView2依存の挙動は継続確認が必要。

## 想定される攻撃シナリオ

1. ユーザーが不具合報告用データの確認用zipをGitHub Issueへ添付したとき、zip内にWebhook URLや履歴tokenが残っていると漏えいする。
2. ローカル攻撃者が `history.json` の `outputPath` を改ざんし、ユーザーに履歴画面から削除操作をさせる。
3. 細工画像により画像デコーダやQRライブラリの脆弱性を突く。
4. Release workflowやAction更新経路を侵害し、成果物を改ざんする。
5. 将来のUI変更でQR URLを直接開くようになり、フィッシングURLへ誘導される。

## 優先的に対策すべきリスク

1. 依存関係既知脆弱性の解消。
2. 診断データの確認用zipから秘密情報を除外し、Webhook URLやtokenをダミー化する。
3. Discord tokenの保存最小化またはOS保護。
4. Release workflow権限・Action pinningの硬化。
5. 任意URLオープンの許可リスト化。
