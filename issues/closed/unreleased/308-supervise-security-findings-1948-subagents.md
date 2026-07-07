# Codex Security findings 19:48 CSVの即時修正をサブエージェントで監督する

## 指示

> '/home/user/work/ClipForVRChat/reports/security/2026-07-07T17-53-16.203Z/codex-security-findings-2026-07-07T19-48-25.660Z.csv'　を底本にして同じ作業を実施

## 文脈

`307` の分類で、重大な仕様変更不要かつ追加情報不要としたfindingを修正対象にする。
前回同様、作業単位ごとにサブエージェントを使用し、メインエージェントが監督、統合、レビュー、テスト、コミットを担当する。
19:48 CSVの追加findingでは、設定インポート時に `OpenSettings(path)` が実行中の `configPath` を即時に切り替えてしまう点と、ディレクトリpickerがユーザー指定のUNC/device/network系パスへ `os.Stat` を先に当ててしまう点が対象になった。
19:48 CSVで追加された `Auto-photo scan cap can leak old photos and starve new ones` もこの修正対象に含める。

## 解釈

このissueでは、19:48 CSVで追加・再指摘された即時修正対象を、競合しにくい範囲へ分けてサブエージェントへ委任する。
仕様変更や追加情報が必要なfindingは修正対象から除外し、保留理由をレポートへ残す。

## 問題

- 新規追加findingは `src/app.go` と `src/frontend/src/main.js` に集中しており、設定インポートとpickerの境界を慎重に扱う必要がある。
- Release workflow系は実装ではなく運用・供給網信頼モデルの判断が必要であり、即時修正に混ぜると方針が曖昧になる。
- サブエージェント作業は監督なしだと、前回修正済みの境界を壊す可能性がある。

追加対象の具体的な修正点:

- インポートしたJSON設定はレビュー用のプレビューとして表示し、明示的な保存まで本番の `configPath` と runtime config を切り替えない。
- 閉じる/破棄操作では、インポート前のアクティブ設定へ戻す。
- ディレクトリpickerの `DefaultDirectory` は、安全なローカル候補だけを使い、UNC/device/network/default系の不明パスは事前に `os.Stat` しない。
- auto-photo の scan cap は、見つからなかった古い画像を後から新規扱いにしない形へ分離して修正する必要がある。

## 期待する挙動

- 即時修正対象findingが作業単位へ分割される。
- 各サブエージェントに明確な担当findingと担当ファイルを指定する。
- 修正差分を統合し、必要なテストを実行する。
- 作業単位ごとにコミットする。
- auto-photo は、全件ベースラインと per-tick の処理上限を分け、古い画像が後から新規投稿されないようにする。

## 受け入れ条件

- [x] 修正対象ごとにサブエージェントへ作業を委任する。
- [x] 仕様変更が必要なfindingと追加情報が必要なfindingを修正対象から除外する。
- [x] 修正差分を統合し、必要なテストを実行する。
- [x] 作業単位ごとにコミットする。
- [x] auto-photo scan cap の finding を修正する。
