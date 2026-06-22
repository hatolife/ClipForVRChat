# 034 GitHub ActionsビルドのバージョンにコミットIDを含める

## 状態

完了

## 背景

GitHub ActionsでRelease用exeをビルドしたとき、「このアプリについて」に表示されるバージョンで、リリース番号とビルド元コミットを同時に確認できるようにしたい。

## 対応内容

- Release workflow の `main.version` に `${tag}.${short_sha}` を渡す。
- 例: `v0.1.3.39531cf`

