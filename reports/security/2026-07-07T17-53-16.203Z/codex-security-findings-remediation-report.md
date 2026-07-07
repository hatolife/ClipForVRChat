# Codex Security findings 全件対処レポート

## 概要

- 底本CSV: `reports/security/2026-07-07T17-53-16.203Z/codex-security-findings-2026-07-07T19-00-29.151Z.csv`
- 底本CSV timestamp: `2026-07-07T19-00-29.151Z`
- 作成日時: `2026-07-07T19:16:35Z`
- finding件数: 44件
- severity内訳: `high` 3件, `medium` 25件, `low` 10件, `informational` 6件
- 目的: 全件を最終的に対処するため、重大仕様変更が必要なもの、追加情報が必要なもの、即時修正対象を分け、作業結果を追跡する。

## 運用方針

- 重大な仕様変更が必要なfindingは、ユーザー承認なしで仕様縮小・延期・代替仕様化しない。
- 追加情報が必要なfindingは、必要情報を明記して今回の即時修正対象から外す。
- 仕様変更不要かつ追加情報不要のfindingは、内容別に分類してサブエージェントへ修正を委任する。
- 作業完了後、対応内容、テスト、コミット、残課題を本レポートへ追記する。

## 初期finding台帳

| ID | Severity | Title | Commit | Relevant paths | 初期状態 |
| --- | --- | --- | --- | --- | --- |
| F01 | `low` | Remote-player mask opt-out is silently reset | `057843a2e7e8` | src/internal/appcore/config.go <br> src/frontend/src/main.js <br> src/internal/appcore/autocapture.go | 未分類 |
| F02 | `medium` | Auto-capture Discord opt-out bypass leaks captures | `823aeee36395` | src/internal/appcore/autocapture.go <br> src/frontend/src/main.js | 未分類 |
| F03 | `low` | Auto-photo skip path bypasses per-tick processing cap | `d94f511fbdea` | src/internal/appcore/autophoto.go <br> src/app.go | 未分類 |
| F04 | `medium` | Unbounded CLI encryption can exhaust memory | `c26a344a629e` | src/main.go <br> src/diagnostic_package.go | 未分類 |
| F05 | `informational` | New diagnostic links are blocked by URL allow-list | `84591e3fe3e9` | src/frontend/src/main.js <br> src/app.go | 未分類 |
| F06 | `low` | Unrelated OSC parameters can suppress fallback | `004bf2362a0e` | src/app.go | 未分類 |
| F07 | `medium` | Raw issue quoting can leak user secrets | `46386b682cb4` | AGENTS.md <br> issues/closed/unreleased/236-sort-closed-issue-index-and-preserve-user-instructions.md | 未分類 |
| F08 | `medium` | Unsaved draft can activate auto-post watchers | `2010d3e6f56d` | src/frontend/src/main.js <br> src/main.go <br> src/app.go | 未分類 |
| F09 | `low` | Debug OSC input is written verbatim to diagnostics | `05c675b8b89a` | src/app.go <br> src/internal/appcore/diagnostic.go <br> src/diagnostic_package.go <br> src/frontend/src/main.js | 未分類 |
| F10 | `low` | Spout diagnose trusts unbounded sender dimensions | `5d8bef2f993e` | tools/spout-capture/main.cpp | 未分類 |
| F11 | `medium` | Hidden camera auto-start settings remain active | `215f5dea51eb` | src/frontend/src/main.js <br> src/internal/appcore/autocapture.go | 未分類 |
| F12 | `informational` | New VRChat feedback help link is blocked | `e6825c2dd549` | src/frontend/src/main.js <br> src/app.go | 未分類 |
| F13 | `medium` | Diagnostic logs leak VRChat private instance IDs | `f5635ac75621` | src/app.go <br> src/internal/appcore/autocapture.go <br> src/internal/appcore/diagnostic.go <br> src/diagnostic_package.go | 未分類 |
| F14 | `medium` | OSC forwarding can self-loop and cause UDP DoS | `da42af372fd6` | src/app.go | 未分類 |
| F15 | `low` | Unbounded localhost IPC request can exhaust app memory | `f62c38923cee` | src/single_instance.go | 未分類 |
| F16 | `medium` | Unbounded OSC packet logging enables log DoS | `3a593468ddac` | src/app.go <br> src/internal/appcore/diagnostic.go <br> src/diagnostic_package.go | 未分類 |
| F17 | `medium` | Auto-post warning suppressed when webhook is configured | `6f250c5054cd` | src/frontend/src/main.js <br> src/app.go <br> src/internal/appcore/autophoto.go <br> src/internal/appcore/autocapture.go | 未分類 |
| F18 | `medium` | Unbounded OSC avatar parameter cache enables DoS | `6b711fd33e45` | src/app.go | 未分類 |
| F19 | `medium` | Metadata failures can lead to unvalidated Discord uploads | `88bab9bd0623` | src/internal/appcore/metadata.go <br> src/internal/appcore/autocapture.go <br> src/internal/appcore/autophoto.go | 未分類 |
| F20 | `informational` | Release notes step uses wrong env syntax on Windows | `36e242ccaacf` | github/workflows/release.yml <br> scripts/extract-release-notes.mjs | 未分類 |
| F21 | `medium` | Unbounded ZIP CLI encryption can exhaust memory | `62d856115af0` | src/main.go <br> src/diagnostic_package.go | 未分類 |
| F22 | `low` | Unbounded ZIP encryption can exhaust memory | `410470383078` | src/main.go <br> src/diagnostic_package.go | 未分類 |
| F23 | `informational` | COM shell calls are made without OS-thread pinning | `c3c9ac1cc82a` | src/reveal_windows.go | 未分類 |
| F24 | `medium` | Startup diagnostics log unredacted local paths | `ec42e072d897` | src/app_diagnostic.go | 未分類 |
| F25 | `medium` | Untrusted history paths are passed to Explorer | `8057aba91e88` | src/app.go <br> src/frontend/src/main.js <br> src/internal/appcore/history.go | 未分類 |
| F26 | `medium` | Explorer reveal trusts history-controlled file paths | `42503e6d5d55` | src/main.go <br> src/app.go <br> src/frontend/src/main.js <br> src/internal/appcore/history.go | 未分類 |
| F27 | `informational` | Output format disabled for Discord-only processing | `5491315a78df` | src/frontend/src/main.js <br> src/internal/appcore/processor.go | 未分類 |
| F28 | `informational` | Output format controls disabled for Discord-only output | `b0a75b48ac84` | src/frontend/src/main.js <br> src/internal/appcore/processor.go | 未分類 |
| F29 | `medium` | Hidden auto-processing folders can exfiltrate images | `11cdfbff7b2b` | src/frontend/src/main.js <br> src/app.go <br> src/internal/appcore/config.go <br> src/internal/appcore/autophoto.go | 未分類 |
| F30 | `high` | Tabbed settings hide webhook auto-post settings | `c10d11bef19f` | src/frontend/src/main.js <br> src/app.go <br> src/internal/appcore/autophoto.go | 未分類 |
| F31 | `low` | Auto-photo scan errors are appended without deduplication | `31df92b2dbea` | src/internal/appcore/autophoto.go <br> src/app.go <br> src/frontend/src/main.js | 未分類 |
| F32 | `low` | Auto-photo scan errors cause unbounded UI result growth | `45068ab64063` | src/internal/appcore/autophoto.go <br> src/app.go <br> src/frontend/src/main.js | 未分類 |
| F33 | `medium` | Screenshots are sent to the auto-photo Discord webhook | `a4edd98a33c9` | src/app.go <br> src/internal/appcore/autophoto.go | 未分類 |
| F34 | `medium` | Screenshots can upload to the wrong Discord webhook | `f759c2f5f2f8` | src/app.go <br> src/internal/appcore/autophoto.go <br> src/internal/appcore/processor.go | 未分類 |
| F35 | `medium` | Update banner opens unvalidated release URL | `866ae90d3892` | src/internal/appcore/update.go <br> src/frontend/src/main.js <br> src/app.go | 未分類 |
| F36 | `medium` | Unvalidated update URL can open attacker-controlled links | `40d856a59b89` | src/internal/appcore/update.go <br> src/frontend/src/main.js <br> src/app.go | 未分類 |
| F37 | `high` | Release signing key exposed to prior build-step code execution | `96dedeadff94` | github/workflows/release.yml | 未分類 |
| F38 | `high` | Diagnostic log can persist Discord webhook tokens | `df8bd7f4a27b` | src/internal/appcore/diagnostic.go <br> src/internal/appcore/processor.go <br> src/internal/appcore/discord.go | 未分類 |
| F39 | `medium` | QR upscaling can exhaust memory on crafted tall images | `fc88b5d3139a` | src/internal/appcore/qrcode.go <br> src/internal/appcore/processor.go <br> src/internal/appcore/image.go | 未分類 |
| F40 | `medium` | Untrusted QR URLs can trigger Discord mentions | `39531cf3a7e5` | src/internal/appcore/qrcode.go <br> src/internal/appcore/processor.go <br> src/internal/appcore/discord.go | 未分類 |
| F41 | `medium` | Unbounded native clipboard PNG copy can crash app | `21cc973660d7` | src/internal/appcore/clipboard.go <br> src/internal/appcore/clipboard_native_windows.go <br> src/internal/appcore/processor.go <br> src/internal/appcore/image.go | 未分類 |
| F42 | `low` | Accepted drops are not cancelled, allowing WebView navigation | `ec99a084bfbf` | src/frontend/src/main.js <br> src/main.go | 未分類 |
| F43 | `medium` | Auto-photo scan cap can leak old photos and starve new ones | `d9a7f58be756` | src/internal/appcore/autophoto.go <br> src/internal/appcore/processor.go | 未分類 |
| F44 | `medium` | History file stores Discord tokens with weak permissions | `709370d231b9` | src/internal/appcore/history.go <br> src/internal/appcore/processor.go <br> src/internal/appcore/types.go | 未分類 |

## 分類結果

### 重大な仕様変更が必要

- 未分類

### 仕様変更不要だが追加情報が必要

- 未分類

### 仕様変更不要かつ追加情報不要

- 未分類

## 作業割り当て

| Work | 対象finding | 担当 | 状態 | 主なファイル | メモ |
| --- | --- | --- | --- | --- | --- |

## 作業結果

- 未着手

## テスト結果

- 未実行

## 残課題

- 未分類
