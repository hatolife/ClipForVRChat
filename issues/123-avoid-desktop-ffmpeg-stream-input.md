# 自動撮影Stream方式でデスクトップ全体を撮らない

## 問題

v0.1.8-rc10でも既存configの `ffmpeg入力引数` が `-f gdigrab -framerate 30 -i desktop` のままだと、Stream方式の出力がデスクトップ全体のスクリーンショットになる。
この入力は確認用の暫定値であり、自動撮影の初期挙動としては不適切。

## 期待する挙動

- 初期設定ではデスクトップ全体ではなくVRChatウィンドウを取得対象にする。
- rc9/rc10で保存済みのデスクトップ全体取得設定は、自動的にVRChatウィンドウ取得へ移行する。
- デスクトップ全体取得を使う場合は、ユーザーが明示的に入力引数を変更した場合だけにする。

## 受け入れ条件

- `ffmpeg入力引数` の初期値が `desktop` ではない。
- 既存configの旧初期値 `-f gdigrab -framerate 30 -i desktop` はNormalize時にVRChatウィンドウ取得へ移行される。
- 設定画面の説明が「デスクトップ全体取得」を初期値として案内しない。
- 既存のStream撮影、ffmpeg確認、Photo方式を壊さない。

## 対応内容

- ffmpeg入力引数の初期値を `-f gdigrab -framerate 30 -i title=VRChat` に変更した。
- rc9/rc10の旧初期値 `-f gdigrab -framerate 30 -i desktop` は設定読み込み時に新初期値へ移行するようにした。
- 設定画面の説明とplaceholderから、デスクトップ全体取得を初期値として案内する文言を削除した。
