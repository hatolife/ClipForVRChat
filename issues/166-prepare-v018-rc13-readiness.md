# v0.1.8-rc13作成可能状態まで未完了項目を解消する

## 問題

[#164](164-audit-v018-completed-items.md) の再監査で、v0.1.8向けissueの一部が完了扱いにできない状態だと判明している。
このままでは `v0.1.8-rc13` を作成しても、Stream Camera映像保存、player_local構図、埋め込みメタデータ、検証導線に不完全な部分が残る。

## 期待する挙動

未完了issueを実装・検証し、`v0.1.8-rc13` を作成できる状態にする。
作業は可能な範囲で安価なモデルのサブエージェントに分担し、親エージェントが完了判定と不足分の再チケット化を行う。

## 受け入れ条件

- [ ] [#121](121-stream-camera-local-view-and-error-ux.md)、[#129](129-add-spout-capture-helper.md) から [#135](135-add-spout-stream-camera-verification-guide.md) が完了扱いにできる。
- [ ] [#138](138-define-player-local-coordinate-spec.md)、[#140](140-implement-player-local-pose-transform.md)、[#141](141-integrate-player-local-coordinate-ui.md)、[#142](142-verify-player-local-camera-compositions.md) が完了扱いにできる。
- [ ] [#143](143-define-autocapture-embedded-metadata-schema.md)、[#144](144-implement-image-metadata-writer.md)、[#145](145-integrate-exif-presence-metadata.md)、[#147](147-verify-embedded-metadata-output.md) が完了扱いにできる。
- [ ] [#156](156-add-wails-api-surface-check.md)、[#157](157-sync-v018-autocapture-docs-and-specs.md)、[#158](158-populate-or-remove-autocapture-world-instance-metadata.md)、[#163](163-show-autocapture-test-results-in-settings.md) が完了扱いにできる。
- [ ] Goテスト、フロントエンドビルド、追加した検証スクリプトが成功する。
- [ ] 不完全な部分が残った場合は、完了扱いにせず新規または既存issueへ問題・期待する挙動・受け入れ条件を記録する。
- [ ] RC13作成前の差分が署名付きコミットで整理されている。

## 進行中メモ

- 2026-07-01: サブエージェントが Stream Camera(Spout) 経路の helper / Go / CI・Release 検証を担当開始。
- 2026-07-01: 対象は `tools/spout-capture/main.cpp`、`src/internal/appcore/spout.go`、`.github/workflows/ci.yml`、`.github/workflows/release.yml` の Spout 関連差分。
