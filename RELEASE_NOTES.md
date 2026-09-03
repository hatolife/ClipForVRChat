# Release Notes

## v0.1.8

## ダウンロード

- [プログラムのダウンロード](https://github.com/hatolife/ClipForVRChat/releases/download/v0.1.8/ClipForVRChat-v0.1.8-windows-amd64.zip)

### 更新内容

- VRChat User CameraをOSCで操作し、複数構図をStream Camera(Spout)方式またはPhoto方式で順番に撮影する自動撮影機能を追加しました。
- 自動撮影画像にsidecar JSONとPNG/JPEGメタデータを保存し、VRChat output logから推定した同席ユーザー、world ID、instance IDも記録できるようにしました。
- 自動撮影結果を既存の結果/履歴画面、ローカル保存、Discord投稿へ統合しました。専用Webhook URLが空欄の場合は通常投稿用Webhook URLへフォールバックします。
- AvatarBeaconからOSC Avatar Parametersで受け取る `avatar_osc` basisを使い、プレイヤー基準の構図をアバター位置と向きに追従できるようにしました。
- 各 `player_local` 構図のCamera Pose送信直前に最新のAvatarBeacon basisを再取得し、撮影開始後の移動にも追従しやすくしました。
- 自動撮影完了時に設定画面を開いている場合は結果画面へ切り替えず、編集中の未保存設定を保持するようにしました。
- 自動撮影終了時はCamera LockとFlyingをOFFにし、待機カメラ位置への移動は任意かつ既定OFFにしました。
- 撮影、Photoボタン、Spout sender復旧、Camera Mode復元などの待機時間をミリ秒単位で調整できるようにしました。
- AvatarBeaconの `AvatarBeacon_v0.0.1.unitypackage` を通常利用者向けzipと分離版zipへ同梱するようにしました。
- 設定画面に「OSC」タブを追加し、OSC送受信ポート、カメラOSCリセット、AvatarBeacon受信状態、player_local basis取得元をまとめました。
- 設定画面を左側の縦ナビゲーションと右側の設定ペインへ変更し、「機能 ON/OFF」「縮小処理」「投稿処理」「撮影処理」「OSC」「その他」に整理しました。
- 「自動撮影スケジュール」を「自動定期撮影」へ変更し、主スイッチと間隔設定を「機能 ON/OFF」へ移動しました。
- VRChatから受信したOSC packetを別UDPポートへ転送できるようにし、他OSC受信アプリとのポート競合を避けやすくしました。
- Stream Camera(Spout)取得の起動待ち、空フレーム待ち、PNG書き出し、失敗画像の隔離を改善し、Stream方式の撮影失敗を減らしました。
- 自動撮影後にUser CameraのMode、Pose、Streaming、Zoom、Exposure、mask類などをできるだけ撮影前の状態へ戻すようにしました。
- 自動撮影の通常モードで、撮影後に送る待機カメラ位置を設定できるようにしました。既定ではOFFで、ONにした場合の初期位置はプレイヤー基準の足元方向です。
- 待機カメラ位置の詳細設定をフォールバックモード中でも開けるようにし、構図ごとの明るさと表示対象マスクを設定できるようにしました。既定ではVRChatの初期明るさに合わせ、他ユーザーも写る設定にしました。
- カメラ構図へ移動するときの飛行モードを、Pose送信直前だけONにして移動後すぐOFFへ戻すようにしました。
- 構図の「現在Poseを取得」で、直近のUser Camera Zoomも拡大率として保存するようにしました。
- ローカルアンカー配置済みカメラを使うフォールバックモードを追加し、自動ON/OFFを個別に設定できるようにしました。自動ON/OFFはどちらも既定OFFです。
- 自動撮影の詳細設定画面でも「保存」「閉じる」ボタンを表示するようにしました。
- Release workflowで同梱するAvatarBeacon packageのPGP署名を検証し、Spout2取得元をSHA256検証付きの固定archiveにしました。
- Release workflowの署名処理をbuild jobから分離し、署名secretを署名専用jobだけで扱うようにしました。
- 単一exe配布では既定で埋め込みSpout helperを使い、同じフォルダに置かれた未署名helperを優先しないようにしました。
- 設定由来のlegacy ffmpeg実行パスを制限し、外部設定から任意実行ファイルを起動しにくくしました。
- 自動処理の監視フォルダとDiscord送信先を保存前に確認できるようにし、広すぎる監視フォルダには警告を表示するようにしました。
- Webhook、自動投稿、監視フォルダ、自動撮影スケジュールなどの重要な設定差分を保存前に具体的に確認できるようにしました。
- 自動撮影のDiscord投稿を専用設定がONの場合だけ実行するようにし、通常Discord投稿ONだけでは自動撮影結果を投稿しないようにしました。
- Discord履歴にWebhook tokenを保存せず、Discord側の削除時だけ現在設定から一時的にtokenを解決するようにしました。
- 別フォルダや別バージョンのClipForVRChatでも同じWindowsユーザーでは単一起動になるようにし、OSC port競合を避けました。
- 起動中の進捗表示、frontend template検査、Wails API surface検査を追加し、GUI起動不具合を切り分けやすくしました。

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
