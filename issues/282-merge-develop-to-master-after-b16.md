# v0.1.8-b16確認後にdevelopをmasterへマージする

## 指示

> よさげなのでmasterにマージ

> push

## 文脈

`v0.1.8-b16` betaで、待機カメラ位置の詳細設定、構図ごとの明るさ、表示対象マスク修正のRelease workflowと配布物確認が完了している。

## 解釈

現在の `develop` を `master` へ通常mergeし、GitHubへpushする。正式リリースタグ作成は指示されていないため行わない。

## 問題

- `master` には `develop` の直近beta確認済み変更がまだ反映されていない。

## 期待する挙動

- `master` が `develop` の `v0.1.8-b16` 確認済み変更を含む。
- `master` がGitHubへpushされる。

## 受け入れ条件

- [ ] `develop` の作業ツリーがcleanであることを確認する。
- [ ] `master` へ `develop` をmergeする。
- [ ] `master` をGitHubへpushする。
- [ ] push後の状態を確認する。
