# AvatarBeaconからHeadForwardAnchorを削除する

## 指示

> HeadForwardAnchor これ無くてもいいなら削除したい

## 文脈

AvatarBeaconは直前の変更で、位置用 `point` もHeadへ追従する方針になった。
この状態では `HeadForwardAnchor` もHeadへ追従しており、向き用に別Transformを持つ理由が薄い。

## 解釈

`offset_rot` のRotationConstraint参照先を `HeadForwardAnchor` から `point` へ差し替え、`HeadForwardAnchor` GameObjectをPrefabから削除する。
位置と向きはどちらも `point` をHeadへ追従させて取得する。

## 問題

同じHeadへ追従するTransformが2つあり、導入手順とPrefab構造が余計に複雑になっている。

## 期待する挙動

- `AvatarBeacon_main.prefab` と `AvatarBeacon_12.prefab` に `HeadForwardAnchor` が存在しない。
- `avatar_beacon/forward/*` の入力元はHeadへ追従する `point` になる。
- 導入手順は `point: Head` のみになる。
- README、NOTICE、仕様書が `HeadForwardAnchor` なしの構造と一致する。

## 受け入れ条件

- [x] Prefab内に `HeadForwardAnchor` GameObjectが残っていない。
- [x] `offset_rot` のRotationConstraintが `point` を参照する。
- [x] `point` のBone Proxy targetがHead相当になっている。
- [x] README、NOTICE、仕様書から `HeadForwardAnchor` 必須説明が消えている。
- [x] AvatarBeacon source tree検証が通る。

## 完了メモ

- `AvatarBeacon_main.prefab` と `AvatarBeacon_12.prefab` から `HeadForwardAnchor` GameObjectを削除した。
- `offset_rot` のRotationConstraint参照先を `HeadForwardAnchor` から `point` へ差し替えた。
- `point` のBone Proxy targetはHead相当の `boneReference: 10` とし、位置とforward/yawの入力元を統合した。
- README、NOTICE、仕様書、OSS表示文言を `HeadForwardAnchor` なしの構造へ更新した。
- `node scripts/check-source-tree.mjs` が通った。
