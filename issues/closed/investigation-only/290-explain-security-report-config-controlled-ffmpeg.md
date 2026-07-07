# 設定由来ffmpeg実行セキュリティレポートを解説する

## 指示

> '/home/user/work/ClipForVRChat/reports/security/2026-07-07T17-53-16.203Z/9b49f259a284819193c25862c885e028'を日本語で解説して

## 文脈

指定されたファイルは、Codex Securityが生成した `Imported config can execute arbitrary ffmpeg command` というhigh findingのレポートである。
外部JSON設定から読み込まれる自動撮影Stream設定がffmpeg実行パスと引数を保持し、それが `exec.CommandContext` に渡されることで、細工された設定ファイルを開いた利用者の権限で任意ローカル実行ファイルを起動できると指摘している。

## 解釈

レポート本文を読み、指摘内容、攻撃シナリオ、根拠、影響、現在HEADでの残存状況、推奨対応を日本語で分かりやすく説明する。
この依頼ではコード修正までは行わず、現時点の理解と対応方針を整理する。

## 問題

- 外部 `.json` 設定を `LoadConfig` で読み込める。
- 設定内のStream撮影用ffmpeg実行名と引数がJSON由来の値として扱われる。
- 正規化処理は文字列のtrimや既定値補完が中心で、実行ファイルの信頼境界を作っていない。
- `ResolveFFmpegPath` は存在確認を行うが、絶対パスや相対パスの任意実行ファイルも許可し、さらに呼び出し側は解決済みパスではなく設定値を実行している。
- テスト撮影またはスケジュール撮影で、設定由来の実行ファイルが利用者権限で起動され得る。

## 期待する挙動

- レポートの意味、影響、根拠、対応方針が日本語で説明される。
- 現在HEADで同種の経路が残っているか分かる。
- 今後修正する場合の要点が分かる。

## 受け入れ条件

- [x] findingの要旨を説明する。
- [x] なぜhigh扱いかを説明する。
- [x] レポート内の検証根拠を説明する。
- [x] 現在HEADでの残存状況を確認する。
- [x] 提案される修正方針を説明する。

## 調査結果

- findingは、外部JSON設定の `autoCapture.stream.ffmpegPath` / `inputArgs` が `exec.CommandContext` へ到達し、任意ローカル実行ファイルを攻撃者指定の引数で起動できる点を指摘している。
- 現在HEADではStream撮影の主経路はSpout helper中心へ移行しているが、互換用の `legacyFfmpegPath` / `legacyInputArgs` と `captureStreamFrameWithFFmpeg` が残っている。
- `src/internal/appcore/config.go` では `legacyFfmpegPath` がJSON-backed fieldとして残り、Normalizeは空欄時に `ffmpeg` を補うだけである。
- `src/internal/appcore/autocapture.go` の `ResolveFFmpegPath` は、パス区切りや絶対パスを含む値でも存在すれば許可する。
- 同じ関数でPATH解決した値を返しているが、`captureStreamFrameWithFFmpeg` は戻り値を捨て、設定値の `cfg.LegacyFFmpegPath` をそのまま `exec.CommandContext` に渡している。
- そのため、レポートの元の `ffmpegPath` 名称は現在の設定名と異なるが、同種の「設定由来の任意実行ファイル起動」経路は残っていると判断する。
- 修正方針は、設定由来のffmpeg実行名を `ffmpeg` / `ffmpeg.exe` のような安全な既定名だけに制限し、絶対パス、相対パス、パス区切り、別名コマンドを拒否し、実行時は `ResolveFFmpegPath` の解決済みパスを使うことである。
