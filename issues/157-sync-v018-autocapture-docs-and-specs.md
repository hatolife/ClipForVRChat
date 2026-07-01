# v0.1.8自動撮影とRelease成果物仕様をREADME/SPECへ反映する

## 問題

READMEと `src/SETTINGS_SPEC.md` は、通常画像処理、VRChat写真自動処理、スクリーンショット自動処理の説明が中心で、自動撮影タブとv0.1.8の制約を十分に説明していない。
また、Release workflowの添付/同梱物やビルドメタデータと仕様文書の差分も残っている。

## 期待する挙動

ユーザー向けREADME、設定画面仕様、アプリ仕様、Release仕様が、v0.1.8で実際に提供する自動撮影機能と制約に一致する。
未実装のSpout、EXIF、player_local、multiなどは、実装済みのように読めない状態にする。

## 受け入れ条件

- READMEに自動撮影タブ、Stream/Spout前提、Photo方式フォールバック、同席ユーザー情報、EXIF未実装/実装状態を記載する。
- `src/SETTINGS_SPEC.md` に自動撮影カテゴリと設定項目を追加する。
- `src/SPEC.md` のv0.1.8実装済み/未実装/将来対応の境界を更新する。
- Release成果物一覧がworkflowの実際の添付/同梱物と一致する。
- Spout helperを追加する場合は、README/SPEC/Release notesに同梱物と前提を追記する。
- ドキュメントが未実装機能を実装済みとして案内していないことを確認する。

## 対応内容

- v0.1.8実装に合わせて対応済み。詳細は `feat/v0.1.8-resolve-issues` の実装、README、SPEC、RELEASE_NOTESを参照。

## 再監査メモ

- 2026-07-01: [#164](164-audit-v018-completed-items.md) の再監査で未達が見つかったため、完了扱いを取り消して `要対応` に戻した。
