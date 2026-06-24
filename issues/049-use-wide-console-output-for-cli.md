# CLIヘルプをWindowsコンソールで文字化けさせない

## 問題

Windows GUIサブシステムexeから `CONOUT$` にUTF-8バイト列を書き込むと、cmdなどCP932のコンソールで日本語ヘルプが文字化けする。

## 期待する挙動

`--version` / `--help` の出力は、PowerShell/cmd/Git Bash のコンソール設定に依存せず日本語を読める形で表示される。

## 受け入れ条件

- cmdで `ClipForVRChat.exe --help` を実行しても日本語説明が文字化けしない。
- PowerShellで `ClipForVRChat.exe --help` を実行しても日本語説明が文字化けしない。
- 通常のGUI起動ではコンソールを表示しない。

## 確認

- PowerShell、cmd、Git Bash で `--version` / `--help` の出力を確認済み。
- cmd の日本語ヘルプは `WriteConsoleW` により文字化けしないことを確認済み。
- GUIサブシステムexeのため、環境によっては出力が次のプロンプト表示後に出る。この挙動は通常GUI起動でコンソールを表示しないことを優先して許容する。
