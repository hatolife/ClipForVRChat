# AvatarBeacon mainと_12のfloat精度を調査する

## 指示

> AvatarBeaconのmainと_12のfloatの精度を調査

## 文脈

AvatarBeaconは通常利用向けの `AvatarBeacon_main.prefab` と、高精度版の `AvatarBeacon_12.prefab` を提供している。
`main` は6parameterのcentered float、`_12` は12parameterのmagnitude/sign分離方式として設計されているが、実際にどの程度の座標分解能になるかを整理する必要がある。

## 解釈

Prefab YAML、仕様ドキュメント、ClipForVRChat側の復元処理を静的に確認し、OSC floatそのものの精度ではなく、Contact Receiverの検出半径と復元スケールから決まる有効精度を比較する。
Unity/VRChat実機でのContact出力の量子化やOSC送信ゆらぎは、静的調査だけでは確定できないため、必要に応じて未確認範囲として分ける。

## 問題

- `main` と `_12` の精度差がREADME上の「高精度」という表現だけでは定量的に分からない。
- Contact Receiverの半径、center offset、ClipForVRChat側の復元式が分散しており、精度の根拠を追いにくい。
- 実機確認が必要な範囲と静的に判断できる範囲を分けておく必要がある。

## 期待する挙動

- `AvatarBeacon_main.prefab` と `AvatarBeacon_12.prefab` のfloat表現方式と分解能が説明できる。
- ClipForVRChat側の復元式と設定値が、Prefab側の意図と一致しているか確認される。
- 静的調査で言えることと実機確認が必要なことが区別される。

## 受け入れ条件

- [x] `main` のcentered float方式の復元範囲と理論分解能を確認する。
- [x] `_12` のmagnitude/sign方式の復元範囲と理論分解能を確認する。
- [x] Prefab側のContact Receiver半径・offsetとClipForVRChat側の復元設定の対応を確認する。
- [x] 実機未確認の精度要因を整理する。

## 調査結果

- `AvatarBeacon_main.prefab` はmagnitude用Contact Receiverを各軸 `radius: 2`、`position: +1` にして `0.0..1.0` のcentered floatを作る。ClipForVRChat側では `raw * 2 - 1` で符号込み値へ戻し、`PositionScale = 1000` を掛けるため、復元範囲は各軸 `-1000..+1000`。
- `AvatarBeacon_12.prefab` はmagnitude用Contact Receiverを各軸 `radius: 1`、sign用を `radius: 0.5` / `position: +0.5` にして、絶対値と符号を分ける。ClipForVRChat側では `(1 - magnitude) * 1000` にsignを適用するため、復元範囲は同じく各軸 `-1000..+1000`。
- 同じraw float精度が出る前提では、`_12` は1floatを絶対値 `0..1000` に使えるため、符号込み `-1000..+1000` を1floatへ詰める `main` より理論分解能が約2倍細かい。
- OSC/Animatorローカルfloatを32-bit floatとして見積もると、raw値1 ULPはおおむね `5.96e-8..1.19e-7`。座標換算では `_12` が約 `0.06..0.12mm`、`main` が約 `0.12..0.24mm`。
- ただし、Contact Proximity内部計算、Animator更新周期、OSC送信タイミング、Constraint評価順による実効精度は静的調査では確定できない。Unity/VRChat実機で固定座標を複数点に置き、OSC raw値と復元値の誤差分布を取る必要がある。

## 反映

- `docs/avatarbeacon-spec.md` に `main` と `_12` のfloat精度比較を追記した。
