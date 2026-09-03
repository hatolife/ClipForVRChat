# Release公開Assetを規約どおり3種類に限定する

## 指示

> releaseで確認できる状態に

## 文脈

`v0.1.8-rc5` のGitHub Releaseには通常利用者向けzip、単一exe署名asc、separated zipに加えて、AvatarBeacon `.unitypackage` が単独Assetとして添付された。リポジトリ規約ではAvatarBeacon packageは配布zip内へ同梱し、ClipForVRChat側の公開Assetには添付しない。

## 解釈

`v0.1.8-rc5` の公開Assetを規約どおり3種類へ修正し、以後のRelease workflowも同じ構成だけを公開する。AvatarBeacon packageの取得・署名検証・zip内同梱は維持する。

## 問題

Release workflowのartifactと`gh release upload`対象にAvatarBeacon packageが残っており、公開Assetの許可リストにも含まれている。このためRelease作成時に規約外の4番目のAssetが公開される。

## 期待する挙動

- GitHub Releaseの公開Assetは通常zip、exe署名asc、separated zipの3種類だけになる。
- AvatarBeacon packageは通常zip内への同梱を継続する。
- 既存の`v0.1.8-rc5`から単独AvatarBeacon Assetを削除する。
- Release workflowの検証と公開処理が3 Asset構成を強制する。

## 受け入れ条件

- [ ] Release artifactから単独AvatarBeacon packageを除外する。
- [ ] 公開Asset許可リストとupload対象を3種類へ限定する。
- [ ] `v0.1.8-rc5`の公開Assetが3種類であることを確認する。
- [ ] Release workflow変更後のCIが成功する。
