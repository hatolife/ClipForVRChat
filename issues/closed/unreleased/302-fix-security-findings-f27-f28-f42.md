# F27/F28/F42 の即時修正

## 指示

> Task: fix immediate no-spec findings F27/F28 and F42 from reports/security/2026-07-07T17-53-16.203Z/codex-security-findings-remediation-report.md.
> - Output format and JPEG quality controls should be editable whenever an output path uses encoded image data, not only when local save is enabled. Discord-only and clipboard-only workflows still use backend encoding.
> - Drop handling must cancel accepted drop events with preventDefault/stopPropagation so WebView cannot navigate to dropped URL/file. Keep existing Wails file-drop behavior intact.
> - Run `node scripts/check-frontend-template-literals.mjs` and `node scripts/check-wails-api-surface.mjs` if main.js changes. Edit files directly and report changed files/tests.

## 文脈

`reports/security/2026-07-07T17-53-16.203Z/codex-security-findings-remediation-report.md` で即時修正対象に分類された F27/F28/F42 を、フロントエンドの設定表示と drop イベント処理で解消する。

## 解釈

出力形式と JPEG 品質は、ローカル保存の有無ではなく「ローカル保存またはDiscord投稿でエンコード済み画像データを出力するかどうか」で編集可否を決める。
drop は Wails の file-drop を残しつつ、WebView 側の予期しないナビゲーションを止める。

## 問題

- Discord-only の処理でも backend は画像をエンコードするのに、UI では出力形式と JPEG 品質が編集できない。
- 受理済みの drop イベントを cancel しておらず、WebView がドロップされた URL / file へ遷移し得る。

## 期待する挙動

- 出力形式と JPEG 品質は、ローカル保存またはDiscord投稿のエンコード画像出力が関係する場面で常に編集できる。
- drop で受理したイベントは `preventDefault` / `stopPropagation` され、WebView のナビゲーションが起きない。
- Wails の file-drop による既存のファイル受け渡し挙動は壊さない。

## 受け入れ条件

- [x] Discord-only のワークフローでも出力形式と JPEG 品質を編集できる。
- [x] 受理した drop イベントで WebView のナビゲーションが抑止される。
- [x] Wails file-drop の既存動作を維持する。
- [x] `node scripts/check-frontend-template-literals.mjs` と `node scripts/check-wails-api-surface.mjs` を実行して確認する。
