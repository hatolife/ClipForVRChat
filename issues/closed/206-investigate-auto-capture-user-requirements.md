# 自動撮影機能を使用するための条件を整理する

## 問題

ユーザーが自動撮影機能を使うために必要なVRChat側、ClipForVRChat側、撮影方式別の条件が複数の実装・検証資料に分散している。

## 期待する挙動

自動撮影を使うための前提条件、Stream方式の条件、Photo方式の条件、構図基準の条件、任意機能の条件を列挙できる。

## 受け入れ条件

- [x] 現行実装と検証資料から利用条件を確認する。
- [x] ユーザー向けに条件を簡潔に列挙する。

## 調査結果

2026-07-05に、`src/internal/appcore/autocapture.go`、`src/internal/appcore/config.go`、`src/internal/appcore/player_local.go`、`src/app.go`、`docs/v0.1.8-stream-spout-verification.md`、`docs/v0.1.8-avatar-osc-basis-verification.md`、`docs/v0.1.8-player-local-verification.md`、`avatar-gimmicks/AvatarBeacon/Assets/PoppoWorks/AvatarBeacon/README.md` を確認した。

自動撮影の利用条件は、ユーザー向け回答として前提条件、Stream方式、Photo方式、構図基準、任意機能に分けて整理した。
