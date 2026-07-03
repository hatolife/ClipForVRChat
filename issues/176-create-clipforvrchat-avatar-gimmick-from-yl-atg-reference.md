# YL-ATGを参考にAvatarBeaconアバターギミックを作成する

## 問題

`avatar_osc` による `player_local` basis受信はアプリ側に実装済みだが、VRChatアバターへ導入するアバターギミック本体がまだリポジトリに存在しない。
ユーザーが参考資料として `ATG_ForAvatar_V0.0.3.unitypackage` と展開済みディレクトリを配置したため、これを調査材料にしつつ、汎用アバターギミック `AvatarBeacon` のPrefab/Assetを作成する必要がある。
2026-07-02時点の `v0.1.8-rc16` 実機確認では、AvatarBeacon導入時にVRChatからOSC Avatar Parametersが送信されていないように見えるため、PrefabがOSC output条件を満たしているか追加確認が必要。

YL-ATGはMIT Licenseだが、Prefab、FBX、Material、構成、パラメータ設計など元プロジェクトに権利がある部分をコピーまたは改変する場合、著作権表示とライセンス表記を正しく保持する必要がある。

## 期待する挙動

`AvatarBeacon` を、アプリ本体と分かりやすく分離したUnity asset sourceとしてリポジトリに配置する。
CIは元ファイルzipを配布し、リリース担当者が手作業で作成した `.unitypackage` をUnityへimportすると、Assetは次のようなパスへ入る。

```text
Assets/PoppoWorks/AvatarBeacon/...
```

Prefabをアバターへ導入すると、VRChat OSC Avatar Parameters経由で外部ツールが受信できるbasis用のposition/forward情報を送信できる。
既定のparameterは `coord/*` と `forward/*` とし、ギミック名・配布名だけでなくOSC parameterにもClipForVRChat固有名を含めない。
必要であれば検証用にYL-ATG互換の `ATG` 系入力も扱えるようにする。

YL-ATG由来のAssetまたは実装をコピー・改変する場合は、YozoraKurage/YL-ATGのMIT License表示をソース、package、Release同梱ライセンスに反映する。
コピーせず仕様参考として再実装する場合でも、調査・参考元としてREADMEに明記するかどうかを実装時に判断し、コピー範囲を記録する。

## 参考資料

### ユーザー配置済みファイル

- `/mnt/c/Users/user/Downloads/ATG_ForAvatar_V0.0.3.unitypackage`
- `/mnt/c/Users/user/Downloads/ATG_ForAvatar_V0.0.3`

### `ATG_ForAvatar_V0.0.3.unitypackage` の確認結果

`unitypackage` はgzip圧縮tar形式で、package内の主なAsset pathは次の通り。

```text
Assets/YozoLab/YL-ATG_ForAvatar/FBXs/arrow.prefab
Assets/YozoLab/YL-ATG_ForAvatar/FBXs/Center.mat
Assets/YozoLab/YL-ATG_ForAvatar/FBXs/Z.mat
Assets/YozoLab/YL-ATG_ForAvatar/FBXs/Y.mat
Assets/YozoLab/YL-ATG_ForAvatar/FBXs/X.mat
Assets/YozoLab/YL-ATG_ForAvatar/FBXs/arrow.fbx
Assets/YozoLab/YL-ATG_ForAvatar/YL-ATG_ForAvatar.prefab
```

展開済みディレクトリでは、少なくとも次のPrefabが確認できる。

```text
Assets/YozoLab/YL-ATG_ForAvatar/YL-ATG_ForAvatar.prefab
```

`YL-ATG_ForAvatar.prefab` はUnity YAMLで、`ParentConstraint`、Contact系と思われる `MonoBehaviour`、`ATG/p/x`、`ATG/p/x+`、`ATG/p/z`、`ATG/p/z+` などのContact tag/parameter名が含まれている。

### YL-ATG repository

- Repository: https://github.com/YozoraKurage/YL-ATG
- License: MIT License, Copyright (c) 2024 YozoraKurage
- README上の導入要件:
  - `ATG_ForAvatar.unitypacage` はアバター用package。
  - Modular Avatarを利用する。
  - Prefabをアバター直下などへ配置し、`ATG/point` のMA Bone Proxy targetを指定する。
  - 既定targetはHead。
  - 開発環境はUnity 2022.3.6f1。

## 実装方針メモ

### リポジトリ配置

Issue #175 のCI元ファイルzip配布方針と合わせ、初期候補は次の構造にする。

```text
avatar-gimmicks/
  AvatarBeacon/
    Assets/
      PoppoWorks/
        AvatarBeacon/
          Prefabs/
          Materials/
          README.md
          LICENSES/
            YL-ATG-MIT.txt
```

`Assets/YozoLab/...` をClipForVRChat配布packageの公開import pathにはしない。
YL-ATGのAssetをそのまま検証用に保持する必要がある場合は、`third_party/`、`references/`、または明示的に配布対象外の場所へ分ける。

### OSC parameter

AvatarBeaconでは、ClipForVRChat専用名を避け、次の汎用parameterを既定にする。

```text
/avatar/parameters/coord/x
/avatar/parameters/coord/xSign
/avatar/parameters/coord/y
/avatar/parameters/coord/ySign
/avatar/parameters/coord/z
/avatar/parameters/coord/zSign
/avatar/parameters/forward/x
/avatar/parameters/forward/xSign
/avatar/parameters/forward/y
/avatar/parameters/forward/ySign
/avatar/parameters/forward/z
/avatar/parameters/forward/zSign
```

アプリ側は `coord/*` / `forward/*` を既定受信経路にし、`ATG/*` 互換は検証・切り分け用として残す。
過去RCで使った `CFVRC/basis/*` はClipForVRChat由来で汎用性が低いため、AvatarBeaconの既定parameterから外す。

### ライセンス方針

実装時に次のいずれに該当するかを明確に記録する。

- YL-ATGのPrefab/FBX/Material/Unity YAML/構成をコピーまたは改変した。
- YL-ATGの仕組みだけを参考にし、ClipForVRChat用にAssetを新規作成した。
- 一部Assetだけをコピーし、残りは新規作成した。

コピーまたは改変した部分がある場合は、最低限次を満たす。

- YL-ATGのMIT License全文を `LICENSES/YL-ATG-MIT.txt` などへ保存する。
- AvatarBeacon側READMEに、YozoraKurage/YL-ATG、URL、MIT License、利用または改変範囲を書く。
- 手動作成する `.unitypackage` 内にもライセンスまたはNOTICEを含める。
- CI生成元ファイルzip、手動作成 `.unitypackage`、Release配布物のライセンス/NOTICEにYL-ATG由来部分を含める。
- PrefabやAssetのメタ情報、README、NOTICEなど、実務上確認しやすい場所に由来を残す。

コピーせず新規作成した場合でも、実装PRまたはチケット追記で「何を参考にし、何をコピーしていないか」を明記する。

## 受け入れ条件

- [x] リポジトリ内に、AvatarBeaconのUnity asset sourceを追加する。
- [x] 配布packageのimport先が `Assets/PoppoWorks/AvatarBeacon/...` になる。
- [x] アバターへ導入できるPrefabを用意し、追跡対象TransformまたはBoneを設定できる。
- [x] 既定の追跡対象をHips相当にし、player rootそのものではないがHeadよりプレイヤー位置に近いbasisであることをREADMEに明記する。
- [x] VRChat Expressions Menuから手動でOSC疎通確認できるデバッグ用parameter/menuを追加する。
- [x] Prefabが `coord/*` と `forward/*` のposition/forward情報をOSC Avatar Parametersへ出せる。
- [x] AvatarBeacon Prefab内GameObjectの役割、必要性、削除判断を仕様書に記録する。
- [x] Expression Parameter枠、Contact数、Modular Avatar/VRCSDK依存、Performance Rankへの影響をREADMEに書く。
- [ ] Unity上でimportし、`Assets/PoppoWorks/AvatarBeacon/...` に展開されることを確認する。
- [ ] `v0.1.8-rc16` でOSCが送信されない原因を特定し、Prefab/導入手順/VRChat OSC設定のいずれが原因か切り分ける。
- [x] `avatar_osc` statusにraw受信件数と最後に受けたAvatar Parameter addressを表示し、OSC未送信とparameter不一致を切り分けられるようにする。
- [ ] ClipForVRChatが実機VRChatから新鮮なposition/yawを受信し、`player_local` 構図の追従撮影に使えることを確認する。
- [x] YL-ATGからコピー・改変した部分の有無と範囲を記録する。
- [x] YL-ATG由来部分がある場合、MIT License全文、著作権表示、NOTICE/README表記をソース、CI生成元ファイルzip、手動作成 `.unitypackage`、Release配布物へ含める。
- [x] #175 のCI生成対象として、このアバターギミックの元ファイルzipを出力できる。
- [ ] 手動作成した `.unitypackage` を #174 のRelease導線から利用者へ案内できる。

## 非対象

- VRChat client改造、メモリ読み、非公式APIによる位置取得。
- YL-ATG upstreamへ変更を取り込むこと。
- Quest単体でClipForVRChatなしに完結する仕組み。
- `Assets/YozoLab/...` をAvatarBeacon配布packageの主import pathとして使うこと。

## メモ

- AvatarBeacon元ファイルは作成済み。Unity/VRChat実機確認は未実施。
- 実装時は、まずユーザー配置済みpackageを参照して挙動と構成を再確認する。
- YL-ATG由来Assetを直接取り込む場合、MIT License上の再配布条件を満たすだけでなく、ClipForVRChat側の配布物から利用者が由来を確認できる状態にする。
- 2026-07-02: `CFVRC/basis/*` は汎用ギミック名として不適切なため、AvatarBeacon既定parameterを `coord/*` と `forward/*` へ変更する。
- 2026-07-02: `docs/avatarbeacon-spec.md` にPrefab構造、GameObjectごとの役割、削除候補、Unity実機確認なしに削るべきでない要素を記録した。
- 2026-07-02: ユーザー実機確認で `v0.1.8-rc16` のAvatarBeaconからOSCが送信されていないように見えるとの報告あり。`localOnly` parameter、Expression Parameters登録、VRChat OSC config生成、Contact receiverの動作条件を優先して調査する。
- 2026-07-02: VRChat公式OSC仕様ではPublished AvatarのOSC config JSONにある `output.address` が値変化時に送信される。実機切り分けでは `Reset OSC Config`、avatar ID別JSONの `coord/*` / `forward/*` 出力、Avatar Dynamics Contact / Avatar Interactions有効化、ClipForVRChatのraw受信件数を確認する。
- 2026-07-02: rc16時点の `SaveObject` はPrefab上、`avatar_beacon/save` を操作するMA Menu Itemであり、`coord/*` / `forward/*` のContact経路やClipForVRChatのbasis復元から参照されていない。用途不明な名前のまま残さず、OSC疎通確認用のデバッグメニューとして作り直す。
- 2026-07-02: `player_local` basis用途ではHeadよりHipsの方がプレイヤー位置に近いため、`point` の既定MA Bone ProxyをHeadからHipsへ変更する。
- 2026-07-02: `point` はMA Bone Proxy付きの追跡アンカーであり、複数のConstraintが参照するため残す。`arrow` mesh/materialは可視化用途のみでOSC送信に直接関与しないため削除する。
- 2026-07-02: YL-ATGの座標取得方式を `docs/avatarbeacon-spec.md` に追記した。VRChat clientやワールドAPIから座標を直接読むのではなく、Constraintで配置したContact Sender/ReceiverのProximity値を使い、magnitude/sign parameterへ分解してOSC Avatar Parametersとして外部へ出す方式。
- 2026-07-04: ユーザー実機確認で `avatar_osc 受信状態: stale / prefix: coord`、`raw: 55 / last: Ahoge_Angle`、position/yawが0、エラー `stale avatar OSC basis` の報告あり。VRChat OSC自体は他Avatar Parameterを受信できているが、AvatarBeaconの `coord/*` / `forward/*` が鮮度切れになっている状態と判断する。次の修正では、内部エラー名ではなく、AvatarBeacon座標parameterが止まっていること、OSC config reset、アバター再読み込み、Avatar Dynamics Contact / Avatar Interactions確認へ誘導する表示にする。
- 2026-07-04: AvatarBeacon実機切り分けのため、basis対象外も含めOSC受信内容をすべて診断ログへ出す。parseできたpacketはaddress/type/value/source/bytes、parse不能packetはbytesとhex previewを出す。
- 2026-07-04: `v0.1.8-rc22` 実機確認で `avatar_osc 受信状態: stale / prefix: coord`、`raw: 15 / last: forward/y`、basis age約21秒の報告あり。OSC経路と一部AvatarBeacon basis parameter受信は確認できているため、ログから `coord/*` / `forward/*` の全12parameterが揃っているか、どのparameterだけ更新停止しているかを確認する。
- 2026-07-04: rc22ログでは `coord/x,y,z` と `forward/x,y,z` は受信済み。`coord/zSign` と `forward/xSign` はログ上未確認だが、符号parameterは省略時に正方向として扱うためstaleの直接原因ではない。直接原因は、設定画面を開いた直後に同じ `127.0.0.1:9001` へOSC受信器を再bindしようとして `Only one usage of each socket address` で失敗し、その後既存受信器も `context canceled` で停止して受信器が残らなかったこと。加えて、basis鮮度時刻をbasis構成parameterの最古時刻で計算しており、VRChat OSCの値変化時送信と相性が悪く、変化しない軸だけでstaleになり得るため、basis対象parameterの最新受信時刻を見るよう修正する。
- 2026-07-04: AvatarBeaconをv0.1.8で使う前提にするため、新規設定、空文字、basisSource欠落設定のプレイヤー基準取得元初期値を `avatar_osc` に変更する。未知値のfallbackは従来通り `manual` とし、壊れたconfigを暗黙にOSC依存へ切り替えない。
