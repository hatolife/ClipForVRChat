# rc25でSpout PNG保存が失敗し、AvatarBeacon受信状態の説明が冗長

## 問題

`v0.1.8-rc25` のStream方式テスト撮影で、Spout helperが `PNG encoder does not support RGBA` を返し、撮影結果が失敗する。
ログでは `auto-capture spout capture helper error: code="png_write_error"` として記録されている。

また、自動撮影タブのAvatarBeacon受信状態に `coord/* と forward/* がログsummaryに出るかで確認します。専用ギミック未導入、OSC無効、parameter欠落、鮮度切れでは自動追従できません。` のような常時説明が表示され、通常操作時には情報量が多すぎる。

## 期待する挙動

Spout helperはVRChat Stream Cameraから取得したRGBAフレームを、Windows WIC PNG encoderで保存できる形式へ変換してPNGとして保存する。

AvatarBeacon受信状態は、状態、最終受信時刻、position/yaw、必要時のエラーだけを表示し、ログsummary確認を促す常時説明は出さない。

## 受け入れ条件

- [x] Spout helperがWIC PNG encoderのRGBA非対応環境でもPNGを書き出せる。
- [ ] `PNG encoder does not support RGBA` でStream方式テスト撮影が失敗しない。
- [x] AvatarBeacon受信状態から、常時表示のログsummary確認説明を削除する。
- [x] エラー時の警告表示やstale診断は、必要な範囲で維持する。
- [x] frontend template literal check、Wails API surface check、Go test、frontend buildが通る。
- [ ] Release workflowでSpout helperビルドと成果物検証が通る。

## 調査メモ

- 2026-07-04: ユーザー提供ログでは `07:42:08` に `code="png_write_error" message="PNG encoder does not support RGBA"` が記録されている。
- Spout helperは `ReceiveImage(..., GL_RGBA, ...)` でRGBAバイト列を受け取り、WIC PNG frameへ `GUID_WICPixelFormat32bppRGBA` を要求している。環境によってPNG encoderがこのpixel formatを受け付けず、`SetPixelFormat` 後のformat不一致で失敗している。
- 2026-07-04: 修正ではSpoutから受け取るRGBAは維持し、PNG書き出し時にWICが受け付けやすい `32bppBGRA`、`32bppPBGRA`、`24bppBGR` へ変換して `WritePixels` する。
- 2026-07-04: ローカルMinGWではSpout2側の `Shellapi.h` 大文字includeで全体ビルドが失敗するため、最終的なSpout helper全体ビルド確認はGitHub ActionsのWindows/MSVC Release workflowで行う。今回変更した `tools/spout-capture/main.cpp` のobject compileは通過した。
