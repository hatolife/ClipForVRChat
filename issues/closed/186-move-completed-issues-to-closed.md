# 完了済みissueをclosedディレクトリへ移動する

## 問題

`issues/` 直下に完了済みのチケットと未完了・要確認のチケットが混在しており、現在対応が必要なissueを把握しづらい。

## 期待する挙動

完了済みのissueは `issues/closed/` 配下へ移動し、完了済み一覧は `issues/closed/README.md` で管理する。

## 受け入れ条件

- `issues/README.md` に未完了、要対応、要確認のissueだけが残っている。
- `issues/closed/README.md` に完了済みissue一覧が移動されている。
- `issues/closed/README.md` の完了済みissueリンクが同ディレクトリ内のファイルを指している。
- 未完了、要対応、要確認のissueは `issues/` 直下に残っている。
- 今後の完了済みissue移動ルールが `AGENTS.md` に記載されている。
