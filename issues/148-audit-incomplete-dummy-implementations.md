# 未実装/ダミー/簡易実装を洗い出す

## 問題

v0.1.8の自動撮影まわりでは、実機確認でデスクトップ撮影、白画像、ワールド座標扱い、EXIF未接続など、仕様上は存在するが完全なプログラムとして成立していない点が見つかっている。
既存コード全体にも、未実装、ダミー実装、仮実装、簡易実装、設定だけ存在して動作しない項目が残っている可能性がある。

## 期待する挙動

既存コードを確認し、完全な実装として扱えない箇所を根拠付きで洗い出す。
見つかった項目は、実装すべきもの、UIから隠すべきもの、将来対応として明示すべきものへ分類し、必要に応じて個別issueを作成する。

## 受け入れ条件

- [x] `TODO`、`FIXME`、`panic("TODO")`、`not implemented`、`dummy`、`placeholder`、仮実装を示す文言を検索する。
- [x] 設定値/API/UIが存在するのに実装へ接続されていない項目を確認する。
- [x] 自動撮影、画像処理、Discord投稿、履歴、診断、Release/CIの主要経路を確認する。
- [x] 見つかった未完成点について、ファイル/関数/現象/期待する修正方針を記録する。
- [x] 対応が必要な項目は日本語issueとして分割し、`issues/README.md` に追加する。
- [x] この調査では実装コード変更を行わない。

## 監査結果

明示的な `TODO` / `FIXME` / `not implemented` は、依存パッケージや生成物を除く実装コードにはほぼ残っていなかった。
一方で、設定型、UI文言、仕様に存在するが実処理へ接続されていない項目が複数見つかった。

既存issueでカバー済みの主な未完成点:

- Stream方式は現状ffmpeg/gdigrabの画面キャプチャであり、Stream Camera/Spout映像そのものではない。対応先: [#129](129-add-spout-capture-helper.md), [#131](131-integrate-spout-helper-into-auto-capture.md), [#132](132-validate-spout-capture-output-and-metadata.md), [#133](133-update-auto-capture-stream-ui-and-docs-for-spout.md), [#135](135-add-spout-stream-camera-verification-guide.md)。
- Spout helperはRelease workflowへ同梱されていない。対応先: [#134](134-package-spout-helper-in-ci-release.md)。
- `CoordinateSpace` は送信前Pose変換に使われず、`player_local` は未実装。対応先: [#138](138-define-player-local-coordinate-spec.md) から [#142](142-verify-player-local-camera-compositions.md)。
- `WriteEXIF` / `WriteUserListToEXIF` は設定だけ存在し、画像埋め込みメタデータwriterに接続されていない。対応先: [#143](143-define-autocapture-embedded-metadata-schema.md) から [#147](147-verify-embedded-metadata-output.md)。
- output log同席ユーザー取得はsnapshot方式であり、継続監視やworld/instance metadataは未完了。対応先: [#111](111-vrc-output-log-presence-users.md)、[#158](158-populate-or-remove-autocapture-world-instance-metadata.md)。

今回新規に分割した未完成点:

- [#149](149-implement-or-remove-autocapture-multi-camera-settings.md) 自動撮影multi/Camera Dolly設定が保存されるが、撮影処理は常にSequential。
- [#150](150-connect-autocapture-scheduler-overlap-controls.md) `skipIfPreviousBatchRunning` と `maxBatches` の挙動が設定どおりではない。
- [#151](151-implement-or-remove-autocapture-discord-post-options.md) `postMode` と `includeImages` がDiscord投稿処理に接続されていない。
- [#152](152-apply-per-view-capture-delay.md) 構図ごとの `captureDelayMs` がログ出力だけで待機処理に使われていない。
- [#153](153-expose-autocapture-output-format-and-filename-template.md) 自動撮影の出力形式/ファイル名テンプレート設定がUIに出ていない。
- [#154](154-decouple-autocapture-discord-user-id-from-sidecar-user-id.md) DiscordユーザーID出力がsidecarユーザーID設定に依存している。
- [#155](155-define-sidecar-json-lifecycle-with-history-delete.md) sidecar JSONの履歴/削除ライフサイクルが未定義。
- [#156](156-add-wails-api-surface-check.md) Wails公開APIとフロント呼び出しの同期をCIで検出できない。
- [#157](157-sync-v018-autocapture-docs-and-specs.md) README/設定仕様とv0.1.8自動撮影機能が同期していない。
- [#158](158-populate-or-remove-autocapture-world-instance-metadata.md) sidecarのworld/instance metadataフィールドが常に空になる。
- [#159](159-disable-discord-allowed-mentions.md) Discord投稿payloadで `allowed_mentions` を無効化していない。
- [#160](160-decide-autocapture-local-database-requirement.md) 自動撮影のSQLite/ローカルDB要件が未実装のまま仕様に残っている。
- [#161](161-implement-or-defer-oscquery-discovery.md) OSCQueryによるポート/endpoint検出が未実装。
- [#162](162-align-autocapture-interval-ui-validation.md) 撮影間隔UIは1秒を許可するが、Normalizeで10秒に丸められる。
- [#163](163-show-autocapture-test-results-in-settings.md) テスト撮影成功結果が設定画面へ十分に表示されない。

この調査では、実装コードの変更は行っていない。
