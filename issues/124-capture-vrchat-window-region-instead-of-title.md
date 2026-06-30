# 自動撮影Stream方式で白画像になるtitle取得を避ける

## 問題

v0.1.8-rc11では `ffmpeg -f gdigrab -i title=VRChat` でVRChatウィンドウを取得しているが、実機では白画像として保存される。
VRChatのGPU描画をウィンドウtitle指定で取得できず、ウィンドウ面だけが白く取得されている可能性が高い。

## 期待する挙動

- `title=VRChat` 直接取得ではなく、VRChatウィンドウの画面座標を取得する。
- ffmpegではデスクトップ全体ではなく、VRChatウィンドウ範囲だけを `gdigrab` の `offset_x` / `offset_y` / `video_size` で取得する。
- VRChatウィンドウが見つからない場合は、白画像を成功扱いで投稿せず、分かるエラーを表示する。

## 受け入れ条件

- 初期のffmpeg入力引数が `title=VRChat` 直接取得ではない。
- 旧 `title=VRChat` 初期値は新しいウィンドウ範囲取得へ自動移行される。
- 実行時に `{window_x}` / `{window_y}` / `{window_width}` / `{window_height}` をVRChatウィンドウ座標で置換する。
- ウィンドウ座標取得失敗時は撮影失敗として扱い、履歴とDiscord投稿へ成功画像として残さない。
- 既存のPhoto方式、ffmpeg確認、sidecar JSON、ユーザー情報保持を壊さない。

## 対応内容

- ffmpeg入力引数の初期値を `-f gdigrab -framerate 30 -offset_x {window_x} -offset_y {window_y} -video_size {window_width}x{window_height} -i desktop` に変更した。
- 実行時にPowerShellでVRChatプロセスのメインウィンドウ矩形を取得し、`{window_x}` / `{window_y}` / `{window_width}` / `{window_height}` を置換するようにした。
- 旧 `title=VRChat` 初期値は新しいウィンドウ範囲切り出しへ自動移行するようにした。
