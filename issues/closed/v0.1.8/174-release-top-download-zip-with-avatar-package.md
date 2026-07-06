# Release本文先頭にダウンロード導線を置き、通常配布zipへ必要物をまとめる

## 問題

GitHub Release本文では、利用者が最初にどのファイルをダウンロードすればよいかが分かりにくい。
また、`avatar_osc` を使うにはアプリ本体だけでなく専用アバターギミック用 `unitypackage` と説明書が必要になるが、現状のRelease Assetsと本文導線はこの配布単位を前提にしていない。
アバターギミックの `.unitypackage` はCIで生成せず、リリース担当者がUnity上で手作業作成してGitHub Releaseへ添付する方針にする。

## 期待する挙動

GitHub Release本文の一番上に、通常利用者向けの最優先導線として次のような見出しリンクを表示する。

```md
# [vX.Y.Zのダウンロード](zipのリンク)
```

リンク先の通常配布zipには、単一化した `ClipForVRChat.exe`、説明書 `README.md`、必要なライセンスファイルを含める。
アバターギミック用 `.unitypackage` は手動作成・手動添付のRelease Assetとして通常配布zipの近くに案内し、CI生成のアバターギミック元ファイルzipは検証・再生成用Assetとして残す。
分離版zipや個別exe、署名、sha256、build metadataなどの検証用Assetsは必要に応じて残すが、Release本文の主導線は通常配布zipと手動添付 `.unitypackage` の組み合わせが分かる形にする。

## 受け入れ条件

- [x] Release本文の最上部に `# [vX.Y.Zのダウンロード](...)` 形式の通常配布zipリンクが表示される。
- [x] RCタグ `vX.Y.Z-rcN` でも、リンク文言とリンク先ファイル名がRC版に正しく置換される。
- [x] 通常配布zipに、単一化した `ClipForVRChat.exe` が含まれる。
- [x] アバターギミック用 `.unitypackage` は、リリース担当者が手動作成してGitHub Release Assetへ添付する。
- [x] Release本文で、通常配布zipと手動添付 `.unitypackage` の両方を利用者向け導線として分かりやすく表示する。
- [x] CI生成のアバターギミック元ファイルzipは、`.unitypackage` の再生成・検証用Assetとして添付される。
- [x] 通常配布zipに、利用者向け説明書 `README.md` が含まれる。
- [x] 通常配布zipに、アプリ本体と同梱物に必要なライセンスファイルが含まれる。
- [x] Release workflowで通常配布zip、zipのsha256、必要ならzip署名を生成・添付する。
- [x] Release workflowの検証で、通常配布zip内の必須ファイルと不要ファイル混入を確認する。
- [x] Release本文では通常配布zipを主導線にし、分離版zipや個別exeは検証・切り分け用として区別して説明する。
- [x] README / SPEC / RELEASE_NOTES の配布物説明が、実際のRelease Assetsとzip内ファイル一覧に一致する。

## メモ

- 2026-07-04: Release workflowで `ClipForVRChat-vX.Y.Z-windows-amd64.zip` と `.zip.sha256` を生成・添付するようにした。zipには `ClipForVRChat.exe`、`README.md`、`LICENSE`、`Spout2-LICENSE.txt`、`Release-signing-public-key.url` を含め、`spout-capture.exe` / `SpoutLibrary.dll` の直接混入は検証で拒否する。
- 2026-07-04: Release本文先頭に通常配布zipリンクを追加し、RC向け置換も `scripts/extract-release-notes.mjs` で確認した。
- 2026-07-04: アバターギミック用 `.unitypackage` のファイル名は `AvatarBeacon-vX.Y.Z.unitypackage` とし、Release担当者がUnityで手動作成・手動添付する。CIで作る `AvatarBeacon-vX.Y.Z-source.zip` は再生成・検証用Assetとして残す。source zip生成と手動 `.unitypackage` 作成運用の詳細は [#175](../../175-package-avatar-gimmick-source-zip.md) で扱う。
- 通常配布zipを主導線にする場合でも、改竄確認の導線が分かりにくくならないよう、zip sha256や署名対象を合わせて整理する。
