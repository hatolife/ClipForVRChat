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
| F44 | History file stores Discord tokens with weak permissions | 再起動後のDiscord削除機能維持とtoken非保存をどう両立するかのデータモデル判断が必要。権限自体は現HEADでprivate file化済み。 | 保留 |

### 重大な仕様変更・運用判断待ちから対応済みへ移した項目

- F43 `Auto-photo scan cap can leak old photos and starve new ones`: 追加情報 `e2dbfdd09f68819187b842db37a4421f` を確認した結果、19:48-F30と同じ根本原因だった。`53b74d5 fix(autophoto): baseline full scans before tick limits` で、起動時baseline/current判定を「scan capに入った先頭N件だけを見る」処理から「対象画像を全件scanして既存/新規を判定し、処理件数だけper-tick上限で制限する」処理へ変更済みのため、旧画像が後から新規扱いで投稿される経路と、先頭fillerにより後方新規画像が永久に見えない経路は解消済みとして扱う。

### 仕様変更不要だが追加情報・方針確認が必要

| ID | Title | 必要な追加情報・方針 | 今回の扱い |
| --- | --- | --- | --- |
| F05 | New diagnostic links are blocked by URL allow-list | 19:48再分析で追加情報を確認し、検索エンジンhostは許可せず、GnuPG公式/keys.openpgp.orgだけを許可する方針に決定した。 | 19:48 W10で対応 |
| F12 | New VRChat feedback help link is blocked | 19:48再分析で追加情報を確認し、固定ヘルプリンク先として `feedback.vrchat.com` を許可する方針に決定した。 | 19:48 W10で対応 |

### 仕様変更不要かつ追加情報不要: 即時修正対象

#### 設定・自動投稿・Webhook routing

- F01 `Remote-player mask opt-out is silently reset`: explicit `false` をNormalizeで既定trueへ戻さない。
- F08 `Unsaved draft can activate auto-post watchers`: 画像引数処理でresults modeへ入る場合、runtime watcher configを保存済みconfigへ戻す。
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
- F27/F28 `Output format ...`: Discord-only出力でも実際に使われるformat/qualityをUIで変更できるようにする。
- F35/F36 `Update URL ...`: API応答の `html_url` をそのまま使わず、公式GitHub Release URLを構築または厳格検証する。
- F42 `Accepted drops are not cancelled...`: drop eventで `preventDefault` / `stopPropagation` し、意図しないWebView navigationを防ぐ。

### HEAD確認済み・重複・軽減済み

| ID | Title | 確認内容 | 今回の扱い |
| --- | --- | --- | --- |
| F09 | Debug OSC input is written verbatim to diagnostics | 現HEADの `SendDebugOSC` はraw lineを診断ログへ書かず、target/address/typesと結果だけを記録している。 | 対応済み扱い |
| F16 | Unbounded OSC packet logging enables log DoS | 現HEADではOSCログはメモリ上 `maxOSCLogEntries=1000` に制限され、永続化も専用ログへthrottleされている。追加強化はF18/F14側で扱う。 | 軽減済み |
| F17 | Auto-post warning suppressed when webhook is configured | 現HEADの `autoPostConfirmationItems` は、VRChat写真自動処理とスクリーンショット自動処理について、Webhook設定済みでも `output.uploadDiscord` かつ自動処理ONなら保存前確認へ出す。自動撮影側の同意境界はF02の仕様判断に従属する。 | 写真/スクショは対応済み扱い、自動撮影はF02従属 |
| F19 | Metadata failures can lead to unvalidated Discord uploads | 後続実装確認が必要だが、画像decode/metadataの扱いとDiscord upload経路が複数に分かれるため、F38/F33/F34/F41完了後に再確認する。 | 後続再確認 |
| F20 | Release notes step uses wrong env syntax on Windows | 現HEADのRelease workflowはPowerShellで `$env:...`、bashで環境変数を使い分けている。 | 対応済み扱い |
| F29 | Hidden auto-processing folders can exfiltrate images | `issues/293` の対応で監視フォルダUI表示と保存前確認が追加済み。 | 対応済み扱い |
| F40 | Untrusted QR URLs can trigger Discord mentions | Discord payloadに `allowed_mentions.parse=[]` が入り、テストも存在する。 | 対応済み扱い |

## 作業割り当て

| Work | 対象finding | 担当 | 状態 | 主なファイル | メモ |
| --- | --- | --- | --- | --- | --- |
| W1 | F01, F03, F31, F32, F33, F34 | Raman | 完了 | `src/internal/appcore/config.go`, `src/internal/appcore/autophoto.go`, tests | 自動投稿・RemotePlayer・scan status |
| W2 | F39, F41 | Harvey | 完了 | `src/internal/appcore/qrcode.go`, `src/internal/appcore/clipboard*.go`, tests | QR upscale/clipboard PNG上限 |
| W3 | F04, F21, F22 | Heisenberg | 完了 | `src/diagnostic_package.go`, `src/diagnostic_package_test.go` | CLI/ZIP暗号化の入力上限 |
| W4 | F15, F23, F25, F26, F35, F36 | Newton | 完了 | `src/single_instance.go`, `src/reveal_windows.go`, `src/app.go`, `src/internal/appcore/update.go`, tests | IPC/Explorer/update境界 |
| W5 | F27, F28, F42 | Sagan | 完了 | `src/frontend/src/main.js` | 出力形式UI/drop cancellation |
| W6 | F10 | Hooke | 完了 | `tools/spout-capture/main.cpp` | Spout sender dimension上限 |
| W7 | F06, F08, F13, F14, F18, F24, F38 | メインエージェント | 完了 | `src/main.go`, `src/app.go`, `src/app_diagnostic.go`, tests | draft/runtime境界、OSC/cache、診断ログredaction |

## 作業結果

### W1: 自動投稿・RemotePlayer・scan status

- F01のために、`AutoCaptureConfig.Normalize` の `RemotePlayer` 補完条件を「nil のときだけ既定値を入れる」処理へ変更し、ユーザーが明示した `RemotePlayer=false` を保存後も維持するようにした。
- F03のために、auto-photo watcher の tick 処理を「stable判定後にskipされたファイルも `processed` に数える」処理へ変更し、skip対象を大量投入して `MaxAutoPhotoProcessPerTick` を迂回できないようにした。
- F31/F32のために、auto-photo watcher に直前のscan statusを保持する処理を追加し、同じ missing/inaccessible/scan-limit error を連続tickで結果へ重複追加しないようにした。
- F33/F34のために、`AutoPhotoWatcher.webhookURL()` を「明示された `WebhookURL` だけをoverrideとして返す」処理へ変更し、screenshot auto-post の空Webhookが `AutoPhoto.WebhookURL` へ誤ってfallbackしないようにした。

### W2: QR upscale / clipboard PNG 上限

- F39のために、QR検出の拡大variant生成を「拡大後の推定総ピクセル数が `MaxImagePixels` 以下の場合だけ生成する」処理へ変更し、細長い画像で2x/3x/4x画像が巨大化しないようにした。
- F41のために、Windows native clipboard PNG読み出しを「`GlobalSize` を確認して既存の入力byte上限を超える場合は `unsafe.Slice` / `make` / `copy` の前に拒否する」処理へ変更し、巨大clipboard dataでアプリがメモリを使い切らないようにした。
- F41のために、native clipboard PNGで上限超過を検出した場合は汎用clipboard読み取りへフォールスルーせずエラーを返す処理へ変更し、巨大データの再試行経路を閉じた。

### W3: CLI/ZIP暗号化入力上限

- F04/F21/F22のために、CLI暗号化のroot入力検証を `os.Stat` から `os.Lstat` ベースへ変更し、FIFO、device、symlinkなどの非通常ファイルを暗号化前に拒否するようにした。
- F04/F21/F22のために、単一ファイル暗号化を「通常ファイルかつ上限サイズ以下の場合だけ `os.ReadFile` へ進む」処理へ変更し、巨大ファイルを全量メモリへ読み込まないようにした。
- F04/F21/F22のために、ディレクトリ暗号化を「zip化前に深さ、件数、単体サイズ、総サイズを検証する」処理へ変更し、巨大ツリーを `bytes.Buffer` へ展開しないようにした。

### W4: IPC / Explorer / update URL / COM境界

- F15のために、single-instance IPCの受信処理を「`io.LimitReader` 相当の上限付きdecodeとtoken/command/path件数/長さ検証を通す」処理へ変更し、localhostから巨大JSONを送られてもメモリを使い切らないようにした。
- F15のために、IPC送信側もpath件数/長さを送信前に検証する処理へ変更し、自プロセス経由でも上限外requestを作らないようにした。
- F23のために、Windows Shell APIによるExplorer表示処理を `runtime.LockOSThread` / `UnlockOSThread` で囲む処理へ変更し、COM初期化と解放が同じOS threadで行われるようにした。
- F25/F26のために、`RevealFileInExplorer` のpath解決を `ResolveHistoryOutputPath` から `ResolveManagedHistoryOutputPath` へ変更し、rendererから渡された任意pathではなく、設定された管理output配下の既存ファイルだけをExplorer表示できるようにした。
- F35/F36のために、update確認結果のURL生成を「API応答の `html_url` を採用する」処理から「取得したtag名で公式GitHub Release URLを構築する」処理へ変更し、改ざん応答で任意URLを開かないようにした。

### W5: frontend出力形式UI / drop cancellation

- F27/F28のために、出力形式とJPEG品質の編集可否を「ローカル保存ONのときだけ有効」から「ローカル保存またはDiscord投稿がONのときに有効」へ変更し、Discord-onlyでも実際に投稿へ使われるformat/qualityを設定できるようにした。
- F42のために、window drop handler を `preventDefault()` と `stopPropagation()` を呼ぶ処理へ変更し、Wails file-drop後にWebViewがドロップされたURL/fileへ遷移しないようにした。

### W6: Spout sender dimension上限

- F10のために、Spout senderの幅/高さを使う前に `validate_sender_dimensions` 相当の共通検証を通す処理へ変更し、capture/diagnoseの両方で異常寸法を拒否するようにした。
- F10のために、`width * height * 4` の計算をoverflow-safeなhelper経由へ変更し、巨大値や乗算overflowで小さいbufferを確保しないようにした。
- F10のために、異常寸法時の戻り値をthrow/未定義動作から `sender_dimension_error` のJSON diagnostic errorへ変更し、利用者と診断ログで原因を判別できるようにした。

### W7: draft/runtime境界・OSC・診断ログ

- F08のために、画像引数処理へ入る直前のUI stateを「未保存draftの `state.Config` のまま」から「保存済み `cfg` を `state.Config` に戻し、draft metadataを解除する」処理へ変更し、results modeで起動するwatcherが未保存draftのWebhook/監視設定を使わないようにした。
- F06のために、avatar OSC fallback判定を「任意の `/avatar/parameters` 受信があればfallbackを止める」処理から「設定されたbasis parameterを受信した場合だけfallbackを止める」処理へ変更し、無関係parameterでpreplaced fallbackが抑止されないようにした。
- F13のために、world安定待ちとavatar OSC session resetの診断ログを「raw `instance_id` を出す」処理から「`instance_id` を `<redacted-vrchat-instance-id>` に置き換える」処理へ変更し、private instance情報を診断ログへ残さないようにした。
- F14のために、OSC forwardのself判定を「wildcard bindではloopback/unspecifiedだけself扱い」から「wildcard bindでは同じportのローカルinterface宛もself扱い」へ変更し、`0.0.0.0:<port>` 受信時にLAN IP同portへ転送して自己ループしないようにした。
- F18のために、avatar OSC sample保存を「mapへ無制限に追加」から「最大 `maxAvatarOSCSamples` 件へ制限し、古い非basis sampleを優先削除する」処理へ変更し、任意parameter名の大量送信でメモリが増え続けないようにした。
- F24のために、起動診断summaryへ入れるoutput/auto-photo/screenshot/auto-capture directoryを「raw path文字列」から「`RedactDiagnosticText` 適用済み文字列」へ変更し、ユーザー名や個人フォルダ名が診断ログに残りにくいようにした。
- F38のために、Discord webhook URLを含むエラー文字列のredactionをテストで固定し、`url.Error` 形式の `Post "https://discord.com/api/webhooks/..."` が診断ログへ生tokenとして残らないことを確認できるようにした。

## テスト結果

- `cd src && go test ./...` を実行し、成功した。
- `node scripts/check-frontend-template-literals.mjs` を実行し、成功した。
- `node scripts/check-wails-api-surface.mjs` を実行し、成功した。
- `node scripts/check-closed-issue-index.mjs` を実行し、成功した。
- `cmake -S tools/spout-capture -B /tmp/cfvrc-spout-logic-test && cmake --build /tmp/cfvrc-spout-logic-test --target spout-capture-logic-test -j2 && ctest --test-dir /tmp/cfvrc-spout-logic-test --output-on-failure` をHookeが実行し、成功した。
- `x86_64-w64-mingw32-g++ -std=c++17 -Wall -Wextra -fsyntax-only -I/tmp/clipforvrchat-spout-mingw/_deps/spout2-src/SPOUTSDK/SpoutLibrary tools/spout-capture/main.cpp` をHookeが実行し、成功した。

## コミット

- `9161f8a fix(autoprocess): harden watcher config boundaries`: F01/F03/F08/F31/F32/F33/F34のために、config normalize、watcher tick上限、scan status dedup、Webhook解決、results mode draft解除を変更した。
- `c9f41d8 fix(inputs): bound diagnostic qr and clipboard reads`: F04/F21/F22/F39/F41のために、CLI暗号化入力検証、QR upscale上限、clipboard PNG byte上限を追加した。
- `a2ba76f fix(app): harden ipc explorer osc and diagnostics`: F06/F13/F14/F15/F18/F23/F24/F25/F26/F35/F36/F38のために、IPC検証、Explorer管理path制限、OSC self-loop/cache制限、update URL固定、診断ログredactionを変更した。
- `0e1c7f0 fix(ui): enable discord output format controls`: F27/F28/F42のために、Discord-only時のformat/quality編集条件とdrop event cancellationを変更した。
- `5af14b1 fix(spout): validate sender dimensions`: F10のために、Spout sender寸法検証とoverflow-safe buffer計算を追加した。
- `docs(issues): record security remediation work`: 本レポート、監督issue、closed issue indexを更新し、完了/保留/後続再確認の状態を記録する。

## 残課題

- 重大な仕様変更・運用判断待ち: F02, F07, F11, F30, F37, F44。今回の即時修正からは除外し、上記分類表の理由に沿って別途仕様判断する。F43は追加情報再検証によりW9修正で対応済みへ移した。
- 追加情報・方針確認待ちだったF05/F12は、19:48再分析で方針を決定しW10で対応した。
- 後続再確認: F19。F33/F34/F38/F41の修正後に、metadata失敗時のDiscord upload経路が未検証投稿にならないか再確認する。
- HEAD確認済み・重複・軽減済み: F09, F16, F17, F20, F29, F40。追加コード修正は不要として扱った。

## 19:48 CSV再分析

### 概要

- 底本CSV: `reports/security/2026-07-07T17-53-16.203Z/codex-security-findings-2026-07-07T19-48-25.660Z.csv`
- 底本CSV timestamp: `2026-07-07T19-48-25.660Z`
- finding件数: 50件
- severity内訳: `high` 5件, `medium` 29件, `low` 10件, `informational` 6件
- 前回19:00 CSVからの追加finding: 6件

### 追加finding

| 19:48 ID | Severity | Title | Commit | Relevant paths | 分類 |
| --- | --- | --- | --- | --- | --- |
| F04 | `high` | Closing imported settings keeps malicious config active | `0be5ba76b565` | src/frontend/src/main.js <br> src/app.go | 即時修正対象 |
| F05 | `high` | Release workflow trusts mutable Actions with write token | `c9ff4a107af2` | github/workflows/release.yml | 重大な運用判断待ち |
| F29 | `medium` | Unbounded native clipboard PNG copy can exhaust memory | `933c8de17c46` | src/internal/appcore/clipboard.go <br> src/internal/appcore/clipboard_native_windows.go <br> src/internal/appcore/processor.go <br> src/internal/appcore/image.go | HEAD確認済み |
| F30 | `medium` | Auto-photo scan cap can skip and later upload old images | `919ed595ef5f` | src/internal/appcore/autophoto.go | 即時修正対象 |
| F32 | `medium` | Directory picker stats attacker-controlled UNC paths | `94419e78d92e` | src/frontend/src/main.js <br> src/app.go | 即時修正対象 |
| F34 | `medium` | Dropped config becomes active without confirmation | `547f84a75e38` | frontend/src/main.js <br> app.go <br> internal/appcore/processor.go | 即時修正対象 |

### 19:48版の分類差分

#### 重大な仕様変更・運用判断が必要

- 19:48-F05のために、Release workflowのmutable Action参照、GITHUB_TOKENのwrite権限、Release asset公開stepをどう分離・固定するか判断が必要である。Action SHA pinningだけでよいか、署名job分離や権限分割まで含めるかでRelease運用が変わるため、今回の即時修正対象から除外する。
- 前回から継続して、19:48-F01/F02/F06/F08/F10/F16/F33相当は、タブ式設定確認、署名鍵隔離、auto-capture Discord同意境界、原文引用redaction、隠しcamera設定、metadata失敗時upload、history token保存モデルの判断待ちとして扱う。

#### 仕様変更不要だが追加情報・方針確認が必要

- 19:48-F45/F46相当の外部URL許可リストは、追加情報 `e9333711518881918aca4b091f8c2741` と `72a14374ab94819184eedb836966f795` により、セキュリティではなく固定リンクの機能退行と確認した。方針は「検索エンジンhostは許可せず、固定用途の公式/参照先hostだけを許可する」とし、即時修正対象へ移した。

#### 仕様変更不要かつ追加情報不要: 即時修正対象

- 19:48-F04/F34のために、Dropped/OpenSettingsで読み込んだ外部configを「即座に `a.configPath` へ採用する」処理から「保存前previewとして保持し、明示保存時だけ既存のactive configへ反映する」処理へ変更する。
- 19:48-F32のために、directory pickerの既定ディレクトリ決定を「ユーザー制御文字列へ直接 `os.Stat` する」処理から「UNC/device/networkなど危険なpathをstat前に既定値なしへ落とす」処理へ変更する。
- 19:48-F30のために、auto-photo watcherの初期baselineを「scan capに入った先頭N件だけをseen扱いする」処理から「既存の対象画像を後から新規扱いしない」処理へ変更する。
- 19:48-F29は、前回F41対応でWindows native clipboard PNGを `GlobalSize` 検証後にcopyする処理へ変更済みのため、HEAD確認済みとして扱う。

### 19:48版 作業割り当て

| Work | 対象finding | 担当 | 状態 | 主なファイル | メモ |
| --- | --- | --- | --- | --- | --- |
| W8 | 19:48-F04, 19:48-F32, 19:48-F34 | Pascal | 完了 | `src/app.go`, tests | imported config preview / directory picker境界 |
| W9 | 19:48-F30 | Plato | 完了 | `src/internal/appcore/autophoto.go`, tests | auto-photo scan cap baseline |
| W10 | 旧F05/旧F12, 19:48-F45, 19:48-F46 | メインエージェント | 完了 | `src/app.go`, `src/frontend/src/main.js`, tests | 外部URL許可リスト方針決定 |

### 19:48版 作業結果

#### W8: imported config preview / directory picker境界

- 19:48-F04/F34のために、`OpenSettings(path)` を「渡されたJSON pathを即座に `a.configPath` へ採用する」処理から「active `configPath` を維持したままimported configを未保存previewとして `state.Config` へ表示する」処理へ変更し、Dropped configを保存前にruntime設定として使わないようにした。
- 19:48-F04/F34のために、`CloseSettings()` を「設定画面を閉じるだけ」から「`SettingsBaselineConfig` がある場合はactive configへ戻してdraft metadataを解除する」処理へ変更し、閉じる/破棄でimported configが残らないようにした。
- 19:48-F32のために、directory pickerの既定ディレクトリ決定を「ユーザー制御のcurrent pathへ直接 `os.Stat` する」処理から「UNC/device/network形式を `pickerDefaultDirectory` で空にしてからstatする」処理へ変更し、pickerボタンで攻撃者管理のUNC pathへ接続しないようにした。

#### W9: auto-photo scan cap baseline

- 19:48-F30のために、auto-photo watcherの起動時baseline scanを「`MaxAutoPhotoScanFiles` で打ち切った先頭N件だけを `seen` に入れる」処理から「既存の対象画像を全件 `seen` に入れる」処理へ変更し、上限外にあった古い画像が後から新規扱いで投稿されないようにした。
- 19:48-F30のために、tick時のscanを「処理候補の発見自体をscan capで打ち切る」処理から「全件scanで既存/新規を判定し、実際の処理件数は既存のper-tick上限で制限する」処理へ変更し、古い画像の誤投稿防止と1tickあたりの処理上限を分離した。
- 旧F43の追加情報 `e2dbfdd09f68819187b842db37a4421f` のために、W9修正後の `Run` / `tick` が全件scanを使うことを再確認し、追加情報で指摘された旧画像再浮上と新規画像starvationが同じ修正で解消済みであることを管理レポートへ追記した。

#### W10: 外部URL許可リスト方針

- 旧F05/19:48-F45のために、診断説明内のGPG案内リンクを「`www.google.com` の検索結果を開く」処理から「GnuPG公式サイト `https://gnupg.org/` を開く」処理へ変更し、検索エンジンhostを外部URL許可リストへ追加しない方針にした。
- 旧F05/19:48-F45のために、`trustedExternalURL` を「既存4hostのみ許可」から「`gnupg.org` と `keys.openpgp.org` も許可する」処理へ変更し、診断説明のGPG公式情報と作者公開鍵リンクだけを開けるようにした。
- 旧F12/19:48-F46のために、`trustedExternalURL` を「VRChat Feedback hostを拒否する」処理から「`feedback.vrchat.com` を許可する」処理へ変更し、Camera OSC不具合報告の固定ヘルプリンクを開けるようにした。
- 旧F05のために、`www.google.com` はテストで拒否されるhostとして固定し、検索エンジン全体を許可リストへ広げないことを確認できるようにした。

### 19:48版 テスト結果

- `cd src && go test ./...` を実行し、成功した。
- `node scripts/check-frontend-template-literals.mjs` を実行し、成功した。
- `node scripts/check-wails-api-surface.mjs` を実行し、成功した。
- `cd src && go test ./internal/appcore -run 'TestScanPhotoFilesAllWithExcludesStatusIgnoresScanCap|TestAutoPhotoWatcherFullScanBaselineDoesNotResurfaceDeletedOldFiles|TestAutoPhotoWatcherFullScanBaselineProcessesLateArrivingFile' -count=1` を追加実行し、旧F43追加情報の再検証として成功した。
- Pascalが `go test ./... -run 'TestAppSaveConfigAndOpenSettings|TestAppOpenSettingsImportedConfigIsPreviewUntilSaved|TestAppOpenSettingsKeepsConfigPathOnLoadError|TestPickerDefaultDirectoryRejectsUnsafeNetworkPaths'` を実行し、成功した。
- Platoが `go test ./internal/appcore -run '^(TestScanPhotoFilesLimitsFileCount|TestScanPhotoFilesReportsMissingDirectory|TestScanPhotoFilesAllWithExcludesStatusIgnoresScanCap|TestAutoPhotoWatcherFullScanBaselineDoesNotResurfaceDeletedOldFiles|TestAutoPhotoWatcherFullScanBaselineProcessesLateArrivingFile)$'` を実行し、成功した。

### 19:48版 コミット

- `9be15d1 fix(settings): keep imported configs as preview`: 19:48-F04/F32/F34のために、imported config preview、CloseSettingsのbaseline復帰、picker default pathのUNC/device/network拒否を実装した。
- `53b74d5 fix(autophoto): baseline full scans before tick limits`: 19:48-F30のために、auto-photoの全件baseline scanとper-tick処理上限を分離した。
- `65c580e fix(help): keep diagnostic links on trusted hosts`: 旧F05のために、GPG説明リンクをGoogle検索からGnuPG公式へ変更した。旧F05/旧F12のhost許可とテストは `9be15d1` に含まれる。
- `docs(issues): record 1948 security remediation work`: 19:48 CSV、追加情報、分類、作業結果、残課題を記録する。

### 19:48版 残課題

- 方針決定済み・実装待ち: `issues/310` は自動撮影Discord投稿を専用ONだけで有効化する。`issues/311` は非表示camera自動起動/終了をバックエンド側で無効化する。`issues/314` はissue原文引用よりsecret redactionを優先する。`issues/315` は非表示タブのWebhook/自動投稿差分を保存前に具体表示する。`issues/316` はRelease署名secretをbuild jobからsign jobへ分離する。
- 方針未決・追加解説済み: `issues/313` はhistory JSON/UI stateにDiscord tokenを保存しない短期案、Discord削除機能を落とす案、OS秘密情報ストアへ分離する案の影響を追加解説した。ユーザー判断待ち。
- 将来提案: `docs/security-future-recommendations.md` に、313のOS秘密情報ストア分離、314のredaction検査、316のkeyless/KMS release signingを記録した。
- 旧F43/19:48-F30相当のauto-photo scan cap policyは、追加情報再検証によりW9修正で対応済みへ移した。
- 追加情報・方針確認待ちは、旧F05/旧F12の方針決定により今回分では解消した。
- HEAD確認済み・重複・軽減済み: 19:48-F13/F14/F27/F29/F38/F47相当は前回修正または確認済みとして扱う。
