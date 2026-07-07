# AvatarBeaconを専用リポジトリへ配置する

## 指示

> Avatar Beaconを別リポジトリに配置しようと思います。
> git@github-ht:hatolife/AvatarBeacon.git
> 上記は空です
> ~/work/AvatarBeaconを新規作成してそこに既存コード配置、CI設定
> push
> ciチェックまでやってほしい versionはv0.0.1

追加指示:

> readmeの# AvatarBeacon 元ファイル packageが意味不明なので# AvatarBeaconにして

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
- 専用リポジトリのルートREADME見出しが `# AvatarBeacon` になっている。

## 受け入れ条件

- [x] `~/work/AvatarBeacon` に既存AvatarBeaconコードを配置する。
- [x] 専用リポジトリ用のCI設定を追加する。
- [x] `v0.0.1` のバージョンメタデータ検証が通る。
- [x] `hatolife/AvatarBeacon` へ初回pushする。
- [x] ルートREADME見出しを `# AvatarBeacon` にする。
- [x] GitHub ActionsのCI成功を確認する。

## 完了メモ

- `~/work/AvatarBeacon` を新規Gitリポジトリとして作成し、`Assets/PoppoWorks/AvatarBeacon`、README、仕様書、バージョン更新スクリプト、source tree検証スクリプトを配置した。
- AvatarBeaconのメタデータを `v0.0.1` に更新し、`node scripts/update-avatarbeacon-version.mjs v0.0.1 --check`、`node scripts/check-source-tree.mjs .`、source zip生成、展開後検証、sha256検証を通した。
- `main` と `v0.0.1` タグを `hatolife/AvatarBeacon` へpushした。
- 初回CIでsha256検証の作業ディレクトリ誤りが見つかったため修正し、最新 `main` CI成功を確認した。
- `v0.0.1` タグは初回コミットを指しており、共有済みタグを書き換えない方針のため、タグCIには修正前CIの失敗履歴が残る。
