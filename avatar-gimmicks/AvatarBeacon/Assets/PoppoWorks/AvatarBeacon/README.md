# AvatarBeacon

`AvatarBeacon` は、VRChat アバターからOSC Avatar Parametersへ、アバター基準の位置と向きの情報を送るための汎用ギミック元ファイルです。
ClipForVRChat は `coord/*` と `forward/*` 形式の受信側ツールのひとつとして利用できます。

## 前提

- Unity 2022.3 系
- VRChat SDK3 Avatar
- Modular Avatar
- Avatar Dynamics の Contact / Constraint 系機能

## 配置

Unity プロジェクトへ import すると、配布先は `Assets/PoppoWorks/AvatarBeacon` になります。
導入用 Prefab は `Assets/PoppoWorks/AvatarBeacon/Prefabs/AvatarBeacon.prefab` です。
Prefab内GameObjectの役割と削除判断は、リポジトリ側の `docs/avatarbeacon-spec.md` に記録しています。
座標測定の基準GameObjectは `WorldOriginAnchor` です。これは原点基準のContact Receiver群であり、手動で移動・回転させないでください。

## 使い方

1. アバターの root 配下に `AvatarBeacon.prefab` を配置します。
2. Modular Avatar の Bone Proxy target は `AvatarBeacon/point` を Hips 基準、`AvatarBeacon/HeadForwardAnchor` を Head 基準で割り当てます。
   - これは player root そのものではありません。
   - `point` は位置用で、Headよりプレイヤー位置に近いbasisとして扱います。
   - `HeadForwardAnchor` はyaw/forward用で、顔の向きに近いbasisとして扱います。
3. 必要に応じて、位置対象または向き対象の Transform を差し替えます。

## 送信 parameter

既定の basis 用 OSC parameter は次の 12 個です。

- `coord/x`
- `coord/xSign`
- `coord/y`
- `coord/ySign`
- `coord/z`
- `coord/zSign`
- `forward/x`
- `forward/xSign`
- `forward/y`
- `forward/ySign`
- `forward/z`
- `forward/zSign`

## 前提コスト

- basis 用に Expression Parameter 枠を 12 個使います。
- Contact / Constraint / Animator の追加分だけ、Avatar Performance Rank に影響します。
- Modular Avatar と VRChat SDK の依存が必要です。

## ClipForVRChat での確認

1. VRChat を起動し、このギミックを入れたアバターを選びます。
2. ClipForVRChat 側で avatar OSC basis の受信状態を確認します。
3. 新鮮な `coord/*` と `forward/*` の値が入った状態で `player_local` 撮影を行い、前後左右移動はHips基準、yaw回転はHead基準に追従することを確認します。
4. 値が古い、欠落している、または別アバターへ切り替えた場合は、追従しないことを確認します。

OSCが届かない場合は、VRChatの `Options > OSC > Reset OSC Config` を実行し、avatar IDごとのOSC config JSONに `coord/*` と `forward/*` の `output.address` があるか確認してください。
ClipForVRChatのログsummaryに `coord/*` または `forward/*` が出るか確認すると、OSC送信経路とAvatarBeaconのbasis parameter出力を切り分けやすくなります。
Avatar Dynamics Contact / Avatar Interactions が無効な場合も値が変化しない可能性があります。

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
- `ATG/p/*` を `coord/*` に置換
- `ATG/r/*` を `forward/*` に置換
- ClipForVRChatのbasis復元に使わない `ATG/SaveObject` / menu item を削除
- position用の既定 Bone Proxy target を Head から Hips へ変更
- yaw/forward用の `HeadForwardAnchor` を追加し、Head基準の向きを `forward/*` へ出す構成に変更
- `point` の可視化だけに使われていた `arrow` mesh / material を削除
- basis復元に使わない `AvatarBeacon Debug` / `Debug OSC Ping` menu item を削除
- 座標基準GameObject名を `WorldC` から `WorldOriginAnchor` に変更
- 誤差に見える微小なTransform値を正規化
- 既定送信名を ClipForVRChat 側の受信実装に合わせて整理
