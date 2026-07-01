# 自動撮影テスト結果を設定画面に表示する

## 問題

構図カードの「テスト撮影」は `TestAutoCaptureView()` を呼び、backend側では成功結果をstate/historyへ追加している。
しかしfrontend側は戻り値の最初のエラーと短いtoastだけを見ており、成功した画像、保存先、Discord投稿結果、sidecar作成結果を設定画面上で確認しづらい。
テスト撮影ボタンは視点確認用の重要操作なので、結果が見えないと実機確認が難しい。

## 期待する挙動

テスト撮影後、設定画面上で対象構図の撮影結果を確認できる。
成功時は保存された画像やDiscord投稿結果へ辿れ、失敗時は原因を構図単位で確認できる。

## 受け入れ条件

- `TestAutoCaptureView()` の戻り値をfrontend stateへ反映し、成功結果を画面に表示する。
- 設定画面を閉じなくても、テスト撮影した画像のサムネイル、保存先、エラーを確認できる。
- 複数回テストした場合に、どの構図の結果か分かる。
- 成功時と失敗時でtoastだけに依存しない表示にする。
- 履歴へ追加する/しない方針を明確にし、通常の自動撮影履歴と混同しない。
- 既存の結果画面、履歴追加、Discord投稿処理を壊さない。

## 対応内容

- v0.1.8実装に合わせて対応済み。詳細は `feat/v0.1.8-resolve-issues` の実装、README、SPEC、RELEASE_NOTESを参照。

## 再監査メモ

- 2026-07-01: [#164](164-audit-v018-completed-items.md) の再監査で未達が見つかったため、完了扱いを取り消して `要対応` に戻した。

## 実装前調査メモ

実装方針:

- `TestAutoCaptureView()` は既に `[]appcore.Result` を返しており、結果には `outputPath`、`url`、`thumbnail`、`error` を入れられる。新しい保存経路は作らず、frontendが戻り値を保持して表示する。
- 自動撮影側の `finalizeAutoCaptureImage()` では現状thumbnailを生成していないため、設定画面上のサムネイル表示を必須にするなら、Stream/Photo保存後に `DecodeImageFile` + `ThumbnailDataURL` を実行して `Result.Thumbnail` を入れる必要がある。
- テスト撮影は通常履歴にも入る現行仕様だが、設定画面では「テスト撮影結果」として構図カード内に別表示し、通常スケジュール結果と混同しない文言にする。

対象ファイル:

- `src/frontend/src/main.js`
- `src/frontend/src/style.css`
- `src/app.go`
- `src/internal/appcore/autocapture.go`

小タスク:

- `autoCaptureTestResults[view.id]` を `{ ok, message, results, updatedAt }` 形式にする。
- 成功時は保存先 `outputPath`、Discord URL、sidecar path候補 `outputPath + ".json"` を表示する。
- `outputPath` がある場合は「保存先で表示」ボタンを追加し、既存 `RevealFileInExplorer` を使う。
- `thumbnail` がある場合はカード内に小さく表示する。ない場合も保存先表示で確認できるようにする。
- 失敗時は最初のエラーだけでなく、全resultの構図名/エラーを表示する。
- `TestAutoCaptureView()` 結果で `Thumbnail` を生成するか、frontendはthumbnailなしでも成立するUIにする。

確認方法:

- テスト撮影成功後、設定画面を閉じずに保存先とDiscord URLを確認できる。
- テスト撮影失敗後、該当構図カードにエラー全文が残る。
- 複数構図で連続テストしても、各カードの結果が混ざらない。
