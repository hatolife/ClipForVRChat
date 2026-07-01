# `player_local` 座標仕様を確定する

## 問題

自動撮影の構図は現在ワールド座標として扱われており、正面/背後/斜めの値がプレイヤー位置や向きに追従しない。
`CoordinateSpace` に `dolly_local` や `template_relative` という未接続の値もあり、仕様と実装の境界が曖昧である。

## 期待する挙動

`player_local` を明示的な座標系として定義し、構図Poseが撮影時点のローカルプレイヤー基準で解釈される。
正面は顔を中心に、背後は後頭部側から少し見下ろし、斜めは正面寄りの横側から少し見下ろす構図になる。

## 受け入れ条件

- `player_local` の原点を、root基準ではなく顔/目線付近を狙えるアンカーとして定義する。
- ローカル軸、単位、Yaw基準、Pitch/Rollの扱い、Euler角の適用順を明文化する。
- 構図は「カメラ位置offset」と「注視点offset」から回転を算出する方式にするか、固定rotationを保持する方式にするかを決める。
- 正面/背後/斜めの初期offset、注視点、ズーム初期値を具体値で決める。
- `world` 既存設定との互換性と、未接続の `dolly_local` / `template_relative` の扱いを決める。
- 仕様が `issues/136-investigate-player-local-camera-coordinate-spec.md` と実装issueから参照できる状態になっている。

## 対応内容

- v0.1.8実装に合わせて対応済み。詳細は `feat/v0.1.8-resolve-issues` の実装、README、SPEC、RELEASE_NOTESを参照。

## 再監査メモ

- 2026-07-01: [#164](164-audit-v018-completed-items.md) の再監査で未達が見つかったため、完了扱いを取り消して `要対応` に戻した。

## 実装前調査メモ

実装方針:

- v0.1.8では、標準OSCだけでローカルプレイヤーroot位置/Yawを自動取得しない。`player_local` は「手動保存したプレイヤー基準Pose」を原点/向きとして扱うローカル座標系と定義する。
- 原点は `AutoCapture.PlayerLocal.BasisPose.Position`、Yaw基準は `BasisPose.Rotation.Y` とする。Pitch/Rollは基準Poseへローカル回転を加算する。
- ローカル位置はメートル相当のUser Camera Pose単位で扱う。
- 軸定義は現行 `TransformPlayerLocalPose()` に合わせ、基準Yaw 0度で `+X=右`, `+Y=上`, `+Z=基準前方` とする。Yawが変わる場合はXZ平面だけをYaw回転する。
- 初期構図は `player_local` に変更し、以下の目安値を採用する。
  - 正面: `position=(0, 0.0, 1.0)`, `rotation=(0, 180, 0)`, `zoom=1.0`
  - 背後: `position=(0, 0.35, -1.6)`, `rotation=(12, 0, 0)`, `zoom=1.0`
  - 斜め: `position=(0.8, 0.2, 1.1)`, `rotation=(8, -145, 0)`, `zoom=1.0`
- 上記は「BasisPoseが顔/目線付近を向くように保存されている」前提の初期値である。真のプレイヤーroot追従ではないことをREADME/SPECへ明記する。

対象ファイル:

- `src/internal/appcore/config.go`
- `src/internal/appcore/player_local.go`
- `src/SPEC.md`
- `README.md`

小タスク:

- `defaultCameraView()` を `coordinateSpace` 引数付きにするか、初期構図作成後に `CoordinateSpace: "player_local"` を設定する。
- `DefaultCameraViewByID()` と `ResetCameraPoseToDefault()` が同じ初期値を返すことを確認する。
- `template_relative` など旧値はNormalizeで初期構図へ移行する現行処理を維持する。
- `player_local` の制限を「未実装」ではなく「手動基準Pose方式」として仕様化する。

確認方法:

- 新規configの3構図が `player_local` になる。
- 旧configで `template_relative` やゼロPoseの初期構図が `player_local` 初期値へ移行する。
- 未キャリブレーション時の撮影失敗は、手動基準Pose保存を促す明確な文言になる。
