# 039 バージョンとリビジョンを分けてビルドに埋め込む

## 状態

完了

## 問題

現在は `main.version` にタグ名とコミットIDを連結して埋め込んでおり、バージョンとリビジョンを個別に扱いにくい。

## 期待する挙動

ビルド時に `version` と `revision` を別々に埋め込み、表示時は `version.revision` のように組み立てられる。

## 受け入れ条件

- `version` と `revision` を別々の変数として持つ。
- ldflagsで `main.version` と `main.revision` を設定できる。
- ldflagsがない場合はGoのビルド情報からrevisionを取得できる。
- ローカルWSLビルド、GitHub ActionsのReleaseビルド、CIビルドが新しいldflagsを使用する。

## 対応内容

- `version` と `revision` を別変数に分け、表示時は `version.revision` に組み立てる。
- Goのビルド情報から `vcs.revision` と `vcs.modified` をfallbackとして取得する。
- WSLローカルビルド、GitHub Actions Release、CIのldflagsを `main.version` / `main.revision` に分離する。
