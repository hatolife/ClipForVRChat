# 自動撮影Discord opt-outを自動撮影アップロードへ厳密に適用する

## 指示

> 現在の作業をできるだけサブエージェントに任せて次の作業を実施 
> 重大な仕様変更・運用判断が必要　の項目について追加情報があります 
> '/home/user/work/ClipForVRChat/reports/security/2026-07-07T17-53-16.203Z/b825b8f570a881919546b22faf5da6cf' '/home/user/work/ClipForVRChat/reports/security/2026-07-07T17-53-16.203Z/394a0dc95c688191a1617b65fcb2befd' '/home/user/work/ClipForVRChat/reports/security/2026-07-07T17-53-16.203Z/e2dbfdd09f68819187b842db37a4421f' '/home/user/work/ClipForVRChat/reports/security/2026-07-07T17-53-16.203Z/958fd884b9588191ab138b574a6bd5eb' '/home/user/work/ClipForVRChat/reports/security/2026-07-07T17-53-16.203Z/69b57f2de0b081919021d7b078b7778c' '/home/user/work/ClipForVRChat/reports/security/2026-07-07T17-53-16.203Z/fa5bc68dae4c8191ad9140983e5ac66f' '/home/user/work/ClipForVRChat/reports/security/2026-07-07T17-53-16.203Z/f23b0a6320c481919bc02541a59d01da' 
> ７件について個別にissue作成 それぞれ対応方針を３パターンissueに記載

## 文脈

追加情報 `b825b8f570a881919546b22faf5da6cf` は、自動撮影専用のDiscord投稿OFFが実効的なopt-outになっていないと指摘している。
現在は通常のDiscord投稿設定がONなら、自動撮影側のDiscord投稿がOFFでも自動撮影画像が投稿され得る。

## 解釈

自動撮影は通常の手動処理よりプライバシー影響が大きいため、自動撮影専用のDiscord投稿ON/OFFを実効ゲートにするか、通常Discord投稿ONとの連動を明確な仕様として再承認する必要がある。

## 問題

- `autoCapture.discord.enabled=false` でも `output.uploadDiscord=true` により自動撮影投稿パスへ入る。
- 自動撮影専用Webhook URLが残っている場合、専用投稿OFFでもそのURLが優先され得る。
- UI文言も「通常Discord投稿ON、またはこの設定ON」と説明しており、opt-outとして弱い。

## 期待する挙動

自動撮影のDiscord投稿は、ユーザーが自動撮影用の投稿を明示的にONにした場合だけ実行される。
自動撮影投稿OFF時は、自動撮影専用Webhook URLや通常投稿用Webhook URLへ画像・本文を送信しない。

## 対応方針案

- A: 自動撮影Discord投稿は `autoCapture.discord.enabled` のみで判定し、通常投稿用 `output.uploadDiscord` では有効化しない。
- B: 通常投稿ONと自動撮影Discord投稿ONの両方がONの場合だけ自動撮影を投稿し、通常投稿ON単独では投稿しない。
- C: 互換移行期間として、既存設定では警告を出し、次回保存時に自動撮影Discord投稿ON/OFFを明示選択させる。

## 方針評価

- A: 最も明確。自動撮影はプライバシー影響が大きいため、専用opt-inを単独の実行条件にするのが妥当。
- B: `output.uploadDiscord` を追加条件として残す理由が弱く、ユーザーに二重ONを要求して混乱を増やす。
- C: 互換には配慮できるが、警告期間中に漏えい経路を残すなら不十分。採用する場合も最終状態はAに寄せる。

## 推奨方針

Aを採用する。
通常処理の投稿ONと自動撮影の投稿ONは、利用者の同意範囲が異なるため分離する。
自動撮影専用Webhook URLは、投稿ONの場合だけ送信先候補として扱う。

## 受け入れ条件

- [ ] 自動撮影投稿の実行条件が `autoCapture.discord.enabled` の明示ONだけになる。
- [ ] 自動撮影投稿OFF時に専用Webhook URLが残っていても送信されない。
- [ ] `output.uploadDiscord=true` かつ `autoCapture.discord.enabled=false` の回帰テストを追加する。
- [ ] staleな `autoCapture.discord.webhookUrl` がOFF時に参照されないことをテストする。
- [ ] UI文言と無効化条件が新しい仕様に一致する。
- [ ] 既存の通常Discord投稿機能は手動処理で維持される。
- [ ] Goテストとフロントエンドテンプレート検査が通る。
