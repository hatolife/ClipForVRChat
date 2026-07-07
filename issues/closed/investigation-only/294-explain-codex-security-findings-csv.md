# Codex Security findings CSVを解説する

## 指示

> '/mnt/c/Users/user/Downloads/codex-security-findings-2026-07-07T18-53-06.315Z.csv' これ解説

## 文脈

指定されたCSVは、Codex Security Findingsからエクスポートされた `hatolife/ClipForVRChat` 向けの指摘一覧である。
24件すべてが `severity=medium`、`status=new` で、設定インポート、自動投稿、診断ログ、OSC、履歴、QR、クリップボード、リソース消費などの経路に関する情報漏えい、同意回避、DoS、危険な外部URL/mention誘発を指摘している。

## 解釈

CSV全体を読み、各findingの意味、影響、優先度感、重複・関連する指摘を日本語で分かりやすく説明する。
この依頼ではコード修正までは行わず、現時点の理解と対応方針を整理する。

## 問題

- CSVには24件のmedium findingが含まれている。
- 多くは「ユーザーが明示的にONにしたつもりのない自動投稿・監視・Webhookが有効になる」問題である。
- 一部は巨大入力や大量OSC packetによるメモリ・CPU・ログDoSである。
- 一部は診断ログ、履歴、issue運用に秘密情報やプライバシー情報が残る問題である。
- 既存issueと重なる指摘も含まれており、対応時には統合管理が必要である。

## 期待する挙動

- CSVの列と全体像が説明される。
- 24件のfindingが日本語で要約される。
- どれから優先して対応すべきか分かる。
- 既存issueと重なるものが分かる。

## 受け入れ条件

- [x] CSVの性質と列の意味を説明する。
- [x] 全24件のfindingの要旨を説明する。
- [x] リスクの種類ごとに分類する。
- [x] 優先対応候補を整理する。
- [x] 既存issueと関連する指摘を明記する。

## 調査結果

- CSVは `finding_url`、`repository`、`title`、`description`、`severity`、`status`、`detected_at`、`committed_at`、`commit_hash`、`relevant_paths` などを持つCodex Security Findingsの一覧である。
- 24件すべてが `medium/new` で、現時点では「ツール上は未解決扱い」の指摘として出力されている。
- 最重要グループは、自動投稿・Webhook・監視フォルダ・未保存draft・隠れた設定により、利用者の意図しないDiscord投稿が起きる可能性がある指摘である。
- 既存の `293-guard-imported-auto-processing-watch-directories.md` は、`Hidden auto-processing folders can exfiltrate images` と強く関連する。
- `Raw issue quoting can leak user secrets` は、このリポジトリのissue運用ルール自体が秘密情報を永続化し得るというメタな指摘である。
- DoS系は、CLI暗号化、ZIP暗号化、QR拡大、クリップボードPNG、OSCログ、OSC cache、OSC forwarding self-loopに分かれる。
- プライバシー/秘密情報系は、VRChat private instance ID、起動診断のローカルパス、history内Discord token、issue原文引用に分かれる。
- UI/同意系は、設定タブ整理や隠れた設定により、保存・自動処理・Webhook投稿の実際の挙動をユーザーが確認しにくくなっていることが中心である。

## Finding別要約

1. `Auto-capture Discord opt-out bypass leaks captures`: 自動撮影専用Discord投稿をOFFにしていても、全体のDiscord投稿ONや古い専用Webhook設定により自動撮影画像が投稿され得る。
2. `Unbounded CLI encryption can exhaust memory`: 画像以外の起動引数を暗号化処理へ回し、任意ファイル/フォルダを全量メモリ展開するため、大きい入力でクラッシュし得る。
3. `Raw issue quoting can leak user secrets`: ユーザー発言を原文のままissueへ保存する運用により、Webhook URL、token、ログ、ローカルパスなどがGit履歴へ残り得る。
4. `Unsaved draft can activate auto-post watchers`: 未保存設定draftが起動時のruntime configとして使われ、自動投稿watcherが保存済み設定ではなくdraftのWebhookで動き得る。
5. `Hidden camera auto-start settings remain active`: UIから隠したカメラ自動起動/終了設定が既存設定や外部config由来では残り、意図せずOSCでカメラ操作され得る。
6. `Diagnostic logs leak VRChat private instance IDs`: 診断ログにVRChatのworld/instance IDが生で残り、private instanceの所有者やnonce相当情報が共有され得る。
7. `OSC forwarding can self-loop and cause UDP DoS`: `0.0.0.0` bind時に同じローカルportへ転送する設定を検出しきれず、受信と転送のループでCPU/ログを消費し得る。
8. `Unbounded OSC packet logging enables log DoS`: OSC packetごとに診断ログへ追記し、回転・容量制限・throttleがないため、UDP floodでディスクや処理時間を消費し得る。
9. `Auto-post warning suppressed when webhook is configured`: Webhookが設定済みの場合に保存時の自動投稿警告が出なくなり、外部config保存時の同意確認が弱くなる。
10. `Unbounded OSC avatar parameter cache enables DoS`: 任意の `/avatar/parameters/...` を無制限にmapへ蓄積し、メモリ、CPU、ログを消費し得る。
11. `Metadata failures can lead to unvalidated Discord uploads`: メタデータ書き込み失敗を警告だけにした結果、画像として不正な拡張子偽装ファイルでもDiscordへ送られ得る。
12. `Unbounded ZIP CLI encryption can exhaust memory`: zip暗号化CLIがzip全体と暗号化結果をメモリに載せるため、大きなzipや特殊ファイルで停止し得る。
13. `Startup diagnostics log unredacted local paths`: 起動診断ログにexe/config/history/output/VRChat写真/Screenshotフォルダなどのローカルパスが生で残る。
14. `Explorer reveal trusts history-controlled file paths`: rendererから渡された任意pathをExplorerへ渡せるため、改ざんhistory経由でUNC pathアクセスや任意ファイル表示が起き得る。
15. `Untrusted history paths are passed to Explorer`: 14と同種で、history由来の `outputPath` を信頼してExplorer表示する設計が問題。
16. `Hidden auto-processing folders can exfiltrate images`: UIに出ない監視フォルダ設定が外部configから残り、攻撃者指定フォルダの画像をWebhookへ送信し得る。既存issue 293と関連。
17. `Tabbed settings hide webhook auto-post settings`: 設定がタブ分割されたことで、ユーザーがWebhook/自動投稿タブを見ないまま外部config全体を保存し得る。
18. `Screenshots can upload to the wrong Discord webhook`: screenshot自動投稿が通常WebhookではなくVRChat自動投稿Webhookへfallbackし、想定外チャンネルへ送信され得る。
19. `Unvalidated update URL can open attacker-controlled links`: update API由来のURLを検証せず開くため、改ざん応答時に任意URLやcustom protocolを開き得る。
20. `QR upscaling can exhaust memory on crafted tall images`: QR検出用の2x/3x/4x拡大が高さや総pixel数を見ず、縦長画像で巨大メモリを使い得る。
21. `Untrusted QR URLs can trigger Discord mentions`: QR内URLをDiscord本文にそのまま入れるため、URL path/query内のmention文字列がrole/user/everyone通知を発火し得る。
22. `Unbounded native clipboard PNG copy can crash app`: Windows clipboardのPNG HGLOBALをサイズ制限前に全量コピーし、大きいclipboard dataでクラッシュし得る。
23. `Auto-photo scan cap can leak old photos and starve new ones`: 起動時baselineが先頭5000件だけなので、古い写真が後から新規扱いで投稿されたり、新規写真が永久に処理されない可能性がある。
24. `History file stores Discord tokens with weak permissions`: historyにDiscord webhook token等を保存し、権限も弱いため、他プロセス/ユーザーから投稿・削除権限が漏れ得る。
