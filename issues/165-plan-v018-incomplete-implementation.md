# v0.1.8未完了項目の実装前調査

## 問題

[#164](164-audit-v018-completed-items.md) の再監査で、v0.1.8の完了扱いから `要対応` に戻した項目が複数ある。
このまま実装へ入ると、外部仕様、既存コードの接続点、テスト方針が曖昧なまま修正範囲が広がる可能性がある。

## 期待する挙動

未完了項目を実装できる状態にするため、各issueへ実装方針、対象ファイル、小タスク、検証方法を追記する。

## 受け入れ条件

- [#121](121-stream-camera-local-view-and-error-ux.md)、[#129](129-add-spout-capture-helper.md) から [#135](135-add-spout-stream-camera-verification-guide.md)、[#138](138-define-player-local-coordinate-spec.md)、[#140](140-implement-player-local-pose-transform.md) から [#145](145-integrate-exif-presence-metadata.md)、[#147](147-verify-embedded-metadata-output.md)、[#156](156-add-wails-api-surface-check.md)、[#157](157-sync-v018-autocapture-docs-and-specs.md)、[#158](158-populate-or-remove-autocapture-world-instance-metadata.md)、[#163](163-show-autocapture-test-results-in-settings.md) を確認する。
- 必要な外部仕様や既存コードの制約を確認する。
- 各issueに実装方針、対象ファイル、小タスク、テスト/確認方法を追記する。
- この調査では実装コード変更を行わない。

## 調査メモ

参照した外部仕様:

- VRChat OSC Overview: デフォルトOSC portは入力9000、出力9001。`--osc=inPort:senderIP:outPort` で変更可能。
- VRChat Camera OSC endpoints: `/usercamera/Mode` は 0=Off、1=Photo、2=Stream。`/usercamera/Pose`、`/usercamera/Capture`、`/usercamera/Close`、`/usercamera/Streaming`、Zoom/Exposure等が提供される。
- Spout2: Windows向けのtexture/frame sharing SDK。Spout2本体はBSD-2-Clauseで、SpoutLibraryを使う方針はGPLのSpoutRecorder流用を避けられる。
- PNG eXIf extension: PNGの `eXIf` chunkはExif profileを保持し、JPEG APP1 marker/length/`Exif\0\0` を含めない。IDAT列の途中以外に1つだけ置ける。
- Discord API: webhookのファイル添付は `multipart/form-data` の `files[n]` と `payload_json` を使う。添付画像内metadataの保持はAPI仕様として保証されていない。

実装順序案:

1. Spout helperとSpout設定/検証を完了する: [#129](129-add-spout-capture-helper.md), [#130](130-add-spout-sender-settings-and-diagnostics.md), [#131](131-integrate-spout-helper-into-auto-capture.md), [#132](132-validate-spout-capture-output-and-metadata.md), [#133](133-update-auto-capture-stream-ui-and-docs-for-spout.md), [#134](134-package-spout-helper-in-ci-release.md)。
2. `player_local` の仕様と保存/逆変換/UIを完了する: [#138](138-define-player-local-coordinate-spec.md), [#140](140-implement-player-local-pose-transform.md), [#141](141-integrate-player-local-coordinate-ui.md)。
3. metadata writer/readerと保存処理を完了する: [#143](143-define-autocapture-embedded-metadata-schema.md), [#144](144-implement-image-metadata-writer.md), [#145](145-integrate-exif-presence-metadata.md), [#147](147-verify-embedded-metadata-output.md)。
4. 検証支援と周辺品質を完了する: [#156](156-add-wails-api-surface-check.md), [#158](158-populate-or-remove-autocapture-world-instance-metadata.md), [#163](163-show-autocapture-test-results-in-settings.md)。
5. ドキュメントと実機確認手順を同期する: [#135](135-add-spout-stream-camera-verification-guide.md), [#142](142-verify-player-local-camera-compositions.md), [#157](157-sync-v018-autocapture-docs-and-specs.md)。

実装時の注意:

- `player_local` は手動基準Pose方式であり、標準OSCだけでプレイヤーrootへ自動追従する実装として扱わない。
- 埋め込みmetadataの失敗は、画像保存全体の失敗にしない。警告としてsidecar/ログ/UIへ残す。
- Discord添付画像のmetadata保持は保証しない。Discord本文とsidecar JSONを正本にする。
- Release zip検証は `dist/` だけでなく展開済みzip内も検証する。

## 対応内容

- 2026-07-01: 対象issueへ実装方針、対象ファイル、小タスク、確認方法を追記した。実装コード変更は行っていない。
