# b10で3枚撮影想定なのに6枚投稿される

## 指示

> '/mnt/c/Users/user/Downloads/ClipForVRChat-v0.1.8-b10-windows-amd64/history.json''/mnt/c/Users/user/Downloads/ClipForVRChat-v0.1.8-b10-windows-amd64/logs' 3枚とる想定で6枚投稿される

> 二重に投稿される不具合あって抑制処理入れたと思うんだけどフォールバックモードで自動撮影した場合ってちゃんと処理走る？

> photoでどうなるかコードからチェック

> 257に対処したい

## 文脈

`v0.1.8-b10` の実行結果で、3枚撮影される想定の自動撮影がDiscordへ6枚投稿された。ユーザーから実行環境の `history.json` と `logs/` が提示された。

## 解釈

自動撮影バッチ内で撮影処理またはDiscord投稿処理が重複している可能性がある。提示された履歴とログから、撮影枚数自体が6枚なのか、3枚の撮影結果が二重投稿されたのかを切り分け、原因を修正する。

## 問題

- 3枚撮影想定に対してDiscord投稿が6枚発生する。
- 重複が撮影側、履歴登録側、Discord投稿側のどこで起きているか未確定。

## 期待する挙動

- 3構図の自動撮影では、Discord投稿も3件だけ発生する。
- 同じ自動撮影結果が重複投稿されない。

## 受け入れ条件

- [x] 提示された `history.json` と `logs/` から重複の発生箇所を特定する。
- [x] 自動撮影のDiscord投稿が重複しないよう修正する。
- [x] 関連テストを追加または更新する。
- [x] 関連チェックが通る。

## 追加受け入れ条件

- [x] Photo方式の自動撮影で作成されたVRChat写真を、通常のVRChat写真自動処理が二重投稿しない。
- [x] フォールバックモード時も同じ抑制が効く。
- [x] 関連テストを追加または更新する。
- [x] 関連チェックが通る。

## 調査

- `history.json` には6件あり、3件は `outputPath` が `VRC-AutoCapture` の元画像、残り3件は通常処理の `output/*_2048.png` だった。
- 通常ログでは自動撮影バッチ `batch-20260707-155241` の成功は3件で、各画像は自動撮影側でDiscord投稿済みだった。
- `autoPhoto.photoDirectory` が `%USERPROFILE%\Pictures\VRChat`、`autoCapture.output.directory` がその配下の `%USERPROFILE%\Pictures\VRChat\VRC-AutoCapture` だったため、VRChat写真自動投稿が自動撮影の出力画像を新規写真として拾い、同じ画像を再投稿していた。

## 確認

- `cd src && GOCACHE=/tmp/clipforvrchat-go-cache go test ./internal/appcore -run 'TestAutoPhotoWatcher|TestScanPhotoFiles'`
- `cd src && GOCACHE=/tmp/clipforvrchat-go-cache go test ./...`
- `git diff --check`

## 追加確認

- 2026-07-07: フォールバックモード時の自動撮影出力が、VRChat写真自動処理の重複取り込み抑制対象になるかコードで確認する。
- `restartAutoPhotoWatcher()` はVRChat写真自動処理の `AutoPhotoWatcher` に `ExcludeDirectories: []string{cfg.AutoCapture.Output.Directory}` を渡している。
- `AutoPhotoWatcher` は `scanPhotoFilesWithExcludesStatus()` で除外ディレクトリを `WalkDir` から `SkipDir` する。
- Stream方式の自動撮影はフォールバックモードでも `autoCaptureOutputPath()` で `AutoCapture.Output.Directory` へ保存するため、VRChat写真自動処理側の二重取り込み抑制が効く。
- Photo方式の自動撮影はVRChat写真フォルダに保存された写真を検出して処理するため、`AutoCapture.Output.Directory` 除外とは別経路になる。フォールバックモードとPhoto方式を組み合わせる場合、VRChat写真自動処理も同時ONなら別途重複確認が必要。
- 確認: `go test ./internal/appcore -run 'TestAutoPhotoWatcherExcludesAutoCaptureOutputDirectory|TestPreplacedLocalAnchorSkipsPoseResolve'`
- 2026-07-07: Photo方式のフォールバックモード時に、VRChat写真自動処理の二重取り込み抑制が効くかコードから確認する。
- Photo方式は `/usercamera/Capture` 後に `AutoPhoto.PhotoDirectory` の新規写真を `waitForNewPhoto()` で検出し、その `photoPath` を `finalizeAutoCaptureImage()` へ渡して自動撮影側でDiscord投稿する。
- フォールバックモードの `preplaced_local_anchor` はカメラPose/Options送信をスキップするだけで、Photo方式の検出対象ディレクトリや `photoPath` は変えない。
- VRChat写真自動処理の `AutoPhotoWatcher` は `ExcludeDirectories` として `AutoCapture.Output.Directory` のみ受け取り、`AutoPhoto.PhotoDirectory` 直下の写真は除外しない。
- そのためPhoto方式でVRChat写真自動処理も有効な場合、同じVRChat写真を自動撮影側で投稿した後、通常のVRChat写真自動処理が再処理してDiscordへ投稿し得る。Stream方式向けに入れた出力ディレクトリ除外だけではPhoto方式の二重投稿抑制にならない。

## 追加対応

- `AutoPhotoWatcher` に `ShouldSkip` を追加し、外部から抑制された写真を処理済み扱いにできるようにした。
- Photo方式の自動撮影バッチ中は通常のVRChat写真自動処理を一時抑制し、自動撮影ランナーが検出した写真パスも予約してバッチ直後の監視tickで二重投稿しないようにした。
- フォールバックモードは同じPhoto方式の検出経路を通るため、同じ抑制が適用される。

## 追加対応後の確認

- `cd src && go test ./internal/appcore -run 'TestAutoPhotoWatcher'`
- `cd src && TMPDIR=/tmp GOCACHE=/tmp/clipforvrchat-go-cache go test . -run 'TestAppAutoCapturePhotoSuppression|TestAppStartupStartsAutoPhotoWatcherWithoutDiscordUpload'`
- `cd src && go test ./...`
- `git diff --check`
