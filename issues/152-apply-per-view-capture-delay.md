# 構図ごとのcaptureDelayMsを撮影待機へ反映する

## 問題

`CameraViewConfig.CaptureDelayMS` は設定型、正規化、診断ログに存在する。
しかしPhoto方式でもStream方式でも、settle後に `captureDelayMs` を待たずに撮影処理へ進む。
設定値があるのに動作へ反映されない。

## 期待する挙動

各構図の `captureDelayMs` が、Pose/カメラパラメータ適用とsettle完了後、実際のCaptureまたはStreamフレーム取得の直前に待機時間として適用される。

## 受け入れ条件

- Photo方式で `/usercamera/Capture` を送る前に `captureDelayMs` を待つ。
- Stream方式でフレーム取得を始める前に `captureDelayMs` を待つ。
- 待機中にcontextがキャンセルされた場合は、撮影を中断し、診断ログと結果に理由を出す。
- `captureDelayMs=0` の場合は追加待機しない。
- 単体テストまたは差し替え可能なsleep処理で、待機が呼ばれることを確認する。
- UIに出す場合は、settle delayとの違いが分かる文言にする。出さない場合は内部設定として扱う理由を明記する。
