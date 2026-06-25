# Release NotesのダウンロードURLをMarkdownリンクにする

## 問題

GitHub Release本文のダウンロード欄がURLをそのまま表示しており、読みづらい。

## 期待する挙動

プログラム本体と署名確認用ファイルは、ファイル名をリンクテキストにしたMarkdownリンクで表示する。

## 受け入れ条件

- `RELEASE_NOTES.md` のv0.1.7ダウンロード欄がMarkdownリンク形式になっている。
- `v0.1.7-rc1` のRelease本文でも、RC版ファイル名をリンクテキストにしたMarkdownリンクとして表示される。
