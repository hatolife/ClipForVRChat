# v0.1.8-b14 betaを作成する

## 指示

> 終わったらbeta作成

## 文脈

`v0.1.8-b13` は既に作成済みで、現在の `develop` にはその後の変更として、Stream自動撮影後のサムネイル/Discord投稿/カメラ終了方向の修正と、自動撮影通常モードの待機カメラ位置機能が含まれている。

## 解釈

現在の `develop` の変更を含めて、次の beta タグ `v0.1.8-b14` を作成し、GitHubへpushしてRelease workflowによる配布物作成を確認する。

## 問題

- `v0.1.8-b13` には直近のStream自動撮影修正と待機カメラ位置機能が含まれていない。
- beta配布物としてGitHub Releaseとzip/署名/分離zip/AvatarBeacon同梱状態を確認する必要がある。

## 期待する挙動

- `v0.1.8-b14` タグが直近変更を含むcommitを指す。
- GitHub ActionsのRelease workflowが成功する。
- GitHub Release `v0.1.8-b14` が prerelease として作成される。
- 公開Assetが仕様通りの種類に絞られている。

## 受け入れ条件

- [x] `RELEASE_NOTES.md` の `v0.1.8` 更新内容が今回のbeta内容を含む。
- [x] `develop` と `v0.1.8-b14` タグをGitHubへpushする。
- [x] Release workflowが成功する。
- [x] GitHub Release `v0.1.8-b14` が prerelease として作成される。
- [x] Release添付ファイル一覧とzip内ファイル一覧を確認する。

## 作業結果

- `v0.1.8-b14` タグを `73be6f83a4b2ae7183fd64a19c37de471f8df9ee` に作成し、GitHubへpushした。
- Release workflow: https://github.com/hatolife/ClipForVRChat/actions/runs/28877124672
- GitHub Release: https://github.com/hatolife/ClipForVRChat/releases/tag/v0.1.8-b14
- GitHub Releaseは prerelease、draftなし。
- 公開Assetは次の3件。
  - `ClipForVRChat-v0.1.8-b14-windows-amd64.zip`
  - `ClipForVRChat-v0.1.8-b14-windows-amd64.exe.asc`
  - `ClipForVRChat-v0.1.8-b14-windows-amd64-separated.zip`
- 通常zip内ファイル一覧を確認した。
  - `AvatarBeacon_v0.0.1.unitypackage`
  - `ClipForVRChat-v0.1.8-b14-windows-amd64.exe.asc`
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
