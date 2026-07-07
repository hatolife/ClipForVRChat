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
- [ ] リポジトリ内の全Markdownファイルを横断し、歴史的issue/reportを除く現行ドキュメントに追加の差分がないか確認する。
- [ ] 更新後のMarkdownリンク切れや明らかな表記崩れがない。
- [x] 実装ファイルは変更せず、ドキュメント更新に限定する。

## 作業ログ

- 2026-07-07: 実装の `DefaultConfig()`、設定UI、Release workflow、AvatarBeacon source packageを確認した。
- 2026-07-07: READMEへ自動撮影スケジュール、フォールバックモード、OSC転送/ログ、未保存変更、保存ファイルを追記した。
- 2026-07-07: `src/SETTINGS_SPEC.md` を現行のタブ構成に合わせ、自動撮影、OSC、その他、Discord投稿を含めて更新した。
- 2026-07-07: `src/SPEC.md` の設定画面保存挙動、設定例、Release成果物名、AvatarBeacon source zip説明を実装・workflowに合わせた。
- 2026-07-07: `docs/v0.1.8-stream-spout-verification.md` の通常版配布物前提、Camera自動起動/終了、OSCログ収集手順を更新した。
- 2026-07-07: `avatar-gimmicks/AvatarBeacon/README.md` で、GitHub Releaseへ公開添付する標準Assetが source zip であり、`.unitypackage` は必要時の手動派生成果物であることを明記した。
