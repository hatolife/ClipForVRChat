# 既定Spout helperは隣接exeより埋め込みhelperを優先する

## 指示

> 対応して

## 文脈

直前に解説した `Signed single exe can run unsigned adjacent Spout helper` findingへの対応。
単一exe配布では `spout-capture.exe` と `SpoutLibrary.dll` を本体に埋め込むが、現在の `ResolveSpoutHelperPath` は既定helper名でもアプリ横の `spout-capture.exe` を埋め込みhelperより先に返す。

## 解釈

署名済み単一exeの既定動作では埋め込みhelperを優先し、同じフォルダに置かれた未署名 `spout-capture.exe` が暗黙に実行されないようにする。
ユーザーが明示的に外部helperパスを設定した場合の互換性は維持する。

## 問題

- 署名済み単一exeを使っても、同じフォルダの未署名 `spout-capture.exe` が優先され得る。
- helper確認、sender一覧取得、Stream撮影の各経路でそのhelperが実行される。
- PGP署名確認済みexeの信頼境界が、隣接バイナリで迂回される。

## 期待する挙動

- 埋め込みhelperが利用可能なビルドでは、helper未指定または既定名 `spout-capture.exe` の場合に埋め込みhelperを使う。
- 明示的な外部helperパス指定は従来通り使える。
- 埋め込みhelperがないビルドでは、分離版互換として従来通り隣接helperやPATH解決を使える。

## 受け入れ条件

- [x] `ResolveSpoutHelperPath` が埋め込みhelper利用可能時の既定helper名で隣接exeを優先しない。
- [x] 明示的な外部helperパス指定の既存テストが通る。
- [x] 回帰テストで埋め込み優先の解決順を確認する。
- [x] 関連Goテストが通る。

## 実施内容

- `ResolveSpoutHelperPath` で、helper未指定または既定名 `spout-capture.exe` の場合は、埋め込みhelperが利用可能なら隣接exe確認より先に埋め込みhelperを返すようにした。
- ユーザーが明示的にパス区切りを含むhelperパスや絶対パスを指定した場合は従来通りそのパスを使う。
- 埋め込みhelperがないビルドでは、分離版互換として隣接helper解決を維持した。
- 回帰テストとして、埋め込みhelper利用可能時に隣接helper lookupが呼ばれないことと、埋め込みhelperなしでは隣接helperを使うことを追加した。
