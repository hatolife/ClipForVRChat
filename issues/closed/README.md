# Closed Issues

このディレクトリは、完了済みissueを保管するためのものです。

| No. | Issue | 状態 | 対応バージョン | 概要 |
| --- | --- | --- | --- | --- |
| 001 | [アイコン品質と複数サイズ生成](v0.1.0/001-icon-quality-and-sizes.md) | 完了 | `v0.1.0` | アプリアイコン生成と複数サイズの品質改善。 |
| 002 | [Release zipの内容とバージョン付きファイル名](v0.1.0/002-release-zip-contents-and-versioned-name.md) | 完了 | `v0.1.0` | 配布zipの内容整理、バージョン情報の埋め込み、Release workflow 整備。 |
| 003 | [メインウィンドウUIと情報表示](v0.1.0/003-main-window-ui-and-about.md) | 完了 | `v0.1.0` | メイン画面、情報表示、アプリ情報取得の追加。 |
| 004 | [初回設定フロー](v0.1.0/004-initial-config-flow.md) | 完了 | `v0.1.0` | GUIアプリとして設定を扱える初期フローを実装。 |
| 005 | [ウィンドウへのドラッグ&ドロップ](v0.1.0/005-window-drag-and-drop.md) | 完了 | `v0.1.0` | 画像/config.json のドラッグ&ドロップ処理と画面全体のドロップ対応。 |
| 006 | [クリップボード画像の保存ファイル名に時刻を含める](v0.1.0/006-clipboard-output-filename.md) | 完了 | `v0.1.0` | クリップボード入力の保存ファイル名を時刻付きに変更。 |
| 007 | [クリップボードのスクリーンショット画像が崩れる](v0.1.1/007-clipboard-screenshot-corruption.md) | 完了 | `21cc973` / `v0.1.1` | Win+Shift+S 由来画像の崩れ対策として Windows 登録PNG読み取りを優先。 |
| 008 | [UI仕様の明文化](v0.1.0/008-ui-specification.md) | 完了 | `v0.1.0` | UI仕様と運用ルールをドキュメント化。 |
| 009 | [ソース配置を src 配下へ移動する](v0.1.0/009-move-source-under-src.md) | 完了 | `v0.1.0` | Wails/Go/frontend ソースを `src` 配下へ整理。 |
| 010 | [設定画面・結果画面・情報画面のUI改善](v0.1.0/010-settings-results-about-ui-improvements.md) | 完了 | `v0.1.0` | 設定/結果/情報画面のレイアウト、進捗、クリア、未保存確認を改善。 |
| 011 | [使い方画面とWebhook発行案内](v0.1.0/011-help-screen-and-webhook-guide.md) | 完了 | `v0.1.0` | 使い方画面と Discord Webhook 公式案内リンクを追加。 |
| 012 | [画像URL履歴とDiscord削除確認画面](v0.1.0/012-image-history-and-discord-delete-review.md) | 完了 | `v0.1.0`, `v0.1.2` | 履歴保存、クリア済み表示、Discord削除確認画面、Ctrl/Shift選択、矩形範囲選択を実装。 |
| 013 | [ユーザーフレンドリーなREADMEと設定画面改善](v0.1.0/013-user-friendly-readme-and-settings-ui.md) | 完了 | `v0.1.0` | README、設定画面、WebHook説明、出力先選択などを改善。 |
| 014 | [VRChat写真の自動検知とDiscord投稿](v0.1.0/014-vrchat-photo-auto-post.md) | 完了 | `v0.1.0` | VRChat写真フォルダの定期スキャンと自動Discord投稿を実装。 |
| 015 | [セキュリティチェック報告書の作成](v0.1.0/015-security-review-report.md) | 完了 | `v0.1.0` | `reports/2026-06-21/security-review.md` を作成し、リスクと推奨対応を整理。 |
| 016 | [Webhook URL と履歴 URL の検証強化](v0.1.0/016-validate-webhook-and-history-urls.md) | 完了 | `v0.1.0` | Discord Webhook URL と履歴画像URLの検証を強化。 |
| 017 | [設定・履歴ファイルの権限と排他制御](v0.1.0/017-harden-local-secret-storage-and-history-locking.md) | 完了 | `v0.1.0` | 設定/履歴ファイルの権限、履歴更新の排他制御を強化。 |
| 018 | [画像入力のサイズ上限とピクセル数上限](v0.1.0/018-limit-image-input-resource-usage.md) | 完了 | `v0.1.0` | 入力画像のファイルサイズ上限とピクセル数上限を追加。 |
| 019 | [VRChat写真自動投稿の走査・処理件数制限](v0.1.0/019-limit-auto-photo-scanning.md) | 完了 | `v0.1.0` | 自動投稿の走査量と処理件数を制限。 |
| 020 | [CI/Release のセキュリティチェック追加](v0.1.0/020-add-security-checks-to-ci-release.md) | 完了 | `v0.1.0` | CI/Release に監査や脆弱性チェックを追加。 |
| 021 | [開発用アイコン生成ツールの gosec 指摘対応](v0.1.0/021-address-gosec-tooling-findings.md) | 完了 | `v0.1.0` | gosec 指摘を受けた開発用ツールの処理を修正。 |
| 022 | [アプリの複数起動制限](v0.1.0/022-prevent-multiple-app-instances.md) | 完了 | `v0.1.0` | 多重起動防止を実装。 |
| 023 | [Goテストカバレッジの拡充](v0.1.0/023-expand-go-test-coverage.md) | 完了 | `v0.1.0` | appcore とアプリ主要処理の Go テストを拡充。 |
| 024 | [Discord Webhook URLエラーの案内改善](v0.1.0/024-improve-discord-webhook-error-message.md) | 完了 | `v0.1.0` | 空/不正な Webhook URL のエラー文言をユーザー向けに改善。 |
| 025 | [設定画面の未保存変更確認](v0.1.0/025-confirm-unsaved-settings-navigation.md) | 完了 | `f17eea0` / `v0.1.0` | 設定画面を離れる前に保存/破棄/キャンセルを選べる確認ダイアログを追加。 |
| 026 | [v0.1.1 不具合修正](v0.1.1/026-fix-v0.1.1-regressions.md) | 完了 | `21cc973` / `v0.1.1`, 028で導線更新 | 履歴画面への導線、Win+Shift+S 画像崩れを修正。 |
| 027 | [画像履歴のマウス矩形範囲選択](v0.1.2/027-history-drag-selection.md) | 完了 | `v0.1.2` | 履歴画面でマウスドラッグによる矩形範囲選択を追加。 |
| 028 | [履歴画面への長押し導線をボタンへ変更](v0.1.2/028-replace-history-long-press-with-button.md) | 完了 | `v0.1.2` | クリア長押しを廃止し、履歴ボタンと各ボタンの説明 tooltip を追加。 |
| 029 | [画像履歴の全選択とピン止め](v0.1.2/029-history-select-all-and-pin.md) | 完了 | `v0.1.2` | Ctrl+A/全選択ボタンと、削除対象外にするピン止めを追加。 |
| 030 | [Discord公式ヘルプリンクが開けない](v0.1.2/030-fix-discord-help-link.md) | 完了 | `v0.1.2` | Discord公式ヘルプURLからテキストフラグメントを削除。 |
| 031 | [Discord削除済み履歴の削除条件とoutput削除](v0.1.3/031-purge-discord-deleted-history.md) | 完了 | `v0.1.3` | Discord削除済みフラグで履歴を削除し、設定に応じてoutputも削除。 |
| 032 | [QRコードURLの検出とDiscord投稿](v0.1.4/032-detect-qr-url-and-post.md) | 完了 | `v0.1.4` | 画像内QRコードのURLを検出し、Discord投稿本文と結果画面に表示する。 |
| 033 | [WSLからのWindowsビルドスクリプト](v0.1.4/033-wsl-windows-build-script.md) | 完了 | `v0.1.4` | WSLからWindows向けexeをコミットハッシュ版としてローカルビルドするスクリプトを追加。 |
| 034 | [GitHub ActionsビルドのバージョンにコミットIDを含める](v0.1.4/034-release-version-includes-commit.md) | 完了 | `v0.1.4` | Release workflow で表示バージョンをリリース番号.コミットIDにする。 |
| 035 | [処理診断ログの追加](v0.1.4/035-add-processing-diagnostic-log.md) | 完了 | `v0.1.4` | QR検出などの原因調査用に処理時の診断ログを出力する。 |
| 037 | [Release exe の外部署名ファイル](v0.1.4/037-release-detached-gpg-signature.md) | 完了 | `v0.1.4` | ReleaseでexeのPGP外部署名ファイルを生成して添付する。 |
| 038 | [情報画面に公式配布元とPGP確認方法を追加](v0.1.4/038-about-official-distribution-and-pgp.md) | 完了 | `v0.1.4` | 公式配布元とexe署名による改竄確認方法を情報画面に記載する。 |
| 039 | [バージョンとリビジョンを分けてビルドに埋め込む](v0.1.4/039-separate-version-and-revision-build-info.md) | 完了 | `v0.1.4` | ビルド時にバージョンとリビジョンを別々に埋め込む。 |
| 040 | [GitHub Release のアップデート通知](v0.1.4/040-check-github-release-updates.md) | 完了 | `v0.1.4` | GitHub Releases の最新Releaseを確認し、更新があればUI内に通知する。 |
| 041 | [Win+Shift+S スクリーンショット自動処理](v0.1.5/041-screenshot-auto-post.md) | 完了 | `v0.1.5` | Screenshotsフォルダを定期スキャンし、Win+Shift+Sで保存された画像を自動処理する。 |
| 042 | [Windows exe のプロパティ表示改善](v0.1.6/042-windows-exe-version-info.md) | 完了 | `v0.1.6` | Windowsのファイルプロパティに製品名、説明、著作権、バージョンを追加。 |
| 043 | [v0.1.6 リリースノート作成](v0.1.6/043-release-notes-v0.1.6.md) | 完了 | `v0.1.6` | `RELEASE_NOTES.md` にv0.1.6の更新内容と配布URLを追加。 |
| 044 | [v0.1.6 Release workflow のVersionInfo生成失敗修正](v0.1.6/044-fix-v0.1.6-release-workflow-version-info.md) | 完了 | `v0.1.6` | Release workflowのVersionInfo生成先ディレクトリ作成漏れを修正。 |
| 045 | [Release本文と成果物の再発防止](v0.1.6/045-fix-release-assets-and-notes.md) | 完了 | `v0.1.6` | Release本文と成果物一覧が仕様通りになるようworkflowと検査を修正。 |
| 046 | [初回設定を保存せず終了した場合に次回も初期設定画面を出す](v0.1.6/046-fix-initial-config-save-timing.md) | 完了 | `v0.1.6` | 初回設定画面を開いただけで `config.json` が作成されないようにする。 |
| 047 | [CLIでversion/help引数に対応する](v0.1.6/047-cli-version-help-args.md) | 完了 | `v0.1.6` | `go-arg` で `--version` / `--help` の早期終了に対応する。 |
| 048 | [Windows GUI exeのCLI出力をPowerShellに表示する](v0.1.6/048-cli-output-from-windows-gui-exe.md) | 完了 | `v0.1.6` | GUIサブシステムのexeでもCLI出力を親コンソールへ表示する。 |
| 049 | [CLIヘルプをWindowsコンソールで文字化けさせない](v0.1.6/049-use-wide-console-output-for-cli.md) | 完了 | `v0.1.6` | `WriteConsoleW` でCLI出力し、cmdのコードページ差による文字化けを避ける。 |
| 050 | [プロダクト全体の問題点チェック報告書を作成する](v0.1.6/050-product-issue-review-report.md) | 完了 | `v0.1.6` | `reports/2026-06-24/product-issue-report.md` を作成し、Release/UX/履歴/設定/OSS表示などの残課題を整理。 |
| 051 | [監査報告書を日付付きで専用ディレクトリへ移動する](v0.1.6/051-date-audit-reports-directory.md) | 完了 | `v0.1.6` | 監査報告書を `reports/` 配下の日付ディレクトリへ整理。 |
| 052 | [Release署名ファイルの説明を実際の成果物に合わせる](v0.1.6/052-align-release-signing-docs.md) | 完了 | `v0.1.6` | README/About/SPEC のPGP確認手順をRelease成果物に合わせる。 |
| 053 | [Discord削除の部分成功を履歴へ保存する](v0.1.6/053-persist-partial-discord-delete-results.md) | 完了 | `v0.1.6` | 複数削除の一部失敗時も成功済み削除状態を保存する。 |
| 054 | [config読み込み失敗時にアクティブ設定パスを変更しない](v0.1.6/054-keep-config-path-on-load-error.md) | 完了 | `v0.1.6` | 不正configの読み込み失敗時に既存の設定パスを維持する。 |
| 055 | [URL自動コピー失敗をユーザーへ表示する](v0.1.6/055-report-copy-single-url-failure.md) | 完了 | `v0.1.6` | URL取得成功とクリップボードコピー失敗を分けて表示する。 |
| 056 | [OSSライセンス表示を依存関係に合わせて更新する](v0.1.6/056-complete-oss-license-list.md) | 完了 | `v0.1.6` | direct dependency のライセンス表示漏れを補う。 |
| 057 | [自動投稿の監視フォルダ異常を通知する](v0.1.6/057-report-auto-post-watch-diagnostics.md) | 完了 | `v0.1.6` | 監視フォルダ異常やスキャン上限到達をUIまたは診断ログへ出す。 |
| 058 | [ローカルWindowsビルドのバージョンを指定できるようにする](v0.1.6/058-allow-local-build-version-override.md) | 完了 | `v0.1.6` | ローカルビルド時にリリース候補バージョンを指定可能にする。 |
| 059 | [診断ログのgosec指摘を整理する](v0.1.6/059-resolve-diagnostic-log-gosec-finding.md) | 完了 | `v0.1.6` | 診断ログの可変パス指摘を設計上明確にする。 |
| 060 | [issue一覧のv0.1.6対応状況を整理する](v0.1.6/060-update-issue-index-release-status.md) | 完了 | `v0.1.6` | issue一覧の未リリース表記と欠番を整理。 |
| 061 | [リリース時にissue一覧へ対応バージョンを記録する](v0.1.6/061-record-release-version-in-issue-index.md) | 完了 | `v0.1.6` | リリース時の対応バージョン記録ルールを追加し、既存一覧を補完。 |
| 062 | [issue一覧の対応バージョン列名を整理する](v0.1.6/062-rename-issue-index-version-column.md) | 完了 | `v0.1.6` | issue一覧の列名を対応バージョン記録の運用に合わせる。 |
| 063 | [v0.1.6履歴のコミットメッセージと署名を整理する](v0.1.6/063-rewrite-v016-history-angular-signed.md) | 完了 | `v0.1.6` | v0.1.6向け履歴のコミットメッセージ、署名、相殺revertを整理。 |
| 064 | [master以外のブランチCIと非正式タグのdraft release対応](v0.1.6/064-run-ci-on-branches-and-draft-nonstable-tags.md) | 完了 | `v0.1.6` | ブランチCI対象と非正式タグのdraft Release作成を整理。 |
| 065 | [開発ブランチ運用とRCリリースを整備する](v0.1.6/065-dev-flow-and-rc-release.md) | 完了 | `v0.1.6` | 作業ブランチ前提の開発フローと `vX.Y.Z-rcW` prerelease 発行を整理。 |
| 066 | [GitHub ActionsのNode 20警告とGo cache警告を解消する](v0.1.6/066-update-actions-runtime-and-go-cache.md) | 完了 | `v0.1.6` | Actions実行系とGo cache設定を更新し、CI/Releaseの警告を解消。 |
| 067 | [todo.mdを直近作業用に整理する](v0.1.6/067-maintain-short-term-todo.md) | 完了 | `v0.1.6` | `todo.md` を短期チェックリストとして整理し、運用ルールを明文化。 |
| 068 | [GitHub Actions CIを高速化する](v0.1.7/068-speed-up-github-actions-ci.md) | 完了 | `v0.1.7` | Wails CLI cacheを追加し、CI/Releaseのインストール時間短縮を試行。 |
| 069 | [更新通知の開き先選択と通知設定](v0.1.7/069-update-notification-destination-and-settings.md) | 完了 | `v0.1.7` | 更新通知からGitHub/BOOTHを選んで開き、更新確認と通知を設定でON/OFFできるようにする。 |
| 070 | [設定画面をカテゴリ別タブに分ける](v0.1.7/070-tabbed-settings-page.md) | 完了 | `v0.1.7` | 設定画面上部にカテゴリタブを追加し、設定項目を分類して表示する。 |
| 071 | [Discord投稿OFFでも自動投稿でDiscordへ投稿される](v0.1.7/071-prevent-discord-post-when-disabled.md) | 完了 | `v0.1.7` | Discord投稿OFF時に自動投稿経路からDiscordへ投稿されないようにする。 |
| 072 | [設定画面仕様を独立ファイルへ分離する](v0.1.7/072-split-settings-screen-spec.md) | 完了 | `v0.1.7` | 設定画面仕様を専用Markdownへ分離し、カテゴリと依存関係を整理する。 |
| 073 | [設定画面を機能/処理/Webhook/更新へ再編する](v0.1.7/073-reorganize-settings-screen-categories.md) | 完了 | `v0.1.7` | 設定画面を利用者視点のカテゴリへ再編する。 |
| 074 | [結果画像のクリック領域をURLコピーと保存先表示に分ける](v0.1.7/074-result-card-split-actions.md) | 完了 | `v0.1.7` | 結果画像の上半分でURLコピー、下半分で保存先表示を行えるようにする。 |
| 075 | [履歴画面をDiscord/ローカル/履歴の状態別操作に作り直す](v0.1.7/075-rebuild-history-screen.md) | 完了 | `v0.1.7` | 履歴画面を状態表示とDiscord/ローカル/履歴削除の独立操作へ作り直す。 |
| 076 | [履歴と結果表示を実際に行った処理に合わせる](v0.1.7/076-history-and-result-display-follow-actual-work.md) | 完了 | `v0.1.7` | 結果画面と履歴画面を実際に行った処理内容に合わせて表示する。 |
| 077 | [処理結果がない場合に理由を表示する](v0.1.7/077-explain-no-result-processing-message.md) | 完了 | `v0.1.7` | 明示的な処理で結果がない場合、設定により処理結果が出ない理由を表示する。 |
| 078 | [履歴のローカル保存パスをconfig基準で解決する](v0.1.7/078-resolve-history-local-paths-relative-to-config.md) | 完了 | `v0.1.7` | 履歴の相対ローカル保存パスをconfig基準で解決し、Explorer表示と削除可否を修正する。 |
| 079 | [ユーザー操作を診断ログへ記録する](v0.1.7/079-record-user-actions-in-diagnostic-log.md) | 完了 | `v0.1.7` | ボタンクリック、画面遷移、処理判断、処理結果表示を診断ログへ記録する。 |
| 080 | [暗号化診断パッケージ作成](v0.1.7/080-create-encrypted-diagnostic-package.md) | 完了 | `v0.1.7` | 起動ログと診断パッケージ作成導線を追加。暗号化用公開鍵で診断パッケージを作成する。 |
| 081 | [診断ログを日付付きログフォルダへ出力する](v0.1.7/081-write-diagnostic-logs-to-dated-log-directory.md) | 完了 | `v0.1.7` | 診断ログを `logs/YYYY-MM-DD.log` に出力し、診断パッケージがログフォルダを収集するようにする。 |
| 082 | [診断パッケージ暗号化鍵をpoppo@hato.lifeへ変更する](v0.1.7/082-use-poppo-openpgp-key-for-diagnostics.md) | 完了 | `v0.1.7` | 診断パッケージ暗号化用の公開鍵を暗号化サブキー付きの `poppo@hato.life` に変更する。 |
| 083 | [不具合報告用データ生成UIを改善する](v0.1.7/083-improve-diagnostic-data-generation-ui.md) | 完了 | `v0.1.7` | 情報画面の不具合報告項目に生成ボタンを配置し、生成中オーバーレイと完了後Explorer表示を追加する。 |
| 084 | [Explorerで生成ファイルを選択表示できない](v0.1.7/084-fix-explorer-file-selection.md) | 完了 | `v0.1.7` | `explorer.exe /select,` ではなく Windows Shell API でファイル選択表示する。 |
| 085 | [診断データ復号後のzipが破損する](v0.1.7/085-encrypt-diagnostic-package-as-binary.md) | 完了 | `v0.1.7` | 診断データをOpenPGP binary literal dataとして暗号化し、復号後zipの破損を防ぐ。 |
| 086 | [診断zip内のログとoutput構成を整理する](v0.1.7/086-normalize-diagnostic-zip-layout-and-output-log.md) | 完了 | `v0.1.7` | zip内ログを `logs/` に統一し、output画像は含めず一覧を診断ログへ記録する。 |
| 087 | [zipファイル引数を公開鍵で暗号化する](v0.1.7/087-encrypt-zip-from-cli-argument.md) | 完了 | `v0.1.7` | zip単体引数でUIを起動せず、同じ公開鍵で `<zip>.gpg` を生成する。 |
| 088 | [診断ログ実フォルダ名をlogsへ統一する](v0.1.7/088-use-logs-directory-for-diagnostic-logs.md) | 完了 | `v0.1.7` | アプリ実フォルダの診断ログ出力先も `logs/YYYY-MM-DD.log` に統一する。 |
| 089 | [不具合報告用データの説明を情報画面に追加する](v0.1.7/089-explain-diagnostic-data-in-about.md) | 完了 | `v0.1.7` | 情報画面に含まれる情報、暗号化方法、利用目的の説明を追加する。 |
| 090 | [不具合報告用データを段階作成しパスを置換する](v0.1.7/090-stage-and-sanitize-diagnostic-data.md) | 完了 | `v0.1.7` | `diagnostics/<timestamp>/` に確認用zipと暗号化zipを作成し、テキスト内パスを環境変数表記へ置換する。 |
| 091 | [セキュリティ監査報告書を作成する](commit-only/091-security-audit-2026-06-25.md) | 完了 | `8e46ffd` | `reports/2026-06-25/security-audit-prompt.md` に基づき、現行リポジトリのセキュリティ監査報告書を作成した。 |
| 092 | [監査報告書の配置と説明を整理する](v0.1.7/092-organize-reports-and-clarify-security-notes.md) | 完了 | `v0.1.7` | 監査報告書を日付ディレクトリへ移動し、診断zip、Release workflow、Windows ACLの説明を補足する。 |
| 093 | [人間が確認する必要がある作業を手順化する](v0.1.7/093-human-verification-guide.md) | 完了 | `v0.1.7` | 監査後に人間が確認する作業、現時点の確認結果、判断基準をMarkdown化する。 |
| 094 | [設定画面の初期タブと初期値を調整する](v0.1.7/094-settings-initial-tab-and-defaults.md) | 完了 | `v0.1.7` | 初回起動時の設定タブ選択とDiscord投稿/投稿URL自動コピーの初期値を調整する。 |
| 095 | [設定画面のDiscord投稿分類を整理する](v0.1.7/095-settings-discord-post-tab-reclassification.md) | 完了 | `v0.1.7` | WebhookタブをDiscord投稿へ変更し、Discord投稿関連設定をまとめる。 |
| 096 | [develop版バージョンに親コミットIDを含める](v0.1.7/096-develop-version-include-parent-commit.md) | 完了 | `v0.1.7` | develop版バージョン表記に親コミットIDを含めて追跡しやすくする。 |
| 097 | [govulncheckのGO-2026-4550を解消する](v0.1.7/097-fix-circl-govulncheck-finding.md) | 完了 | `v0.1.7` | `cloudflare/circl` の既知脆弱性を解消し、govulncheckを成功させる。 |
| 098 | [診断データからWebhook URLとDiscord tokenを除外する](v0.1.7/098-redact-diagnostic-secrets.md) | 完了 | `v0.1.7` | 確認用zipと暗号化zipからWebhook URLとDiscord tokenの生値を除外する。 |
| 099 | [Release workflowの権限を最小化する](v0.1.7/099-harden-release-workflow-permissions.md) | 完了 | `v0.1.7` | Release作成権限を必要なjobへ限定する。 |
| 100 | [OpenURLで開けるURLを許可ホストに制限する](v0.1.7/100-restrict-open-url-hosts.md) | 完了 | `v0.1.7` | アプリから開ける外部URLを信頼済みHTTPSホストへ限定する。 |
| 101 | [履歴のローカル削除対象を管理output配下に制限する](v0.1.7/101-restrict-history-local-delete-to-output.md) | 完了 | `v0.1.7` | 履歴改ざん時にoutput外ファイルを削除できないようにする。 |
| 102 | [診断ログのローカルパスと秘密情報を抑制する](v0.1.7/102-redact-diagnostic-log-paths.md) | 完了 | `v0.1.7` | 診断ログ出力時点でパスと秘密情報を可能な範囲で抑制する。 |
| 103 | [Windows実機でACLを確認する](v0.1.8/103-check-windows-acl.md) | 完了 | `v0.1.8` | Windows実機でconfig/history/logs/diagnosticsのACLを確認する。 |
| 104 | [ReleaseにSBOMまたはビルドメタデータを追加する](v0.1.7/104-add-release-sbom-or-build-metadata.md) | 完了 | `v0.1.7` | Release成果物に依存関係やビルド環境の追跡情報を追加する。 |
| 105 | [Windows GoテストでWailsイベント送信ガードが効かない](v0.1.7/105-fix-windows-test-wails-event-guard.md) | 完了 | `v0.1.7` | Windowsの `.test.exe` テストバイナリでもWailsイベント送信を抑制する。 |
| 106 | [Release NotesのダウンロードURLをMarkdownリンクにする](v0.1.7/106-release-notes-download-links.md) | 完了 | `v0.1.7` | Release本文のダウンロード欄でファイル名をリンクテキストとして表示する。 |
| 107 | [Git Flow運用をAGENTS.mdへ明文化する](v0.1.7/107-update-git-flow-agents.md) | 完了 | `v0.1.7` | `develop` 基準の通常開発と `master` 基準のリリースブランチ運用を明文化する。 |
| 108 | [Discord投稿ONで通常投稿用Webhook URLが空欄の場合に保存時警告を出す](v0.1.7/108-warn-empty-discord-webhook-on-save.md) | 完了 | `v0.1.7` | Discord投稿ONで通常投稿用Webhook URLが空欄の場合、保存後も画面上部に警告を表示する。 |
| 109 | [Release workflow のタグ名取り扱いを安全化する](v0.1.8/109-harden-release-tag-handling.md) | 完了 | `v0.1.8` | タグ名を検証し、シェルコマンドへ環境変数経由で渡してRelease workflowのコマンド注入を防ぐ。 |
| 110 | [VRChat自動構図撮影](v0.1.8/110-vrc-auto-composition-capture.md) | 完了 | `v0.1.8` | OSCでUser Cameraを制御し、構図ごとの自動撮影とsidecar JSON保存を実装する。 |
| 111 | [VRChat output logからの同席ユーザー保持](v0.1.8/111-vrc-output-log-presence-users.md) | 完了 | `v0.1.8` | output_log監視で撮影時点の同席ユーザー情報を保持し、画像メタデータやDiscord投稿へ紐づける。 |
| 112 | [自動撮影でVRChat写真が保存されない](v0.1.8/112-fix-auto-capture-action-osc.md) | 完了 | `v0.1.8` | User CameraのCapture/CloseをAction OSCとして送信し、写真保存されない問題を修正する。 |
| 113 | [自動撮影の診断ログを詳細化する](v0.1.8/113-add-auto-capture-diagnostics.md) | 完了 | `v0.1.8` | 自動撮影のスケジューラ、OSC送信、写真検出状況を診断ログで追跡できるようにする。 |
| 114 | [自動撮影タブに機能説明を追加する](v0.1.8/114-explain-auto-capture-settings-tab.md) | 完了 | `v0.1.8` | 自動撮影タブの先頭に機能概要と使い方の説明枠を追加する。 |
| 115 | [自動撮影のCapture送信とCloseタイミングを調整する](v0.1.8/115-fix-auto-capture-button-action-and-close.md) | 完了 | `v0.1.8` | Capture/Closeを押下・解放OSCとして送信し、全失敗時にカメラを閉じないようにする。 |
| 116 | [自動撮影の写真検出ずれとカメラ未表示時の案内を改善する](v0.1.8/116-fix-auto-capture-photo-detection-and-camera-open-note.md) | 完了 | `v0.1.8` | 遅れて保存された写真の検出ずれを防ぎ、Photo方式でUser Camera表示が必要なことを案内する。 |
| 117 | [自動撮影の現在Pose保存と構図管理を実装する](v0.1.8/117-implement-camera-pose-preset-calibration.md) | 完了 | `v0.1.8` | VRChatから現在のUser Camera Poseを受信し、構図プリセットとして保存・管理できるようにする。 |
| 118 | [自動撮影の解像度一時変更はv0.1.8で断念する](v0.1.8/118-defer-auto-capture-resolution-control.md) | 完了 | `v0.1.8` | v0.1.8ではVRChatの現在のカメラ解像度設定を使用し、未完成の解像度変更設定を出さない。 |
| 119 | [未実装の自動撮影方式を設定画面から外す](v0.1.8/119-remove-unimplemented-auto-capture-mode-options.md) | 完了 | `v0.1.8` | 実装済みのStream(ffmpeg)方式とPhoto方式だけを表示し、未実装方式を設定画面から外す。 |
| 120 | [自動撮影のPose操作と初期構図設定を修正する](v0.1.8/120-fix-auto-capture-pose-controls-and-defaults.md) | 完了 | `v0.1.8` | 初期Pose、拡大率、撮影トグル、現在Pose保存/追加APIを修正する。 |
| 122 | [自動撮影Stream方式のffmpeg確認と導入導線を追加する](v0.1.8/122-auto-capture-ffmpeg-status-and-install.md) | 完了 | `v0.1.8` | Stream方式に必要なffmpegの確認、未導入表示、winget導入ボタンを追加する。 |
| 123 | [自動撮影Stream方式でデスクトップ全体を撮らない](v0.1.8/123-avoid-desktop-ffmpeg-stream-input.md) | 完了 | `v0.1.8` | ffmpeg入力引数の初期値と既存設定移行でデスクトップ全体撮影を避ける。 |
| 124 | [自動撮影Stream方式で白画像になるtitle取得を避ける](v0.1.8/124-capture-vrchat-window-region-instead-of-title.md) | 完了 | `v0.1.8` | VRChatウィンドウtitle直接取得ではなく画面上のウィンドウ範囲を切り出す。 |
| 125 | [構図ごとに現在Pose追加とカメラ移動ボタンを配置する](v0.1.8/125-add-per-view-pose-add-and-move-camera.md) | 完了 | `v0.1.8` | 構図カード内に現在Pose追加と設定Poseへのカメラ移動ボタンを追加する。 |
| 126 | [自動撮影OSCの押下状態を解除するデバッグボタンを追加する](v0.1.8/126-add-auto-capture-osc-recovery-button.md) | 完了 | `v0.1.8` | User Camera関連OSCをfalse/Offへ戻すデバッグボタンを追加する。 |
| 127 | [v0.1.8最小要件のtodoを再整理する](v0.1.8/127-restate-v018-minimum-requirements-todo.md) | 完了 | `v0.1.8` | v0.1.8の最小要件を大項目/小項目でtodoへ再整理する。 |
| 128 | [VRChat Stream Camera/Spout映像そのものを保存する方式を調査する](v0.1.8/128-investigate-vrchat-stream-camera-spout-capture.md) | 完了 | `v0.1.8` | Stream Camera/Spout出力を直接受信して静止画保存する方式を調査し、実装issueへ分割する。 |
| 133 | [Stream方式UIとドキュメントをSpout前提へ更新する](v0.1.8/133-update-auto-capture-stream-ui-and-docs-for-spout.md) | 完了 | `v0.1.8` | 自動撮影タブとドキュメントからFFmpeg主経路の誤解をなくし、Stream Camera(Spout)前提にする。 |
| 135 | [Stream Camera/Spout方式の実機確認手順を整備する](v0.1.8/135-add-spout-stream-camera-verification-guide.md) | 完了 | `v0.1.8` | RC実機確認でStream Camera映像そのものを保存できているか確認する手順を作る。 |
| 136 | [プレイヤー中心ローカル座標系の仕様を確定する](v0.1.8/136-investigate-player-local-camera-coordinate-spec.md) | 完了 | `v0.1.8` | 自動撮影構図をプレイヤー位置/向き基準のローカル座標として扱う仕様と実装方針を調査し、実装issueへ分割する。 |
| 137 | [EXIFへ同席ユーザー情報を書き込む方式を調査する](v0.1.8/137-investigate-exif-presence-user-metadata.md) | 完了 | `v0.1.8` | 自動撮影画像へ同席ユーザー情報をEXIF/XMP/PNGメタデータとして保持する方式を調査し、実装issueへ分割する。 |
| 138 | [`player_local` 座標仕様を確定する](v0.1.8/138-define-player-local-coordinate-spec.md) | 完了 | `v0.1.8` | プレイヤー中心ローカル構図の原点、軸、回転、初期値、既存設定互換を定義する。 |
| 139 | [ローカルプレイヤー基準Poseの取得可否を実機調査する](v0.1.8/139-investigate-player-basis-source.md) | 完了 | `v0.1.8` | 標準OSC/OSCQuery/ログ等からプレイヤー位置と向きを取得できるか確認する。 |
| 140 | [`player_local` からUser Camera Poseへの変換処理を実装する](v0.1.8/140-implement-player-local-pose-transform.md) | 完了 | `v0.1.8` | プレイヤー基準Poseがある場合にローカル構図をワールドPoseへ変換する。 |
| 142 | [プレイヤー中心構図の実機確認手順を整備する](v0.1.8/142-verify-player-local-camera-compositions.md) | 完了 | `v0.1.8` | Desktop/VR/アバター差分を含めた正面・背後・斜め構図の確認手順を作る。 |
| 143 | [自動撮影の埋め込みメタデータschemaを確定する](v0.1.8/143-define-autocapture-embedded-metadata-schema.md) | 完了 | `v0.1.8` | EXIF/PNGメタデータへ入れる撮影情報と同席ユーザー情報のschema、サイズ、プライバシー方針を決める。 |
| 144 | [PNG/JPEGの埋め込みメタデータwriterを追加する](v0.1.8/144-implement-image-metadata-writer.md) | 完了 | `v0.1.8` | JPEG EXIF APP1、PNG eXIf/iTXtへ非破壊でメタデータを書き込むwriterを追加する。 |
| 145 | [自動撮影保存処理へEXIF/PNGメタデータ書き込みを統合する](v0.1.8/145-integrate-exif-presence-metadata.md) | 完了 | `v0.1.8` | `finalizeAutoCaptureImage` に埋め込みメタデータ書き込みを接続し、sidecar/Discordと整合させる。 |
| 146 | [EXIF/埋め込みメタデータ設定とプライバシー説明をUIへ出す](v0.1.8/146-expose-exif-privacy-settings.md) | 完了 | `v0.1.8` | 自動撮影タブで埋め込みメタデータとユーザーID埋め込みを設定できるようにする。 |
| 147 | [埋め込みメタデータの読み戻し/Discord投稿後/実機検証を整備する](v0.1.8/147-verify-embedded-metadata-output.md) | 完了 | `v0.1.8` | PNG/JPEG読み戻し、Discord投稿後、ユーザー数過多時のメタデータ検証手順を整備する。 |
| 148 | [未実装/ダミー/簡易実装を洗い出す](v0.1.8/148-audit-incomplete-dummy-implementations.md) | 完了 | `v0.1.8` | 既存コードの未実装、仮実装、設定だけ存在する項目を洗い出して個別issueへ分割する。 |
| 149 | [自動撮影multi/Camera Dolly設定を実装または削除する](v0.1.8/149-implement-or-remove-autocapture-multi-camera-settings.md) | 完了 | `v0.1.8` | 保存/正規化されるmulti撮影設定を実装するか、未実装設定として隠す。 |
| 150 | [自動撮影スケジュールの重複実行制御を接続する](v0.1.8/150-connect-autocapture-scheduler-overlap-controls.md) | 完了 | `v0.1.8` | `skipIfPreviousBatchRunning` と `maxBatches` を実際のスケジューラ挙動へ反映する。 |
| 151 | [自動撮影DiscordのpostMode/includeImagesを実装または削除する](v0.1.8/151-implement-or-remove-autocapture-discord-post-options.md) | 完了 | `v0.1.8` | `postMode` と `includeImages` が常にShot画像投稿になる未接続状態を解消する。 |
| 152 | [構図ごとのcaptureDelayMsを撮影待機へ反映する](v0.1.8/152-apply-per-view-capture-delay.md) | 完了 | `v0.1.8` | 構図ごとの撮影直前待機設定をPhoto/Stream両方の撮影処理へ接続する。 |
| 153 | [自動撮影の出力形式/ファイル名テンプレートを設定画面へ出す](v0.1.8/153-expose-autocapture-output-format-and-filename-template.md) | 完了 | `v0.1.8` | `imageFormat` と `filenameTemplate` を自動撮影タブから設定できるようにする。 |
| 154 | [DiscordユーザーID出力をsidecar JSON設定から独立させる](v0.1.8/154-decouple-autocapture-discord-user-id-from-sidecar-user-id.md) | 完了 | `v0.1.8` | sidecar/Discord/EXIFのユーザーID出力制御を出力先ごとに独立させる。 |
| 155 | [sidecar JSONの履歴・削除ライフサイクルを定義する](v0.1.8/155-define-sidecar-json-lifecycle-with-history-delete.md) | 完了 | `v0.1.8` | 画像削除や履歴削除時にsidecar JSONをどう扱うか仕様化し実装へ接続する。 |
| 156 | [Wails公開APIとフロント呼び出しの同期チェックを追加する](v0.1.8/156-add-wails-api-surface-check.md) | 完了 | `v0.1.8` | フロントが呼ぶWails APIとGo公開メソッドの不一致をCI/検証で検出する。 |
| 157 | [v0.1.8自動撮影とRelease成果物仕様をREADME/SPECへ反映する](v0.1.8/157-sync-v018-autocapture-docs-and-specs.md) | 完了 | `v0.1.8` | README/設定仕様/SPECをv0.1.8自動撮影の実装済み範囲と制約に同期する。 |
| 158 | [自動撮影sidecarのworld/instance metadataを取得または削除する](v0.1.8/158-populate-or-remove-autocapture-world-instance-metadata.md) | 完了 | `v0.1.8` | sidecarに存在するworld/instance metadataフィールドを実データへ接続するかschemaから外す。 |
| 159 | [Discord投稿でallowed_mentionsを無効化する](v0.1.8/159-disable-discord-allowed-mentions.md) | 完了 | `v0.1.8` | Webhook投稿payloadで意図しないメンションを防ぐ。 |
| 160 | [自動撮影ローカルDB要件を実装または仕様から外す](v0.1.8/160-decide-autocapture-local-database-requirement.md) | 完了 | `v0.1.8` | SQLite/ローカルDB要件と現行sidecar/history方式の差分を解消する。 |
| 161 | [OSCQueryによるVRChat OSC検出を実装または延期明示する](v0.1.8/161-implement-or-defer-oscquery-discovery.md) | 完了 | `v0.1.8` | OSCQuery未実装の扱いを決め、実装または手動設定のみの仕様へ整理する。 |
| 162 | [自動撮影間隔のUI最小値をNormalizeと一致させる](v0.1.8/162-align-autocapture-interval-ui-validation.md) | 完了 | `v0.1.8` | 撮影間隔UIの最小値とNormalizeの丸め値を一致させる。 |
| 163 | [自動撮影テスト結果を設定画面に表示する](v0.1.8/163-show-autocapture-test-results-in-settings.md) | 完了 | `v0.1.8` | 構図ごとのテスト撮影成功/失敗結果を設定画面上で確認できるようにする。 |
| 164 | [v0.1.8完了扱い項目の再監査](v0.1.8/164-audit-v018-completed-items.md) | 完了 | `v0.1.8` | 完了扱いのv0.1.8項目に不完全な実装や未確認項目が残っていないか確認する。 |
| 165 | [v0.1.8未完了項目の実装前調査](v0.1.8/165-plan-v018-incomplete-implementation.md) | 完了 | `v0.1.8` | 未完了項目を実装可能な粒度に調査し、各issueへ実装方針を追記する。 |
| 167 | [Codex Security findingsを現在HEADで再検証し修正する](v0.1.8/167-remediate-codex-security-findings-2026-07-01.md) | 完了 | `v0.1.8` | Codex Security findingsを現在HEADで再検証し、未修正または部分修正の問題だけを安全に修正する。 |
| 169 | [Spout同梱バイナリの必要性と安全性説明を追加する](v0.1.8/169-document-spout-binary-necessity-and-safety.md) | 完了 | `v0.1.8` | `spout-capture.exe` と `SpoutLibrary.dll` の必要性、安全性の根拠、利用者の確認観点をREADME/SPEC/Release Notesへ追記する。 |
| 170 | [Spout同梱バイナリを単一exe配布に戻せるか調査する](investigation-only/170-investigate-single-exe-distribution-for-spout.md) | 完了 | 調査のみ | v0.1.8-a13で増えた `spout-capture.exe` / `SpoutLibrary.dll` を利用者から見て単一exeへ寄せられるか調査した。 |
| 172 | [player_local基準Poseの挙動を見直す](v0.1.8/172-clarify-or-redesign-player-local-basis-behavior.md) | 完了 | `v0.1.8` | `player_local` の basis source を AvatarBeacon の `avatar_osc` 既定と `manual` フォールバックで整理し、README/verification docsを更新した。 |
| 174 | [Release本文先頭にダウンロード導線を置き、通常配布zipへ必要物をまとめる](v0.1.8/174-release-top-download-zip-with-avatar-package.md) | 完了 | `v0.1.8` | Release本文の最上部に通常配布zipへの大きなダウンロードリンクを置き、手動添付するアバターギミック用unitypackageも利用者向けに案内する。 |
| 186 | [完了済みissueをclosedディレクトリへ移動する](maintenance-only/186-move-completed-issues-to-closed.md) | 完了 | 整理のみ | `issues/README.md` で完了扱いのissueを `issues/closed/` へ移動し、一覧リンクを更新する。 |
| 190 | [カメラ撮影機能終了時に一時OSC状態を確実に解除する](v0.1.8-a31/190-make-camera-osc-reset-streaming-compat.md) | 完了 | `v0.1.8-a31` | カメラOSCリセットと撮影終了時の `/usercamera/Streaming=false` などを通常のStream制御と同じbool+numeric互換送信に揃え、成否にかかわらず一時OSC状態を解除する。 |
| 193 | [a30 Release workflowのWindowsテストcleanup失敗を安定化する](v0.1.8-a34/193-stabilize-a30-windows-release-test-cleanup.md) | 完了 | `v0.1.8-a34` | Windows CIでOSC受信器テスト終了時の一時ディレクトリcleanupが診断ログ/UDP listener停止と競合しないよう、receiver終了待ちを追加した。 |
| 194 | [CIのWails Build application工程が長時間完了しない](v0.1.8-a34/194-investigate-ci-wails-build-hang.md) | 完了 | `v0.1.8-a34` | Windows `wails build` がbindings生成で止まる問題に対し、bindings生成skip、WebView2 strategy、mod同期抑止、timeout/verboseを追加した。 |
| 198 | [world ID を Avatar OSC で送れるか調査する](v0.1.8/198-world-id-avatar-osc-feasibility-investigation.md) | 断念 | `v0.1.8` | AvatarBeacon単体で現在world ID文字列をOSC送信する公式手段が確認できないため断念し、VRChat output log由来のworld/instance取得を代替とする。 |
| 203 | [Spout2ライセンス表記の配布方法を確認する](investigation-only/203-confirm-spout2-license-distribution.md) | 完了 | 調査のみ | `Spout2-LICENSE.txt` は公開zipへ同梱継続が妥当で、アプリ内OSS表示だけで代替するにはBSD-2-Clause本文表示が必要。 |
| 204 | [アプリ内OSSライセンス表示でライセンス本文を確認できるようにする](unreleased/204-show-full-oss-license-text-in-app.md) | 完了 | 未定 | Spout2のBSD-2-Clause本文をアプリ内OSSライセンス画面で確認できるようにし、zip内個別ライセンスファイル代替の前提を整えた。 |
| 205 | [Release note書式と公開Assetを整理する](v0.1.8/205-simplify-release-notes-and-assets.md) | 完了 | `v0.1.8` | Release noteを指定書式へ揃え、通常zip内容とGitHub Release添付Assetを必要な4種類へ絞った。 |
| 205 | [Release build metadata assetの意図を確認する](investigation-only/205-investigate-release-build-metadata-asset.md) | 完了 | 調査のみ | Release assetに添付されるbuild metadata JSONの生成箇所、内容、添付意図を確認した。 |
| 206 | [自動撮影機能を使用するための条件を整理する](investigation-only/206-investigate-auto-capture-user-requirements.md) | 完了 | 調査のみ | 自動撮影を使うための前提条件、方式別条件、構図基準、任意機能の条件を整理した。 |
| 207 | [設定保存時に変更後の設定値を診断ログへ出力する](unreleased/207-log-config-values-after-settings-save.md) | 完了 | 未定 | 設定保存成功/失敗時に診断ログを出し、Webhook設定有無とフォールバック状態を秘密情報なしで追えるようにした。 |
| 209 | [構図カードの未キャリブレーション表示を分かりやすくする](unreleased/209-clarify-camera-view-calibration-label.md) | 完了 | 未定 | 構図カードから内部状態のキャリブレーション表示を削除し、撮影対象と座標系だけを表示する。 |
| 210 | [構図設定の拡大率が常に最低値になる](unreleased/210-fix-auto-capture-view-zoom-range.md) | 完了 | 未定 | 構図ZoomのUI範囲、初期値、既存設定補正をUser Camera Zoomの扱いに合わせた。 |
| 211 | [ローカルアンカー配置済みカメラを使うフォールバックモードを追加する](unreleased/211-add-preplaced-local-anchor-fallback-mode.md) | 完了 | 未定 | VRChat内でローカルアンカー配置済みのカメラを使い、ClipForVRChatは撮影だけ操作するフォールバックを追加した。 |
| 212 | [Stream方式撮影失敗をSpout helper録画から切り分ける](v0.1.8-a43/212-investigate-stream-capture-spout-helper-debug-recording.md) | 完了 | `v0.1.8-a43` | Spout helper単体と本体経由の実機確認で、`IsUpdated()` 未処理による透明フレームを修正し、Stream方式撮影成功を確認した。 |
| 213 | [完了済みissueを対応バージョンごとのフォルダに整理する](maintenance-only/213-organize-closed-issues-by-version.md) | 完了 | 整理のみ | 完了済みissueを対応バージョン別フォルダへ移動し、READMEリンクと運用ルールを更新した。 |
| 214 | [Spout senderが空のときStreamingをOFF/ONして再確認する](v0.1.8-a43/214-retry-spout-sender-after-streaming-toggle.md) | 完了 | `v0.1.8-a43` | sender一覧が空のときOSCでStreamingをOFF/ONして再確認し、本体経由でsender復帰からcapture成功まで確認した。 |
| 215 | [alpha/beta/rcを分けたバージョニング規則を定義する](maintenance-only/215-define-alpha-beta-rc-versioning.md) | 完了 | 整理のみ | RCを正式リリース直前確認に限定し、alpha/beta/rc/正式版のタグ規則、CI/CD対象、勝手なスコープ縮小禁止を明文化した。 |
| 216 | [v0.1.8既存RCタグをalpha/beta扱いへ整理する](maintenance-only/216-reclassify-v018-rc-tags-to-alpha-beta.md) | 完了 | 整理のみ | 旧 `v0.1.8-rc1..rc43` を `v0.1.8-a1..a43`、旧 `v0.1.8-rc44` を `v0.1.8-b1` としてタグ・記録・分類を整理した。 |
| 217 | [v0.1.8-a44とv0.1.8-b2でRelease workflowを確認する](maintenance-only/217-confirm-a43-b1-release-creation.md) | 完了 | 整理のみ | CI設定変更を含む `v0.1.8-a44` と `v0.1.8-b2` でalpha/betaタグのRelease作成を確認した。 |
| 218 | [直近コミットメッセージがConventional Commitsから外れた原因を調査する](maintenance-only/218-investigate-conventional-commit-regression.md) | 完了 | 整理のみ | 直近130コミットの件名をConventional Commits形式へ整理し、再発防止ルールと検査スクリプトを追加した。 |
| 219 | [構図設定の右下に＋/－ボタンを追加する](unreleased/219-add-plus-minus-composition-buttons.md) | 完了 | 未定 | 構図設定リストの右下に、構図追加と末尾削除用の `＋` / `－` ボタンを追加した。 |
| 220 | [現在Pose取得で拡大率が保存されない](v0.1.8-b3/220-save-current-camera-pose-captures-zoom.md) | 完了 | `v0.1.8-b3` | 現在Pose取得時に直近のUser Camera Zoomも構図へ保存し、45以外の拡大率を反映できるようにした。 |
| 221 | [詳細設定画面でも保存/閉じるボタンを表示する](v0.1.8-b4/221-show-settings-actions-in-detail-view.md) | 完了 | `v0.1.8-b4` | 自動撮影の詳細設定画面でも設定画面上部の保存/閉じる操作を表示するようにした。 |
| 222 | [フォールバックモード自動切替設定の解釈誤り](v0.1.8-b4/222-default-enable-preplaced-local-anchor-toggle.md) | 完了 | `v0.1.8-b4` | フォールバックモード既定ONと自動切替廃止は解釈誤りとして、#235で自動ON/OFF設定へ修正する。 |
| 223 | [保存/閉じるボタンを上部ナビ右へ移動する](unreleased/223-move-settings-actions-to-header-nav.md) | 完了 | 未定 | 設定画面中の保存/閉じる操作を、設定/使い方/情報ボタンの右側へ移動した。 |
| 224 | [終了時にSpout helper展開キャッシュを削除する](unreleased/224-clean-spout-helper-cache-on-exit.md) | 完了 | 未定 | `%LOCALAPPDATA%\ClipForVRChat\spout-helper` に展開された内蔵Spout helperの管理キャッシュをアプリ終了時に削除するようにした。 |
| 225 | [自動撮影用Webhook URLをDiscord投稿タブへ移動する](unreleased/225-move-autocapture-webhook-to-discord-tab.md) | 完了 | 未定 | 自動撮影用Webhook URLをDiscord投稿タブへ移し、投稿先設定を同じ場所にまとめた。 |
| 226 | [フォールバックモードを構図設定の上へ移動する](unreleased/226-move-fallback-mode-above-composition.md) | 完了 | 未定 | 自動撮影タブ概要で構図設定の上にフォールバックモードを表示し、詳細画面に入らず切り替えられるようにした。 |
| 227 | [OSC送信ログを分離し任意OSC送信欄を追加する](unreleased/227-split-osc-send-log-and-debug-send.md) | 完了 | 未定 | OSCタブで受信ログと送信ログを分離し、デバッグ用に任意OSCを送信できる入力欄を追加した。 |
| 229 | [ウィンドウ終了時にも未保存設定の保存確認を出す](unreleased/229-confirm-unsaved-settings-on-window-close.md) | 完了 | 未定 | 設定画面で未保存変更がある状態の×ボタンや複数起動による既存終了で、保存確認ダイアログを出すようにした。 |
| 230 | [保存せず終了した設定変更を一時変更として復元する](unreleased/230-restore-unsaved-settings-draft-from-single-instance.md) | 完了 | 未定 | 未保存変更をsingle-instance配下へ一時保存し、次回起動時に保存せず一時状態として復元するようにした。 |
| 231 | [未保存変更がある設定項目名をハイライトする](unreleased/231-highlight-unsaved-settings-fields.md) | 完了 | 未定 | 保存前に変更された設定項目名を設定画面上部へハイライト表示するようにした。 |
| 232 | [GitHub Release本文からバージョン見出しを除外する](v0.1.8/232-strip-release-body-version-heading.md) | 完了 | `v0.1.8` | Release本文生成時に先頭の単独バージョン見出しを除外し、タイトルとの二重表記を防ぐ。 |
| 233 | [OSCデバッグ送信が送信ログに表示されない](unreleased/233-fix-debug-osc-send-log.md) | 完了 | 未定 | OSC受信/転送ログとOSC送信ログを別バッファに分け、受信ログ増加で送信ログが消えないようにした。 |
| 234 | [未保存変更の項目名が大分類だけで分かりにくい](unreleased/234-show-specific-unsaved-setting-labels.md) | 完了 | 未定 | 未保存変更表示でOSCホストやOSC送信ポートなどの具体的な設定項目名を表示するようにした。 |
| 235 | [自動フォールバック切替設定を追加する](unreleased/235-add-auto-fallback-mode-toggle.md) | 完了 | 未定 | フォールバックモードの自動ON/OFF制御を別フラグで追加し、どちらも既定OFFにした。 |
| 236 | [closed issue一覧の昇順とissue記録形式を修正する](unreleased/236-sort-closed-issue-index-and-preserve-user-instructions.md) | 完了 | 未定 | closed issue一覧を昇順に直し、昇順チェックと原文付きissue記録ルールを追加した。 |
| 237 | [構図削除後に初期3構図へ戻すとグレーアウトする](unreleased/237-fix-reset-default-views-disabled-after-delete.md) | 完了 | 未定 | 初期3構図へ戻した後も、フォールバックOFFなら構図が操作できる状態になるようにした。 |
| 238 | [未保存変更のある設定項目を行としてハイライトする](unreleased/238-highlight-unsaved-setting-rows.md) | 完了 | 未定 | 未保存変更のある主要な設定行と構図カードを画面上でハイライトするようにした。 |
| 239 | [AvatarBeaconにバージョン情報を埋め込みOSCで通知する](unreleased/239-embed-avatarbeacon-version-and-send-osc-version.md) | 完了 | 未定 | AvatarBeaconのPrefab/source内へ固定名でバージョン情報を埋め込み、OSC送信時に1回だけバージョン番号を送信し、CI/Releaseで自動更新する。 |
| 240 | [b6でGUIが表示されないVue compiler errorを調査する](unreleased/240-investigate-b6-gui-vue-compiler-error.md) | 完了 | 未定 | `v0.1.8-b6` のVue runtime templateタグ不整合を修正し、compiler検査をCI/Releaseに追加した。 |
| 242 | [設定内のカメラ姿勢Pose表記を日本語へ変更する](unreleased/242-localize-settings-camera-pose-labels.md) | 完了 | 未定 | 設定UI内のカメラ姿勢を表す `Pose` を日本語表記へ置き換えた。 |
| 243 | [`/usercamera/Close` が送信されないようにする](v0.1.8-b10/243-stop-sending-usercamera-close.md) | 完了 | `v0.1.8-b10` | カメラOSCリセットと旧 `closeCameraAfterBatch` 設定値由来の `/usercamera/Close` 送信を止めた。 |
| 246 | [OSCデバッグ送信説明、AvatarBeaconログフィルター、b7 basis受信判定を修正する](unreleased/246-fix-osc-debug-avatarbeacon-filter-and-b7-basis-status.md) | 完了 | 未定 | OSCデバッグ送信説明、AvatarBeacon以外ログフィルター、README日本語化、version GameObjectのEditorOnly化、b7 main受信判定とPose受信ログ量を修正した。 |
| 247 | [その他タブにPC起動時自動起動トグルを追加する](unreleased/247-add-startup-shortcut-toggle-to-other-settings.md) | 完了 | 未定 | 更新タブをその他へ変更し、Startupフォルダの固定名ショートカットでPC起動時自動起動をON/OFFできるようにした。 |
| 248 | [未保存変更の対象設定ラベルをハイライトする](unreleased/248-highlight-unsaved-setting-labels.md) | 完了 | 未定 | 未保存変更チップに対応する設定行とタブボタンをハイライトし、変更箇所を分かりやすくした。 |
| 249 | [AvatarBeaconバージョンOSC受信をログに記録する](unreleased/249-log-avatarbeacon-version-osc-receive.md) | 完了 | 未定 | AvatarBeacon version OSCを受信したら診断ログへ記録し、導入済みギミックのバージョンを切り分けやすくした。 |
| 250 | [初回起動直後に未保存変更が表示される](unreleased/250-fix-initial-unsaved-settings-after-first-launch.md) | 完了 | 未定 | configなし初回起動で既定値のOSC転送先が差分扱いになる問題を修正し、未保存変更の具体項目と行ハイライトを改善した。 |
| 251 | [AvatarBeaconフォールバック判定とOSCデバッグUIを調整する](unreleased/251-adjust-avatarbeacon-fallback-and-osc-debug-ui.md) | 完了 | 未定 | OSCデバッグ既定入力、AvatarBeaconのみログフィルター、フォールバック自動ON待機時間、構図値丸め、保存ボタン挙動を調整した。 |
| 252 | [起動中にexeへドラッグしたファイルを既存インスタンスへ渡す](unreleased/252-forward-exe-drop-paths-to-running-instance.md) | 完了 | 未定 | 単一起動中に2つ目のexeへ渡されたファイル引数を既存インスタンスへ転送し、通常のドロップ処理で扱うようにした。 |
| 253 | [不具合報告zipとgpgのファイル名に添付可否を付ける](unreleased/253-prefix-diagnostic-package-filenames.md) | 完了 | 未定 | 不具合報告用の暗号化前zipと添付用gpgのファイル名先頭に、添付可否が分かる文言を付けるようにした。 |
| 254 | [不具合報告についての説明文を更新する](unreleased/254-rewrite-diagnostic-report-description.md) | 完了 | 未定 | 情報タブの不具合報告説明を、zip/gpgの役割、添付対象、公開鍵による手動暗号化案内が分かる文面へ更新した。 |
| 255 | [OSC送受信ログを分けて保持し不具合報告zipに含める](v0.1.8-b10/255-persist-osc-send-receive-logs-in-diagnostic-zip.md) | 完了 | `v0.1.8-b10` | `logs/osc_send.jsonl` と `logs/osc_recieve.jsonl` に最新1000件のOSC送受信ログをJSON Linesで保存し、診断zipへ含めるようにした。 |
| 256 | [v0.1.8-b10 betaを作成する](maintenance-only/256-create-v018-b10-beta.md) | 完了 | 整理のみ | #243 と #255 を含む `v0.1.8-b10` betaタグを作成し、Release workflowとGitHub Releaseを確認した。 |
| 257 | [b10で3枚撮影想定なのに6枚投稿される](unreleased/257-fix-b10-autocapture-discord-duplicate-posts.md) | 完了 | 未定 | VRChat写真自動投稿が自動撮影出力ディレクトリを拾って二重投稿する問題を修正した。 |
| 258 | [保存ボタンで設定タブとスクロール位置を移動しない](unreleased/258-preserve-settings-tab-and-scroll-on-save.md) | 完了 | 未定 | 設定保存後も表示中の設定タブ、詳細画面、スクロール位置を維持するようにした。 |
| 259 | [不具合報告についての説明文の改行を維持する](unreleased/259-preserve-diagnostic-report-description-line-breaks.md) | 完了 | 未定 | 情報タブの不具合報告説明を指定文面へ更新し、表示上の改行を維持するようにした。 |
| 260 | [OSSライセンス表記を監査し本文を表示する](unreleased/260-audit-oss-license-notices-and-show-full-text.md) | 完了 | 未定 | OSS表示を実依存に合わせて監査し、不足・誤表記・テスト専用項目を整理して全項目のライセンス本文を表示できるようにした。 |
| 261 | [AvatarBeaconライセンスをOSS画面の最下部に分けて表示する](unreleased/261-show-avatarbeacon-license-separately.md) | 完了 | 未定 | AvatarBeacon関連ライセンスを通常OSS一覧とは別セクションとしてOSSライセンス画面の最下部に表示するようにした。 |
| 262 | [処理タブを機能と自動撮影の間に移動する](unreleased/262-move-process-settings-tab-before-autocapture.md) | 完了 | 未定 | 設定画面のタブ順を、機能、処理、自動撮影、OSC、Discord投稿、その他に変更した。 |
| 263 | [Discord投稿タブを処理と自動撮影の間に移動する](unreleased/263-move-discord-settings-tab-between-process-and-autocapture.md) | 完了 | 未定 | 設定画面のタブ順を、機能、処理、Discord投稿、自動撮影、OSC、その他に変更した。 |
| 264 | [v0.1.8-b11 betaを作成する](maintenance-only/264-create-v018-b11-beta.md) | 完了 | 整理のみ | 現在のdevelopを `v0.1.8-b11` betaとしてタグ付けし、Release workflowと配布Assetを確認した。 |
| 265 | [AvatarBeaconを専用リポジトリへ配置する](avatarbeacon-v0.0.1/265-extract-avatarbeacon-to-dedicated-repository.md) | 完了 | `AvatarBeacon v0.0.1` | 既存AvatarBeacon元ファイルを専用リポジトリへ初期配置し、CI設定、push、最新main CI確認を行った。 |
| 266 | [AvatarBeacon READMEを作り直し、ClipForVRChat側をsubmodule化する](unreleased/266-rewrite-avatarbeacon-readme-and-use-submodule.md) | 完了 | 未定 | AvatarBeacon READMEを利用者向けに全面改稿し、ClipForVRChat側のAvatarBeaconをsubmodule参照へ置き換えた。 |
