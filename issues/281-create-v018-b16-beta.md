# v0.1.8-b16 betaを作成する

## 指示

> 待機カメラ位置
> の
> 詳細設定ボタンが押せない  同様にbeta発行しながら治るまでチェック

> あとなぜかカメラの明るさが最大になる そもそも今まで明るさ設定が構図設定になかったのも問題なので追加　vrchatの明るさの初期値を規定値に採用

> あとマスクで他ユーザーがoffになる

## 文脈

`v0.1.8-b15` は既に作成済みで、現在の `develop` にはその後の変更として、待機カメラ位置の詳細設定ボタン修正、構図ごとの明るさ設定、他ユーザーマスク既定ONが含まれている。

## 解釈

現在の `develop` の変更を含めて、次の beta タグ `v0.1.8-b16` を作成し、GitHubへpushしてRelease workflowによる配布物作成を確認する。

## 問題

- `v0.1.8-b15` には待機カメラ位置の詳細設定ボタン修正、構図ごとの明るさ設定、他ユーザーマスク既定ONが含まれていない。
- beta配布物としてGitHub Releaseとzip/署名/分離zipの作成を確認する必要がある。

## 期待する挙動

- `v0.1.8-b16` タグが直近変更を含むcommitを指す。
- GitHub ActionsのRelease workflowが成功する。
- GitHub Release `v0.1.8-b16` が prerelease として作成される。
- 公開Assetが仕様通りの種類に絞られている。

## 受け入れ条件

- [ ] `RELEASE_NOTES.md` の `v0.1.8` 更新内容が今回のbeta内容を含む。
- [ ] `develop` と `v0.1.8-b16` タグをGitHubへpushする。
- [ ] Release workflowが成功する。
- [ ] GitHub Release `v0.1.8-b16` が prerelease として作成される。
- [ ] Release添付ファイル一覧とzip内ファイル一覧を確認する。
