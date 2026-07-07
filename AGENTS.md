# AGENTS.md

## 作業ルール

- ユーザーから不具合、改善、仕様変更、調査依頼などの指示があった場合は、まず `issues/` に日本語のチケットを作成または追記してから作業する。
- チケットには先頭から「指示」「文脈」「解釈」を書き、その後に「問題」「期待する挙動」「受け入れ条件」を簡潔に書く。「指示」にはユーザー発言を誤字も含めて原文のまま引用する。ただし、秘密情報・プライバシー情報は原文引用よりredactionを優先し、コミット前に `[REDACTED: 種類]` 形式へ置換する。対象例は Discord Webhook URL/token、ユーザー名を含むローカルパス、QRコード由来URL、診断ログ内のsecret/token/key/API credentialなどとする。例: `[REDACTED: discord webhook url]`、`[REDACTED: discord webhook token]`、`[REDACTED: local user path]`、`[REDACTED: qr url]`、`[REDACTED: diagnostic log secret]`。redaction後も指示の意図と作業判断に必要な文脈が残るように、秘密値の周辺文は可能な限り保持する。
- 既存チケットに該当する場合は、新規作成ではなく該当チケットを更新する。
- 作業が完了したチケットは `issues/README.md` から該当行を外し、`issues/closed/README.md` へ状態 `完了` の行として移動する。チケットファイルも対応バージョンに応じて `issues/closed/vX.Y.Z/` などの分類フォルダへ移動し、`issues/closed/README.md` のリンクは移動後のファイルを指す形にする。
- GUIが表示されない、起動画面で止まる、フロントエンドエラー、Wails/HTML/Vue/API初期化不具合を扱う場合は、作業前に `docs/frontend-runtime-troubleshooting.md` を確認する。新しい既知パターンを同定した場合は、原因、同定方法、対処方法を同ドキュメントへ追記する。
- `src/frontend/src/main.js` のVue templateを変更した場合は、`node scripts/check-frontend-template-literals.mjs` と `node scripts/check-wails-api-surface.mjs` を実行する。
- `issues/closed/README.md` を変更した場合は、`node scripts/check-closed-issue-index.mjs` を実行する。
- 作業は意味のある単位で細かくコミットする。
- コミットメッセージの件名はConventional Commits形式にする。例: `fix(ui): prevent empty state overlap`、`docs(issues): close release checklist`、`ci(release): validate prerelease tags`。履歴を書き換えた後や複数コミットを作った後は、必要に応じて `node scripts/check-commit-subjects.mjs <range>` で確認する。
- Release成果物の確認が必要な変更では、GitHub ActionsのCI/Releaseステータスを確認する。
- リリース時は `issues/README.md` の対象チケットについて、対応バージョン欄をリリースするバージョンに更新する。
- Release workflowを変更した場合は、GitHub Release本文が `RELEASE_NOTES.md` の該当バージョンから作成されること、Release添付ファイル一覧、zip内ファイル一覧を確認する。不要な公開鍵ファイルなど、仕様外の成果物を添付・同梱しない。

## todo.md運用

- `todo.md` は直近作業の短期チェックリストとして使う。長期的な課題管理や完了済み作業の記録は `issues/` とコミット履歴に寄せる。
- 新しいまとまった作業を始めるときは、不要になった項目を削除し、現在の作業に必要な項目だけへ更新する。
- 作業中は進捗に合わせてチェックを更新し、作業完了時点で次の作業に不要な完了済み項目を残し続けない。

## ブランチ運用

- `master` は正式リリース済み、または正式リリース可能な安定状態を保つ。通常開発を直接 `master` で進めない。
- `develop` は次リリースへ向けた通常開発の統合ブランチとする。通常の不具合修正、改善、仕様変更、調査は `develop` を基準に `fix/...`、`feat/...`、`docs/...`、`chore/...` など目的が分かる作業ブランチを切って進める。
- 作業ブランチは作業完了後にCIが通っていることを確認してから `develop` へ取り込む。`master` へ直接取り込むのはリリース反映や緊急修正など、ユーザーが明示した場合に限る。
- リリース前の安定化作業は `master` から `release/vX.Y.Z` ブランチを作成して行う。`develop` のリリース対象変更を `release/vX.Y.Z` へ取り込み、RC確認とリリース向け修正をそのブランチで行う。
- `release/vX.Y.Z` で問題がないことを確認したら `master` へ取り込み、正式タグを打つ。リリースブランチ上の修正は `develop` にも反映し、次開発へ差分を残さない。
- 履歴を書き換える必要がある場合は、作業ブランチ上で整理してから `master` へ反映する。共有済みブランチやタグの force push は、ユーザーの明示指示がある場合のみ行う。

## リリース運用

- ユーザーが「`vX.Y.Z` でこの機能を実装する」と指定または合意した場合、その機能は `vX.Y.Z` のリリース対象として扱う。難易度、残作業、実装リスクを理由に、ユーザーの明示承認なしで未実装部分を次バージョンへ延期したり、仕様を縮小して完了扱いにしたりしない。
- 対象バージョンから機能を外す、仕様を縮小する、または次バージョンへ送る必要がある場合は、理由、未実装範囲、影響を明示し、ユーザーの承認を得てから `issues/` とリリースノート方針を更新する。
- alpha、beta、rc、正式版への段階移行は、ユーザーの明示的な合意に基づいて行う。特に、機能開発中や最低機能確認段階の成果物を、作業者判断だけで `vX.Y.Z-rcW` として扱わない。
- 開発版タグは、段階に応じて `vX.Y.Z-aW`、`vX.Y.Z-bW`、`vX.Y.Z-rcW` を使い分ける。`W` は同じ段階内の連番とし、タグの付け直しは共有済みブランチやタグの force push 禁止ルールに従う。
- `vX.Y.Z-aW` は alpha として扱う。未完成機能、実装途中の機能、破壊的変更の早期確認、CI/CDと配布物生成の確認に使う。正式リリース品質を期待しない。
- `vX.Y.Z-bW` は beta として扱う。対象機能が概ね揃い、主要フローの機能確認と不具合洗い出しを行う段階で使う。既知の未解決問題や実機確認待ちが残っていてよい。
- `vX.Y.Z-rcW` は release candidate として扱う。既知の重大問題がなく、Release成果物、ドキュメント、手作業チェックを通せばそのまま `vX.Y.Z` として正式リリースできる候補だけに使う。機能追加や大きな仕様変更を含む確認版には使わない。
- `vX.Y.Z-aW` / `vX.Y.Z-bW` / `vX.Y.Z-rcW` / `vX.Y.Z` の各タグでは CI/CD と Release workflow を走らせ、GitHub Releaseと配布成果物を作成する。
- `vX.Y.Z-aW` / `vX.Y.Z-bW` / `vX.Y.Z-rcW` タグのGitHub Releaseは prerelease として作成される。Release本文は原則として `RELEASE_NOTES.md` の同じバージョン見出しから作成し、専用見出しがない場合は対応する `vX.Y.Z` の見出しを使う。
- 正式リリースは `vX.Y.Z` 形式のタグを打つ。GitHub Releaseは draft でも prerelease でもない通常Releaseとして作成される。
- `vX.Y.Z` / `vX.Y.Z-aW` / `vX.Y.Z-bW` / `vX.Y.Z-rcW` 以外の `v*` タグは原則として作成しない。検証目的で一時タグが必要な場合は、GitHub Releaseを公開しない運用に留め、確認後にタグとReleaseを削除する。
- `RELEASE_NOTES.md` の各バージョン本文は、GitHub Release本文として使われるため、次の書式に揃える。GitHub Releaseのタイトルにバージョン番号が表示されるため、本文先頭に `# vX.Y.Z` のようなバージョン見出しを置かない。
  - 先頭は `## ダウンロード` を置き、`- [プログラムのダウンロード](https://github.com/hatolife/ClipForVRChat/releases/download/vX.Y.Z/ClipForVRChat-vX.Y.Z-windows-amd64.zip)` の1行を置く。
  - 続けて `### 更新内容` を置き、直前の正式バージョンから見て結果的に何が変わったかを端的に列挙する。RCごとの細かい修正差分や内部作業の羅列にしない。
  - 最後に `### 比較` を置き、`https://github.com/hatolife/ClipForVRChat/compare/vA.B.C...vX.Y.Z` を置く。
- GitHub Releaseへ公開添付するCI生成Assetは、通常利用者向けzip、単一exe署名asc、検証・切り分け用separated zip、AvatarBeacon source zipの4種類に絞る。sha256、個別exe、build metadataは公開Assetとして添付しない。

## Spout2 revision更新運用

- Spout2の取得元revisionは、Releaseごとに自動更新しない。固定したcommit/archiveを使い、更新理由がある場合だけ明示的なissueで変更する。
- `vX.Y.Z-a1`、`vX.Y.Z-b1`、`vX.Y.Z-rc1` を作成する前に、Spout2 revisionを更新すべきか調査する。更新すべきと判断した場合は、その段階の `a2`、`b2`、`rc2` で更新が反映されるように作業する。
- `vX.Y.Z-a2`、`vX.Y.Z-b2`、`vX.Y.Z-rc2` を作成する前に、直前の `a1`、`b1`、`rc1` でSpout2更新要否の確認が正しく行われたかを確認する。確認漏れ、判断誤り、またはVRChat側更新などで再確認が必要な場合は、そのタイミングで確認と更新処理を行う。
- VRChatの現在バージョンに更新があった場合は、Stream Camera/Spout挙動への影響を前向きに疑い、Spout2更新要否を積極的に確認する。更新すべきと判断した場合は、固定revision、archive hash、Release metadata、ライセンス表記、Spout helper動作確認をまとめて更新する。
- Spout2 revisionを更新する場合は、旧revision、新revision、対応するupstream tag、archive URL、archive SHA256、更新理由、CI/Release build結果、`spout-capture.exe --help` / `--version` / `--list-senders`、Stream撮影の実機確認結果、ライセンス差分確認をissueへ記録する。
