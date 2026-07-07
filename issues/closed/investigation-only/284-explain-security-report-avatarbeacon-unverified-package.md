# 外部AvatarBeacon package同梱セキュリティレポートを解説する

## 指示

> '/home/user/work/ClipForVRChat/reports/security/2026-07-07T17-53-16.203Z/fa08ed3c46e4819185948dbcba8e4fa0'
> を日本語で解説

## 文脈

指定されたファイルは、Codex Securityが生成した `Unverified external Unity package bundled in releases` というhigh findingのレポートである。
Release workflowが外部AvatarBeacon `.unitypackage` をダウンロードして配布zipへ同梱する処理について、サプライチェーン上の信頼境界が指摘されている。

## 解釈

レポート本文を読み、指摘内容、攻撃シナリオ、根拠、提案修正を日本語で分かりやすく説明する。
この依頼ではコード修正までは行わず、現時点の理解と対応方針を整理する。

## 問題

- 外部Release assetをRelease workflow内で取得し、期待SHA256や署名を検証せずに公式配布zipへ入れている可能性がある。
- exe署名はzip内のUnity packageまでは保護しないため、利用者がexe署名を確認しても同梱packageの真正性までは保証されない。

## 期待する挙動

- レポートの意味、影響、根拠、提案修正が日本語で説明される。
- 今後修正する場合の要点が分かる。

## 受け入れ条件

- [x] findingの要旨を説明する。
- [x] なぜhigh扱いかを説明する。
- [x] レポート内の検証根拠を説明する。
- [x] 提案patchの意味を説明する。

## 調査結果

- findingは、Release workflowが外部の `AvatarBeacon_v0.0.1.unitypackage` をダウンロードし、固定済みの期待SHA256・署名・attestationと照合せずに配布zipへ同梱している点を指摘している。
- 現状の `.github/workflows/release.yml` では、ダウンロード後にSHA256を算出してbuild metadataへ記録しているが、取得前に信頼した期待値とは比較していない。
- ダウンロードしたpackageは通常利用者向けzipとseparated zipの両方へコピーされる。
- GPG署名は単体exeに対してのみ作成・検証され、zip全体やzip内のUnity packageは署名対象ではない。
- レポートの提案patchは、workflow環境変数へ固定SHA256を追加し、ダウンロード後のhashがその値と一致しなければRelease buildを失敗させる内容である。
