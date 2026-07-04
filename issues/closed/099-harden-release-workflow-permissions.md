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

## 追加調査 2026-06-29

### 問題

Release作成jobは `contents: write` に限定済みだが、GitHub Release作成に `softprops/action-gh-release@v3` の可変タグ参照を使っている。第三者Actionのタグやメンテナが侵害された場合、write権限を持つRelease作成jobで任意コードが実行されるリスクがある。

### 期待する挙動

write権限を持つRelease作成処理では、可変タグ参照の第三者Actionに依存しない。

### 受け入れ条件

- `contents: write` を持つRelease作成jobで第三者Actionの可変タグ参照を実行しない。
- Release本文、draft/prerelease判定、添付ファイル一覧は従来通り維持する。
- build jobは引き続き `contents: read` のまま動作する。

## 対応内容（2026-06-29）

- `softprops/action-gh-release` を使わず、GitHub CLIでRelease作成・更新・添付ファイルアップロードを行うようにした。
- Release作成jobの `contents: write` はRelease操作に必要な最小権限として維持し、build jobは `contents: read` のまま維持した。
