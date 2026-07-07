# Auto-photo scan capによる旧画像漏えいを再検証する

## 指示

> 現在の作業をできるだけサブエージェントに任せて次の作業を実施 
> 重大な仕様変更・運用判断が必要　の項目について追加情報があります 
> '/home/user/work/ClipForVRChat/reports/security/2026-07-07T17-53-16.203Z/b825b8f570a881919546b22faf5da6cf' '/home/user/work/ClipForVRChat/reports/security/2026-07-07T17-53-16.203Z/394a0dc95c688191a1617b65fcb2befd' '/home/user/work/ClipForVRChat/reports/security/2026-07-07T17-53-16.203Z/e2dbfdd09f68819187b842db37a4421f' '/home/user/work/ClipForVRChat/reports/security/2026-07-07T17-53-16.203Z/958fd884b9588191ab138b574a6bd5eb' '/home/user/work/ClipForVRChat/reports/security/2026-07-07T17-53-16.203Z/69b57f2de0b081919021d7b078b7778c' '/home/user/work/ClipForVRChat/reports/security/2026-07-07T17-53-16.203Z/fa5bc68dae4c8191ad9140983e5ac66f' '/home/user/work/ClipForVRChat/reports/security/2026-07-07T17-53-16.203Z/f23b0a6320c481919bc02541a59d01da' 
> ７件について個別にissue作成 それぞれ対応方針を３パターンissueに記載

## 文脈

追加情報 `e2dbfdd09f68819187b842db37a4421f` は、Auto-photoのscan capが初期baselineにも適用されることで、古い画像の後日投稿や新規画像のstarvationが起きると指摘している。
19:48 CSV対応では、baseline/current scanを全件scanへ分離する修正を既に入れている。

## 解釈

この項目は追加情報に基づき、既存修正がfindingの要件を満たしているかを再検証する対象とする。
未解決なら追加修正し、解決済みならレポート上で対応済みへ移す。

## 問題

- 初期baselineがcap付きだと、cap外の既存画像が後から新規扱いになる。
- cap前方に大量ファイルがあると、後方の新規画像が処理されない。
- 自動投稿はDiscord uploadを強制するため、誤処理が画像漏えいになる。

## 期待する挙動

baselineと現在状態の把握ではcapにより既存ファイルを見落とさず、1 tickあたりの処理量だけを制限する。
古いファイルの後日投稿と新規ファイルのstarvationを再現テストで防ぐ。

## 対応方針案

- A: 既存修正どおりbaseline/currentは全件scanし、処理件数だけをtick単位で制限する。
- B: 全件scanに加えて未処理queueを導入し、大量新規ファイルも順次処理する。
- C: scan capを大幅に増やし、警告ログだけ追加して仕様上の制限として扱う。

## 方針評価

- A: findingの根本原因である不完全baselineを消せるため、今回の追加情報に対する正しい対応。
- B: 大量新規ファイルの公平性改善には有効だが、永続queue設計が入り別仕様になる。
- C: cap外ファイルが残る構造を温存するため、セキュリティ対応として不適切。

## 推奨方針

Aを採用する。
既に入れた全件baseline修正がfindingの根本条件を消しており、queue導入は別の仕様変更になる。
このissueでは追加情報との突合とテスト確認を完了条件にする。

## 受け入れ条件

- [x] 初期baselineがcap外の既存画像もseenに含めることを確認する。
- [x] 旧画像がcap変動で新規扱いされないテストが通る。
- [x] cap前方の大量ファイルで後方の新規画像が永久starveしないことを確認する。
- [x] 管理レポートに「既存修正で対応済み」または追加修正内容を具体的に追記する。

## 作業ログ

- 2026-07-08: `AutoPhotoWatcher.Run` と `tick` が、初期baseline/current判定に `scanPhotoFilesAllWithExcludesStatus` を使い、`MaxAutoPhotoScanFiles` で候補発見を打ち切らないことを確認した。
- 2026-07-08: `cd src && go test ./internal/appcore -run 'TestScanPhotoFilesAllWithExcludesStatusIgnoresScanCap|TestAutoPhotoWatcherFullScanBaselineDoesNotResurfaceDeletedOldFiles|TestAutoPhotoWatcherFullScanBaselineProcessesLateArrivingFile' -count=1` を実行し、成功した。
- 2026-07-08: 追加情報 `e2dbfdd09f68819187b842db37a4421f` の根本原因は `53b74d5 fix(autophoto): baseline full scans before tick limits` で解消済みと判断し、管理レポートのF43を保留から対応済みへ移した。
