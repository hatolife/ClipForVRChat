# Webhook URL保存時の確認ダイアログ循環を防ぐ

## 指示

> webhook url 設定したときダイアログが連続して大量に出る

## 文脈

`v0.1.8-rc4` で通常投稿用Webhook URLを変更して保存すると、「重要な設定変更」と「自動処理の投稿内容」の確認ダイアログが繰り返し表示される。

## 解釈

Webhook URL変更は両方の確認対象になるが、各確認は1回だけ表示する。自動投稿確認を確定した後は、先に完了した重要設定確認を再要求せず保存または画面移動を続行する。

## 問題

自動投稿確認の確定処理が、自動投稿確認だけをスキップして保存処理を再開している。未保存状態では重要設定変更の判定が引き続き真になるため、重要設定確認へ戻り、両ダイアログ間で循環する。

## 期待する挙動

- Webhook URL変更時に必要な確認をそれぞれ最大1回だけ表示する。
- 自動投稿確認を確定すると保存処理が完了する。
- 設定画面から移動する場合も同じ確認循環が発生しない。
- キャンセルや該当設定を開く既存動作は変えない。

## 受け入れ条件

- [x] 通常保存で重要設定確認と自動投稿確認を確定した後、保存処理を確認済み状態で再開する。
- [x] 保存後の画面移動でも確認済み状態を引き継ぐ。
- [x] 確認再開時に両確認を完了済みとして扱うことを自動テストする。
- [x] Frontend test、template検査、Wails API surface検査が成功する。
- [x] Frontend buildが成功する。
- [x] Windows CIが成功する。

## 確認

- ローカル: Frontend test 9件成功。
- ローカル: Frontend template literal check成功。
- ローカル: Wails API surface check成功。
- ローカル: Vite依存が未展開のためFrontend buildは未実施。GitHub Actionsの `npm ci` 後に確認する。
- GitHub Actions CI #268: Frontend test/build、Goテスト、Spout CTest、Wails build、GUI lifecycle smokeを含めて成功。
