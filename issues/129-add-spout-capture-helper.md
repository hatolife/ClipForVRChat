# Spout受信ヘルパーを追加する

## 問題

現在のStream方式はFFmpegの `gdigrab` で画面領域を取得しており、VRChat Stream CameraのSpout映像そのものを保存できない。
Go/Wails本体へ直接DirectX/OpenGL/Spout処理を組み込むと、CGOやネイティブクラッシュの影響範囲が大きい。

## 期待する挙動

Windows専用の小さな `spout-capture.exe` を追加し、Spout2 SDKを使ってSpout senderを列挙し、指定senderまたは自動選択senderから1フレームをPNGとして保存できる。

## 受け入れ条件

- `spout-capture.exe --list-senders` で利用可能なSpout sender名、解像度、可能ならhost pathをJSONで出力できる。
- `spout-capture.exe --capture --sender <name> --output <png> --timeout-ms <ms>` で指定senderの1フレームをPNGへ保存できる。
- sender未指定時は、単一senderまたはVRChat/Stream Cameraらしいsenderを自動選択し、複数候補で判断不能な場合は候補一覧付きで失敗する。
- 成功時はsender名、幅、高さ、フレーム番号または取得時刻をJSONで標準出力へ返す。
- 失敗時は機械可読なエラーコードと、人間が読める日本語/英語メッセージを返す。
- ヘルパーはWindows x64でビルドでき、通常のアプリ起動とは独立して単体実行できる。
- Spout2のBSD 2-Clause Licenseを遵守し、同梱するライセンス表記を追加する。
- SpoutRecorderのGPL-3.0コードをコピー、リンク、派生利用しない。

## 実装メモ

- 第一候補はSpout2の `SpoutLibrary` またはSpout2 SDKの受信APIを使う。
- `ReceiveImage` でCPU側ピクセルを受け取り、PNGエンコードして保存する。
- PNG保存はWindows標準API、WIC、またはライセンス確認済みの小さいPNGライブラリを使う。
- CLIはGo本体から呼びやすいよう、標準出力JSON、標準エラー診断、終了コードを固定する。
