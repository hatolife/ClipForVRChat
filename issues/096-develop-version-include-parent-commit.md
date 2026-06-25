# develop版バージョンに親コミットIDを含める

## 問題

ローカルビルドや開発ビルドでバージョンが `develop` だけになることがあり、どのコミット由来のexeか判断しにくい。

## 期待する挙動

開発ビルドでは `vX.Y.Z.aaaaaaa.develop` のように、ベースバージョン、親コミットID、develop識別子を含むバージョン表記になる。

## 受け入れ条件

- develop版の表示バージョンに親コミットIDの短縮ハッシュが含まれる。
- 形式は `vX.Y.Z.aaaaaaa.develop` を基本とする。
- Release/RCビルドの既存バージョン表記を壊さない。
- `ClipForVRChat.exe --version` とアプリ情報表示で同じバージョン情報を参照できる。

## 対応内容

- 開発ビルドの既定バージョンを `v0.1.7.<commit>.develop` 形式にした。
- ローカルWindowsビルドスクリプトで、`--version` 未指定時だけdevelop識別子を付けるようにした。
- Release/RC/CIでは `buildChannel` を明示し、Release表記に `.develop` が混ざらないようにした。
