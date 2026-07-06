# 終了時にSpout helper展開キャッシュを削除する

## 問題

単一exe版で内蔵Spout helperを利用すると、`%LOCALAPPDATA%\ClipForVRChat\spout-helper` 配下に `spout-capture.exe` と `SpoutLibrary.dll` の展開データが残る。

## 期待する挙動

ClipForVRChat終了時に、アプリが管理しているSpout helper展開キャッシュを削除する。

明示設定された外部helperや分離版zipの同梱helperは削除しない。

## 受け入れ条件

- [x] Wailsの終了処理でSpout helperキャッシュ削除が実行される。
- [x] 削除対象は `ClipForVRChat/spout-helper` の管理ルートに限定される。
- [x] 削除失敗時は終了を妨げず、診断ログに記録される。
- [x] Go testが通る。
