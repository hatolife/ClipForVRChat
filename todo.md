# rc29後のStream Camera/Spout取得修正チェックリスト

- [x] rc29ログでカメラ未起動時のMode/Streaming送信とsender状態を確認する
- [x] rc29 separated2ログでsenderあり/blank-frame失敗を確認する
- [x] issueを作成する
- [x] StreamingのOSC bool送信にnumeric 1/0互換送信を追加する
- [x] Spout helperのblank-frameエラーへフレーム統計を追加する
- [x] Go側のhelper errorログへフレーム統計を追加する
- [x] テストと静的チェックを実行する
- [ ] 変更をコミットする
