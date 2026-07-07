# 全Markdownドキュメントを現状仕様に合わせて更新する

## 指示

> md更新するチケット探して

> 228を進めたい 実装とリポジトリ内のmdに差がないか確認して 仕様変更を適用してほしい

## 文脈

リポジトリ内のMarkdown更新作業に着手するため、既存チケットの中から該当する管理単位を確認する。
実装やリリース運用、設定UI、AvatarBeacon配布物などの変更が継続しているため、Markdown側に古い説明が残っていないかを実装基準で確認する。

## 解釈

Markdown全般の更新は、新規チケットではなく既存の `#228` を使って進める。
実装ファイルは変更せず、README、仕様書、運用ドキュメント、検証手順、AvatarBeacon関連ドキュメントなどのMarkdownを現状仕様に合わせて更新する。

## 問題

実装やリリース運用の変更が進んでおり、リポジトリ内のMarkdownドキュメント全体に古い説明、重複、未反映の仕様が残っている可能性がある。

## 期待する挙動

README、仕様書、運用ドキュメント、issue関連ドキュメントなど、リポジトリ内のMarkdownを横断的に確認し、現状仕様とリリース運用に合う内容へ更新する。

## 受け入れ条件

- [x] README、SPEC、SETTINGS_SPEC、v0.1.8検証手順、AvatarBeacon source READMEを対象に、古い記述や矛盾を洗い出す。
- [x] 現状仕様、配布物、リリース運用、設定UIの説明と合わない記述を更新する。
- [x] 重複している説明は必要に応じて整理し、参照先を明確にする。
- [x] リポジトリ内の全Markdownファイルを横断し、歴史的issue/reportを除く現行ドキュメントに追加の差分がないか確認する。
- [x] 更新後のMarkdownリンク切れや明らかな表記崩れがない。
- [x] 実装ファイルは変更せず、ドキュメント更新に限定する。

## 作業ログ

- 2026-07-07: 実装の `DefaultConfig()`、設定UI、Release workflow、AvatarBeacon source packageを確認した。
- 2026-07-07: READMEへ自動撮影スケジュール、フォールバックモード、OSC転送/ログ、未保存変更、保存ファイルを追記した。
- 2026-07-07: `src/SETTINGS_SPEC.md` を現行のタブ構成に合わせ、自動撮影、OSC、その他、Discord投稿を含めて更新した。
- 2026-07-07: `src/SPEC.md` の設定画面保存挙動、設定例、Release成果物名、AvatarBeacon source zip説明を実装・workflowに合わせた。
- 2026-07-07: `docs/v0.1.8-stream-spout-verification.md` の通常版配布物前提、Camera自動起動/終了、OSCログ収集手順を更新した。
- 2026-07-07: `avatar-gimmicks/AvatarBeacon/README.md` で、GitHub Releaseへ公開添付する標準Assetが source zip であり、`.unitypackage` は必要時の手動派生成果物であることを明記した。

## 追記 2026-07-07

### 指示

> issues/228を実施

### 文脈

既存の #228 を継続し、未完了の全Markdown横断確認とリンク・表記確認を実施する。

### 解釈

歴史的issue/reportは現行仕様追従の対象外としつつ、現行ドキュメントと管理用Markdownのリンク、配布物、リリース運用、設定説明に追加の差分がないか確認して完了させる。

- 2026-07-07: 全Markdownを横断し、現行ドキュメント対象のローカルMarkdownリンク切れがないことをスクリプトで確認した。
- 2026-07-07: Release公開Assetの運用が現行workflowでは3種類であるため、`AGENTS.md` のAvatarBeacon source zip添付説明を更新した。

## 追記 2026-07-07（仕様書再確認）

### 指示

> 仕様書などは実装と一致していましたか ズレていたら仕様書を修正して

### 文脈

前回の #228 完了処理ではRelease公開Asset運用とMarkdownリンク確認を中心に整理したが、仕様書本文と実装の一致確認について追加確認が必要になった。

### 解釈

`src/SPEC.md`、`src/SETTINGS_SPEC.md`、関連README/検証手順を実装コードと照合し、ズレがある場合は仕様書側を現行実装へ合わせて修正する。

- 2026-07-07: `DefaultConfig()` と設定UI templateを `src/SPEC.md` / `src/SETTINGS_SPEC.md` / `README.md` と再照合し、設定例の不足項目とWebhook入力欄の無効化条件のズレを仕様書側で修正した。
