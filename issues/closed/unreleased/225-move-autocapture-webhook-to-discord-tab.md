# 自動撮影用Webhook URLをDiscord投稿タブへ移動する

## 問題

自動撮影用Webhook URLが自動撮影の詳細設定内にあり、Discord投稿先の設定が複数タブに分散している。

## 期待する挙動

自動撮影用Webhook URLを、通常投稿用Webhook URL、VRChat写真用Webhook URL、スクリーンショット用Webhook URLと同じDiscord投稿タブで編集できる。

## 受け入れ条件

- [x] 自動撮影用Webhook URLの入力欄がDiscord投稿タブに表示される。
- [x] 自動撮影詳細画面から同入力欄が外れる。
- [x] 既存のフォールバック説明とdisabled条件が維持される。
- [x] frontend template/API検査が通る。
