# AvatarBeaconの位置取得をHips基準からHead基準へ変更する

## 指示

> Hips基準をやめて 位置座標を取得もHeadにしたい　これはチケット作成して作業

## 文脈

AvatarBeaconは現状、`point` をHipsへ追従させて `avatar_beacon/coord/*` の位置を出し、`HeadForwardAnchor` をHeadへ追従させて `avatar_beacon/forward/*` の向きを出す構成になっている。

## 解釈

位置と向きの入力元を分離せず、`avatar_beacon/coord/*` もHead由来の座標として出す。`avatar_beacon/forward/*` は従来通りHead由来の向きとして維持する。

## 問題

Hips基準では、顔や視点を中心にした自動撮影構図で意図した位置からずれる場合がある。

## 期待する挙動

- `AvatarBeacon_main.prefab` と `AvatarBeacon_12.prefab` の位置用 `point` がHeadへ追従する。
- README、仕様書、アプリ内説明が「位置もHead基準」と一致する。
- 既存のOSC parameter pathは変更しない。

## 受け入れ条件

- [x] `point` のBone Proxy targetがHead相当になっている。
- [x] `avatar_beacon/coord/*` の説明がHead基準位置として更新されている。
- [x] `avatar_beacon/forward/*` の説明がHead基準向きとして維持されている。
- [x] Vue template変更に必要な検証スクリプトが通る。

## 完了メモ

- `AvatarBeacon_main.prefab` と `AvatarBeacon_12.prefab` の位置用 `point` Bone ProxyをHead相当の `boneReference: 10` に変更した。
- `HeadForwardAnchor` はforward/yaw用Constraintの参照先として残した。`point` もHead基準になったため概念上は統合余地があるが、既存Prefab構造を壊さないため今回は削除しない。
- README、仕様書、検証手順、アプリ内説明、OSSライセンス表示内の変更点説明をHead基準へ更新した。
- `node scripts/check-frontend-template-literals.mjs`、`node scripts/check-wails-api-surface.mjs`、`avatar-gimmicks/AvatarBeacon` 側の `node scripts/check-source-tree.mjs` が通った。
