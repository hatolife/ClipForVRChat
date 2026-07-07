# Spout helper隣接exe優先セキュリティレポートを解説する

## 指示

> '/home/user/work/ClipForVRChat/reports/security/2026-07-07T17-53-16.203Z/0a3a04c8e55c819191605ec086b69951'
> 日本語で解説

## 文脈

指定されたファイルは、Codex Securityが生成した `Signed single exe can run unsigned adjacent Spout helper` というhigh findingのレポートである。
単一exe配布で `spout-capture.exe` と `SpoutLibrary.dll` を埋め込むようにした一方、Spout helper解決処理が埋め込みhelperより先にアプリ横の `spout-capture.exe` を採用する点が指摘されている。

## 解釈

レポート本文を読み、指摘内容、攻撃シナリオ、根拠、提案される修正方針を日本語で分かりやすく説明する。
この依頼ではコード修正までは行わず、現時点の理解と対応方針を整理する。

## 問題

- 署名済み単一exeを利用しても、同じフォルダにある未署名の `spout-capture.exe` が優先実行され得る。
- 既定helper名のまま使う通常利用者経路で、埋め込み済みhelperの信頼性が迂回される可能性がある。
- helper確認、sender一覧取得、Stream撮影の各経路で解決済みhelperが実行される。

## 期待する挙動

- レポートの意味、影響、根拠、対応方針が日本語で説明される。
- 今後修正する場合の要点が分かる。

## 受け入れ条件

- [x] findingの要旨を説明する。
- [x] なぜhigh扱いかを説明する。
- [x] レポート内の検証根拠を説明する。
- [x] 提案される修正方針を説明する。

## 調査結果

- findingは、`ResolveSpoutHelperPath` が既定名 `spout-capture.exe` でも埋め込みhelperより先にアプリ横の同名ファイルを返す点を指摘している。
- `CheckSpoutHelper`、`ListSpoutSenders`、Stream撮影処理は、解決されたhelper pathを `exec.CommandContext` で実行する。
- そのため、利用者が `ClipForVRChat-vX.Y.Z-windows-amd64.exe` のPGP署名を確認していても、同じフォルダに攻撃者が用意した `spout-capture.exe` があれば、その未署名helperが実行される可能性がある。
- 修正方針は、埋め込みhelperが利用できるビルドでは既定helper名に対して埋め込みhelperを優先し、外部helperはユーザーが明示的に設定したパス、または埋め込みhelperが存在しない分離版の補助経路に限定することである。
