# rc18でGUIが表示されずavatar_oscエラーが大量出力される

## 問題

`v0.1.8-rc18` の単一exe版、separated版のどちらでも起動してもGUIが表示されない。ユーザー提供ログでは、起動直後から `avatar osc basis resolve error: status="partial" err="partial avatar OSC basis sample: missing coord/x"` などが大量に連続出力され、短時間でログが肥大化している。

## 期待する挙動

- `avatar_osc` の受信状態がpartial/missingでも、GUIは通常通り表示される。
- AvatarBeaconやVRChat OSCの状態に問題がある場合も、ログ出力は抑制され、UI上の受信状態で確認できる。
- 起動直後に大量ログ出力や高負荷でアプリが固まらない。

## 受け入れ条件

- rc18相当の設定/OSC受信状態でもGUI表示を阻害しない。
- `avatar osc basis resolve error` がループで大量出力されない。
- Wails起動、DOM準備、frontend初期化、初期状態取得、手動終了開始、shutdown完了が診断ログで追える。
- WebViewがHTML表示まで進んだ場合は、Vue初期化前でも簡易的な起動進捗が表示される。
- frontend template内の説明文でJS template literalを壊さず、起動時に `avatar is not defined` などのReferenceErrorが出ない。
- 同種のtemplate内バッククォート混入をCI/Releaseで検出できる。
- `go test ./...` とフロントエンドビルドが通る。
- 必要なら次RCで修正を確認できる。

## 調査メモ

- 提供ログは43,275行中43,263行が `avatar osc basis resolve error` だった。
- `basis_source="manual"` でもAvatar OSC受信処理が全Avatar Parameter受信ごとにbasis再構築を試みていた。
- `latestAvatarOSCBasisSnapshotLocked` がmanual設定のまま `ResolvePlayerLocalBasisPose` を呼び、manual基準Pose未設定エラーをAvatar OSC診断にも混ぜていた。
- 対策として、Avatar OSC診断snapshotでは復元時だけ `BasisSource=avatar_osc` を強制し、basisに関係するOSC address以外では再構築しない。非readyログも状態変化時だけに抑制する。

## rc19再現ログ

- `v0.1.8-rc19` separated版でもGUIが表示されない。
- ログは6行のみで、大量ログは解消している。
- 20:31:17にstartupし、20:31:22に `auto-capture osc receiver stop: err=context canceled` で終了している。ユーザー操作で手動終了した可能性があるため、この終了時刻だけで自動終了とは判断しない。
- Avatar OSC basisは20:31:20にready化しており、Avatar OSC処理自体は停止原因ではなさそう。
- 次の調査対象はWailsウィンドウ生成、起動オプション、Windows GUI subsystem、またはfrontend asset読み込み失敗。

## 追加デバッグ方針

- `runUI` 開始/終了、embedded frontend asset概要、Wails `OnStartup` / `OnDomReady` / `OnBeforeClose` / `OnShutdown` を診断ログへ出す。
- frontend側はscript load、Vue mount開始、mounted開始、各Wails API呼び出し、mounted完了、JS error/unhandled rejectionを `LogUserAction` 経由で診断ログへ出す。
- `index.html` に静的な起動画面を入れ、JSやVue mount前にWebView表示まで進んでいるかを視覚的に確認できるようにする。
- Vue初期化中も起動オーバーレイを出し、Go API取得で時間がかかっているのか、表示処理自体が止まっているのかを切り分ける。
- 過去の同種確認として `issues/003-main-window-ui-and-about.md` に「通常起動時に空白画面にならない」要件がある。現状の `index.html` は空の `#app` のため、JS実行前に空白になり得る点は同根のUX不足として扱う。一方で `issues/048` はGUI subsystem exeの標準出力問題、`issues/105` はGoテスト中のWailsイベント送信問題であり、今回のウィンドウ表示停止とは直接同根ではなさそう。

## rc20再現ログ

- `v0.1.8-rc20` separated版で、簡易起動画面に `フロントエンドエラー: Uncaught ReferenceError: avatar is not defined at http://wails.localhost/assets/index-DvTLSiJp.js:397:68` が表示された。
- 診断ログでは `frontend_script_loaded api=available` の直後に同じ `frontend_error` が出ており、Wails APIとHTML表示までは到達している。
- 該当箇所はVue `template: \`...\`` 内の説明文 `Debug OSC Pingは \`/avatar/parameters/avatar_beacon/debug/ping\`` で、説明用バッククォートがJS template literalを途中で閉じ、`/avatar/...` が不正に評価されてReferenceErrorになっていた。

## rc21確認と再発防止

- `v0.1.8-rc21` でGUIが表示されることを確認した。
- 直接原因は `src/frontend/src/main.js` のVue template literal内に、説明文用の生バッククォートを入れたことだった。
- 原因同定は、簡易起動画面の `フロントエンドエラー` 表示、診断ログの `frontend_script_loaded api=available` 直後の `frontend_error`、ビルド済みassetの該当行を組み合わせて行えた。
- 同種の不具合をCI/Releaseで検出するため、`scripts/check-frontend-template-literals.mjs` を追加した。
- 再発時に最初に確認する資料として、`docs/frontend-runtime-troubleshooting.md` に原因、同定方法、対処方法、確認コマンドをまとめる。
