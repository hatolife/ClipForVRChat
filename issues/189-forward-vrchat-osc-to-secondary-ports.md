# VRChatから受信したOSCを他アプリ向けに別ポートへ転送できるようにする

## 問題

VRChatから外部アプリへ送られるOSCは、通常ローカルの特定UDPポートで受信する。ClipForVRChatがその受信ポートを確保している状態では、同じポートを受信したい別アプリと競合する。

UDPの同一IP/同一ポートを複数アプリが安定して同時受信できる前提にはできないため、ClipForVRChatと他のOSC受信アプリを同時に使いたい場合、どちらかが受信できない、起動に失敗する、環境依存で片方だけ受け取る、といった問題が起きる。

## 期待する挙動

ClipForVRChatがVRChatからのOSC受信ポートを代表してlistenし、ClipForVRChat自身に必要なOSCは従来どおり処理する。

そのうえで、ClipForVRChatに関係しないOSC、またはユーザーが指定したOSCを、設定した別ポートへ転送できるようにする。他アプリはVRChatの元ポートではなく、その転送先ポートを受信することで競合を回避できる。

## 受け入れ条件

- 設定画面に `OSC` タブを追加し、自動撮影タブにあるOSCホスト、送信ポート、受信ポート、現在Pose保存の有効秒数、カメラOSCリセット、AvatarBeacon/player_local basis関連設定、AvatarBeacon受信状態を移動する。
- 自動撮影タブの説明文は、AvatarBeaconが別途アバターへ導入する必要のあるアバターギミックであること、Photo方式は連射時にシャッター音が鳴ることを明記する。
- OSC転送を有効/無効にできる設定がある。
- 転送先として `host` と `port` を1つ以上設定できる。
- VRChatから受信したOSC packetを、設定された転送先へUDP送信できる。
- ClipForVRChatが利用するOSC addressを転送するか、ClipForVRChatに関係ないOSCだけを転送するかを設定または仕様として明確にする。
- 転送先ポートがClipForVRChat自身の受信ポートと同じ場合は保存または起動時に警告し、ループや自己転送を避ける。
- 転送先が応答しない、送信に失敗する、port設定が不正、などの状態を診断ログで確認できる。
- 転送機能を無効にしている場合、既存のOSC受信、自動撮影、AvatarBeacon受信の挙動が変わらない。
- 他アプリ側が受信ポートを変更できない場合は、この機能では競合回避できないことをREADME/SPECまたはUIで説明する。

## メモ

- VRChatへOSCを送信する入力ポートと、VRChatから外部アプリへ送られる出力OSCの受信ポートは別問題として扱う。
- 競合しやすいのは、VRChatからのOSCを外部アプリがlistenするローカル受信ポート。
- ClipForVRChatが単一起動になっても、他のOSC受信アプリとのポート競合は残るため、本issueは #188 とは別の競合回避策として扱う。
- 他アプリが転送先ポートを指定できる場合に有効な回避策であり、他アプリがVRChat標準ポート固定の場合は解決できない。
- 2026-07-07: OSCタブ冒頭の同一ポート競合と転送構成の説明が短く、ユーザーにとって「何をどのポートへ向けるべきか」が分かりにくいため、UI文言をより丁寧にする。

## 実装案

- 既存のOSC受信処理でpacketを受け取った直後に、転送設定を参照して別UDP送信先へ複製送信する。
- OSC message単位ではなくpacket単位で転送すれば、bundleや型tagをできるだけ維持できる。
- 転送フィルタは初期実装では単純にする。
  - 全転送
  - ClipForVRChatが処理したaddressも含めて全転送
  - 将来対応として、address prefix allow/deny listを追加する
- 設定例:
  - `enabled`
  - `targets: [{ host: "127.0.0.1", port: 9101 }]`
  - `mode: "all"` または `mode: "unhandled_only"`
- 診断ログには、転送有効状態、転送先、直近送信件数、直近エラーを出す。

## 注意点

- 転送先に送ったOSCを再びClipForVRChatが受信する設定にするとループする可能性がある。受信ポートと同一の転送先、または明らかな自己転送は拒否または強い警告にする。
- UDP転送は配送保証がない。転送先アプリが起動していない場合でも送信自体は成功することがあるため、「相手が受け取った」ことまでは保証できない。
- 複数転送先を許可する場合、1つの転送先エラーで他の転送を止めない。
- ClipForVRChat内部処理より先に転送するか後に転送するかで、ログ上の順序や異常時の見え方が変わる。まずは「受信後に内部処理と転送をそれぞれ行う」方針でよい。
- addressを加工せず転送する。VRChat標準のaddressやAvatar Parameter名をClipForVRChat固有に変換しない。
- プライバシー上、Avatar Parameterやカメラ状態などVRChatから出るOSCを別アプリへ渡す機能になるため、設定UIでは有効化時に転送先を明示する。

## 検証観点

- ClipForVRChatがVRChat出力OSCポートをlistenしている状態で、別ポートへOSCが転送される。
- 他アプリまたはテスト受信器が転送先ポートでOSCを受信できる。
- 転送無効時はOSC packetが転送されない。
- 転送先を複数設定した場合、それぞれへ送信される。
- 転送先の1つが不正でも、他の転送先とClipForVRChat本体処理が継続する。
- 受信ポートと同じ転送先を設定した場合に警告または拒否される。
- AvatarBeacon受信、自動撮影、User Camera Pose受信の既存挙動が転送設定で壊れない。

## 実装内容

- `autoCapture.osc.forward` 設定を追加した。
  - `enabled`: OSC転送の有効/無効。
  - `mode`: `all` または `unhandled_only`。
  - `targets`: `host` / `port` の転送先配列。
- 転送はOSC messageを再構築せず、VRChatから受信したUDP payloadをpacket単位でそのまま複製送信する。
- `all` はClipForVRChatが処理するaddressも含めて転送する。
- `unhandled_only` はClipForVRChatが使う `/usercamera/*` 相当のUser Camera sampleと `/avatar/parameters/*` を転送しない。
- 転送先がClipForVRChat自身の受信endpointと同じ場合は起動時に除外し、診断ログへ `self_forward` として出す。
- 転送開始、無効状態、転送先除外、10秒ごとのsummary、終了時summaryを診断ログへ出す。
- 設定画面に `OSC` タブを追加し、自動撮影タブにあった以下を移動した。
  - OSCホスト
  - OSC送信ポート
  - OSC受信ポート
  - 現在Pose保存の有効秒数
  - カメラOSCリセット
  - プレイヤー基準の取得元
  - AvatarBeacon受信状態
- `OSC` タブにOSC転送の有効/無効、転送モード、複数転送先の追加/削除UIを追加した。
- 自動撮影タブの説明文に、AvatarBeaconはアバターへ導入する専用ギミックであること、Photo方式では連射時にVRChat側のシャッター音が鳴ることを明記した。

## 実装上の仕様

- 転送先アプリは、VRChat標準の出力OSCポートではなく、ClipForVRChatで設定した転送先ポートをlistenする。
- 他アプリ側の受信ポートを変更できない場合、この機能ではポート競合を解消できない。
- UDP転送のため、転送先アプリが起動していなくても送信側では成功扱いになる場合がある。受信成功までは保証しない。
- 1つの転送先で送信に失敗しても、他の転送先とClipForVRChat本体処理は継続する。
- OSCタブの冒頭説明では、VRChatからの受信ポートはClipForVRChatが代表してlistenし、他アプリは転送先として指定した別ポートをlistenする構成を案内する。

## 受け入れ条件の実装状況

- [x] 設定画面に `OSC` タブを追加し、自動撮影タブにあるOSCホスト、送信ポート、受信ポート、現在Pose保存の有効秒数、カメラOSCリセット、AvatarBeacon/player_local basis関連設定、AvatarBeacon受信状態を移動する。
- [x] 自動撮影タブの説明文は、AvatarBeaconが別途アバターへ導入する必要のあるアバターギミックであること、Photo方式は連射時にシャッター音が鳴ることを明記する。
- [x] OSC転送を有効/無効にできる設定がある。
- [x] 転送先として `host` と `port` を1つ以上設定できる。
- [x] VRChatから受信したOSC packetを、設定された転送先へUDP送信できる。
- [x] ClipForVRChatが利用するOSC addressを転送するか、ClipForVRChatに関係ないOSCだけを転送するかを設定または仕様として明確にする。
- [x] 転送先ポートがClipForVRChat自身の受信ポートと同じ場合は保存または起動時に警告し、ループや自己転送を避ける。
- [x] 転送先が応答しない、送信に失敗する、port設定が不正、などの状態を診断ログで確認できる。
- [x] 転送機能を無効にしている場合、既存のOSC受信、自動撮影、AvatarBeacon受信の挙動が変わらない。
- [x] 他アプリ側が受信ポートを変更できない場合は、この機能では競合回避できないことをUIで説明する。

## 実行済み検証

- `cd src && GOCACHE=/tmp/clipforvrchat-go-cache go test ./...`
- `node scripts/check-frontend-template-literals.mjs`
- `node scripts/check-wails-api-surface.mjs`
- `git diff --check`
- `cd src && GOOS=windows GOARCH=amd64 GOCACHE=/tmp/clipforvrchat-go-cache-windows go test -c -o /tmp/clipforvrchat-main.test .`
- `cd src && GOOS=windows GOARCH=amd64 GOCACHE=/tmp/clipforvrchat-go-cache-windows go test -c -o /tmp/clipforvrchat-appcore.test ./internal/appcore`

## 実機確認が必要な項目

- Windows実機で、VRChatから受信したOSCが設定した転送先ポートへ届くこと。
- 他OSC受信アプリ側の受信ポートを転送先へ変更し、ClipForVRChatと同時利用できること。
- `all` と `unhandled_only` の動作が期待どおりであること。
- 転送先を自己受信ポートにした場合にUI警告と診断ログが確認できること。
