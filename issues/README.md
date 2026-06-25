# Issues

このディレクトリは、実装タスクや調査タスクを Markdown で管理するためのものです。

状態はリポジトリ上の実装状況をもとに整理しています。実機確認が必要なものは「要確認」としています。

| No. | Issue | 状態 | 対応バージョン | 概要 |
| --- | --- | --- | --- | --- |
| 001 | [アイコン品質と複数サイズ生成](001-icon-quality-and-sizes.md) | 完了 | `v0.1.0` | アプリアイコン生成と複数サイズの品質改善。 |
| 002 | [Release zipの内容とバージョン付きファイル名](002-release-zip-contents-and-versioned-name.md) | 完了 | `v0.1.0` | 配布zipの内容整理、バージョン情報の埋め込み、Release workflow 整備。 |
| 003 | [メインウィンドウUIと情報表示](003-main-window-ui-and-about.md) | 完了 | `v0.1.0` | メイン画面、情報表示、アプリ情報取得の追加。 |
| 004 | [初回設定フロー](004-initial-config-flow.md) | 完了 | `v0.1.0` | GUIアプリとして設定を扱える初期フローを実装。 |
| 005 | [ウィンドウへのドラッグ&ドロップ](005-window-drag-and-drop.md) | 完了 | `v0.1.0` | 画像/config.json のドラッグ&ドロップ処理と画面全体のドロップ対応。 |
| 006 | [クリップボード画像の保存ファイル名に時刻を含める](006-clipboard-output-filename.md) | 完了 | `v0.1.0` | クリップボード入力の保存ファイル名を時刻付きに変更。 |
| 007 | [クリップボードのスクリーンショット画像が崩れる](007-clipboard-screenshot-corruption.md) | 完了 | `21cc973` / `v0.1.1` | Win+Shift+S 由来画像の崩れ対策として Windows 登録PNG読み取りを優先。 |
| 008 | [UI仕様の明文化](008-ui-specification.md) | 完了 | `v0.1.0` | UI仕様と運用ルールをドキュメント化。 |
| 009 | [ソース配置を src 配下へ移動する](009-move-source-under-src.md) | 完了 | `v0.1.0` | Wails/Go/frontend ソースを `src` 配下へ整理。 |
| 010 | [設定画面・結果画面・情報画面のUI改善](010-settings-results-about-ui-improvements.md) | 完了 | `v0.1.0` | 設定/結果/情報画面のレイアウト、進捗、クリア、未保存確認を改善。 |
| 011 | [使い方画面とWebhook発行案内](011-help-screen-and-webhook-guide.md) | 完了 | `v0.1.0` | 使い方画面と Discord Webhook 公式案内リンクを追加。 |
| 012 | [画像URL履歴とDiscord削除確認画面](012-image-history-and-discord-delete-review.md) | 完了 | `v0.1.0`, `v0.1.2` | 履歴保存、クリア済み表示、Discord削除確認画面、Ctrl/Shift選択、矩形範囲選択を実装。 |
| 013 | [ユーザーフレンドリーなREADMEと設定画面改善](013-user-friendly-readme-and-settings-ui.md) | 完了 | `v0.1.0` | README、設定画面、WebHook説明、出力先選択などを改善。 |
| 014 | [VRChat写真の自動検知とDiscord投稿](014-vrchat-photo-auto-post.md) | 完了 | `v0.1.0` | VRChat写真フォルダの定期スキャンと自動Discord投稿を実装。 |
| 015 | [セキュリティチェック報告書の作成](015-security-review-report.md) | 完了 | `v0.1.0` | `reports/2026-06-21/security-review.md` を作成し、リスクと推奨対応を整理。 |
| 016 | [Webhook URL と履歴 URL の検証強化](016-validate-webhook-and-history-urls.md) | 完了 | `v0.1.0` | Discord Webhook URL と履歴画像URLの検証を強化。 |
| 017 | [設定・履歴ファイルの権限と排他制御](017-harden-local-secret-storage-and-history-locking.md) | 完了 | `v0.1.0` | 設定/履歴ファイルの権限、履歴更新の排他制御を強化。 |
| 018 | [画像入力のサイズ上限とピクセル数上限](018-limit-image-input-resource-usage.md) | 完了 | `v0.1.0` | 入力画像のファイルサイズ上限とピクセル数上限を追加。 |
| 019 | [VRChat写真自動投稿の走査・処理件数制限](019-limit-auto-photo-scanning.md) | 完了 | `v0.1.0` | 自動投稿の走査量と処理件数を制限。 |
| 020 | [CI/Release のセキュリティチェック追加](020-add-security-checks-to-ci-release.md) | 完了 | `v0.1.0` | CI/Release に監査や脆弱性チェックを追加。 |
| 021 | [開発用アイコン生成ツールの gosec 指摘対応](021-address-gosec-tooling-findings.md) | 完了 | `v0.1.0` | gosec 指摘を受けた開発用ツールの処理を修正。 |
| 022 | [アプリの複数起動制限](022-prevent-multiple-app-instances.md) | 完了 | `v0.1.0` | 多重起動防止を実装。 |
| 023 | [Goテストカバレッジの拡充](023-expand-go-test-coverage.md) | 完了 | `v0.1.0` | appcore とアプリ主要処理の Go テストを拡充。 |
| 024 | [Discord Webhook URLエラーの案内改善](024-improve-discord-webhook-error-message.md) | 完了 | `v0.1.0` | 空/不正な Webhook URL のエラー文言をユーザー向けに改善。 |
| 025 | [設定画面の未保存変更確認](025-confirm-unsaved-settings-navigation.md) | 完了 | `f17eea0` / `v0.1.0` | 設定画面を離れる前に保存/破棄/キャンセルを選べる確認ダイアログを追加。 |
| 026 | [v0.1.1 不具合修正](026-fix-v0.1.1-regressions.md) | 完了 | `21cc973` / `v0.1.1`, 028で導線更新 | 履歴画面への導線、Win+Shift+S 画像崩れを修正。 |
| 027 | [画像履歴のマウス矩形範囲選択](027-history-drag-selection.md) | 完了 | `v0.1.2` | 履歴画面でマウスドラッグによる矩形範囲選択を追加。 |
| 028 | [履歴画面への長押し導線をボタンへ変更](028-replace-history-long-press-with-button.md) | 完了 | `v0.1.2` | クリア長押しを廃止し、履歴ボタンと各ボタンの説明 tooltip を追加。 |
| 029 | [画像履歴の全選択とピン止め](029-history-select-all-and-pin.md) | 完了 | `v0.1.2` | Ctrl+A/全選択ボタンと、削除対象外にするピン止めを追加。 |
| 030 | [Discord公式ヘルプリンクが開けない](030-fix-discord-help-link.md) | 完了 | `v0.1.2` | Discord公式ヘルプURLからテキストフラグメントを削除。 |
| 031 | [Discord削除済み履歴の削除条件とoutput削除](031-purge-discord-deleted-history.md) | 完了 | `v0.1.3` | Discord削除済みフラグで履歴を削除し、設定に応じてoutputも削除。 |
| 032 | [QRコードURLの検出とDiscord投稿](032-detect-qr-url-and-post.md) | 完了 | `v0.1.4` | 画像内QRコードのURLを検出し、Discord投稿本文と結果画面に表示する。 |
| 033 | [WSLからのWindowsビルドスクリプト](033-wsl-windows-build-script.md) | 完了 | `v0.1.4` | WSLからWindows向けexeをコミットハッシュ版としてローカルビルドするスクリプトを追加。 |
| 034 | [GitHub ActionsビルドのバージョンにコミットIDを含める](034-release-version-includes-commit.md) | 完了 | `v0.1.4` | Release workflow で表示バージョンをリリース番号.コミットIDにする。 |
| 035 | [処理診断ログの追加](035-add-processing-diagnostic-log.md) | 完了 | `v0.1.4` | QR検出などの原因調査用に処理時の診断ログを出力する。 |
| 037 | [Release exe の外部署名ファイル](037-release-detached-gpg-signature.md) | 完了 | `v0.1.4` | ReleaseでexeのPGP外部署名ファイルを生成して添付する。 |
| 038 | [情報画面に公式配布元とPGP確認方法を追加](038-about-official-distribution-and-pgp.md) | 完了 | `v0.1.4` | 公式配布元とexe署名による改竄確認方法を情報画面に記載する。 |
| 039 | [バージョンとリビジョンを分けてビルドに埋め込む](039-separate-version-and-revision-build-info.md) | 完了 | `v0.1.4` | ビルド時にバージョンとリビジョンを別々に埋め込む。 |
| 040 | [GitHub Release のアップデート通知](040-check-github-release-updates.md) | 完了 | `v0.1.4` | GitHub Releases の最新Releaseを確認し、更新があればUI内に通知する。 |
| 041 | [Win+Shift+S スクリーンショット自動処理](041-screenshot-auto-post.md) | 完了 | `v0.1.5` | Screenshotsフォルダを定期スキャンし、Win+Shift+Sで保存された画像を自動処理する。 |
| 042 | [Windows exe のプロパティ表示改善](042-windows-exe-version-info.md) | 完了 | `v0.1.6` | Windowsのファイルプロパティに製品名、説明、著作権、バージョンを追加。 |
| 043 | [v0.1.6 リリースノート作成](043-release-notes-v0.1.6.md) | 完了 | `v0.1.6` | `RELEASE_NOTES.md` にv0.1.6の更新内容と配布URLを追加。 |
| 044 | [v0.1.6 Release workflow のVersionInfo生成失敗修正](044-fix-v0.1.6-release-workflow-version-info.md) | 完了 | `v0.1.6` | Release workflowのVersionInfo生成先ディレクトリ作成漏れを修正。 |
| 045 | [Release本文と成果物の再発防止](045-fix-release-assets-and-notes.md) | 完了 | `v0.1.6` | Release本文と成果物一覧が仕様通りになるようworkflowと検査を修正。 |
| 046 | [初回設定を保存せず終了した場合に次回も初期設定画面を出す](046-fix-initial-config-save-timing.md) | 完了 | `v0.1.6` | 初回設定画面を開いただけで `config.json` が作成されないようにする。 |
| 047 | [CLIでversion/help引数に対応する](047-cli-version-help-args.md) | 完了 | `v0.1.6` | `go-arg` で `--version` / `--help` の早期終了に対応する。 |
| 048 | [Windows GUI exeのCLI出力をPowerShellに表示する](048-cli-output-from-windows-gui-exe.md) | 完了 | `v0.1.6` | GUIサブシステムのexeでもCLI出力を親コンソールへ表示する。 |
| 049 | [CLIヘルプをWindowsコンソールで文字化けさせない](049-use-wide-console-output-for-cli.md) | 完了 | `v0.1.6` | `WriteConsoleW` でCLI出力し、cmdのコードページ差による文字化けを避ける。 |
| 050 | [プロダクト全体の問題点チェック報告書を作成する](050-product-issue-review-report.md) | 完了 | `v0.1.6` | `reports/2026-06-24/product-issue-report.md` を作成し、Release/UX/履歴/設定/OSS表示などの残課題を整理。 |
| 051 | [監査報告書を日付付きで専用ディレクトリへ移動する](051-date-audit-reports-directory.md) | 完了 | `v0.1.6` | 監査報告書を `reports/` 配下の日付ディレクトリへ整理。 |
| 052 | [Release署名ファイルの説明を実際の成果物に合わせる](052-align-release-signing-docs.md) | 完了 | `v0.1.6` | README/About/SPEC のPGP確認手順をRelease成果物に合わせる。 |
| 053 | [Discord削除の部分成功を履歴へ保存する](053-persist-partial-discord-delete-results.md) | 完了 | `v0.1.6` | 複数削除の一部失敗時も成功済み削除状態を保存する。 |
| 054 | [config読み込み失敗時にアクティブ設定パスを変更しない](054-keep-config-path-on-load-error.md) | 完了 | `v0.1.6` | 不正configの読み込み失敗時に既存の設定パスを維持する。 |
| 055 | [URL自動コピー失敗をユーザーへ表示する](055-report-copy-single-url-failure.md) | 完了 | `v0.1.6` | URL取得成功とクリップボードコピー失敗を分けて表示する。 |
| 056 | [OSSライセンス表示を依存関係に合わせて更新する](056-complete-oss-license-list.md) | 完了 | `v0.1.6` | direct dependency のライセンス表示漏れを補う。 |
| 057 | [自動投稿の監視フォルダ異常を通知する](057-report-auto-post-watch-diagnostics.md) | 完了 | `v0.1.6` | 監視フォルダ異常やスキャン上限到達をUIまたは診断ログへ出す。 |
| 058 | [ローカルWindowsビルドのバージョンを指定できるようにする](058-allow-local-build-version-override.md) | 完了 | `v0.1.6` | ローカルビルド時にリリース候補バージョンを指定可能にする。 |
| 059 | [診断ログのgosec指摘を整理する](059-resolve-diagnostic-log-gosec-finding.md) | 完了 | `v0.1.6` | 診断ログの可変パス指摘を設計上明確にする。 |
| 060 | [issue一覧のv0.1.6対応状況を整理する](060-update-issue-index-release-status.md) | 完了 | `v0.1.6` | issue一覧の未リリース表記と欠番を整理。 |
| 061 | [リリース時にissue一覧へ対応バージョンを記録する](061-record-release-version-in-issue-index.md) | 完了 | `v0.1.6` | リリース時の対応バージョン記録ルールを追加し、既存一覧を補完。 |
| 062 | [issue一覧の対応バージョン列名を整理する](062-rename-issue-index-version-column.md) | 完了 | `v0.1.6` | issue一覧の列名を対応バージョン記録の運用に合わせる。 |
| 063 | [v0.1.6履歴のコミットメッセージと署名を整理する](063-rewrite-v016-history-angular-signed.md) | 完了 | `v0.1.6` | v0.1.6向け履歴のコミットメッセージ、署名、相殺revertを整理。 |
| 064 | [master以外のブランチCIと非正式タグのdraft release対応](064-run-ci-on-branches-and-draft-nonstable-tags.md) | 完了 | `v0.1.6` | ブランチCI対象と非正式タグのdraft Release作成を整理。 |
| 065 | [開発ブランチ運用とRCリリースを整備する](065-dev-flow-and-rc-release.md) | 完了 | `v0.1.6` | 作業ブランチ前提の開発フローと `vX.Y.Z-rcW` prerelease 発行を整理。 |
| 066 | [GitHub ActionsのNode 20警告とGo cache警告を解消する](066-update-actions-runtime-and-go-cache.md) | 完了 | `v0.1.6` | Actions実行系とGo cache設定を更新し、CI/Releaseの警告を解消。 |
| 067 | [todo.mdを直近作業用に整理する](067-maintain-short-term-todo.md) | 完了 | `v0.1.6` | `todo.md` を短期チェックリストとして整理し、運用ルールを明文化。 |
| 068 | [GitHub Actions CIを高速化する](068-speed-up-github-actions-ci.md) | 完了 | `v0.1.7` | Wails CLI cacheを追加し、CI/Releaseのインストール時間短縮を試行。 |
| 069 | [更新通知の開き先選択と通知設定](069-update-notification-destination-and-settings.md) | 完了 | `v0.1.7` | 更新通知からGitHub/BOOTHを選んで開き、更新確認と通知を設定でON/OFFできるようにする。 |
| 070 | [設定画面をカテゴリ別タブに分ける](070-tabbed-settings-page.md) | 完了 | `v0.1.7` | 設定画面上部にカテゴリタブを追加し、設定項目を分類して表示する。 |
| 071 | [Discord投稿OFFでも自動投稿でDiscordへ投稿される](071-prevent-discord-post-when-disabled.md) | 完了 | `v0.1.7` | Discord投稿OFF時に自動投稿経路からDiscordへ投稿されないようにする。 |
| 072 | [設定画面仕様を独立ファイルへ分離する](072-split-settings-screen-spec.md) | 完了 | `v0.1.7` | 設定画面仕様を専用Markdownへ分離し、カテゴリと依存関係を整理する。 |
| 073 | [設定画面を機能/処理/Webhook/更新へ再編する](073-reorganize-settings-screen-categories.md) | 完了 | `v0.1.7` | 設定画面を利用者視点のカテゴリへ再編する。 |
| 074 | [結果画像のクリック領域をURLコピーと保存先表示に分ける](074-result-card-split-actions.md) | 完了 | `v0.1.7` | 結果画像の上半分でURLコピー、下半分で保存先表示を行えるようにする。 |
| 075 | [履歴画面をDiscord/ローカル/履歴の状態別操作に作り直す](075-rebuild-history-screen.md) | 完了 | `v0.1.7` | 履歴画面を状態表示とDiscord/ローカル/履歴削除の独立操作へ作り直す。 |
| 076 | [履歴と結果表示を実際に行った処理に合わせる](076-history-and-result-display-follow-actual-work.md) | 完了 | `v0.1.7` | 結果画面と履歴画面を実際に行った処理内容に合わせて表示する。 |
| 077 | [処理結果がない場合に理由を表示する](077-explain-no-result-processing-message.md) | 完了 | `v0.1.7` | 明示的な処理で結果がない場合、設定により処理結果が出ない理由を表示する。 |
| 078 | [履歴のローカル保存パスをconfig基準で解決する](078-resolve-history-local-paths-relative-to-config.md) | 完了 | `v0.1.7` | 履歴の相対ローカル保存パスをconfig基準で解決し、Explorer表示と削除可否を修正する。 |
| 079 | [ユーザー操作を診断ログへ記録する](079-record-user-actions-in-diagnostic-log.md) | 完了 | `v0.1.7` | ボタンクリック、画面遷移、処理判断、処理結果表示を診断ログへ記録する。 |
| 080 | [暗号化診断パッケージ作成](080-create-encrypted-diagnostic-package.md) | 完了 | `v0.1.7` | 起動ログと診断パッケージ作成導線を追加。暗号化用公開鍵で診断パッケージを作成する。 |
| 081 | [診断ログを日付付きログフォルダへ出力する](081-write-diagnostic-logs-to-dated-log-directory.md) | 完了 | `v0.1.7` | 診断ログを `logs/YYYY-MM-DD.log` に出力し、診断パッケージがログフォルダを収集するようにする。 |
| 082 | [診断パッケージ暗号化鍵をpoppo@hato.lifeへ変更する](082-use-poppo-openpgp-key-for-diagnostics.md) | 完了 | `v0.1.7` | 診断パッケージ暗号化用の公開鍵を暗号化サブキー付きの `poppo@hato.life` に変更する。 |
| 083 | [不具合報告用データ生成UIを改善する](083-improve-diagnostic-data-generation-ui.md) | 完了 | `v0.1.7` | 情報画面の不具合報告項目に生成ボタンを配置し、生成中オーバーレイと完了後Explorer表示を追加する。 |
| 084 | [Explorerで生成ファイルを選択表示できない](084-fix-explorer-file-selection.md) | 完了 | `v0.1.7` | `explorer.exe /select,` ではなく Windows Shell API でファイル選択表示する。 |
| 085 | [診断データ復号後のzipが破損する](085-encrypt-diagnostic-package-as-binary.md) | 完了 | `v0.1.7` | 診断データをOpenPGP binary literal dataとして暗号化し、復号後zipの破損を防ぐ。 |
| 086 | [診断zip内のログとoutput構成を整理する](086-normalize-diagnostic-zip-layout-and-output-log.md) | 完了 | `v0.1.7` | zip内ログを `logs/` に統一し、output画像は含めず一覧を診断ログへ記録する。 |
| 087 | [zipファイル引数を公開鍵で暗号化する](087-encrypt-zip-from-cli-argument.md) | 完了 | `v0.1.7` | zip単体引数でUIを起動せず、同じ公開鍵で `<zip>.gpg` を生成する。 |
| 088 | [診断ログ実フォルダ名をlogsへ統一する](088-use-logs-directory-for-diagnostic-logs.md) | 完了 | `v0.1.7` | アプリ実フォルダの診断ログ出力先も `logs/YYYY-MM-DD.log` に統一する。 |
| 089 | [不具合報告用データの説明を情報画面に追加する](089-explain-diagnostic-data-in-about.md) | 完了 | `v0.1.7` | 情報画面に含まれる情報、暗号化方法、利用目的の説明を追加する。 |
| 090 | [不具合報告用データを段階作成しパスを置換する](090-stage-and-sanitize-diagnostic-data.md) | 完了 | `v0.1.7` | `diagnostics/<timestamp>/` に確認用zipと暗号化zipを作成し、テキスト内パスを環境変数表記へ置換する。 |
| 091 | [セキュリティ監査報告書を作成する](091-security-audit-2026-06-25.md) | 完了 | `8e46ffd` | `reports/2026-06-25/security-audit-prompt.md` に基づき、現行リポジトリのセキュリティ監査報告書を作成した。 |
| 092 | [監査報告書の配置と説明を整理する](092-organize-reports-and-clarify-security-notes.md) | 完了 | `v0.1.7` | 監査報告書を日付ディレクトリへ移動し、診断zip、Release workflow、Windows ACLの説明を補足する。 |
| 093 | [人間が確認する必要がある作業を手順化する](093-human-verification-guide.md) | 完了 | `v0.1.7` | 監査後に人間が確認する作業、現時点の確認結果、判断基準をMarkdown化する。 |
| 094 | [設定画面の初期タブと初期値を調整する](094-settings-initial-tab-and-defaults.md) | 完了 | `v0.1.7` | 初回起動時の設定タブ選択とDiscord投稿/投稿URL自動コピーの初期値を調整する。 |
| 095 | [設定画面のDiscord投稿分類を整理する](095-settings-discord-post-tab-reclassification.md) | 完了 | `v0.1.7` | WebhookタブをDiscord投稿へ変更し、Discord投稿関連設定をまとめる。 |
| 096 | [develop版バージョンに親コミットIDを含める](096-develop-version-include-parent-commit.md) | 完了 | `v0.1.7` | develop版バージョン表記に親コミットIDを含めて追跡しやすくする。 |
| 097 | [govulncheckのGO-2026-4550を解消する](097-fix-circl-govulncheck-finding.md) | 完了 | `v0.1.7` | `cloudflare/circl` の既知脆弱性を解消し、govulncheckを成功させる。 |
| 098 | [診断データからWebhook URLとDiscord tokenを除外する](098-redact-diagnostic-secrets.md) | 完了 | `v0.1.7` | 確認用zipと暗号化zipからWebhook URLとDiscord tokenの生値を除外する。 |
| 099 | [Release workflowの権限を最小化する](099-harden-release-workflow-permissions.md) | 完了 | `v0.1.7` | Release作成権限を必要なjobへ限定する。 |
| 100 | [OpenURLで開けるURLを許可ホストに制限する](100-restrict-open-url-hosts.md) | 完了 | `v0.1.7` | アプリから開ける外部URLを信頼済みHTTPSホストへ限定する。 |
| 101 | [履歴のローカル削除対象を管理output配下に制限する](101-restrict-history-local-delete-to-output.md) | 完了 | `v0.1.7` | 履歴改ざん時にoutput外ファイルを削除できないようにする。 |
| 102 | [診断ログのローカルパスと秘密情報を抑制する](102-redact-diagnostic-log-paths.md) | 完了 | `v0.1.7` | 診断ログ出力時点でパスと秘密情報を可能な範囲で抑制する。 |
| 103 | [Windows実機でACLを確認する](103-check-windows-acl.md) | 要確認 | 未定 | Windows実機でconfig/history/logs/diagnosticsのACLを確認する。 |
| 104 | [ReleaseにSBOMまたはビルドメタデータを追加する](104-add-release-sbom-or-build-metadata.md) | 完了 | `v0.1.7` | Release成果物に依存関係やビルド環境の追跡情報を追加する。 |
| 105 | [Windows GoテストでWailsイベント送信ガードが効かない](105-fix-windows-test-wails-event-guard.md) | 完了 | `v0.1.7` | Windowsの `.test.exe` テストバイナリでもWailsイベント送信を抑制する。 |
| 106 | [Release NotesのダウンロードURLをMarkdownリンクにする](106-release-notes-download-links.md) | 完了 | `v0.1.7` | Release本文のダウンロード欄でファイル名をリンクテキストとして表示する。 |
| 107 | [Git Flow運用をAGENTS.mdへ明文化する](107-update-git-flow-agents.md) | 完了 | `v0.1.7` | `develop` 基準の通常開発と `master` 基準のリリースブランチ運用を明文化する。 |

## 状態の意味

| 状態 | 意味 |
| --- | --- |
| 完了 | 実装、テスト、またはドキュメント作成が完了している。 |
| 要対応 | 問題を整理済みで、実装またはドキュメント更新が必要。 |
| 要確認 | 修正は入っているが、対象環境や実機での再現確認が必要。 |
| 一部将来対応あり | 主要な受け入れ条件は満たしているが、明示的に将来対応へ回した項目が残っている。 |
