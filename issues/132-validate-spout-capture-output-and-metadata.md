# Spout取得画像とメタデータを検証する

## 問題

これまでのStream方式では、白画像、デスクトップ周辺、VRChatウィンドウ領域など、ユーザーが意図しない画像でも成功扱いになることがあった。
Spoutヘルパー導入後も、sender誤選択や空フレームを成功扱いすると同じ問題が残る。

## 期待する挙動

Spout取得後に画像と取得元メタデータを検証し、VRChat Stream Cameraのsenderから有効なフレームを保存できた場合だけ成功扱いにする。
sidecar JSONやDiscord投稿にも、Stream方式で取得したsender情報を紐づける。

## 受け入れ条件

- Spout取得後にPNGをデコードし、0バイト、読めない画像、幅/高さ0の画像を失敗扱いにする。
- ほぼ単色の白画像/黒画像など、明らかに無効なフレームを検出した場合は失敗扱いにし、診断ログに理由を記録する。
- helperから返ったsender名、幅、高さ、フレーム番号または取得時刻をsidecar JSONへ保存する。
- Discord投稿本文または添付メタデータに、設定で許可されている範囲で撮影方式 `stream/spout` とsender情報を含められる。
- 失敗時は履歴、Discord、通常画像削除処理へ成功画像として混入しない。
- 画像検証が過剰に厳しすぎる場合に備え、診断ログから閾値と失敗理由を追える。
- 既存のPresenceユーザー情報保持と画像紐づけを壊さない。

## 実装メモ

- 「デスクトップを撮ったか」を完全自動判定するのは難しいため、まずSpout helper経由であること、senderメタデータがあること、画像として有効であることを必須にする。
- 単色検出は白画像再発防止のガードであり、実機確認で閾値調整する。

## 対応内容

- v0.1.8実装に合わせて対応済み。詳細は `feat/v0.1.8-resolve-issues` の実装、README、SPEC、RELEASE_NOTESを参照。

## 再監査メモ

- 2026-07-01: [#164](164-audit-v018-completed-items.md) の再監査で未達が見つかったため、完了扱いを取り消して `要対応` に戻した。

## 実装前調査メモ

実装方針:

- `validateCapturedImage()` は `DecodeConfig` だけでは白画像/黒画像を検出できないため、画像をdecodeしてサンプル検査を追加する。
- 全ピクセル走査は4K画像で重くなるため、最大16384点程度の等間隔サンプリングで判定する。
- 判定は過剰に厳しくしない。以下を「明らかな無効フレーム」として失敗扱いにする。
  - 有効サンプルの99%以上が白近傍 (`R/G/B >= 250`) かつRGBの分散が極小。
  - 有効サンプルの99%以上が黒近傍 (`R/G/B <= 5`) かつRGBの分散が極小。
  - ほぼ全透明またはRGBAが完全固定で、標準偏差が閾値未満。
- 失敗時は `mean`, `stddev`, `white_ratio`, `black_ratio`, `sample_count` を診断ログへ出す。
- `SpoutCaptureResult` のsender名/幅/高さ/フレーム/取得時刻はすでにsidecarへ入るため、Discord本文にもsender情報を追加する。

対象ファイル:

- `src/internal/appcore/spout.go`
- `src/internal/appcore/autocapture.go`
- `src/internal/appcore/autocapture_test.go`
- `src/internal/appcore/spout_test.go` を新設してもよい。

小タスク:

- `validateCapturedImage(path, logPath)` または `validateCapturedImage(path) (diagnostic, error)` に変更し、診断値をログに残せる形にする。
- テスト用に白PNG、黒PNG、単色でないPNGを生成する。
- JPEG変換後ではなく、Spout helperが保存したPNGを検証する。
- `autoCaptureDiscordContent()` へStream metadataを渡せるようにし、`stream/spout` とsender名を本文に含める。

確認方法:

- 白一色/黒一色PNGは失敗する。
- 通常の小さなカラーPNGは成功する。
- 失敗時に履歴、Discord投稿、成功sidecarへ混入しない。
