# AvatarBeacon 詳細仕様

## 目的

AvatarBeacon は、VRChat アバター上の追跡対象 Transform から、外部ツールが受信できる位置と向きの情報を OSC Avatar Parameters として出すアバターギミックです。
ClipForVRChat 専用ではなく、同じ parameter を読める外部ツールでも使えるように、公開OSC parameterに ClipForVRChat 固有名を含めません。

このギミックは YozoraKurage/YL-ATG `ATG_ForAvatar_V0.0.3` を元に、import path、Prefab名、parameter名、配布説明を変更したものです。
YL-ATG由来部分は `avatar-gimmicks/AvatarBeacon/Assets/PoppoWorks/AvatarBeacon/LICENSES/YL-ATG-MIT.txt` と `NOTICE.md` に記録します。

## 依存

- Unity 2022.3 系
- VRChat SDK3 Avatar
- Modular Avatar
- VRChat Avatar Dynamics Contact / Constraint

CIで作る `AvatarBeacon-*-source.zip` には VRCSDK と Modular Avatar 本体を含めません。
Unity側で `.unitypackage` を手作業作成するときも、外部依存本体は同梱しません。

## 追跡対象

Prefab root直下の `point` が追跡対象です。
READMEの手順では Modular Avatar の Bone Proxy target を `AvatarBeacon/point` に設定し、既定は Head 相当として扱います。

これは player root ではありません。
ClipForVRChat 側では `coord/*` から復元した位置を `player_local` basis の位置、`forward/*` から復元した水平forward vectorを basis yaw として使います。

## OSC Parameter

VRChatから外部へ送られるOSC addressは `/avatar/parameters/` が付いた形になります。
Prefab内のparameter名はその下の相対名です。

| 用途 | Avatar parameter | OSC address |
| --- | --- | --- |
| position X magnitude | `coord/x` | `/avatar/parameters/coord/x` |
| position X sign | `coord/xSign` | `/avatar/parameters/coord/xSign` |
| position Y magnitude | `coord/y` | `/avatar/parameters/coord/y` |
| position Y sign | `coord/ySign` | `/avatar/parameters/coord/ySign` |
| position Z magnitude | `coord/z` | `/avatar/parameters/coord/z` |
| position Z sign | `coord/zSign` | `/avatar/parameters/coord/zSign` |
| forward X magnitude | `forward/x` | `/avatar/parameters/forward/x` |
| forward X sign | `forward/xSign` | `/avatar/parameters/forward/xSign` |
| forward Y magnitude | `forward/y` | `/avatar/parameters/forward/y` |
| forward Y sign | `forward/ySign` | `/avatar/parameters/forward/ySign` |
| forward Z magnitude | `forward/z` | `/avatar/parameters/forward/z` |
| forward Z sign | `forward/zSign` | `/avatar/parameters/forward/zSign` |
| internal save/control | `avatar_beacon/save` | `/avatar/parameters/avatar_beacon/save` |

`coord/*` と `forward/*` は float、`*Sign` と `avatar_beacon/save` はbool相当の制御値として扱います。
現Prefabでは Modular Avatar Parameters に全13個を `localOnly` として登録しています。

## 値の復元

AvatarBeacon は YL-ATG のContact方式を引き継いでおり、各軸を magnitude と sign に分けて出します。
ClipForVRChat 側の既定復元は次の通りです。

- `value = (1 - magnitude) * 1000`
- `sign > 0` なら正、それ以外なら負
- positionは `coord/x,y,z` を使う
- yawは `forward/x,z` から `atan2(forward.x, forward.z)` で求める
- `forward/y` は受信完全性の確認には使うが、`player_local` のyaw計算には直接使わない

`ATG/*` 互換受信は、既存YL-ATGとの切り分け用にClipForVRChat側へ残します。
AvatarBeaconの既定出力は `ATG/*` ではありません。

## Prefab 構造

静的に確認できるGameObject構造は次の通りです。

```text
AvatarBeacon
├── WorldC
│   ├── const_x
│   ├── const_y
│   ├── const_z
│   └── rot
│       ├── offset_rot
│       │   └── X
│       ├── get_rot
│       │   ├── X
│       │   ├── Y
│       │   └── Z
│       └── Rot_flag
├── point
└── AvatarBeacon
    └── SaveObject
```

### Root `AvatarBeacon`

Prefab rootです。
Modular Avatar Parameters相当のComponentで、`coord/*`、`forward/*`、`avatar_beacon/save` をAvatar Expression Parameterへ登録します。
同じroot上に、Modular Avatar系のmenu/install補助と思われるComponentも付いています。

このrootを削るとparameter登録とPrefab導入単位が壊れるため必須です。

### `point`

追跡対象の基準Transformです。
Modular Avatar Bone Proxy相当のComponentを持ち、Headなど任意のBone/Transformへ追従させるための入口です。

このオブジェクトがないと、どのアバターTransformを座標出力対象にするかを指定できないため必須です。

### `WorldC`

position出力の中心です。
`coord/x,y,z` と `coord/xSign,ySign,zSign` のContact receiver相当Componentを持ち、`point` 由来の位置を magnitude/sign に分けてAvatar Parameterへ書き込みます。

このオブジェクトを削ると `coord/*` が出なくなるため必須です。

### `const_x` / `const_y` / `const_z`

各軸の基準点を作る補助Transformです。
ParentConstraint と Contact sender相当Componentを持ち、`WorldC` 側のContact receiverへ衝突タグを渡す役割と推定しています。

単体では出力parameterを直接持ちませんが、`WorldC` の軸magnitude/sign算出に使われるため、Unity/VRChat実機で代替確認するまでは削除しません。

### `rot`

forward vector出力の中心です。
`forward/x,y,z` のContact receiver相当Componentを持ちます。
ClipForVRChat はこのうち水平成分 `forward/x,z` からyawを復元します。

このオブジェクトを削ると `forward/*` が出なくなるため必須です。

### `get_rot` / `get_rot/X,Y,Z`

forward vectorの各軸を取り出す補助階層です。
各軸子要素は Contact sender相当Componentを持ち、`rot` の `forward/*` receiverへ入力する構成と推定しています。

Yaw追従の主要経路なので削除しません。

### `Rot_flag`

`forward/xSign,ySign,zSign` のContact receiver相当Componentを持つ符号出力用オブジェクトです。
forward vectorの各軸が正方向か負方向かを外部で復元するために必要です。

ClipForVRChatのyaw計算は `forward/x,z` の符号復元に依存するため必須です。

### `offset_rot` / `offset_rot/X`

回転/forward算出のためのオフセット補助階層です。
Prefab YAML上は直接parameter名を持ちませんが、`rot`、`get_rot`、`Rot_flag` のContact/Constraint構成に組み込まれています。

YL-ATG由来の座標変換グラフの一部であり、削るとforwardの符号や向きが変わる可能性があるため、実機検証前には削除しません。

### 子 `AvatarBeacon` / `SaveObject`

YL-ATG由来の保存/制御用補助階層です。
`avatar_beacon/save` を使うメニュー制御相当Componentを持ちます。
現時点では、元Prefabと同じ制御グラフを保つため残しています。

Unity実機で `avatar_beacon/save` が不要と確認できた場合は、将来の削除候補です。
削除する場合は、Modular Avatar Parametersから `avatar_beacon/save` を外し、メニュー/制御Componentを同時に整理する必要があります。

### `FBXs/arrow.*` とMaterial

`FBXs/arrow.prefab` は main Prefab 内で `point` 配下にPrefabInstanceとして参照されています。
視覚的な軸表示だけに見える可能性がありますが、追跡対象の子TransformやConstraint参照に関わっている可能性があります。

Unity上で非表示化/削除してOSC出力が変わらないことを確認するまでは、配布物から削除しません。

## 必要性の判断

現状のAvatarBeaconは、YL-ATG_ForAvatarを大きく簡略化した新規実装ではなく、YL-ATGのPrefab/FBX/Material/構成を改変した配布物です。
静的監査では、ClipForVRChat専用名のparameter、YozoLab import path、不要なVRCSDK/Modular Avatar本体の同梱はありません。

ただし、Contact/ConstraintグラフはGameObject間の相互作用で成立するため、Unity/VRChat実機確認なしに「見た目だけ」「不要そう」に見えるGameObjectを削るのは危険です。
v0.1.8では、不要機能を追加しないことよりも、YL-ATGで成立しているグラフを壊さずに汎用parameter名・ライセンス・配布導線を整えることを優先します。

削除候補として扱うのは次の2点です。

- `avatar_beacon/save` と子 `AvatarBeacon/SaveObject`
- `FBXs/arrow.*` と関連Material

これらは、Unity上で削除後も `coord/*` と `forward/*` が正しく出続けることを確認できた場合だけ削除します。

## 実機確認項目

- Unity import後、`Assets/PoppoWorks/AvatarBeacon/...` として配置される。
- `AvatarBeacon.prefab` をアバターroot配下へ置ける。
- `point` をHead基準に設定できる。
- VRChat OSCで `/avatar/parameters/coord/*` と `/avatar/parameters/forward/*` が更新される。
- 前後左右移動で `coord/*` が変化する。
- yaw回転で `forward/*` が変化する。
- アバター切り替えやOSC停止時にClipForVRChat側が鮮度切れとして扱う。

## OSC送信されない場合の確認

`v0.1.8-rc16` 実機確認では、AvatarBeacon導入時にOSC Avatar Parametersが送信されていないように見える報告があった。
この場合、次を優先して切り分ける。

- VRChatのOSCが有効で、送信先がClipForVRChatの受信ポートに合っているか。
- avatar IDごとのOSC config JSONがAvatarBeacon導入後のparameterを含んでいるか。
- `coord/*` と `forward/*` に `output.address` があるか。
- VRChat Action Menuの `Reset OSC Config` 後にアバターを再読み込みしても変わらないか。
- Avatar Dynamics Contact / Avatar Interactions が有効か。
- ClipForVRChatの `avatar_osc` raw受信件数が0か、他parameterだけ届いているか。

raw受信件数が0なら、Prefabの復元ロジック以前にVRChatからOSC packetが出ていない。
raw受信件数があるのに `coord/*` / `forward/*` が欠けるなら、OSC configまたはPrefab parameter名の不一致を疑う。
