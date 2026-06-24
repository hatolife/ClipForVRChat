# Windows GUI exeのCLI出力をPowerShellに表示する

## 問題

`ClipForVRChat.exe --version` を PowerShell から実行しても、Wails の Windows GUI サブシステム exe では標準出力が親コンソールへ接続されず、何も表示されない。

## 期待する挙動

Windows配布exeでも、`--version` / `--help` のときは親コンソールへ出力が表示される。

## 受け入れ条件

- PowerShell から `.\ClipForVRChat.exe --version` を実行するとバージョンが表示される。
- PowerShell から `.\ClipForVRChat.exe --help` を実行するとヘルプが表示される。
- 通常のGUI起動では余計なコンソール接続やウィンドウ表示を増やさない。

## 判断

PowerShell/cmd/Git Bashで出力を次のプロンプトより前にきれいに表示するには、Windows側でコンソールサブシステムexeとしてビルドする必要がある。しかしその方式では、通常のGUIアプリ起動時にコンソールが表示される可能性がある。

ClipForVRChatは通常のGUI起動体験を優先するため、`-windowsconsole` は使わない。`--version` / `--help` は親コンソールへ出力するが、GUIサブシステムexeの仕様によりプロンプト表示後に出力される場合がある。
