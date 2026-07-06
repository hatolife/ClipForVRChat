# Release build metadata assetの意図を確認する

## 問題

CI/Release assetsに添付される `ClipForVRChat-v0.1.8-a37-build-metadata.json` が、利用者向け成果物なのか検証向け成果物なのか分かりにくい。

## 期待する挙動

このJSONの生成箇所、含まれる情報、添付される意図を説明できる。

## 受け入れ条件

- Release workflow上の生成・添付箇所を確認する。
- 既存issue/SPECに残っている追加意図を確認する。
- ユーザーへ用途と通常利用者への影響を説明する。

## 調査結果

- `.github/workflows/release.yml` の `Write build metadata` で生成し、`Upload release assets artifact` と `gh release upload` でRelease assetへ添付している。
- 元の追加意図は #104「ReleaseにSBOMまたはビルドメタデータを追加する」で、依存関係やビルド環境の追跡性を確保するため。
- `v0.1.8` 系では単一exeへ埋め込んだ `spout-capture.exe` と `SpoutLibrary.dll` のSHA-256も記録し、単一exe版と分離版が同じ元バイナリから作られたことを確認しやすくしている。
