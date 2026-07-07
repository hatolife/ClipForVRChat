# v0.1.8-b11 betaを作成する

## 指示

終わったらbeta作成

## 文脈

`v0.1.8-b10` は既に作成済みで、GitHub Releaseも prerelease として公開されている。`v0.1.8-b10` 以降、`develop` には設定タブ順の変更、OSSライセンス表示の整理、使い方画面文面草案の更新、VRChat写真/自動撮影まわりの修正記録などが入っている。

## 解釈

現在の `develop` の変更を含めて、次の beta タグ `v0.1.8-b11` を作成し、GitHubへpushしてRelease workflowによる配布物作成を確認する。

## 問題

- `v0.1.8-b10` には直近のUI/ドキュメント/設定タブ順/ライセンス表示関連の変更が含まれていない。
- beta配布物を作成するには、タグを現在の変更込みのcommitへ付ける必要がある。

## 期待する挙動

- `v0.1.8-b11` タグが直近変更を含むcommitを指す。
- Release workflowが起動し、GitHub Releaseと規定の配布Assetが作成される。
- GitHub Releaseは prerelease として作成される。

## 受け入れ条件

- [x] `develop` の未push変更をGitHubへpushする。
- [x] `v0.1.8-b11` タグを作成してGitHubへpushする。
- [x] Release workflowが成功する。
- [x] GitHub Release `v0.1.8-b11` が prerelease として作成される。
- [x] 公開Assetが通常利用者向けzip、単一exe署名asc、separated zip、AvatarBeacon source zipの4種類であることを確認する。

## 対応メモ

- `v0.1.8-b11` タグを `163d961447fdf38479a31c7fe6d83eb7a0ba1ec4` に作成し、GitHubへpushした。
- Release workflow run `28854546164` が成功した。
- GitHub Release: https://github.com/hatolife/ClipForVRChat/releases/tag/v0.1.8-b11
- GitHub Releaseは `prerelease: true`、`draft: false` で作成された。
- 公開Assetは以下4件のみであることを確認した。
  - `ClipForVRChat-v0.1.8-b11-windows-amd64.zip`
  - `ClipForVRChat-v0.1.8-b11-windows-amd64.exe.asc`
  - `ClipForVRChat-v0.1.8-b11-windows-amd64-separated.zip`
  - `AvatarBeacon-v0.1.8-b11-source.zip`
