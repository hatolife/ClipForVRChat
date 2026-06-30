# DiscordユーザーID出力をsidecar JSON設定から独立させる

## 問題

UIには「サイドカーJSONにユーザーIDを含める」と「DiscordにユーザーIDを含める」が別トグルとして表示される。
しかし `RunOnce()` では `IncludeUserIDsInSidecar=false` の時点でユーザーIDを消した一覧を作り、その一覧をDiscord本文にも使っている。
そのためDiscord側だけユーザーIDを出す設定が実際には効かない。

## 期待する挙動

sidecar JSON、Discord本文、EXIF/埋め込みメタデータのユーザーID出力は、それぞれの設定に従って独立して制御される。

## 受け入れ条件

- Presenceの元snapshotはユーザーIDを保持したまま扱い、出力先ごとにマスクしたコピーを作る。
- `IncludeUserIDsInSidecar=false` でも `IncludeUserIDsInDiscord=true` ならDiscord本文にはユーザーIDが入る。
- `IncludeUserIDsInDiscord=false` ならsidecar設定に関係なくDiscord本文にはユーザーIDが入らない。
- EXIF/埋め込みメタデータ実装時も、専用設定に従ってユーザーIDを制御できる構造にする。
- 出力先ごとのユーザーID有無を単体テストで確認する。
- UI文言で、各トグルがどの出力先に効くかを明確にする。
