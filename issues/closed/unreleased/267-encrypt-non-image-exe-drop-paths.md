# exeへD&Dされた画像以外のパスをGPG暗号化する

## 指示

> ClipForVRChatのexeにD&Dされたときの動作の仕様を変更
> あたえられたデータの拡張子で処理を分岐するが、
> 現在の処理は
> 画像ならClipForVRChatの入力画像とする。
> zipならgpg暗号化する。
> と思われる。
>
> これを
> - 画像ならClipForVRChatの入力画像とする。
> - それ以外のファイルやフォルダはgpg暗号化する。
> に変更。
>
> 実装と仕様書に反映

## 文脈

`ClipForVRChat.exe` へファイルをドラッグ&ドロップすると、Windowsでは渡されたパスが起動引数として扱われる。現状は単一 `.zip` 引数だけをGPG暗号化し、それ以外の位置引数は画像処理へ渡している。

## 解釈

exeへ渡されたパスは、画像ファイルなら従来どおりClipForVRChatの入力画像として扱う。画像ではない通常ファイルやフォルダは、zip拡張子に限らず公開鍵GPG暗号化対象として扱う。フォルダは直接GPGへ渡せないため、フォルダ内容をzip化したデータを暗号化する。

## 問題

- `.zip` 以外のファイルやフォルダをexeへD&Dすると画像処理へ回り、画像読み込みエラーになる。
- 不具合報告や任意データの暗号化入口として、zip以外のデータを扱えない。

## 期待する挙動

- 画像ファイルは従来どおり画像処理へ渡される。
- 画像以外のファイルは、そのファイル内容をGPG暗号化して `<元パス>.gpg` を作成する。
- フォルダは、フォルダ内容をzip化したデータをGPG暗号化して `<フォルダパス>.zip.gpg` を作成する。

## 実装

- 起動引数を画像パスと暗号化対象パスへ分岐し、画像以外のファイルやフォルダは単一起動制御やUI起動より前にGPG暗号化する。
- 画像扱いは既存デコーダに合わせて `.png`、`.jpg`、`.jpeg`、`.webp` とする。
- フォルダはメモリ上でzip化し、平文zipを成果物として残さず `<フォルダパス>.zip.gpg` を作成する。
- Windows GUI exeのコンソール出力判定も、画像以外の暗号化時に結果やエラーを表示できるよう更新する。
- `src/SPEC.md` と `README.md` にexe D&D時の分岐を記載する。

## 確認

- `TMPDIR=/tmp/clipforvrchat-go-build GOCACHE=/tmp/clipforvrchat-go-cache XDG_CACHE_HOME=/tmp/clipforvrchat-cache go test . -run 'TestEncryptFileWithPublicKey|TestEncryptDirectoryWithPublicKey|TestSplitStartupPathsRoutesImagesAndEncryptionTargets|TestHandleCLIArgs'`
- `TMPDIR=/tmp/clipforvrchat-go-build GOCACHE=/tmp/clipforvrchat-go-cache GOOS=windows GOARCH=amd64 go test -c -o /tmp/clipforvrchat-main.test.exe .`
- `git diff --check -- README.md issues/closed/README.md issues/closed/unreleased/267-encrypt-non-image-exe-drop-paths.md src/SPEC.md src/cli_console_windows.go src/main.go src/diagnostic_package.go src/main_test.go src/app_test.go`

## 受け入れ条件

- [x] exeへ画像ファイルだけが渡された場合は、画像処理が実行される。
- [x] exeへ画像以外のファイルが渡された場合は、UIを起動せずGPG暗号化ファイルが作成される。
- [x] exeへフォルダが渡された場合は、フォルダ内容を含むGPG暗号化ファイルが作成される。
- [x] 仕様書にexe D&D時の画像/非画像/フォルダ分岐が記載されている。
- [x] 関連テストが通る。
