# 自動撮影でVRChat写真が保存されない

## 問題

v0.1.8-rc1 の自動撮影で「撮影後のVRChat写真ファイルを検出できませんでした」と表示され、VRChat側にも写真が保存されない。`/usercamera/Capture` をbool押下として送っているため、VRChatのAction OSCとして認識されていない可能性がある。

## 期待する挙動

自動撮影時にVRChat User CameraがPhotoモードへ切り替わり、Action OSCとして撮影指示が送られ、VRChat標準写真が保存される。保存後は検出した写真へsidecar JSONを作成する。

## 受け入れ条件

- `/usercamera/Capture` と `/usercamera/Close` を引数なしAction OSCとして送信する。
- Photoモード切替後に短い待機を入れ、VRChat側のカメラ起動前に撮影指示を送らない。
- 写真検出対象はVRChat写真フォルダを優先し、設定出力先との混同を避ける。
- OSC Actionのパケット生成をテストで確認する。

## 実装メモ

- `/usercamera/Capture` と `/usercamera/Close` は引数なしのOSC Actionとして送信する。
- Photoモード切替後、撮影前に800ms待機する。
- 写真検出は既存の「VRChat写真フォルダ」設定を優先し、補助的に自動撮影出力先も確認する。

