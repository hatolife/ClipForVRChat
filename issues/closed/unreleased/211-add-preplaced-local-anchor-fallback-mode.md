# ローカルアンカー配置済みカメラを使うフォールバックモードを追加する

## 問題

`player_local` 構図は通常 AvatarBeacon などの基準Poseが必要で、アバターギミックなしでは使いにくい。
一方、VRChat内でユーザーがUser Cameraをローカルアンカーにして事前配置すれば、ClipForVRChatはカメラPoseを動かさず撮影操作だけ送る運用ができる。

## 期待する挙動

構図設定の上にフォールバックモードのON/OFFを置き、ONのときはユーザーがVRChat内でローカルアンカー配置したカメラをそのまま使う。
ClipForVRChatは構図PoseやZoomなどを送らず、撮影に必要な最小限のOSCだけを送る。

## 受け入れ条件

- 構図設定の上にフォールバックモードON/OFFのトグルを表示する。
- フォールバックONではPhoto方式で `/usercamera/Capture` のみ送る。
- フォールバックONではStream方式でPose/Zoom/Mode/Streaming等を送らず、現在のSpout映像を取得する。
- フォールバックONでは `player_local` basis未受信でも撮影に進める。
- フォールバックONでは「このPoseへカメラ移動」を送信しない。
- フォールバックONでは、OFF時だけ使う構図Pose、Zoom、座標系、現在Pose保存、初期Poseリセットの操作をグレーアウトして混乱を防ぐ。
- AvatarBeaconからのbasis受信が30秒以上ない場合は、撮影実行時の有効モードを自動でフォールバックにする。
- AvatarBeaconのbasis受信が戻った場合は、撮影実行時の有効モードを通常モードに戻す。
- 既定はOFFで、従来の構図制御は維持する。

## 対応内容

- `autoCapture.capture.preplacedLocalAnchor` を追加した。2026-07-07に #222 で既定ONへ変更した。
- 構図設定画面の先頭にフォールバックモードのトグルを追加した。
- フォールバックONではCamera Mode、Streaming、Pose、Zoom、表示マスク、Close、復元送信をスキップする。
- Photo方式では撮影時に `/usercamera/Capture` だけ送る。
- Stream方式ではOSCでカメラ操作をせず、現在のSpout映像を取得する。
- フォールバックONではAvatarBeacon basis待ちとbasis解決をスキップする。
- フォールバックONでは「このPoseへカメラ移動」をUI/API双方で送信しない。

## 追加対応

- 2026-07-07: フォールバックON時に、構図カード内のPose/Zoom/座標系など、撮影時に送信されないOFF用設定をグレーアウトする。
- 2026-07-07: フォールバックモードの説明を、アバターギミック未導入時の利用条件とVRChat内でカメラを出しっぱなしにする手順が分かる文言へ変更する。
- 2026-07-07: AvatarBeaconのbasis最終受信から30秒以上経過した場合の自動切り替えは #222 で廃止し、ユーザーが選んだフォールバックON/OFFを尊重する。OFF時はbasis不足なら通常撮影エラーとして扱う。
