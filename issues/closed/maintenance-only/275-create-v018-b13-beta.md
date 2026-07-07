# v0.1.8-b13 betaを作成する

## 指示

> 終わったらbeta作成

## 文脈

`v0.1.8-b12` は既に作成済みで、現在の `develop` にはその後の変更として、使い方画面の本文化、AvatarBeacon v0.0.1 unitypackage同梱Release workflow変更、関連issue整理が含まれている。

## 解釈

現在の `develop` の変更を含めて、次の beta タグ `v0.1.8-b13` を作成し、GitHubへpushしてRelease workflowによる配布物作成を確認する。

## 問題

- `v0.1.8-b12` には直近の使い方画面実装とRelease workflow変更が含まれていない。
- beta配布物を作成するには、タグを現在の変更込みのcommitへ付ける必要がある。

## 期待する挙動

- `v0.1.8-b13` タグが直近変更を含むcommitを指す。
- Release workflowが成功し、GitHub Releaseがprereleaseとして作成される。
- 公開Assetが通常利用者向けzip、単一exe署名asc、separated zipの3種類に絞られ、zip内に `AvatarBeacon_v0.0.1.unitypackage` が含まれる。

## 受け入れ条件

- [x] `develop` がGitHubへpushされている。
- [x] `v0.1.8-b13` タグを作成してGitHubへpushする。
- [x] Release workflowが成功する。
- [x] GitHub Release `v0.1.8-b13` が prerelease として作成される。
- [x] 公開Assetが規定の3種類であることを確認する。
- [x] 通常zipとseparated zipに `AvatarBeacon_v0.0.1.unitypackage` が含まれることを確認する。

## 対応メモ

- `develop` を `5151103fa852cfa9d17258905988a0bcfca220de` までpushした。
- `v0.1.8-b13` タグを `5151103fa852cfa9d17258905988a0bcfca220de` に作成し、GitHubへpushした。
- Release workflow run `28869782471` が成功した。
- develop CI run `28869774200` が成功した。
- GitHub Release: https://github.com/hatolife/ClipForVRChat/releases/tag/v0.1.8-b13
- GitHub Releaseは prerelease、draftなし。
- 公開Assetは次の3件。
  - `ClipForVRChat-v0.1.8-b13-windows-amd64.zip`
  - `ClipForVRChat-v0.1.8-b13-windows-amd64.exe.asc`
  - `ClipForVRChat-v0.1.8-b13-windows-amd64-separated.zip`
- 通常zipの中身は次の6ファイル。
  - `AvatarBeacon_v0.0.1.unitypackage`
  - `ClipForVRChat-v0.1.8-b13-windows-amd64.exe.asc`
  - `ClipForVRChat.exe`
  - `LICENSE`
  - `README.md`
  - `Release-signing-public-key.url`
- separated zipの中身は次の8ファイル。
  - `AvatarBeacon_v0.0.1.unitypackage`
  - `ClipForVRChat.exe`
  - `LICENSE`
  - `README.md`
  - `Release-signing-public-key.url`
  - `Spout2-LICENSE.txt`
  - `SpoutLibrary.dll`
  - `spout-capture.exe`
