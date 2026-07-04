# Discord投稿でallowed_mentionsを無効化する

## 問題

Discord投稿本文には、QRコードURL、ファイル名、同席ユーザー表示名など、ユーザー入力や外部データ由来の文字列が含まれる可能性がある。
現行の `UploadDiscordWithContent()` は `content` だけをpayloadに入れており、`allowed_mentions` を指定していない。
表示名などに `@everyone`、`@here`、role/user mention相当の文字列が含まれると、意図しないメンションにつながる可能性がある。

## 期待する挙動

Discord Webhook投稿では、既定で全メンションを無効化する。
ユーザーが投稿本文に含めた文字列はテキストとして表示され、Discord通知を発生させない。

## 受け入れ条件

- [x] Discord投稿payloadに `allowed_mentions: { "parse": [] }` を含める。
- [x] 通常画像投稿、VRChat写真自動投稿、スクリーンショット自動投稿、自動撮影Discord投稿の全経路で適用される。
- [x] payloadの単体テストで `allowed_mentions` が入ることを確認する。
- [x] 既存の画像添付、QRコードURL本文、同席ユーザー本文を壊さない。
- [x] READMEまたは仕様に、メンションを無効化して投稿することを記載する。

## 対応内容

- Discord Webhook payload生成を共通化し、全投稿経路で `allowed_mentions.parse=[]` を含めるようにした。
- READMEに、Discord投稿時はメンションを無効化することを追記した。
