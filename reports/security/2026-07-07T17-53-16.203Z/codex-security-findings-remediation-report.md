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

### 重大な仕様変更・運用判断が必要

| ID | Title | 理由 | 今回の扱い |
| --- | --- | --- | --- |
| F02 | Auto-capture Discord opt-out bypass leaks captures | `autoCapture.discord.enabled` と全体 `output.uploadDiscord` の同意境界をどう定義するかの仕様判断が必要。現HEADには全体ONで自動撮影Discord投稿もONになることを期待するテストがある。 | 保留 |
| F07 | Raw issue quoting can leak user secrets | AGENTS運用として「ユーザー原文を誤字含め保存」と「秘密情報redaction」のどちらを優先するかの運用判断が必要。 | 保留 |
| F11 | Hidden camera auto-start settings remain active | 既存configのカメラ自動起動/終了設定を強制無効化するか、UIへ戻すかの仕様判断が必要。 | 保留 |
| F30 | Tabbed settings hide webhook auto-post settings | 外部config保存時にセキュリティ関連タブの確認を必須化するか、現在のタブ+確認ダイアログ方式でよいかのUX/仕様判断が必要。 | 保留 |
| F37 | Release signing key exposed to prior build-step code execution | 署名job分離、鍵管理方式、hardware/KMS/keyless signing等のRelease信頼モデル判断が必要。 | 保留 |
| F43 | Auto-photo scan cap can leak old photos and starve new ones | scan cap時に古い写真をどう扱うか、baselineを全件保持するか、監視対象を制限するかの仕様判断が必要。 | 保留 |
| F44 | History file stores Discord tokens with weak permissions | 再起動後のDiscord削除機能維持とtoken非保存をどう両立するかのデータモデル判断が必要。権限自体は現HEADでprivate file化済み。 | 保留 |

### 仕様変更不要だが追加情報・方針確認が必要

| ID | Title | 必要な追加情報・方針 | 今回の扱い |
| --- | --- | --- | --- |
| F05 | New diagnostic links are blocked by URL allow-list | `www.google.com` や `keys.openpgp.org` をアプリの外部URL許可リストへ入れるか、リンク先を許可済みhostへ変えるか判断が必要。 | 保留 |
| F12 | New VRChat feedback help link is blocked | `feedback.vrchat.com` を許可リストへ入れるか、アプリ外導線を別形式にするか判断が必要。 | 保留 |

### 仕様変更不要かつ追加情報不要: 即時修正対象

#### 設定・自動投稿・Webhook routing

- F01 `Remote-player mask opt-out is silently reset`: explicit `false` をNormalizeで既定trueへ戻さない。
- F08 `Unsaved draft can activate auto-post watchers`: 画像引数処理でresults modeへ入る場合、runtime watcher configを保存済みconfigへ戻す。
- F17 `Auto-post warning suppressed when webhook is configured`: 自動写真/スクショ側は現HEADで大部分軽減済み。自動撮影側はF02の仕様判断に従属するため、今回の修正対象は残存する単純警告漏れが確認できた場合のみ。
- F33/F34 `Screenshots ... wrong Discord webhook`: screenshot watcherの空WebhookがAutoPhoto webhookへfallbackしないようにする。

#### 診断ログ・秘密情報redaction

- F13 `Diagnostic logs leak VRChat private instance IDs`: instance IDをログ/診断用にredactまたはhashする。
- F24 `Startup diagnostics log unredacted local paths`: 起動診断に残るpathをredaction済み形式へ寄せる。
- F38 `Diagnostic log can persist Discord webhook tokens`: URL/error redactionのテストを追加し、漏れがあればredactionを強化する。

#### 入力サイズ・DoS上限

- F03 `Auto-photo skip path bypasses per-tick processing cap`: skipもtick上限に数える。
- F04/F21/F22 `Unbounded CLI/ZIP encryption can exhaust memory`: CLI暗号化入力に通常ファイル確認とサイズ/件数上限を入れる。
- F10 `Spout diagnose trusts unbounded sender dimensions`: Spout sender dimensionsに上限とoverflow checkを入れる。
- F15 `Unbounded localhost IPC request can exhaust app memory`: IPC decode前にrequest size、command/token/path件数/長さ上限を入れる。
- F18 `Unbounded OSC avatar parameter cache enables DoS`: avatar OSC sample cacheにallowlistまたは上限/evictionを入れる。
- F31/F32 `Auto-photo scan errors ...`: 同一scan errorの重複追加を抑制し、UI/result増殖を防ぐ。
- F39 `QR upscaling can exhaust memory on crafted tall images`: QR upscale後の総pixel数上限を入れる。
- F41 `Unbounded native clipboard PNG copy can crash app`: clipboard PNGをcopyする前にGlobalSize上限を確認する。

#### 境界チェック・URL・WebView

- F14 `OSC forwarding can self-loop and cause UDP DoS`: wildcard bind時のsame-port local targetを保守的にself-forward扱いにする。
- F23 `COM shell calls are made without OS-thread pinning`: `reveal_windows.go` のCOM呼び出しを `runtime.LockOSThread` で囲む。
- F25/F26 `Explorer reveal ...`: Wails APIでも管理output配下のファイルだけExplorer表示できるようにする。
- F27/F28 `Output format ...`: Discord/clipboard-only出力でも実際に使われるformat/qualityをUIで変更できるようにする。
- F35/F36 `Update URL ...`: API応答の `html_url` をそのまま使わず、公式GitHub Release URLを構築または厳格検証する。
- F42 `Accepted drops are not cancelled...`: drop eventで `preventDefault` / `stopPropagation` し、意図しないWebView navigationを防ぐ。

### HEAD確認済み・重複・軽減済み

| ID | Title | 確認内容 | 今回の扱い |
| --- | --- | --- | --- |
| F09 | Debug OSC input is written verbatim to diagnostics | 現HEADの `SendDebugOSC` はraw lineを診断ログへ書かず、target/address/typesと結果だけを記録している。 | 対応済み扱い |
| F16 | Unbounded OSC packet logging enables log DoS | 現HEADではOSCログはメモリ上 `maxOSCLogEntries=1000` に制限され、永続化も専用ログへthrottleされている。追加強化はF18/F14側で扱う。 | 軽減済み |
| F19 | Metadata failures can lead to unvalidated Discord uploads | 後続実装確認が必要だが、画像decode/metadataの扱いとDiscord upload経路が複数に分かれるため、F38/F33/F34/F41完了後に再確認する。 | 後続再確認 |
| F20 | Release notes step uses wrong env syntax on Windows | 現HEADのRelease workflowはPowerShellで `$env:...`、bashで環境変数を使い分けている。 | 対応済み扱い |
| F29 | Hidden auto-processing folders can exfiltrate images | `issues/293` の対応で監視フォルダUI表示と保存前確認が追加済み。 | 対応済み扱い |
| F40 | Untrusted QR URLs can trigger Discord mentions | Discord payloadに `allowed_mentions.parse=[]` が入り、テストも存在する。 | 対応済み扱い |

## 作業割り当て

| Work | 対象finding | 担当 | 状態 | 主なファイル | メモ |
| --- | --- | --- | --- | --- | --- |

## 作業結果

- 未着手

## テスト結果

- 未実行

## 残課題

- 未分類
