# OpenURLで開けるURLを許可ホストに制限する

## 問題

Wailsブリッジの `OpenURL` が受け取ったURLをそのままOSブラウザへ渡す。現在の呼び出し元は定数中心だが、将来任意URLが混入した場合のリスクがある。

## 期待する挙動

アプリから開けるURLは、GitHub、BOOTH、Discord公式ヘルプ、作者Twitterなどの信頼済みHTTPS URLに限定される。

## 受け入れ条件

- `https` 以外のURLは拒否される。
- 未許可ホストのURLは拒否される。
- 既存のGitHub、BOOTH、Discordヘルプ、Twitterリンクは開ける。
- Go側にテストがある。

## 対応内容

- `OpenURL` で `https` と許可ホストを検証するようにした。
- 許可ホストを GitHub、BOOTH、Discord公式ヘルプ、X/Twitter に限定した。
- 許可/拒否URLのGoテストを追加した。
