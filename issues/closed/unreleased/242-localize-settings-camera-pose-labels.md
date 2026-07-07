# 設定内のカメラ姿勢Pose表記を日本語へ変更する

## 指示

> 設定内のカメラ姿勢を表すPoseという表記を全て日本語にしたい 基本的に位置でいいと思う

## 文脈

設定画面にはカメラ位置・回転を表す `Pose` 表記が残っている。利用者向け表示としては英語のままだと意味が伝わりにくい。

## 解釈

内部コード、OSC address `/usercamera/Pose`、JSON key、ログなどの技術名は維持し、設定UIや利用者向け説明に出るカメラ姿勢の `Pose` を日本語へ置き換える。基本は「位置」とし、回転も明示した方が誤解が少ない箇所は「位置・向き」とする。

## 問題

- 設定内に `Pose` という英語表記が残っている。
- カメラの位置・向き設定なのか、OSC endpoint名なのかが混ざって見える。

## 期待する挙動

- 設定UI上のカメラ姿勢を表す `Pose` 表記が日本語になる。
- 技術的なOSC endpoint名 `/usercamera/Pose` は必要な箇所で維持される。

## 受け入れ条件

- [x] 設定UI内のカメラ姿勢を表す `Pose` が日本語表記になる。
- [x] 必要な検査を実行する。

## 確認

- `node scripts/check-frontend-template-literals.mjs`
- `node scripts/check-wails-api-surface.mjs`
- `cd src && GOCACHE=/tmp/clipforvrchat-go-cache go test ./...`
- `git diff --check`
