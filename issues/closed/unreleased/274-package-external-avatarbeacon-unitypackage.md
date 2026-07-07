# 外部AvatarBeacon unitypackageを配布zipへ同梱する

## 指示

> AvatarBeaconを別リポジトリに分けたのでciを変更したい
> assetsにはAvatarBeaconのsourceを含めない
> 当面下記のデータを配布物として使用したい
> https://github.com/hatolife/AvatarBeacon/releases/download/v0.0.1/AvatarBeacon_v0.0.1.unitypackage
> これを
> ClipForVRChat-v***-windows-amd64.zip 
> ClipForVRChat-v***-windows-amd64-separated.zip 
> などに含めるようにしたい

## 文脈

AvatarBeaconは別リポジトリへ分離済みで、ClipForVRChat側のRelease workflowがリポジトリ内sourceから `AvatarBeacon-*-source.zip` を生成・公開添付する運用は現状に合わない。

## 解釈

GitHub Releaseの公開AssetからAvatarBeacon source zipを外し、Release workflowで外部Releaseの `.unitypackage` を取得して通常利用者向けzipと分離版zipの両方へ同梱する。

## 問題

- Release workflowが `avatar-gimmicks/AvatarBeacon` からsource zipを生成している。
- GitHub Release公開AssetにもAvatarBeacon source zipを含めている。
- 通常利用者向けzipと分離版zipに、利用者がAvatarBeaconを導入するための `.unitypackage` が入らない。

## 期待する挙動

- Release workflowは `AvatarBeacon_v0.0.1.unitypackage` を外部URLから取得する。
- `ClipForVRChat-v***-windows-amd64.zip` と `ClipForVRChat-v***-windows-amd64-separated.zip` の両方に `.unitypackage` を含める。
- GitHub Releaseの公開AssetはAvatarBeacon source zipを含めない。

## 受け入れ条件

- [x] Release workflowのAvatarBeacon source zip生成・公開添付が削除されている。
- [x] Release workflowで外部AvatarBeacon unitypackageをダウンロードし、通常zipと分離zipへ同梱する。
- [x] Release asset検査でunitypackage同梱とsource zip非公開を検証する。
- [x] README/SPECなどのRelease成果物説明が新しい配布形態と一致している。

## 対応メモ

- Release workflowで `AvatarBeacon_v0.0.1.unitypackage` を外部Releaseから取得し、通常zipと分離版zipへコピーするようにした。
- GitHub Release公開Assetの許可リストとupload対象から `AvatarBeacon-*-source.zip` を削除した。
- Release asset検査で、source zip混入禁止、unitypackage同梱、通常zip/分離zipの許可ファイル一覧を確認するようにした。
- CI workflowからAvatarBeacon source zip artifact生成を削除した。
