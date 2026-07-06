# 自動撮影機能を使用するための条件を整理する

## 問題

ユーザーが自動撮影機能を使うために必要なVRChat側、ClipForVRChat側、撮影方式別の条件が複数の実装・検証資料に分散している。

## 期待する挙動

自動撮影を使うための前提条件、Stream方式の条件、Photo方式の条件、構図基準の条件、任意機能の条件を列挙できる。

## 受け入れ条件

- [x] 現行実装と検証資料から利用条件を確認する。
- [x] ユーザー向けに条件を簡潔に列挙する。

## 調査結果

2026-07-05に、`src/internal/appcore/autocapture.go`、`src/internal/appcore/config.go`、`src/internal/appcore/player_local.go`、`src/app.go`、`docs/v0.1.8-stream-spout-verification.md`、`docs/v0.1.8-avatar-osc-basis-verification.md`、`docs/v0.1.8-player-local-verification.md`、`avatar-gimmicks/AvatarBeacon/Assets/PoppoWorks/AvatarBeacon/README.md` を確認した。

自動撮影の利用条件は、ユーザー向け回答として前提条件、Stream方式、Photo方式、構図基準、任意機能に分けて整理した。



## 調査結果

現行実装上、ユーザーが自動撮影を使う条件は次の通りです。

### 共通の必須条件

- VRChat側でOSCを有効にする。
- ClipForVRChatがVRChatへOSC送信できること。既定は 127.0.0.1:9000。
- ClipForVRChatがVRChatからOSC受信できること。既定は 127.0.0.1:9001。
- 自動撮影で撮る構図が1つ以上ONになっていること。
- 撮影方式が Stream または Photo のどちらかであること。
- player_local 構図を使う場合、撮影前にプレイヤー基準Poseが解決できること。

### Stream方式の条件

- 撮影方式を Stream にする。
- VRChatのStream Camera / Spoutが使える状態であること。
- spout-capture.exe と SpoutLibrary.dll が利用可能であること。通常版は内蔵helper、分離版は同じフォルダ配置が必要。
- helper確認 が成功すること。
- Spout senderが検出できること。
- senderが複数あり自動選択できない場合は、VRChat Stream Cameraのsender名を選ぶこと。
- 取得画像が黒一色、白一色、透明、未更新フレームではなく、有効なStream Camera映像であること。

### Photo方式の条件

- 撮影方式を Photo にする。
- VRChat標準写真の保存先が正しく設定されていること。既定は %USERPROFILE%/Pictures/VRChat。
- User Cameraが表示され、/usercamera/Capture によりVRChat側で新しい写真ファイルが保存されること。
- 撮影後30秒以内に新規写真ファイルを検出できること。

### player_local構図の条件

- 既定の avatar_osc を使う場合、AvatarBeacon導入済みアバターが必要。
- VRChat OSC Avatar Parametersで coord/* と forward/* が届くこと。
- AvatarBeacon受信状態が ready かつ新鮮であること。既定の鮮度は3秒。
- 専用ギミックなしで使う場合は manual に切り替え、現在Poseをmanual基準として保存しておくこと。

### 開始時撮影の追加条件

- 自動撮影スケジュールがONであること。
- 開始時撮影 がONであること。
- 開始時撮影ではAvatarBeacon basis readyを最大30秒待つ。
- VRChat output log監視がONの場合、world/instance情報が取得でき、3秒安定することも待つ。

### 任意機能の条件

- 同席ユーザー、world/instance、sidecar JSON用メタデータを取りたい場合はVRChat output log監視が必要。
- Discord自動投稿を使う場合は、自動撮影用Webhook URL、または通常投稿用Webhook URLが必要。
- 画像埋め込みメタデータやsidecar JSONは自動撮影自体の必須条件ではない。

