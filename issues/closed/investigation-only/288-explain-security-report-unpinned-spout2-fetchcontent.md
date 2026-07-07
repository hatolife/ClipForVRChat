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

## 追加調査

### 指示

> Spout2取得を不変commit SHAへ固定する、hash検証付きアーカイブにする、またはvendor済みソースとして管理することです。　これについてより詳しく

### 対策案の比較

| 対策 | 内容 | 防げること | 残るリスク | 運用負荷 |
| --- | --- | --- | --- | --- |
| commit SHA固定 | `GIT_TAG` にタグ名ではなく40桁のcommit SHAを指定する。 | tag移動やbranch更新で別内容を取る問題を防ぐ。 | GitHubから取得する経路自体への依存は残る。 | 低い。更新時にSHAを差し替える。 |
| hash検証付きarchive | GitHubのarchive tar/zip等を `URL` で取得し、`URL_HASH SHA256=...` で検証する。 | ダウンロードされたarchiveが想定byte列と違えばビルドを失敗させる。 | GitHub archiveの再生成仕様やURL選定に注意が必要。 | 中程度。hash更新手順が必要。 |
| vendor済みソース | Spout2の必要ソースを `third_party/` 等へ取り込み、ビルド時に外部取得しない。 | Releaseビルド時の外部依存取得をなくせる。 | 取り込み時のレビュー、ライセンス、更新管理が必要。 | 高い。差分管理と更新作業が増える。 |

### このリポジトリでの現実的な優先順

1. まずは `GIT_TAG` をSpout2の実commit SHAへ固定する。
2. 可能ならCIで、現在使うSpout2 revisionや取得元をbuild metadataへ記録する。
3. より強い再現性が必要になったら、hash検証付きarchiveまたはvendor化へ進める。
4. 分離版zipを今後も配るなら、helper単体またはzip全体の署名も別issueで検討する。

CMakeのローカルドキュメントでも、外部サーバーから取得する場合はbranch/tag名よりcommit hashの利用が推奨されている。
また、URL downloadでは `URL_HASH` がdownload内容のintegrity確認として強く推奨されている。

### Spout2 revisionを更新するタイミング

`GIT_TAG` をcommit SHAへ固定した後は、通常の依存ライブラリ更新と同じ扱いにする。
Releaseごとに自動更新するものではなく、理由があるときだけ明示的なPR/issueで更新する。

- Spout2にセキュリティ修正、クラッシュ修正、Windows/DirectX/OpenGL/Spout互換性の修正が入ったとき。
- VRChat Stream Cameraや利用者環境で、現在固定しているSpout2由来と見られる不具合が確認されたとき。
- Visual Studio、Windows SDK、CMake、GitHub Actions runner更新などで、現在のSpout2 revisionがビルドできなくなったとき。
- Spout2側のライセンス、配布形態、依存DLL構成に変更があり、同梱物やライセンス表記を見直す必要があるとき。
- 定期棚卸しで、長期間固定したrevisionが古くなりすぎていると判断したとき。ただし自動追従ではなく、差分確認と実機確認を伴う。

更新しないケースも明確にする。
「新しいタグが出たから」「最新にしておきたいから」だけでは更新しない。
commit SHA固定は再現性と監査性を保つための措置なので、更新時は旧SHA、新SHA、対応するupstream tag、更新理由、ビルド結果、helper動作確認、ライセンス差分確認を記録する。

build metadataにSpout2 revisionを記録する目的は、Release成果物から「どのSpout2 commitでhelperを作ったか」を後から追跡できるようにすることである。
これは更新頻度を上げるためではなく、不具合調査、セキュリティ監査、Release再現確認で参照するための情報である。
