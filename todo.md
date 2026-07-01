# v0.1.8-rc13準備チェックリスト

- [ ] Stream Camera(Spout)経路を完了扱いにできる状態へ修正する
  - [ ] helperのsender選択/エラーJSON/バージョン表示を完成させる
  - [ ] Go側のSpout画像検証、ログ、sidecar/Discord紐づけを確認する
  - [ ] Stream設定UIとCI/Release同梱検証をSpout前提へ揃える
- [ ] player_local構図を完了扱いにできる状態へ修正する
  - [ ] 初期構図をplayer_localへ変更し、妥当な初期値へ更新する
  - [ ] world/player_localの順変換・逆変換と保存/追加/移動APIを確認する
  - [ ] 構図ごとのテスト撮影結果表示を確認する
- [ ] 埋め込みメタデータを完了扱いにできる状態へ修正する
  - [ ] metadata schema、PNG/JPEG writer/reader、サイズ超過処理を完成させる
  - [ ] 保存処理でmetadata失敗を警告扱いにし、sidecar/Discord/履歴と整合させる
  - [ ] 読み戻しテストと検証手順を追加する
- [ ] 周辺検証とドキュメントをRC13向けに揃える
  - [ ] Wails API同期チェックをCIへ追加する
  - [ ] world/instance metadata取得テストを追加する
  - [ ] README/SPEC/SETTINGS_SPEC/RELEASE_NOTES/検証手順を同期する
  - [ ] 未完了が残る場合はissue化し、完了扱いにしない
