# Release Security Checklist

## リリース可否判定

- [ ] SEC-20260708-001 を修正する、または分離版/検証版の既知リスクとして明記する。
- [ ] CIの `npm audit --omit=dev` が成功している。
- [ ] CIの `govulncheck ./...` が成功している。
- [ ] Windows buildとRelease workflowが成功している。

## ブロッカー項目

- Critical/High findingが新規に出た場合はRelease停止。
- Discord Webhook URLや署名secretがログ、artifact、Release assetに露出した場合はRelease停止。
- 通常利用者向けzipに仕様外ファイルが混入した場合はRelease停止。

## 修正必須項目

- Spout helper default解決のPATH fallbackを廃止または明確に抑止する。
- Release workflowのasset allowlist検証が通ること。
- AvatarBeacon package署名検証が通ること。

## 修正推奨項目

- `powershell.exe` / `winget` を信頼済み絶対パスで起動する。
- GitHub Actionsをcommit SHA pinningへ移行する。
- SBOM生成を検討する。

## 既知リスク

- 外部helperをユーザーが明示指定した場合、そのhelperの信頼性はユーザー判断に依存する。
- PGP署名はexe単体のdetached signatureであり、Windows Authenticode署名ではない。
- SmartScreen警告は配布実績や署名方式に依存する。

## 署名確認

- [ ] Release assetの `.exe.asc` が生成されている。
- [ ] 固定fingerprint `BE40 AA8D 082F 493F 613B C072 21DC 3486 1B40 E77D` と照合している。
- [ ] `gpg --verify` が成功している。

## ハッシュ確認

- [ ] workflow内のartifact digest検証が成功している。
- [ ] zip作成後の内部SHA256検証が成功している。
- [ ] Spout helper/DLLとAvatarBeacon packageのmetadata hashが記録されている。

## SBOM確認

- [ ] 現状はSBOMなし。必要ならCycloneDX等で生成する。

## 依存関係確認

- [ ] Go moduleが意図しない差分を含まない。
- [ ] frontend lock fileが更新意図どおり。
- [ ] Spout2 revision/SHA256が意図どおり。

## CI/CD確認

- [ ] build/sign/package/release jobの権限が最小化されている。
- [ ] sign jobだけが署名secretを参照している。
- [ ] Release jobだけが `contents: write`。

## インストーラ確認

- 該当なし。zip配布。

## アップデート経路確認

- [ ] アプリはGitHub Releases通知のみで自動更新しない。
- [ ] READMEに公式入手先が記載されている。

## README・ユーザー向け文書確認

- [ ] Webhook URLを秘密として扱う注意がある。
- [ ] 非公式配布元から取得しない注意がある。
- [ ] 診断データ共有時は暗号化済み `.gpg` を使う注意がある。
- [ ] 外部helper/分離版zip利用時の注意がある。

## リリースノートに記載すべき注意点

- Spout helperの同梱/埋め込み方式。
- AvatarBeacon unitypackageの同梱と署名検証。
- PGP署名検証方法。
- 既知制限と実機確認が必要な機能。
