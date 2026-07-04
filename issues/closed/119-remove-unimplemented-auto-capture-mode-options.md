# 未実装の自動撮影方式を設定画面から外す

## 問題

Spout2直接受信やCamera Dolly Multiなど、まだ実装していない方式を設定画面から選べると、実際に使える機能と誤解される。

## 期待する挙動

設定画面には実際に動作するStream(ffmpeg)方式とPhoto方式だけを表示し、未実装のSpout2直接受信やCamera Dolly Multi設定をダミーとして出さない。

## 受け入れ条件

- 自動撮影タブでStream(ffmpeg)方式とPhoto方式だけを選択できる。
- Spout2直接受信やCamera Dolly Multiを選択できない。
- 既存configに未対応方式が残っていても正規化でStreamへ戻る。

## 対応内容

- v0.1.8実装に合わせて対応済み。Stream方式はSpout主経路へ置き換え、ffmpeg主導線は廃止。詳細はREADME/SPEC/RELEASE_NOTESを参照。
