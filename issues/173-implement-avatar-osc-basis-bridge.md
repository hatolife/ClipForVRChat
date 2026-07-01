# 専用アバターギミックOSCでhead/avatar基準Poseを自動取得する

## 問題

`player_local` は現在、手動保存した基準Poseを使うため、プレイヤー移動や回転に自動追従できない。
標準VRChat OSCだけではローカルプレイヤーroot位置/Yawを直接取得できないが、YL-ATGの調査により、Avatar Dynamics Contact/Constraintを使った専用アバターギミックからworld座標相当をOSC Avatar Parametersとして外部へ送れる可能性がある。

## 期待する挙動

専用アバターギミックを導入したユーザーは、ClipForVRChatがVRChatから送られるOSC Avatar Parametersを受信し、Headまたは任意avatar transform由来のworld position/forwardを `player_local` のbasisとして自動更新できる。
通常の標準OSC利用者には影響せず、専用ギミック未導入時は従来の手動基準Poseまたはworld構図へ安全にフォールバックできる。

## 受け入れ条件

- [ ] 自動撮影設定に、player basis sourceとして `manual` / `avatar_osc` を選べる設定を追加する。
- [ ] `avatar_osc` 有効時、VRChat OSC送信ポートから `/avatar/parameters/...` のbasis用parameterを受信できる。
- [ ] 受信したpositionとforward/yawを、既存の `AutoCapture.PlayerLocal.BasisPose` 相当へ変換できる。
- [ ] basis値には鮮度を持たせ、古い値では撮影せず、理由をUI/診断ログ/テスト撮影結果へ表示する。
- [ ] 専用アバターギミック未導入、OSC無効、parameter不足、値欠落、鮮度切れ、異常値を区別して診断できる。
- [ ] `player_local` 構図は `avatar_osc` basisが新鮮な場合にプレイヤー移動/回転へ追従して撮影できる。
- [ ] 受信basisとresolved world poseをsidecar/埋め込みmetadataへ区別して記録できる。
- [ ] 専用ギミックが必要な機能であること、通常の標準OSCだけでは動かないことをREADME/SPEC/UIへ明記する。
- [ ] Windows実機で、専用アバター導入済みのVRChatからposition/yawを受け、正面/背後/斜め構図の追従撮影を確認する。

## 実装に必要な情報

### 参照調査

- [#139](139-investigate-player-basis-source.md)
- [#172](172-clarify-or-redesign-player-local-basis-behavior.md)
- YL-ATG: https://github.com/YozoraKurage/YL-ATG

### 推奨スコープ

v0.1.8本線ではなく、v0.1.9以降の実験機能または高度機能として扱う。
v0.1.8では、初期構図とUIを誤解しにくくする短期修正を優先する。

### 用語

- `manual basis`: 現在の「現在Poseをプレイヤー基準に保存」で保持する手動基準Pose。
- `avatar OSC basis`: 専用アバターギミックからOSC Avatar Parametersとして受けるHead/avatar transform基準Pose。
- `head basis`: Head付近を原点にするbasis。player rootではない。
- `avatar basis`: Head以外の任意avatar transformを原点にするbasis。

### OSC受信parameter案

YL-ATG互換を優先する場合:

- position magnitude: `/avatar/parameters/ATG/p/x`
- position sign: `/avatar/parameters/ATG/p/x+`
- position magnitude: `/avatar/parameters/ATG/p/y`
- position sign: `/avatar/parameters/ATG/p/y+`
- position magnitude: `/avatar/parameters/ATG/p/z`
- position sign: `/avatar/parameters/ATG/p/z+`
- rotation vector magnitude: `/avatar/parameters/ATG/r/x`
- rotation vector sign: `/avatar/parameters/ATG/r/x+`
- rotation vector magnitude: `/avatar/parameters/ATG/r/y`
- rotation vector sign: `/avatar/parameters/ATG/r/y+`
- rotation vector magnitude: `/avatar/parameters/ATG/r/z`
- rotation vector sign: `/avatar/parameters/ATG/r/z+`

ClipForVRChat専用に新規設計する場合は、prefixを衝突しにくい名前にする。

- `/avatar/parameters/CFVRC/basis/p/x`
- `/avatar/parameters/CFVRC/basis/p/xSign`
- `/avatar/parameters/CFVRC/basis/p/y`
- `/avatar/parameters/CFVRC/basis/p/ySign`
- `/avatar/parameters/CFVRC/basis/p/z`
- `/avatar/parameters/CFVRC/basis/p/zSign`
- `/avatar/parameters/CFVRC/basis/f/x`
- `/avatar/parameters/CFVRC/basis/f/xSign`
- `/avatar/parameters/CFVRC/basis/f/y`
- `/avatar/parameters/CFVRC/basis/f/ySign`
- `/avatar/parameters/CFVRC/basis/f/z`
- `/avatar/parameters/CFVRC/basis/f/zSign`

`ATG/*` 互換と `CFVRC/*` 専用の両方を設定で選べると検証しやすい。

### 復元ロジック案

YL-ATG互換値は、Float parameterの0..1値と符号flagを組み合わせて座標を復元する。
調査時点のYL-ATGでは、受信値から約1000倍した絶対値を作り、`*+` parameterで正負を決めていた。
ClipForVRChat側では、この変換を設定可能にしておく。

- `scale`: 既定 `1000`
- `invertMagnitude`: 既定 `true`
- `positiveFlagThreshold`: 既定 `0`
- `maxAbsPosition`: 異常値除外用
- `freshnessSec`: 既定 `1` から `3`

forward/yawは、受信したforward vector相当の水平成分から算出する。

- `yaw = atan2(forward.x, forward.z)` を基本候補にする。
- 座標軸と符号は実機で確認し、設定または実装修正できるようにする。
- pitch/rollは当初使わず、`player_local` 変換にはYawだけを使う。

### Go側実装箇所候補

- `src/internal/appcore/config.go`
  - `AutoCapturePlayerLocalConfig` にbasis sourceとavatar OSC設定を追加する。
  - 既存configの後方互換Normalizeを追加する。
- `src/internal/appcore/player_local.go`
  - `AvatarOSCBasis` から `CameraPoseConfig` へ変換する関数を追加する。
  - Yaw算出、異常値除外、鮮度判定をテスト可能に分離する。
- `src/app.go`
  - 既存のOSC pose receiverとは別に、avatar basis receiverを管理する。
  - 受信状態を設定画面へ返すAPIを追加する。
  - `ResolveCameraViewPose` に渡すbasisを、source設定に応じてmanual/avatar OSCから選ぶ。
- `src/frontend/src/main.js`
  - basis source選択、受信状態、最終受信時刻、position/yaw、診断エラーを表示する。
  - 専用アバターギミックが必要なことを短く明示する。

### テスト方針

- 受信OSC addressからaxis/signへ正しくマップできる。
- `scale`、`invertMagnitude`、sign flagの組み合わせでpositionを復元できる。
- forward vectorからYawを算出できる。
- stale/partial/invalid basisでは撮影に使わない。
- manual basis設定は既存挙動を壊さない。
- `avatar_osc` sourceが無効な環境でも通常起動・通常撮影に影響しない。

### 専用アバターギミック側の必要条件

- Modular Avatarで導入できるprefabにする。
- 既定の追跡対象はHead相当とし、必要に応じて対象Bone/Transformを変更可能にする。
- Avatar Dynamics Contact/Constraintでpositionとforward vectorをFloat parametersへエンコードする。
- Parameter prefixは `CFVRC/` を既定にし、YL-ATG互換も検証用に受けられるようにする。
- Expression Parameter枠、Contact数、Performance Rankへの影響をREADMEへ書く。
- 既存アバターギミックとのparameter/tag衝突を避ける命名にする。

### UI/ドキュメントで明記する制限

- 標準OSCだけでは使えない。
- 専用アバターギミックを導入したアバター使用中だけ有効。
- デフォルトはHead基準であり、player root基準ではない。
- アバター切り替え、OSC無効化、parameter出力停止、値の鮮度切れで自動追従は停止する。
- 座標精度、遅延、更新頻度はVRChat/Avatar Dynamics/OSCの挙動に依存する。

### 実機確認項目

- 専用アバター導入後、ClipForVRChatにposition/yawが表示される。
- プレイヤーが前後左右へ移動するとpositionが追従する。
- プレイヤーがYaw回転するとbasis yawが追従する。
- `player_local` 正面/背後/斜め構図が、移動後も概ね同じ相対構図で撮影される。
- アバターを未導入アバターへ切り替えると、鮮度切れで安全に失敗する。
- OSCを無効にすると、診断表示が分かる。
- Stream/Spout方式とPhoto方式の両方でbasis解決が同じになる。

## 非対象

- VRChat client改造、メモリ読み、非公式APIによる位置取得。
- 任意ワールドのUdonなしでplayer rootを完全自動取得すること。
- Quest単体で外部Windowsアプリなしに完結すること。
