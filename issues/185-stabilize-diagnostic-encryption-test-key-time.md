# 診断zip暗号化テストの一時鍵時刻を安定化する

## 問題

`go test ./...` で `TestDiagnosticZipDoesNotIncludeOutputImages` が、テスト内で生成したOpenPGP一時鍵を暗号化可能な鍵として扱えず失敗する場合がある。

鍵の作成時刻と暗号化時刻がほぼ同時で、環境や時刻解決によって有効開始前の鍵として扱われることがある。

## 期待する挙動

テスト用OpenPGP鍵は暗号化時刻より十分過去の作成時刻で生成し、ローカル環境でもCIでも安定して暗号化/復号テストが通る。

## 受け入れ条件

- [x] `TestDiagnosticZipDoesNotIncludeOutputImages` が時刻揺らぎで失敗しない。
- [x] `go test ./...` が通る。

## 対応メモ

- 2026-07-04: テスト用OpenPGP entityの作成時刻を暗号化時刻より十分過去に固定した。
