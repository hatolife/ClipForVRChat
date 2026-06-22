# 040 GitHub Release のアップデート通知

## 状態

完了

## 問題

アプリ起動時に新しいGitHub Releaseがあるか分からず、ユーザーが手動で確認する必要がある。

## 期待する挙動

GitHub Releases の最新Releaseを確認し、現在のアプリより新しいバージョンがある場合はUI内に通知を表示する。通知を押すとGitHub Releaseページを開く。

## 受け入れ条件

- GitHub Releases の latest API を確認できる。
- 現在バージョンより新しいReleaseがある場合だけ通知する。
- 同じバージョンでも、GitHub側のRelease公開時刻がアプリに埋め込まれた時刻より明確に新しい場合は通知する。
- 通知をクリックするとGitHub Releaseページを開く。
- 確認に失敗しても通常操作を妨げない。

## 対応内容

- GitHub Releases latest API から最新タグ、公開時刻、Release URLを取得する。
- アプリに埋め込まれたバージョンとリリース時刻を比較する。
- 新しいReleaseがある場合、ヘッダー下に通知を表示する。
- 通知を押すとGitHub Releasesを開く。
