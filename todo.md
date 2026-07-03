# rc20 frontend起動エラー修正チェックリスト

- [x] rc20ログと `avatar is not defined` の発生箇所を確認する
- [x] Vue template内の説明文バッククォートを除去する
- [x] template内バッククォート混入を検出するCI/Release検査を追加する
- [x] frontend build、静的検査、Goテストを実行する
