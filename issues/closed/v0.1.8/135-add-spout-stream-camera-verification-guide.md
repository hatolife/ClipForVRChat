# Stream Camera/Spout方式の実機確認手順を整備する

## 問題

Spout senderの有無、sender名、VRChat側のCamera/OSC状態は実機環境でしか確認できない。
確認手順が曖昧だと、デスクトップ周辺取得や白画像の再発を見落とす可能性がある。

## 期待する挙動

v0.1.8 RCを実機確認するときに、Stream Camera(Spout)映像そのものを保存できているかを再現性のある手順で確認できる。

## 受け入れ条件

- VRChat側でOSCを有効化する手順を記載する。
- 自動撮影タブでSpout helper確認、sender一覧取得、sender選択を行う手順を記載する。
- カメラを手動で出した場合と、自動撮影が起動した場合の確認観点を分けて記載する。
- 正面/背面/斜めなど複数構図のテスト撮影で、Stream Camera映像そのものが保存されることを確認する手順を記載する。
- デスクトップ周辺、VRChatウィンドウ領域、白画像が保存された場合のログ取得と確認箇所を記載する。
- sidecar JSON、Discord投稿、Presenceユーザー情報、履歴にSpout取得画像が紐づくことを確認する手順を記載する。
- 実機確認で不具合が出た場合に、提出してほしいログ/config/sidecar/画像を明記する。

## 実装メモ

- 実装完了後にREADMEまたは専用 `docs/` 配下へ手順化する。
- 受け入れ条件は、ユーザーがRCで確認する最小要件と対応させる。

## 対応内容

- v0.1.8実装に合わせて対応済み。詳細は `feat/v0.1.8-resolve-issues` の実装、README、SPEC、RELEASE_NOTESを参照。

## 再監査メモ

- 2026-07-01: [#164](164-audit-v018-completed-items.md) の再監査で未達が見つかったため、完了扱いを取り消して `要対応` に戻した。

## 実装前調査メモ

実装方針:

- 実機確認手順はREADME本文へ詰め込まず、`docs/v0.1.8-stream-spout-verification.md` のような専用文書を追加する。
- 手順は「事前準備」「手動Camera表示中」「Cameraを閉じた状態から自動起動」「失敗時収集物」に分ける。
- Discord/sidecar/履歴の紐づけ確認も同じ手順に含める。

対象ファイル:

- `docs/v0.1.8-stream-spout-verification.md` 新設
- `README.md`
- `src/SPEC.md`

小タスク:

- VRChat側でOSCを有効にする手順を記載する。
- 自動撮影タブでhelper確認、sender一覧更新、sender選択、start delay設定を行う手順を記載する。
- 正面/背後/斜めのテスト撮影で、保存画像がデスクトップ周辺/白画像ではなくStream Camera映像そのものか確認する観点を記載する。
- sidecar JSONの `stream.backend=spout`、sender名、同席ユーザー、world/instance、画像SHA256を確認する手順を記載する。
- Discord投稿画像のメタデータ保持はDiscord側仕様が保証しないため、投稿本文とsidecarを正本として確認する方針を記載する。
- 失敗時に提出する `logs/*.log`、`config.json`、`history.json`、画像、画像`.json`、sender一覧スクリーンショットを明記する。

確認方法:

- 手順だけでRC利用者が同じログ/config/sidecar/画像を収集できる。
- 手動Camera表示中と自動起動時の差分を切り分けられる。
