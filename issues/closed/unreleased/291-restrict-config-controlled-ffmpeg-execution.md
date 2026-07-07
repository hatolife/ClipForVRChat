# 設定由来のffmpeg実行パスを制限する

## 指示

> 対処

## 文脈

`reports/security/2026-07-07T17-53-16.203Z/9b49f259a284819193c25862c885e028` は、外部JSON設定から読み込まれる自動撮影Stream設定がffmpeg実行パスと引数を保持し、そのまま `exec.CommandContext` に渡されることで、細工された設定ファイルを開いた利用者の権限で任意ローカル実行ファイルを起動できると指摘している。
現在HEADでは主経路はSpout helperへ移行済みだが、互換用の `legacyFfmpegPath` / `legacyInputArgs` と `captureStreamFrameWithFFmpeg` が残っている。

## 解釈

調査済みfindingへの対処として、設定ファイル由来の値で任意実行ファイルを選択できないよう最小修正する。
既存の通常用途であるPATH上の `ffmpeg` / `ffmpeg.exe` 実行は維持する。

## 問題

`legacyFfmpegPath` が設定ファイルから読み込まれ、絶対パスや区切り文字を含む任意パスでも存在すれば許可される。
さらに実行時は `ResolveFFmpegPath` の解決済みパスではなく設定値が `exec.CommandContext` に渡される。

## 期待する挙動

- 設定由来のffmpeg実行名は安全な既定の `ffmpeg` / `ffmpeg.exe` のみに限定される。
- 絶対パス、相対パス、パス区切り文字を含む値、ffmpeg以外のコマンド名は拒否される。
- 実行時は検証・解決済みのPATH上ffmpegを使う。

## 受け入れ条件

- [x] `ResolveFFmpegPath` が任意パスや任意コマンド名を拒否するテストがある。
- [x] `captureStreamFrameWithFFmpeg` が解決済みパスを実行に使う。
- [x] appcoreテストが成功する。

## 対応内容

- `ResolveFFmpegPath` で、設定由来の値に絶対パス、相対パス、パス区切り文字、`ffmpeg` / `ffmpeg.exe` 以外のコマンド名が指定された場合は拒否するようにした。
- `captureStreamFrameWithFFmpeg` では、検証・解決済みのffmpegパスを `exec.CommandContext` に渡すようにした。
- 任意パスや任意コマンド名の拒否、解決済みパスの使用を回帰テストで確認した。
