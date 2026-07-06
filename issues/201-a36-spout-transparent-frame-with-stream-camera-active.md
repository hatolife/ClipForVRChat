# a36でStream Camera起動済みでもSpout取得が透明フレームになる

## 問題

`v0.1.8-a36` のテスト撮影で、VRChat内のStream Cameraを起動しSpoutもONにしているにもかかわらず、Stream方式のSpout取得が次のエラーで失敗する。

```text
Spout取得に失敗しました: Spoutフレームは取得できましたが、ほぼ透明です。VRChat Stream Cameraの映像ではない可能性があります。
```

ユーザー提供ログ:

- `/mnt/c/Users/user/Downloads/ClipForVRChat-v0.1.8-a36-windows-amd64/logs/2026-07-05.log`

ログ上は `VRCSender1` senderを検出し、`1920x1080` のSpoutフレームも受信している。
ただし helper JSON は `frame=0`, `firstFrame=0`, `lastReceivedFrame=0`, `receiveAttempts=265`, `receiveSuccesses=264`, `transparentRatio=1.0` で、timeout内にframe番号が進まず全透明フレームしか取得できていない。

## 期待する挙動

- VRChat Stream Cameraが起動済みでSpout ONの場合、有効な映像フレームを取得できる。
- 透明フレーム継続の場合、VRChat senderが更新されていないのか、helperの受信形式/初期化順序が悪いのか、Stream Cameraの状態が不正なのかをログで切り分けられる。
- 透明フレームを保存成功扱いせず、必要な復旧操作や次に見るべきログが分かるエラーにする。

## 受け入れ条件

- [x] a36ログを確認し、sender検出済み・receive成功・frame番号不変・全透明で失敗していることを記録する。
- [ ] `spout-capture` helperの受信形式、frame counter、receiver初期化、透明判定を再調査する。
- [x] frame番号が進まない全透明フレーム継続時のエラー文を、単なる「映像ではない可能性」より原因切り分けしやすくする。
- [ ] 可能なら、受信前後の待機/再初期化/別受信形式などで有効フレーム取得率を改善する。
- [x] ローカルテストまたは静的検証可能な範囲で、透明フレーム分類を確認する。

## 実装メモ

- 2026-07-05: a36ログでは `VRCSender1` を検出し、`receiveAttempts=265` / `receiveSuccesses=264` まで進んでいる。一方で `frame=0` / `firstFrame=0` / `lastReceivedFrame=0` / `transparentRatio=1.0` で、sender frame番号が進まない全透明受信で失敗していた。
- 2026-07-05: Go側のエラー分類で、全透明かつsender frame番号不変の場合は「senderは見つかり受信も成功しているが、frame番号が進まず全透明」と分かる文言へ分岐するようにした。
- helperの受信形式や再初期化による改善は、Windows/VRChat実機での追加確認が必要。
