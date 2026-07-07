# AvatarBeaconを専用リポジトリへ配置する

## 指示

> Avatar Beaconを別リポジトリに配置しようと思います。
> git@github-ht:hatolife/AvatarBeacon.git
> 上記は空です
> ~/work/AvatarBeaconを新規作成してそこに既存コード配置、CI設定
> push
> ciチェックまでやってほしい versionはv0.0.1

## 文脈

現在のAvatarBeacon元ファイルは `avatar-gimmicks/AvatarBeacon` に配置され、ClipForVRChat側CI/Releaseでsource zip化している。
専用リポジトリ `hatolife/AvatarBeacon` は作成済みだが空で、初回配置とCI設定が必要。

## 解釈

既存のAvatarBeacon元ファイルを `~/work/AvatarBeacon` へコピーし、専用リポジトリとして初期化する。
初回バージョンは `v0.0.1` としてメタデータへ埋め込み、GitHub Actionsでsource zip生成と基本検証を行えるようにする。

## 問題

- AvatarBeaconがClipForVRChat本体リポジトリ内にあり、単独でCIや配布物確認をしにくい。
- 専用リポジトリが空のため、初期コミット、タグ、CI確認が未実施。

## 期待する挙動

- `~/work/AvatarBeacon` が独立したGitリポジトリとして作成される。
- 既存AvatarBeacon元ファイル、README、ライセンス/NOTICEが保持される。
- `v0.0.1` のバージョン情報が `Version.txt` とPrefab内バージョン表示へ反映される。
- GitHub Actionsでsource zipとsha256を生成し、最低限の構成検証が通る。
- `hatolife/AvatarBeacon` へpushされ、CI結果を確認できる。

## 受け入れ条件

- [ ] `~/work/AvatarBeacon` に既存AvatarBeaconコードを配置する。
- [ ] 専用リポジトリ用のCI設定を追加する。
- [ ] `v0.0.1` のバージョンメタデータ検証が通る。
- [ ] `hatolife/AvatarBeacon` へ初回pushする。
- [ ] GitHub ActionsのCI成功を確認する。
