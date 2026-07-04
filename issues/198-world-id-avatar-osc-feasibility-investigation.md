# world ID を Avatar OSC で送れるか調査する

## 問題

VRChat の Avatar / OSC / Expression Parameter / Built-in Parameter / World API を確認し、アバター側から現在の world ID を文字列として外部 OSC に送れるか不明である。
AvatarBeacon へ world ID 送信ギミックを追加できるか、また不可能なら代替案へ切り替える必要がある。

## 期待する挙動

- アバター単体、または AvatarBeacon の拡張で現在 world ID を外部へ伝えられるか判定できる。
- 公式手段で不可能なら、その理由と代替案を明確にできる。
- AvatarBeacon に追加できる範囲がある場合、どの Prefab / GameObject / 参照ファイルを拡張対象にするか分かる。

## 受け入れ条件

- [x] VRChat 公式ドキュメント上、avatar から外部 OSC へ world ID を string で送る経路の有無を確認する。
- [x] Built-in Parameters、Parameter Driver、OSC Avatar Parameters、World 側 API のどこで world ID が扱えるか確認する。
- [x] 公式手段がない場合、output log 監視、写真 metadata、world 側ギミック、AvatarBeacon への token 送信などの代替案を比較する。
- [x] AvatarBeacon に実装可能な範囲があるなら、影響するファイル / Prefab / GameObject 名を明記する。
- [x] 判定結果を issues と docs に反映するか、少なくとも実装可否判断の根拠をまとめる。

## 調査結果

- 2026-07-04時点で確認したVRChat公式OSC仕様では、アバターから現在world ID文字列を外部OSCへ送る標準経路は確認できなかった。
- Avatar Parameters / Parameter Driver / Contact Receiverで扱える値は主にbool/int/floatであり、`wrld_...` のような文字列をAvatarBeacon単体で生成・送信する前提にはできない。
- AvatarBeaconはHips/Head/WorldOriginAnchorの位置差分や向きの数値をOSC parameterとして外へ出す用途なら実装できるが、現在world IDそのものを知るsourceを持たない。
- 現在world/instance情報の取得は、既存の `SnapshotVRChatWorld()` によるVRChat output log監視を使うのが現実的。
- 撮影後の写真metadataにはworld IDが入る可能性があるが、開始時撮影前の移動中判定には使いにくい。
- world側ギミックでworld情報を扱える可能性はあるが、AvatarBeaconとは別配布物になり、全worldで使える汎用解ではない。

## 結論

AvatarBeaconからworld IDをOSC送信する実装は、現時点では公式手段がなく不可と判断する。

v0.1.8の開始時撮影・移動中判定では、AvatarBeaconのOSC basis ready/freshとVRChat output log由来のworld/instance安定判定を組み合わせる。
