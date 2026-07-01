# v0.1.8-rc13作成可能状態まで未完了項目を解消する

## 問題

[#164](164-audit-v018-completed-items.md) の再監査で、v0.1.8向けissueの一部が完了扱いにできない状態だと判明している。
このままでは `v0.1.8-rc13` を作成しても、Stream Camera映像保存、player_local構図、埋め込みメタデータ、検証導線に不完全な部分が残る。

## 期待する挙動

未完了issueを実装・検証し、`v0.1.8-rc13` を作成できる状態にする。
作業は可能な範囲で安価なモデルのサブエージェントに分担し、親エージェントが完了判定と不足分の再チケット化を行う。

## 受け入れ条件

- [x] [#121](121-stream-camera-local-view-and-error-ux.md)、[#129](129-add-spout-capture-helper.md) から [#135](135-add-spout-stream-camera-verification-guide.md) が完了扱いまたは実機確認待ちにできる。
- [x] [#138](138-define-player-local-coordinate-spec.md)、[#140](140-implement-player-local-pose-transform.md)、[#141](141-integrate-player-local-coordinate-ui.md)、[#142](142-verify-player-local-camera-compositions.md) が完了扱いまたは実機確認待ちにできる。
- [x] [#143](143-define-autocapture-embedded-metadata-schema.md)、[#144](144-implement-image-metadata-writer.md)、[#145](145-integrate-exif-presence-metadata.md)、[#147](147-verify-embedded-metadata-output.md) が完了扱いにできる。
- [x] [#156](156-add-wails-api-surface-check.md)、[#157](157-sync-v018-autocapture-docs-and-specs.md)、[#158](158-populate-or-remove-autocapture-world-instance-metadata.md)、[#163](163-show-autocapture-test-results-in-settings.md) が完了扱いにできる。
- [x] Goテスト、フロントエンドビルド、追加した検証スクリプトが成功する。
- [x] 不完全な部分が残った場合は、完了扱いにせず新規または既存issueへ問題・期待する挙動・受け入れ条件を記録する。
- [x] RC13作成前の差分が署名付きコミットで整理されている。

## 対応内容

- 2026-07-01: Spout helperに `--version`、候補付きエラーJSON、sender自動選択を追加した。
- 2026-07-01: Go側Spout連携でhelper version確認、候補付き失敗JSON解析、白/黒/透明フレーム検出、統計ログを追加した。
- 2026-07-01: `player_local` 初期構図、world/player_local逆変換、現在Pose保存/追加、resolved pose sidecar/metadata反映を追加した。
- 2026-07-01: PNG iTXt/eXIf、JPEG EXIF APP1の書き込み/読み戻し、idempotent差し替え、ユーザー数過多時の切り詰めを実装した。
- 2026-07-01: metadata書き込み失敗を警告扱いにし、sidecar JSON、Discord投稿、履歴追加を継続するようにした。
- 2026-07-01: Wails API同期チェック、world/instance token保持テスト、設定画面のテスト撮影結果詳細表示を追加した。
- 2026-07-01: Stream/Spout、player_local、埋め込みmetadataの実機確認手順を `docs/` に追加した。

## 残確認

- [#121](121-stream-camera-local-view-and-error-ux.md)、[#129](129-add-spout-capture-helper.md)、[#130](130-add-spout-sender-settings-and-diagnostics.md)、[#131](131-integrate-spout-helper-into-auto-capture.md)、[#132](132-validate-spout-capture-output-and-metadata.md)、[#134](134-package-spout-helper-in-ci-release.md)、[#141](141-integrate-player-local-coordinate-ui.md) はWindows実機またはGitHub Actionsでの確認が必要なため `要確認` とした。
- 2026-07-01: `v0.1.8-rc13` は修正前コミットでタグ作成・push済みだが、CI/Releaseが `SpoutLibrary_static.lib` リンク失敗で落ちた。`release/v0.1.8` は追加コミット `01db7b0` で修正し、CI run `28494855637` は成功した。公開済みタグを書き換えるには明示許可が必要なため、Release成果物は `v0.1.8-rc13` タグ再作成許可または次RCタグ作成待ち。
- 2026-07-01: ユーザーから `rc13を作成` の明示指示を受けたため、既存の失敗済み `v0.1.8-rc13` タグを削除し、現在の `release/v0.1.8` HEADへ署名付きタグを付け直してRelease workflowを再実行する。

## 進行中メモ

- 2026-07-01: サブエージェントが Stream Camera(Spout) 経路の helper / Go / CI・Release 検証を担当開始。
- 2026-07-01: 対象は `tools/spout-capture/main.cpp`、`src/internal/appcore/spout.go`、`.github/workflows/ci.yml`、`.github/workflows/release.yml` の Spout 関連差分。
