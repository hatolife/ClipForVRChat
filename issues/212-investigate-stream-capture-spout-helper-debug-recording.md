# Stream方式撮影失敗をSpout helper録画から切り分ける

## 問題

現在Streamモードで撮影できない。
まずはClipForVRChat本体の自動撮影処理ではなく、`spout-capture.exe` 自体がVRChat Stream CameraのSpoutフレームを録画/保存できていない可能性を優先して切り分ける。

関連する既存チケット:

- `199-rc35-spout-black-frame-timeout-with-stream-camera-active.md`
- `201-rc36-spout-transparent-frame-with-stream-camera-active.md`
- `202-add-spout-helper-debug-recording-and-hide-console.md`

## 期待する挙動

- `spout-capture.exe` 単体でsender列挙、1フレーム取得、デバッグ録画出力の成否を確認できる。
- helper単体で失敗する場合、sender未検出、受信停止、frame番号不変、RGBA生フレーム空/黒/透明、PNG書き出し失敗のどこで止まっているかを特定できる。
- helper単体で成功する場合、ClipForVRChatから渡す設定、helper実行パス、出力パス、OSCでのStream Camera起動、PNG後続変換/検証のどこで失敗するかを段階的に切り分けられる。

## 受け入れ条件

- [ ] 実機環境で `spout-capture.exe --list-senders` の結果を記録する。
- [ ] 実機環境で `spout-capture.exe --capture` の単体実行結果、終了コード、JSON、出力PNGの有無を記録する。
- [ ] `--debug-dir` / `--debug-frames` の出力有無と、`frames.jsonl`、frame metadata、`.rgba` の内容傾向を確認する。
- [ ] `spout-capture.exe --diagnose` でsender一覧、frame番号、ReceiveImage成否、RGBA統計を時系列で記録できる。
- [ ] ClipForVRChatのテスト撮影ログから、helper実行引数、debug_dir、debug_frames、helper JSON、出力ファイル検証結果を確認する。
- [ ] helper単体失敗か、ClipForVRChat統合部失敗か、VRChat Stream Camera/sender側失敗かを判断する。
- [ ] 原因に応じて修正チケットまたは既存チケットへの追記内容を決める。

## 調査計画

1. `spout-capture.exe` の場所とバージョンを確認する。
   - Release同梱版、埋め込み展開版、開発ビルド版のどれを実行しているかを混同しない。
   - `--version` のJSONとファイルパスを記録する。
2. VRChat側の前提を固定する。
   - VRChat内でStream Cameraを手動表示し、Spout/StreamingをONにする。
   - ClipForVRChatの自動起動OSCは一旦使わず、VRChat側がsenderを出している状態から始める。
3. helper単体でsender列挙を確認する。
   - `spout-capture.exe --list-senders` を実行し、`VRCSender1` などのsender名、解像度、hostPathを記録する。
   - senderが0件なら、ClipForVRChatではなくVRChat Stream Camera/Spout出力側を調査対象にする。
4. helper単体で1フレーム取得を確認する。
   - `--capture --output <png> --timeout-ms 10000` を実行する。
   - 成功JSON、失敗JSON、終了コード、PNGが生成されたかを記録する。
   - 失敗コードが `capture_no_new_frame` / `capture_receive_stalled` / `capture_blank_frame` のどれかを確認する。
5. helper単体で録画デバッグを確認する。
   - `--debug-dir <dir> --debug-frames 30` を付け、`spout-capture-debug.log`、`frames.jsonl`、`*.json`、`*.rgba` が出るか確認する。
   - `.rgba` が複数出る場合は、frame番号、alpha、RGB統計が変化しているかを見る。
   - `.rgba` が出ない場合は、Spout受信前またはdebug出力処理で止まっている可能性を優先する。
6. ClipForVRChat経由のテスト撮影を確認する。
   - 設定の「Spout録画デバッグ」をONにしてテスト撮影する。
   - アプリログの `spout capture begin` に `--debug-dir` と `--debug-frames` が含まれるか確認する。
   - helper単体と同じsender/timeout/debug frame countで実行されているか比較する。
7. 結果で分岐する。
   - helper単体でも失敗する: `spout-capture.exe` / Spout receiver / VRChat sender状態を主因として調査する。
   - helper単体は成功し、ClipForVRChat経由だけ失敗する: helperパス、実行権限、出力先、引数、timeout、起動タイミング、後続PNG検証/変換を調査する。
   - helper単体もClipForVRChat経由もPNGは取れるが内容が黒/透明: VRChat Stream Cameraの表示状態、sender更新、Spout受信形式、frame counterの扱いを調査する。

## 実装メモ

- 2026-07-07: 現在の仮説は「`spout-capture.exe` 自体が有効フレームを録画できていない可能性が高い」。まずhelper単体のsender列挙、capture、debug録画を確認してから、ClipForVRChat統合部へ進む。
- 2026-07-07: `v0.1.8-rc40` 分離版の `/mnt/c/Users/user/Downloads/ClipForVRChat-v0.1.8-rc40-windows-amd64-separated/spout-capture.exe` を使って、helper単体の `--version`、sender列挙、capture、debug録画を確認する。
- 2026-07-07: `--version` は `v0.1.8-7f242a9`。初回 `--list-senders` は `senders=[]` だったが、ゲーム内カメラ起動済み状態で再実行すると `VRCSender1` (`1920x1080`, hostPath=`C:\Program Files (x86)\Steam\steamapps\common\VRChat\VRChat.exe`) を検出した。
- 2026-07-07: `--capture --sender VRCSender1 --timeout-ms 10000 --debug-dir ... --debug-frames 30` は exit code 3、`code=capture_blank_frame`。`receiveAttempts=314`、`receiveSuccesses=313` で受信自体は成功しているが、`firstFrame=0`、`lastReceivedFrame=0`、`frame=0` のまま進まない。`frameStats` は `mean=0`、`stddev=0`、`transparentRatio=1` で全透明。
- 2026-07-07: debug録画は30フレーム分の `.rgba` と `.json` を出力した。各 `.rgba` は `8294400` bytes (`1920*1080*4`)。1枚目と30枚目は `cmp` で完全一致し、先頭64 bytesも全て `00`。helperはsenderへ接続してフレーム受信に成功しているが、VRChat側senderが全ゼロRGBAかつframe番号0のまま更新されていない状態と判断する。
- 2026-07-07: 調査中に `frames.jsonl` / per-frame `.json` の `session` フィールドに余分な `"` が入り、不正JSONになるバグを確認した。撮影失敗の主因ではないが、デバッグ解析を阻害するため `tools/spout-capture/main.cpp` で修正した。
- 2026-07-07: 続きの調査補助として、1回のcaptureではなく一定時間のsender/receive/frame stats推移を記録する `spout-capture.exe --diagnose` を追加した。stdoutはsummary JSONのみ、`--debug-dir` 配下へ `diagnose.jsonl`、`diagnose-summary.json`、既存形式のframe dumpを出力する。
- 2026-07-07: ローカル検証は `ctest --test-dir /tmp/cfvrc-spout-logic-diagnose --output-on-failure` と `x86_64-w64-mingw32-g++ -std=c++17 -Wall -Wextra -fsyntax-only -I/tmp/cfvrc-spout-mingw-diagnose/_deps/spout2-src/SPOUTSDK/SpoutLibrary tools/spout-capture/main.cpp`。MinGWでの完全ビルドはSpout2側include `Shellapi.h` / `Commctrl.h` の大文字小文字問題で停止したため、CI/MSVCまたは実Windowsビルドでの確認が必要。
