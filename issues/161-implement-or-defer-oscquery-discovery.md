# OSCQueryによるVRChat OSC検出を実装または延期明示する

## 問題

Codex用仕様書では、OSCQueryで利用可能なOSCアドレス、型、ポートを検出し、使えない場合に手動設定へフォールバックする設計になっている。
現行実装は手動設定された `127.0.0.1:9000/9001` 相当のポートを使い、OSCQuery clientやendpoint検出処理を持たない。

## 期待する挙動

OSCQueryを実装する場合は、VRChatのOSCQuery情報からポートとUser Camera endpointの有無を確認し、設定画面や診断ログへ反映する。
v0.1.8で実装しない場合は、手動設定のみ対応であることをUI/仕様へ明記する。

## 受け入れ条件

- VRChat OSCQueryを使うか、v0.1.8では延期するかを決める。
- 実装する場合は、OSCQuery discovery、endpoint型検証、失敗時フォールバックを追加する。
- `/usercamera/Pose`、`/usercamera/Mode`、`/usercamera/Streaming` など自動撮影に必要なendpointの存在を診断できる。
- 実装しない場合は、自動撮影タブとREADME/SPECに「OSCホスト/ポートは手動設定」と明記する。
- OSCQuery未対応環境でも既存の手動OSC設定が壊れない。
