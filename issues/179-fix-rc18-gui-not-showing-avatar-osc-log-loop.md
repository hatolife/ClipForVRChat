# rc18でGUIが表示されずavatar_oscエラーが大量出力される

## 問題

`v0.1.8-rc18` の単一exe版、separated版のどちらでも起動してもGUIが表示されない。ユーザー提供ログでは、起動直後から `avatar osc basis resolve error: status="partial" err="partial avatar OSC basis sample: missing coord/x"` などが大量に連続出力され、短時間でログが肥大化している。

## 期待する挙動

- `avatar_osc` の受信状態がpartial/missingでも、GUIは通常通り表示される。
- AvatarBeaconやVRChat OSCの状態に問題がある場合も、ログ出力は抑制され、UI上の受信状態で確認できる。
- 起動直後に大量ログ出力や高負荷でアプリが固まらない。

## 受け入れ条件

- rc18相当の設定/OSC受信状態でもGUI表示を阻害しない。
- `avatar osc basis resolve error` がループで大量出力されない。
- `go test ./...` とフロントエンドビルドが通る。
- 必要なら次RCで修正を確認できる。

## 調査メモ

- 提供ログは43,275行中43,263行が `avatar osc basis resolve error` だった。
- `basis_source="manual"` でもAvatar OSC受信処理が全Avatar Parameter受信ごとにbasis再構築を試みていた。
- `latestAvatarOSCBasisSnapshotLocked` がmanual設定のまま `ResolvePlayerLocalBasisPose` を呼び、manual基準Pose未設定エラーをAvatar OSC診断にも混ぜていた。
- 対策として、Avatar OSC診断snapshotでは復元時だけ `BasisSource=avatar_osc` を強制し、basisに関係するOSC address以外では再構築しない。非readyログも状態変化時だけに抑制する。
