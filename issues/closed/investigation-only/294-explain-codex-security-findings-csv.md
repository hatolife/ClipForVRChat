# Codex Security findings CSVを解説する

## 指示

> '/mnt/c/Users/user/Downloads/codex-security-findings-2026-07-07T18-53-06.315Z.csv' これ解説

> '/home/user/work/ClipForVRChat/reports/security/2026-07-07T17-53-16.203Z/codex-security-findings-2026-07-07T19-00-29.151Z.csv'　数件追加されたのでこれを底本として　再度分析

## 文脈

指定されたCSVは、Codex Security Findingsからエクスポートされた `hatolife/ClipForVRChat` 向けの指摘一覧である。
再分析の底本は `reports/security/2026-07-07T17-53-16.203Z/codex-security-findings-2026-07-07T19-00-29.151Z.csv` で、44件のfindingを含む。
内訳は `high` 3件、`medium` 25件、`low` 10件、`informational` 6件である。
設定インポート、自動投稿、診断ログ、OSC、履歴、QR、クリップボード、リソース消費、Release署名、CI/Release、WebView drop、URL allow-listなどの経路に関する情報漏えい、同意回避、DoS、危険な外部URL/mention誘発、Release信頼性低下を指摘している。

## 解釈

CSV全体を読み、各findingの意味、影響、優先度感、重複・関連する指摘を日本語で分かりやすく説明する。
この依頼ではコード修正までは行わず、現時点の理解と対応方針を整理する。

## 問題

- 新CSVには44件のfindingが含まれている。
- 旧CSV24件に加えて20件が追加され、`Tabbed settings hide webhook auto-post settings` は `medium` から `high` へ上がっている。
- 新しい `high` は、設定タブ分割によるWebhook/自動投稿設定の見落とし、Release署名秘密鍵のCI runner露出、診断ログへのDiscord webhook token永続化である。
- 多くは「ユーザーが明示的にONにしたつもりのない自動投稿・監視・Webhookが有効になる」問題である。
- 一部は巨大入力や大量OSC packetによるメモリ・CPU・ログDoSである。
- 一部は診断ログ、履歴、issue運用に秘密情報やプライバシー情報が残る問題である。
- 既存issueと重なる指摘も含まれており、対応時には統合管理が必要である。

## 期待する挙動

- 新CSVの列と全体像が説明される。
- 44件のfindingが日本語で要約される。
- どれから優先して対応すべきか分かる。
- 既存issueと重なるものが分かる。

## 受け入れ条件

- [x] CSVの性質と列の意味を説明する。
- [x] 新CSV全44件のfindingの要旨を説明する。
- [x] リスクの種類ごとに分類する。
- [x] 優先対応候補を整理する。
- [x] 既存issueと関連する指摘を明記する。

## 調査結果

- CSVは `finding_url`、`repository`、`title`、`description`、`severity`、`status`、`detected_at`、`committed_at`、`commit_hash`、`relevant_paths` などを持つCodex Security Findingsの一覧である。
- 新CSVは44件すべてが `status=new` で、現時点では「ツール上は未解決扱い」の指摘として出力されている。
- severity内訳は `high` 3件、`medium` 25件、`low` 10件、`informational` 6件である。
- 旧CSVから削除されたfindingはない。
- 新CSVで追加されたfindingは20件である。
- 旧CSVからseverityが変わったfindingは `Tabbed settings hide webhook auto-post settings` の1件で、`medium` から `high` へ上がっている。
- 最重要グループは、自動投稿・Webhook・監視フォルダ・未保存draft・隠れた設定により、利用者の意図しないDiscord投稿が起きる可能性がある指摘である。
- 新CSVではRelease署名秘密鍵と診断ログ内Webhook tokenの2件も最優先対応候補に入る。
- 既存の `293-guard-imported-auto-processing-watch-directories.md` は、`Hidden auto-processing folders can exfiltrate images` と強く関連する。
- `Raw issue quoting can leak user secrets` は、このリポジトリのissue運用ルール自体が秘密情報を永続化し得るというメタな指摘である。
- DoS系は、CLI暗号化、ZIP暗号化、QR拡大、クリップボードPNG、OSCログ、OSC cache、OSC forwarding self-loop、localhost IPC、Spout sender dimensions、auto-photo scan error増殖に分かれる。
- プライバシー/秘密情報系は、VRChat private instance ID、起動診断のローカルパス、history内Discord token、診断ログ内Discord webhook token、debug OSC引数、issue原文引用に分かれる。
- UI/同意系は、設定タブ整理、隠れた設定、カメラmask設定のsilent reset、URL allow-list不整合により、保存・自動処理・Webhook投稿・操作導線の実際の挙動をユーザーが確認しにくくなっていることが中心である。

## 旧CSVからの差分

- 件数は24件から44件へ増えた。
- 追加は20件で、内訳は `high` 2件、`medium` 2件、`low` 10件、`informational` 6件である。
- `Tabbed settings hide webhook auto-post settings` は既存findingだが、severityが `medium` から `high` へ上がった。
- 旧CSVの24件はすべて新CSVにも残っている。
- 新しい最重要追加は `Release signing key exposed to prior build-step code execution` と `Diagnostic log can persist Discord webhook tokens` である。
- `Screenshots are sent to the auto-photo Discord webhook` と `Update banner opens unvalidated release URL` は、旧CSVの同種findingと重複または補強関係にある。

## Finding別要約

1. `Remote-player mask opt-out is silently reset` (`low`): 他ユーザーを写さないmask設定をOFFにしてもNormalizeで既定値trueへ戻り、Discord投稿時に意図せず他ユーザーが写り得る。
2. `Auto-capture Discord opt-out bypass leaks captures` (`medium`): 自動撮影専用Discord投稿をOFFにしていても、全体のDiscord投稿ONや古い専用Webhook設定により自動撮影画像が投稿され得る。
3. `Auto-photo skip path bypasses per-tick processing cap` (`low`): skip対象ファイルがtickごとの処理上限に数えられず、安定判定sleepを大量発生させてwatcherを長時間止め得る。
4. `Unbounded CLI encryption can exhaust memory` (`medium`): 画像以外の起動引数を暗号化処理へ回し、任意ファイル/フォルダを全量メモリ展開するため、大きい入力でクラッシュし得る。
5. `New diagnostic links are blocked by URL allow-list` (`informational`): 新しい診断案内リンクがURL許可リストに含まれず、クリックしても開けない機能退行。
6. `Unrelated OSC parameters can suppress fallback` (`low`): 関係ないOSC avatar parameterがAvatarBeacon受信扱いになり、自動fallbackを抑止して撮影失敗を招き得る。
7. `Raw issue quoting can leak user secrets` (`medium`): ユーザー発言を原文のままissueへ保存する運用により、Webhook URL、token、ログ、ローカルパスなどがGit履歴へ残り得る。
8. `Unsaved draft can activate auto-post watchers` (`medium`): 未保存設定draftが起動時のruntime configとして使われ、自動投稿watcherが保存済み設定ではなくdraftのWebhookで動き得る。
9. `Debug OSC input is written verbatim to diagnostics` (`low`): debug OSC入力の引数値がそのまま診断ログへ残り、ユーザーが送った秘密文字列が診断パッケージへ入る可能性がある。
10. `Spout diagnose trusts unbounded sender dimensions` (`low`): Spout senderが返す幅/高さを上限なしでbuffer確保に使い、helper停止やunsafeなbuffer sizingにつながり得る。
11. `Hidden camera auto-start settings remain active` (`medium`): UIから隠したカメラ自動起動/終了設定が既存設定や外部config由来では残り、意図せずOSCでカメラ操作され得る。
12. `New VRChat feedback help link is blocked` (`informational`): VRChat feedbackリンクがURL許可リストに含まれず、ヘルプボタンが開けない。
13. `Diagnostic logs leak VRChat private instance IDs` (`medium`): 診断ログにVRChatのworld/instance IDが生で残り、private instanceの所有者やnonce相当情報が共有され得る。
14. `OSC forwarding can self-loop and cause UDP DoS` (`medium`): `0.0.0.0` bind時に同じローカルportへ転送する設定を検出しきれず、受信と転送のループでCPU/ログを消費し得る。
15. `Unbounded localhost IPC request can exhaust app memory` (`low`): 単一起動用localhost IPCがrequest size制限なしでJSON decodeし、巨大token等でメモリを消費し得る。
16. `Unbounded OSC packet logging enables log DoS` (`medium`): OSC packetごとに診断ログへ追記し、回転・容量制限・throttleがないため、UDP floodでディスクや処理時間を消費し得る。
17. `Auto-post warning suppressed when webhook is configured` (`medium`): Webhookが設定済みの場合に保存時の自動投稿警告が出なくなり、外部config保存時の同意確認が弱くなる。
18. `Unbounded OSC avatar parameter cache enables DoS` (`medium`): 任意の `/avatar/parameters/...` を無制限にmapへ蓄積し、メモリ、CPU、ログを消費し得る。
19. `Metadata failures can lead to unvalidated Discord uploads` (`medium`): メタデータ書き込み失敗を警告だけにした結果、画像として不正な拡張子偽装ファイルでもDiscordへ送られ得る。
20. `Release notes step uses wrong env syntax on Windows` (`informational`): Windows Actions上で `$TAG_NAME` を使っており、Release notes抽出が失敗してRelease buildが止まり得る。
21. `Unbounded ZIP CLI encryption can exhaust memory` (`medium`): zip暗号化CLIがzip全体と暗号化結果をメモリに載せるため、大きなzipや特殊ファイルで停止し得る。
22. `Unbounded ZIP encryption can exhaust memory` (`low`): 21と同種で、`.zip` 拡張子の巨大ファイルや非zip偽装ファイルを全量読み込みしてOOMになり得る。
23. `COM shell calls are made without OS-thread pinning` (`informational`): Windows Shell APIを呼ぶCOM初期化/終了を同一OS threadへ固定しておらず、Explorer表示で不安定化し得る。
24. `Startup diagnostics log unredacted local paths` (`medium`): 起動診断ログにexe/config/history/output/VRChat写真/Screenshotフォルダなどのローカルパスが生で残る。
25. `Untrusted history paths are passed to Explorer` (`medium`): history由来の `outputPath` を信頼してExplorerへ渡し、UNC pathアクセスや任意path probeにつながり得る。
26. `Explorer reveal trusts history-controlled file paths` (`medium`): 25と同種で、rendererから渡された任意pathをExplorerへ渡せる設計が問題。
27. `Output format disabled for Discord-only processing` (`informational`): Discordのみ出力時にもbackendは画像format/qualityを使うのに、UI上は設定できない機能不一致。
28. `Output format controls disabled for Discord-only output` (`informational`): 27と同種で、local save OFF時にDiscord/clipboard出力のformat/quality調整ができない。
29. `Hidden auto-processing folders can exfiltrate images` (`medium`): UIに出ない監視フォルダ設定が外部configから残り、攻撃者指定フォルダの画像をWebhookへ送信し得る。既存issue 293と関連。
30. `Tabbed settings hide webhook auto-post settings` (`high`): 設定タブを見ないまま外部config全体を保存でき、Webhook/自動投稿設定を確認しないまま画像流出につながり得る。
31. `Auto-photo scan errors are appended without deduplication` (`low`): missing/inaccessible/scan limitエラーがscanごとに結果へ追加され続け、メモリ/UIが増え続ける。
32. `Auto-photo scan errors cause unbounded UI result growth` (`low`): 31と同種で、同じscan errorがDOM/result listへ無制限に増殖する。
33. `Screenshots are sent to the auto-photo Discord webhook` (`medium`): screenshot自動投稿が通常WebhookではなくVRChat自動投稿Webhookへfallbackし、想定外チャンネルへ送信され得る。
34. `Screenshots can upload to the wrong Discord webhook` (`medium`): 33と同種で、スクリーンショットが想定外のDiscord webhookへ送られる。
35. `Update banner opens unvalidated release URL` (`medium`): update API由来のURLを検証せず開くため、改ざん応答時に任意URLやcustom protocolを開き得る。
36. `Unvalidated update URL can open attacker-controlled links` (`medium`): 35と同種で、update bannerが信頼済み表示のまま攻撃者URLへ誘導し得る。
37. `Release signing key exposed to prior build-step code execution` (`high`): Release署名秘密鍵がbuild/test/dependency実行後の同一jobで使われ、先行ステップがrunner状態を汚染すると鍵を読まれ得る。
38. `Diagnostic log can persist Discord webhook tokens` (`high`): Discord送信失敗時の `err.Error()` にwebhook URL/tokenが含まれ、診断ログへ永続化され得る。
39. `QR upscaling can exhaust memory on crafted tall images` (`medium`): QR検出用の2x/3x/4x拡大が高さや総pixel数を見ず、縦長画像で巨大メモリを使い得る。
40. `Untrusted QR URLs can trigger Discord mentions` (`medium`): QR内URLをDiscord本文にそのまま入れるため、URL path/query内のmention文字列がrole/user/everyone通知を発火し得る。
41. `Unbounded native clipboard PNG copy can crash app` (`medium`): Windows clipboardのPNG HGLOBALをサイズ制限前に全量コピーし、大きいclipboard dataでクラッシュし得る。
42. `Accepted drops are not cancelled, allowing WebView navigation` (`low`): drop eventをcancelしておらず、WebViewが外部URL/ファイルへnavigateしてWails bridge露出につながり得る。
43. `Auto-photo scan cap can leak old photos and starve new ones` (`medium`): 起動時baselineが先頭5000件だけなので、古い写真が後から新規扱いで投稿されたり、新規写真が永久に処理されない可能性がある。
44. `History file stores Discord tokens with weak permissions` (`medium`): historyにDiscord webhook token等を保存し、権限も弱いため、他プロセス/ユーザーから投稿・削除権限が漏れ得る。

## 優先対応候補

1. `Diagnostic log can persist Discord webhook tokens`: live tokenが診断ログへ残るため、ログ書き込み前のURL/token redactionを最優先で入れる。
2. `Release signing key exposed to prior build-step code execution`: Release署名の信頼性に直結するため、署名jobをbuild jobから分離し、artifact入力だけを署名する構成へ変える。
3. `Tabbed settings hide webhook auto-post settings` / `Hidden auto-processing folders can exfiltrate images` / `Auto-post warning suppressed when webhook is configured`: 外部config保存時のWebhook、監視フォルダ、自動投稿ONを必ず確認対象にする。
4. `Unsaved draft can activate auto-post watchers` / `Auto-capture Discord opt-out bypass leaks captures` / `Screenshots ... wrong webhook`: runtime configとWebhook選択の同意境界を整理する。
5. `History file stores Discord tokens with weak permissions`: historyへtokenを保存しない、保存する場合は最小化・権限制限・UI state非露出にする。
6. DoS系: 入力サイズ、件数、ログ、cache、UI result、IPC requestに上限とrate limitを入れる。
