# ソース配置を src 配下へ移動する

## 問題

アプリ実装ファイルがリポジトリ直下に散らばっている。

## 期待する挙動

README、LICENSE、issues、GitHub Actions、AGENTS.md を除き、アプリ本体と関連ソースを `src/` 配下へまとめる。

## 受け入れ条件

- Go/Wailsアプリ本体は `src/` 配下にある。
- CI/Release は `src/` を作業ディレクトリとしてビルドできる。
- Release zip には引き続き exe、README、LICENSE が含まれる。
