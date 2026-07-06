# player_local基準Poseの挙動を見直す

## 問題

`player_local` の basis source が `manual` だけに見えると、プレイヤー移動に自動追従するものか、手動保存した基準Poseで固定するものかが分かりにくい。
実際には、AvatarBeacon導入済みアバターから受け取る `avatar_osc` basis を既定にし、専用ギミックなしの場合だけ `manual` basis をフォールバックとして使う。

## 期待する挙動

初期状態でも自動撮影を開始した利用者が、`avatar_osc` 既定と `manual` フォールバックの関係を迷わず理解できる。
`player_local` が使える条件と、`manual` に切り替えて基準Poseを保存する条件がUIとドキュメントで明確になっている。

## 受け入れ条件

- [x] `player_local` の既定 basis source が `avatar_osc`、フォールバックが `manual` だと分かる。
- [x] `manual` basis が手動保存した基準Poseを使うことが、README/verification docsで明確になっている。
- [x] プレイヤー移動後に再保存が必要な範囲がREADME/verification docsへ反映されている。
- [x] 既存の AvatarBeacon / `avatar_osc` 実装で、manual-only という誤解は解消できている。

## メモ

- 標準OSCの `/usercamera/Pose` はUser Cameraのworld poseであり、ローカルプレイヤーrootの位置/Yawを返すものではない。
- a14時点の問題提起は手動基準Pose方式だったが、v0.1.8では AvatarBeacon により `avatar_osc` を既定basis sourceとして使える。
- 現在のプレイヤー位置取得方法の再調査は [#139](139-investigate-player-basis-source.md) に追記した。
- 2026-07-04: AvatarBeacon + docs で `player_local` の basis source を既定 `avatar_osc` / fallback `manual` として整理できたため、issue 172 は解消済みとして閉じる。
