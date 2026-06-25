# Release workflowの権限を最小化する

## 問題

Release workflowでRelease作成に必要な `contents: write` がworkflow全体に付与されている。外部Actionもメジャーバージョンタグ参照であり、サプライチェーンリスクが残る。

## 期待する挙動

GitHub Actionsの権限をjob単位で必要最小限にし、Release作成権限をRelease作成jobへ限定する。

## 受け入れ条件

- build/test相当のjobは `contents: read` で動作する。
- Release作成に必要なjobだけ `contents: write` を持つ。
- Release用secretが不要なjobに露出しない。
- RC/正式Release workflowが従来通り成果物を作成できる。

## 対応内容

- Release workflowをビルド/パッケージ作成jobとGitHub Release作成jobに分けた。
- workflow既定権限とビルドjobを `contents: read` にした。
- GitHub Release作成jobだけ `contents: write` を持つようにした。
