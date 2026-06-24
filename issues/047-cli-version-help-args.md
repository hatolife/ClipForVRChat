# CLIでversion/help引数に対応する

## 問題

`ClipForVRChat.exe --version` のようにコマンドラインからバージョンを確認する手段がない。ヘルプ表示も明示的に用意されていない。

## 期待する挙動

`go-arg` を使って、現時点では `--version` と `--help` のみを処理する。通常起動や既存の画像/config.json 引数処理は維持する。

## 受け入れ条件

- `ClipForVRChat.exe --version` でアプリ名とバージョンを標準出力へ表示して終了する。
- `ClipForVRChat.exe --help` で対応引数が分かるヘルプを標準出力へ表示して終了する。
- `version`、`help` 以外の新しい引数機能は追加しない。
- `--version` / `--help` ではアプリの起動ロックやGUI起動を行わない。
