# Spout helperに録画デバッグ出力を追加し起動コンソールを隠す

## 問題

`v0.1.8-rc36` ではSpout sender検出と `ReceiveImage` 成功までは進んでいるが、frame番号が進まず全透明フレームになる。
この状態では、最終的な静止画保存処理だけを見ても、Spout helper自体が正しく映像を受け取れているのか、受信直後のRGBAデータがどうなっているのかを確認しにくい。

また、`spout-capture.exe` 起動時にコンソールウィンドウが見えるため、通常利用時の体験が悪い。

## 期待する挙動

- `spout-capture.exe` にデバッグ引数を追加し、受信した元RGBAフレームを一定数保存できる。
- 保存先にはhelper自身のログ、受信frameごとのmetadata、RGBA生データが残り、静止画抽出前の状態を確認できる。
- 設定画面にSpout録画デバッグ用フラグを追加し、ONの場合はテスト撮影時にデバッグ引数を付けてhelperを起動する。
- 通常のhelper起動ではWindowsのコンソールウィンドウを表示しない。

## 受け入れ条件

- [x] `spout-capture.exe --capture` にデバッグ出力先引数を追加する。
- [x] デバッグ有効時、受信したRGBAフレームを数フレーム保存する。
- [x] デバッグ有効時、helperログとframe metadataを保存する。
- [x] Go側がテスト撮影時だけデバッグ引数を渡せる。
- [x] 設定画面に録画デバッグ用フラグを追加する。
- [x] Windowsでhelper起動時のコンソール表示を抑止する。
- [x] ローカルで可能なテストを追加・更新する。

## 実装メモ

- `spout-capture.exe --capture` に `--debug-dir` と `--debug-frames` を追加した。
- デバッグ有効時は `spout-capture-debug.log`、`frames.jsonl`、`<session>_frame_000001.rgba`、`<session>_frame_000001.json` を保存する。
- `.rgba` はPNG化前、かつalphaを強制不透明化する前の受信直後RGBA8を保存する。
- 設定画面の「Spout録画デバッグ」はテスト撮影時だけ有効で、通常の自動撮影にはdebug引数を渡さない。
- Windowsでは `SysProcAttr.HideWindow` でhelper起動時のコンソール表示を抑止する。
