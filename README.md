# ClipForVRChat

ClipForVRChat は、VRChat で外部画像URLを使いやすくするための Windows アプリです。

画像を VRChat の 2048x2048 制限に収まるように縮小し、必要に応じて Discord に投稿して、画像URLをクリップボードへコピーします。

## できること

- 画像を縦横 2048px 以内に縮小
- 縦横比を維持
- 出力形式を PNG または JPG から選択
- 縮小画像をローカル保存
- Discord Webhook に画像を投稿
- 1枚だけなら投稿後の画像URLを自動でクリップボードへコピー
- 複数枚ならサムネイル付きのURL一覧を表示
- サムネイルをクリックして画像URLを再コピー

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
- 出力形式
- JPEG品質
- UI表示モード

PNG出力のとき、JPEG品質は使用されません。

## 画面

起動後の画面では、左側に画像ドロップ領域、右側に結果が表示されます。

結果が多い場合も、スクロールするのは結果部分だけです。

ドラッグ中はウィンドウ全体にドロップ案内が表示され、どの画面からでも画像や `config.json` をドロップできます。

複数画像を入力した場合、処理開始時点で全画像分の枠が表示され、完了したものからサムネイルへ切り替わります。

結果サムネイルにマウスを重ねるとURLコピーの案内が表示されます。結果一覧は右上のクリアボタンで消せます。

使い方画面では、画像の処理手順やDiscord Webhook URLの発行方法へのリンクを確認できます。

情報画面では、バージョン、GitHub URL、使用しているOSSライセンスを確認できます。GitHub URLはクリックするとブラウザで開きます。

## 配布zip

Release の `ClipForVRChat-vX.Y.Z-windows-amd64.zip` には、exe、README、LICENSE が含まれます。

## Discord設定

投稿先チャンネルで Discord Webhook を作成し、設定画面で Webhook URL を指定してください。

コピーされるURLは Discord 添付画像の直リンクです。
