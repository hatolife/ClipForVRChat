# AGENTS.md

## 作業ルール

- ユーザーから不具合、改善、仕様変更、調査依頼などの指示があった場合は、まず `issues/` に日本語のチケットを作成または追記してから作業する。
- チケットには「問題」「期待する挙動」「受け入れ条件」を簡潔に書く。
- 既存チケットに該当する場合は、新規作成ではなく該当チケットを更新する。
- 作業が完了したチケットは `issues/README.md` から該当行を外し、`issues/closed/README.md` へ状態 `完了` の行として移動する。チケットファイルも対応バージョンに応じて `issues/closed/vX.Y.Z/` などの分類フォルダへ移動し、`issues/closed/README.md` のリンクは移動後のファイルを指す形にする。
- GUIが表示されない、起動画面で止まる、フロントエンドエラー、Wails/HTML/Vue/API初期化不具合を扱う場合は、作業前に `docs/frontend-runtime-troubleshooting.md` を確認する。新しい既知パターンを同定した場合は、原因、同定方法、対処方法を同ドキュメントへ追記する。
- `src/frontend/src/main.js` のVue templateを変更した場合は、`node scripts/check-frontend-template-literals.mjs` と `node scripts/check-wails-api-surface.mjs` を実行する。
- 作業は意味のある単位で細かくコミットする。
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

- 正式リリース前に release candidate を確認する場合は、`vX.Y.Z-rcW` 形式のタグを打つ。例: `v0.1.7-rc1`。
- `vX.Y.Z-rcW` タグのGitHub Releaseは prerelease として作成される。Release本文は原則として `RELEASE_NOTES.md` の同じバージョン見出しから作成し、専用見出しがない場合は対応する `vX.Y.Z` の見出しを使う。
- 正式リリースは `vX.Y.Z` 形式のタグを打つ。GitHub Releaseは draft でも prerelease でもない通常Releaseとして作成される。
- `vX.Y.Z` / `vX.Y.Z-rcW` 以外の `v*` タグは検証用の draft Release として扱う。検証目的で作成したタグとReleaseは、確認後に削除する。
- `RELEASE_NOTES.md` の各バージョン本文は、GitHub Release本文としてそのまま使われるため、次の書式に揃える。
  - 先頭は `# vX.Y.Z`。
  - 続けて `## ダウンロード` を置き、`- [プログラムのダウンロード](https://github.com/hatolife/ClipForVRChat/releases/download/vX.Y.Z/ClipForVRChat-vX.Y.Z-windows-amd64.zip)` の1行を置く。
  - 続けて `### 更新内容` を置き、直前の正式バージョンから見て結果的に何が変わったかを端的に列挙する。RCごとの細かい修正差分や内部作業の羅列にしない。
  - 最後に `### 比較` を置き、`https://github.com/hatolife/ClipForVRChat/compare/vA.B.C...vX.Y.Z` を置く。
- GitHub Releaseへ公開添付するCI生成Assetは、通常利用者向けzip、単一exe署名asc、検証・切り分け用separated zip、AvatarBeacon source zipの4種類に絞る。sha256、個別exe、build metadataは公開Assetとして添付しない。
