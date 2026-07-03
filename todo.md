# rc26後のSpout sender待機修正チェックリスト

- [x] Spout senderなし即失敗のissueを作成する
- [x] Spout helperでsender出現をtimeout内で待つ
- [x] Stream撮影直前にMode/Streamingを補助再送する
- [x] テストと静的チェックを実行する
- [ ] 変更をコミットしてpushする
- [ ] rcを作成してRelease workflowを確認する
