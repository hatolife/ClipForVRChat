# v0.1.8-rc13準備チェックリスト

- [x] Stream Camera(Spout)経路を完了扱いにできる状態へ修正する
  - [x] helperのsender選択/エラーJSON/バージョン表示を完成させる
  - [x] Go側のSpout画像検証、ログ、sidecar/Discord紐づけを確認する
  - [x] Stream設定UIとCI/Release同梱検証をSpout前提へ揃える
- [x] player_local構図を完了扱いにできる状態へ修正する
  - [x] 初期構図をplayer_localへ変更し、妥当な初期値へ更新する
  - [x] world/player_localの順変換・逆変換と保存/追加/移動APIを確認する
  - [x] 保存済みworld poseとplayer_local設定の区別をsidecar/診断ログへ反映する
- [x] 埋め込みメタデータを完了扱いにできる状態へ修正する
  - [x] metadata schema、PNG/JPEG writer/reader、サイズ超過処理を完成させる
  - [x] 保存処理でmetadata失敗を警告扱いにし、sidecar/Discord/履歴と整合させる
  - [x] 読み戻しテストと検証手順を追加する
- [x] 周辺検証とドキュメントをRC13向けに揃える
  - [x] Wails API同期チェックをCIへ追加する
  - [x] world/instance metadata取得テストを追加する
  - [x] README/SPEC/SETTINGS_SPEC/RELEASE_NOTES/検証手順を同期する
  - [x] 実機/Windows CIが必要な項目は `要確認` とし、完了扱いにしない
