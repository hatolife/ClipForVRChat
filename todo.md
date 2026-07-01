# issue 173 作業チェックリスト

- [x] frontend の自動撮影タブに `manual` / `avatar_osc` の basis source UI を追加する
- [x] `avatar_osc` の受信状態、最終受信、position/yaw、エラー表示を追加する
- [x] backend で `/avatar/parameters/...` を受信し、`CFVRC/basis` / `ATG` 系のbasisへ復元する
- [x] `avatar_osc` basis を自動撮影と player_local 構図保存へ適用する
- [x] sidecar / 埋め込みmetadataへ basis source と basis pose を記録する
- [x] Wails wrapper と frontend API 呼び出しを合わせる
- [x] README / SPEC / docs に専用アバターギミック必須と head/avatar 基準の注意を追記する
- [x] Go test / Wails API surface / frontend build を通す
- [ ] Windows実機で専用アバターOSCから position / yaw / stale / player_local追従を確認する
