# 自動撮影説明を改善し、カメラ自動起動/終了を既定OFFにする

## 問題

自動撮影タブの説明が長く、詳細設定画面でもタブ、保存ボタン、閉じるボタン、概要説明が表示されるため、詳細設定だけを編集したいときに見通しが悪い。

また、現在は自動撮影開始時にUser CameraのModeやStream CameraのStreamingをOSCで自動操作しているが、VRChat側ではCamera UIを閉じるとCamera OSCが効かない不具合が報告されており、自動起動/自動終了が予期しない挙動になる場合がある。

## 期待する挙動

- 自動撮影の説明文を短くし、AvatarBeacon受信状態を同じ説明領域で分かりやすく表示する。
- AvatarBeacon受信状態が正常なときは詳細説明を折り畳み、エラー時またはユーザーが開いたときだけ詳細説明を表示する。
- 詳細設定画面では設定タブ、設定保存ボタン、設定閉じるボタン、自動撮影の概要説明を表示しない。
- 撮影方式の説明を、Stream方式とPhoto方式の説明に改行して分ける。
- カメラ自動起動と自動終了のON/OFF設定を用意し、既定OFFにする。
- カメラ自動起動OFFのときは、ユーザーがゲーム内でカメラを起動しておく前提の動作にする。

## 受け入れ条件

- [x] 自動撮影タブの概要説明が新しい短い説明に置き換わる。
- [x] 概要説明内に「現在のAvatarBeacon受信状態: エラー」または「現在のAvatarBeacon受信状態: 成功 x, y, z, yaw°」が表示される。
- [x] AvatarBeacon受信状態が正常なときは詳細説明が折り畳まれる。
- [x] AvatarBeacon受信状態がエラーのとき、または詳細情報をクリックしたときは詳細説明が表示される。
- [x] 自動撮影の詳細設定画面では設定タブ、保存ボタン、閉じるボタン、概要説明が表示されない。
- [x] 撮影方式説明がStream方式とPhoto方式の2行説明になる。
- [x] カメラ自動起動、自動終了の設定が追加され、既定OFFになる。
- [x] カメラ自動起動OFF時に `/usercamera/Mode`、`/usercamera/SmoothMovement`、`/usercamera/Streaming=true` を事前送信しない。
- [x] カメラ自動終了OFF時に `/usercamera/Streaming=false` と `/usercamera/Close` を撮影後に送信しない。

## 実装メモ

- `openCameraBeforeBatch` を追加し、既定OFFにした。
- 既存の `closeCameraAfterBatch` も既定OFFに変更した。
- `openCameraBeforeBatch=false` のとき、撮影前のCamera Mode切り替え、SmoothMovement、Streaming開始、Stream撮影直前の再送を行わない。
- AvatarBeacon受信状態は自動撮影タブでもポーリングする。
- Camera OSC不具合の説明はVRChat Feedbackの報告を参照する。
