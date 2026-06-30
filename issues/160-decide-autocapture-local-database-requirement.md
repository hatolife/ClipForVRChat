# 自動撮影ローカルDB要件を実装または仕様から外す

## 問題

Codex用仕様書では、自動撮影の構図、Batch、Shot、同席ユーザー情報をSQLiteへ保存する設計になっている。
現行実装は `config.json`、sidecar JSON、既存 `history.json` で保持しており、SQLiteやローカルDB層は存在しない。
仕様と実装の差分が残ると、保持・検索・削除・復旧の期待値が曖昧になる。

## 期待する挙動

v0.1.8でSQLiteを実装するか、sidecar JSON + history JSONを正式な保持方式として仕様からDB要件を外すかを決める。
決定に従って、画像とメタデータの紐づけ、履歴、削除、診断、将来移行方針を整理する。

## 受け入れ条件

- v0.1.8でSQLiteを使うかどうかを明文化する。
- SQLiteを使う場合は、Batch/Shot/User/Viewのschema、マイグレーション、バックアップ方針を実装issueへ分割する。
- SQLiteを使わない場合は、sidecar JSONとhistory JSONを正本/索引として扱う仕様へREADME/SPEC/todoを更新する。
- sidecar JSONとhistory JSONだけで、最小要件の「画像と同席ユーザー情報の紐づけ」が満たせることを確認する。
- 将来SQLiteへ移行する場合の互換方針を記録する。

## 対応内容

- v0.1.8実装に合わせて対応済み。詳細は `feat/v0.1.8-resolve-issues` の実装、README、SPEC、RELEASE_NOTESを参照。
