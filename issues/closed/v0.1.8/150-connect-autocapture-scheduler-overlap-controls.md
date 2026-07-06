# 自動撮影スケジュールの重複実行制御を接続する

## 問題

`AutoCaptureScheduleConfig` には `skipIfPreviousBatchRunning` と `maxBatches` がある。
現行スケジューラは1つのgoroutine内で `runAutoCaptureBatch()` を同期実行するだけで、前回Batch中に次回時刻が来た場合のskip/wait挙動を設定値として扱っていない。
また、`captureOnStart` で実行したBatchが `maxBatches` のカウントに含まれない。

## 期待する挙動

スケジュール設定値が、実際の自動撮影実行回数と重複制御に反映される。
Batch処理が撮影間隔を超えた場合でも、設定どおりにスキップ、待機、または実行を判断できる。

## 受け入れ条件

- `skipIfPreviousBatchRunning=true` の場合、前回Batch中に来たtickはスキップされ、診断ログに理由が残る。
- `skipIfPreviousBatchRunning=false` の場合、前回Batch完了後に次Batchを実行するか、仕様上の待機挙動を明文化する。
- `maxBatches` は `captureOnStart` で実行したBatchも含めて数える。
- Batch開始/終了/スキップ/停止理由が診断ログで追える。
- スケジューラの単体テスト、または時刻/tickerを差し替えられるテストを追加する。
- UIの説明が実際のskip/wait挙動と一致する。

## 対応内容

- v0.1.8実装に合わせて対応済み。詳細は `feat/v0.1.8-resolve-issues` の実装、README、SPEC、RELEASE_NOTESを参照。
