# Stream方式UIとドキュメントをSpout前提へ更新する

## 問題

現在の自動撮影タブやRelease Notesには、Stream方式をFFmpegで切り出す説明が残っている。
しかしFFmpeg/gdigrabは画面キャプチャであり、VRChat Stream Camera映像そのものを取得する主方式としては不適切。

## 期待する挙動

自動撮影タブ、README、SPEC、RELEASE_NOTESの説明を、Stream Camera(Spout)方式が主経路である内容へ更新する。
ユーザーがFFmpeg設定を調整すればStream Cameraが直接取れると誤解しないUIにする。

## 受け入れ条件

- 自動撮影タブの説明で、Stream方式はVRChat Stream CameraのSpout映像を直接受信して保存する方式だと説明する。
- `ffmpegパス`、`ffmpeg入力引数`、`ffmpegをインストール` をStream方式の主設定として表示しない。
- FFmpeg経路を残す場合は、legacy/debugの画面キャプチャfallbackとして明示し、通常利用の導線から外す。
- `RELEASE_NOTES.md` から「Stream方式はffmpegでVRChat Stream Camera映像を切り出す」という誤解を招く記述を修正する。
- README/SPECに、Stream Camera(Spout)の前提、VRChat OSC有効化、sender選択、Photo方式との違いを記載する。
- 既存機能の説明、Discord投稿、Presence、Photo方式の説明を壊さない。

## 実装メモ

- v0.1.8の本命要件はStream Camera(Spout)であり、Photo方式はシャッター音ありのフォールバックとして扱う。
- FFmpeg導入ボタンを残すかどうかは実装時に判断するが、残す場合も主導線にはしない。

## 対応内容

- v0.1.8実装に合わせて対応済み。詳細は `feat/v0.1.8-resolve-issues` の実装、README、SPEC、RELEASE_NOTESを参照。

## 再監査メモ

- 2026-07-01: [#164](164-audit-v018-completed-items.md) の再監査で未達が見つかったため、完了扱いを取り消して `要対応` に戻した。
