# Release Security Checklist

調査日: 2026-06-25

## ビルド環境

- [ ] Go versionがworkflowとローカル手順で一致している。
- [ ] Node.js versionがworkflowとローカル手順で一致している。
- [ ] Wails CLI versionが `WAILS_VERSION` と一致している。
- [ ] Windows runnerでWails buildが成功する。
- [ ] ローカルビルド成果物をRelease成果物と混同しない。

## 依存関係

- [ ] `go mod tidy` 後に不要差分がない。
- [ ] `npm ci` が成功する。
- [ ] `npm audit --omit=dev` が0件。
- [ ] `govulncheck ./...` が0件。
- [ ] `go list -m -u all` の重要更新候補を確認した。
- [ ] `npm outdated --long` のメジャー更新候補を確認した。

## テスト

- [ ] `go test ./...` が成功する。
- [ ] `npm run build` が成功する。
- [ ] 診断データzip作成・復号・zip展開のテストが成功する。
- [ ] Windows実機でCLI、GUI起動、Explorer選択表示、診断データ作成を確認する。

## 静的解析

- [ ] `gosec ./...` が0件。
- [ ] `rg` で秘密情報やdebug残存を確認した。
- [ ] `innerHTML` / `v-html` などHTML直接挿入がないことを確認した。

## 動的解析

- [ ] 大きな画像、壊れた画像、非画像ファイルの処理を確認する。
- [ ] Webhook URL未設定・不正・失効時の表示を確認する。
- [ ] Discord投稿OFF時に外部送信されないことを確認する。
- [ ] ローカル保存OFF時にoutputへ書かれないことを確認する。
- [ ] QRコードURL検出ON/OFFの挙動を確認する。

## 署名

- [ ] exe detached PGP署名が作成される。
- [ ] Release workflow内で署名検証が成功する。
- [ ] 公開鍵URLがREADME/About/Release zip内で整合している。
- [ ] 署名鍵secretが必要なjob以外へ露出していない。

## ハッシュ

- [ ] zip SHA-256が作成される。
- [ ] `.sha256` のファイル名とzip名が一致する。
- [ ] ユーザー向けREADMEに検証手順がある。

## SBOM

- [ ] SBOM作成要否を判断する。
- [ ] 作成する場合はGo/npm双方の依存関係を含める。
- [ ] Release assetへSBOMを添付する。

## 配布物

- [ ] zip内に `ClipForVRChat.exe`、`README.md`、`LICENSE`、`Release-signing-public-key.url` が含まれる。
- [ ] zip内に `.asc` や不要な公開鍵実体ファイルが混入していない。
- [ ] Release assetsにzip、sha256、exe.ascが添付される。
- [ ] Release本文が `RELEASE_NOTES.md` の対象バージョンから作成される。

## 更新経路

- [ ] アプリは自動更新しない。
- [ ] 更新通知のGitHub/BOOTHリンクが正しい。
- [ ] 不正なRelease APIレスポンスで任意URLを開かない設計になっている。

## ドキュメント

- [ ] READMEの設定項目とUI文言が一致している。
- [ ] SPEC/SETTINGS_SPECの用語が現行UIと一致している。
- [ ] 不具合報告用データの説明が、zipと `.zip.gpg` の違いを明確にしている。
- [ ] Webhook URLが秘密情報であることを明記している。

## 既知の制限

- Windows実機ACLは別途確認が必要。
- Discord添付URLの長期永続性は保証しない。
- 自動更新は行わず、ユーザー操作で更新する。
- 診断データは秘密情報を含み得る。

## 未修正の既知問題

- `GO-2026-4550` が検出されている場合はリリース不可。
- 診断データ確認用zipに秘密情報が残る場合は、少なくともRelease notesやUIで強く注意喚起する。可能ならRelease前に修正する。

## リリース可否判断

- [ ] Critical/High findingsが0件。
- [ ] Medium findingsのうちRelease blocking項目が解消済み。
- [ ] CIとRelease workflowが成功。
- [ ] GitHub Release本文と成果物を人間が確認済み。

2026-06-25時点の暫定判断: `GO-2026-4550` が未解消なら正式リリース不可。
