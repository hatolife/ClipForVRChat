# このPoseへカメラ移動でavatar_osc基準のplayer_local変換を使う

## 問題

`v0.1.8-a36` で自動撮影タブの「このPoseへカメラ移動」を押すと、`playerLocal.basisSource` が `avatar_osc` でAvatarBeacon受信もreadyになっているにもかかわらず、次のmanual基準Poseエラーが出る。

```text
プレイヤー基準Poseが未設定のため、player_local構図を撮影できません。自動撮影タブで現在Poseをプレイヤー基準として保存してください
```

また、カメラ移動先が常にワールド座標のように見え、プレイヤーの位置と向きを基準にしたローカル構図として移動できていない。

a36ログでは `auto-capture move camera open begin` 後に `auto-capture camera pose resolve error` が出ており、撮影本体の `run_once` では `player_local` の `resolved_pose` が計算されている。
このため、カメラ移動APIだけが古いmanual基準の解決経路を使っている可能性が高い。

## 期待する挙動

- `playerLocal.basisSource=avatar_osc` の場合、「このPoseへカメラ移動」はAvatarBeaconから受信した最新のposition/yawをbasisとして使う。
- `player_local` 構図は、プレイヤー位置と向きを基準にworld poseへ変換してから `/usercamera/Pose` へ送信する。
- manual基準Pose未設定エラーは、basis source が `manual` の場合だけ出す。
- カメラ移動時のログに、basis source、basis pose、local pose、resolved world pose、view id/name、送信先を出し、変換不具合を追跡できる。

## 受け入れ条件

- [x] `MoveCameraToView` が撮影本体と同じ `avatar_osc` basis解決を使う。
- [x] `avatar_osc` ready時にmanual基準Pose未設定エラーが出ない。
- [x] `player_local` のlocal poseとAvatarBeacon basisからresolved world poseが計算される。
- [x] カメラ移動の診断ログに「プレイヤーがここで、設定がこうだから、カメラをここへ移動する」が分かる情報を出す。
- [x] basisがmissing/partial/staleの場合は、AvatarBeacon受信状態に即したエラー文にする。
- [x] Goテストを追加または更新する。

## 実装メモ

- 2026-07-05: `MoveCameraToView` で対象viewが `player_local` かつ basis source が `avatar_osc` の場合、App層で最新AvatarBeacon basisを解決して `cfg.AutoCapture.PlayerLocal.BasisPose` に注入してから appcore のカメラ移動処理へ渡すようにした。
- 2026-07-05: `camera pose resolved` ログへ `basis_source`、`basis_pose`、`local_pose`、`resolved_pose` を出すようにした。
- 2026-07-05: manual未保存かつavatar_osc readyの状態で `/usercamera/Pose` が送信されるGoテストを追加した。
