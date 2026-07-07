# F15/F23/F25/F26/F35/F36 の即時修正を行う

## 指示

> Task: fix immediate no-spec findings F15, F23, F25/F26, F35/F36 from reports/security/2026-07-07T17-53-16.203Z/codex-security-findings-remediation-report.md.
> - Bound localhost single-instance IPC request size and token/command/path lengths/counts before/while decoding.
> - Lock the OS thread around Windows COM Shell calls in reveal_windows.go.
> - Ensure RevealFileInExplorer only reveals app-managed output files under the configured output directory; reject arbitrary renderer paths including UNC/device-like paths.
> - Ensure update info does not expose/open arbitrary html_url; construct or validate official GitHub release URLs.
> Add focused tests. Run relevant go tests from src if feasible. Edit files directly and report changed files/tests.

## 文脈

`297` の分類で即時修正対象となった F15, F23, F25, F26, F35, F36 をまとめて対処する。
単一起動IPC、Windows Explorer表示、GitHub Release更新URLの3経路をそれぞれ境界強化する。

## 解釈

このissueでは、既存の正常系を保ちながら、外部入力のサイズ・形式・配布元を厳格化する。
実装後は、境界値と拒否系のテストを追加して、再発しやすい入力経路を固定する。

## 問題

- localhost IPC 受信で巨大JSONや長い配列/文字列を与えると、不要にメモリを使う可能性がある。
- Windows COM Shell 呼び出しが OS スレッド依存なのに、固定されていない。
- Explorer 表示 API が、アプリ管理外の任意パスを受ける可能性がある。
- 更新通知が API 応答の `html_url` をそのまま開くと、信頼境界を越える。

## 期待する挙動

- single-instance IPC は、受信サイズと各フィールドの長さ/件数に上限がある。
- reveal は configured output directory 配下のアプリ管理ファイルのみを対象にする。
- update URL は GitHub Release の公式URLとして構築または検証される。
- それぞれに対応した focused tests がある。

## 受け入れ条件

- [x] `src/single_instance.go` で IPC request size と token/command/path の上限が入る。
- [x] `src/reveal_windows.go` の COM Shell 呼び出しが OS thread に固定される。
- [x] `src/app.go` / `src/reveal_other.go` で Explorer reveal が managed output files に限定される。
- [x] `src/internal/appcore/update.go` で update URL が公式 GitHub Release URL に限定される。
- [x] 関連する tests が追加され、対象 `go test` が通る。
