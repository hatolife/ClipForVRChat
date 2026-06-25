# 041 Win+Shift+S スクリーンショット自動処理

## 状態

完了

## 問題

Win+Shift+S で保存された画像は `%USERPROFILE%\Pictures\Screenshots` に保存されるが、現在はVRChat写真自動処理の対象外で、手動でドラッグ&ドロップする必要がある。

## 期待する挙動

設定でONにすると、Screenshotsフォルダに新しく保存された画像を検知し、既存の画像処理と同じように縮小・Discord投稿・履歴追加を行う。初期値はOFFにする。

## 受け入れ条件

- 設定画面でスクリーンショット自動処理をON/OFFできる。
- 初期値はOFF。
- 監視対象フォルダの初期値は `%USERPROFILE%\Pictures\Screenshots`。
- スクリーンショット自動処理用のスキャン間隔を設定画面で変更できる。
- スクリーンショット用Webhook URLを設定画面で変更できる。
- スクリーンショット用Webhook URLが空の場合は通常投稿用Webhook URLへ投稿できる。
- 既存のVRChat写真自動処理と併用できる。
- 保存後に監視設定が反映される。

## 対応内容

- `screenshotAutoPost` 設定を追加した。
- 初期Screenshotsフォルダを `%USERPROFILE%\Pictures\Screenshots` にした。
- スクリーンショット自動処理用のスキャン間隔とWebhook URLを設定値として追加した。
- VRChat写真自動処理とスクリーンショット自動処理を同時に起動できるようにした。
- 設定画面、README、仕様書を更新した。
