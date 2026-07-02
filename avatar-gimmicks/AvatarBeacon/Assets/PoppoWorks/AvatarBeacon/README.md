# AvatarBeacon

`AvatarBeacon` は、VRChat アバターからOSC Avatar Parametersへ、アバター基準の位置と向きの情報を送るための汎用ギミック元ファイルです。
ClipForVRChat は `CFVRC/basis` 形式の受信側ツールのひとつとして利用できます。

## 前提

- Unity 2022.3 系
- VRChat SDK3 Avatar
- Modular Avatar
- Avatar Dynamics の Contact / Constraint 系機能

## 配置

Unity プロジェクトへ import すると、配布先は `Assets/PoppoWorks/AvatarBeacon` になります。
導入用 Prefab は `Assets/PoppoWorks/AvatarBeacon/Prefabs/AvatarBeacon.prefab` です。

## 使い方

1. アバターの root 配下に `AvatarBeacon.prefab` を配置します。
2. Modular Avatar の Bone Proxy target は `AvatarBeacon/point` を Head 基準で割り当てます。
   - これは player root 基準ではありません。
   - 既定の追跡対象は Head 相当です。
3. 必要に応じて、対象 Transform を Head 以外へ差し替えます。

## 送信 parameter

既定の OSC parameter は次の 12 個です。

- `CFVRC/basis/p/x`
- `CFVRC/basis/p/xSign`
- `CFVRC/basis/p/y`
- `CFVRC/basis/p/ySign`
- `CFVRC/basis/p/z`
- `CFVRC/basis/p/zSign`
- `CFVRC/basis/f/x`
- `CFVRC/basis/f/xSign`
- `CFVRC/basis/f/y`
- `CFVRC/basis/f/ySign`
- `CFVRC/basis/f/z`
- `CFVRC/basis/f/zSign`

`SaveObject` は配布元 YL-ATG の補助制御を `AvatarBeacon/SaveObject` に置き換えたものです。

## 前提コスト

- Expression Parameter 枠を 12 個使います。`SaveObject` を残す構成では追加の保存用 parameter が 1 個増えます。
- Contact / Constraint / Animator の追加分だけ、Avatar Performance Rank に影響します。
- Modular Avatar と VRChat SDK の依存が必要です。

## ClipForVRChat での確認

1. VRChat を起動し、このギミックを入れたアバターを選びます。
2. ClipForVRChat 側で avatar OSC basis の受信状態を確認します。
3. 新鮮な `CFVRC/basis` 値が入った状態で `player_local` 撮影を行い、前後左右移動と yaw 回転に追従することを確認します。
4. 値が古い、欠落している、または別アバターへ切り替えた場合は、追従しないことを確認します。

## 検証限界

この作業環境では Unity Editor / VRChat 実機を起動しての import と動作確認までは行っていません。
そのため、Prefab の完全な動作確認、Modular Avatar の実適用、OSC の実送受信は Unity/VRChat 側で最終確認してください。

## 手動で unitypackage にする手順

1. Unity へ import した `Assets/PoppoWorks/AvatarBeacon` を開きます。
2. 必要な Prefab と依存 asset が揃っていることを確認します。
3. `Assets > Export Package...` を選び、`Assets/PoppoWorks/AvatarBeacon` 配下を export します。
4. Export対象に `LICENSES/YL-ATG-MIT.txt` と `NOTICE.md` が含まれることを確認します。
5. VRCSDK本体、Modular Avatar本体、Unity `Library/` や `Temp/` は unitypackage に含めません。

## 由来

この asset source は YozoraKurage/YL-ATG `ATG_ForAvatar_V0.0.3` を参考にしています。
コピー・改変した範囲は以下です。

- Prefab 名と公開 import path を `AvatarBeacon` に変更
- `ATG/p/*` と `ATG/r/*` を `CFVRC/basis/*` に置換
- `ATG/SaveObject` を `AvatarBeacon/SaveObject` に置換
- 既定送信名を ClipForVRChat 側の受信実装に合わせて整理
