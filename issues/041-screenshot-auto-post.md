# 041 Win+Shift+S スクリーンショット自動投稿

## 状態

完了

## 問題

Win+Shift+S で保存された画像は `%USERPROFILE%\Pictures\Screenshots` に保存されるが、現在はVRChat写真自動投稿の対象外で、手動でドラッグ&ドロップする必要がある。

## 期待する挙動

設定でONにすると、Screenshotsフォルダに新しく保存された画像を検知し、既存の画像処理と同じように縮小・Discord投稿・履歴追加を行う。初期値はOFFにする。

## 受け入れ条件

- 設定画面でスクリーンショット自動投稿をON/OFFできる。
- 初期値はOFF。
- 監視対象フォルダの初期値は `%USERPROFILE%\Pictures\Screenshots`。
- 自動投稿のスキャン間隔を設定画面で変更できる。
- 通常のDiscord Webhook URLへ投稿できる。
- 既存のVRChat写真自動投稿と併用できる。
- 保存後に監視設定が反映される。

## 対応内容

- `screenshotAutoPost` 設定を追加した。
- 初期Screenshotsフォルダを `%USERPROFILE%\Pictures\Screenshots` にした。
- 自動投稿のスキャン間隔を設定値として追加した。
- VRChat写真自動投稿とスクリーンショット自動投稿を同時に起動できるようにした。
- 設定画面、README、仕様書を更新した。
