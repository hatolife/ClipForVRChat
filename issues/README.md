# Issues

このディレクトリは、実装タスクや調査タスクを Markdown で管理するためのものです。

状態はリポジトリ上の実装状況をもとに整理しています。実機確認が必要なものは「要確認」としています。

| No. | Issue | 状態 | 対応バージョン/コミット | 概要 |
| --- | --- | --- | --- | --- |
| 001 | [アイコン品質と複数サイズ生成](001-icon-quality-and-sizes.md) | 完了 | `a800610` | アプリアイコン生成と複数サイズの品質改善。 |
| 002 | [Release zipの内容とバージョン付きファイル名](002-release-zip-contents-and-versioned-name.md) | 完了 | `6263083`, `9a816fe` | 配布zipの内容整理、バージョン情報の埋め込み、Release workflow 整備。 |
| 003 | [メインウィンドウUIと情報表示](003-main-window-ui-and-about.md) | 完了 | `08db9f3`, `bc70012` | メイン画面、情報表示、アプリ情報取得の追加。 |
| 004 | [初回設定フロー](004-initial-config-flow.md) | 完了 | `ecd8c98` 以降 | GUIアプリとして設定を扱える初期フローを実装。 |
| 005 | [ウィンドウへのドラッグ&ドロップ](005-window-drag-and-drop.md) | 完了 | `d2ce5aa`, `ec99a08` | 画像/config.json のドラッグ&ドロップ処理と画面全体のドロップ対応。 |
| 006 | [クリップボード画像の保存ファイル名に時刻を含める](006-clipboard-output-filename.md) | 完了 | `2bf08ca` | クリップボード入力の保存ファイル名を時刻付きに変更。 |
| 007 | [クリップボードのスクリーンショット画像が崩れる](007-clipboard-screenshot-corruption.md) | 完了 | `21cc973` / `v0.1.1` | Win+Shift+S 由来画像の崩れ対策として Windows 登録PNG読み取りを優先。 |
| 008 | [UI仕様の明文化](008-ui-specification.md) | 完了 | `ff917f0` | UI仕様と運用ルールをドキュメント化。 |
| 009 | [ソース配置を src 配下へ移動する](009-move-source-under-src.md) | 完了 | `88961d1` | Wails/Go/frontend ソースを `src` 配下へ整理。 |
| 010 | [設定画面・結果画面・情報画面のUI改善](010-settings-results-about-ui-improvements.md) | 完了 | `8f0e267`, `ec99a08`, `f17eea0` | 設定/結果/情報画面のレイアウト、進捗、クリア、未保存確認を改善。 |
| 011 | [使い方画面とWebhook発行案内](011-help-screen-and-webhook-guide.md) | 完了 | `bd5acf2`, `05d0105` | 使い方画面と Discord Webhook 公式案内リンクを追加。 |
| 012 | [画像URL履歴とDiscord削除確認画面](012-image-history-and-discord-delete-review.md) | 完了 | `709370d` 以降 | 履歴保存、クリア済み表示、Discord削除確認画面、Ctrl/Shift選択、矩形範囲選択を実装。 |
| 013 | [ユーザーフレンドリーなREADMEと設定画面改善](013-user-friendly-readme-and-settings-ui.md) | 完了 | `05d0105` | README、設定画面、WebHook説明、出力先選択などを改善。 |
| 014 | [VRChat写真の自動検知とDiscord投稿](014-vrchat-photo-auto-post.md) | 完了 | `7efdf8d` | VRChat写真フォルダの定期スキャンと自動Discord投稿を実装。 |
| 015 | [セキュリティチェック報告書の作成](015-security-review-report.md) | 完了 | `4c77718` | `SECURITY_REVIEW.md` を作成し、リスクと推奨対応を整理。 |
| 016 | [Webhook URL と履歴 URL の検証強化](016-validate-webhook-and-history-urls.md) | 完了 | `73fd6bf` | Discord Webhook URL と履歴画像URLの検証を強化。 |
| 017 | [設定・履歴ファイルの権限と排他制御](017-harden-local-secret-storage-and-history-locking.md) | 完了 | `6df353c` | 設定/履歴ファイルの権限、履歴更新の排他制御を強化。 |
| 018 | [画像入力のサイズ上限とピクセル数上限](018-limit-image-input-resource-usage.md) | 完了 | `bbfc9a3` | 入力画像のファイルサイズ上限とピクセル数上限を追加。 |
| 019 | [VRChat写真自動投稿の走査・処理件数制限](019-limit-auto-photo-scanning.md) | 完了 | `d9a7f58` | 自動投稿の走査量と処理件数を制限。 |
| 020 | [CI/Release のセキュリティチェック追加](020-add-security-checks-to-ci-release.md) | 完了 | `5c935b1` | CI/Release に監査や脆弱性チェックを追加。 |
| 021 | [開発用アイコン生成ツールの gosec 指摘対応](021-address-gosec-tooling-findings.md) | 完了 | `afda4bf` | gosec 指摘を受けた開発用ツールの処理を修正。 |
| 022 | [アプリの複数起動制限](022-prevent-multiple-app-instances.md) | 完了 | `ef82b87` | 多重起動防止を実装。 |
| 023 | [Goテストカバレッジの拡充](023-expand-go-test-coverage.md) | 完了 | `01a8400` | appcore とアプリ主要処理の Go テストを拡充。 |
| 024 | [Discord Webhook URLエラーの案内改善](024-improve-discord-webhook-error-message.md) | 完了 | `e3146b3` | 空/不正な Webhook URL のエラー文言をユーザー向けに改善。 |
| 025 | [設定画面の未保存変更確認](025-confirm-unsaved-settings-navigation.md) | 完了 | `f17eea0` / `v0.1.0` | 設定画面を離れる前に保存/破棄/キャンセルを選べる確認ダイアログを追加。 |
| 026 | [v0.1.1 不具合修正](026-fix-v0.1.1-regressions.md) | 完了 | `21cc973` / `v0.1.1`, 028で導線更新 | 履歴画面への導線、Win+Shift+S 画像崩れを修正。 |
| 027 | [画像履歴のマウス矩形範囲選択](027-history-drag-selection.md) | 完了 | 未リリース | 履歴画面でマウスドラッグによる矩形範囲選択を追加。 |
| 028 | [履歴画面への長押し導線をボタンへ変更](028-replace-history-long-press-with-button.md) | 完了 | 未リリース | クリア長押しを廃止し、履歴ボタンと各ボタンの説明 tooltip を追加。 |
| 029 | [画像履歴の全選択とピン止め](029-history-select-all-and-pin.md) | 完了 | 未リリース | Ctrl+A/全選択ボタンと、削除対象外にするピン止めを追加。 |
| 030 | [Discord公式ヘルプリンクが開けない](030-fix-discord-help-link.md) | 完了 | 未リリース | Discord公式ヘルプURLからテキストフラグメントを削除。 |
| 031 | [Discord削除済み履歴の削除条件とoutput削除](031-purge-discord-deleted-history.md) | 完了 | 未リリース | Discord削除済みフラグで履歴を削除し、設定に応じてoutputも削除。 |
| 032 | [QRコードURLの検出とDiscord投稿](032-detect-qr-url-and-post.md) | 完了 | 未リリース | 画像内QRコードのURLを検出し、Discord投稿本文と結果画面に表示する。 |
| 033 | [WSLからのWindowsビルドスクリプト](033-wsl-windows-build-script.md) | 完了 | 未リリース | WSLからWindows向けexeをコミットハッシュ版としてローカルビルドするスクリプトを追加。 |
| 034 | [GitHub ActionsビルドのバージョンにコミットIDを含める](034-release-version-includes-commit.md) | 完了 | 未リリース | Release workflow で表示バージョンをリリース番号.コミットIDにする。 |
| 035 | [処理診断ログの追加](035-add-processing-diagnostic-log.md) | 完了 | 未リリース | QR検出などの原因調査用に処理時の診断ログを出力する。 |
| 037 | [Release exe の外部署名ファイル](037-release-detached-gpg-signature.md) | 完了 | 未リリース | ReleaseでexeのPGP外部署名ファイルを生成して添付する。 |
| 038 | [情報画面に公式配布元とPGP確認方法を追加](038-about-official-distribution-and-pgp.md) | 完了 | 未リリース | 公式配布元とexe署名による改竄確認方法を情報画面に記載する。 |
| 039 | [バージョンとリビジョンを分けてビルドに埋め込む](039-separate-version-and-revision-build-info.md) | 完了 | 未リリース | ビルド時にバージョンとリビジョンを別々に埋め込む。 |
| 040 | [GitHub Release のアップデート通知](040-check-github-release-updates.md) | 完了 | 未リリース | GitHub Releases の最新Releaseを確認し、更新があればUI内に通知する。 |
| 041 | [Win+Shift+S スクリーンショット自動投稿](041-screenshot-auto-post.md) | 完了 | 未リリース | Screenshotsフォルダを定期スキャンし、Win+Shift+Sで保存された画像を自動投稿する。 |
| 046 | [初回設定を保存せず終了した場合に次回も初期設定画面を出す](046-fix-initial-config-save-timing.md) | 完了 | 未リリース | 初回設定画面を開いただけで `config.json` が作成されないようにする。 |
| 047 | [CLIでversion/help引数に対応する](047-cli-version-help-args.md) | 完了 | 未リリース | `go-arg` で `--version` / `--help` の早期終了に対応する。 |
| 048 | [Windows GUI exeのCLI出力をPowerShellに表示する](048-cli-output-from-windows-gui-exe.md) | 完了 | 未リリース | GUIサブシステムのexeでもCLI出力を親コンソールへ表示する。 |
| 049 | [CLIヘルプをWindowsコンソールで文字化けさせない](049-use-wide-console-output-for-cli.md) | 完了 | 未リリース | `WriteConsoleW` でCLI出力し、cmdのコードページ差による文字化けを避ける。 |
| 050 | [プロダクト全体の問題点チェック報告書を作成する](050-product-issue-review-report.md) | 完了 | 未リリース | `PRODUCT_ISSUE_REPORT.md` を作成し、Release/UX/履歴/設定/OSS表示などの残課題を整理。 |

## 状態の意味

| 状態 | 意味 |
| --- | --- |
| 完了 | 実装、テスト、またはドキュメント作成が完了している。 |
| 要確認 | 修正は入っているが、対象環境や実機での再現確認が必要。 |
| 一部将来対応あり | 主要な受け入れ条件は満たしているが、明示的に将来対応へ回した項目が残っている。 |
