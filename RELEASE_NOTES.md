# Release Notes

## v0.1.8

### 更新内容

- 設定画面に「自動撮影」タブを追加し、OSC、撮影間隔、撮影方式、出力、Presence、Discord投稿設定をまとめました。
- VRChat User CameraへOSCで構図を送り、Stream Camera(Spout)方式またはPhoto方式で有効な構図を順番に撮影する自動撮影機能を追加しました。
- Stream方式では内蔵の `spout-capture.exe` がVRChat Stream CameraのSpout senderから1フレームをPNGとして受信し、必要に応じてJPGへ変換して保存します。通常版exeには `spout-capture.exe` と `SpoutLibrary.dll` を埋め込み、初回使用時に管理フォルダへ展開して呼び出します。
- `spout-capture.exe` はSpout受信だけを担当するWindows helperです。Spout/DirectX/OpenGL/DLLロードをClipForVRChat本体プロセスから隔離するために別プロセスとして実行しており、sender列挙、1フレーム受信、指定先へのPNG保存、結果JSON出力だけを行います。ネットワーク送信やWebhook URL/設定ファイルの読み取りは行いません。
- VRChat output logから同じインスタンスにいるユーザー情報、world ID、instance IDを推定し、撮影画像に対応するsidecar JSONへ保存するようにしました。
- 自動撮影画像へPNG iTXt/eXIfまたはJPEG EXIF APP1で撮影メタデータを埋め込めるようにしました。ユーザーID埋め込みは設定で独立して制御できます。
- 自動撮影した画像を既存の結果/履歴画面で扱えるようにし、設定で有効化した場合はDiscord Webhookへ投稿できるようにしました。画像添付なしの本文のみ投稿にも対応しました。
- 埋め込みメタデータ書き込みに失敗した場合も、sidecar JSON、Discord投稿、履歴追加は可能な限り継続し、警告として記録します。
- 構図カード内に「現在Poseから追加」と「このPoseへカメラ移動」を追加し、設定済みPoseをゲーム内カメラへ送れるようにしました。
- User Camera関連OSCをfalse/Offへ戻す「カメラOSCをリセット」ボタンを自動撮影タブに追加しました。
- 専用アバターギミックからOSC Avatar Parametersで送られるPoseを `player_local` のbasisに使う `avatar_osc` 取得元を追加しました。標準OSC単体では動作せず、専用アバターギミック導入済みアバターでの実機確認が必要です。
- AvatarBeacon を `avatar_osc` basis 確認用の汎用アバターギミックとして整理しました。OSC parameterは `coord/*` と `forward/*` を既定にし、CI は source zip だけを作成します。Unity で作る `.unitypackage` はリリース担当者が手作業で作成して GitHub Release に手動添付します。source zip は `Assets/PoppoWorks/AvatarBeacon/...` に展開できる形です。
- AvatarBeaconのbasisを、positionはHips基準、yaw/forwardはHead基準で送る構成にしました。プレイヤー位置はHips寄り、向きは顔の向き寄りとして扱います。
- プレイヤー基準の取得元の初期値を `avatar_osc` に変更しました。新規設定、または `basisSource` が未保存の設定では、AvatarBeaconからのOSC basis受信を既定で使います。
- VRChat写真、スクリーンショット、自動撮影の専用Webhook URLが空欄の場合に通常投稿用Webhook URLへフォールバックすることを設定画面で明示し、通常投稿用Webhook URLが設定済みなら専用Webhook未設定だけでは保存時確認を出さないようにしました。
- `avatar_osc` 受信状態で、Avatar OSC未受信時にmanual基準Pose未設定エラーが表示されないようにし、AvatarBeaconの `coord/*` / `forward/*` 受信確認先を表示するようにしました。
- `avatar_osc` 受信処理がVRChatのAvatar Parameter受信ごとに大量ログを出し、起動直後のGUI表示を阻害する場合がある問題を修正しました。
- 自動撮影タブの `avatar_osc` 受信状態を自動更新にし、手動更新ボタンを削除しました。最終受信、position、yawを確認できます。
- `avatar_osc` 受信状態のyaw表示がGo APIの `pose.rotation.y` を読んでおらず0度に見える問題を修正しました。yawはAvatarBeaconの `forward/*` から復元した現在の向きです。
- GUI起動の切り分け用に、Wails起動/DOM準備/frontend初期化/手動終了/shutdownの診断ログと、起動中の簡易進捗表示を追加しました。
- 自動撮影タブの説明文がfrontendのtemplate literalを壊し、起動時に `avatar is not defined` が出る問題を修正しました。同種の混入をCI/Releaseで検出する検査も追加しました。
- `avatar_osc` が `stale` の場合に、内部エラー名ではなく、AvatarBeaconの `coord/*` / `forward/*` 更新停止、最後に受信したparameter、OSC config resetやContact確認へ誘導する診断文を表示するようにしました。
- AvatarBeacon切り分け用のOSC診断ログを、全packet連続出力から10秒ごとのsummaryへ変更しました。status、raw、last、position/yaw、`coord/*` / `forward/*` の最新値と受信時刻を確認できます。
- 設定画面を開いたときに同じOSC受信ポートへ再bindして失敗し、`avatar_osc` 受信器が止まって `stale` になる問題を修正しました。同じhost/port/log pathでは既存の受信器を維持します。
- `avatar_osc` の鮮度判定を、basis対象parameterの最古受信時刻ではなく最新受信時刻で見るようにしました。VRChat OSCが値変化時だけ送る挙動でも、変化していない軸だけを理由にstaleになりにくくなります。
- 自動撮影タブのAvatarBeacon受信状態から、ログsummary確認を促す常時説明文を削除しました。
- Spout helperのPNG書き出しでWindows WIC PNG encoderがRGBA pixel formatを受け付けず、Stream方式テスト撮影が `PNG encoder does not support RGBA` で失敗する問題を修正しました。
- 自動撮影後にUser Cameraの状態をできるだけ元へ戻す復元処理を追加しました。撮影前に受信したUser Camera OSCのMode、Pose、Streaming、SmoothMovement、Zoom、Exposure、mask類などを優先し、受信できない項目は設定画面末尾付近のフォールバック値を使います。撮影失敗やキャンセルでも、OSC接続が開けた後なら復元処理を試みます。
- Stream方式でVRChatへStream Camera/Spout ONを送った直後、Spout sender生成が間に合わず `Spout senderがありません` で即失敗する問題を修正しました。Spout helperはtimeout内でsender出現を待ち、自動撮影側も各shotのSpout取得直前にStream起動OSCを再送します。
- Stream方式でSpout起動直後の空フレームを透明PNGとして保存し、`取得画像がほぼ透明です` で失敗する問題を修正しました。Spout helperは非空の有効フレームまでtimeout内で待ち、アプリ側は検証成功後だけ最終ファイル名へ確定するため、失敗画像が自動処理やDiscord投稿へ流れにくくなりました。

### 既知の制限

- player_local構図は標準OSCだけでプレイヤーrootを自動取得できないため、AvatarBeacon導入済みアバターからの `avatar_osc` basisを基本にします。専用ギミックなしの場合は、手動で保存したmanual基準Poseをフォールバックとして使います。
- `avatar_osc` basisはpositionがHips基準、yawがHead基準であり、player root基準そのものではありません。専用アバターギミック未導入、OSC無効、parameter欠落、鮮度切れの場合は撮影前に失敗します。
- AvatarBeacon の `.unitypackage` は CI では作らず、手作業で作成して Release に添付します。source zip は再生成・検証用の配布物です。
- output log由来のユーザー一覧やworld/instance情報は、VRChatログの内容によって取得できない場合があります。
- Camera Dolly Multi、解像度一時変更、SQLiteローカルDB、OSCQuery自動検出はv0.1.8の対象外です。v0.1.8ではsidecar JSONと履歴JSONを正本/索引として扱います。

### ダウンロード

- プログラムのダウンロード: [ClipForVRChat-v0.1.8-windows-amd64.exe](https://github.com/hatolife/ClipForVRChat/releases/download/v0.1.8/ClipForVRChat-v0.1.8-windows-amd64.exe)
- 署名確認用ファイル: [ClipForVRChat-v0.1.8-windows-amd64.exe.asc](https://github.com/hatolife/ClipForVRChat/releases/download/v0.1.8/ClipForVRChat-v0.1.8-windows-amd64.exe.asc)
- 破損確認用ファイル: [ClipForVRChat-v0.1.8-windows-amd64.exe.sha256](https://github.com/hatolife/ClipForVRChat/releases/download/v0.1.8/ClipForVRChat-v0.1.8-windows-amd64.exe.sha256)
- 検証・切り分け用分離版zip: [ClipForVRChat-v0.1.8-windows-amd64-separated.zip](https://github.com/hatolife/ClipForVRChat/releases/download/v0.1.8/ClipForVRChat-v0.1.8-windows-amd64-separated.zip)
- 署名確認用公開鍵: https://keys.openpgp.org/search?q=release-signing@hato.life
- 署名確認用fingerprint: `BE40 AA8D 082F 493F 613B C072 21DC 3486 1B40 E77D`

通常版exeにはStream Camera(Spout)方式用の `spout-capture.exe` と `SpoutLibrary.dll` を埋め込みます。分離版zipはhelper単体確認や不具合切り分け用です。

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
