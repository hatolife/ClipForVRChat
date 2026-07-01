# v0.1.8-rc13 CI修正チェックリスト

- [ ] Spout helperのWindows CIビルド失敗を修正する
  - [x] `SpoutLibrary_static.lib` リンク失敗の原因をissueへ追記する
  - [x] CMakeを上流Spout2の実ターゲットに合わせる
  - [x] `SpoutLibrary.dll` をCI/Release同梱・検証対象にする
  - [ ] GitHub ActionsのCI/Release成功を確認する
