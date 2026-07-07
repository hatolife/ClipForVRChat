# issue原文引用で秘密情報を永続化しない

## 指示

> 現在の作業をできるだけサブエージェントに任せて次の作業を実施 
> 重大な仕様変更・運用判断が必要　の項目について追加情報があります 
> '/home/user/work/ClipForVRChat/reports/security/2026-07-07T17-53-16.203Z/b825b8f570a881919546b22faf5da6cf' '/home/user/work/ClipForVRChat/reports/security/2026-07-07T17-53-16.203Z/394a0dc95c688191a1617b65fcb2befd' '/home/user/work/ClipForVRChat/reports/security/2026-07-07T17-53-16.203Z/e2dbfdd09f68819187b842db37a4421f' '/home/user/work/ClipForVRChat/reports/security/2026-07-07T17-53-16.203Z/958fd884b9588191ab138b574a6bd5eb' '/home/user/work/ClipForVRChat/reports/security/2026-07-07T17-53-16.203Z/69b57f2de0b081919021d7b078b7778c' '/home/user/work/ClipForVRChat/reports/security/2026-07-07T17-53-16.203Z/fa5bc68dae4c8191ad9140983e5ac66f' '/home/user/work/ClipForVRChat/reports/security/2026-07-07T17-53-16.203Z/f23b0a6320c481919bc02541a59d01da' 
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

## 受け入れ条件

- [ ] AGENTS.mdに、秘密情報・プライバシー情報は原文引用からredactする例外を明記する。
- [ ] redaction対象例にDiscord Webhook URL/token、ローカルユーザー名を含むパス、QR URL、診断ログのsecretを含める。
- [ ] redaction済み箇所は `[REDACTED: discord webhook url]` のように種類が分かる表記へ統一する。
- [ ] redaction後もユーザー指示の意図と作業判断に必要な文脈が残る。
- [ ] 今後のissue作成時にこの例外を参照できる。
