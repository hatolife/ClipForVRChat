# Discord履歴のWebhook削除トークン保存を保護する

## 指示

> 現在の作業をできるだけサブエージェントに任せて次の作業を実施 
> 重大な仕様変更・運用判断が必要　の項目について追加情報があります 
> '/home/user/work/ClipForVRChat/reports/security/2026-07-07T17-53-16.203Z/b825b8f570a881919546b22faf5da6cf' '/home/user/work/ClipForVRChat/reports/security/2026-07-07T17-53-16.203Z/394a0dc95c688191a1617b65fcb2befd' '/home/user/work/ClipForVRChat/reports/security/2026-07-07T17-53-16.203Z/e2dbfdd09f68819187b842db37a4421f' '/home/user/work/ClipForVRChat/reports/security/2026-07-07T17-53-16.203Z/958fd884b9588191ab138b574a6bd5eb' '/home/user/work/ClipForVRChat/reports/security/2026-07-07T17-53-16.203Z/69b57f2de0b081919021d7b078b7778c' '/home/user/work/ClipForVRChat/reports/security/2026-07-07T17-53-16.203Z/fa5bc68dae4c8191ad9140983e5ac66f' '/home/user/work/ClipForVRChat/reports/security/2026-07-07T17-53-16.203Z/f23b0a6320c481919bc02541a59d01da' 
> ７件について個別にissue作成 それぞれ対応方針を３パターンissueに記載

## 文脈

追加情報 `958fd884b9588191ab138b574a6bd5eb` は、履歴JSONにDiscord webhook削除用のID/token/message IDが保存され、UI stateにも露出すると指摘している。
削除機能を維持するには何らかの削除メタデータが必要だが、保存範囲と露出範囲を見直す必要がある。

## 解釈

履歴からDiscord投稿を削除できる利便性と、webhook tokenを永続化しない安全性のどちらを優先するかの仕様判断が必要。
履歴ファイル権限の強化だけでは、tokenを複製して保存する根本問題が残る。

## 問題

- historyにDiscord webhook tokenがJSONとして永続化される。
- POSIX環境では広いpermissionで保存されると他ユーザーから読まれる可能性がある。
- UI state/results/historyにもtokenが含まれ、renderer側の漏えい面が広がる。

## 期待する挙動

Discord webhook tokenはhistory JSONにもフロントエンドへ返すstateにも保存しない。
履歴ファイルはprivate permissionで保存し、削除機能は現在の設定から安全にtokenを解決できる場合だけ実行する。

## 対応方針案

- A: historyには `messageID`、`webhookID`、Webhook種別/参照だけを保存し、tokenは現在のconfigのWebhook URLからID一致時にだけ解決する。history/UI stateへtokenは保存しない。
- B: Discord投稿の履歴削除機能を廃止または現セッション限定にし、永続履歴には削除用metadataを保存しない。
- C: tokenをOS credential storeまたは暗号化ストアへ分離し、historyにはopaque参照IDだけ保存する。

## 方針評価

- A: 短期でtoken重複保存を止められる。Webhook URL自体はconfigに残るが、historyへの追加複製とrenderer露出を消せる。
- B: 最も安全だが、履歴からDiscord投稿を削除する既存機能を落とすため利用者影響が大きい。
- C: 最もきれいな設計だが、OS別実装・移行・バックアップ時の扱いが大きく、今回の即時修正範囲を超える。

## 推奨方針

Aを短期対応、Cを将来改善として扱う。
「private permissionにしたうえでtokenをhistoryへ残す」案は不十分なため採用しない。
現在configからtokenを解決できない古い履歴のDiscord削除は失敗を許容し、ローカル履歴削除だけ可能にする。

## 方針決定

2026-07-08: ユーザー判断によりAを採用する。
history JSONとUI stateにはDiscord tokenを保存しない。
履歴には `messageID`、`webhookID`、Webhook種別/参照だけを保存し、Discord削除時は現在configのWebhook URLからID一致でtokenを一時解決できる場合だけ削除APIを呼ぶ。
tokenを解決できない古い履歴ではDiscord削除不可を明示し、ローカル履歴削除だけ可能にする。

## 追加解説

このissueで守りたい対象は、Discord webhook URLそのものではなく、履歴機能が追加で保持する削除用tokenの複製である。
Discord webhook URLは現在の設定に存在するため、同じユーザー権限のプロセスから完全に秘匿する設計ではない。
ただし、history JSONとUI stateへtokenを重複保存すると、履歴ファイル共有、診断出力、renderer state漏えい、古いバックアップからの復元など、configとは別の経路でtokenが残り続ける。

A案は、履歴には削除対象を識別する `messageID` と `webhookID`、および「通常Webhook由来」「自動撮影Webhook由来」などの参照情報だけを保存する。
削除実行時は現在のconfig内のWebhook URLを解析し、履歴の `webhookID` と一致した場合だけtokenを一時的に使う。
そのため、Webhook URLを変更した後の古い履歴はDiscord側の削除ができなくなるが、historyにtokenを永続化しない代償として許容する設計になる。

B案は最も安全だが、アプリ再起動後にDiscord投稿を履歴から削除する機能を失う。
既存ユーザーにとっては機能削除に近く、誤投稿時の撤回導線も弱くなる。
「Discord上の削除はアプリ終了まで」と明示するなら成立するが、仕様変更の影響はAより大きい。

C案は長期的には最も整っている。
historyにはopaque IDだけを置き、tokenはWindows Credential ManagerなどOS側の秘密情報ストアへ分離できる。
一方で、OS別実装、portable運用、バックアップ/移行、credential削除、テスト環境での扱いが増えるため、今回の即時修正としては大きい。

判断の基準は、再起動後のDiscord削除機能を維持したいかである。
今回はA採用のため、再起動後の削除機能は「現在configから同じWebhookのtokenを解決できる範囲で維持する」扱いにする。
OS別の秘密情報ストア実装は将来提案として扱う。

## 受け入れ条件

- [ ] history保存先ディレクトリとhistoryファイルがprivate permissionで作成される。
- [ ] `history.json` にDiscord tokenの生値が保存されない。
- [ ] UIへ返す `Result` / `HistoryEntry` にDiscord tokenの生値が含まれない。
- [ ] バックエンドのDiscord削除処理は、現在configのWebhook URLからID一致でtokenを解決できる場合だけ実行する。
- [ ] tokenを解決できない履歴では、Discord削除不可を明示し、ローカル履歴削除は可能にする。
- [ ] 既存履歴の読み込み互換を維持する。
- [ ] tokenがJSON/API stateへ不要に露出しないことをテストで確認する。
