# 085 診断データ復号後のzipが破損する

## 問題

不具合報告用データを `gpg -d` で復号すると、生成されたzipが破損している。
OpenPGPのliteral data packetがテキスト扱いになっており、復号時にzipのバイト列が変換される。

## 期待する挙動

- 診断データ内のzipはOpenPGP literal data packetでbinaryとして暗号化される。
- `gpg -d diagnostics.zip.gpg > diagnostics.zip` で復号したzipを正常に展開できる。

## 受け入れ条件

- OpenPGP暗号化時に `FileHints.IsBinary` を指定する。
- テストで復号後のOpenPGP messageがbinary literal dataであることを確認する。
- テストで暗号化前zipが正常なzipであることを確認する。
