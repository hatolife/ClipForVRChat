# 外部AvatarBeacon unitypackageのPGP署名をRelease workflowで検証する

## 指示

> https://github.com/hatolife/AvatarBeacon/releases/tag/v0.0.1のunitypackageをrelease-signing@hato.lifeで署名しました これを検証する形にしてほしい

## 文脈

ClipForVRChatのRelease workflowは、AvatarBeacon別リポジトリの `AvatarBeacon_v0.0.1.unitypackage` をRelease時にダウンロードし、通常利用者向けzipと分離版zipへ同梱している。
前回のセキュリティレポートでは、外部Release assetを検証せずに公式配布zipへ入れるサプライチェーン上の問題が指摘された。

## 解釈

AvatarBeacon `v0.0.1` Releaseには `AvatarBeacon_v0.0.1.unitypackage.asc` が追加されているため、ClipForVRChatのRelease workflowでpackage本体と署名ファイルを取得し、`release-signing@hato.life` の信頼済みfingerprintに対応する公開鍵で検証してからpackageへ進む。

## 問題

- 現状は `.unitypackage` の存在と非空だけを確認しており、署名検証していない。
- ダウンロード後SHA256は記録しているが、改ざん防止の検証条件としては使っていない。
- exe署名はzip内のAvatarBeacon packageを保護しない。

## 期待する挙動

- Release workflowが `AvatarBeacon_v0.0.1.unitypackage.asc` も取得する。
- Release workflowが `release-signing@hato.life` の公開鍵fingerprintを固定して確認する。
- PGP署名検証に成功したAvatarBeacon packageだけが通常zipと分離版zipへ同梱される。

## 受け入れ条件

- [x] `.github/workflows/release.yml` でAvatarBeacon package署名をダウンロードする。
- [x] `.github/workflows/release.yml` で公開鍵fingerprintとUIDを確認してから `gpg --verify` する。
- [x] build metadataまたは仕様ドキュメントで署名検証済みであることが分かる。
- [x] Release zipと公開Assetの許可リストに不要な `.asc` を混入させない。

## 対応内容

- Release workflowで `AvatarBeacon_v0.0.1.unitypackage.asc` をダウンロードし、非空確認とSHA256記録を行うようにした。
- `keys.openpgp.org` から `release-signing@hato.life` の公開鍵を固定fingerprintで取得し、import後にもfingerprintとUIDを確認してから `gpg --verify` するようにした。
- `VALIDSIG` のfingerprintも確認し、期待したrelease signing key以外の署名ではbuildを止めるようにした。
- build metadataへAvatarBeacon package署名のURL、SHA256、検証fingerprintを記録するようにした。
- 検証後の `.unitypackage.asc` は最終 `dist` から削除し、Release zipや公開Assetへ混入しないようにした。

## 確認

- AvatarBeacon `v0.0.1` Releaseのasset一覧に `AvatarBeacon_v0.0.1.unitypackage` と `AvatarBeacon_v0.0.1.unitypackage.asc` があることを確認した。
- ローカルで実assetと署名を取得し、追加したworkflow相当のfingerprint/UID/PGP署名検証が成功することを確認した。
