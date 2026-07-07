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

- [x] AGENTS.mdに `vX.Y.Z-a1` / `vX.Y.Z-b1` / `vX.Y.Z-rc1` と `a2` / `b2` / `rc2` の確認運用を記載する。
- [x] AGENTS.mdにVRChat現在バージョン更新時はSpout2更新要否を積極的に確認する方針を記載する。
- [x] Spout2取得をhash検証付きarchiveへ変更する。
- [x] Release build metadataへSpout2 revision/source/hashを記録する。
- [x] 関連検査を実行する。

## 実装結果

- AGENTS.mdへSpout2 revision更新運用を追加した。
- Spout2 `2.007.015` が指す `fc63ac2c16918ffa1d3384cddfba2bbc70650bb4` のarchiveを取得対象にした。
- archive SHA256は `087730642decf46d71fa1a96c1f63f4f6a99feb5e82e9b911b6b70c65d8a6305` としてCMakeの `URL_HASH` に固定した。
- Release build metadataの `spout.spout2` にrevision、upstream tag、archive URL、archive SHA256を記録するようにした。

## 検証

- `git ls-remote https://github.com/leadedge/Spout2.git refs/tags/2.007.015 refs/tags/2.007.015^{}` でtag revisionを確認した。
- `curl` と `sha256sum --check` でarchive hashを確認した。
- `cmake -P` の `file(DOWNLOAD ... EXPECTED_HASH SHA256=...)` でarchive hash検証を確認した。
- `cmake -S tools/spout-capture -B /tmp/clipforvrchat-spout-cmake-check`
- `cmake --build /tmp/clipforvrchat-spout-cmake-check`
- `ctest --test-dir /tmp/clipforvrchat-spout-cmake-check --output-on-failure`
- `go test ./...`

## 残確認

- ローカル環境に `pwsh` がないため、Release workflow内のPowerShell metadata生成処理は実行確認できていない。YAML読み込みとCMake変数抽出相当の確認は実施した。
- Windows専用のSpout helper本体ビルドはGitHub Actions上での確認が必要。
