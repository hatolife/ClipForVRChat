# 実機確認を減らす自動検証基盤を追加する

## 指示

> いけるところまで進めてほしい
> サブエージェントも使って並列に進められるなら進めてほしい

> https://github.com/xcd0/ClipForVRChat
>
> これ新規で作ったからここにpush

> 元のリポジトリにpush

> releaseで確認できる状態に

## 文脈

自動撮影はOSC、AvatarBeacon、Spout helper、Wails GUI、Discordなど複数の外部境界をまたぐ。現状は各部品の単体テストやWindows buildがある一方、既存のSpout C++テストがCIで実行されず、Frontendの自動テストと実Wails GUIの起動確認も不足している。

旧リポジトリへのGitHub接続は読み取り専用だったため、書き込み可能な `xcd0/ClipForVRChat` を新しい公開先として作業を継続する。

その後、元の `hatolife/ClipForVRChat` への書き込み権限が復旧したため、同じ検証済み差分を元リポジトリの作業ブランチへ反映する。

元リポジトリの `develop` へ取り込んだ後、`v0.1.8-rc4` のpre-releaseを作成し、Windows成果物を取得して確認できる状態にする。

## 解釈

実VRChat、実GPU、実アバターでしか確認できない境界を残し、その直前までをCIで検証できるようにする。まず既存資産を使える低リスクな層から追加し、標準Windows runnerでの成立性が未確認なGUI/Spout実動作は観測可能なsmoke testまたは後続スパイクとして分離する。

## 問題

- `spout-capture-logic-test` がCMake/CTestへ登録済みだがCIで実行されていない。
- Frontendの保存前確認条件などをブラウザやWails本体から独立してテストしにくい。
- AvatarBeacon受信、basis確定、OSCログ、forwardを一連の実UDPシナリオとして検証していない。
- Windows CIで実Wails exeがWebView2を生成できるか自動確認していない。
- Real Spout E2EはGPUおよび共有texture環境への依存があり、標準runnerでの成立性が不明。
- 新しい公開先でCIを実行したところ、既存のFrontend依存とGo toolchain/moduleに公開後の脆弱性が検出され、追加したCTestとGUI smokeへ到達する前に停止する。

## 期待する挙動

- Spout C++ロジックテストがWindows CIでbuild・実行される。
- Frontendの重要な純粋ロジックが依存追加なしで自動テストされる。
- 実UDPを使ったOSC統合シナリオがGoテストで再現される。
- 実Wails exeのWebView2 DOM ready、Frontend script読込、初期Wails API呼び出しをWindows CIから観測できる。
- 未確定のReal Spout E2EとGUI操作E2Eが、成功済み項目と混同されず残件として記録される。

## 受け入れ条件

- [x] CIが `spout-capture-logic-test` をbuildして `ctest` を実行する。
- [x] Frontend純粋ロジックのテストを `npm test` で実行できる。
- [x] 設定保存前の自動投稿確認条件を含むFrontendテストが通る。
- [x] AvatarBeacon受信からOSC forwardまでの実UDP統合テストが通る。
- [x] Windows CIで実Wails exeの `dom_ready`、`frontend_script_loaded`、`GetInitialState complete` を検出するsmoke testを実行する。
- [x] GUI smoke testは初回CI結果を確認するまで非blockingとし、結果と残件を記録する。
- [x] Real Spout E2Eの未確認範囲を明記する。
- [x] `npm audit` と `govulncheck` が到達可能な脆弱性なしで完了する。

## 確認

- ローカル: `cd src/frontend && npm test` で6件成功。
- ローカル: `node scripts/check-frontend-template-literals.mjs` 成功。
- ローカル: CI YAML parse成功。
- ローカル: `node scripts/check-commit-subjects.mjs develop..HEAD` 成功。
- ローカル: `git diff --check develop...HEAD` 成功。
- GitHub Actions run #9: Frontend依存導入、`npm audit`、Frontend test/build/static check、Wails API surface、`go test ./...` は成功。
- GitHub Actions run #9: `govulncheck` がGo 1.26.4の標準ライブラリと `golang.org/x/image v0.43.0` に到達可能な脆弱性を検出したため、Go patch versionとmoduleを更新して再確認する。
- GitHub Actions run #13: Go 1.26.7と `golang.org/x/image v0.45.0` への更新後、Frontend test/build/static check、`npm audit`、`go test ./...`、`govulncheck`、Spout helper build、CTest、埋め込みhelper抽出、Wails buildが成功。DevTools endpoint smokeは `continue-on-error` によりstep結論が成功表示になったが、ログ上はtimeoutしていた。
- GitHub Actions run #15: DevTools endpoint smokeをblocking化した結果、同じtimeoutをCI失敗として検出した。CDP endpointの有無をGUI起動成否と同一視せず、アプリ自身のlifecycle診断ログでDOM ready、Frontend script読込、初期Wails API完了を検証する方式へ変更する。
- GitHub Actions run #17: lifecycle診断ログ方式でもログ生成前にtimeoutした。GitHub-hosted Windows runnerが管理者権限で動き、通常版の管理者起動拒否ダイアログでWails開始前に待機するため、配布版には含まれない `ciguismoke` build tagで作るCI専用exeをsmoke対象にする。
- GitHub Actions run #19: Documentation checksとWindows build jobが成功。Frontend test/build/static check、`npm audit`、`go test ./...`、`govulncheck`、Spout helper build、CTest、埋め込みhelper抽出、Wails build、blocking GUI lifecycle smokeがすべて成功した。

## 残件

- synthetic Spout senderから既知画像を送信し、helper出力PNGのpixel/metadataを検証する。
- 標準Windows runnerで共有textureが成立しない場合はWindows GPU runnerまたはself-hosted runnerを評価する。
- Browser + mock Wails APIのGUI操作E2Eと、実WebView2へCDP接続する少数シナリオを追加する。
