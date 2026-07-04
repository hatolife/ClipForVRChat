# 064 master以外のブランチCIと非正式タグのdraft release対応

## 状態

完了

## 問題

CI workflow の push 対象が `master` / `main` に限定されているため、作業ブランチへ push してもCIが走らない。
また、Release workflow は `v*` タグで起動するが、`v0.1.6a` のような `vX.Y.Z` 以外のタグでも通常Releaseとして公開される可能性がある。

## 期待する挙動

master以外のブランチへ push した場合もCIが走る。
`vX.Y.Z` 形式の正式タグは通常Releaseとして作成し、それ以外の `v*` タグはdraft Releaseとして作成する。

## 受け入れ条件

- CI workflow の push がすべてのブランチを対象にする。
- Release workflow が `vX.Y.Z` 形式のタグかどうかを判定する。
- `vX.Y.Z` 以外の `v*` タグでは GitHub Release が draft になる。
- `RELEASE_NOTES.md` に該当タグの項目がないdraftタグでもRelease本文を生成できる。
- 新しいチケットが `issues/README.md` に掲載されている。
