# Review Log

## 実施日時

- 2026-07-08 06:13 JST頃。

## 確認したファイル・ディレクトリ

- `README.md`
- `RELEASE_NOTES.md`
- `.github/workflows/ci.yml`
- `.github/workflows/release.yml`
- `src/go.mod`, `src/go.sum`
- `src/frontend/package.json`, `src/frontend/package-lock.json`
- `src/main.go`, `src/app.go`
- `src/internal/appcore/*.go`
- `src/diagnostic_package.go`
- `src/single_instance*.go`
- `src/startup_shortcut_windows.go`
- `src/reveal_windows.go`
- `tools/spout-capture/*`
- `reports/2026-06-25/*`
- `reports/security/*`
- `issues/closed/unreleased/296-316` 周辺のセキュリティ対応チケット

## 実行したコマンド

- `find . -maxdepth 2 -type f`
- `git ls-files`
- `rg -n "TODO|FIXME|HACK|SECURITY|password|secret|token|key|credential|unsafe|system\\(|exec|CreateProcess|ShellExecute|LoadLibrary|WinExec|SearchPath|registry|pipe|socket|listen|bind|http|https|tmp|AppData|ProgramData|Program Files|ACL|service|COM|update|installer|download|signature|verify|hash|memcpy|strcpy|sprintf|reinterpret_cast|new |delete |malloc|free"`
- `npm run build`
- `node scripts/check-frontend-template-literals.mjs`
- `node scripts/check-vue-runtime-template.mjs`
- `node scripts/check-wails-api-surface.mjs`
- `cmake -S tools/spout-capture -B /tmp/clipforvrchat-spout-logic -DCMAKE_BUILD_TYPE=Release`
- `cmake --build /tmp/clipforvrchat-spout-logic --target spout-capture-logic-test`
- `ctest --test-dir /tmp/clipforvrchat-spout-logic --output-on-failure`
- `npm ls --omit=dev --all --package-lock-only`
- `GOPROXY=off GOSUMDB=off GOCACHE=/tmp/clipforvrchat-go-cache go test ./...`
- `GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off GOCACHE=/tmp/clipforvrchat-go-cache go test ./...`
- `[REDACTED: local user path]/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.26.4.linux-amd64/bin/go test ./...`
- 対応する `go vet ./...`

## 実行結果の概要

- `npm run build`: 成功。
- frontend template literal check: 成功。
- Vue runtime template check: 成功。
- Wails API surface check: 成功。46 frontend calls matched App methods。
- Spout capture logic C++ test: 成功。
- `npm ls --omit=dev --all --package-lock-only`: 依存ツリーを確認。optional dependency不足表示のみ。
- `go test ./...` / `go vet ./...`: ローカルGo toolchain不整合で未完了。`go version` は `go1.26.4` を表示するが、実行時に `go.mod requires go >= 1.26.4 (running go 1.26.1; GOTOOLCHAIN=local)` または `compile: version "go1.26.1" does not match go tool version "go1.26.4"` が出た。

## ビルド可否

- frontend buildは可。
- C++ helperの非Windowsロジックテストbuildは可。
- Windows Wails buildは未実施。

## テスト可否

- frontend関連チェックとC++ロジックテストは成功。
- Goテストは環境のGo標準ライブラリ/toolchain不整合で未完了。

## 静的解析の実施状況

- `rg` による危険API/secret/外部実行/ネットワーク/ファイル操作検索を実施。
- `go vet` はGo環境不整合で未完了。
- CodeQL、gosec、golangci-lint、clang-tidy、cppcheckは未実施。

## 動的解析の実施状況

- Windows実機動作、Spout helper実行、OSC実通信、Discord投稿、Release workflow実行は未実施。

## 確認できなかった項目

- 最新依存脆弱性DB照会。
- Windows binary mitigation確認。
- PGP署名検証とRelease asset一覧。
- SmartScreen表示。
- Spout/DirectX/OpenGL実機挙動。

## 追加調査が必要な項目

- `ResolveSpoutHelperPath` のPATH fallback廃止可否。
- `powershell.exe` / `winget` の絶対パス解決。
- GitHub Actions SHA pinningの更新運用。
