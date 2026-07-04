# AvatarBeaconのOSC送信頻度と負荷影響を評価し、送信頻度を設定可能にする

## 問題

AvatarBeacon は `coord/*` と `forward/*` の basis parameter をOSC Avatar Parametersとして出力する。
現状、Prefab上に明示的な送信頻度設定はなく、VRChatのContact Receiver / Expression Parameter / OSC出力の更新挙動に依存している。

そのため、AvatarBeacon導入時にOSCが大量送信される可能性、VRChat本体・ネットワーク・ClipForVRChat・他OSC受信アプリへの悪影響、低頻度化した場合の自動撮影精度への影響が未評価。

特に、Unity Inspector上で「どの程度の頻度で送信するか」を調整できないため、利用環境に応じた安全側の設定ができない。

## 期待する挙動

AvatarBeaconの実際のOSC送信頻度と負荷を把握できる。

必要であれば、Unity Inspector上で送信頻度または更新頻度を変更できる。

低頻度設定にした場合でも、ClipForVRChatの `avatar_osc` basisとして実用可能な範囲、破綻する範囲、推奨値が分かる。

## 受け入れ条件

- [ ] 現状のAvatarBeaconが1秒あたり何packet程度送るかを実測する。
  - 静止時
  - 歩行/移動時
  - 頭を動かした時
  - `coord/*` / `forward/*` それぞれ
- [ ] 現状Prefabに明示的なrate limitや周期設定が存在しないこと、または存在する場合はその値と設定箇所を仕様へ書く。
- [ ] OSC大量送信による悪影響を評価する。
  - VRChat本体の負荷
  - OSC受信ポートを共有する他アプリへの影響
  - ClipForVRChatの受信処理、OSCログUI、forward機能への影響
  - 診断ログやメモリ使用量への影響
- [ ] Unity Inspector上で送信頻度を設定できる設計案を作る。
  - 例: `UpdateRateHz`、`SendIntervalSec`、`Quality` など
  - 設定可能範囲、既定値、最小値、最大値
  - 既存Prefab/Animator/Contact構成で実現可能か
- [ ] 1Hzで自動撮影の構図精度やyaw追従に問題が出るか確認する。
- [ ] 0.1Hzで自動撮影の構図精度やyaw追従に問題が出るか確認する。
- [ ] ClipForVRChat側のfreshness判定と整合する推奨送信頻度を決める。
- [ ] 推奨値と注意点をAvatarBeacon READMEまたは仕様mdへ追記する。

## 実装メモ

- 現状確認時点では、`avatar-gimmicks/AvatarBeacon/Assets/PoppoWorks/AvatarBeacon/Prefabs/AvatarBeacon.prefab` に `coord/*` と `forward/*` のExpression Parameter定義はあるが、明示的な送信周期設定は見当たらない。
- VRChat OSCはAvatar Parameterの変化を出力するため、実際のpacket rateはAvatar Dynamics Contact / Expression Parameter更新頻度、値の変化量、VRChat側OSC出力仕様に依存する可能性が高い。
- ClipForVRChat側の `avatar_osc` freshnessは短すぎる送信間隔を前提にしすぎると低頻度設定でstaleになりやすい。低頻度対応をする場合は、送信頻度設定とfreshness設定の整合が必要。
- 単に送信頻度を下げるだけだと、撮影直前のプレイヤー位置・Head yawが古い値になり、player_local構図がずれる可能性がある。

## 検証観点

- VRChat実機でOSC MonitorまたはClipForVRChatのOSCログUIを使い、1秒あたりの受信件数を測る。
- 1Hz、0.5Hz、0.2Hz、0.1Hzのような候補で、静止時/移動時/撮影直前に必要な精度を比較する。
- ClipForVRChat側で `avatar_osc` がstaleにならないか確認する。
- OSC forward有効時に他アプリが過負荷にならないか確認する。
- Inspector設定変更後、AvatarBeacon source zipと手動unitypackage化手順に反映されることを確認する。
