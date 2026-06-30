# VRChat Stream Camera/Spout映像そのものを保存する方式を調査する

## 問題

v0.1.8-rc12時点のStream方式は、ffmpegでデスクトップやVRChatウィンドウ周辺を取得しており、ユーザーが求める「VRChat Stream Camera映像そのもの」の保存になっていない。
そのため、デスクトップ周辺画像、白画像、VRChatウィンドウキャプチャになり、シャッター音なしでStream Camera出力を静止画化する本命要件を満たしていない。

## 期待する挙動

VRChatのStream Camera/Spout出力を直接受信し、そのフレームを画像として保存する実装方式を調査し、v0.1.8で実装可能な方針と作業チケットへ分割する。

## 受け入れ条件

- [x] VRChat Stream Camera/Spout出力を取得する一次情報を確認する。
- [x] 現行のGo/Wailsアプリへ組み込める方式を比較する。
- [x] v0.1.8で採用する実装方針を決める。
- [x] 実装作業を細かい日本語issueへ分割する。
- [x] この調査では実装コード変更を行わない。

## 調査メモ

- VRChat公式OSCドキュメントでは、VRChatのOSC受信ポート初期値は9000、送信ポート初期値は9001とされている。
- VRChat WikiのUser Camera OSC表では、`/usercamera/Mode` の `2` がStream、`/usercamera/Streaming` がSpout streamとして定義されている。
- FFmpeg公式devicesドキュメントの `gdigrab` はWindowsのデスクトップ全体、デスクトップ領域、またはウィンドウ内容を取得する入力デバイスであり、VRChat Stream CameraのSpout senderを直接受信するものではない。
- Spout2はWindows向けのビデオフレーム共有システムで、DirectX/OpenGLテクスチャ共有、SDK、サンプルを提供している。
- Spout2本体はBSD 2-Clause License。`SpoutLibrary` はC/C++互換DLLとして利用でき、`ReceiveImage`、`GetSenderList`、`GetSenderCount`、`GetSender`、`GetSenderInfo` などの受信/列挙APIを持つ。
- SpoutRecorderはGPL-3.0 Licenseのため、実装参考として調査するだけに留め、コードのコピーやリンク対象にはしない。

確認した一次情報:

- VRChat OSC Overview: https://docs.vrchat.com/docs/osc-overview.md
- VRChat Wiki OSC User Camera: https://wiki.vrchat.com/wiki/OSC#User_Camera
- FFmpeg devices / gdigrab: https://ffmpeg.org/ffmpeg-devices.html#gdigrab
- Spout2 README: https://github.com/leadedge/Spout2
- Spout2 License: https://github.com/leadedge/Spout2/blob/master/LICENSE
- Spout2 SpoutLibrary: https://github.com/leadedge/Spout2/tree/master/SPOUTSDK/SpoutLibrary
- SpoutRecorder License: https://github.com/leadedge/SpoutRecorder/blob/master/LICENSE

## 比較した方式

### 1. FFmpeg gdigrab継続

不採用。
`gdigrab` は画面領域やウィンドウを撮る方式であり、ユーザーが求めているStream Camera映像そのものではない。
これを続ける限り、デスクトップ周辺、VRChatウィンドウ周辺、白画像の問題を設計上解消できない。

### 2. FFmpegにSpout入力を期待する

不採用。
FFmpeg公式devicesドキュメントに標準Spout入力は確認できない。
ユーザー環境で追加ビルドや非標準プラグインを要求する構成は、v0.1.8の配布アプリとして不安定。

### 3. Go/Wails本体へSpout/DirectX/OpenGLを直接組み込む

v0.1.8では不採用。
CGO、Windowsネイティブ描画API、DLLロード、Wails本体プロセスの安定性を同時に扱う必要があり、クラッシュ時にアプリ全体へ影響する。
将来的な統合候補にはなるが、最短で実機確認可能な実装経路としては重い。

### 4. Windows専用Spout受信ヘルパーを同梱する

採用。
`spout-capture.exe` のような小さいWindowsヘルパーを作り、VRChat Stream CameraのSpout senderを列挙/選択し、1フレームをPNGへ保存する。
Go/Wails本体はOSCでUser CameraをStreamモードにし、ヘルパーをタイムアウト付きで呼び出し、成功した画像を既存のsidecar JSON、Discord投稿、履歴処理へ渡す。

採用理由:

- VRChat Stream Cameraの実体であるSpout streamを直接受信できる。
- Wails本体にDirectX/OpenGL処理を混ぜず、ネイティブ依存とクラッシュ範囲を隔離できる。
- Windows CI/Release workflowでヘルパーをビルド/同梱/検証する単位を分けやすい。
- `--list-senders`、`--capture` のCLIにすれば、設定画面の診断や実機ログが明確になる。

## 決定した実装方針

v0.1.8のStream方式は、FFmpeg/gdigrabの画面キャプチャ経路を主経路から外し、Windows同梱のSpout受信ヘルパーを主経路にする。

処理の責務:

1. Go/Wails本体はOSCで `/usercamera/Mode=2`、`/usercamera/Streaming=true`、構図Pose、Zoomなどを送る。
2. Spout受信ヘルパーは利用可能なsenderを列挙し、設定されたsenderまたは自動選択されたVRChat系senderから1フレームを取得する。
3. ヘルパーはPNG画像と、sender名、解像度、フレーム情報を機械可読な結果として返す。
4. Go/Wails本体は画像を検証し、既存の `finalizeAutoCaptureImage` 経路でsidecar JSON、Discord投稿、履歴へ紐づける。
5. 失敗時は履歴/Discordへ成功扱いで残さず、設定画面と診断ログへsender一覧、選択sender、タイムアウト、出力先を出す。

## 分割issue

- [129](129-add-spout-capture-helper.md): Spout受信ヘルパーを追加する。
- [130](130-add-spout-sender-settings-and-diagnostics.md): Spout sender設定と診断UIを追加する。
- [131](131-integrate-spout-helper-into-auto-capture.md): 自動撮影Stream方式をSpoutヘルパーへ統合する。
- [132](132-validate-spout-capture-output-and-metadata.md): Spout取得画像とメタデータを検証する。
- [133](133-update-auto-capture-stream-ui-and-docs-for-spout.md): Stream方式UIとドキュメントをSpout前提へ更新する。
- [134](134-package-spout-helper-in-ci-release.md): CI/ReleaseでSpoutヘルパーをビルド/同梱する。
- [135](135-add-spout-stream-camera-verification-guide.md): Stream Camera/Spout方式の実機確認手順を整備する。
