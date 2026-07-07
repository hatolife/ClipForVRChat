# カメラ移動時の飛行モードON時間を短くする

## 指示

> カメラを移動させた後飛行モードをOFFにするよ変更したい 飛行モードonの時間を極力短時間にするようにしたい 実装完了したらrc1を作成

## 文脈

自動撮影と「この構図へカメラ移動」では `/usercamera/Pose` でUser Camera位置を送信する。User Cameraの飛行モードがONのまま残ると、撮影後や手動移動後の操作に影響する可能性がある。

## 解釈

Pose送信に必要な最短範囲だけ `/usercamera/Flying` をONにし、Pose送信直後にOFFへ戻す。自動撮影中の各構図移動と手動の構図移動の両方を対象にする。実装と検証が完了したら `v0.1.8-rc1` を作成する。

## 問題

- カメラ移動後にUser Cameraの飛行モードがONのまま残る可能性がある。
- 飛行モードONの時間が長いと、ユーザー操作やVRChat側のカメラ状態へ影響しやすい。

## 期待する挙動

- カメラPose送信直前にだけ飛行モードをONにする。
- Pose送信が成功または失敗した後、すぐ飛行モードをOFFにする。
- 自動撮影と手動カメラ移動のどちらでも同じ制御になる。
- 実装後に `v0.1.8-rc1` のGitHub Releaseが作成される。

## 受け入れ条件

- [x] `/usercamera/Pose` 送信直前に `/usercamera/Flying=true` を送信する。
- [x] `/usercamera/Pose` 送信直後に `/usercamera/Flying=false` を送信する。
- [x] Pose送信失敗時も可能な限り `/usercamera/Flying=false` を送信する。
- [x] 既存の自動撮影テストが通る。
- [x] `v0.1.8-rc1` タグを作成し、Release workflowと配布Assetを確認する。

## 結果

- `applyCameraView` でPose送信直前に `/usercamera/Flying=true`、Pose送信直後に `/usercamera/Flying=false` を送るようにした。
- Pose送信が失敗した場合も `/usercamera/Flying=false` を送信するテストを追加した。
- `v0.1.8-rc1` のRelease workflowが成功し、GitHub Releaseがprereleaseとして公開された。
- 通常zipと分離zipの内容を確認した。
