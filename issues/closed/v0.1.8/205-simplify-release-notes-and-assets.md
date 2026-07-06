# Release note書式と公開Assetを整理する

## 問題

Release本文とGitHub Release Assetsに、通常利用者向け主導線と検証用ファイルが混在している。次回以降のRelease noteは、直前正式版からの変更点を端的に示し、公開Assetも必要な配布物だけに絞りたい。

## 期待する挙動

Release noteは、バージョン見出し、ダウンロード、更新内容、比較の順で統一される。GitHub Releaseへ添付されるファイルは、通常zip、単一exe署名asc、separated zip、AvatarBeacon source zipだけになる。通常zipには単一exe、署名asc、署名用公開鍵URL、README、LICENSEだけが入る。

## 受け入れ条件

- リリースノート記載書式の文書が更新され、次回以降の書式が明確になっている。
- Release workflowが通常zipへ `README.md`、`ClipForVRChat.exe`、単一exe署名asc、`Release-signing-public-key.url`、`LICENSE` を含める。
- 通常zipに `Spout2-LICENSE.txt`、`spout-capture.exe`、`SpoutLibrary.dll` を直接含めない。
- GitHub ReleaseへuploadするAssetは `ClipForVRChat-<バージョン>-windows-amd64.zip`、`ClipForVRChat-<バージョン>-windows-amd64.exe.asc`、`ClipForVRChat-<バージョン>-windows-amd64-separated.zip`、`AvatarBeacon-<バージョン>-source.zip` だけにする。
- 同じtagでRelease workflowを再実行しても、上記以外の既存Assetを削除してからuploadする。

## 対応内容

- `AGENTS.md` のリリース運用に、Release note本文の固定書式と更新内容の粒度を追記した。
- `RELEASE_NOTES.md` の `v0.1.8` を指定書式へ整理し、直前正式版から見た変更点に圧縮した。
- `scripts/extract-release-notes.mjs` を、本文中の `## ダウンロード` で節を切らず、次のバージョン見出しだけを境界にするよう修正した。
- Release workflowで通常利用者向けzipに単一exe署名ascを含め、`Spout2-LICENSE.txt`、`spout-capture.exe`、`SpoutLibrary.dll` を通常zipから除外する検証を追加した。
- Release workflowのartifact/upload対象を通常zip、単一exe署名asc、separated zip、AvatarBeacon source zip、release bodyだけにし、GitHub Release公開Assetは4種類だけuploadするようにした。
- 同じtagでRelease workflowを再実行した場合、許可リスト外の既存Release Assetを削除してからuploadするようにした。
- README、SPEC、アプリ内PGP確認手順を新しい配布物構成に合わせた。

## 検証

- `node scripts/check-frontend-template-literals.mjs`
- `node scripts/check-wails-api-surface.mjs`
- `node scripts/extract-release-notes.mjs v0.1.8 RELEASE_NOTES.md /tmp/release-body.md`
- `node scripts/extract-release-notes.mjs v0.1.8-rc99 RELEASE_NOTES.md /tmp/release-body-rc.md --fallback-rc-notes`
- `npm run build` in `src/frontend`
