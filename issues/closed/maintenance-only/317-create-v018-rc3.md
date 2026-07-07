# v0.1.8-rc3 release candidateを作成する

## 指示

> 終わったらrc3作成

## 文脈

`v0.1.8-rc2` 作成後、Codex Security findings 19:48版の追加対応として、Discord履歴token非保存、issue secret redaction運用、Release署名secret分離、自動撮影Discord opt-in厳密化、非表示camera設定無効化、非表示タブ重要設定の保存前確認を実装した。
これらを含めて次のrelease candidate `v0.1.8-rc3` を作成する。

## 解釈

現在のHEADを `v0.1.8-rc3` としてタグ付けし、GitHubへpushしてRelease workflowを走らせる。
GitHub Releaseがprereleaseとして作成され、Release本文、公開Asset、zip内容が仕様通りであることを確認する。

## 問題

- `v0.1.8-rc2` には直近のセキュリティ対処が含まれていない。
- Release workflowの署名job分離を変更したため、rc3で実際のRelease workflow成功と成果物確認が必要である。
- rc3作成前に、Release notesとissue状態を現在のリリース候補に合わせる必要がある。

## 期待する挙動

- `v0.1.8-rc3` タグが直近のリリース準備commitを指す。
- GitHub ActionsのRelease workflowが成功する。
- GitHub Release `v0.1.8-rc3` がprereleaseとして作成される。
- GitHub Release本文が `RELEASE_NOTES.md` の `v0.1.8` から作成される。
- 公開Assetが仕様通りの種類に絞られている。

## 受け入れ条件

- [x] `RELEASE_NOTES.md` の `v0.1.8` 更新内容がrc3の重要なセキュリティ対処を含む。
- [x] rc3作成前のSpout2 revision更新要否確認が、rc2前後の固定revision運用で満たされていることを確認する。
- [x] リリース準備commitを作成する。
- [x] `v0.1.8-rc3` タグを作成し、GitHubへpushする。
- [x] Release workflowが成功する。
- [x] GitHub Release `v0.1.8-rc3` がprereleaseとして作成される。
- [x] GitHub Release本文が `RELEASE_NOTES.md` の `v0.1.8` 内容から作成されている。
- [x] Release添付ファイル一覧とzip内ファイル一覧を確認する。

## 作業メモ

2026-07-08: `1721073 docs(release): prepare v018 rc3` で、rc3作成のために `RELEASE_NOTES.md` の `v0.1.8` 更新内容を「rc2時点のセキュリティ対処」から「署名secret分離、自動撮影Discord opt-in厳密化、重要設定保存前確認、履歴token非保存を含む内容」へ変更した。

2026-07-08: `v0.1.8-rc3` を `1721073a2d207046f8f31aa7335831c59235d090` に署名付きで作成し、GitHubへpushした。
Release workflow `https://github.com/hatolife/ClipForVRChat/actions/runs/28897907603` は、`Build unsigned executable and payload` jobの `Build Spout helper` stepで失敗した。
失敗原因は、Windows headerの `max` macroが `std::numeric_limits<...>::max()` に干渉し、MSVCで `C2589` / `C2059` が発生したことだった。

2026-07-08: ユーザー承認により、失敗した `v0.1.8-rc3` タグを削除し、`4df3b617ff486980f7da35d22e2fa619832a9062` へ署名付きで付け直してGitHubへpushした。
`fix(spout): avoid windows max macro collisions` で、Release workflowのMSVC buildを通すために、Spout helperのWindows header includeを「`windows.h` をそのままincludeする」処理から「`NOMINMAX` を定義してから `windows.h` をincludeする」処理へ変更した。

2026-07-08: Release workflow `https://github.com/hatolife/ClipForVRChat/actions/runs/28898287169` が成功した。
`Build unsigned executable and payload`、`Sign executable`、`Package release assets`、`Create GitHub release` がすべて成功した。
GitHub Release `https://github.com/hatolife/ClipForVRChat/releases/tag/v0.1.8-rc3` はprerelease、draftなしで作成された。

2026-07-08: GitHub Release本文が `RELEASE_NOTES.md` の `v0.1.8` 内容から作成され、ダウンロードURL、更新内容、比較URLが `v0.1.8-rc3` 向けに置換されていることを確認した。
公開Assetはworkflowの許可リストどおり、次の4件だった。

- `ClipForVRChat-v0.1.8-rc3-windows-amd64.zip`
- `ClipForVRChat-v0.1.8-rc3-windows-amd64.exe.asc`
- `ClipForVRChat-v0.1.8-rc3-windows-amd64-separated.zip`
- `AvatarBeacon_v0.0.1.unitypackage`

2026-07-08: 通常zip内ファイル一覧を確認した。

- `AvatarBeacon_v0.0.1.unitypackage`
- `ClipForVRChat-v0.1.8-rc3-windows-amd64.exe.asc`
- `ClipForVRChat.exe`
- `LICENSE`
- `README.md`
- `Release-signing-public-key.url`

2026-07-08: separated zip内ファイル一覧を確認した。

- `AvatarBeacon_v0.0.1.unitypackage`
- `ClipForVRChat.exe`
- `LICENSE`
- `README.md`
- `Release-signing-public-key.url`
- `spout-capture.exe`
- `Spout2-LICENSE.txt`
- `SpoutLibrary.dll`
