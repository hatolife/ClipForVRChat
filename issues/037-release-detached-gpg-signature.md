# 037 Release exe の外部署名ファイル

## 状態

対応中

## 問題

Releaseで配布される実行ファイルが公式配布物から改竄されていないか、ユーザーが確認できる外部署名ファイルがない。

## 期待する挙動

GitHub Releases に、zipとは別に `ClipForVRChat.exe` を検証するための `ClipForVRChat.exe.asc` を添付する。

## 受け入れ条件

- CIで `ClipForVRChat.exe` に対するPGP detached signatureを生成する。
- Release asset に `ClipForVRChat.exe.asc` をzipとは別に添付する。
- 生成した署名をCI内で検証する。
- READMEと情報画面の説明がexe署名の確認手順になっている。
