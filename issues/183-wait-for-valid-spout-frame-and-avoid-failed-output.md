# Spout取得直後の空フレームが透明PNGとして保存される

## 問題

`v0.1.8-rc26` のStream方式テスト撮影で、カメラ起動後にSpout senderは取得できているが、保存されたPNGが完全に透明になり、アプリ側で `取得画像がほぼ透明です` と失敗する。

実機出力PNGを確認すると、1920x1080 RGBA の全チャンネルが0で、VRChat映像ではなくSpout起動直後の空フレームを保存している状態だった。

また、失敗した透明PNGが自動処理の監視対象として残り、通常の画像処理/Discord投稿に流れる可能性がある。

## 期待する挙動

Stream方式では、Spoutの `ReceiveImage` が成功しただけで終了せず、timeout内で非空の有効フレームを待って保存する。

取得失敗または画像検証失敗時は、失敗画像や一時ファイルを監視対象の最終ファイル名で残さない。

## 受け入れ条件

- [x] Spout helperが全透明/全黒/全白に近い初期フレームを即保存せず、timeout内で次フレームを待つ。
- [x] RGBが有効でalphaだけが0に近いフレームは、スクリーンショット用途として不透明PNGへ正規化する。
- [x] Go側はhelper出力を一時ファイルへ保存し、検証成功後だけ最終ファイル名へ移動する。
- [x] 検証失敗時のログにsender名、サイズ、frame、画像統計を出し、原因を切り分けやすくする。
- [x] 失敗したSpout取得画像が自動処理/Discord投稿へ流れない。
- [x] Go test、frontend check、Spout helper build/checkが通る。

## 実装メモ

- `tools/spout-capture/main.cpp` で受信ピクセルの簡易統計を取り、空フレームならtimeoutまで再受信する。
- `write_png_wic()` に渡す前にalphaを255へ正規化し、VRChat Stream Camera映像を不透明画像として保存する。
- `src/internal/appcore/spout.go` で `outputPath + ".tmp"` などへhelper出力し、検証後にrenameする。

## 対応メモ

- 2026-07-04: rc26の実機出力PNGを確認し、1920x1080 RGBAの全チャンネルが0の空フレームであることを確認した。
- 2026-07-04: Spout helperで受信フレームのluma統計を取り、全黒/全白に近い初期フレームはtimeout内で待ち直すようにした。
- 2026-07-04: 保存前にalphaを255へ正規化し、RGBが有効なVRChat映像を透明PNGにしないようにした。
- 2026-07-04: Go側は `.tmp` へhelper出力し、画像検証成功後だけ最終パスへrenameするようにした。
- 2026-07-04: rc28のRelease workflowでMSVCの `max` macro衝突によりSpout helper buildが失敗したため、`(std::max)(...)` でmacro展開を避けるようにした。
- ローカル検証: `go test ./...`、frontend build、template literal check、Wails API surface check、Spout helper object compile。
