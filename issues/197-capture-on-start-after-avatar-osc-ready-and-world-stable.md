# 開始時に撮影を初期ONにし、OSC基準確定後かつワールド移動中でない時だけ実行する

## 問題

自動撮影スケジュールの `開始時に撮影` は、自動撮影開始直後に1回撮影するための設定だが、現在の初期値はOFFになっている。

AvatarBeaconを前提にした `player_local` 構図では、撮影直後の位置と向きの正確さが重要であり、開始時撮影はOSCでプレイヤー基準座標が確定したタイミングを基準に実行したい。

一方で、ワールド移動中やワールド読み込み直後に撮影すると、Avatar OSC basis、User Camera、Spout sender、Presence/world metadataが不安定な状態で撮影される可能性がある。

現在、VRChat公式OSCで「現在いるworld ID」を直接取得できる標準addressは確認できていない。

## 期待する挙動

`開始時に撮影` の初期値をONにする。

ただし、自動撮影開始直後に即撮影するのではなく、次の条件を満たしてから開始時撮影を実行する。

- AvatarBeacon由来の `avatar_osc` basisがreadyである。
- 必要なOSC値の鮮度が十分である。
- ワールド移動中ではない。
- ワールド/インスタンス情報が取得できる場合は、撮影直前のworld metadataとして安定している。

ワールド移動中は開始時撮影を抑止または延期する。

## 受け入れ条件

- [x] `AutoCaptureScheduleConfig.CaptureOnStart` の新規初期値をONにする。
- [x] 既存ユーザー設定の移行方針を決める。
  - 既存configで明示的にOFFにしている場合は尊重する。
  - 古いconfigで未設定の場合だけONへ移行する案を優先する。
- [x] 開始時撮影は、自動撮影開始直後ではなく `avatar_osc` basis ready後に実行する。
- [x] `avatar_osc` basisがreadyにならない場合のtimeout、エラー表示、ログを決める。
- [x] ワールド移動中は開始時撮影を実行しない。
- [x] ワールド移動中の検出方法を実装する。
- [x] ワールド移動完了後、必要なら開始時撮影を延期実行する。
- [x] ワールド移動中に無効化された/延期されたことが診断ログで分かる。
- [x] UI上で「開始時に撮影」はOSC基準確定後に実行されることが分かる文言にする。
- [x] `node scripts/check-frontend-template-literals.mjs` と `node scripts/check-wails-api-surface.mjs` を実行する。
- [x] `cd src && GOCACHE=/tmp/clipforvrchat-go-cache go test ./...` を実行する。

## world ID / ワールド移動検出メモ

- 2026-07-04時点で確認したVRChat公式OSC docsでは、OSC OverviewはInputとAvatar Parametersを主なAPIとして説明している。
- 公式OSC Avatar Parameters docsでは、VRChatから外部へ送られる標準的な通知として `/avatar/change` によるavatar ID通知が説明されているが、現在world ID通知は確認できていない。
- 公式OSC Avatar Parameters docsでは、Avatar Parameterの出力はparameter値が変化したときに送信される説明であり、world IDを標準parameterとして外部へ送る仕様は確認できていない。
- VRChat写真メタデータには `vrc:WorldID` が入るリリースノートがあるが、これは撮影後の写真metadataであり、開始時撮影の前にワールド移動中かを判定する用途には使いにくい。
- 既存実装には `SnapshotVRChatWorld()` があり、VRChat output logから直近の `wrld_...` world/instance tokenを取得している。まずはこの経路でワールド移動/安定判定を設計する案が現実的。
- 公式に現在world IDを取得できるOSC/OSCQuery/APIが確認できた場合は、それを優先する。
- 公式手段がない場合、アバターギミックでworld IDを取得できるか調査する。ただし、アバターギミック単体でworld ID文字列を知る手段は限定的な可能性があるため、実現性を先に確認する。

参考:

- https://docs.vrchat.com/docs/osc-overview
- https://docs.vrchat.com/docs/osc-avatar-parameters
- https://docs.vrchat.com/docs/vrchat-202531

## 実装メモ

- 2026-07-07追記: ユーザー指示「開始時に撮影を規定値でONに」に基づき、backend既定値だけでなくfrontend側の欠損補正でも `captureOnStart` をONにする。
- `CaptureOnStart` の初期値変更は、単純なdefault変更だけだと既存configの明示OFFと区別できない可能性がある。
- 現状の `AutoCaptureScheduleConfig.CaptureOnStart` は plain `bool` で、未設定と明示 `false` をJSON上で区別できない。
  - 定義: `src/internal/appcore/config.go`
  - default: `DefaultConfig()` / `DefaultAutoCaptureConfig()` 周辺
  - load/save: `LoadConfig()` / `SaveConfig()`
  - 既存の `Normalize()` は `CaptureOnStart` を補正していない。
- 既存の `LoadConfig()` は `DefaultConfig()` を土台に `json.Unmarshal` しているため、`CaptureOnStart` の default を ON に変更すれば、未設定は ON、明示 `false` は false のまま扱える。今回の最小実装ではこの経路を優先する。
- 現在のスケジューラは `src/app.go` の `runAutoCaptureScheduler()` で `InitialDelaySec` 後に `CaptureOnStart` を見て即 `runAutoCaptureBatch()` を呼ぶ。
- バッチ直前の `prepareAutoCaptureConfigForRunLocked()` と `freshPlayerLocalBasisLocked()` は、`avatar_osc` basisがready/freshでない場合にエラーにできる。開始時撮影では、ここで失敗させるだけでなく、初回専用のready待ちを挟む最小実装を優先する。
- `latestAvatarOSCBasisSnapshotLocked()` は `ready/stale/missing/partial/invalid`、鮮度、最終受信parameter、raw countを持つため、開始時撮影のready待ちに再利用できる。
- `avatar_osc` basis ready待ちは、既存のAvatar OSC受信状態判定とfreshness設定を再利用する。
- ワールド移動中判定は候補を比較する。
  - VRChat output logのworld join行を監視する。
  - Presence user集合リセットの既存処理を使う。
  - world ID/instance IDの変化から一定秒数は不安定期間として扱う。
  - Avatar OSC basisが一定時間readyになったことを移動完了の補助条件にする。
- 今回は world 安定判定を大規模化しすぎず、開始時撮影の待機ヘルパー内で `SnapshotVRChatWorld()` を使った診断・猶予判定から始める。
- `WatchOutputLog` がONの場合は、VRChat output logから `WorldID` / `InstanceID` が取得でき、同じ値が3秒以上維持されるまで開始時撮影を待つ。
- `WatchOutputLog` がONでworld metadataを取得できない場合は、危険なタイミングで即撮影せずtimeout時に開始時撮影をスキップし、診断ログとUIメッセージに理由を出す。
- `WatchOutputLog` がOFFの場合は、world安定待ちは明示的に無効とみなし、`avatar_osc` basis ready/freshだけを開始時撮影条件にする。
- ワールド移動中に自動撮影スケジュールが動いている場合、開始時撮影だけでなく定期撮影も抑止すべきかは別途判断する。

## 調査メモ（2026-07-04）

- `parseVRChatWorldMetadata()` は最新output log内の最後の `wrld_...` tokenから `WorldID` と `InstanceID` を取り出す。
- `parseVRChatPresenceLog()` は `Entering Room` / `Joining wrld_` で現在ユーザー集合をリセットするため、ワールド移動シグナルとして使える。
- ただし現状はbatch開始時のsnapshot読み込みであり、output logのtail監視やfsnotifyによる継続的な「移動中」状態管理はない。
- 実装案:
  - `waitForStartReadiness(ctx, cfg)` のような初回専用helperを追加する。
  - `avatar_osc` basisがready/freshになるまで短周期で待つ。
  - output log由来のworld/instance tokenが一定秒数変化しないことを確認する。
  - worldが不安定、またはtimeoutした場合は開始時撮影をスキップ/延期し、診断ログへ理由を出す。
- 関連テスト候補:
  - `src/internal/appcore/config_test.go` に新規config/既存config移行のテストを追加する。
  - `src/app_test.go` に開始時撮影ready待ちのsuccess/timeout/一回だけ実行のテストを追加する。
  - `src/internal/appcore/autocapture_test.go` に複数world tokenや `Joining wrld_` 検出のテストを追加する。

## 検証観点

- 新規configで `開始時に撮影` がONになる。
- 既存configで明示OFFにしていた場合、勝手にONへ変わらない。
- 自動撮影開始後、AvatarBeacon basisがreadyになるまで開始時撮影しない。
- AvatarBeacon basis ready直後に開始時撮影が1回だけ実行される。
- VRChatワールド移動中は開始時撮影しない。
- ワールド移動完了後に、延期された開始時撮影が必要に応じて1回だけ実行される。
- world ID取得に失敗しても、診断ログで原因が分かり、危険なタイミングの撮影を避ける。
