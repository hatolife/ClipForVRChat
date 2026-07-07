# Security future recommendations

この文書は、2026-07-08時点のsecurity findings整理で、今回の即時方針からは外したが将来検討すべき推奨事項を記録する。
実装判断済みの短期方針は各issueを正とし、この文書は後続改善の候補として扱う。

## Issue 313: Discord履歴tokenのOS秘密情報ストア分離

- 対象: `issues/313-protect-history-discord-token-storage.md`
- 将来案: Discord webhook tokenをhistory JSONへ保存せず、Windows Credential ManagerなどOSの秘密情報ストアへ分離する。
- 目的: 再起動後のDiscord履歴削除機能を維持しつつ、history JSONとUI stateにtokenを置かない。
- 検討事項: OS別実装、portable運用、バックアップ/移行、credential削除、テスト環境でのmock、既存履歴からの移行。

## Issue 314: issue secret redaction検査

- 対象: `issues/314-redact-secrets-in-verbatim-issue-quotes.md`
- 将来案: issue作成・更新時に、Discord Webhook URL、token、ローカルユーザー名を含むパス、QR URL、診断ログsecretらしき値を検出する検査スクリプトを追加する。
- 目的: AGENTS.mdのredaction例外を運用だけに頼らず、コミット前に漏えいを検出する。
- 検討事項: 誤検知時の例外指定、既存issue内の歴史的記録、Windows pathの扱い、検査対象ディレクトリ、CIで必須にするかローカル補助にするか。

## Issue 316: keyless/KMS release signing

- 対象: `issues/316-isolate-release-signing-secrets-from-build-job.md`
- 将来案: 長期GPG秘密鍵をGitHub Actions runnerへ渡す方式から、Sigstore/keyless signing、KMS、HSMなどへ移行する。
- 目的: Release signing secretのrunner露出をなくし、署名鍵漏えい時の長期影響を減らす。
- 検討事項: 利用者向け検証手順、公開Asset名、GitHub Release本文、既存GPG署名との移行期間、Windows利用者が検証しやすい導線。

