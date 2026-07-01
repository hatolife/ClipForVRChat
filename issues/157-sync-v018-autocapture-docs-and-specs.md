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

## 実装前調査メモ

実装方針:

- このissueは、実装系issueが完了した後のドキュメント同期ゲートとして扱う。
- README/SPEC/Release Notesは、実装前に先行して「できる」と書かない。実装後の実際の制約に合わせる。
- v0.1.8で提供する `player_local` は手動基準Pose方式であり、プレイヤーroot自動追従ではないことを明記する。
- Discord添付画像内metadataは保持保証がないため、sidecar JSONとDiscord本文を正本として説明する。

対象ファイル:

- `README.md`
- `src/SPEC.md`
- `src/SETTINGS_SPEC.md`
- `RELEASE_NOTES.md`
- `docs/v0.1.8-stream-spout-verification.md`
- `docs/v0.1.8-player-local-verification.md`
- `docs/v0.1.8-embedded-metadata-verification.md`

小タスク:

- 自動撮影タブの全設定項目を `src/SETTINGS_SPEC.md` に追加する。特に `startDelayMs`、Spout sender、画像埋め込みmetadata、Discord画像添付を含める。
- READMEの自動撮影説明を、Stream/Spout、Photoフォールバック、手動player_local、sidecar/Discord/埋め込みmetadataの実装済み範囲に合わせる。
- Release成果物一覧に `spout-capture.exe`、`Spout2-LICENSE.txt`、`Release-signing-public-key.url` を含める。
- 既知制限に、標準OSCではplayer root自動取得しないこと、Discord添付metadata保持は保証しないことを入れる。
- `rg "ffmpeg入力設定|未実装|ダミー|Spout.*未"` で不整合文言を確認する。

確認方法:

- README/SPEC/Release Notesが [#129](129-add-spout-capture-helper.md) から [#163](163-show-autocapture-test-results-in-settings.md) の実装済み範囲と一致する。
- Release workflowの実際のzip内容とドキュメントの成果物一覧が一致する。
