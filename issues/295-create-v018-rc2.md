# v0.1.8-rc2 release candidateを作成する

## 指示

> rc2作成

## 文脈

`v0.1.8-rc1` 作成後、AvatarBeacon package署名検証、既定Spout helper選択、Spout2取得固定、設定由来ffmpeg実行制限、インポート設定由来の自動処理監視フォルダ確認など、複数のセキュリティ対処が入った。
これらを含めて次のrelease candidate `v0.1.8-rc2` を作成する。

## 解釈

現在のHEADを `v0.1.8-rc2` としてタグ付けし、GitHubへpushしてRelease workflowを走らせる。
GitHub Releaseがprereleaseとして作成され、公開Assetとzip内容が仕様通りであることを確認する。

## 問題

- `v0.1.8-rc1` には直近のセキュリティ対処が含まれていない。
- rc2作成前に、Release notesとissue状態を現在のリリース候補に合わせる必要がある。
- `v0.1.8-rc2` のRelease workflowと成果物確認が必要である。

## 期待する挙動

- `v0.1.8-rc2` タグが直近のリリース準備commitを指す。
- GitHub ActionsのRelease workflowが成功する。
- GitHub Release `v0.1.8-rc2` がprereleaseとして作成される。
- 公開Assetが仕様通りの種類に絞られている。

## 受け入れ条件

- [ ] `RELEASE_NOTES.md` の `v0.1.8` 更新内容がrc2の重要なセキュリティ対処を含む。
- [ ] rc2作成前のSpout2 revision更新要否確認が、rc1後の対応で満たされていることを確認する。
- [ ] リリース準備commitを作成する。
- [ ] `v0.1.8-rc2` タグを作成し、GitHubへpushする。
- [ ] Release workflowが成功する。
- [ ] GitHub Release `v0.1.8-rc2` がprereleaseとして作成される。
- [ ] Release添付ファイル一覧とzip内ファイル一覧を確認する。
