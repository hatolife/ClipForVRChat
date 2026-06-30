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
