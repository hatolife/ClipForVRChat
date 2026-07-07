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

- [x] `RELEASE_NOTES.md` の `v0.1.8` 更新内容がrc2の重要なセキュリティ対処を含む。
- [x] rc2作成前のSpout2 revision更新要否確認が、rc1後の対応で満たされていることを確認する。
- [x] リリース準備commitを作成する。
- [x] `v0.1.8-rc2` タグを作成し、GitHubへpushする。
- [x] Release workflowが成功する。
- [x] GitHub Release `v0.1.8-rc2` がprereleaseとして作成される。
- [x] Release添付ファイル一覧とzip内ファイル一覧を確認する。

## 作業結果

- `v0.1.8-rc2` タグを `732d720dd43bd95e64634701c136fe69b935970b` に署名付きで作成し、GitHubへpushした。
- rc2作成前のSpout2 revision更新要否確認は、rc1後に `289` でSpout2取得を固定commit archiveとSHA256検証へ変更済みであることを確認した。
- Release workflow: https://github.com/hatolife/ClipForVRChat/actions/runs/28891663277
- Branch CI: https://github.com/hatolife/ClipForVRChat/actions/runs/28891664954
- GitHub Release: https://github.com/hatolife/ClipForVRChat/releases/tag/v0.1.8-rc2
- GitHub Releaseは prerelease、draftなし。
- 公開Assetは現行 `274` 方針通り、AvatarBeacon source zipを含めず次の3件。
  - `ClipForVRChat-v0.1.8-rc2-windows-amd64.zip`
  - `ClipForVRChat-v0.1.8-rc2-windows-amd64.exe.asc`
  - `ClipForVRChat-v0.1.8-rc2-windows-amd64-separated.zip`
- 通常zip内ファイル一覧を確認した。
  - `AvatarBeacon_v0.0.1.unitypackage`
  - `ClipForVRChat-v0.1.8-rc2-windows-amd64.exe.asc`
  - `ClipForVRChat.exe`
  - `LICENSE`
  - `README.md`
  - `Release-signing-public-key.url`
- separated zip内ファイル一覧を確認した。
  - `AvatarBeacon_v0.0.1.unitypackage`
  - `ClipForVRChat.exe`
  - `LICENSE`
  - `README.md`
  - `Release-signing-public-key.url`
  - `spout-capture.exe`
  - `Spout2-LICENSE.txt`
  - `SpoutLibrary.dll`
- ローカル確認として次を実行した。
  - `GOCACHE=/tmp/clipforvrchat-go-cache go test ./...` (`src/`)
  - `node scripts/check-frontend-template-literals.mjs`
  - `node scripts/check-wails-api-surface.mjs`
  - `node scripts/check-vue-runtime-template.mjs`
  - `node scripts/check-auto-processing-confirmation.mjs`
  - `node scripts/check-closed-issue-index.mjs`
  - `npm run build` (`src/frontend/`)
