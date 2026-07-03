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
READMEの手順では Modular Avatar の Bone Proxy target を `AvatarBeacon/point` に設定し、既定は Hips 相当として扱います。

これは player root そのものではありませんが、Headよりプレイヤー位置に近いbasisとして扱います。
ClipForVRChat 側では `coord/*` から復元した位置を `player_local` basis の位置、`forward/*` から復元した水平forward vectorを basis yaw として使います。

`point` は単なる見た目用オブジェクトではありません。
MA Bone Proxyの付いた追跡アンカーであり、`WorldOriginAnchor`、`const_x`、`const_z` などのConstraintがこのTransformを参照して座標エンコードの基準にします。
そのため `point` 自体は必要です。

## OSC Parameter

VRChatから外部へ送られるOSC addressは `/avatar/parameters/` が付いた形になります。
Prefab内のparameter名はその下の相対名です。

| 用途 | Avatar parameter | OSC address | 詳細 |
| --- | --- | --- | --- |
| position X magnitude | `coord/x` | `/avatar/parameters/coord/x` | 追跡対象 `point` のworld X座標の絶対値側。ClipForVRChatは既定で `1 - value` に `1000` を掛けて復元する。 |
| position X sign | `coord/xSign` | `/avatar/parameters/coord/xSign` | X座標の符号。既定では `0` 以下を負、`0` より大きい値を正として扱う。 |
| position Y magnitude | `coord/y` | `/avatar/parameters/coord/y` | 追跡対象 `point` のworld Y座標の絶対値側。高さ方向のbasis位置に使う。 |
| position Y sign | `coord/ySign` | `/avatar/parameters/coord/ySign` | Y座標の符号。復元時に `coord/y` と組み合わせる。 |
| position Z magnitude | `coord/z` | `/avatar/parameters/coord/z` | 追跡対象 `point` のworld Z座標の絶対値側。前後方向のbasis位置に使う。 |
| position Z sign | `coord/zSign` | `/avatar/parameters/coord/zSign` | Z座標の符号。復元時に `coord/z` と組み合わせる。 |
| forward X magnitude | `forward/x` | `/avatar/parameters/forward/x` | 追跡対象の向きから作るforward vectorのX成分の絶対値側。yaw算出に使う。 |
| forward X sign | `forward/xSign` | `/avatar/parameters/forward/xSign` | forward X成分の符号。yaw算出前に `forward/x` と組み合わせる。 |
| forward Y magnitude | `forward/y` | `/avatar/parameters/forward/y` | forward vectorのY成分の絶対値側。ClipForVRChatでは受信完全性の確認に使い、yaw計算には直接使わない。 |
| forward Y sign | `forward/ySign` | `/avatar/parameters/forward/ySign` | forward Y成分の符号。`forward/y` と組み合わせて完全なforward vectorとして復元できるように残す。 |
| forward Z magnitude | `forward/z` | `/avatar/parameters/forward/z` | forward vectorのZ成分の絶対値側。yaw算出に使う。 |
| forward Z sign | `forward/zSign` | `/avatar/parameters/forward/zSign` | forward Z成分の符号。yaw算出前に `forward/z` と組み合わせる。 |

`coord/*` と `forward/*` のmagnitudeは float、`*Sign` はbool相当の制御値として扱います。
magnitudeだけでは正負が分からないため、各軸は必ずmagnitude parameterとsign parameterを組み合わせて復元します。
現Prefabでは Modular Avatar Parameters に全12個を `localOnly` として登録しています。

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

## YL-ATG の座標取得方式

YL-ATG は、VRChat clientやワールドAPIからプレイヤー座標を直接取得しているわけではありません。
アバター内に配置した Transform、Unity Constraint、VRChat Avatar Dynamics Contact を組み合わせ、追跡対象Transformの位置と向きを Avatar Expression Parameter に変換します。
VRChat のOSC Avatar ParametersはそのExpression Parameterを外部へ送るため、ClipForVRChatはOSC経由で値を受け取れます。

前提になるVRChat公式仕様は次の通りです。

- Contact Sender と Contact Receiver は、同じcollision tagを持つと接触として扱われる。
- アバター上の Contact Receiver は、接触結果をAnimator Parameterへ書き込める。
- Receiver Type `Proximity` は、SenderがReceiver中心にどれだけ近いかを `0.0` から `1.0` のfloatとして出す。
- Expression Parameterに登録されたAvatar parameterは、OSC有効時に `/avatar/parameters/<name>` として送信される。

YL-ATG のPrefabでは、これを利用して次のような変換をしています。

1. `point` を MA Bone Proxy でHeadなどのアバターBoneへ追従させる。
2. `WorldOriginAnchor` を座標の受け側にし、`ATG/p/x`、`ATG/p/y`、`ATG/p/z` の Proximity Receiver を置く。
3. `const_x`、`const_y`、`const_z` を `point` と `WorldOriginAnchor` の間にConstraintで配置し、対応するContact Sender tagを出す。
4. ReceiverのProximity値を `ATG/p/*` のmagnitudeとして使う。
5. 同じ軸の `ATG/p/x+`、`ATG/p/y+`、`ATG/p/z+` を符号判定用parameterとして使う。
6. `rot`、`offset_rot`、`get_rot`、`Rot_flag` で同じ仕組みを回転方向に適用し、`ATG/r/*` と `ATG/r/*+` を出す。

つまり、座標そのものを文字列やOSC packetとしてアバターから任意送信しているのではなく、Contact ReceiverのProximity floatを「距離センサー」として使い、座標の絶対値と符号へ分解しています。
値のスケールは `value = (1 - magnitude) * 1000` として外部側で復元します。
このため、絶対値が小さい座標はmagnitudeが1に近く、絶対値が大きい座標はmagnitudeが0に近い値として表現され、符号は別parameterで補います。

AvatarBeacon では、この仕組み自体は維持しつつ、公開parameter名だけを次のように置き換えます。

| YL-ATG | AvatarBeacon | 意味 |
| --- | --- | --- |
| `ATG/p/x` | `coord/x` | X座標magnitude |
| `ATG/p/x+` | `coord/xSign` | X座標符号 |
| `ATG/p/y` | `coord/y` | Y座標magnitude |
| `ATG/p/y+` | `coord/ySign` | Y座標符号 |
| `ATG/p/z` | `coord/z` | Z座標magnitude |
| `ATG/p/z+` | `coord/zSign` | Z座標符号 |
| `ATG/r/x` | `forward/x` | forward X magnitude |
| `ATG/r/x+` | `forward/xSign` | forward X符号 |
| `ATG/r/y` | `forward/y` | forward Y magnitude |
| `ATG/r/y+` | `forward/ySign` | forward Y符号 |
| `ATG/r/z` | `forward/z` | forward Z magnitude |
| `ATG/r/z+` | `forward/zSign` | forward Z符号 |

参考:

- VRChat Contacts: https://creators.vrchat.com/common-components/contacts/
- VRChat OSC Avatar Parameters: https://docs.vrchat.com/docs/osc-avatar-parameters

## Prefab 構造

静的に確認できるGameObject構造は次の通りです。

```text
AvatarBeacon
├── WorldOriginAnchor
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
```

### Root `AvatarBeacon`

Prefab rootです。
Modular Avatar Parameters相当のComponentで、`coord/*` と `forward/*` をAvatar Expression Parameterへ登録します。

このrootを削るとparameter登録とPrefab導入単位が壊れるため必須です。

### `point`

追跡対象の基準Transformです。
Modular Avatar Bone Proxy相当のComponentを持ち、既定ではHipsへ追従します。
必要に応じてHeadなど任意のBone/Transformへ差し替えられます。

このオブジェクトがないと、どのアバターTransformを座標出力対象にするかを指定できないため必須です。

### `WorldOriginAnchor`

position出力の中心です。
`coord/x,y,z` と `coord/xSign,ySign,zSign` のContact receiver相当Componentを持ち、`point` 由来の位置を magnitude/sign に分けてAvatar Parameterへ書き込みます。
Transformは原点、無回転、scale 1 を保つ前提です。

このオブジェクトを削ると `coord/*` が出なくなるため必須です。

### `const_x` / `const_y` / `const_z`

各軸の基準点を作る補助Transformです。
ParentConstraint と Contact sender相当Componentを持ち、`WorldOriginAnchor` 側のContact receiverへ衝突タグを渡す役割と推定しています。

単体では出力parameterを直接持ちませんが、`WorldOriginAnchor` の軸magnitude/sign算出に使われるため、Unity/VRChat実機で代替確認するまでは削除しません。

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

### 削除したデバッグメニュー

YL-ATG由来の `ATG/SaveObject` は、Prefab上ではMA Menu Itemで保存/制御用parameterを操作するだけで、`coord/*` / `forward/*` のContact経路やClipForVRChatのbasis復元には使われていません。
AvatarBeaconでは用途が曖昧な `SaveObject` を残さず、basis復元に使わない `AvatarBeacon Debug` / `Debug OSC Ping` も削除します。
OSC疎通確認は、ClipForVRChatの `raw` / `last` と10秒summaryに `coord/*` / `forward/*` が出るかで行います。

### `arrow` mesh

YL-ATG由来の `arrow` mesh / material は、`point` の位置と向きをUnity Scene上で見やすくする可視化用の子Prefabでした。
Contact sender/receiver、MA Bone Proxy、ParentConstraintの参照先ではなく、OSC parameter出力にも直接関与しません。

AvatarBeaconでは配布物を最小化するため、`arrow` のPrefabInstanceと `FBXs/arrow.*`、関連materialを削除します。
`point` は残しますが、mesh表示は持ちません。

## 必要性の判断

現状のAvatarBeaconは、YL-ATG_ForAvatarを大きく簡略化した新規実装ではなく、YL-ATGのPrefab/FBX/Material/構成を改変した配布物です。
静的監査では、ClipForVRChat専用名のparameter、YozoLab import path、不要なVRCSDK/Modular Avatar本体の同梱はありません。

ただし、Contact/ConstraintグラフはGameObject間の相互作用で成立するため、Unity/VRChat実機確認なしに「見た目だけ」「不要そう」に見えるGameObjectを削るのは危険です。
v0.1.8では、不要機能を追加しないことよりも、YL-ATGで成立しているグラフを壊さずに汎用parameter名・ライセンス・配布導線を整えることを優先します。

`arrow` meshと用途不明な `SaveObject` は削除または置換済みです。
残っているGameObjectは、静的監査上はparameter登録、OSC疎通デバッグ、position/forwardのContact/Constraint経路、追跡アンカーのいずれかに分類できます。

## 実機確認項目

- Unity import後、`Assets/PoppoWorks/AvatarBeacon/...` として配置される。
- `AvatarBeacon.prefab` をアバターroot配下へ置ける。
- `point` をHips基準に設定できる。
- ClipForVRChatの `avatar_osc` raw/lastまたは10秒summaryに `coord/*` / `forward/*` が表示される。
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
