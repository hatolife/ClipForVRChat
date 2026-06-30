# Release Notes

## v0.1.8

### 更新内容

- 設定画面に「自動撮影」タブを追加し、OSC、撮影間隔、撮影方式、出力、Presence、Discord投稿設定をまとめました。
- VRChat User CameraへOSCで構図を送り、Stream(ffmpeg)方式またはPhoto方式で有効な初期構図を順番に撮影する自動撮影MVPを追加しました。
- VRChat output logから同じインスタンスにいるユーザー情報を推定し、撮影画像に対応するsidecar JSONへ保存するようにしました。
- 自動撮影した画像を既存の結果/履歴画面で扱えるようにし、設定で有効化した場合はDiscord Webhookへ投稿できるようにしました。
- RC確認で見つかった、VRChat写真が保存されない問題に対して、Capture/CloseをUser CameraのAction OSCとして送るよう修正し、Stream方式ではffmpegで映像から静止画を切り出せるようにしました。
- Stream方式ではVRChat OSCのStream Cameraモードを開くよう修正し、設定画面からffmpegの確認と `winget install ffmpeg` による導入を行えるようにしました。
- Stream方式のffmpeg初期入力をVRChatウィンドウ範囲の切り出しに変更し、旧RCのデスクトップ全体取得やtitle直接取得の設定は自動移行するようにしました。
- 構図カード内に「現在Poseから追加」と「このPoseへカメラ移動」を追加し、設定済みPoseをゲーム内カメラへ送れるようにしました。
- User Camera関連OSCをfalse/Offへ戻す「カメラOSCをリセット」ボタンを自動撮影タブに追加しました。

### 既知の制限

- v0.1.8のStream方式は外部ffmpegを使います。Spout2直接受信、Camera Dolly Multi、解像度一時変更は将来対応です。
- output log由来のユーザー一覧は、アプリ起動前から同じインスタンスにいたユーザーを完全に復元できない場合があります。

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
