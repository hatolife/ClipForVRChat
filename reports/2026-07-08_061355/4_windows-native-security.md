# Windows Native Security

| 項目 | 判定 | 確認結果 |
| --- | --- | --- |
| DLL Hijacking | 要確認 | Spout helperは `SpoutLibrary.dll` を同一ディレクトリから使う想定。通常単一exeではhash付きcacheへ埋め込み展開するが、分離版zipでは同梱DLLの真正性確認はユーザーとRelease検査に依存する。 |
| Search Path | 問題あり | `ResolveSpoutHelperPath` が最後に `exec.LookPath("spout-capture.exe")` を使う。`powershell.exe` / `winget` も名前指定で起動している。 |
| LoadLibrary | 要確認 | Go本体で直接LoadLibraryは見当たらない。SpoutLibrary側の内部ロードは未検証。 |
| CreateProcess | 問題あり | Go `exec.Command*` 経由でhelper、ffmpeg、PowerShell、wingetを起動。ffmpegは名前制限済みだがPATH解決。Spout helperと固定コマンドにPATH依存が残る。 |
| ShellExecute | 問題なし | ファイル表示はShell COM API、URLはWails `BrowserOpenURL`。URLはhttpsかつallowlistで制限されている。 |
| COM | 要確認 | Spout helperはWIC/COMを使う。Startup shortcutはPowerShell COM経由で `.lnk` を作成する。COM初期化失敗処理はあるがWindows実機未確認。 |
| Named Pipe | 該当なし | 名前付きパイプ実装は見当たらない。単一起動IPCはlocalhost TCP。 |
| Registry | 該当なし | 直接Registry操作は見当たらない。 |
| ACL | 要確認 | Go側は `0600`/`0700` 相当で保存するが、Windows ACLの実効制御は実機確認が必要。 |
| UAC | 問題なし | Windows起動時にTokenElevationを確認し、管理者権限起動を拒否する。 |
| Temp File | 問題なし | Spout出力は `.tmp` へ保存後に検証してrename。診断作業dirはapp-owned配下またはtempで生成。 |
| TOCTOU | 要確認 | ファイル存在確認後の書き込み/renameは一般的な競合余地がある。管理ディレクトリ内処理が中心で重大な悪用経路は未確認。 |
| Auto Update | 該当なし | 自動更新はなく、GitHub Releasesの通知のみ。 |
| Installer | 該当なし | インストーラはなくzip配布。ffmpeg補助で `winget install ffmpeg` を実行するUI操作がある。 |
| Code Signing | 要確認 | ReleaseではPGP detached signatureを作成する。Windows Authenticode署名は確認できない。 |
| SmartScreen | 要確認 | 実機未確認。zip配布のためSmartScreen表示はReleaseごとに確認が必要。 |
| ASLR | 要確認 | Wails/Go/Windows build artifactをローカルLinuxで確認できない。 |
| DEP | 要確認 | 同上。 |
| CFG | 要確認 | 同上。 |
| Unicode | 問題なし | Windows native APIではUTF-16変換を使用。C++ helperはpathを `std::filesystem::path` とWIC wide stringへ渡す。 |
| Long Path | 要確認 | 明示的なlong path対応は未確認。Windowsのpath長制限に依存する可能性がある。 |
