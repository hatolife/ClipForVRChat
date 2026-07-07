# v0.1.8-rc3 release candidateを作成する

## 指示

> 終わったらrc3作成

## 文脈

`v0.1.8-rc2` 作成後、Codex Security findings 19:48版の追加対応として、Discord履歴token非保存、issue secret redaction運用、Release署名secret分離、自動撮影Discord opt-in厳密化、非表示camera設定無効化、非表示タブ重要設定の保存前確認を実装した。
これらを含めて次のrelease candidate `v0.1.8-rc3` を作成する。

## 解釈

現在のHEADを `v0.1.8-rc3` としてタグ付けし、GitHubへpushしてRelease workflowを走らせる。
GitHub Releaseがprereleaseとして作成され、Release本文、公開Asset、zip内容が仕様通りであることを確認する。

## 問題

- `v0.1.8-rc2` には直近のセキュリティ対処が含まれていない。
- Release workflowの署名job分離を変更したため、rc3で実際のRelease workflow成功と成果物確認が必要である。
- rc3作成前に、Release notesとissue状態を現在のリリース候補に合わせる必要がある。

## 期待する挙動

- `v0.1.8-rc3` タグが直近のリリース準備commitを指す。
- GitHub ActionsのRelease workflowが成功する。
- GitHub Release `v0.1.8-rc3` がprereleaseとして作成される。
- GitHub Release本文が `RELEASE_NOTES.md` の `v0.1.8` から作成される。
- 公開Assetが仕様通りの種類に絞られている。

## 受け入れ条件

- [x] `RELEASE_NOTES.md` の `v0.1.8` 更新内容がrc3の重要なセキュリティ対処を含む。
- [x] rc3作成前のSpout2 revision更新要否確認が、rc2前後の固定revision運用で満たされていることを確認する。
- [x] リリース準備commitを作成する。
- [x] `v0.1.8-rc3` タグを作成し、GitHubへpushする。
- [ ] Release workflowが成功する。
- [ ] GitHub Release `v0.1.8-rc3` がprereleaseとして作成される。
- [ ] GitHub Release本文が `RELEASE_NOTES.md` の `v0.1.8` 内容から作成されている。
- [ ] Release添付ファイル一覧とzip内ファイル一覧を確認する。

## 作業メモ

2026-07-08: `1721073 docs(release): prepare v018 rc3` で、rc3作成のために `RELEASE_NOTES.md` の `v0.1.8` 更新内容を「rc2時点のセキュリティ対処」から「署名secret分離、自動撮影Discord opt-in厳密化、重要設定保存前確認、履歴token非保存を含む内容」へ変更した。

2026-07-08: `v0.1.8-rc3` を `1721073a2d207046f8f31aa7335831c59235d090` に署名付きで作成し、GitHubへpushした。
Release workflow `https://github.com/hatolife/ClipForVRChat/actions/runs/28897907603` は、`Build unsigned executable and payload` jobの `Build Spout helper` stepで失敗した。
失敗原因は、Windows headerの `max` macroが `std::numeric_limits<...>::max()` に干渉し、MSVCで `C2589` / `C2059` が発生したことだった。
