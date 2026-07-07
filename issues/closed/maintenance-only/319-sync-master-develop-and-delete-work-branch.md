# master/developを同一状態にして作業ブランチを削除する

## 指示

> おわったらmaster,developと一致する状態にして作業ブランチ削除

## 文脈

`v0.1.8-rc3` 作成、Release workflow確認、GitHub Release確認、再発防止メモ追記まで完了した。
現在の作業ブランチは `fix/prefer-embedded-spout-helper` であり、これまでのセキュリティ対応とrc3整理commitを含んでいる。

## 解釈

現在の作業成果を `master` と `develop` の両方へ反映し、両ブランチが同じcommitを指す状態にする。
反映後、作業ブランチ `fix/prefer-embedded-spout-helper` を削除する。

## 問題

- 作業成果が作業ブランチ上に残ったままだと、`master` / `develop` とリリース後の整理状態が分岐する。
- 作業ブランチを残すと、完了済み作業の管理対象が増える。

## 期待する挙動

- `master` と `develop` が同一commitを指す。
- `fix/prefer-embedded-spout-helper` がlocal/remoteから削除される。
- 作業ツリーに今回作業由来の未コミット差分が残らない。

## 受け入れ条件

- [x] `master` が現在の完了commitを指す。
- [x] `develop` が現在の完了commitを指す。
- [x] `origin/master` と `origin/develop` が同一commitを指す。
- [x] local/remoteの `fix/prefer-embedded-spout-helper` が削除される。
- [x] 作業ツリーがcleanである。ただし今回作業外の既存未追跡ファイルは対象外とする。

## 作業メモ

2026-07-08: rc3後の整理状態を通常開発ブランチへ反映するために、localの `master` と `develop` を「旧remote commit `3916c74f88b5fc9e8a558d51d9a6820981d6bd7a`」から「本issue完了記録を含む最終commit」へ進める。
このcommit作成後に `origin/master` と `origin/develop` へ同一commitをpushし、local/remoteの `fix/prefer-embedded-spout-helper` を削除する。
