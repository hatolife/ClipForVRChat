# Issues

このディレクトリは、実装タスクや調査タスクを Markdown で管理するためのものです。

状態はリポジトリ上の実装状況をもとに整理しています。実機確認が必要なものは「要確認」としています。

| No. | Issue | 状態 | 対応バージョン | 概要 |
| --- | --- | --- | --- | --- |
| 121 | [自動撮影をローカル視点基準とStream方式中心に修正する](121-stream-camera-local-view-and-error-ux.md) | 要確認 | `v0.1.8` | ローカル視点、テスト撮影、失敗表示、UI文言、Stream方式を修正する。 |
| 129 | [Spout受信ヘルパーを追加する](129-add-spout-capture-helper.md) | 要確認 | `v0.1.8` | Windows同梱のSpout受信ヘルパーでsender列挙と1フレームPNG保存を行う。 |
| 130 | [Spout sender設定と診断UIを追加する](130-add-spout-sender-settings-and-diagnostics.md) | 要確認 | `v0.1.8` | 自動撮影タブからSpout helper確認、sender一覧取得、sender選択/自動選択を設定できるようにする。 |
| 131 | [自動撮影Stream方式をSpoutヘルパーへ統合する](131-integrate-spout-helper-into-auto-capture.md) | 要確認 | `v0.1.8` | Stream方式でFFmpeg画面キャプチャではなくSpoutヘルパーから画像を取得する。 |
| 132 | [Spout取得画像とメタデータを検証する](132-validate-spout-capture-output-and-metadata.md) | 要確認 | `v0.1.8` | Spout取得画像の有効性とsender情報を検証し、sidecar/Discordへ紐づける。 |
| 134 | [CI/ReleaseでSpoutヘルパーをビルド/同梱する](134-package-spout-helper-in-ci-release.md) | 要確認 | `v0.1.8` | Windows CI/ReleaseでSpoutヘルパーと必要DLL/ライセンスをビルド・同梱・検証する。 |
| 141 | [自動撮影構図UIへ `player_local` を統合する](141-integrate-player-local-coordinate-ui.md) | 要確認 | `v0.1.8` | 構図設定、現在Pose保存/追加、リセット、移動、テスト撮影を座標系に対応させる。 |
| 166 | [v0.1.8-rc13作成可能状態まで未完了項目を解消する](166-prepare-v018-rc13-readiness.md) | 要確認 | `v0.1.8` | RC13作成前に未完了issueを実装・検証し、残課題を再チケット化する。 |
| 171 | [Spout helperを本体exeへ埋め込み単一exe配布にする](171-embed-spout-helper-single-exe-release.md) | 要確認 | `v0.1.8` | C案採用に基づき、通常利用者向けReleaseを単一exe主導線へ戻し、分離版zipも検証用に残す。 |
| 172 | [player_local基準Poseの挙動を見直す](172-clarify-or-redesign-player-local-basis-behavior.md) | 未着手 | 未定 | `player_local` が自動追従ではなく手動基準Pose方式であるため、初期構図とUI説明の期待ズレを解消する。 |
| 173 | [専用アバターギミックOSCでHips/avatar基準Poseを自動取得する](173-implement-avatar-osc-basis-bridge.md) | 要確認 | `v0.1.8` | YL-ATG方式を参考に、専用アバターギミックからOSCでHips/avatar基準Poseを受け取り `player_local` basisへ使う機能を実装する。 |
| 175 | [CIでアバターギミック元ファイルzipを配布する](175-package-avatar-gimmick-source-zip.md) | 要確認 | `v0.1.8-rc17` | `Assets/PoppoWorks/AvatarBeacon/...` に配置したPrefab等をCIで元ファイルzip化し、`.unitypackage` は手作業で作成してGitHub Releaseへ添付する。 |
| 176 | [YL-ATGを参考にAvatarBeaconアバターギミックを作成する](176-create-clipforvrchat-avatar-gimmick-from-yl-atg-reference.md) | 要確認 | `v0.1.8-rc24` | ユーザー配置済みのATG_ForAvatar packageを参考に汎用アバターギミックAvatarBeaconを作成し、`coord/*` と `forward/*` の汎用OSC parameter、YL-ATG由来部分のMITライセンス表記、stale診断、10秒ごとのOSC summaryログ、`avatar_osc` 初期値化、受信器維持、受信状態の自動更新とyaw表示を整備する。 |
| 177 | [自動処理Webhookの通常投稿フォールバックを明確にする](177-clarify-discord-webhook-fallback-for-auto-processing.md) | 要確認 | `v0.1.8-rc18` | 自動処理専用Webhookが空欄の場合に通常投稿用Webhookへフォールバックする表示と保存前確認条件を整理する。 |
| 178 | [avatar_osc受信状態でmanual基準Pose未設定エラーが出る](178-fix-avatar-osc-status-misleading-manual-basis-error.md) | 要確認 | `v0.1.8-rc18` | Avatar OSC未受信時にmanual基準Pose未設定エラーが出る表示を修正し、AvatarBeaconのbasis parameter確認先を分かりやすくする。 |
| 179 | [rc18でGUIが表示されずavatar_oscエラーが大量出力される](179-fix-rc18-gui-not-showing-avatar-osc-log-loop.md) | 要確認 | `v0.1.8-rc21` | rc18起動直後にAvatar OSC受信処理が大量ログを出し、GUI表示を阻害する問題を修正し、GUI起動診断ログと起動進捗表示で切り分けやすくする。rc20で見つかったfrontend template由来の `avatar is not defined` も修正する。 |
| 180 | [rc25でSpout PNG保存が失敗し、AvatarBeacon受信状態の説明が冗長](180-fix-rc25-spout-png-encoder-and-avatarbeacon-status-note.md) | 要確認 | `v0.1.8-rc26` | Spout helperのWIC PNG書き出しをRGBA非対応環境でも動く形式へ直し、AvatarBeacon受信状態の常時説明文を削除する。 |
| 181 | [自動撮影後にUser Camera状態をできるだけ元へ戻す](181-restore-user-camera-state-after-auto-capture.md) | 要確認 | `v0.1.8-rc26` | 撮影前のUser Camera状態を可能な範囲で取得し、撮影後にMode/Streaming/Pose/Zoom/Exposure/mask類を復元する。取得できない項目は設定画面の復元用デフォルト値で戻せるようにする。 |
| 182 | [Stream Camera起動直後にSpout senderが出る前に失敗する](182-wait-for-spout-sender-after-stream-camera-osc.md) | 要確認 | `v0.1.8-rc27` | Stream方式でOSC送信後、Spout sender生成をtimeout内で待ち、必要に応じてStream起動OSCを再送する。 |
| 183 | [Spout取得直後の空フレームが透明PNGとして保存される](183-wait-for-valid-spout-frame-and-avoid-failed-output.md) | 要確認 | `v0.1.8-rc29` | Stream方式でSpout起動直後の空フレームを保存せず、有効フレーム待機と失敗出力の隔離を行う。 |
| 184 | [情報画面の他所取得ファイル注意文を削除する](184-remove-info-screen-third-party-file-warning.md) | 要確認 | `v0.1.8-rc29` | 情報画面の公式配布場所案内から、他所取得ファイルに関する責任否認文を削除する。 |
| 185 | [診断zip暗号化テストの一時鍵時刻を安定化する](185-stabilize-diagnostic-encryption-test-key-time.md) | 要確認 | `v0.1.8-rc29` | 診断zip暗号化テストのOpenPGP一時鍵作成時刻を安定化し、ローカルgo testの時刻揺らぎ失敗を防ぐ。 |
| 187 | [カメラ未起動/起動直後のStream Camera Spout取得を安定化する](187-stream-camera-start-and-blank-spout-frame-diagnostics.md) | 要確認 | 未定 | Streamingの互換OSC送信とblank-frame統計ログで、Stream Camera/Spout取得失敗を切り分ける。 |
| 188 | [別パス/別バージョンを含めてClipForVRChatを単一起動にする](188-global-single-instance-across-install-paths.md) | 要確認 | `v0.1.8` | OSC port競合を避けるため、配置パスやバージョンが違っても単一起動にし、既存を閉じる/アクティブ化する選択肢を出す。 |
| 189 | [VRChatから受信したOSCを他アプリ向けに別ポートへ転送できるようにする](189-forward-vrchat-osc-to-secondary-ports.md) | 要確認 | `v0.1.8` | ClipForVRChatが代表して受信したVRChat OSC packetを設定した別ポートへ転送し、他OSC受信アプリとのポート競合を避ける。 |
| 190 | [カメラ撮影機能終了時に一時OSC状態を確実に解除する](190-make-camera-osc-reset-streaming-compat.md) | 要対応 | 未定 | カメラOSCリセットと撮影終了時の `/usercamera/Streaming=false` などを通常のStream制御と同じbool+numeric互換送信に揃え、成否にかかわらず一時OSC状態を解除する。 |
| 191 | [UI上のすべてのボタンにマウスオーバー説明を追加する](191-add-hover-descriptions-to-all-ui-buttons.md) | 要対応 | 未定 | すべてのUIボタンに、機能・対象・注意点が分かるマウスオーバー説明を追加する。 |
| 192 | [OSCタブに送受信OSCのリアルタイムログ表示を追加する](192-add-realtime-osc-log-panel-to-osc-tab.md) | 要対応 | 未定 | OSCタブ最下部に送受信/forwardの一時リアルタイムログを表示し、正規表現フィルタとコピー機能を追加する。 |

## 状態の意味

| 状態 | 意味 |
| --- | --- |
| 完了 | 実装、テスト、またはドキュメント作成が完了している。 |
| 要対応 | 問題を整理済みで、実装またはドキュメント更新が必要。 |
| 要確認 | 修正は入っているが、対象環境や実機での再現確認が必要。 |
| 一部将来対応あり | 主要な受け入れ条件は満たしているが、明示的に将来対応へ回した項目が残っている。 |
