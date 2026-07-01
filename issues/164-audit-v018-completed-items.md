# v0.1.8完了扱い項目の再監査

## 問題

v0.1.8のissue一覧では多くの項目が完了扱いになっているが、実装・テスト・ドキュメント・Release同梱の観点で、完了扱いにできない不完全な項目が残っている可能性がある。

## 期待する挙動

v0.1.8完了扱いの各項目について、実装済み・未実装・仮実装・実機/CI未確認・仕様と実装の不一致を区別して把握できる。

## 受け入れ条件

- v0.1.8完了扱いのissueと実装を照合する。
- 不完全な項目は、根拠となるファイル/処理と不足内容を記録する。
- 完了扱いが妥当な項目と、要対応/要確認へ戻すべき項目を分ける。
- 調査結果をこのissueへ追記する。

## 再監査結果

2026-07-01時点の `feat/v0.1.8-resolve-issues` を確認した結果、以下は完了扱いにできないため `issues/README.md` で `要対応` に戻した。
この調査では実装コードの変更は行わない。

| issue | 判定 | 根拠 | 不足内容 |
| --- | --- | --- | --- |
| #121 | 要対応 | `src/internal/appcore/config.go` の初期構図は `world` のまま。 | 初期構図をプレイヤー中心ローカル視点として扱う要件が満たせていない。 |
| #129 | 要対応 | `tools/spout-capture/main.cpp` の複数sender失敗時は固定メッセージのみ。 | 自動選択不能時に候補一覧付きで失敗する受け入れ条件が未達。 |
| #130 | 要対応 | `src/internal/appcore/config.go` には `startDelayMs` があるが、`src/frontend/src/main.js` に設定UIがない。 | Stream起動後待機時間を自動撮影タブから設定できない。 |
| #131 | 要対応 | `src/internal/appcore/spout.go` の成功ログは処理時間を記録していない。 | Spout取得の処理時間診断、#129の未達、実機経路の確認が残る。 |
| #132 | 要対応 | `src/internal/appcore/spout.go` の `validateCapturedImage()` はDecodeConfigまで。 | ほぼ単色の白/黒画像検出と閾値診断が未実装。 |
| #133 | 要対応 | `src/internal/appcore/autocapture.go` のPhoto失敗メッセージに `ffmpeg入力設定` が残る。 | Stream方式の主経路がSpoutであるUI/エラー文言に同期できていない。 |
| #134 | 要対応 | `.github/workflows/release.yml` は展開済みzip内の必須ファイルを個別検証していない。 | Release zip内の `ClipForVRChat.exe`、`spout-capture.exe`、README、LICENSE、Spoutライセンス確認が不十分。 |
| #135 | 要対応 | README/SPECは概要のみ。 | 手動Camera/自動起動、複数構図、白画像/デスクトップ周辺時のログ収集手順が不足。 |
| #138 | 要対応 | `src/SPEC.md` の `player_local` 説明は最小限。 | 原点、アンカー、軸、Euler適用順、初期offset/注視点/ズームの詳細仕様が不足。 |
| #140 | 要対応 | `src/internal/appcore/player_local.go` はYawのみで位置変換し、回転は単純加算。 | PlayerBasis/アンカー表現、Yaw 0/90/180/-90、正面/背後/斜め、Pitch/Rollのテストが不足。 |
| #141 | 要対応 | `src/app.go` の `SaveCurrentCameraPoseToView()` と `AddCurrentCameraPoseAsView()` が `world` を強制する。 | `player_local` の現在Pose保存/追加、初期リセット、sidecar上のresolved pose区別が未達。 |
| #142 | 要対応 | README/SPECに実機確認手順が詳細化されていない。 | Desktop/VR、身長差、座り、移動/回転後、手動Camera/自動起動の確認手順が不足。 |
| #143 | 要対応 | `src/internal/appcore/metadata.go` のschemaは最小フィールド中心。 | app version、ユーザー件数、最大サイズ、省略ルール、プライバシー方針の定義が不足。 |
| #144 | 要対応 | `src/internal/appcore/metadata.go` はPNG `iTXt` のみで、拡張子で形式分岐する。 | PNG `eXIf`、マジックバイト判定、既存メタデータ重複方針、自前reader読み戻しが未実装。 |
| #145 | 要対応 | `finalizeAutoCaptureImage()` は埋め込み失敗時にResultを返し、sidecar/Discordへ進まない。 | 埋め込み失敗時も画像保存、sidecar、Discordを可能な範囲で継続する条件が未達。 |
| #147 | 要対応 | `metadata_test.go` はバイト列存在確認と画像decodeのみ。 | PNG/JPEGの自前reader読み戻し、Discord投稿後、ユーザー数過多時の検証手順が不足。 |
| #156 | 要対応 | `scripts/` と CI に frontend `api.*` と Go `App` 公開メソッドの照合がない。 | Wails API surfaceの静的チェックが未実装。 |
| #157 | 要対応 | README/SPEC/Release Notesが上記未達項目を実装済みとして読める。 | 実装済み範囲と制約のドキュメント同期が未完了。 |
| #158 | 要対応 | `parseVRChatWorldMetadata()` の正規表現は `~region(us)` の括弧以降を落とす。 | instance IDの完全保持とログ形式揺れのテストが不足。 |
| #163 | 要対応 | `src/frontend/src/main.js` はテスト撮影結果を成功/失敗メッセージだけで表示する。 | サムネイル、保存先、Discord投稿結果、sidecar確認導線が設定画面にない。 |

## 完了扱いを維持した代表項目

- #120: 初期Pose数値、拡大率、撮影トグル、初期Poseリセット、Wails API呼び出し自体は存在する。ただし `player_local` 座標対応の不足は #121/#141 に戻した。
- #125: 構図カード内の現在Pose追加ボタンとカメラ移動ボタンは配置済み。座標系対応不足は #141 に戻した。
- #126: User Camera関連OSCをOff/falseへ戻すリセットAPI/UIは存在する。
- #149: multi/Camera Dollyはv0.1.8対象外として通常導線から外されている。
- #151: `includeImages=false` は本文のみ投稿へ接続済み。`postMode` は `shot` に正規化される。
- #152: `captureDelayMs` はPhoto/Streamの撮影前待機へ接続済み。
- #153: `imageFormat` と `filenameTemplate` は自動撮影タブへ露出済み。
- #154: sidecar/Discord/画像埋め込みのユーザーID出力設定は分離済み。
- #155: ローカル削除/purge時に画像隣接sidecarを削除する処理とテストがある。
- #159: Discord payloadの `allowed_mentions.parse=[]` は実装・テスト済み。
- #160: SQLite/ローカルDBはv0.1.8対象外としてsidecar/history方式へ整理済み。
- #161: OSCQueryはv0.1.8対象外として延期明示済み。
- #162: 撮影間隔UIとNormalizeの最小値は整合済み。

## 補足確認

- Spout2の `2.007.015` tag は upstream に存在することを `git ls-remote --tags https://github.com/leadedge/Spout2` で確認した。
- Windows実機のSpout capture、GitHub Actionsの直近実行、Release zipの実物はこの調査では実行確認していない。
