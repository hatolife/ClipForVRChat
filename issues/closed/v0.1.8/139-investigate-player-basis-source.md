# ローカルプレイヤー基準Poseの取得可否を実機調査する

## 問題

`player_local` を実装するには、撮影時点のローカルプレイヤー位置と向きが必要である。
現行調査では、VRChat公式OSCの `/usercamera/Pose` はカメラPoseであり、標準Avatar Parametersにもプレイヤーroot位置/yawは見当たらない。

## 期待する挙動

標準OSC/OSCQuery/VRChatログ/実機挙動から、ローカルプレイヤー基準Poseを取得できるかを確認し、取得できない場合は自動 `player_local` を未対応として扱う判断材料を残す。

## 受け入れ条件

- `/usercamera/Pose` のGet値が、Camera ModeごとにカメラPoseのみを返すのか、プレイヤー基準として使える値を返す余地があるのかを実機で確認する。
- `/usercamera/Mode=2`、`/usercamera/Streaming`、`/usercamera/LookAtMe` の組み合わせで、カメラ起動直後のPoseを基準Poseとして使えるか確認する。
- OSCQueryで公開されるendpoint一覧を取得し、ローカルプレイヤー位置/向きに相当するendpointがないか確認する。
- VRChat output logに、撮影時点で使えるローカルプレイヤー位置/向き情報が出ていないか確認する。
- 標準機能だけで取得できない場合、手動基準Pose、専用ワールド/アバター補助OSC、または未対応表示のどれを採用するか決める。
- 調査結果に基づき、`player_local` 実装を進めてよい条件と、進めてはいけない条件を明記する。

## 対応内容

- v0.1.8実装に合わせて対応済み。詳細は `feat/v0.1.8-resolve-issues` の実装、README、SPEC、RELEASE_NOTESを参照。

## 2026-07-01 再調査: 現在のプレイヤー位置取得方法

### 結論

標準VRChat OSCだけで、外部アプリが現在のローカルプレイヤーroot位置/Yawを直接読む方法は見つからない。
`/usercamera/Pose` で読めるのはUser Cameraの位置/回転であり、`/tracking/...` は外部トラッカー情報をVRChatへ入力する用途で、VRChatから現在のプレイヤー位置を出力する用途ではない。

Udonの `VRCPlayerApi` では `GetPosition()`、`GetRotation()`、`GetTrackingData()` によりプレイヤー位置、回転、head/hand等のtracking dataを取得できる。
ただしこれはワールド内スクリプトのAPIであり、ClipForVRChatのような外部Windowsアプリが、任意の通常ワールドでその値を直接読む経路ではない。

### 確認した経路

#### 1. VRChat標準OSC

- OSC Overviewでは、VRChatは既定で受信9000、送信9001を使う。
- Avatar Parameters APIは、avatar parameterを `/avatar/parameters/name` として入出力する仕組みであり、ローカルプレイヤーrootのworld position/yawを標準出力するものではない。
- User Camera OSCの `/usercamera/Pose` はCamera position & rotationであり、プレイヤーroot positionではない。
- Tracking OSCの `/tracking/trackers/...` と `/tracking/trackers/head/...` はwrite-onlyで、外部プログラムからVRChatへtracker/head位置を供給するもの。VRChatから現在のHMD/プレイヤー位置を読む用途ではない。

判断: 標準OSCのみで「現在のプレイヤー位置を取得」は不可。

#### 2. OSCQuery

OSCQueryはOSCアプリ同士が互いのcapabilityやportを発見するための仕組み。
VRChatが公開するOSC endpointの発見には使えるが、標準OSCにプレイヤー位置出力endpointが存在しない限り、位置取得そのものは解決しない。

判断: 自動検出には有用だが、player_local自動追従の根本解決にはならない。

#### 3. Udon / ワールド補助

VRChat Creator DocsのPlayer Positionsでは、Udonから以下を取得できる。

- `GetPosition()`: Playerのworld space位置。
- `GetRotation()`: Playerのworld space回転。
- `GetTrackingData()`: Local PlayerではTrackingManager由来、Remote Playerでは骨由来のHead/Hand等。`AvatarRoot` はplayer capsuleが付くavatar root transform。

このため、専用ワールドやワールド内Udon prefabを前提にできるなら、ワールド内ではローカルプレイヤー基準Poseを得られる。
しかしUdonが任意の外部WindowsアプリへUDP/OSCで直接送信できる公式経路は、この調査範囲では確認できていない。
また、ClipForVRChatを任意ワールドで使う用途では、専用ワールド導入を前提にできない。

判断: 「専用ワールド/ワールドギミック前提」なら有望。ただし汎用アプリの標準機能としては採用しづらい。

#### 4. アバター補助

Avatar Parametersは外部OSCへ値を出せるが、アバター側が自身のworld position/yawを任意に計算してparameterへ流す仕組みは標準アバター機能だけでは見当たらない。
アバターにスクリプトを持たせることも通常できない。

判断: 一般配布アプリの解決策としては弱い。専用アバターギミック前提の可能性は残るが、汎用性は低い。

##### 2026-07-02補足: 専用アバターギミックからOSCでworld座標を送れるか

結論として、「Avatar Parameterとして既に持っている値をOSCで外部へ送る」ことは可能だが、「アバター単体で自分のworld座標を取得してAvatar Parameterへ入れ、それをOSC送信する」ことは標準機能だけでは成立しない可能性が高い。

根拠:

- OSC Avatar Parametersは、Avatar Parameterを外部OSCアプリへ送信できる。configの `output.address` を設定すれば、parameter値の変化を任意OSC addressへ送れる。
- Animator Parametersのbuilt-in一覧には、`VelocityX/Y/Z`、`VelocityMagnitude`、`Upright`、`Grounded`、`TrackingType`、avatar scale系などはあるが、world positionやworld yawはない。
- custom Expression Parameterは作れるが、制御手段はExpressions Menu、Avatar Parameter Driver、OSC入力などであり、アバター自身のworld transformを数値化してparameterへ書き込むAPIではない。
- Contact Receiverはavatar上ではAnimator Parameterを更新できるが、主に接触/近接の0..1値であり、world座標そのものではない。world上のUdonではContact Senderのworld position等を読めるが、それはワールド側Udon機能であり、任意の通常ワールドでアバター単体が外部アプリへ座標を出す仕組みではない。

したがって、専用アバターを同梱するだけでは、ClipForVRChatへローカルプレイヤーworld座標を安定送信する設計にはならない。
実現性があるとすれば以下のどちらか。

1. 専用ワールド/Udon連携
   - Udonで `VRCPlayerApi.GetPosition()` / `GetRotation()` / `GetTrackingData()` を取得する。
   - ただしUdonから任意UDP/OSCで外部アプリへ直接送信する公式経路は確認できていない。別の橋渡し手段が必要。

2. 外部VRランタイム連携
   - SteamVR/OpenXR等からHMD poseを取得し、ClipForVRChat側でVRChat world座標へキャリブレーションする。
   - VR専用かつ座標合わせが必要で、Desktop/Quest単体の汎用解決にはならない。

#### 5. SteamVR/OpenVR/OpenXR等の外部VRランタイム

VRモードでは外部ランタイムからHMD poseを取れる可能性がある。
ただし得られる座標系はVRランタイムのtracking spaceであり、VRChat world座標と原点/Yaw/scaleが一致する保証はない。
Desktop mode、Quest単体、SteamVR以外の構成では使えない場合がある。
VRChat world座標へ変換するには別途キャリブレーションが必要。

判断: 高度なオプションとしては検討可能だが、v0.1.8の一般向け自動撮影の標準経路には不向き。

#### 6. User Camera LookAtMe等からの推定

`/usercamera/LookAtMe` 有効時のUser Camera rotationと、既知のcamera positionから、ローカルプレイヤー方向のrayを推定できる可能性はある。
複数のcamera positionからrayを取れば、頭部または注視点らしき位置を三角測量できる可能性もある。

ただしこれは公式に「プレイヤー位置取得」として保証されたendpointではなく、User Cameraを一時的に動かす必要がある。
撮影対象の構図を壊す、カメラ挙動差やモード差に影響される、rootではなくLookAt対象点になる、距離推定が不安定、という問題がある。

判断: 実験価値はあるが、信頼できる標準仕様として扱うべきではない。採用するなら「近似/実験的な自動基準推定」と明示する。

### 実装方針候補

1. v0.1.8対応
   - 初期構図を `world` に戻す、または `player_local` 構図を初期OFFにする。
   - `player_local` は手動基準Pose方式として、UI文言を「プレイヤー追従」ではなく「手動基準」に寄せる。
   - YL-ATG方式を参考にしたAvatarBeaconを用意し、専用アバター導入時はOSC Avatar Parametersからhead/avatar基準Poseを受信できるようにする。

2. 追加検証候補
   - User Camera `LookAtMe` とPose readbackによる近似基準推定を実機検証する。
   - OSCQueryでUser Camera/Tracking endpoint一覧を診断表示し、利用可能なendpointを可視化する。

3. 長期候補
   - 専用ワールド/Udon prefabと連携する方式を別機能として検討する。
   - VRランタイムpose取得とVRChat座標へのキャリブレーションを、VRユーザー向けの高度なオプションとして検討する。

### 参照

- VRChat OSC Overview: https://docs.vrchat.com/docs/osc-overview
- VRChat OSC Avatar Parameters: https://docs.vrchat.com/docs/osc-avatar-parameters
- VRChat Wiki OSC User Camera / Tracking: https://wiki.vrchat.com/wiki/OSC
- VRChat OSCQuery: https://docs.vrchat.com/docs/oscquery
- VRChat Creator Docs Player API: https://creators.vrchat.com/worlds/udon/players/
- VRChat Creator Docs Player Positions: https://creators.vrchat.com/worlds/udon/players/player-positions/

## 2026-07-02 追加調査: YL-ATG

`https://github.com/YozoraKurage/YL-ATG` が、アバターギミック併用で現在のプレイヤー位置を外部OSCへ出せる根拠になるか確認する。

### 調査結果

YL-ATGは、アバター上の任意オブジェクトの絶対座標をOSC経由で外部Unityプロジェクトへ同期するためのギミック/プログラム。
READMEでは、`ATG_ForAvatar.unitypackage` をアバターへ入れ、`ATG/point` のMA Bone Proxyにトラッキングしたいオブジェクトを指定する。デフォルトはHead。
外部側は `ATG_ForUnity.unitypackage` とOscCoreを入れ、Play Modeで `YL-ATG/TrackingObject` がVRChatから取得した位置になる、と説明している。

実装上の構成:

- アバター側prefabはModular Avatarで `ATG/p/x`, `ATG/p/y`, `ATG/p/z`, `ATG/p/x+`, `ATG/p/y+`, `ATG/p/z+`, `ATG/r/x`, `ATG/r/y`, `ATG/r/z`, `ATG/r/x+`, `ATG/r/y+`, `ATG/r/z+`, `ATG/SaveObject` のparametersを追加している。
- 座標値はアバター上のscriptではなく、Avatar Dynamics Contact Sender/ReceiverとConstraintを使ってFloat parameterへエンコードしている。
- `ATG_OSCHub.py` はVRChatのOSC送信ポート9001をlistenし、`/avatar/parameters/ATG/*` を `ports.json` の `9010`, `9003`, `9004` へ転送するだけ。
- 外部Unity側 `YL_ATGCore.cs` は受信値を `vrc_p_x = (1 - value) * 1000` のように復元し、`ATG/p/*+` のflagで正負を決める。復元値を `TrackingObject.transform.position` へ反映する。
- rotation系も同様に `ATG/r/*` と `ATG/r/*+` で方向ベクトル相当を復元し、`RotationVecObject.transform.position` へ反映している。

このため、前回の「アバター単体でworld座標を取得してAvatar Parameterへ入れる標準APIは見当たらない」という判断は、通常のAnimator Parameter/APIという意味では正しい。
一方でYL-ATGはContact/Constraintを使ったエンコードにより、アバターギミックだけでworld座標相当を外部OSCへ出す実例になっている。

ClipForVRChatへの応用可能性:

- 専用アバターギミック併用を許容するなら、現在のHead等のworld positionをClipForVRChatがOSCで受け取り、`player_local` のbasisとして使える可能性がある。
- `ATG/point` の対象をHeadにすると、取得できるのはplayer rootではなくHead付近。player root基準が必要なら、アバターroot相当や腰/足元相当を安定して取れる対象を別途検証する必要がある。
- yawはrotationそのものではなく、YL-ATGのrotation vector相当から水平成分を取って推定する実装になりそう。
- ClipForVRChat側には `/avatar/parameters/ATG/*` のOSC受信、値の復元、鮮度管理、basis更新、診断表示が必要。

制限とリスク:

- 利用者のアバターへ専用ギミックを入れる必要がある。一般ユーザー向けの初期導線としては重い。
- Modular Avatar、Avatar Dynamics Contact、OSC Avatar Parameter出力に依存する。
- Expression Parameter枠、Contact数、Avatar Performance、既存ギミックとのtag/parameter衝突を考慮する必要がある。
- Float parameterのOSC出力頻度・精度・遅延・VRChat更新による挙動差は実機確認が必要。
- default Head trackingでは、プレイヤーの足元/rootではなく頭部基準になる。撮影構図としてはむしろ有用な可能性があるが、仕様名を「player root」ではなく「avatar/head basis」などへ分けるべき。

判断:

「専用アバターギミックを合わせて開発し、アバターからOSCでworld座標を送る」案は可能性あり。
ただし、標準OSCだけで誰でも自動追従できる機能ではなく、アバター改変を前提にした拡張機能として扱うのが妥当。

次に進めるなら、YL-ATG方式を参考に最小ギミックを作り、ClipForVRChatが `head position + forward vector` を受け取って `player_local` basisを自動更新できるかを実機検証する。

## 2026-07-02 追加調査: 類似プロジェクトと取得手段

YL-ATGと似た仕組みの他プロジェクト、およびVRChatのプレイヤー位置座標を取得する手段を調査する。

### 類似/周辺プロジェクト

#### YL-ATG

- URL: https://github.com/YozoraKurage/YL-ATG
- 種別: VRChat内のアバター上オブジェクト座標をOSC Avatar Parameters経由で外部へ出す。
- 方式: Avatar Dynamics Contact Sender/Receiver + Constraintで座標をFloat parameterへエンコードし、`/avatar/parameters/ATG/*` として送信する。
- ClipForVRChatへの関連度: 高い。専用アバターギミック前提なら、head/avatar基準の自動basis取得に応用できる可能性がある。

#### VRCPrismStudio

- YL-ATG READMEで主な併用先として言及されている。
- 調査時点では公開リポジトリや詳細実装を確認できなかった。
- ClipForVRChatへの関連度: 中。YL-ATGの受け側/撮影側ワークフローの参考になり得るが、実装根拠としてはYL-ATG本体を見る方が確実。

#### VRCThumbParamsOSC

- URL: https://github.com/I5UCC/VRCThumbParamsOSC
- 種別: SteamVR controller action、Tracker button action、XInput actionをAvatar Parametersへ送るOSCツール。
- 方式: pyopenvr等で外部デバイス状態を読み、VRChat OSC serverへparameterとして入力する。
- 座標取得: trackerの存在/ボタン等は扱うが、VRChat内のplayer world positionを取得して外へ出すものではない。
- ClipForVRChatへの関連度: 低から中。OpenVR/SteamVRを読む外部アプリ設計やOSC parameter送信設定の参考にはなるが、VRChat world座標取得の直接解ではない。

#### VRCFaceTracking

- URL: https://github.com/benaclejames/VRCFaceTracking
- 種別: eye/lip/face tracking hardwareとVRChat OSC serverのbridge。
- 方式: 外部tracking hardwareの値をAvatar Parametersへ送る。
- 座標取得: 顔/目/口の表情入力であり、VRChat world position取得ではない。
- ClipForVRChatへの関連度: 低。外部データをOSCでAvatar Parametersへ流す成熟例として参考になる。

#### ContactGloveOSC

- URL: https://github.com/Diver-X/ContactGloveOSC
- 種別: ContactGloveのhand trackingをVRChat OSC APIで使うAvatar Gimmick + 自動セットアップツール。
- 方式: 外部グローブ情報をVRChatへ入力し、専用アバターギミックで扱う。
- 座標取得: VRChat内player world positionを外へ出すものではない。
- ClipForVRChatへの関連度: 低から中。専用アバターギミックと外部OSCアプリをセットで配る設計の参考にはなる。

#### OSC-Syncer

- URL: https://github.com/itsrivervrc/OSC-Syncer
- 種別: README上は別VRChatアカウント追従/ダンス同期を謳う個人プロジェクト。
- 方式: 公開内容は説明文中心で、位置取得の実装根拠は確認できない。
- ClipForVRChatへの関連度: 低。コンセプト例としては近いが、採用判断の根拠にはしない。

### 取得手段の整理

#### A. VRChat標準OSC readback

- `/usercamera/Pose` はUser Camera poseでありplayer rootではない。
- `/tracking/trackers/...` は外部tracker poseをVRChatへ入れる用途で、VRChatから現在位置を読む用途ではない。
- Avatar Parametersは標準built-inにworld position/yawがない。

評価: 汎用のplayer position取得手段としては不可。

#### B. Avatar Dynamics Contact/Constraintエンコード

- YL-ATG方式。
- アバター上にContact Sender/ReceiverとConstraintを組み、対象transformのworld座標をcontact値としてFloat parameterへ変換し、OSC Avatar Parametersで外部へ送る。
- スクリプトなしアバターギミックとして成立している点が重要。

評価: 専用アバター改変を許容するなら最有力。ClipForVRChatでは「avatar/head basis OSC bridge」として別機能化するのが妥当。

#### C. Udon/専用ワールド

- Udonなら `VRCPlayerApi.GetPosition()` / `GetRotation()` / `GetTrackingData()` を取得できる。
- ただし任意ワールドでは使えず、外部アプリへOSC送信する経路も別途設計が必要。

評価: 専用ワールド/撮影スタジオ用途なら候補。一般利用の標準機能には不向き。

#### D. 外部VRランタイム

- OpenVR/SteamVR/OpenXRからHMD/tracker poseを取る。
- VRCThumbParamsOSCのように外部VR runtimeを読むOSCツールは存在する。
- ただし取得できるのはVR runtime tracking spaceであり、VRChat world座標とは原点/Yawが一致しない。Desktop/Quest単体にも弱い。

評価: VR専用の高度機能なら候補。必ずキャリブレーションが必要。

#### E. 外部ハードウェアbridge

- VRCFaceTrackingやContactGloveOSCのように、外部hardware stateをOSCでVRChatへ送る方式。
- プレイヤー位置取得ではなく、外部入力をavatar parameterへ反映する用途。

評価: ClipForVRChatのplayer_local自動追従には直接使えないが、専用OSCアプリ+専用アバターギミックの配布設計は参考になる。

#### F. User Camera LookAtMe/Poseからの推定

- User CameraのLookAtMe挙動とPose readbackからプレイヤー方向を推定する案。
- rootではなく注視点推定になり、安定性や副作用が不明。

評価: 実験候補。標準仕様として扱うには弱い。

#### G. クライアント改造/メモリ読み/非公式API

- 技術的には位置情報へ到達できる可能性はあるが、VRChat ToSや安全性、配布リスクの問題が大きい。

評価: ClipForVRChatでは採用しない。

### 判断

公開例として最もClipForVRChatに近いのはYL-ATG。
他の多くのOSCプロジェクトは「外部デバイスからVRChatへ入力」するもので、「VRChat内のworld座標を外部アプリへ出す」ものではない。

ClipForVRChatで現実的に進めるなら、次の2段階がよい。

1. v0.1.8: 初期構図とUIを誤解しにくくし、`manual` basis と `avatar_osc` basis の両方を確認対象にする。
2. v0.1.8: YL-ATG方式を参考にしたAvatarBeaconを別配布し、`head position + forward vector` をOSCで受ける。
