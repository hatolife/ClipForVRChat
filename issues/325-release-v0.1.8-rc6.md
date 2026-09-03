# v0.1.8-rc6を公開して実機確認可能にする

## 指示

> `v0.1.8-rc6`

## 文脈

自動撮影と設定画面の改善をPR #16で `develop` へ取り込み、Windows CI #294が成功した。ユーザーはRelease成果物で確認するprereleaseとして `v0.1.8-rc6` を明示的に選択した。

## 解釈

- 今回の変更を `v0.1.8-rc6` の対象として記録する。
- `release/v0.1.8` でRelease Notesと対象Issueを更新し、CI成功後にタグを作成する。
- GitHub Releaseがprereleaseとして公開され、規約どおり3種類のAssetを持つことを確認する。
- Camera Lock OFFによるアンカー固定解除など、VRChat実機でしか確認できない項目はRelease成果物で確認する。

## 問題

実装は `develop` へ取り込まれているが、対象バージョンの記録とRelease成果物がまだないため、Windows実機で変更を確認できない。

## 期待する挙動

- `v0.1.8-rc6` のWindows成果物をダウンロードできる。
- GitHub Releaseがprereleaseとして公開される。
- 通常zip、exe署名asc、separated zipの3 Assetだけが公開される。
- 通常zipにAvatarBeacon packageが同梱される。

## 受け入れ条件

- [x] `v0.1.8-rc6` を対象バージョンとしてIssue一覧へ記録する。
- [x] `v0.1.8` Release Notesへ今回の変更を追記する。
- [ ] `release/v0.1.8` のCIが成功する。
- [ ] `v0.1.8-rc6` タグからRelease workflowが成功する。
- [ ] GitHub Releaseがprereleaseで、公開Assetが規約どおり3種類である。
- [ ] 通常zipにAvatarBeacon packageが同梱される。

