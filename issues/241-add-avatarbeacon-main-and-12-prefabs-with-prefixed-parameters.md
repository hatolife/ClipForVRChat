# AvatarBeacon_mainとAvatarBeacon_12を追加しparameter pathをavatar_beacon配下へ変更する

## 指示

> 1パラ圧縮は無しで_6まで実装

> ClipForVRChat側で_6かどうかを判定する処理を入れたい xSignなどSignが来てないことで判定できると思う

> prefab名を_6を_mainに 同じディレクトリに中身からのテキストファイルを配置し、そのファイル名を`_通常はmainをご使用ください。mainは6パラメータ使用します。_12は高精度ですが12パラメータ使用します`のようにしてほしい ファイル名が説明文になっているイメージ

## 文脈

AvatarBeaconはClipForVRChat専用ではないため、現行12parameter方式の情報を削って用途を限定するのではなく、通常利用向けの `main` Prefab と、高精度な `_12` Prefab を分けて提供する方針になった。

既存ユーザーはいないため、旧 `coord/*` / `forward/*` へのfallback互換は不要。

## 解釈

現行Prefabは構造を維持し、名称と公開parameterを `AvatarBeacon_12` / `avatar_beacon/...` に変更する。

`AvatarBeacon_main` は1parameter packではなく、位置3軸とforward3軸をそれぞれ1floatで表す6parameter方式として追加する。ClipForVRChat側は `*Sign` が届かないことで `main` 方式と判定する。ただし、Contact Proximity由来値から符号込みfloatをアバター内だけで安定生成できるかはUnity/VRChat実機確認が必要なため、静的に可能な範囲を実装し、仕様に検証限界を書く。

## 問題

- 現行parameter名が `coord/*` / `forward/*` で、他ギミックや汎用parameterと衝突しやすい。
- 現行Prefab名だけでは12parameter方式であることが分からない。
- 通常利用向けの6parameter方式Prefabがまだない。
- ClipForVRChat側の既定受信prefixが旧 `coord` 前提である。

## 期待する挙動

- 通常利用向けPrefabが `AvatarBeacon_main` として提供される。
- 現行高精度方式のPrefabが `AvatarBeacon_12` として提供される。
- `AvatarBeacon_12` のparameter pathは `avatar_beacon/coord/*` と `avatar_beacon/forward/*` になる。
- `AvatarBeacon_main` のPrefabが追加され、6parameter方式の意図と制約が分かる。
- Prefabディレクトリに、通常は `main` を使うことがファイル名で分かる説明用テキストファイルがある。
- ClipForVRChatの既定受信prefixが `avatar_beacon` になり、旧fallbackを前提にしない。

## 受け入れ条件

- [x] `AvatarBeacon_main.prefab` が通常利用向けとして追加され、旧 `AvatarBeacon.prefab` ではなく導入対象として説明される。
- [x] `AvatarBeacon_12.prefab` が高精度12parameter方式として追加される。
- [x] 12方式のparameterが `avatar_beacon/coord/x` などへ変更される。
- [x] `AvatarBeacon_main.prefab` が追加され、6parameter名が仕様に記録される。
- [x] Prefabディレクトリに `_通常はmainをご使用ください。mainは6パラメータ使用します。_12は高精度ですが12パラメータ使用します.txt` がある。
- [x] ClipForVRChat側の既定prefix、UI、診断文言、テストが `avatar_beacon` 前提になる。
- [x] 旧 `coord/*` / `forward/*` fallbackを追加しない。
- [x] Unity/VRChat実機で確認が必要な範囲をREADMEまたは仕様に明記する。

## 確認

- ローカル: `cd src && GOCACHE=/tmp/clipforvrchat-go-cache go test ./...`
- ローカル: `node scripts/check-frontend-template-literals.mjs`
- ローカル: `node scripts/check-wails-api-surface.mjs`
- ローカル: `git diff --check`
- 未実施: Unity Editor / VRChat 実機での import、Contact出力、OSC実送受信確認
