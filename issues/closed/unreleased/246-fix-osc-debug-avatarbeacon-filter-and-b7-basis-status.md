# OSCデバッグ送信説明、AvatarBeaconログフィルター、b7 basis受信判定を修正する

## 指示

OSCデバッグ送信
の説明文更新したい 改行は維持して
```
現在のVRChat OSC送信先へ任意のOSCを送信します。
例: /chatbox/input test
```

---
OSC受信ログ
AvatarBeaconが出していない物だけ表示するフィルターを作って新設したボタンを押すとそれが適用されるようにしたい
AvatarBeaconが出すOSCが変更されても正しく変更を追従するようにしたい

---
下記を日本語化
avatar-gimmicks/AvatarBeacon/README.md
---

AvatarBeacon Version v0.1.8-b7のようなバージョン情報のGameObjectはEditorOnlyにしてほしい

---
b7のAvatarBeacon_mainを設定した 受信も良さそう でもエラーになる

現在のAvatarBeacon受信状態: エラー
あとログに大量に
2026-07-07T09:51:27+09:00 auto-capture pose received: x=-4.694 y=0.865 z=-2.089 rx=4.878 ry=59.575 rz=0.050
2026-07-07T09:51:27+09:00 auto-capture pose received: x=-4.694 y=0.865 z=-2.089 rx=4.878 ry=59.575 rz=0.050
のようなものが出る

'/mnt/c/Users/user/Downloads/ClipForVRChat-v0.1.8-b7-windows-amd64/logs/2026-07-07.log'

## 文脈

`v0.1.8-b7` で AvatarBeacon は `AvatarBeacon_main` と `AvatarBeacon_12` に分かれ、parameter path は `avatar_beacon/...` 配下へ変更された。ユーザーは b7 の `AvatarBeacon_main` を導入し、OSC受信自体は見えているが、アプリ側の AvatarBeacon 受信状態がエラーになっている。提示ログでは旧 `coord/p/*` 形式を待っている summary が出ており、既存設定の移行または既定補完が新 path に追従できていない可能性がある。

## 解釈

OSCデバッグ送信の説明文を指定文へ差し替える。OSC受信ログには、現在アプリが AvatarBeacon basis として認識する address 群を除外するフィルターボタンを追加し、AvatarBeacon 側の address 定義変更時にも同じ判定を使って追従させる。AvatarBeacon README は英語混在をなくして日本語化する。Prefab内のバージョン表示 GameObject はタグを `EditorOnly` にする。b7 の `AvatarBeacon_main` を受信してもエラーになる問題は、旧 prefix/旧schema待ちを新 `avatar_beacon` schemaへ移行し、通常の `/usercamera/Pose` 受信ログを大量出力しないようにする。

## 問題

- OSCデバッグ送信の説明文が古い例を含んでいる。
- OSC受信ログで AvatarBeacon 以外だけを確認する導線がない。
- AvatarBeacon README に英語の説明が残っている。
- バージョン確認用 GameObject が `EditorOnly` ではない。
- b7 の `AvatarBeacon_main` 導入時にも、既存設定が旧 `coord/p/*` schema を待ち続けて受信状態がエラーになる場合がある。
- `/usercamera/Pose` 受信ログが高頻度で診断ログへ出続け、原因確認を妨げる。

## 期待する挙動

- OSCデバッグ送信の説明文は指定された2行表示になる。
- OSC受信ログで「AvatarBeacon以外」フィルターをボタンで適用できる。
- フィルター対象はアプリ側の AvatarBeacon address 判定と同じ定義から生成される。
- AvatarBeacon README は日本語で読める。
- Prefab内の `AvatarBeacon Version ...` GameObject は `EditorOnly` タグになる。
- 既存設定でも b7 の `avatar_beacon/...` basis を正しく認識する。
- `/usercamera/Pose` の通常受信は必要以上に診断ログへ連続出力されない。

## 受け入れ条件

- [x] OSCデバッグ送信説明文が指定の改行を維持している。
- [x] OSC受信ログに AvatarBeacon 以外フィルターボタンがあり、押すと該当フィルターが適用される。
- [x] AvatarBeacon basis address 判定がUIフィルターと受信処理で共通化されている。
- [x] `avatar-gimmicks/AvatarBeacon/README.md` が日本語化されている。
- [x] AvatarBeacon Prefab内の version GameObject の tag が `EditorOnly` になり、version更新scriptでも維持される。
- [x] 旧 `coord` / `forward` prefix 設定が `avatar_beacon` へ補正され、b7 `AvatarBeacon_main` の6parameter basisを認識できる。
- [x] `/usercamera/Pose` の受信診断ログが連続大量出力されない。
