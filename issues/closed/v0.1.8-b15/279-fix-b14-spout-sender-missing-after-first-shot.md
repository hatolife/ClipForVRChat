# b14で1枚目成功後に2枚目以降のSpout取得が失敗する

## 指示

> '/mnt/c/Users/user/Downloads/ClipForVRChat-v0.1.8-b14-windows-amd64/logs'
> １枚目成功後何かいらない子として２枚目以降失敗してる気がする

> 解決するまで継続して作業 betaどんどん作っていい　C:\Users\user\Downloads　にci生成物ダウンロードして実行して確認　configは'/mnt/c/Users/user/Downloads/ClipForVRChat-v0.1.8-b14-windows-amd64/config.json'をコピーして使用

## 文脈

- `v0.1.8-b14` の自動撮影ログでは、1枚目のStream/Spout撮影は成功している。
- 同じバッチ内の2枚目以降でSpout sender一覧が空になり、`sender_not_found` で失敗している。
- ユーザーは1枚目成功後に何かが不要扱いされ、2枚目以降の撮影に影響している可能性を示唆している。

## 解釈

1枚目のSpout取得成功後、Spout senderまたはStream Cameraの状態を維持できず、後続の構図撮影でSpout senderを再検出できない不具合として調査・修正する。

## 問題

- 1枚目の撮影後に2枚目以降の `auto-capture spout list` が `senders=0` になる。
- recoveryで `/usercamera/Streaming=true` を再送してもsenderが復帰せず、後続撮影が失敗する。

## 期待する挙動

- 同一バッチ内で1枚目成功後もStream Camera/Spout senderを維持し、2枚目以降も連続して撮影できる。
- senderが一時的に消えた場合でも、実際に復帰を待ってから撮影するか、原因が分かるログを出す。

## 受け入れ条件

- b14ログの失敗原因を実装上の挙動と照合して説明できる。
- 必要な修正を入れ、連続撮影時に1枚目成功後のsender消失で即失敗しない。
- 関連する自動テストまたは静的検証を実行して結果を記録する。

## 対応メモ

- 2026-07-08: b14ログを確認。1枚目は `spout sender recovery` のON再送後に `senders=1` となり `VRCSender1` の取得に成功。2枚目以降は1枚目成功後の次回 `spout list` が `senders=0` となり、`open_before_batch=false` のためOFF/ONではなくON再送のみ実施、その後も `senders=0` のまま `sender_not_found` で失敗していた。
- 2026-07-08: `openCameraBeforeBatch=false` でも、ON再送後の再確認でsenderが戻らない場合は `/usercamera/Streaming=false` → `true` のトグルへ段階的に上げるよう修正した。ONだけで復帰した場合は従来どおりOFFを送らない。
- 2026-07-08: 検証: `go test ./internal/appcore -run 'TestRecoverEmptySpoutSenderList|TestAutoCaptureRunnerRunOnceSkipsCameraAutoOpenWhenDisabled'`、`go test ./...`。実機VRChatでの連続撮影確認は未実施のため要確認。
- 2026-07-08: ユーザー指示に基づき、解決確認までbeta作成とCI生成物の実行確認を継続する。`v0.1.8-b15` を本修正の確認betaとして扱う。
- 2026-07-08: `v0.1.8-b15` タグのRelease workflowとbranch CIが成功。Release asset `ClipForVRChat-v0.1.8-b15-windows-amd64.zip` を `C:\Users\user\Downloads` へダウンロードし、SHA256 `bc321c27fcf2180d9aaf9628ea18097314b4389ea47776560c1776f46b4d6a4f` がGitHub Release asset digestと一致することを確認した。
- 2026-07-08: `C:\Users\user\Downloads\ClipForVRChat-v0.1.8-b15-windows-amd64` へ展開し、b14の `config.json` をコピーして起動確認した。開始時撮影はAvatarBeacon basis未readyで一度skipしたが、その後readyになり次回スケジュールで3構図撮影が実行された。
- 2026-07-08: 実機ログで、正面/背後/斜めの全3枚が `senders=0` → ON再送後も `senders=0` → `spout sender recovery escalate` → OFF/ONトグル後 `senders=1` → `spout capture success` へ進み、`batch complete: results=3 errors=0` を確認した。b14で失敗していた2枚目以降も成功したため完了とする。
