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

- [ ] `RELEASE_NOTES.md` の `v0.1.8` 更新内容が今回のbeta内容を含む。
- [ ] `develop` と `v0.1.8-b14` タグをGitHubへpushする。
- [ ] Release workflowが成功する。
- [ ] GitHub Release `v0.1.8-b14` が prerelease として作成される。
- [ ] Release添付ファイル一覧とzip内ファイル一覧を確認する。
