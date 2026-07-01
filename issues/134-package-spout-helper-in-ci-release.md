# CI/ReleaseでSpoutヘルパーをビルド/同梱する

## 問題

Spout受信ヘルパーを追加しても、GitHub ActionsのCI/Releaseでビルド、テスト、zip同梱、ライセンス同梱ができなければ、v0.1.8 RCでユーザーが実機確認できない。

## 期待する挙動

Windows CIとRelease workflowでSpoutヘルパーをビルドし、配布zipにアプリ本体、Spoutヘルパー、必要DLL、ライセンスを同梱する。
zip内容検証でも必要ファイルの存在と不要ファイルの不在を確認する。

## 受け入れ条件

- Windows CIでSpoutヘルパーをビルドし、単体の `--help` または `--list-senders` が実行できることを確認する。
- Release workflowでSpoutヘルパーをビルドし、配布zipへ同梱する。
- 動的リンクDLLを使う場合は、必要なDLLをアプリ実行フォルダへ同梱する。
- Spout2のBSD 2-Clause Licenseを配布物内またはOSSライセンス表示に含める。
- SpoutRecorderのGPL-3.0成果物や不要なサンプル/デバッグファイルをzipへ混入させない。
- Release zip検証で `ClipForVRChat.exe`、Spoutヘルパー、必要DLL、README、LICENSE、必要なNOTICE/OSSライセンスが存在することを確認する。
- ローカルWindowsビルド手順でもSpoutヘルパーを含めて確認できる。

## 実装メモ

- GitHub Actionsは `windows-latest` で動いているため、MSVC/CMake/Visual Studio Build Toolsを使う案を優先する。
- Wailsビルドとは別ステップでヘルパーをビルドし、`dist/ClipForVRChat/` へコピーする。
- 署名やzipハッシュ生成の対象変更が必要か確認する。

## 対応内容

- v0.1.8実装に合わせて対応済み。詳細は `feat/v0.1.8-resolve-issues` の実装、README、SPEC、RELEASE_NOTESを参照。

## 再監査メモ

- 2026-07-01: [#164](164-audit-v018-completed-items.md) の再監査で未達が見つかったため、完了扱いを取り消して `要対応` に戻した。
