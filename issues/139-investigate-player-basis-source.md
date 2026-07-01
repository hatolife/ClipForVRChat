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

1. v0.1.8短期対応
   - 初期構図を `world` に戻す、または `player_local` 構図を初期OFFにする。
   - `player_local` は手動基準Pose方式として、UI文言を「プレイヤー追従」ではなく「手動基準」に寄せる。

2. v0.1.9以降の実験候補
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
