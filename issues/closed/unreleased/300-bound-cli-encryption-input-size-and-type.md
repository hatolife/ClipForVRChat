# CLI暗号化の入力サイズとファイル種別を上限付きで検証する

## 指示

> CLI encryption for non-image files/directories/zip must reject non-regular special files and enforce conservative max input size, total directory size, file count, and depth before buffering.  
> Keep behavior compatible for normal diagnostic zip/file encryption.  
> Prefer simple bounded validation over large rewrites/streaming.  
> Add focused tests for oversized file and directory limit behavior if feasible. Run relevant go tests from src if feasible. Edit files directly and report changed files/tests.

## 文脈

`reports/security/2026-07-07T17-53-16.203Z/codex-security-findings-remediation-report.md` の F04/F21/F22 について、CLI経由の暗号化入力を先に検証し、巨大入力や特殊ファイルでメモリを圧迫しないようにする。

## 解釈

通常の単一ファイル暗号化、ディレクトリのzip化暗号化、診断zipの作成を壊さない範囲で、事前の bounded validation を入れる。
大きな設計変更やストリーミング化は行わず、入力の種別確認とサイズ上限で止める。

## 問題

- 非通常ファイルや特殊ファイルをCLI暗号化へ通すと、想定外の入力で失敗や資源浪費が起きる。
- ディレクトリ入力は buffering 前に件数、深さ、総サイズを抑えないと、巨大ツリーでメモリを使い切る。
- 既存の正常系である通常ファイル暗号化と診断zip作成は維持する必要がある。

## 期待する挙動

- CLI暗号化は通常ファイルまたはディレクトリのみ受け付ける。
- 特殊ファイル、巨大ファイル、巨大ディレクトリは明確なエラーで拒否される。
- 診断zipと通常の小さな暗号化入力は従来どおり処理できる。
- 代表的な上限制御についてテストがある。

## 受け入れ条件

- [x] CLI暗号化で非通常ファイル/特殊ファイルを拒否する。
- [x] 単一ファイル暗号化に保守的なサイズ上限を入れる。
- [x] ディレクトリ暗号化に総サイズ、件数、深さの上限を入れる。
- [x] 既存の通常ファイル暗号化と診断zip作成の正常系が維持される。
- [x] 巨大ファイルとディレクトリ上限制御の focused test を追加する。
