# Spout helperを本体exeへ埋め込み単一exe配布にする

## 問題

`v0.1.8-rc13` では `ClipForVRChat.exe`、`spout-capture.exe`、`SpoutLibrary.dll` が分かれており、利用者が改竄確認する対象が増えている。
[#170](170-investigate-single-exe-distribution-for-spout.md) の追加調査でC案を採用する方針になったため、通常利用者向けの主導線を単一exeへ戻したい。

## 期待する挙動

Releaseの主成果物は `ClipForVRChat-vX.Y.Z-windows-amd64.exe` とし、本体exeに `spout-capture.exe` と `SpoutLibrary.dll` を埋め込む。
実行時はhash検証付きで管理ディレクトリへ展開し、既存の別プロセス呼び出しを維持する。
検証・切り分け用として、分離版zipも同じRelease Assetsへ残す。

## 受け入れ条件

- [x] 既定設定では、外部 `spout-capture.exe` がexe横にない場合でも埋め込みhelperを展開して利用できる。
- [x] 明示的に `spoutHelperPath` が指定された場合は、外部helperを優先して切り分けできる。
- [x] 展開した `spout-capture.exe` と `SpoutLibrary.dll` は期待SHA-256と一致することを検証する。
- [x] Release workflowで単一exe版、単一exe署名、単一exe sha256、分離版zip、分離版zip sha256、build metadataを生成する。
- [x] Release本文の主リンクを単一exe版へ更新し、分離版zipは検証・切り分け用として説明する。
- [x] CI/Releaseで埋め込みhelper経路と分離版helper経路の最低限の動作確認を行う。
- [x] README/SPEC/Release Notesの配布物説明を単一exe主導線へ更新する。

## 調査メモ

- 既定経路はexe横helperを優先し、存在しない場合に埋め込みhelperへフォールバックする。`spoutHelperPath` で明示された外部パスはそのまま外部helperとして扱う。
- 既存の `spout-capture.exe` が exe 横にある分離版zipは、切り分け用としてそのまま動く前提を残す。

## 対応内容

- Go本体に `embeddedspout` build tag 用の埋め込み資産を追加し、Releaseビルド時だけ `spout-capture.exe` と `SpoutLibrary.dll` を `go:embed` する構成にした。
- `ResolveSpoutHelperPath` を、明示パスは外部helper優先、既定名はexe横helper、埋め込みhelper、PATHの順で解決するようにした。
- 埋め込みhelperはユーザーcache配下のhash付きディレクトリへ展開し、展開前後にSHA-256を検証する。
- CI/Release workflowでhelperを先にビルドし、埋め込み用ディレクトリへステージしてから `wails build -tags embeddedspout` を実行するようにした。
- Release workflowで通常版exe、exe署名、exe sha256、分離版zip、分離版zip sha256、build metadataを生成・添付するようにした。
- README、SPEC、Release Notes、Stream Spout実機確認手順、アプリ内AboutのPGP説明を単一exe主導線へ更新した。

## 要確認

- GitHub Actions上でWindows CIが成功すること。
- `v0.1.8-rcN` Releaseで、Release本文の主リンク、添付ファイル一覧、分離版zip内ファイル一覧が仕様通りであること。
- Windows実機で、通常版exe単体から内蔵helperが展開され、Stream Camera(Spout)のsender一覧取得とテスト撮影ができること。
