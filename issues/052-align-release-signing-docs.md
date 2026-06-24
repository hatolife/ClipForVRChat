# Release署名ファイルの説明を実際の成果物に合わせる

## 問題

Release workflow は zip 内に `.asc` を含めず、Release asset として `ClipForVRChat-vX.Y.Z-windows-amd64.exe.asc` を添付する。一方で README、アプリ内About、SPEC は zip 内に署名ファイルや公開鍵 `.asc` があるように説明している。

## 期待する挙動

ユーザー向け説明、仕様書、Release workflow の成果物一覧が一致している。

## 受け入れ条件

- README のPGP確認手順が現行Release成果物と一致している。
- アプリ内AboutのPGP確認手順が現行Release成果物と一致している。
- SPEC のRelease同梱物・添付物一覧が現行Release workflow と一致している。
