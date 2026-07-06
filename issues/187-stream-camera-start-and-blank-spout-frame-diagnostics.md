# カメラ未起動/起動直後のStream Camera Spout取得を安定化する

## 問題

`v0.1.8-a29` で、VRChatのカメラ未起動状態からStream方式のテスト撮影を実行すると、アプリは `/usercamera/Mode=2` と `/usercamera/Streaming=true` を送っているが、Spout senderが作成されず `Spout senderがありません` で失敗する。

また、カメラを起動した状態でテストすると、Spout sender `VRCSender1` は検出でき、1920x1080のフレーム受信までは進むが、helperのblank-frame判定をtimeout内に通過せず `Spoutフレームは取得できましたが、timeout内に有効な映像になりませんでした` で失敗する。

現状のログでは、blank-frameとして捨てたフレームのRGB統計が残らないため、本当に全黒なのか、判定が厳しすぎるのか、Spout画素取得自体が期待と違うのかを切り分けにくい。

## 期待する挙動

カメラ未起動状態でも、Stream方式の自動撮影開始時にVRChat Stream CameraとSpoutが可能な限り自動でONになる。

senderはあるが有効フレームにならない場合は、最後に受信したフレームの統計をログへ出し、次の実機確認で原因を特定できる。

## 受け入れ条件

- [x] Stream起動時の `/usercamera/Streaming` はOSC boolに加え、numeric `1` も互換送信する。
- [x] Stream停止/復元時の `/usercamera/Streaming=false` はnumeric `0` も互換送信する。
- [x] Spout helperの `capture_blank_frame` エラーに最後のフレーム統計を含める。
- [x] Go側のhelper errorログにsender名、サイズ、最後のフレーム統計が出る。
- [x] Go test、frontend check、Spout helper compileが通る。

## 実装メモ

- VRChat公式のCamera OSCは `/usercamera/Mode` をGet/Set、`/usercamera/Streaming` をbool endpointとしている。
- VRChat WikiのOSC説明では、bool parameterに対して `,T`、`,i 1`、`,f 1.0` のような送信がtrue扱いできると説明されている。Camera endpointでもnumeric boolを補助送信することで実機挙動を確認する。

## 対応メモ

- 2026-07-04: a29ログで、カメラ未起動時は `/usercamera/Mode=2` が受信値として戻る一方、Spout senderが0件のままになることを確認した。
- 2026-07-04: a29 separated2ログで、カメラ起動後はsender `VRCSender1` が見えるが、helperがtimeout内に有効フレーム判定できず `capture_blank_frame` になることを確認した。
- 2026-07-04: `/usercamera/Streaming` の起動/停止/復元で、OSC boolに加えてnumeric `1/0` も送る互換送信を追加した。
- 2026-07-04: Spout helperの `capture_blank_frame` JSONにsender名、width/height、frame、最後のフレーム統計を追加し、Go側のhelper errorログにも出すようにした。
- ローカル検証: `go test ./...`、frontend template literal check、Wails API surface check、`tools/spout-capture/main.cpp` のmingw単体compile。
