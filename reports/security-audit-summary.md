# Security Audit Summary

調査日: 2026-06-25

対象コミット: `8e46ffdefaff79eb86fd7ef02606c45ff9df4ec7`

## 全体サマリ

ClipForVRChatは、ローカル画像を縮小し、必要に応じてDiscord Webhookへ投稿するWindows向けWailsアプリである。サーバーとして待ち受ける処理はなく、主な攻撃面はローカル入力ファイル、設定・履歴ファイル、Webhook URL、Discord/GitHub通信、Release workflow、診断データ作成機能である。

画像入力にはファイルサイズとピクセル数上限があり、Discord/GitHub通信にはタイムアウトが設定されている。Webhook URLとDiscord添付URLはホストとパスを検証しており、履歴・設定ファイルはGo側で `0600` 相当の保存を試みている。CIには `npm audit`、`govulncheck`、`go test` が入っている。

一方で、`govulncheck` が暗号関連の既知脆弱性 `GO-2026-4550` を到達可能として検出した。また、不具合報告用データで暗号化前zipを残す仕様はユーザー確認には有用だが、zip内の `config.json` / `history.json` にはDiscord Webhook URLや削除用tokenが含まれ得るため、誤添付時の漏えいリスクが残る。

## 重大度別件数

- Critical: 0
- High: 0
- Medium: 3
- Low: 4
- Info: 2

## 最優先で修正すべき項目

1. `github.com/cloudflare/circl@v1.6.2` の `GO-2026-4550` を解消する。
2. 診断データの確認用zipにWebhook URLやDiscord削除tokenが平文で入らないようにする。
3. `history.json` に保存するDiscord tokenの保持期間・用途・削除導線を見直す。
4. Release workflowの第三者ActionをSHA pinningし、権限をjob単位へ寄せる。
5. `OpenURL` の許可URLをアプリが明示的に管理する。

## リリース可否の暫定判断

現時点では正式リリース不可と判断する。

理由は、`govulncheck` が到達可能な既知脆弱性を検出し、CI/Release workflow上でも同じチェックが失敗する可能性が高いためである。少なくとも `cloudflare/circl` を修正版へ更新し、`govulncheck ./...` が成功することをリリース前条件にするべきである。

診断データ確認用zipの平文秘密情報リスクは、リリースブロッカーにするかは運用判断が必要だが、ユーザーがGitHub Issueへ誤ってzipを添付する可能性を考えると、直近で修正することを推奨する。

## 残存リスク

- Discord Webhook URLは秘密情報であり、config/history/log/diagnosticsの扱いに依存する。
- ローカルユーザーがconfig/historyを改ざんした場合、同一ユーザー権限内で任意の保存先・削除対象へ誘導できる。
- WailsブリッジはフロントエンドからGoメソッドを呼べるため、将来HTML挿入や外部コンテンツ表示を追加する場合は影響が大きくなる。
- Releaseの改ざん耐性はPGP署名とsha256で補強されているが、コード署名証明書やSBOMはない。
- Windowsの実ファイルACLがGoの `0600` 指定だけで期待通り最小化されるかは実機確認が必要。

## 追加調査が必要な項目

- `GO-2026-4550` が実際のcv25519公開鍵暗号化経路へ与える影響範囲。
- Windows実機上での `config.json`、`history.json`、`logs/`、`diagnostics/` のACL。
- 診断データzip内に含めるべき情報と、秘密情報を除外しても調査可能かの運用確認。
- GitHub Actions Release実行時の環境保護、手動承認、tag作成権限。
- WebView2/Wailsのセキュリティ設定と、外部コンテンツを読み込まない前提の継続確認。

## セキュリティ上の総評

ローカルデスクトップアプリとしては、危険になりやすいネットワーク投稿先と画像入力に基本的な制限が入っている。CIにも依存関係監査が組み込まれており、Release成果物の署名・ハッシュも整備されている。

今後の重点は、秘密情報を含む診断・履歴データの扱い、Release supply chainの硬化、依存関係脆弱性の継続的な解消である。
