# セキュリティ監査サマリ

## 全体サマリ

ClipForVRChat は、画像変換、Discord Webhook 投稿、履歴管理、VRChat写真/スクリーンショット監視、OSCによるVRChat User Camera制御、Spout helperによるStream Camera取得を行うWindows向けWailsアプリである。今回の監査では、ローカルコード、設定、CI/Release workflow、依存定義、Spout C++ helper、既存セキュリティ対応履歴を確認した。

過去の主要リスクであるDiscord Webhook任意URL投稿、履歴へのWebhook token保存、診断ログ/診断zipのsecret露出、Release署名secretのbuild job露出、Spout2未固定取得、AvatarBeacon package未検証、config-controlled ffmpeg任意実行、隠れた自動投稿設定は、現行コードでは概ね対策済みだった。

一方で、Windows上で外部実行ファイルを起動する箇所に、PATH検索へ依存する残存リスクがある。特に分離版または埋め込みhelperが利用できない場合の `spout-capture.exe` fallback は、同名の不正実行ファイルがPATH上にある環境で意図しないhelperを起動し得る。また、Startup shortcut作成とffmpegインストール補助では `powershell.exe` / `winget` を固定名で起動している。

## 件数

| Severity | 件数 |
| --- | ---: |
| Critical | 0 |
| High | 0 |
| Medium | 1 |
| Low | 2 |
| Info | 2 |

## リリース可否

条件付き可。通常利用者向け単一exeの主要フローは、埋め込みSpout helperを優先するため直ちにブロックとはしない。ただし分離版zipやhelper欠損環境のPATH fallbackは、リリース前に修正することを推奨する。

## 最優先修正項目

1. `ResolveSpoutHelperPath` の default/helper名指定時に `exec.LookPath` fallbackを使わず、埋め込みhelperまたは本体exe隣接のhelperだけに限定する。
2. Windowsの `powershell.exe` / `winget` 起動を絶対パスまたは信頼済み解決結果に固定する。
3. GitHub Actionsの外部Actionをmajor tagではなくコミットSHA pinningへ移行する運用を検討する。

## 残存リスク

- ユーザーが明示的に外部helperパスを指定した場合、その実行ファイルの真正性はユーザー判断に依存する。
- Discord Webhook URL、QRコードURL、VRChatログ由来の同席ユーザー名/World情報は、設定や有効化状態によってローカル保存またはDiscord投稿される。
- `npm audit` / `govulncheck` は今回ローカルでは外部DB問い合わせを避けたため、最新脆弱性DBに対する結果はCIまたはネットワーク許可環境での確認が必要。
- Windows本番ビルド、ASLR/DEP/CFG、SmartScreen、PGP署名検証はローカルLinux環境では確認していない。

## 未確認事項

- Windows上でのWails本体ビルドとRelease成果物内容。
- `spout-capture.exe --help` / `--version` / `--list-senders` のWindows実行。
- 最新のnpm/Govulncheck脆弱性DB照会結果。
- CodeQL、clang-tidy、cppcheck、ASan、UBSan、golangci-lint、cargo audit。

## 監査時点での前提条件

- 監査日時: 2026-07-08 06:13 JST頃。
- ブランチ: `develop`。
- 外部サービスへの能動的アクセスは禁止と解釈し、GitHub/npm/Go脆弱性DBなどへの問い合わせ型監査は実施していない。
- 実行環境はLinux/WSL相当であり、Windows固有の動的検証は未実施。
