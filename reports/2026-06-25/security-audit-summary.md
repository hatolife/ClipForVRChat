# Security Audit Summary

調査日: 2026-06-25

対象コミット: `8e46ffdefaff79eb86fd7ef02606c45ff9df4ec7`

## 全体サマリ

ClipForVRChatは、ローカル画像を縮小し、必要に応じてDiscord Webhookへ投稿するWindows向けWailsアプリである。サーバーとして待ち受ける処理はなく、主な攻撃面はローカル入力ファイル、設定・履歴ファイル、Webhook URL、Discord/GitHub通信、Release workflow、診断データ作成機能である。

画像入力にはファイルサイズとピクセル数上限があり、Discord/GitHub通信にはタイムアウトが設定されている。Webhook URLとDiscord添付URLはホストとパスを検証しており、履歴・設定ファイルはGo側で `0600` 相当の保存を試みている。CIには `npm audit`、`govulncheck`、`go test` が入っている。

一方で、`govulncheck` が暗号関連の既知脆弱性 `GO-2026-4550` を到達可能として検出した。また、不具合報告用データで暗号化前zipを残す仕様はユーザー確認用として妥当だが、そのzip内の `config.json` / `history.json` にDiscord Webhook URLや削除用tokenがそのまま含まれると誤添付時の漏えいリスクが残る。平文zipは残しつつ、添付前にユーザーが確認できる安全化済みデータとして扱うべきである。

## 重大度別件数

- Critical: 0
- High: 0
- Medium: 3
- Low: 4
- Info: 2

## 最優先で修正すべき項目

1. `github.com/cloudflare/circl@v1.6.2` の `GO-2026-4550` を解消する。
2. 診断データの確認用zipにWebhook URLやDiscord削除tokenが入らないよう、ダミー値へ置換する。
3. `history.json` に保存するDiscord tokenを診断データへ含めない方針を明文化する。
4. Release workflowの第三者ActionをSHA pinningし、権限をjob単位へ寄せ、Environment protectionを検討する。
5. `OpenURL` の許可URLをアプリが明示的に管理する。

## リリース可否の暫定判断

現時点では正式リリース不可と判断する。

理由は、`govulncheck` が到達可能な既知脆弱性を検出し、CI/Release workflow上でも同じチェックが失敗する可能性が高いためである。少なくとも `cloudflare/circl` を修正版へ更新し、`govulncheck ./...` が成功することをリリース前条件にするべきである。

診断データ確認用zip自体は、ユーザーが送信前に内容を確認するために残す仕様でよい。ただし、Webhook URL、Discord token、その他の秘密情報はダミー値に置換し、平文zipを見ても秘密情報が分からない状態にすることを直近で修正するべきである。

## 残存リスク

- Discord Webhook URLは秘密情報であり、config/history/log/diagnosticsの扱いに依存する。
- ローカルユーザーがconfig/historyを改ざんした場合、同一ユーザー権限内で任意の保存先・削除対象へ誘導できる。
- WailsブリッジはフロントエンドからGoメソッドを呼べるため、将来HTML挿入や外部コンテンツ表示を追加する場合は影響が大きくなる。
- Releaseの改ざん耐性はPGP署名とsha256で補強されているが、コード署名証明書やSBOMはない。
- Windowsの実ファイルACLがGoの `0600` 指定だけで期待通り最小化されるかは実機確認が必要。WindowsではUnixの `0600` と同じ意味ではなく、実際のアクセス権はDACLで決まる。

## 追加調査が必要な項目

- `GO-2026-4550` が実際のcv25519公開鍵暗号化経路へ与える影響範囲。
- Windows実機上での `config.json`、`history.json`、`logs/`、`diagnostics/` のACL。
- 診断データzip内のWebhook URL、Discord token、必要に応じたURL類をどこまでダミー化するかの運用確認。
- GitHub Actions Release実行時のEnvironment protection、手動承認、tag作成権限。
- WebView2/Wailsのセキュリティ設定と、外部コンテンツを読み込まない前提の継続確認。

## セキュリティ上の総評

ローカルデスクトップアプリとしては、危険になりやすいネットワーク投稿先と画像入力に基本的な制限が入っている。CIにも依存関係監査が組み込まれており、Release成果物の署名・ハッシュも整備されている。

今後の重点は、秘密情報を含む診断・履歴データの扱い、Release supply chainの硬化、依存関係脆弱性の継続的な解消である。
