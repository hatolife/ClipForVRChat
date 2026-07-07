# AvatarBeaconバージョンOSC受信をログに記録する

## 指示

終わったらでいい　oscでアバターギミックのバージョン情報受信したらログにそれを記載するようにしたい

## 文脈

AvatarBeacon は前回の実装で、OSC送信時に1回だけ `/avatar/parameters/AvatarBeacon/version` へバージョン番号を送信するようになった。受信側でこの値を確認できれば、ユーザーが導入しているアバターギミックのバージョン切り分けに使える。

## 解釈

ClipForVRChat が OSC 受信中に AvatarBeacon の version parameter を受け取った場合、診断ログへ `AvatarBeacon version received` のような専用ログを残す。通常のbasis受信やUserCamera Pose受信の高頻度ログとは別に、version受信時だけ記録する。

## 問題

- AvatarBeacon version OSC を受信しても、ログから導入済みギミックのバージョンを確認しにくい。

## 期待する挙動

- `/avatar/parameters/AvatarBeacon/version` 受信時に診断ログへ version 情報が出る。
- OSCログにも通常の受信行として値が見える。
- version受信ログは高頻度に連続出力されない。

## 受け入れ条件

- [x] AvatarBeacon version OSC受信時、診断ログに address と値が記録される。
- [x] 同じversionが連続受信されても診断ログが大量出力されない。
- [x] 既存のAvatarBeacon basis判定を壊さない。
