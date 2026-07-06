# alpha/beta/rcを分けたバージョニング規則を定義する

## 問題

`v0.1.8-rcW` が44番台まで進み、RCが本来の「正式リリース直前の最終確認版」ではなく、機能開発や最低機能確認の配布版として使われている。
また、`vX.Y.Z` で実装すると決めた機能について、一部だけ実装して難しい残りを次バージョンへ送る判断や、alpha/beta段階を経ずにRC扱いへ進める判断が、明示合意なしに行われている。

## 期待する挙動

- 機能開発中の配布版は `vX.Y.Z-aW` を使う。
- 機能が揃った後の機能チェック版は `vX.Y.Z-bW` を使う。
- 正式リリース可能性を広く確認する段階だけ `vX.Y.Z-rcW` を使う。
- 正式版は `vX.Y.Z` を使う。
- alpha、beta、rc、正式版のすべてでCI/CDが走り、GitHub Releaseと配布成果物が作成される。
- `vX.Y.Z` の対象として合意した機能は、ユーザーが明示的に延期・縮小を承認しない限り、そのバージョンのリリース対象に残す。
- alpha/beta/rc/正式版の段階移行は、ユーザーの明示的な合意に基づいて行う。

## 受け入れ条件

- `AGENTS.md` に alpha/beta/rc/正式版の用途と昇格条件が記載されている。
- Release workflow が `vX.Y.Z-aW`、`vX.Y.Z-bW`、`vX.Y.Z-rcW`、`vX.Y.Z` を受け付ける。
- `vX.Y.Z-aW`、`vX.Y.Z-bW`、`vX.Y.Z-rcW` は GitHub Release の prerelease として作成される。
- `vX.Y.Z` は GitHub Release の通常Releaseとして作成される。
- Release notes 抽出が alpha/beta/rc タグから対応する `vX.Y.Z` の本文へ fallback できる。
- 対象バージョンの未完了機能を、ユーザー合意なしに次バージョン送りにしないルールが明文化されている。
- alpha/beta/rc/正式版への昇格を、ユーザー合意なしに行わないルールが明文化されている。

## 作業メモ

- `AGENTS.md` に、alpha/beta/rc/正式版の用途、CI/CD対象、ユーザー合意なしのスコープ縮小・次バージョン送り禁止、RC昇格禁止を追記した。
- Release workflow が `vX.Y.Z-aW`、`vX.Y.Z-bW`、`vX.Y.Z-rcW`、`vX.Y.Z` を受け付け、alpha/beta/rcを prerelease、正式版を通常Releaseとして扱うようにした。
- Release notes 抽出が alpha/beta/rc から正式版見出しへfallbackできるようにした。
