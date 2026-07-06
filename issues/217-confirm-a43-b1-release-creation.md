# v0.1.8-a44とv0.1.8-b2でRelease workflowを確認する

## 問題

旧RCタグをalpha/betaタグへ移した際、Release workflowを一時停止してタグをpushしたため、`v0.1.8-a43` と `v0.1.8-b1` のGitHub Releaseは作成されなかった。再pushではRelease workflowは起動したが、タグが指す旧commit上のworkflowが `a/b` タグに対応しておらず失敗した。

## 期待する挙動

Release workflowのCI設定変更を含む新しい `v0.1.8-a44` と `v0.1.8-b2` タグを作成し、alpha/betaタグでもGitHub Releaseと配布成果物が作成されることを確認する。

## 受け入れ条件

- Release workflowが `vX.Y.Z-aW` / `vX.Y.Z-bW` / `vX.Y.Z-rcW` / `vX.Y.Z` の全形式を受け付ける。
- `v0.1.8-a44` と `v0.1.8-b2` のタグがCI設定変更を含むcommitを指す。
- `v0.1.8-a44` と `v0.1.8-b2` のRelease workflowが成功する。
- `v0.1.8-a44` と `v0.1.8-b2` のGitHub Releaseがprereleaseとして作成される。
- Release添付ファイルが規定の4種類に絞られている。

## 作業メモ

- 2026-07-07: `gh release view v0.1.8-a43` と `gh release view v0.1.8-b1` はどちらも `release not found`。
- 2026-07-07: `gh run list --workflow Release` では、タグ再分類後の `v0.1.8-a43` / `v0.1.8-b1` に対応する新規runは見当たらない。タグpush時にRelease workflowを一時停止していたため、Releaseは自然作成されていない。
- 2026-07-07: ユーザー指示により、`v0.1.8-a43` と `v0.1.8-b1` をremoteから削除して再pushし、Release workflowが起動してRelease作成されるか確認する。
- 2026-07-07: 再pushにより `v0.1.8-a43` / `v0.1.8-b1` のRelease workflowは起動したが、どちらも `Prepare version` で失敗した。旧commit上のworkflowが `vX.Y.Z` / `vX.Y.Z-rcN` だけを許可しており、`a/b` タグを受け付けないため。
- 2026-07-07: `a43/b1` の旧commitを無理に使わず、CI設定変更を含む新しい `v0.1.8-a44` / `v0.1.8-b2` を作成して確認する方針へ変更する。
