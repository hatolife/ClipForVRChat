# フロントエンド実行時エラーの切り分け

この資料は、GUIが表示されない、起動画面のまま止まる、または `フロントエンドエラー` が表示される場合に最初に確認する。Wails、HTML、Vue、Wails API、Go backendのどこまで進んだかを切り分けるための手順をまとめる。

## まず見るログ

配布版では、exeと同じフォルダ配下の `logs/YYYY-MM-DD.log` を確認する。

- `ui wails run begin`: Wails起動処理へ入っている。
- `ui startup begin`: Go側のWails startupへ到達している。
- `ui lifecycle dom_ready`: WebViewのDOM準備まで進んでいる。
- `ui action="frontend_script_loaded"`: frontendのJavaScript assetが読み込まれている。
- `api GetInitialState begin`: Vue側からGo API呼び出しを開始している。
- `api GetInitialState complete`: 初期状態取得が完了している。
- `ui action="frontend_error"`: JavaScript実行時エラーが発生している。
- `ui lifecycle before_close`: ユーザー操作などで終了処理に入っている。
- `ui shutdown complete`: shutdownが完了している。

`frontend_script_loaded api=available` の後に `frontend_error` が出る場合は、Wails APIとHTML表示までは到達しているため、まずVue/JavaScript実行時エラーを疑う。

## よくある原因: template literal内の生バッククォート

`src/frontend/src/main.js` はVue templateを `template: \`...\`` のJavaScript template literalとして書いている。この中に説明文として生のバッククォートを入れると、template literalが途中で閉じられ、後続の文字列がJavaScriptとして評価される。

例:

```js
template: `
  <p>Debug OSC Pingは `/avatar/parameters/avatar_beacon/debug/ping` を送信します。</p>
`
```

この例では `/avatar/...` が意図せずJavaScript式として評価され、起動時に `Uncaught ReferenceError: avatar is not defined` のようなエラーになることがある。ビルドは通る場合があるため、runtimeで初めて発覚する。

## この原因であるという同定方法

次の条件が揃う場合、このパターンを強く疑う。

1. 起動画面に `フロントエンドエラー: Uncaught ReferenceError: ... is not defined` が表示される。
2. ログに `frontend_script_loaded api=available` が出た直後、`frontend_error` が出ている。
3. エラーURLが `http://wails.localhost/assets/index-*.js:行:列` を指している。
4. `src/frontend/src/main.js` の `template: \`...\`` 内に説明用の生バッククォートがある。

確認コマンド:

```bash
node scripts/check-frontend-template-literals.mjs
rg -n "\x60" src/frontend/src/main.js
```

`rg` の結果を見るときは、Vue templateの開始 `template: \`` から閉じバッククォートまでの範囲に、説明文用のバッククォートが混ざっていないか確認する。

## 対処方法

Vue template内では、説明文用の生バッククォートを使わない。

- コード風に見せたいOSC addressやパスは `<code>...</code>` で書く。
- 単なる強調なら日本語の鉤括弧や通常テキストにする。
- 長い説明文や動的テキストは、必要に応じて `data`、`computed`、または別コンポーネントへ逃がす。

修正後は最低限、次を実行する。

```bash
node scripts/check-frontend-template-literals.mjs
node scripts/check-wails-api-surface.mjs
npm run build
cd src && GOCACHE=/tmp/clipforvrchat-go-cache go test ./...
```

Release workflowでは `Check frontend template literals` が通っていることも確認する。

## UI表示不具合時の基本順序

1. 最新の `logs/YYYY-MM-DD.log` を読む。
2. `ui wails run begin`、`startup begin`、`dom_ready`、`frontend_script_loaded`、`GetInitialState begin/complete`、`frontend_error` の順に到達点を確認する。
3. 起動画面に表示された進捗またはエラー文を確認する。
4. `frontend_error` がある場合は、ビルド済みassetの行番号と `src/frontend/src/main.js` の該当箇所を対応させる。
5. `frontend_error` がなく `GetInitialState begin` で止まる場合は、Go APIまたはbackend初期化を疑う。
6. `before_close` が出ている場合は、手動終了か自動終了かを時刻と操作履歴で切り分ける。

この手順で新しい既知パターンが見つかった場合は、このファイルへ原因、同定方法、対処方法を追記する。
