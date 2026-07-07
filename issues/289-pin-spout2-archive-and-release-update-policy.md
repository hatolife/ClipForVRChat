# Spout2取得をhash検証付きarchiveへ固定し更新運用を明文化する

## 指示

> vX.Y.Zのa1,b1,rc1タイミングで更新すべきか調査して更新すべきと判断した場合に更新するa2,b2,rc2では直前バージョンでチェックが正しく行われたかを確認しダメならそのタイミングで確認と更新処理をする vrchatの現在バージョンに更新があれば積極的に更新してほしい としたい　AGENTS.mdに記載　そのうえで2. hash検証付きアーカイブまで対応

## 文脈

Codex Securityの `Unpinned Spout2 FetchContent in release build` findingへの対応として、Spout2取得を可変タグから不変commit/hash検証へ移行する方針を整理済み。
ユーザーは、Spout2 revision更新判断を各リリース段階の初回タグで行い、必要なら次タグで反映する運用をAGENTS.mdへ記載したうえで、hash検証付きarchive方式まで実装することを求めている。

## 解釈

AGENTS.mdにSpout2 revision更新判断のタイミングとVRChat更新時の積極確認方針を追記する。
そのうえで、`tools/spout-capture/CMakeLists.txt` のSpout2取得を `GIT_TAG 2.007.015` から、不変commitのGitHub archive URLと `URL_HASH SHA256=...` へ変更する。
Release build metadataにもSpout2 revision/source/hashを記録し、後から成果物の依存元を追跡できるようにする。

## 問題

- 現在のSpout2取得はタグ名指定で、タグ移動や依存配送経路の侵害に弱い。
- Spout2 revisionをいつ更新するかの運用がAGENTS.mdに明文化されていない。
- Release metadataにSpout2取得元revision/hashが記録されていない。

## 期待する挙動

- AGENTS.mdにSpout2 revision更新確認のタイミングが明記される。
- Spout2取得は不変commit archive URLとSHA256 hashで検証される。
- Release metadataにSpout2 revision/source/hashが記録される。
- CI/ReleaseのSpout helper buildが引き続き動作する。

## 受け入れ条件

- [ ] AGENTS.mdに `vX.Y.Z-a1` / `vX.Y.Z-b1` / `vX.Y.Z-rc1` と `a2` / `b2` / `rc2` の確認運用を記載する。
- [ ] AGENTS.mdにVRChat現在バージョン更新時はSpout2更新要否を積極的に確認する方針を記載する。
- [ ] Spout2取得をhash検証付きarchiveへ変更する。
- [ ] Release build metadataへSpout2 revision/source/hashを記録する。
- [ ] 関連検査を実行する。
