# AvatarBeacon READMEを汎用OSCギミックとして修正する

## 指示

> readmeを修正
>
> 下記を削除 : なぜかというとAvatarBeaconはClipForVRChat専用ギミックではなく汎用的にアバターの座標値と向きをOSCで取得するギミックでどんなツールからでもそれを取得できることを想定しているから。
>
> ClipForVRChat は、この値を受け取って player_local 構図の基準にできます。名前とOSC parameterは汎用にしているため、同じ出力を読める別ツールでも利用できます。
> 主な用途は、外部ツールが「いまの自分のアバター基準で前後左右どちらへカメラを動かすか」を判断するための基準値を得ることです。
> ClipForVRChatでの確認 の項目全て
> ClipForVRChatのbasis復元に使わない保存用・デバッグ用menu/parameterを削除
>
> 下記を削除
> 前提の項目 指示になかった項目のはず。
>
> 下記を削除
> 精度確認や互換検証が必要な場合だけ AvatarBeacon_12.prefab を使います。
>
> ライセンス表記がおかしい
> 普通に書いてほしい
> AvatarBeaconはMIT Copyright (c) 2026 hatolife
> YL-ATGはMIT Copyright (c) 2024 YozoraKurage
> みたいな書き方になるのが普通なんじゃない？
>
> 全体で下記を変更
> - `。`で改行して

## 文脈

AvatarBeacon READMEは前回改稿でClipForVRChatの利用例を含めたが、AvatarBeacon自体はClipForVRChat専用ではなく、アバターの座標値と向きをOSCで取得する汎用ギミックとして説明すべき。

## 解釈

AvatarBeacon専用リポジトリのREADMEとUnity配布内READMEから、ClipForVRChat固有の用途説明、前提セクション、不要な `AvatarBeacon_12` 利用誘導を削除する。
ライセンスはAvatarBeacon本体とYL-ATG由来部分を分けて、MIT Licenseと著作権者を明記する。
日本語文は `。` で行を分ける。

## 問題

- READMEがClipForVRChat固有の目的に寄って見える。
- ライセンス説明がYL-ATG由来部分だけに見え、AvatarBeacon本体の権利表記が分かりにくい。
- 文の改行単位が読みづらい箇所がある。

## 期待する挙動

- AvatarBeaconが汎用的にアバターの座標値と向きをOSCで取得するギミックだと分かる。
- ClipForVRChat固有の説明がREADMEから外れる。
- ライセンス欄で `AvatarBeacon: MIT License, Copyright (c) 2026 hatolife` と `YL-ATG: MIT License, Copyright (c) 2024 YozoraKurage` が分かる。
- 日本語文が `。` ごとに改行されている。

## 受け入れ条件

- [x] AvatarBeacon READMEから指定されたClipForVRChat固有文と確認セクションを削除する。
- [x] READMEから前提セクションを削除する。
- [x] `AvatarBeacon_12.prefab` の「精度確認や互換検証が必要な場合だけ」文を削除する。
- [x] ライセンス表記をAvatarBeacon本体とYL-ATG由来部分に分ける。
- [x] 日本語文を `。` ごとに改行する。

## 完了メモ

- AvatarBeacon専用リポジトリのルートREADMEとUnity配布内READMEを、ClipForVRChat固有説明なしの汎用OSCギミック説明へ修正した。
- `前提` セクション、指定されたClipForVRChat確認セクション、指定された `AvatarBeacon_12.prefab` 利用誘導文を削除した。
- AvatarBeacon本体のMIT License本文を `LICENSE` と `Assets/PoppoWorks/AvatarBeacon/LICENSES/AvatarBeacon-MIT.txt` に追加した。
- READMEのライセンス欄を `AvatarBeacon: MIT License, Copyright (c) 2026 hatolife` と `YL-ATG: MIT License, Copyright (c) 2024 YozoraKurage` が分かる形へ変更した。
- AvatarBeacon専用リポジトリの最新 `main` は `bddd206`。GitHub Actions CI成功とsource artifact upload成功を確認した。
- ClipForVRChat側のsubmodule pointerを `bddd206` へ更新し、CI/ReleaseのAvatarBeacon source zip検証に `LICENSE` と `AvatarBeacon-MIT.txt` を追加した。
