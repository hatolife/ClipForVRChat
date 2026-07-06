# GitHub Release本文からバージョン見出しを除外する

## 問題

GitHub Releases に表示されるリリース本文へ、`# v0.1.8-b2` のようなバージョン見出しが含まれている。GitHub Releaseのタイトルにも同じバージョンが表示されるため、本文内で二重表記になる。

## 期待する挙動

Release workflow が作成する本文は、先頭の単独バージョン見出しを含めず、`## ダウンロード` など本文内容から始まる。

## 受け入れ条件

- `RELEASE_NOTES.md` に `# vX.Y.Z` 形式の本文用見出しが残っていても、生成される `dist/release-body.md` には含まれない。
- prerelease が正式版のリリースノートへフォールバックする場合も、生成本文に `# vX.Y.Z-aW` / `# vX.Y.Z-bW` / `# vX.Y.Z-rcW` が含まれない。
- Release workflow とローカル抽出スクリプトで同じ本文生成結果になる。

## 対応メモ

- `RELEASE_NOTES.md` の現行 `v0.1.8` から本文用の `# v0.1.8` 見出しを削除する。
- 抽出スクリプト側でも先頭の単独バージョン見出しを除外し、古い形式が残ってもGitHub Release本文へ出ないようにする。

## 完了メモ

- `scripts/extract-release-notes.mjs` で先頭の単独バージョン見出しを除外するようにした。
- Release workflowの通常リリース本文生成を同スクリプトへ統一した。
- `RELEASE_NOTES.md` と `AGENTS.md` のリリースノート書式を、本文が `## ダウンロード` から始まる形へ更新した。
