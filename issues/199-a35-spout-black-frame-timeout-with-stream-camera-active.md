# a35でStream Camera起動済みでもSpout有効映像待ちがtimeoutする

## 問題

`v0.1.8-a35` のテスト撮影で、VRChat内のStream Cameraを起動し、SpoutもONにしているにもかかわらず、Stream方式のSpout取得が次のエラーで失敗する。

```text
Spout取得に失敗しました: Spoutフレームは取得できましたが、timeout内に有効な映像になりませんでした。VRChat Stream Cameraの映像が表示されているか確認してください。
```

ユーザー提供ログ:

- `/mnt/c/Users/user/Downloads/ClipForVRChat-v0.1.8-a35-windows-amd64/logs/2026-07-04.log`

ログ上は `VRCSender1` senderを検出し、`1920x1080` のSpoutフレームも取得している。
ただし helper の `frameStats` は `mean=0`, `stddev=0`, `near_black=1.0000`, `frame=0` で、timeout内に黒フレームしか取得できていない。

## 期待する挙動

- VRChat Stream Cameraが起動済みでSpout ONの場合、黒/透明/初期フレームを避け、有効な映像フレームを取得して保存できる。
- 有効映像にならない場合、VRChat側のStream Camera映像が黒いのか、Spout helperが新しいframeを待てていないのか、sender受信が止まっているのかをログだけで切り分けやすい。
- sender検出済み・frame番号不変・receive失敗・黒フレーム継続・有効フレーム取得を区別して診断できる。

## 受け入れ条件

- [x] a35ログを確認し、sender未検出ではなく黒フレーム継続で失敗していることを記録する。
- [x] `spout-capture` helper側で、frame counter、受信待機、黒フレーム判定、timeout処理を再調査する。
- [x] frame番号が増えていない場合、frame番号は増えているがreceiveできない場合、receiveできたが黒い場合を区別してログ/JSONへ出す。
- [x] Go側でhelper JSONの `code` / `message` / `frameStats` を分類し、sender missing / no new frame / black frame / transparent or blank frame を切り分けやすくする。
- [x] helperロジックのローカルテストを追加し、分類ルールとblank-frame統計を確認する。
- [x] テストまたはローカル検証可能な範囲で、黒フレーム統計とエラー分類を確認する。
- [ ] Windows/VRChat実機で、次回RCのログに `frame_state` / `first_frame` / `last_received_frame` / `receive_attempts` / `receive_successes` / `transparent` が出ることを確認する。

## 調査メモ

- 2026-07-04 20:52頃のログでは、`auto_select=true` で `VRCSender1` を選択している。
- helper結果JSON:
  - `code="capture_blank_frame"`
  - `senderName="VRCSender1"`
  - `width=1920`
  - `height=1080`
  - `frame=0`
  - `samples=15360`
  - `mean=0`
  - `stddev=0`
  - `nearBlackRatio=1`
  - `hostPath="C:\\Program Files (x86)\\Steam\\steamapps\\common\\VRChat\\VRChat.exe"`
- `frame=0` が継続している場合、Spout receiverが新規フレームを待てていない、またはVRChat側senderが黒初期フレームのまま更新されていない可能性がある。
- `frame=0` でも実データが更新される仕様の可能性があるため、frame counterだけに依存せず、統計値の変化や複数サンプルの履歴も見たい。
- 2026-07-04: helper側で `capture_no_new_frame` / `capture_receive_stalled` / `capture_blank_frame` / `capture_success` を分けるため、frame進捗、受信回数、最後のframe、最後の統計を保持する共通ロジックを切り出した。
- 2026-07-04: `tools/spout-capture/capture_logic_test.cpp` を追加し、透明フレーム統計と状態分類のローカルテストを通した。
