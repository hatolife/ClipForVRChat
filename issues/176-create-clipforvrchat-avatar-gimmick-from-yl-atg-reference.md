# YL-ATGを参考にClipForVRChat用アバターギミックを作成する

## 問題

`avatar_osc` による `player_local` basis受信はアプリ側に実装済みだが、VRChatアバターへ導入するClipForVRChat用アバターギミック本体がまだリポジトリに存在しない。
ユーザーが参考資料として `ATG_ForAvatar_V0.0.3.unitypackage` と展開済みディレクトリを配置したため、これを調査材料にしつつ、ClipForVRChat用のPrefab/Assetを作成する必要がある。

YL-ATGはMIT Licenseだが、Prefab、FBX、Material、構成、パラメータ設計など元プロジェクトに権利がある部分をコピーまたは改変する場合、著作権表示とライセンス表記を正しく保持する必要がある。

## 期待する挙動

ClipForVRChat用のアバターギミックを、アプリ本体と分かりやすく分離したUnity asset sourceとしてリポジトリに配置する。
CIは元ファイルzipを配布し、リリース担当者が手作業で作成した `.unitypackage` をUnityへimportすると、Assetは次のようなパスへ入る。

```text
Assets/PoppoWorks/ClipForVRChatAvatarOSCBasis/...
```

Prefabをアバターへ導入すると、VRChat OSC Avatar Parameters経由でClipForVRChatが受信できるbasis用のposition/forward情報を送信できる。
既定のparameter prefixはClipForVRChat専用の `CFVRC/basis` とし、必要であれば検証用にYL-ATG互換の `ATG` 系入力も扱えるようにする。

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
  ClipForVRChatAvatarOSCBasis/
    Assets/
      PoppoWorks/
        ClipForVRChatAvatarOSCBasis/
          Prefabs/
          Materials/
          README.md
          LICENSES/
            YL-ATG-MIT.txt
```

`Assets/YozoLab/...` をClipForVRChat配布packageの公開import pathにはしない。
YL-ATGのAssetをそのまま検証用に保持する必要がある場合は、`third_party/`、`references/`、または明示的に配布対象外の場所へ分ける。

### OSC parameter案

ClipForVRChat専用Prefabでは、次のparameterを既定にする。

```text
/avatar/parameters/CFVRC/basis/p/x
/avatar/parameters/CFVRC/basis/p/xSign
/avatar/parameters/CFVRC/basis/p/y
/avatar/parameters/CFVRC/basis/p/ySign
/avatar/parameters/CFVRC/basis/p/z
/avatar/parameters/CFVRC/basis/p/zSign
/avatar/parameters/CFVRC/basis/f/x
/avatar/parameters/CFVRC/basis/f/xSign
/avatar/parameters/CFVRC/basis/f/y
/avatar/parameters/CFVRC/basis/f/ySign
/avatar/parameters/CFVRC/basis/f/z
/avatar/parameters/CFVRC/basis/f/zSign
```

アプリ側は既に `CFVRC/basis/*` と `ATG/*` 互換の受信を扱うため、Prefab側はまず `CFVRC/basis/*` の送信を優先し、必要に応じて `ATG/*` 互換Prefabまたはmigration手順を追加する。

### ライセンス方針

実装時に次のいずれに該当するかを明確に記録する。

- YL-ATGのPrefab/FBX/Material/Unity YAML/構成をコピーまたは改変した。
- YL-ATGの仕組みだけを参考にし、ClipForVRChat用にAssetを新規作成した。
- 一部Assetだけをコピーし、残りは新規作成した。

コピーまたは改変した部分がある場合は、最低限次を満たす。

- YL-ATGのMIT License全文を `LICENSES/YL-ATG-MIT.txt` などへ保存する。
- ClipForVRChat側READMEに、YozoraKurage/YL-ATG、URL、MIT License、利用または改変範囲を書く。
- 手動作成する `.unitypackage` 内にもライセンスまたはNOTICEを含める。
- CI生成元ファイルzip、手動作成 `.unitypackage`、Release配布物のライセンス/NOTICEにYL-ATG由来部分を含める。
- PrefabやAssetのメタ情報、README、NOTICEなど、実務上確認しやすい場所に由来を残す。

コピーせず新規作成した場合でも、実装PRまたはチケット追記で「何を参考にし、何をコピーしていないか」を明記する。

## 受け入れ条件

- [ ] リポジトリ内に、ClipForVRChat用アバターギミックのUnity asset sourceを追加する。
- [ ] 配布packageのimport先が `Assets/PoppoWorks/ClipForVRChatAvatarOSCBasis/...` になる。
- [ ] アバターへ導入できるPrefabを用意し、追跡対象TransformまたはBoneを設定できる。
- [ ] 既定の追跡対象をHead相当にし、Head基準でありplayer root基準ではないことをREADMEに明記する。
- [ ] Prefabが `CFVRC/basis/*` のposition/forward情報をOSC Avatar Parametersへ出せる。
- [ ] Expression Parameter枠、Contact数、Modular Avatar/VRCSDK依存、Performance Rankへの影響をREADMEに書く。
- [ ] Unity上でimportし、`Assets/PoppoWorks/ClipForVRChatAvatarOSCBasis/...` に展開されることを確認する。
- [ ] ClipForVRChatが実機VRChatから新鮮なposition/yawを受信し、`player_local` 構図の追従撮影に使えることを確認する。
- [ ] YL-ATGからコピー・改変した部分の有無と範囲を記録する。
- [ ] YL-ATG由来部分がある場合、MIT License全文、著作権表示、NOTICE/README表記をソース、CI生成元ファイルzip、手動作成 `.unitypackage`、Release配布物へ含める。
- [ ] #175 のCI生成対象として、このアバターギミックの元ファイルzipを出力できる。
- [ ] 手動作成した `.unitypackage` を #174 のRelease導線から利用者へ案内できる。

## 非対象

- VRChat client改造、メモリ読み、非公式APIによる位置取得。
- YL-ATG upstreamへ変更を取り込むこと。
- Quest単体でClipForVRChatなしに完結する仕組み。
- `Assets/YozoLab/...` をClipForVRChat配布packageの主import pathとして使うこと。

## メモ

- 現時点ではissue化のみ。アバターギミック本体の作成は未実施。
- 実装時は、まずユーザー配置済みpackageを参照して挙動と構成を再確認する。
- YL-ATG由来Assetを直接取り込む場合、MIT License上の再配布条件を満たすだけでなく、ClipForVRChat側の配布物から利用者が由来を確認できる状態にする。
