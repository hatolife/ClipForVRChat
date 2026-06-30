# 自動撮影multi/Camera Dolly設定を実装または削除する

## 問題

`AutoCaptureCaptureConfig` には `concurrentMode`、`requestedCameraCount`、`multiBackend`、`fallbackToSequential` があり、正規化と診断出力の対象になっている。
しかし `AutoCaptureRunner.RunOnce()` は有効構図を常に1つずつ順番に処理しており、multi撮影やCamera Dolly backendは実装されていない。

## 期待する挙動

v0.1.8でmulti撮影を扱うなら、設定値に従って実カメラ数と撮影パスを計算し、Stream/Spout/Dolly連携へ接続する。
扱わないなら、設定値を内部/診断/仕様から未実装機能として見せず、Sequential固定であることを明示する。

## 受け入れ条件

- `concurrentMode=multi` を指定した場合に、実装済みの動作になるか、明確にSequentialへ正規化される。
- `requestedCameraCount`、`multiBackend`、`fallbackToSequential` が、未実装のまま保存/診断されない。
- Multi実装を行う場合は、CapturePlanner相当で `C_eff` とpass数を計算し、テストを追加する。
- Camera Dolly未対応環境では、設定に従ってSequentialへフォールバックし、その理由を診断ログに出す。
- UI、README、SPEC、sidecar JSONが、実際に動く撮影方式と矛盾しない。

## 対応内容

- v0.1.8実装に合わせて対応済み。詳細は `feat/v0.1.8-resolve-issues` の実装、README、SPEC、RELEASE_NOTESを参照。
