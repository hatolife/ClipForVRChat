# ReleaseにSBOMまたはビルドメタデータを追加する

## 問題

Release成果物にSBOMやビルドメタデータが含まれず、依存関係やビルド環境の追跡性が不足している。

## 期待する挙動

Release成果物から、ビルドに使ったコミット、Go/Node/Wailsバージョン、主要依存関係を追跡できる。

## 受け入れ条件

- Release workflowでビルドメタデータまたはSBOMを生成する。
- Release assetsに添付される。
- 成果物検証手順に反映される。

## 備考

監査上はInfo項目。

## 対応内容

- Release workflowで `ClipForVRChat-<tag>-build-metadata.json` を生成するようにした。
- build metadataにtag、commit、app version、Go/Node/Wailsバージョンを記録するようにした。
- build metadataをGitHub Release assetへ添付するようにした。
