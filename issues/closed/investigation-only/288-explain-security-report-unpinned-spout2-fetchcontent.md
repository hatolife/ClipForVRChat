# Spout2 FetchContent未固定セキュリティレポートを解説する

## 指示

> '/home/user/work/ClipForVRChat/reports/security/2026-07-07T17-53-16.203Z/87e03cb4699c8191b7cc0e5a4c3830d3'を日本語で解説して

## 文脈

指定されたファイルは、Codex Securityが生成した `Unpinned Spout2 FetchContent in release build` というhigh findingのレポートである。
`tools/spout-capture/CMakeLists.txt` がReleaseビルド時にSpout2をGitHubから `FetchContent` で取得し、`GIT_TAG 2.007.015` という可変タグ相当の指定で固定している点が指摘されている。

## 解釈

レポート本文を読み、指摘内容、攻撃シナリオ、根拠、影響、既存対策、推奨対応を日本語で分かりやすく説明する。
この依頼ではコード修正までは行わず、現時点の理解と対応方針を整理する。

## 問題

- Release workflowがSpout helperをビルドする際、外部リポジトリのSpout2をビルド時に取得している。
- 取得対象が不変のcommit SHAや検証済みarchive hashではなく、動かされ得るタグ名で指定されている。
- 取得したSpout2側のCMake/build処理はRelease runner上で実行される。
- 生成された `spout-capture.exe` は本体exeへの埋め込みと分離版zip同梱に使われるが、Detached署名は単一の本体exeのみを対象としている。

## 期待する挙動

- レポートの意味、影響、根拠、対応方針が日本語で説明される。
- 今後修正する場合の要点が分かる。

## 受け入れ条件

- [x] findingの要旨を説明する。
- [x] なぜhigh扱いかを説明する。
- [x] レポート内の検証根拠を説明する。
- [x] 提案される修正方針を説明する。

## 調査結果

- findingは、`tools/spout-capture/CMakeLists.txt` が `FetchContent` で `https://github.com/leadedge/Spout2.git` を取得し、`GIT_TAG 2.007.015` で固定している点を指摘している。
- `FetchContent` はRelease workflowのCMake configure/build中に実行されるため、取得先のCMakeコードもRelease runner上で処理される。
- upstream repository、タグ、または依存配送経路が侵害されると、Releaseビルド時に悪意あるCMake/build処理が実行され、`spout-capture.exe` や成果物が改変され得る。
- 生成された `spout-capture.exe` は通常版では本体exeへ埋め込まれ、分離版zipでは同梱され、製品のStream撮影機能から実行される。
- 現在の署名は単一配布用の `ClipForVRChat.exe` のDetached署名であり、分離版helper単体またはzip全体の署名ではない。通常版についても、署名はReleaseビルド後の成果物に対して行われるため、ビルド時に取得する依存関係の改変そのものを防ぐものではない。
- 修正方針は、Spout2取得を不変commit SHAへ固定する、hash検証付きarchiveまたはvendor済みソースにする、加えてhelperまたはzip全体の署名・検証を検討することである。
