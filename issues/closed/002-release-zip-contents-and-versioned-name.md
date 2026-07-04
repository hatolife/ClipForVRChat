# Release zipの内容とバージョン付きファイル名

## 問題

Windows向けRelease zipにはユーザー向けファイルを含める必要があり、zip名にもバージョン情報が必要。

## 期待する挙動

Windows向けRelease zipには以下を含める。

- `ClipForVRChat.exe`
- `README.md`
- `LICENSE`

zipファイル名にはタグのバージョンを含める。

## 受け入れ条件

- `vX.Y.Z` のReleaseで `ClipForVRChat-vX.Y.Z-windows-amd64.zip` がアップロードされる。
- zipにexe、README、LICENSEが含まれる。
- 不要な `ClipForVRChat-source.zip` はReleaseに添付しない。
