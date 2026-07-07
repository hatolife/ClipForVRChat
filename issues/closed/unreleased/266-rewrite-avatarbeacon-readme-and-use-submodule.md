# AvatarBeacon READMEを作り直し、ClipForVRChat側をsubmodule化する

## 指示

> AvatarBeaconのreadmeがひどいので ちゃんと0から書きたい これは何 導入は どういう出力 どのような仕組みなの YL-ATGとの関係は ライセンス のように書いて

追加指示:

> ClipForVRChatのAvatarBeaconはgitのサブモジュールにして

## 文脈

AvatarBeaconは `~/work/AvatarBeacon` に専用リポジトリとして初期配置済み。
ClipForVRChat側には同じ元ファイルが `avatar-gimmicks/AvatarBeacon` に通常ファイルとして残っている。

## 解釈

AvatarBeacon専用リポジトリの利用者向けREADMEを、既存説明の継ぎ足しではなく、目的、導入、出力、仕組み、YL-ATGとの関係、ライセンスが分かる構成で書き直す。
そのうえでClipForVRChat側の `avatar-gimmicks/AvatarBeacon` を専用リポジトリへのGit submoduleに置き換える。

## 問題

- 現在のREADMEはsource zipや作業用packageの説明が先に出て、利用者が「これは何か」「どう使うか」を把握しにくい。
- ClipForVRChat本体リポジトリにAvatarBeaconの実体ファイルが重複して残っており、専用リポジトリと同期しにくい。

## 期待する挙動

- AvatarBeacon READMEだけで、概要、導入手順、OSC出力、内部の仕組み、YL-ATG由来範囲、ライセンス確認先が分かる。
- ClipForVRChat側の `avatar-gimmicks/AvatarBeacon` が `hatolife/AvatarBeacon` のsubmoduleになる。
- 既存のCI/Releaseで参照するパスは `avatar-gimmicks/AvatarBeacon` のまま維持される。

## 受け入れ条件

- [x] `~/work/AvatarBeacon/README.md` を0から読める構成へ書き直す。
- [x] AvatarBeacon専用リポジトリでREADME変更をコミットしてpushする。
- [x] ClipForVRChatの `avatar-gimmicks/AvatarBeacon` をGit submoduleへ置き換える。
- [x] ClipForVRChat側でsubmoduleを含む状態確認と関連検証を行う。

## 完了メモ

- AvatarBeacon専用リポジトリで、ルートREADMEとUnity配布内READMEを「これは何」「導入」「出力」「仕組み」「YL-ATGとの関係」「ライセンス」中心に全面改稿した。
- AvatarBeacon専用リポジトリの最新 `main` は `cbca2ee`。GitHub Actions CI成功とsource artifact upload成功を確認した。
- ClipForVRChat側の `avatar-gimmicks/AvatarBeacon` を `https://github.com/hatolife/AvatarBeacon.git` のsubmoduleに置き換えた。
- CI/Releaseのcheckoutでsubmoduleを取得し、AvatarBeacon source zipは `README.md` と `Assets` だけを含めるようにした。
- `node scripts/update-avatarbeacon-version.mjs v0.0.1 --check` と、AvatarBeacon source zip相当の展開・必須ファイル・sha256検証を通した。
