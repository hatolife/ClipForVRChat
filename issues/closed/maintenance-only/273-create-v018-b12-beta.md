# v0.1.8-b12 betaを作成する

## 指示

> betaつくって

## 文脈

`v0.1.8-b11` は既に作成済みで、GitHub Releaseも prerelease として公開されている。`v0.1.8-b11` 以降、`develop` にはAvatarBeaconの専用リポジトリ/サブモジュール化、AvatarBeacon README/ライセンス/基準Poseの調整、exe引数暗号化、自動撮影Webhook URLの修正、VRChat写真/自動撮影の重複投稿抑止などが入っている。

## 解釈

現在の `develop` の変更を含めて、次の beta タグ `v0.1.8-b12` を作成し、GitHubへpushしてRelease workflowによる配布物作成を確認する。

## 問題

- `v0.1.8-b11` には直近のAvatarBeaconサブモジュール化、ドキュメント更新、UI/自動処理修正が含まれていない。
- beta配布物を作成するには、タグを現在の変更込みのcommitへ付ける必要がある。

## 期待する挙動

- `v0.1.8-b12` タグが直近変更を含むcommitを指す。
- Release workflowが起動し、GitHub Releaseと規定の配布Assetが作成される。
- GitHub Releaseは prerelease として作成される。

## 受け入れ条件

- [x] `develop` の未push変更をGitHubへpushする。
- [x] `v0.1.8-b12` タグを作成してGitHubへpushする。
- [x] Release workflowが成功する。
- [x] GitHub Release `v0.1.8-b12` が prerelease として作成される。
- [x] 公開Assetが通常利用者向けzip、単一exe署名asc、separated zip、AvatarBeacon source zipの4種類であることを確認する。

## 対応メモ

- `v0.1.8-b12` タグを `0309bbae4d7753162330902c62c908654930687c` に作成し、GitHubへpushした。
- Release workflow run `28867151953` が成功した。
- GitHub Release: https://github.com/hatolife/ClipForVRChat/releases/tag/v0.1.8-b12
- GitHub Releaseは `prerelease: true`、`draft: false` で作成された。
- 公開Assetは以下4件のみであることを確認した。
  - `ClipForVRChat-v0.1.8-b12-windows-amd64.zip`
  - `ClipForVRChat-v0.1.8-b12-windows-amd64.exe.asc`
  - `ClipForVRChat-v0.1.8-b12-windows-amd64-separated.zip`
  - `AvatarBeacon-v0.1.8-b12-source.zip`
