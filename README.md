# ClipForVRChat

ClipForVRChat は、VRChat で外部画像URLを使いやすくするための Windows アプリです。

画像を VRChat で扱いやすい 2048x2048px 以内へ縮小し、必要に応じて Discord に投稿して、画像URLをクリップボードへコピーします。

## できること

- 画像を縦横 2048px 以内に縮小
- 複数画像をまとめて処理
- PNG または JPG で保存
- 縮小画像をPCに保存
- Discord Webhook に画像を投稿
- 投稿された画像URLをクリップボードへコピー
- 画像内のQRコードURLを検出してDiscord本文と結果画面に表示
- 取得した画像URLを履歴として保存
- 履歴画面からDiscord上の投稿画像を削除
- VRChatで撮影された写真を自動で検知してDiscordへ投稿
- Win+Shift+Sで保存されたスクリーンショットを自動で検知してDiscordへ投稿
- VRChat User CameraをOSCで操作し、Stream Camera(Spout)方式で自動撮影
- VRChatから受信したOSCを他アプリ向けの別UDPポートへ転送

## 使い方

### 1. ダウンロード

最新リリースから `ClipForVRChat-vX.Y.Z-windows-amd64.zip` をダウンロードします。

https://github.com/hatolife/ClipForVRChat/releases/latest

zipを展開し、フォルダ内の `ClipForVRChat.exe` を起動してください。

同じリリースには、zip本体とは別に `ClipForVRChat-vX.Y.Z-windows-amd64.exe.asc` も添付されています。`.exe.asc` は通常利用者向けzip内の `ClipForVRChat.exe` をPGP署名検証するためのファイルです。通常はzipだけダウンロードすれば使用できます。

PGP署名を確認する場合は、zipを展開し、`ClipForVRChat.exe` と同じフォルダに `.exe.asc` を置いてください。公開鍵は `release-signing@hato.life` の鍵を使用し、取り込んだ鍵のfingerprintが `BE40 AA8D 082F 493F 613B C072 21DC 3486 1B40 E77D` と一致することを、このREADMEや公式配布ページなどRelease assetとは別の信頼経路で確認してください。同じReleaseに同梱されたURLや公開鍵だけでは、公開鍵自体の真正性は確認できません。信頼済みfingerprintの公開鍵で `gpg --verify ClipForVRChat-vX.Y.Z-windows-amd64.exe.asc ClipForVRChat.exe` が成功した場合に限り、exeと署名の組み合わせがその鍵で署名されたものだと判断できます。

### 2. 初回設定

初めて起動したときは設定画面が開きます。

まずはそのまま保存しても使用できます。Discordへ投稿したい場合は、Discord Webhook URL を設定してください。

設定は `ClipForVRChat.exe` と同じ場所の `config.json` に保存されます。通常はこのファイルを直接編集する必要はありません。

### 3. 画像を処理する

アプリ画面に画像ファイルをドラッグ&ドロップします。

複数の画像をまとめてドロップすると、一括で処理できます。処理中は結果欄に進捗が表示され、完了した画像からサムネイルへ切り替わります。

クリップボードにコピーした画像を処理したい場合は、「クリップボード画像を処理」ボタンを押してください。

`ClipForVRChat.exe` にファイルやフォルダを直接ドラッグ&ドロップした場合、画像ファイルは入力画像として処理します。画像以外のファイルは `<元パス>.gpg`、フォルダは `<フォルダパス>.zip.gpg` としてGPG暗号化します。

### 4. URLをコピーする

Discord投稿がONの場合、投稿された画像URLを取得できます。

1枚だけ処理した場合は、設定がONなら自動でクリップボードへコピーされます。結果欄のサムネイルをクリックすると、その画像URLをもう一度コピーできます。

### 5. QRコードURLを確認する

設定画面で「QRコードURL検出」をONにすると、処理した画像内のQRコードを読み取ります。

QRコードにURLが含まれている場合は、Discord投稿の本文にそのURLを追記します。結果画面でも画像ごとに検出したURLを確認でき、URLをクリックするとクリップボードへコピーできます。

1枚の画像に複数のQRコードがある場合も、読み取れたURLをまとめて表示します。QRコードにURL以外の文字列が入っている場合は表示しません。

### 6. 結果をクリアする

結果欄右上の「クリア」ボタンで、画面上の結果を消せます。

クリアしても履歴は削除されません。あとから履歴画面で確認できます。

### 7. 履歴を見る

結果欄右上の「履歴」ボタンから、画像履歴画面を開けます。

履歴画面では、Discord投稿URL、ローカル保存ファイル、QRコードURLを含む過去の処理記録を確認できます。Ctrlキーを押しながらクリックすると複数選択、Shiftキーを押しながらクリックすると範囲選択できます。

選択した履歴は、削除可能なものだけをDiscord、ローカル、履歴から個別に削除できます。ピン止めした履歴は削除対象になりません。

Discord投稿、ローカル保存、QRコードURL取得のいずれも行われなかった画像は、結果画面や履歴には表示されません。

### 8. VRChat写真を自動処理する

設定画面の「VRChat写真自動処理」をONにすると、VRChatの写真フォルダを監視します。

VRChatで写真が保存されると、設定に応じてDiscord投稿、QRコードURL検出、ローカル保存を行います。VRChat写真だけ別のチャンネルへ投稿したい場合は、「VRChat写真用Webhook URL」を設定してください。空のままなら、通常投稿用Webhook URLへ投稿します。

### 9. スクリーンショットを自動処理する

設定画面の「スクリーンショット自動処理」をONにすると、WindowsのScreenshotsフォルダを監視します。

Win+Shift+Sなどで画像が保存されると、設定に応じてDiscord投稿、QRコードURL検出、ローカル保存を行います。スクリーンショットだけ別のチャンネルへ投稿したい場合は、「スクリーンショット用Webhook URL」を設定してください。空のままなら、通常投稿用Webhook URLへ投稿します。

### 10. VRChatを自動撮影する

設定画面の「自動撮影」タブでは、VRChat User CameraをOSCで操作し、指定した構図を順番に撮影できます。

自動撮影は初期状態では無効です。使う場合は自動撮影スケジュールをONにします。スケジュール詳細では、撮影間隔、初回待機時間、開始時撮影、最大実行回数を設定できます。開始時撮影は既定ONですが、AvatarBeacon basisが使える状態、またはフォールバックモードが使える状態になるまで待ってから実行します。

Stream方式はVRChat Stream CameraのSpout映像を内蔵の `spout-capture.exe` で直接受信して保存します。通常利用者向けzip内の `ClipForVRChat.exe` には、`spout-capture.exe` と実行に必要な `SpoutLibrary.dll` が埋め込まれており、初回使用時にClipForVRChatが管理フォルダへ展開して呼び出します。通常はユーザーが直接起動せず、ClipForVRChatが自動撮影時に呼び出します。Photo方式はVRChat標準写真を使うフォールバックで、VRChat側のシャッター音が発生します。

`spout-capture.exe` と `SpoutLibrary.dll` は、VRChat Stream Cameraの実体であるSpout senderから映像フレームを受け取るために必要です。Spout/DirectX/OpenGL/DLLロードなどのWindowsネイティブ処理を本体プロセスに混ぜると、GPUやDLLまわりの失敗がアプリ全体のクラッシュにつながりやすくなります。そのためv0.1.8では、Spout受信だけを小さい別プロセスへ分離し、本体exeにはそのhelper一式を埋め込みます。本体はOSC制御、設定、保存後処理、履歴、Discord投稿を担当します。

安全性を確認しやすいように、`spout-capture.exe` の役割は限定しています。利用可能なSpout senderの列挙、指定senderからの1フレーム受信、指定された保存先へのPNG書き込み、結果JSONの出力だけを行います。Discord Webhook URLや設定ファイルを読み取らず、画像や設定を外部へ送信するネットワーク通信も行いません。`SpoutLibrary.dll` はSpout2 SDKの実行時DLLで、`spout-capture.exe` がsender列挙とフレーム受信に使います。ヘルパーのソースは `tools/spout-capture/main.cpp`、ビルド設定は `tools/spout-capture/CMakeLists.txt` にあり、Spout2はBSD 2-Clause Licenseの `SpoutLibrary` を使用します。

配布物を確認するときは、公式Releaseまたは公式配布元から通常利用者向けzipを取得してください。通常利用者向けzipには `ClipForVRChat.exe`、`ClipForVRChat-vX.Y.Z-windows-amd64.exe.asc`、`Release-signing-public-key.url`、`README.md`、`LICENSE`、`AvatarBeacon_v0.0.1.unitypackage` が入ります。Release buildでは、同梱する `AvatarBeacon_v0.0.1.unitypackage` をAvatarBeacon Releaseの `.asc` 署名と `release-signing@hato.life` の固定fingerprintで検証します。PGP署名まで確認する場合は、zip内の `ClipForVRChat.exe` と `.exe.asc` を検証してください。

Release Assetsには、検証・切り分け用に `ClipForVRChat-vX.Y.Z-windows-amd64-separated.zip` も添付します。このzipには `ClipForVRChat.exe`、`spout-capture.exe`、`SpoutLibrary.dll`、`Spout2-LICENSE.txt`、`AvatarBeacon_v0.0.1.unitypackage` が入ります。helper単体確認や外部helper指定での切り分けが必要な場合だけ使用してください。`spout-capture.exe` だけ、または `SpoutLibrary.dll` だけではStream方式は動作しません。

Stream方式で `Spout helperが見つかりません`、`Spout helperは見つかりましたが実行確認に失敗しました` と表示される場合は、単一exe版を使っているか、分離版zipを使う場合は `spout-capture.exe` と `SpoutLibrary.dll` が両方残っているか確認してください。

自動撮影では、画像と同じ場所にsidecar JSONを保存できます。sidecar JSONには撮影時刻、構図、Stream sender情報、VRChat output logから取得できた同席ユーザー、world ID、instance ID、画像SHA256を記録します。PNG/JPEG画像内にも、設定に応じて自動撮影メタデータを埋め込めます。

プレイヤー基準構図は、AvatarBeacon導入済みアバターから受け取る `avatar_osc` basis を基本にします。標準OSCだけでは動かず、positionとyawはHead基準で、player root そのものではありません。専用ギミックなしの場合は `manual` basis で手動保存した基準位置をフォールバックとして使えます。確認手順は `docs/v0.1.8-avatar-osc-basis-verification.md` と `docs/v0.1.8-player-local-verification.md`、AvatarBeaconの詳細は `docs/avatarbeacon-spec.md` を参照してください。

AvatarBeaconが使えない環境では、フォールバックモードを使えます。フォールバックモードでは、VRChat内であらかじめ配置したローカルアンカーCameraを使うため、ClipForVRChatはCamera Poseを送信しません。AvatarBeacon未受信時の自動ON、受信復帰時の自動OFFも設定できますが、どちらも既定OFFです。

AvatarBeacon は `avatar_osc` basis 確認用の汎用アバターギミックで、ClipForVRChat専用ではありません。AvatarBeaconの元ファイルは別リポジトリで管理し、ClipForVRChatの通常zipと分離版zipには当面 `AvatarBeacon_v0.0.1.unitypackage` を同梱します。Unityへimportすると `Assets/PoppoWorks/AvatarBeacon/...` が追加される形を想定しています。

RC確認時の詳しい確認手順は、`docs/v0.1.8-stream-spout-verification.md`、`docs/v0.1.8-player-local-verification.md`、`docs/v0.1.8-embedded-metadata-verification.md` を参照してください。

### 11. コマンドラインで確認する

PowerShell、cmd、Git Bash からバージョンとヘルプを確認できます。

```pwsh
.\ClipForVRChat.exe --version
.\ClipForVRChat.exe --help
```

`ClipForVRChat.exe` は通常のGUIアプリとしてビルドしているため、コマンドライン出力が次のプロンプト表示後に出る場合があります。これは通常起動時に余計なコンソールを表示しないための仕様です。

## Discord設定

Discord投稿を使うには、投稿先チャンネルの Webhook URL が必要です。

### Webhook URL の作り方

1. Discordで投稿先にしたいサーバーとチャンネルを開きます。
2. チャンネルの設定を開きます。
3. 「連携サービス」または「インテグレーション」から Webhook を作成します。
4. 作成した Webhook のURLをコピーします。
5. ClipForVRChat の設定画面を開き、「Discord Webhook URL」に貼り付けます。
6. 「Discord投稿」をONにして保存します。

Discord公式の説明:

https://support.discord.com/hc/ja/articles/228383668-%E3%82%A6%E3%82%A7%E3%83%96%E3%83%95%E3%83%83%E3%82%AF%E3%81%AE%E3%81%94%E7%B4%B9%E4%BB%8B

Webhook URL は、投稿できる権限を持つURLです。知らない人に共有しないでください。

Discord投稿では、本文に含まれる `@everyone`、`@here`、ユーザー/ロールメンションが通知を発生させないよう、メンションを無効化して送信します。

## 設定

設定画面は、画面右上の「設定」ボタンから開けます。上部のタブでカテゴリを切り替えます。

### 機能

- VRChat写真自動処理: VRChat上で撮影されたときに処理します。
- スクリーンショット自動処理: Win + Shift + Sでスクリーンショットが撮られたときに処理します。
- QRコードURL検出: 画像内のQRコードからURLを取得します。

### 自動撮影

- AvatarBeacon受信状態: `avatar_osc` basisの状態、最終受信、position、yaw、エラーを確認できます。
- フォールバックモード: AvatarBeaconなしで、VRChat内に配置済みのローカルアンカーCameraを使います。
- 自動撮影スケジュール: 撮影間隔、初回待機時間、開始時撮影、最大実行回数を設定できます。
- 撮影と出力: Stream/Photo方式、Spout helper、sender、自動選択、録画デバッグ、保存先、保存形式、ファイル名、sidecar JSON、画像埋め込みメタデータを設定できます。
- Camera状態復元: 撮影後にUser CameraのMode、Pose、Streaming、Zoom、Exposure、表示対象mask類などをできるだけ戻します。
- 構図: 正面、背後、斜めなどの `player_local` 構図を管理できます。

### OSC

- OSCホスト: 通常は `127.0.0.1` です。
- OSC送信ポート: ClipForVRChatからVRChatへ送るUDPポートです。初期値は `9000` です。
- OSC受信ポート: VRChatから外部アプリへ届くUDPポートです。初期値は `9001` です。
- カメラOSCリセット: VRChat User CameraのOSC操作状態を戻したいときに使います。
- プレイヤー基準の取得元: 通常は `avatar_osc`、専用ギミックなしの場合は `manual` を使います。
- OSC転送: VRChatから受信したOSC packetを、他アプリ用の別ポートへ転送できます。
- OSC受信ログ/送信ログ: 一時ログを確認し、フィルタやコピーができます。診断用には `logs/osc_recieve.jsonl` と `logs/osc_send.jsonl` に保存されます。
- OSCデバッグ送信: 現在のVRChat OSC送信先へ任意のOSCを送れます。

### 処理

- ローカル保存: 処理した画像をローカルに保存します。
- 出力先フォルダ: ローカル保存時の保存先です。初期値は `./output` です。
- ファイル名サフィックス: ローカル保存時のファイル名末尾に付ける文字です。
- 出力形式: PNGまたはJPGを選べます。
- JPEG品質: JPG出力時だけ使います。

### その他

- PC起動時に自動起動: Windows Startupフォルダに現在のexeへのショートカットを作成または削除します。
- 更新確認: 起動時にGitHub Releasesを確認し、新しいバージョンがあるか調べます。初期値はONです。
- 更新通知: 新しいバージョンが見つかったとき、画面上部に通知を表示します。初期値はONです。

### Discord投稿

- Discord投稿: 処理した画像をDiscord Webhookへ投稿します。
- 投稿URLの自動コピー: Discordに投稿したURLをクリップボードに保存します。
- 通常投稿用Webhook URL: Discordで作成したWebhook URLを貼り付けます。
- VRChat写真用Webhook URL: VRChat写真だけ別のDiscordチャンネルに投稿したい場合に入力します。空なら通常投稿用Webhook URLを使います。
- スクリーンショット用Webhook URL: スクリーンショットだけ別のDiscordチャンネルに投稿したい場合に入力します。空なら通常投稿用Webhook URLを使います。
- 自動撮影用Webhook URL: 自動撮影だけ別のDiscordチャンネルに投稿したい場合に入力します。空なら通常投稿用Webhook URLを使います。

Discord投稿がOFFの場合、投稿URLの自動コピーとWebhook設定はグレーアウトします。VRChat写真自動処理とスクリーンショット自動処理は、縮小、QRコードURL検出、ローカル保存を行えるためグレーアウトしません。ローカル保存がOFFの場合、出力形式を含むローカル保存に関係する設定はグレーアウトします。
自動撮影の `player_local` は AvatarBeacon の `avatar_osc` basis を既定にし、専用アバターギミックがない場合だけ `manual` basis またはフォールバックモードに切り替えて使います。受信状態や最終受信時刻、position/yaw は設定画面の自動撮影タブとOSCタブで確認できます。

## 画面

### 初期画面

画像をドラッグ&ドロップする入口です。

左側にドロップ領域、右側に結果欄があります。画像をまだ処理していない場合、結果欄には「まだ処理結果はありません。」と表示されます。

新しいバージョンが見つかった場合は画面上部に通知を表示します。通知内のボタンからGitHub ReleasesまたはBOOTHを開けます。通知右端の×ボタンで、その起動中は通知を閉じられます。

### 結果画面

処理した画像のサムネイルを表示します。

サムネイルにマウスを重ねると、利用できる操作が表示されます。URLがある場合はURLコピー、ローカル保存ファイルがある場合は保存先フォルダの表示ができます。

QRコードURL検出がONの場合、画像から読み取れたURLもサムネイル内に表示されます。

### 設定画面

保存先、出力形式、Discord投稿などを変更できます。設定項目は上部タブでカテゴリごとに切り替えて表示します。

出力先フォルダは、入力欄へ貼り付けても、「選択」ボタンからフォルダを選んでも設定できます。

未保存変更がある場合は、変更されたタブや項目が強調されます。設定画面を閉じる、別画面へ移動する、またはウィンドウを閉じる前に、保存、破棄、キャンセルを選べます。

### 使い方画面

アプリ内で基本的な使い方を確認できます。

Discord Webhook の作り方を確認するための公式リンクもあります。

### 情報画面

アプリのバージョン、ライセンス、GitHub、バグ報告先を確認できます。

不具合報告用データを作成できます。作成時には確認用zipと、添付用の暗号化済み `.zip.gpg` が作成されます。

バグ報告:

https://github.com/hatolife/ClipForVRChat/issues

### 画像履歴画面

過去の処理記録を確認できます。Discord投稿URL、ローカル保存ファイル、QRコードURLの有無が表示されます。

履歴画面では、Discordから削除、ローカルから削除、履歴から削除を個別に実行できます。ピン止めした履歴は削除対象になりません。

ローカル保存ファイルが存在しない履歴では、ローカルファイルのパスは表示されません。

### VRChat写真自動処理

設定画面でONにすると、VRChat写真フォルダを監視します。

アプリを起動している間に新しい写真が保存されると、自動で結果画面に追加されます。

### スクリーンショット自動処理

設定画面でONにすると、Screenshotsフォルダを監視します。

アプリを起動している間にWin+Shift+Sなどで新しいスクリーンショットが保存されると、自動で結果画面に追加されます。

## 保存されるファイル

- `config.json`: アプリ設定
- `history.json`: 取得した画像URLの履歴
- `output` フォルダ: ローカル保存した縮小画像
- `logs/YYYY-MM-DD.log`: 診断ログ
- `logs/osc_send.jsonl`: ClipForVRChatから送信したOSCログ
- `logs/osc_recieve.jsonl`: VRChatから受信したOSCとforwardのログ
- `diagnostics/`: 不具合報告用データの確認用zipと暗号化済み `.zip.gpg`
- `%LOCALAPPDATA%` 配下の単一起動管理フォルダ: 未保存設定の下書きと、単一exe版から展開したSpout helperの一時キャッシュ
