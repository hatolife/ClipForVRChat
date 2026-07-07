# AvatarBeacon READMEのYL-ATG関係と謝辞を修正する

## 指示

> readme.mdの修正
>
> YL-ATGとの関係の項目の修正
> > AvatarBeacon は、YozoraKurage/YL-ATG ATG_ForAvatar_V0.0.3 を元にした派生物です。座標をContact / Constraint / Expression Parameterで外部へ出す考え方とPrefab構成の一部を引き継いでいます。
> `。`で改行されてない。指示守って。
> 下記に変更。
> ```
> AvatarBeacon は、YozoraKurage/YL-ATG ATG_ForAvatar_V0.0.3 を元にした派生物です。
>
> 簡単に言うと [ClipForVRChat](https://github.com/hatolife/ClipForVRChat) で都合がよいようにしたものです。
> 中身はほぼYL-ATGです。簡単に変更点を記載します。
> - 精度下げてパラメーター数を半分にしたPrefabを用意。
> 	- `AvatarBeacon_main.prefab` : 6パラメータ版。
> 	- `AvatarBeacon_12.prefab` : 12パラメータ版。YL-ATGのと同じ精度。
> - 向き用に HeadForwardAnchor を追加
> - 可視化用arrow mesh/materialを削除
> - 公開parameterを ATG/* から avatar_beacon/* に変更
> - 配置先を `Assets/PoppoWorks/AvatarBeacon` に変更。
> ```
> ---
>
> ライセンスの部分末尾に
> [夜空くらげ](https://x.com/yozorakurage)さんへの謝辞書きたい
> [ClipForVRChat](https://github.com/hatolife/ClipForVRChat) でカメラをアバター基準のローカル座標系に配置する機能は、YL-ATGがないと実現できなかった。
>
> ---
> 下記は削除
> 配布物 の項目
> 詳細仕様 の項目

## 文脈

AvatarBeacon READMEは汎用OSCギミックとして整理済みだが、YL-ATGとの関係、謝辞、不要セクションについて追加修正が必要。

## 解釈

AvatarBeacon専用リポジトリのルートREADMEとUnity配布内READMEで、YL-ATGとの関係セクションを指定文へ差し替える。
ライセンス欄末尾に夜空くらげさんへの謝辞を追加し、配布物・詳細仕様セクションを削除する。

## 問題

- YL-ATGとの関係説明が指定意図より抽象的。
- 謝辞が不足している。
- `配布物` と `詳細仕様` がREADMEの主目的から外れている。

## 期待する挙動

- YL-ATGとの関係が指定文に近い形で記載される。
- ライセンス欄末尾に夜空くらげさんへの謝辞がある。
- `配布物` と `詳細仕様` セクションがない。

## 受け入れ条件

- [x] ルートREADMEのYL-ATGとの関係を指定内容へ変更する。
- [x] Unity配布内READMEのYL-ATGとの関係を指定内容へ変更する。
- [x] ライセンス欄末尾に夜空くらげさんへの謝辞を追加する。
- [x] `配布物` セクションを削除する。
- [x] `詳細仕様` セクションを削除する。

## 完了メモ

- AvatarBeacon専用リポジトリのルートREADMEとUnity配布内READMEで、YL-ATGとの関係を指定内容へ修正した。
- ライセンス欄末尾に夜空くらげさんへの謝辞とClipForVRChatのアバター基準ローカル座標系機能に関する謝辞を追加した。
- ルートREADMEから `配布物` と `詳細仕様` セクションを削除した。
- 句点複数行チェック、指定削除対象検索、AvatarBeacon source zip相当の展開・必須ファイル・sha256検証を通した。
- AvatarBeacon専用リポジトリの最新 `main` は `65af631`。
- ClipForVRChat側のsubmodule pointerを `65af631` へ更新した。
