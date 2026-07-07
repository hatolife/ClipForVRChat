# issue原文引用で秘密情報を永続化しない

## 指示

> 現在の作業をできるだけサブエージェントに任せて次の作業を実施 
> 重大な仕様変更・運用判断が必要　の項目について追加情報があります 
> '[REDACTED: local user path]/reports/security/2026-07-07T17-53-16.203Z/b825b8f570a881919546b22faf5da6cf' '[REDACTED: local user path]/reports/security/2026-07-07T17-53-16.203Z/394a0dc95c688191a1617b65fcb2befd' '[REDACTED: local user path]/reports/security/2026-07-07T17-53-16.203Z/e2dbfdd09f68819187b842db37a4421f' '[REDACTED: local user path]/reports/security/2026-07-07T17-53-16.203Z/958fd884b9588191ab138b574a6bd5eb' '[REDACTED: local user path]/reports/security/2026-07-07T17-53-16.203Z/69b57f2de0b081919021d7b078b7778c' '[REDACTED: local user path]/reports/security/2026-07-07T17-53-16.203Z/fa5bc68dae4c8191ad9140983e5ac66f' '[REDACTED: local user path]/reports/security/2026-07-07T17-53-16.203Z/f23b0a6320c481919bc02541a59d01da'
> ７件について個別にissue作成 それぞれ対応方針を３パターンissueに記載

## 文脈

追加情報 `69b57f2de0b081919021d7b078b7778c` は、AGENTS.mdの「ユーザー発言を原文のまま引用」ルールに、Webhook URL、token、ローカルパス、ログなどの秘密情報を伏せる例外がないと指摘している。

## 解釈

ユーザー発言の意図や誤字を残す運用は維持しつつ、秘密情報・プライバシー情報はissueへ永続化しない例外を明文化する必要がある。

## 問題

- issueはリポジトリにコミットされ、closed配下にも残る永続記録である。
- 原文引用を厳格に守ると、ユーザーが誤って貼ったWebhook URLやtokenも残る。
- 既存の診断ログではredaction方針があるが、issue作成ルールには同等の方針がない。

## 期待する挙動

issueの「指示」は原文の意味・誤字をできるだけ保つが、秘密情報やプライバシー情報は `[REDACTED: ...]` などへ置換してから記録する。
置換した事実と種類は分かるが、生値は残さない。

## 対応方針案

- A: AGENTS.mdに、原文引用より秘密情報redactionを優先する例外を追記する。
- B: issue作成前にredactionチェックを行うスクリプトを追加し、Webhook/tokenらしき値があると失敗させる。
- C: issue本文にはユーザー発言を要約だけ記載し、原文引用ルールを廃止する。

## 方針評価

- A: すぐ効く運用修正。現行の「原文を残す」要求とセキュリティ例外を両立できる。
- B: 抜け漏れ防止として有効だが、誤検知・対象パターン保守が必要なためAの補助策。
- C: 秘密情報は残りにくいが、ユーザーが求めた原文保持の監査性を失う。

## 推奨方針

Aを即時対応、Bを補強策として後続対応にする。
Cはユーザーが求めた「原文を残す」運用を失うため採用しない。

## 方針決定

2026-07-08: ユーザー判断によりAを採用する。
issueの「指示」では原文引用を維持するが、Discord Webhook URL、token、ローカルユーザー名を含むパス、QR URL、診断ログ内secretなどは原文引用よりredactionを優先する。
検査スクリプト追加は将来提案として `docs/security-future-recommendations.md` に分離する。

2026-07-08追記: ユーザー指示「Working directory: [REDACTED: local user path]. You are not alone in the codebase; do not revert edits made by others. Implement issue 314. Ownership: AGENTS.md, and if you choose a lightweight check script then scripts/ and package-adjacent docs/tests for that script. Requirements: Update AGENTS.md so issue '指示' should quote user wording verbatim except secrets/private data must be redacted before committing. Include examples: Discord Webhook URLs/tokens, local paths with user names, QR URLs, diagnostic log secrets. Prefer a minimal documentation/process fix unless a simple script already fits existing patterns. Run any relevant check if you add/change scripts. Commit your changes with a conventional commit message. Final response: summary, changed files, tests/checks run, commit hash if committed.」に基づき、軽量スクリプトは追加せずAGENTS.mdの運用ルールを更新する。

## 受け入れ条件

- [x] AGENTS.mdに、秘密情報・プライバシー情報は原文引用からredactする例外を明記する。
- [x] redaction対象例にDiscord Webhook URL/token、ローカルユーザー名を含むパス、QR URL、診断ログのsecretを含める。
- [x] redaction済み箇所は `[REDACTED: discord webhook url]` のように種類が分かる表記へ統一する。
- [x] redaction後もユーザー指示の意図と作業判断に必要な文脈が残る。
- [x] 今後のissue作成時にこの例外を参照できる。
