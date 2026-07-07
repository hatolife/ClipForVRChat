# 非表示のカメラ自動起動/終了設定を無効化する

## 指示

> 現在の作業をできるだけサブエージェントに任せて次の作業を実施 
> 重大な仕様変更・運用判断が必要　の項目について追加情報があります 
> '/home/user/work/ClipForVRChat/reports/security/2026-07-07T17-53-16.203Z/b825b8f570a881919546b22faf5da6cf' '/home/user/work/ClipForVRChat/reports/security/2026-07-07T17-53-16.203Z/394a0dc95c688191a1617b65fcb2befd' '/home/user/work/ClipForVRChat/reports/security/2026-07-07T17-53-16.203Z/e2dbfdd09f68819187b842db37a4421f' '/home/user/work/ClipForVRChat/reports/security/2026-07-07T17-53-16.203Z/958fd884b9588191ab138b574a6bd5eb' '/home/user/work/ClipForVRChat/reports/security/2026-07-07T17-53-16.203Z/69b57f2de0b081919021d7b078b7778c' '/home/user/work/ClipForVRChat/reports/security/2026-07-07T17-53-16.203Z/fa5bc68dae4c8191ad9140983e5ac66f' '/home/user/work/ClipForVRChat/reports/security/2026-07-07T17-53-16.203Z/f23b0a6320c481919bc02541a59d01da' 
> ７件について個別にissue作成 それぞれ対応方針を３パターンissueに記載

## 文脈

追加情報 `394a0dc95c688191a1617b65fcb2befd` は、自動撮影の `openCameraBeforeBatch` / `closeCameraAfterBatch` がUIから隠れた後も、既存configやインポートconfigのtrue値をバックエンドが実行し続けると指摘している。

## 解釈

通常UIで確認・変更できないカメラ自動起動/終了は、少なくとも現行UIでは有効設定として扱うべきではない。
復活させる場合は、ユーザーが見える場所で明示設定できる状態に戻してから扱う。

## 問題

- UIが非表示にしたフラグを、フロントエンドは未定義時だけfalseにして既存trueを残す。
- バックエンドはtrue値を尊重し、VRChat camera/stream/close OSCを送信する。
- スケジュールやDiscord投稿と組み合わさると、意図しない撮影・投稿につながる。

## 期待する挙動

UIで操作できないカメラ自動起動/終了設定は保存・実行時に無効化され、既存またはインポート設定のtrue値でOSC camera操作が走らない。

## 対応方針案

- A: config正規化で `openCameraBeforeBatch` / `closeCameraAfterBatch` を常にfalseへ丸める。
- B: UIへ設定項目を復活させ、保存前確認にも含めてユーザーが明示的にONにできるようにする。
- C: 既存true値を読み込んだ場合だけ警告を出して、保存時にfalseへ移行する。

## 方針評価

- A: 現行UIで操作できない値を実行させないため、短期修正として最も安全。
- B: 機能を維持できるが、Camera OSCの挙動説明、保存前確認、実機確認が必要で今回の即時修正より大きい。
- C: 保存前にスケジューラやテスト撮影が走る経路を残し得るため、単独では不十分。

## 推奨方針

Aを採用する。
該当機能は現行UIで隠されており、Camera OSCの安定性やプライバシー影響も大きいため、まず実行面で無効化する。
将来UIを復活させる場合は別issueで再設計する。

## 方針決定

2026-07-08: ユーザー判断によりAを採用する。
UIで操作できない `openCameraBeforeBatch` / `closeCameraAfterBatch` は、既存configやインポートconfigにtrueが残っていても有効化しない。
実装時はフロントエンドだけに依存せず、バックエンド正規化または実行前正規化でfalseへ丸める。

## 受け入れ条件

- [ ] 読み込み・保存・実行のいずれでも非表示カメラ自動起動/終了trueが有効化されない。
- [ ] インポートconfigにtrueが含まれていても保存後はfalseになる。
- [ ] フロントエンドだけでなくバックエンド正規化または実行前正規化で強制される。
- [ ] 自動撮影の通常動作、Stream取得、復元処理を壊さない。
- [ ] 該当挙動を検証するGoテストを追加する。
