# Release本文先頭にダウンロード導線を置き、通常配布zipへ必要物をまとめる

## 問題

GitHub Release本文では、利用者が最初にどのファイルをダウンロードすればよいかが分かりにくい。
また、`avatar_osc` を使うにはアプリ本体だけでなく専用アバターギミック用 `unitypackage` と説明書が必要になるが、現状のRelease Assetsと本文導線はこの配布単位を前提にしていない。

## 期待する挙動

GitHub Release本文の一番上に、通常利用者向けの最優先導線として次のような見出しリンクを表示する。

```md
# [vX.Y.Zのダウンロード](zipのリンク)
```

リンク先の通常配布zipには、単一化した `ClipForVRChat.exe`、アバターギミック用 `unitypackage`、説明書 `README.md`、必要なライセンスファイルを含める。
分離版zipや個別exe、署名、sha256、build metadataなどの検証用Assetsは必要に応じて残すが、Release本文の主導線は通常配布zipにする。

## 受け入れ条件

- [ ] Release本文の最上部に `# [vX.Y.Zのダウンロード](...)` 形式の通常配布zipリンクが表示される。
- [ ] RCタグ `vX.Y.Z-rcN` でも、リンク文言とリンク先ファイル名がRC版に正しく置換される。
- [ ] 通常配布zipに、単一化した `ClipForVRChat.exe` が含まれる。
- [ ] 通常配布zipに、アバターギミック用 `unitypackage` が含まれる。
- [ ] 通常配布zipに、利用者向け説明書 `README.md` が含まれる。
- [ ] 通常配布zipに、アプリ本体と同梱物に必要なライセンスファイルが含まれる。
- [ ] Release workflowで通常配布zip、zipのsha256、必要ならzip署名を生成・添付する。
- [ ] Release workflowの検証で、通常配布zip内の必須ファイルと不要ファイル混入を確認する。
- [ ] Release本文では通常配布zipを主導線にし、分離版zipや個別exeは検証・切り分け用として区別して説明する。
- [ ] README / SPEC / RELEASE_NOTES の配布物説明が、実際のRelease Assetsとzip内ファイル一覧に一致する。

## メモ

- 現時点ではissue化のみ。実装は行わない。
- アバターギミック用 `unitypackage` の生成元、ファイル名、ライセンス表記、同梱可否は実装時に確定する。
- 通常配布zipを主導線にする場合でも、改竄確認の導線が分かりにくくならないよう、zip sha256や署名対象を合わせて整理する。
