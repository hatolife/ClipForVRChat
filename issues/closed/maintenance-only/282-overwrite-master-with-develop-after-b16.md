# v0.1.8-b16確認後にmasterをdevelopで上書きする

## 指示

> よさげなのでmasterにマージ

> push

> fix/b12-spout-sender-recovery-after-batch は削除

> もうmasterをdevelopで上書きしていいよ

## 文脈

`v0.1.8-b16` betaで、待機カメラ位置の詳細設定、構図ごとの明るさ、表示対象マスク修正のRelease workflowと配布物確認が完了している。

## 解釈

現在の `develop` を `master` へ反映し、GitHubへpushする。通常mergeではなく、追加指示に従って `master` を `develop` と同一の履歴へ上書きする。正式リリースタグ作成は指示されていないため行わない。不要になった `fix/b12-spout-sender-recovery-after-batch` branch はlocal/remoteとも削除する。

## 問題

- `master` には `develop` の直近beta確認済み変更がまだ反映されていない。
- 通常mergeでは履歴差分が大きく、追加指示により `develop` で `master` を上書きする方針へ変更された。

## 期待する挙動

- `master` が `develop` と同じ履歴になる。
- `master` がGitHubへpushされる。
- `fix/b12-spout-sender-recovery-after-batch` branch が削除される。

## 受け入れ条件

- [x] `develop` の作業ツリーがcleanであることを確認する。
- [x] `master` を `develop` と同一の履歴へ上書きする。
- [x] `master` をGitHubへpushする。
- [x] `fix/b12-spout-sender-recovery-after-batch` branch をlocal/remoteから削除する。
- [x] push後の状態を確認する。

## 結果

- `origin/master` を `develop` で force-with-lease 更新した。
- ローカル `master` も `develop` と同じコミットへ更新した。
- `develop`、`master`、`origin/develop`、`origin/master` が同じコミットを指すことを確認した。
- `fix/b12-spout-sender-recovery-after-batch` はlocal/remoteとも削除済み。
