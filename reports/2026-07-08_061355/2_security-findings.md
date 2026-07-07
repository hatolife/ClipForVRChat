# セキュリティ findings

## SEC-20260708-001

- ID: SEC-20260708-001
- タイトル: Spout helper 解決がPATH上の同名実行ファイルへfallbackする
- 重大度: Medium
- CWE: CWE-427 Uncontrolled Search Path Element
- 該当ファイル: `src/internal/appcore/spout.go`
- 該当箇所: `ResolveSpoutHelperPath`, `exec.LookPath(helperPath)`
- 説明: `spoutHelperPath` が空または `spout-capture.exe` の場合、埋め込みhelper、exe隣接helperの順に優先される。しかし埋め込みhelperがない、または分離版で隣接helperが欠けている場合、最後に `exec.LookPath("spout-capture.exe")` が呼ばれる。PATH上に同名の不正実行ファイルがある環境では、ユーザーが意図しないhelperを起動する可能性がある。
- 攻撃者視点: ユーザーのPATHまたはアプリ起動環境に同名実行ファイルを配置できる攻撃者は、Spout確認または自動撮影時に任意のローカル実行ファイルを起動させる余地がある。
- 影響: 非管理者権限のローカルコード実行。設定、画像、環境情報などへのアクセスはそのプロセス権限に依存する。
- 再現方法: 防御的確認として、埋め込みhelperを無効化したテストでexe隣接helperがない状態を作り、PATH上の `spout-capture.exe` が解決されるかを単体テストで確認する。
- 修正案: default/helper名指定時は、埋め込みhelperまたはexe隣接helperのみを許可し、PATH fallbackを廃止する。カスタムhelperを許可する場合は絶対パスまたは明示選択済みパスに限定し、UIで信頼できるhelperのみ指定する注意を出す。
- 修正例: `exec.LookPath(helperPath)` へ進む前に、`helperPath == spoutHelperFileName` の場合はエラーを返す。任意名のPATH検索も原則禁止する。
- 修正後確認方法: `ResolveSpoutHelperPath("spout-capture.exe")` が埋め込み/隣接なしでPATH上のhelperを返さずエラーになるテストを追加する。
- 信頼度: High
- 備考: 通常利用者向け単一exeでは埋め込みhelperが優先されるため、主な影響は分離版、不完全な展開、開発/検証環境に寄る。

## SEC-20260708-002

- ID: SEC-20260708-002
- タイトル: Windows固定コマンド起動がPATH検索へ依存している
- 重大度: Low
- CWE: CWE-427
- 該当ファイル: `src/startup_shortcut_windows.go`, `src/app.go`, `src/internal/appcore/autocapture.go`
- 該当箇所: `exec.Command("powershell.exe", ...)`, `exec.LookPath("winget")`, `exec.CommandContext(..., "winget", ...)`, `exec.CommandContext(..., "powershell.exe", ...)`
- 説明: Startup shortcut作成、shortcut読取、VRChatウィンドウ矩形取得、ffmpegインストール補助で `powershell.exe` または `winget` を名前だけで起動している。Windowsの実行ファイル検索順やPATHが汚染された環境では、意図しない同名実行ファイルを起動するリスクが残る。
- 攻撃者視点: ユーザーのPATH、作業ディレクトリ、または起動環境へ同名実行ファイルを置ける場合、該当機能の実行時に不正プログラムを起動させる余地がある。
- 影響: 非管理者権限のローカルコード実行。`winget install ffmpeg` はユーザー操作を伴う補助機能であり、既定自動実行ではない。
- 再現方法: 防御的確認として、該当関数をテスト可能にし、解決される実行ファイルパスがSystem32等の期待パスに固定されることを確認する。
- 修正案: `powershell.exe` は `%SystemRoot%\System32\WindowsPowerShell\v1.0\powershell.exe` または信頼済みの絶対パスを使用する。`winget` も解決済み絶対パスを検証してから同じパスを起動する。
- 修正例: `trustedPowerShellPath()` / `trustedWingetPath()` を追加し、`exec.CommandContext` へ絶対パスを渡す。
- 修正後確認方法: PATH上に同名のテストファイルがある状態でも、期待パスだけが使われる単体テストを追加する。
- 信頼度: Medium
- 備考: アプリは管理者権限起動を拒否しているため、権限昇格には直結しにくい。

## SEC-20260708-003

- ID: SEC-20260708-003
- タイトル: GitHub Actions の外部Actionがmajor tag参照で固定されている
- 重大度: Low
- CWE: CWE-494 Download of Code Without Integrity Check
- 該当ファイル: `.github/workflows/ci.yml`, `.github/workflows/release.yml`
- 該当箇所: `uses: actions/checkout@v7`, `actions/setup-go@v6`, `actions/setup-node@v6`, `actions/upload-artifact@v6`, `actions/download-artifact@v7`
- 説明: Release workflowはjob単位の権限分離、artifact digest検証、署名job分離が入っており改善されている。一方、外部Actionはmajor tag参照であり、tag移動や上流侵害時の供給網リスクを完全には排除できない。
- 攻撃者視点: 外部Actionの配布元やtagが侵害された場合、CI/Release環境で任意コードが実行され得る。
- 影響: Release成果物改ざん、秘密情報への到達、workflow実行環境の悪用。ただし現在は署名secretがsign jobに分離され、release jobだけ `contents: write` を持つ。
- 再現方法: workflowの `uses:` 参照がコミットSHAかどうかを静的確認する。
- 修正案: 外部ActionをコミットSHAへpinし、Dependabot/Renovate等で更新する。少なくともRelease workflowの署名/Release経路から優先する。
- 修正例: `actions/checkout@<full commit sha>` のように固定し、コメントで対応するtagを記録する。
- 修正後確認方法: `rg 'uses: .*@v[0-9]+' .github/workflows` がRelease経路で0件になることを確認する。
- 信頼度: High
- 備考: 運用負荷とのトレードオフがあるため、リリース前必須かはプロジェクト方針で判断する。

## SEC-20260708-004

- ID: SEC-20260708-004
- タイトル: ローカル監査では最新脆弱性DB照会を実施していない
- 重大度: Info
- CWE: 該当なし
- 該当ファイル: `.github/workflows/ci.yml`, `.github/workflows/release.yml`, `src/go.mod`, `src/frontend/package-lock.json`
- 該当箇所: `npm audit --omit=dev`, `govulncheck`
- 説明: 今回のプロンプトでは外部サービスへの能動的アクセスを避ける指示があったため、ローカルで `npm audit` や `govulncheck` のDB問い合わせは実行していない。CI/Release workflowには両チェックが組み込まれている。
- 攻撃者視点: 既知脆弱性DBの更新により新たな脆弱性が判明していても、今回のローカル監査だけでは検出できない。
- 影響: 依存関係の最新既知脆弱性の見落とし。
- 再現方法: ネットワーク許可環境またはGitHub Actionsで `npm audit --omit=dev` と `govulncheck ./...` を実行する。
- 修正案: Release前にCI結果を確認し、失敗時は依存更新または影響調査を行う。
- 修正例: なし。
- 修正後確認方法: CI/Releaseの該当step成功を確認する。
- 信頼度: High
- 備考: 今回のローカル `npm ls --omit=dev --all --package-lock-only` は依存ツリー確認のみであり、脆弱性DB照会ではない。

## SEC-20260708-005

- ID: SEC-20260708-005
- タイトル: Windows本番成果物の緩和機構と署名はローカル未確認
- 重大度: Info
- CWE: 該当なし
- 該当ファイル: `.github/workflows/release.yml`, `src/wails.json`
- 該当箇所: Wails Windows build, PGP detached signature, Release asset packaging
- 説明: ローカル環境はLinuxであり、Windows exeのASLR/DEP/CFG、SmartScreen表示、PGP署名、zip内容、Spout helper実行は未確認である。Release workflowには成果物検査と署名処理があるが、今回のローカル監査では実行していない。
- 攻撃者視点: 成果物検証がReleaseごとに抜けると、意図しないファイル混入や署名不備に気付きにくい。
- 影響: 配布物の真正性確認や利用者案内の信頼性低下。
- 再現方法: GitHub ActionsのRelease workflowを実行し、asset一覧、zip内ファイル、PGP署名、helper動作を確認する。
- 修正案: Release checklistに従い、ReleaseごとにCI/Release結果と署名検証を確認する。
- 修正例: なし。
- 修正後確認方法: `ClipForVRChat-vX.Y.Z-windows-amd64.zip` と `.exe.asc` の検証記録を残す。
- 信頼度: High
- 備考: 今回はRelease作業ではないため未実施。
