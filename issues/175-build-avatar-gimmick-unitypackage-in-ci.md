# CIでアバターギミック用unitypackageを生成する

## 問題

`avatar_osc` は専用アバターギミックが必要だが、このリポジトリには現時点でアバター側のPrefabやUnityプロジェクトが存在しない。
Releaseで通常配布zipへアバターギミック用 `unitypackage` を同梱するには、リポジトリ上の配置、Unity packageの生成方法、CIでの検証方法を決める必要がある。

## 期待する挙動

アバターギミックのデータをリポジトリ内でアプリ本体と分かりやすく分離し、CIで `ClipForVRChatAvatarOSCBasis-vX.Y.Z.unitypackage` のような成果物を生成できる。
`.unitypackage` をUnityへimportしたとき、アセットは次の構造で入る。

```text
Assets/PoppoWorks/<アバターギミック名>/<アバターギミックのデータ Prefabなど>
```

推奨配置案:

```text
avatar-gimmicks/
  ClipForVRChatAvatarOSCBasis/
    Assets/
      PoppoWorks/
        ClipForVRChatAvatarOSCBasis/
          Prefabs/
          Materials/
          Animations/
          README.md
    ProjectSettings/
    Packages/
    Assets/Editor/ExportUnityPackage.cs
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

### CIでunitypackageを作る方法

Unity公式のCLIでは、Editorをbatch modeで起動し、staticなEditorメソッドを `-executeMethod` で呼べる。
`-executeMethod` で呼ぶメソッドはstaticで、Editor scriptは `Editor` フォルダに置く必要がある。
失敗時は例外を投げる、または `EditorApplication.Exit(1)` で非0終了にする。

Unity公式APIの `AssetDatabase.ExportPackage` は、指定したAsset pathを `.unitypackage` へ書き出せる。
フォルダ配下を含めるには `ExportPackageOptions.Recurse` を使う。
`ExportPackageOptions.IncludeDependencies` は依存Assetも含めるため、VRCSDKやModular Avatarなど外部依存が意図せず入らないか注意が必要。

推奨する初期方針:

- export対象は `Assets/PoppoWorks/ClipForVRChatAvatarOSCBasis` に限定する。
- owned assetはすべてこのフォルダ配下へ置く。
- VRCSDK、Modular Avatarなどは同梱せず、利用者側の前提依存としてREADMEに書く。
- 初期CIでは `ExportPackageOptions.Recurse` を基本にし、`IncludeDependencies` を使う場合はzip/package中身検査で外部依存混入を検出する。

想定Editor script:

```csharp
using System;
using System.IO;
using UnityEditor;

namespace PoppoWorks.ClipForVRChatAvatarOSCBasis.Editor
{
    public static class ExportUnityPackage
    {
        public static void Export()
        {
            var output = GetArg("-packageOutput");
            if (string.IsNullOrWhiteSpace(output))
            {
                throw new Exception("Missing -packageOutput");
            }

            Directory.CreateDirectory(Path.GetDirectoryName(output));
            AssetDatabase.ExportPackage(
                "Assets/PoppoWorks/ClipForVRChatAvatarOSCBasis",
                output,
                ExportPackageOptions.Recurse
            );
        }

        private static string GetArg(string name)
        {
            var args = Environment.GetCommandLineArgs();
            for (var i = 0; i < args.Length - 1; i++)
            {
                if (args[i] == name)
                {
                    return args[i + 1];
                }
            }
            return "";
        }
    }
}
```

想定CIコマンド:

```powershell
Unity.exe `
  -quit `
  -batchmode `
  -nographics `
  -projectPath avatar-gimmicks/ClipForVRChatAvatarOSCBasis `
  -executeMethod PoppoWorks.ClipForVRChatAvatarOSCBasis.Editor.ExportUnityPackage.Export `
  -packageOutput dist/ClipForVRChatAvatarOSCBasis-vX.Y.Z.unitypackage
```

GitHub Actionsで実行する場合は、Unity Editorとライセンスが必要になる。
選択肢は主に次の2つ。

- self-hosted runnerにUnity Editor、VRCSDK、Modular Avatar等を入れて実行する。
- GameCIなどのUnity向けAction/Dockerを使い、Unity license secretを設定してbatch modeを実行する。

## 受け入れ条件

- [ ] リポジトリ内に、アプリ本体と分離したアバターギミック用UnityプロジェクトまたはUnity package sourceを配置する。
- [ ] Unity上のアセットパスが `Assets/PoppoWorks/<アバターギミック名>/...` になる。
- [ ] プレイヤー位置座標をOSC Avatar Parametersへ出すPrefab等のアバターギミックデータを配置する。
- [ ] `Assets/Editor/ExportUnityPackage.cs` など、CIから呼べるexport用Editor scriptを追加する。
- [ ] CIで `.unitypackage` を生成し、Release artifact / Release Assetへ含められる。
- [ ] CIで `.unitypackage` の存在、ファイル名、sha256を検証する。
- [ ] 可能であれば生成した `.unitypackage` を空の検証用Unityプロジェクトへimportし、`Assets/PoppoWorks/<アバターギミック名>/...` に展開されることを確認する。
- [ ] VRCSDK、Modular Avatar、その他外部依存を同梱しない方針なら、READMEに導入前提を明記する。
- [ ] 外部依存を同梱する方針なら、ライセンスと再配布可否を確認し、Release zip内ライセンスへ反映する。
- [ ] #174 の通常配布zipへ、生成した `.unitypackage` を同梱する。

## 参考資料

- Unity Manual: Command line arguments (`-batchmode`, `-executeMethod`, error handling)
  - https://docs.unity.cn/2021.1/Documentation/Manual/CommandLineArguments.html
- Unity Scripting API: `AssetDatabase.ExportPackage`
  - https://docs.unity3d.com/ScriptReference/AssetDatabase.ExportPackage.html
- Unity Scripting API: `ExportPackageOptions`
  - https://docs.unity3d.com/550/Documentation/ScriptReference/ExportPackageOptions.html
- Unity Manual: command-line license management
  - https://docs.unity3d.com/6000.4/Documentation/Manual/ManagingYourUnityLicense.html
- GameCI Unity Builder documentation
  - https://game.ci/docs/github/builder/

## メモ

- 現時点では調査とissue化のみ。実装は行わない。
- アバターギミック名の候補は `ClipForVRChatAvatarOSCBasis`。最終名は実装時に確定する。
- `.unitypackage` はAsset pathを保持するため、source側も最初から `Assets/PoppoWorks/<アバターギミック名>/...` に置くのが安全。
- `IncludeDependencies` を安易に使うと、VRCSDKやModular Avatar側のAsset/Packageまで含む可能性があるため、初期実装ではowned assetを1フォルダへ集約する。
