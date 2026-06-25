# 084 Explorerで生成ファイルを選択表示できない

## 問題

不具合報告用データ生成後や履歴の保存先表示で、Explorerが対象ファイルを選択せずマイドキュメントを開くことがある。
現在は `explorer.exe /select,PATH` をコマンドライン引数で起動しており、パス解釈に失敗すると既定フォルダが開かれる。

## 期待する挙動

- Explorerで対象ファイルの親フォルダを開き、そのファイルを選択状態にする。
- `explorer.exe /select,` のコマンドライン解釈に依存しない。
- Windows Shell APIで選択表示できない場合は明確なエラーを返す。

## 受け入れ条件

- Windowsでは `SHOpenFolderAndSelectItems` を使ってファイル選択表示を行う。
- 存在しないパスやディレクトリは従来通りエラーになる。
- LinuxなどWindows以外ではWindows専用機能としてエラーになる。
