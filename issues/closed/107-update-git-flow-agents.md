# Git Flow運用をAGENTS.mdへ明文化する

## 問題

現状の `AGENTS.md` は作業ブランチ運用とRCリリースについて記載しているが、今後の開発基準ブランチを `develop` とする運用が明確ではない。

## 期待する挙動

通常開発は `develop` を基準に `fix/...` や `feat/...` などの作業ブランチを作成し、良好なら `develop` に取り込む。
リリース前は `master` から `release/vX.Y.Z` を作成し、RC確認と安定化を行った上で `master` へ取り込む。

## 受け入れ条件

- `AGENTS.md` に `develop` を通常開発の統合ブランチとするルールが書かれている。
- `AGENTS.md` に作業ブランチを `develop` から作成して `develop` に戻すルールが書かれている。
- `AGENTS.md` にリリース前は `master` から `release/vX.Y.Z` を作成するルールが書かれている。
- `AGENTS.md` に `release/vX.Y.Z` でRC確認を行い、問題なければ `master` と `develop` へ反映するルールが書かれている。
