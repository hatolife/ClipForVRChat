# Threat Model

## 保護資産

- Discord Webhook URL、Webhook token、Discord投稿URL。
- ローカル画像、スクリーンショット、VRChat写真、変換後画像、履歴。
- VRChat output log由来のWorld/Instance/Avatar/User情報。
- 診断ログ、OSCログ、診断パッケージ。
- Release署名鍵、CI artifact、Release asset。
- Spout helperとSpoutLibrary.dll、AvatarBeacon unitypackage。

## 想定ユーザー

- VRChat利用者。
- Discord Webhookを設定する一般ユーザー。
- 自動撮影、AvatarBeacon、OSC転送を検証するユーザー。
- Release成果物を確認する開発者。

## 想定攻撃者

- 悪意ある画像/QRコードをユーザーへ処理させる第三者。
- ユーザーのPC上で同じ権限のファイル配置やPATH汚染ができるローカル攻撃者。
- 不正な設定ファイルをユーザーへ配布する攻撃者。
- CI/CDや外部依存、Release assetの供給網を狙う攻撃者。

## 信頼境界

- ユーザー入力ファイルとアプリ内部処理。
- 設定ファイルと外部実行ファイル起動。
- Wails frontendとGo backend API。
- ローカルPCとDiscord/GitHub/BOOTH等の外部通信。
- GitHub Actions build/sign/package/release jobs。
- 本体exeとSpout helperプロセス。

## 入力経路

- GUI drag & drop、Clipboard画像、CLI引数。
- `config.json`、settings draft、history。
- VRChat写真/スクリーンショット監視フォルダ。
- VRChat output log。
- OSC UDP受信、OSC debug send入力。
- Spout sender名、Spout frame。
- GitHub Releases APIレスポンス。
- Release workflowで取得するAvatarBeacon packageと署名。

## 出力経路

- 変換画像、auto capture画像、sidecar JSON、EXIF/PNG metadata。
- Discord Webhook投稿。
- クリップボードへのURLコピー。
- 履歴、設定、診断ログ、OSCログ、診断zip/gpg。
- UDP OSC送信/転送。
- GitHub Release assets。

## 権限境界

- アプリはWindowsで管理者権限起動を拒否する実装がある。
- ローカルファイルはユーザー権限で読み書きする。
- Discord投稿はWebhook URLの権限で実行される。
- Release署名secretはGitHub Actions sign jobに限定されている。

## Attack Surface

- 画像デコーダとQRコード検出。
- Discord Webhook URL検証とHTTP投稿。
- AutoPhoto/スクリーンショット監視。
- OSC受信、転送、debug送信。
- Spout helper起動、SpoutLibrary.dllロード、PNG書き出し。
- 診断パッケージ生成と暗号化。
- GitHub Actions Release pipeline。

## STRIDE観点

- Spoofing: Discord/GitHub/AvatarBeacon packageの配布元なりすまし、Spout helper同名ファイル。
- Tampering: 設定ファイル、history、Release artifact、helper/DLLの差し替え。
- Repudiation: 診断ログには操作ログが残るが、ユーザー操作と外部変更の完全な証明はできない。
- Information Disclosure: Webhook URL、VRChat同席情報、画像、QRコードURL、診断zip。
- Denial of Service: 巨大画像、巨大フォルダ暗号化、Spout sender異常寸法、OSCログ増加。
- Elevation of Privilege: 管理者権限起動は拒否するが、PATH検索由来の同権限コード実行リスクは残る。

## 優先対策

- Spout helperとWindows固定コマンドのPATH検索依存を削減する。
- 依存脆弱性チェックをCIで継続し、Release前に結果を確認する。
- Release workflowのAction SHA pinningを検討する。
- ユーザー向け文書でWebhook、診断zip、OSC転送、外部helper指定の注意を明確にする。

## セキュリティ上の前提条件

- ユーザーは公式Releaseまたは公式配布元から入手する。
- 設定ファイル、helper、DLLは信頼できるものだけを使う。
- Discord Webhook URLは秘密として扱う。
- 診断データは内容確認後、暗号化された `.gpg` を共有する。

## 想定外としたリスク

- VRChat、Discord、GitHub、BOOTH、Windows自体への攻撃。
- ユーザーPCが既にマルウェアに完全掌握されている場合。
- 外部サービスのアカウント乗っ取り。
