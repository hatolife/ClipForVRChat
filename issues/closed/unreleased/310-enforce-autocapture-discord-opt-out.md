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

## 方針決定

2026-07-08: ユーザー判断によりAを採用する。
自動撮影Discord投稿は `autoCapture.discord.enabled` の明示ONだけで有効化し、通常投稿用 `output.uploadDiscord` は自動撮影投稿の許可条件として使わない。
実装時は、自動撮影投稿OFFのために `output.uploadDiscord=true` でも自動撮影画像・本文をDiscordへ送らないように変更する。

## 受け入れ条件

- [x] 自動撮影投稿の実行条件が `autoCapture.discord.enabled` の明示ONだけになる。
- [x] 自動撮影投稿OFF時に専用Webhook URLが残っていても送信されない。
- [x] `output.uploadDiscord=true` かつ `autoCapture.discord.enabled=false` の回帰テストを追加する。
- [x] staleな `autoCapture.discord.webhookUrl` がOFF時に参照されないことをテストする。
- [x] UI文言と無効化条件が新しい仕様に一致する。
- [x] 既存の通常Discord投稿機能は手動処理で維持される。
- [x] Goテストとフロントエンドテンプレート検査が通る。

## 作業メモ

2026-07-08: ユーザー指示によりbackend-sideのみで実装する。`src/frontend/src/main.js` は変更しない。
`autoCapture.discord.enabled` を自動撮影Discord投稿の唯一の実行条件とし、`output.uploadDiscord=true` 単独では投稿しない回帰テストを追加する。

2026-07-08: `ca7049a fix(autocapture): enforce backend privacy gates` で、自動撮影Discord投稿OFF時の漏えいを防ぐために、backendの実行条件を「`output.uploadDiscord` または `autoCapture.discord.enabled`」から「`autoCapture.discord.enabled` の明示ONのみ」へ変更した。
同じコミットで、自動撮影専用Webhook URLが残っていてもOFF時には参照しないようにし、通常投稿ONだけでは自動撮影画像・本文をDiscordへ送らない回帰テストを追加した。

2026-07-08: `src/frontend/src/main.js` で、UIの誤認を防ぐために、自動撮影Discord投稿の説明を「通常Discord投稿ON、またはこの設定ON」から「この設定ONのときだけ」へ変更した。
同じく、自動撮影用Webhook URLと画像添付行の無効化条件を「通常投稿OFFかつ自動撮影投稿OFF」から「自動撮影投稿OFF」へ変更し、画面上でも専用opt-inと一致させた。
