# Release workflow のタグ名コマンドインジェクション対策

## 問題

Release workflow が `v*` タグで起動し、タグ名を PowerShell / Bash のコマンド文字列へ直接埋め込んでいる。悪意あるタグ名により文字列を抜けて任意コマンドが実行され、Release 成果物が改ざんされる可能性がある。

## 期待する挙動

タグ名は workflow 内で安全な文字種へ検証し、シェルへは環境変数として渡してデータとして扱う。成果物名や Release 作成に使うタグ名も検証済みの値だけを使う。

## 受け入れ条件

- Release workflow の PowerShell / Bash スクリプト内で `${{ github.ref_name }}` を直接展開しない。
- `vX.Y.Z`、`vX.Y.Z-rcW`、検証用の安全な `v*` タグのみを許可し、危険な文字を含むタグでは失敗する。
- 成果物名、署名、ハッシュ、メタデータ、Release 作成は検証済みタグ名を使う。
