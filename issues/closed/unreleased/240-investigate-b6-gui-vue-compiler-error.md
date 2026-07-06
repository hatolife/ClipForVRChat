# b6でGUIが表示されないVue compiler errorを調査する

## 指示

> '/mnt/c/Users/user/Downloads/ClipForVRChat-v0.1.8-b6-windows-amd64/logs'
>
> を参照
> b6でUIが表示されない致命的な不具合発生している
> 過去に同様の不具合が発生したことがある
> 再発防止済みのはず
> まず調査 修正作業は待って

## 文脈

`v0.1.8-a18` から `v0.1.8-a21` の対応で、GUI未表示時のWails/frontend診断ログ、起動進捗表示、Vue template literal内の生バッククォート検出を追加済みだった。`v0.1.8-b6` で再びUIが表示されないため、指定ログを起点に既存の再発防止で検出できなかった原因を切り分ける。

## 解釈

現時点では修正作業に入らず、ログと既存の再発防止策の範囲を調査する。原因候補、既存検査で漏れた理由、修正時に必要な受け入れ条件を整理してから次の作業可否を確認する。

## 問題

`v0.1.8-b6` 起動時にfrontend assetは読み込まれているが、Vue mount中に `SyntaxError: https://vuejs.org/error-reference/#compiler-24` が発生し、UI初期化が完了していない。`api GetInitialState begin` へ到達していないため、Go API初期状態取得前のVue template compile段階で停止している可能性が高い。

## 期待する挙動

- `v0.1.8-b6` 相当のfrontend templateでもGUIが表示される。
- Vue template構文エラーがrelease前のCI/Release検査で検出される。
- 既存のtemplate literal生バッククォート検査だけで拾えないVue template破損も、再発防止対象として扱える。

## 受け入れ条件

- 指定ログから、Wails、frontend script、Vue mount、Go API初期化のどこで止まったかを説明できる。
- `compiler-24` の実際の意味と該当するtemplate破損箇所を特定できる。
- 既存の `scripts/check-frontend-template-literals.mjs` とCI/Release検査で漏れた理由を説明できる。
- 修正作業へ入る前に、必要な修正範囲と検証コマンドを提示する。
- ユーザー承認前にアプリ本体の修正を行わない。

## 調査メモ

- ログ: `/mnt/c/Users/user/Downloads/ClipForVRChat-v0.1.8-b6-windows-amd64/logs/2026-07-07.log`
- `ui wails run begin`、`ui lifecycle startup begin/complete`、`ui action="frontend_script_loaded" detail="... api=available"` までは到達している。
- `ui action="frontend_mount_error"` に `SyntaxError: https://vuejs.org/error-reference/#compiler-24` が記録されている。
- 直後に `ui action="frontend_error"` が `Uncaught SyntaxError: https://vuejs.org/error-reference/#compiler-24 at http://wails.localhost/:0:0` を記録している。
- `api GetInitialState begin` は出ていないため、初期状態API呼び出し前にVue mountが失敗している。
- ログ末尾には `ui lifecycle before_close` と `shutdown complete` があり、アプリプロセス自体は終了処理まで進んでいる。
- `auto-capture pose received` が短時間に大量出力されており、UI停止の直接原因とは別にログ量の多さも確認が必要。
- ログは2761行で、そのうち `auto-capture pose received` が2738行ある。約28秒間の起動で大量に出ている。
- ローカルで `@vue/compiler-dom` に `src/frontend/src/main.js` の `template` 文字列を直接渡すと、次の3件を検出した。
  - [src/frontend/src/main.js:2244](/home/user/work/ClipForVRChat/src/frontend/src/main.js:2244) `code=24 Element is missing end tag.` 自動撮影タブの `<section>`。
  - [src/frontend/src/main.js:2900](/home/user/work/ClipForVRChat/src/frontend/src/main.js:2900) `code=23 Invalid end tag.`。
  - [src/frontend/src/main.js:2901](/home/user/work/ClipForVRChat/src/frontend/src/main.js:2901) `code=23 Invalid end tag.`。
- 直接原因の候補は [src/frontend/src/main.js:2573](/home/user/work/ClipForVRChat/src/frontend/src/main.js:2573) の `</div>` が、[src/frontend/src/main.js:2574](/home/user/work/ClipForVRChat/src/frontend/src/main.js:2574) の `<template v-else-if="autoCaptureDetailView === 'fallback'">` より前にあり、詳細画面内の `v-if / v-else-if / v-else` 連鎖を壊していること。
- `git blame` では、該当するfallback分岐追加と周辺の `</div>` は `474c600 fix(autocapture): split fallback auto controls` で入っている。
- `node scripts/check-frontend-template-literals.mjs` は成功した。生バッククォート混入だけを見る検査のため、Vue templateのタグ整合性は検出対象外。
- `node scripts/check-wails-api-surface.mjs` は成功した。今回の原因はWails API不一致ではない。
- `src/frontend` で `npm run build` は成功した。`createApp({ template: \`...\` })` のランタイムtemplateはVite build時にVue templateとして検査されず、実行時のVue mountで初めてcompiler errorになる。
- 同じ `@vue/compiler-dom` 確認では、`v0.1.8-b4` と `v0.1.8-b5` はOK、`v0.1.8-b6` と `474c600` は同じ3件のエラー。`474c600^` はOKだったため、混入コミットは `474c600 fix(autocapture): split fallback auto controls` と見てよい。

## 対応メモ

- `src/frontend/src/main.js` の自動撮影詳細画面で、`fallback` 分岐が `settings-detail-view` の外へ出ていたため、余分な `</div>` の位置を修正して `v-if / v-else-if / v-else` 連鎖を復元した。
- `scripts/check-vue-runtime-template.mjs` を追加し、`@vue/compiler-dom` で `main.js` のruntime templateを直接compileしてタグ不整合を検出するようにした。
- CIとRelease workflowへ `Check Vue runtime template` を追加した。
- `docs/frontend-runtime-troubleshooting.md` へ、`compiler-24` / runtime templateタグ不整合の原因、同定方法、対処方法を追記した。

## 検証メモ

- `node scripts/check-frontend-template-literals.mjs`
- `node scripts/check-vue-runtime-template.mjs`
- `node scripts/check-wails-api-surface.mjs`
- `npm run build`
- `cd src && GOCACHE=/tmp/clipforvrchat-go-cache go test ./...`
