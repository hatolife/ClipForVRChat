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

- [ ] `develop` の未push変更をGitHubへpushする。
- [ ] `v0.1.8-b11` タグを作成してGitHubへpushする。
- [ ] Release workflowが成功する。
- [ ] GitHub Release `v0.1.8-b11` が prerelease として作成される。
- [ ] 公開Assetが通常利用者向けzip、単一exe署名asc、separated zip、AvatarBeacon source zipの4種類であることを確認する。
