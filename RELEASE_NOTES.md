# Release Notes

## v0.1.8

### 更新内容

- 設定画面に「自動撮影」タブを追加し、OSC、撮影間隔、撮影方式、出力、Presence、Discord投稿設定をまとめました。
- VRChat User CameraへOSCで構図を送り、Stream Camera(Spout)方式またはPhoto方式で有効な構図を順番に撮影する自動撮影機能を追加しました。
- Stream方式では同梱の `spout-capture.exe` がVRChat Stream CameraのSpout senderから1フレームをPNGとして受信し、必要に応じてJPGへ変換して保存します。
- VRChat output logから同じインスタンスにいるユーザー情報、world ID、instance IDを推定し、撮影画像に対応するsidecar JSONへ保存するようにしました。
- 自動撮影画像へPNG iTXtまたはJPEG EXIF APP1で撮影メタデータを埋め込めるようにしました。ユーザーID埋め込みは設定で独立して制御できます。
- 自動撮影した画像を既存の結果/履歴画面で扱えるようにし、設定で有効化した場合はDiscord Webhookへ投稿できるようにしました。画像添付なしの本文のみ投稿にも対応しました。
- 構図カード内に「現在Poseから追加」と「このPoseへカメラ移動」を追加し、設定済みPoseをゲーム内カメラへ送れるようにしました。
- User Camera関連OSCをfalse/Offへ戻す「カメラOSCをリセット」ボタンを自動撮影タブに追加しました。

### 既知の制限

- player_local構図は標準OSCだけでプレイヤーrootを自動取得できないため、手動で保存したプレイヤー基準Poseを使います。
- output log由来のユーザー一覧やworld/instance情報は、VRChatログの内容によって取得できない場合があります。
- Camera Dolly Multi、解像度一時変更、SQLiteローカルDB、OSCQuery自動検出はv0.1.8の対象外です。v0.1.8ではsidecar JSONと履歴JSONを正本/索引として扱います。

### ダウンロード

- プログラムのダウンロード: [ClipForVRChat-v0.1.8-windows-amd64.zip](https://github.com/hatolife/ClipForVRChat/releases/download/v0.1.8/ClipForVRChat-v0.1.8-windows-amd64.zip)
- 署名確認用ファイル: [ClipForVRChat-v0.1.8-windows-amd64.exe.asc](https://github.com/hatolife/ClipForVRChat/releases/download/v0.1.8/ClipForVRChat-v0.1.8-windows-amd64.exe.asc)
- 署名確認用公開鍵: https://keys.openpgp.org/search?q=release-signing@hato.life

### 比較

https://github.com/hatolife/ClipForVRChat/compare/v0.1.7...v0.1.8

## v0.1.7

### 更新内容

- 設定画面のカテゴリを整理し、Discord投稿関連設定を「Discord投稿」タブへまとめました。
- 初期設定でDiscord投稿と投稿URLの自動コピーをOFFにしました。
- Discord投稿がONで通常投稿用Webhook URLが空欄の場合、画面上部の警告と入力欄の注意表示で設定漏れに気づきやすくしました。
- 開発ビルドのバージョン表記にコミットIDと `develop` を含めるようにしました。
- 不具合報告用データと診断ログから、Webhook URLやDiscord tokenなどの秘密情報が残りにくいよう改善しました。
- セキュリティ監査結果を受け、依存関係更新、外部URL制限、履歴ローカル削除範囲の制限、Release workflow権限の最小化、ビルドメタデータ添付を行いました。

### ダウンロード

- プログラムのダウンロード: [ClipForVRChat-v0.1.7-windows-amd64.zip](https://github.com/hatolife/ClipForVRChat/releases/download/v0.1.7/ClipForVRChat-v0.1.7-windows-amd64.zip)
- 署名確認用ファイル: [ClipForVRChat-v0.1.7-windows-amd64.exe.asc](https://github.com/hatolife/ClipForVRChat/releases/download/v0.1.7/ClipForVRChat-v0.1.7-windows-amd64.exe.asc)
- 署名確認用公開鍵: https://keys.openpgp.org/search?q=release-signing@hato.life

### 比較

https://github.com/hatolife/ClipForVRChat/compare/v0.1.6...v0.1.7

## v0.1.6

### 更新内容

- 初回起動時に、設定を保存するまで `config.json` を作成しないよう修正しました。
- Windowsのファイルプロパティにバージョン情報・製品情報を追加しました。
- セキュリティチェック (`gosec`) 対応と品質改善を行いました。

### ダウンロード

- プログラムのダウンロード: https://github.com/hatolife/ClipForVRChat/releases/download/v0.1.6/ClipForVRChat-v0.1.6-windows-amd64.zip
- 署名確認用ファイル: https://github.com/hatolife/ClipForVRChat/releases/download/v0.1.6/ClipForVRChat-v0.1.6-windows-amd64.exe.asc
- 署名確認用公開鍵: https://keys.openpgp.org/search?q=release-signing@hato.life

### 比較

https://github.com/hatolife/ClipForVRChat/compare/v0.1.5...v0.1.6

## v0.1.5

### 更新内容

- Win+Shift+SなどでScreenshotsフォルダに保存された画像を、自動でDiscordへ投稿する機能を追加しました。
- スクリーンショット自動投稿を設定画面でON/OFFできるようにしました。初期値はOFFです。
- VRChat写真自動投稿とスクリーンショット自動投稿のスキャン間隔を、それぞれ設定画面で変更できるようにしました。
- スクリーンショット自動投稿用Webhook URLを設定できるようにしました。空の場合は通常のDiscord Webhook URLへ投稿します。

### ダウンロード

- プログラムのダウンロード: https://github.com/hatolife/ClipForVRChat/releases/download/v0.1.5/ClipForVRChat-v0.1.5-windows-amd64.zip
- 署名確認用ファイル: https://github.com/hatolife/ClipForVRChat/releases/download/v0.1.5/ClipForVRChat-v0.1.5-windows-amd64.exe.asc
- 署名確認用公開鍵: https://keys.openpgp.org/search?q=release-signing@hato.life

### 比較

https://github.com/hatolife/ClipForVRChat/compare/v0.1.4...v0.1.5

## v0.1.4

### 更新内容

- QRコード読み取り機能を追加しました。Discord投稿本文と結果画面に表示します。
- 新しいバージョンの通知機能を追加しました。起動時にチェックします。
- 公式配布ファイルの改竄確認ができるように、PGP署名ファイルと公開鍵を追加しました。
- 情報画面とREADMEに、公式配布場所、PGP確認方法、連絡・要望の案内を追加しました。

### ダウンロード

- プログラムのダウンロード: https://github.com/hatolife/ClipForVRChat/releases/download/v0.1.4/ClipForVRChat-v0.1.4-windows-amd64.zip
- 署名確認用ファイル: https://github.com/hatolife/ClipForVRChat/releases/download/v0.1.4/ClipForVRChat-v0.1.4-windows-amd64.exe.asc
- 署名確認用公開鍵: https://keys.openpgp.org/search?q=release-signing@hato.life

### 比較

https://github.com/hatolife/ClipForVRChat/compare/v0.1.3...v0.1.4
