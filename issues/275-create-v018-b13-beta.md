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

- [ ] `develop` がGitHubへpushされている。
- [ ] `v0.1.8-b13` タグを作成してGitHubへpushする。
- [ ] Release workflowが成功する。
- [ ] GitHub Release `v0.1.8-b13` が prerelease として作成される。
- [ ] 公開Assetが規定の3種類であることを確認する。
- [ ] 通常zipとseparated zipに `AvatarBeacon_v0.0.1.unitypackage` が含まれることを確認する。
