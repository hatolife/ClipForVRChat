# 診断データからWebhook URLとDiscord tokenを除外する

## 問題

確認用zipと暗号化対象zipに、`config.json` のWebhook URLや `history.json` の `discordToken` が含まれ得る。

## 期待する挙動

不具合報告用データには確認用zipを残しつつ、Webhook URLやDiscord削除用tokenなどの秘密情報はダミー値へ置換される。

## 受け入れ条件

- 診断データ内の `config.json` にWebhook URLの生値が含まれない。
- 診断データ内の `history.json` に `discordToken` の生値が含まれない。
- 確認用zipと `.zip.gpg` の復号後zipで同じ安全化済みデータを確認できる。
- Webhook URLやtokenが設定済みかどうかを診断データから判別できなくてもよい。

## 対応内容

- 診断テキストの共通redaction処理を追加した。
- Webhook URL、`discordToken`、`webhookUrl` の生値を診断データへ残さないようにした。
- 確認用zipに入るJSON/logにも同じredactionを適用した。
