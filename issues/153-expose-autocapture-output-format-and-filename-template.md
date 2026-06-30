# 自動撮影の出力形式/ファイル名テンプレートを設定画面へ出す

## 問題

`AutoCaptureOutputConfig` には `imageFormat` と `filenameTemplate` があり、保存処理でも使われている。
しかし自動撮影タブには出力先とsidecar JSON設定だけがあり、出力形式とファイル名テンプレートを編集できない。
`todo.md` では編集可能扱いになっており、設定UIと実装が一致していない。

## 期待する挙動

自動撮影タブから、Stream/Spout保存画像の出力形式とファイル名テンプレートを設定できる。
Photo方式ではVRChat側の保存名を使うことが分かり、設定が適用される範囲を誤解しない。

## 受け入れ条件

- 自動撮影タブに `imageFormat` の選択UIを追加する。対応形式は実装済みの `png` / `jpg` / `jpeg` に限定する。
- 自動撮影タブに `filenameTemplate` の入力UIを追加し、利用可能なplaceholderを説明する。
- 空値や不正拡張子は保存時またはNormalize時に安全な初期値へ戻す。
- Stream/Spout方式では設定した形式とテンプレートが保存ファイルに反映される。
- Photo方式ではこの設定がVRChat標準写真のファイル名には反映されないことをUIで説明する。
- `todo.md` とREADME/SPECの記述を実装状態に合わせる。

## 対応内容

- v0.1.8実装に合わせて対応済み。詳細は `feat/v0.1.8-resolve-issues` の実装、README、SPEC、RELEASE_NOTESを参照。
