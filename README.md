# ClipForVRChat

ClipForVRChat は、VRChat で外部画像URLを使いやすくするための Windows アプリです。

画像を VRChat の 2048x2048 制限に収まるように縮小し、必要に応じて Discord に投稿して、画像URLをクリップボードへコピーします。

## できること

- 画像を縦横 2048px 以内に縮小
- 縦横比を維持
- JPEG は JPEG、それ以外は PNG として出力
- 縮小画像をローカル保存
- Discord Webhook に画像を投稿
- 1枚だけなら投稿後の画像URLを自動でクリップボードへコピー
- 複数枚ならサムネイル付きのURL一覧を表示

## 使い方

クリップボード内の画像を処理:

```txt
ClipForVRChat.exe
```

画像ファイルを処理:

```txt
ClipForVRChat.exe image.png
ClipForVRChat.exe image1.png image2.jpg
```

設定画面を開く:

```txt
ClipForVRChat.exe config.json
```

初回起動時に `config.json` がない場合は、最初に設定画面が開きます。保存すると、そのまま画像処理を続行します。

起動後のウィンドウには画像ファイルをドラッグ&ドロップできます。複数画像もまとめて処理できます。

`config.json` をウィンドウにドロップすると、その設定ファイルの編集画面を開きます。

## 設定

設定は exe と同じ場所にある `config.json` に保存されます。

通常はJSONを手で編集せず、設定画面から変更する想定です。

主な設定:

- ローカル保存のON/OFF
- Discord投稿のON/OFF
- Discord Webhook URL
- 出力先フォルダ
- JPEG品質
- UI表示モード

## Discord設定

投稿先チャンネルで Discord Webhook を作成し、設定画面で Webhook URL を指定してください。

コピーされるURLは Discord 添付画像の直リンクです。
