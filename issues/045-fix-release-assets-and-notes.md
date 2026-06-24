# 045 Release本文と成果物の再発防止

## 状態

完了

## 問題

v0.1.6 Releaseで、`RELEASE_NOTES.md` の更新内容がGitHub Release本文に反映されず、不要な公開鍵ファイルがRelease添付とzip内に含まれた。

## 期待する挙動

GitHub Releaseには該当バージョンのリリースノート本文を記載し、成果物は配布zip、exe署名ファイル、必要な補助ファイルだけにする。公開鍵はGitHub Releaseでは配布せず、keys.openpgp.orgへの案内に統一する。

## 受け入れ条件

- Release workflowが `RELEASE_NOTES.md` の該当バージョン本文をGitHub Release本文に使う。
- Release添付から `*-release-signing-public-key.asc` を削除する。
- zip内から `ClipForVRChat-release-signing-public-key.asc` を削除する。
- zip内に `.asc` 署名ファイルを含めない。
- zip内にkeys.openpgp.orgを開くURLショートカットを同梱する。
- Release作成前に、想定外の公開鍵ファイルが成果物に含まれていないことを検査する。
- `AGENTS.md` にRelease成果物と本文確認の再発防止ルールを追加する。
