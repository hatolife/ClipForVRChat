# フォールバックモードを既定ONにしOFF時は通常撮影失敗を許容する

## 問題

#211で追加したローカルアンカー配置済みカメラのフォールバックモードは手動ON/OFFできるが、既定値がOFFで、AvatarBeacon未受信時には実行時に自動でONへ切り替わる。
ユーザーがOFFを選んだ場合でも自動フォールバックにより、通常のPose/Zoom送信経路の失敗を確認しにくい。

## 期待する挙動

フォールバックモードは既定ONにする。
ユーザーがOFFにした場合は自動フォールバックを使わず、AvatarBeacon basis未受信など通常撮影に必要な条件が不足していれば撮影準備または撮影が失敗する。

## 受け入れ条件

- 新規/未指定設定では `autoCapture.capture.preplacedLocalAnchor` がONになる。
- 設定画面のフォールバックモードトグルは既定ONで表示される。
- ユーザーがOFFにした場合、AvatarBeacon basis未受信でも自動的にONへ戻さない。
- OFF時にbasis不足なら通常のplayer_local basisエラーとして撮影に失敗する。
- ON時は従来通りPose/Zoom/Mode/Streamingなどを送らず、配置済みカメラ前提で撮影だけ行う。
