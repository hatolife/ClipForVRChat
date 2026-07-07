# Release署名鍵をビルドジョブから分離する

## 指示

> 現在の作業をできるだけサブエージェントに任せて次の作業を実施 
> 重大な仕様変更・運用判断が必要　の項目について追加情報があります 
> '/home/user/work/ClipForVRChat/reports/security/2026-07-07T17-53-16.203Z/b825b8f570a881919546b22faf5da6cf' '/home/user/work/ClipForVRChat/reports/security/2026-07-07T17-53-16.203Z/394a0dc95c688191a1617b65fcb2befd' '/home/user/work/ClipForVRChat/reports/security/2026-07-07T17-53-16.203Z/e2dbfdd09f68819187b842db37a4421f' '/home/user/work/ClipForVRChat/reports/security/2026-07-07T17-53-16.203Z/958fd884b9588191ab138b574a6bd5eb' '/home/user/work/ClipForVRChat/reports/security/2026-07-07T17-53-16.203Z/69b57f2de0b081919021d7b078b7778c' '/home/user/work/ClipForVRChat/reports/security/2026-07-07T17-53-16.203Z/fa5bc68dae4c8191ad9140983e5ac66f' '/home/user/work/ClipForVRChat/reports/security/2026-07-07T17-53-16.203Z/f23b0a6320c481919bc02541a59d01da' 
> ７件について個別にissue作成 それぞれ対応方針を３パターンissueに記載

## 文脈

追加情報 `f23b0a6320c481919bc02541a59d01da` は、Release workflowで依存関係インストール・ビルド・テストを実行した同一jobの後段にGPG署名秘密鍵を渡すと、先行ステップがrunner状態を汚染して署名鍵を奪えると指摘している。

## 解釈

署名鍵はビルド済みartifactに署名する最小jobへ分離し、依存関係やリポジトリコードを実行したrunner状態から切り離す必要がある。
Release workflow変更のため、配布AssetとRelease本文の検証も受け入れ条件に含める。

## 問題

- signing secretが、npm/go/wails等を実行した後の同一runner jobに渡される。
- 先行ステップが `GITHUB_PATH` やshell hook、偽 `gpg` で後段を汚染できる。
- 署名鍵が漏れると、以後の悪意あるexeが公式署名に見える。

## 期待する挙動

ビルドjobは署名鍵を持たず、署名jobはビルド済みartifactを取得して最小限のコマンドだけで署名する。
公開Assetの種類とzip内ファイル一覧は既存リリース運用に合う。

## 対応方針案

- A: build/sign/package/releaseをjob分割し、sign jobだけに署名secretを渡す。
- B: 署名だけを保護された別environmentまたはreusable workflowへ分離し、artifact digestを入力として承認後に署名する。
- C: GPG長期鍵を廃止し、Sigstore/keylessまたは外部KMS/HSM署名へ移行する。

## 方針評価

- A: 現行workflowを保ったままtrust boundaryを分離できるため、短期で最も現実的。
- B: セキュリティ境界はAより強くできるが、GitHub environment運用や承認フローの設計が必要。
- C: 長期鍵をrunnerへ渡さない理想形だが、利用者向け検証手順やRelease asset方針の再設計が必要。

## 不採用の軽減策

同一jobのまま `PATH` 固定、絶対パス `gpg`、環境初期化だけで固める案は、先行ステップが同じrunner状態へ影響できる問題を残す。
補助策としては有効だが、このfindingの最終対策にはしない。

## 推奨方針

Aを短期対応として採用し、Cは将来の理想形として別途検討する。
BはAより強い運用分離を狙えるが、承認フローとRelease手順の設計が必要なため次段階の候補にする。

## 方針決定

2026-07-08: ユーザー判断によりAを採用する。
Release workflowはbuild/sign/package/releaseをjob分割し、署名secretはsign jobだけへ渡す。
同一job内のPATH固定や絶対パスGPG指定だけでは最終対策にしない。
keyless/KMS署名への移行は将来提案として `docs/security-future-recommendations.md` に分離する。

## 受け入れ条件

- [ ] ビルドjobに `CI_RELEASE_GPG_PRIVATE_KEY` / passphrase / fingerprint が渡らない。
- [ ] 署名jobはビルド済みexe artifactだけを取得し、依存関係インストールやアプリビルドを行わない。
- [ ] 署名対象exeのartifact名とdigestを検証し、別artifact差し替えを検出できる。
- [ ] 署名jobでは絶対パスのGPG等を使い、runner状態の影響を最小化する。
- [ ] Release workflowの公開Assetは通常zip、単一exe署名asc、separated zip、AvatarBeacon unitypackageの運用に合う。
- [ ] Release本文が `RELEASE_NOTES.md` 由来であること、zip内ファイル一覧、不要な公開鍵/sha256個別公開がないことを確認する。
