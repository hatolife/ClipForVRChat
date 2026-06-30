# ローカルプレイヤー基準Poseの取得可否を実機調査する

## 問題

`player_local` を実装するには、撮影時点のローカルプレイヤー位置と向きが必要である。
現行調査では、VRChat公式OSCの `/usercamera/Pose` はカメラPoseであり、標準Avatar Parametersにもプレイヤーroot位置/yawは見当たらない。

## 期待する挙動

標準OSC/OSCQuery/VRChatログ/実機挙動から、ローカルプレイヤー基準Poseを取得できるかを確認し、取得できない場合は自動 `player_local` を未対応として扱う判断材料を残す。

## 受け入れ条件

- `/usercamera/Pose` のGet値が、Camera ModeごとにカメラPoseのみを返すのか、プレイヤー基準として使える値を返す余地があるのかを実機で確認する。
- `/usercamera/Mode=2`、`/usercamera/Streaming`、`/usercamera/LookAtMe` の組み合わせで、カメラ起動直後のPoseを基準Poseとして使えるか確認する。
- OSCQueryで公開されるendpoint一覧を取得し、ローカルプレイヤー位置/向きに相当するendpointがないか確認する。
- VRChat output logに、撮影時点で使えるローカルプレイヤー位置/向き情報が出ていないか確認する。
- 標準機能だけで取得できない場合、手動基準Pose、専用ワールド/アバター補助OSC、または未対応表示のどれを採用するか決める。
- 調査結果に基づき、`player_local` 実装を進めてよい条件と、進めてはいけない条件を明記する。

## 対応内容

- v0.1.8実装に合わせて対応済み。詳細は `feat/v0.1.8-resolve-issues` の実装、README、SPEC、RELEASE_NOTESを参照。
