# CIでアバターギミック元ファイルzipを配布する

## 問題

`avatar_osc` は専用アバターギミックが必要だが、このリポジトリには現時点でアバター側のPrefabやUnity package sourceが存在しない。
当初はCIで `.unitypackage` を生成する案を検討したが、GitHub Actions上でUnity Editor、Unity license、VRCSDK、Modular Avatar依存を扱う必要があり、運用負荷と失敗要因が大きい。

そのため、CIではアバターギミックに必要な元ファイル群をzip化して配布し、`.unitypackage` 化はリリース担当者がUnity上で手作業で行い、GitHub Releaseへ手動添付する方針へ変更する。
アバターギミック本体の作成とYL-ATG由来部分のライセンス整理は #176 で扱う。

## 期待する挙動

アバターギミックのデータをリポジトリ内でアプリ本体と分かりやすく分離し、CIで `AvatarBeacon-vX.Y.Z-source.zip` のような元ファイルzipを生成できる。

zip内には、Unityへコピーまたはimportして手作業で `.unitypackage` を作るために必要なowned asset、README、ライセンス/NOTICEを含める。
CIではUnity Editorを起動せず、zipの中身、ファイル名、sha256、不要ファイル混入を検証する。

リリース担当者は、そのzipを元にUnity上で `.unitypackage` を作成し、GitHub Release Assetへ手動で配置する。
Release本文と通常配布導線では、CI生成zipと手動添付 `.unitypackage` の違いが利用者に分かるようにする。

想定配置案:

```text
avatar-gimmicks/
  AvatarBeacon/
    Assets/
      PoppoWorks/
        AvatarBeacon/
          Prefabs/
          Materials/
          Animations/
          README.md
          LICENSES/
            YL-ATG-MIT.txt
    README.md
```

`.unitypackage` をUnityへimportしたときの最終Asset pathは、手作業作成後も次の構造になる。

```text
Assets/PoppoWorks/<アバターギミック名>/<アバターギミックのデータ Prefabなど>
```

## 調査結果

### 現在のリポジトリ内配置

2026-07-02時点で、リポジトリ内に `.prefab`、`.unitypackage`、Unity `Assets/`、`ProjectSettings/`、`Packages/` は存在しない。
`avatar_osc` の受信・復元ロジックと検証手順はあるが、プレイヤー位置座標をOSCで送信するアバターギミックそのものはまだ配置されていない。

関連する既存実装・文書:

- `src/app.go`: `/avatar/parameters/...` のOSC受信と `CFVRC/basis` / `ATG` 系parameter復元。
- `src/internal/appcore/player_local.go`: 受信値からbasis poseを復元。
- `docs/v0.1.9-avatar-osc-basis-verification.md`: 実機確認手順。
- `issues/173-implement-avatar-osc-basis-bridge.md`: アバターギミック方式の設計メモ。
- `issues/176-create-clipforvrchat-avatar-gimmick-from-yl-atg-reference.md`: ギミック本体作成とYL-ATGライセンス整理。

### CIでunitypackageを作らない理由

Unity公式のCLIでは、Editorをbatch modeで起動し、`AssetDatabase.ExportPackage` を呼ぶことで `.unitypackage` を生成できる。
ただし、この方式はCIにUnity Editorと有効なUnity licenseを用意する必要がある。
さらに、VRCSDKやModular AvatarなどのUnity依存をCIへ持ち込むと、依存の再配布可否、バージョン固定、license secret管理、import失敗、外部依存混入検査が必要になる。

今回の目的は「利用者に必要なアバターギミックを配布すること」であり、`.unitypackage` の完全自動生成にこだわる必要はない。
CIではリポジトリ内のowned assetと説明書/ライセンスをzip化し、リリース担当者がUnity上で確認したうえで `.unitypackage` を手動作成する方が現実的と判断する。

### CIでzip化する対象

CIの対象はUnity project全体ではなく、配布対象のowned asset rootに限定する。

推奨方針:

- zip対象は `avatar-gimmicks/<アバターギミック名>/Assets/PoppoWorks/<アバターギミック名>` と、その周辺README/ライセンスに限定する。
- VRCSDK、Modular Avatar、Unity `Library/`、`Temp/`、`ProjectSettings/`、`Packages/` は原則同梱しない。
- `.meta` はUnity Asset参照維持に必要なため、owned assetに対応するものは含める。
- YL-ATG由来Assetがある場合、MIT License全文とNOTICE/READMEを必ず含める。
- zip内のパスは、Unity projectへ展開したときに `Assets/PoppoWorks/<アバターギミック名>/...` になる形にする。

## 受け入れ条件

- [ ] リポジトリ内に、アプリ本体と分離したアバターギミック用Unity asset sourceを配置する。
- [ ] Unity上のアセットパスが `Assets/PoppoWorks/<アバターギミック名>/...` になる。
- [ ] プレイヤー位置座標をOSC Avatar Parametersへ出すPrefab等のアバターギミックデータを配置する。
- [ ] CIで `.unitypackage` は生成しない。
- [ ] CIでアバターギミック元ファイルzipを生成し、Release artifact / Release Assetへ含める。
- [ ] CIで元ファイルzipの存在、ファイル名、sha256、zip内ファイル一覧を検証する。
- [ ] zip内に `Assets/PoppoWorks/<アバターギミック名>/...`、README、必要なライセンス/NOTICEが含まれる。
- [ ] zip内にUnity `Library/`、`Temp/`、不要な外部依存、VRCSDK、Modular Avatar本体が混入しない。
- [ ] `.unitypackage` はリリース担当者がUnity上で手作業作成し、GitHub Release Assetへ手動添付する運用をREADMEまたはリリース手順に書く。
- [ ] 手動作成した `.unitypackage` のファイル名、配置、sha256確認手順をリリース手順に書く。
- [ ] VRCSDK、Modular Avatar、その他外部依存を同梱しない方針なら、READMEに導入前提を明記する。
- [ ] 外部依存を同梱する方針に変える場合は、ライセンスと再配布可否を確認し、Release配布物のライセンス/NOTICEへ反映する。
- [ ] #174 のRelease導線で、CI生成元ファイルzipと手動添付 `.unitypackage` の関係が分かる。

## 参考資料

- Unity Manual: Command line arguments (`-batchmode`, `-executeMethod`, error handling)
  - https://docs.unity.cn/2021.1/Documentation/Manual/CommandLineArguments.html
- Unity Scripting API: `AssetDatabase.ExportPackage`
  - https://docs.unity3d.com/ScriptReference/AssetDatabase.ExportPackage.html
- Unity Manual: command-line license management
  - https://docs.unity3d.com/6000.4/Documentation/Manual/ManagingYourUnityLicense.html

## メモ

- 現時点では方針更新のみ。アバターギミック本体とCI実装は未実施。
- アバターギミック名は `AvatarBeacon` とする。ClipForVRChat専用ではなく、同じOSC parameterを受信できる他ツールでも使える汎用アバターギミックとして扱う。
- `.unitypackage` はAsset pathを保持するため、source側も最初から `Assets/PoppoWorks/<アバターギミック名>/...` に置くのが安全。
- CIでUnityを動かさないため、Unity import確認と `.unitypackage` 作成確認は手作業のリリースチェックに寄せる。
