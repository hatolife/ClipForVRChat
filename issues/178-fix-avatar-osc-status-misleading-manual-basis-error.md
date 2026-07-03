# avatar_osc受信状態でmanual基準Pose未設定エラーが出る

## 問題

自動撮影タブの `avatar_osc 受信状態` で `raw: 0` のとき、実際にはVRChatからAvatar OSC parameterを受信していない状態なのに、`プレイヤー基準Poseが未設定です` というmanual basis向けのエラーが表示される。
これにより、AvatarBeaconのbasis parameterが届いていないのか、手動基準Poseを保存すべきなのかが分かりにくい。

## 期待する挙動

`GetAvatarOSCBasisStatus` は現在のplayer_local basis source設定に関係なく、Avatar OSCの受信状態だけを返す。
`raw: 0` の場合は、Avatar OSC parameter未受信として確認すべき内容が分かるエラーを表示する。

## 受け入れ条件

- [x] `GetAvatarOSCBasisStatus` がmanual basis未設定エラーを返さない。
- [x] `raw: 0` の場合、VRChat OSC、専用アバターギミック、`/avatar/parameters` 送信の確認を促す。
- [x] AvatarBeaconのbasis確認先が `coord/*` / `forward/*` だと分かる。
- [x] 自動撮影タブの説明がHips/avatar基準に統一される。

## メモ

- 2026-07-02: `GetAvatarOSCBasisStatus` が `GetLatestPlayerLocalBasis` 経由でmanual basis未設定エラーを返していたため、Avatar OSC専用の `latestAvatarOSCBasisSnapshotLocked` を直接返すようにした。
- 2026-07-02: `raw: 0` 時のエラー文に `coord/*` / `forward/*` の確認を明記した。
- 2026-07-04: AvatarBeaconからDebug Ping menu/parameterを削除する方針に合わせ、未受信時の案内を `coord/*` / `forward/*` 確認へ一本化する。
- 2026-07-02: 回帰テスト `TestAppAvatarOSCBasisStatusIgnoresManualBasisMissing` を追加した。
