# 開発ブランチ運用とRCリリースを整備する

## 問題

現状の開発運用は `master` への直接作業を前提にしやすく、リリース前確認用の release candidate を正式Releaseと区別して発行するルールも明文化されていない。

## 期待する挙動

- 通常開発は作業ブランチを切って進め、`master` はリリース可能な状態を保つ。
- `vX.Y.Z-rcW` 形式のタグを push した場合、GitHub Release は prerelease として作成される。
- RC専用のRelease noteがない場合は、対応する正式バージョンのRelease noteを使える。
- `vX.Y.Z` 形式のタグは通常の正式Releaseとして作成される。
- `vX.Y.Z` / `vX.Y.Z-rcW` 以外の `v*` タグは、検証用として draft Release に留める。

## 受け入れ条件

- `AGENTS.md` に開発ブランチ運用とリリース手順が記載されている。
- Release workflow が `vX.Y.Z-rcW` タグを prerelease として扱う。
- Release note 抽出が `vX.Y.Z-rcW` から `vX.Y.Z` へ fallback できる。
- CI workflow と Release workflow の構文検査が通る。
