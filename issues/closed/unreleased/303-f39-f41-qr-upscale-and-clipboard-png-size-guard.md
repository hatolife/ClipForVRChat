# F39/F41 即時修正

## 指示

> fix immediate no-spec findings F39 and F41 from reports/security/2026-07-07T17-53-16.203Z/codex-security-findings-remediation-report.md.
> - QR detection upscaling must not create variants whose total pixel count exceeds a safe bounded limit; avoid crafted tall image OOM.
> - Native Windows clipboard PNG path must reject too-large GlobalSize before unsafe.Slice/make/copy. Use existing input byte limits where appropriate.
> Add focused tests where possible (non-Windows tests for QR; Windows-specific tests can be helper-unit tests if feasible). Run relevant go tests from src if feasible. Edit files directly and report changed files/tests.

## 文脈

`reports/security/2026-07-07T17-53-16.203Z/codex-security-findings-remediation-report.md` の即時修正対象として、F39 と F41 をまとめて対処する。

## 解釈

QR 検出の拡大画像生成に総ピクセル数の上限を設け、細長い画像でも安全にスキップできるようにする。
Windows のネイティブ clipboard PNG 読み出しでは、`GlobalSize` が大きすぎる場合に `unsafe.Slice` や `make` の前で拒否する。

## 問題

- QR 検出の拡大候補が細長い画像で巨大化し、メモリ消費が跳ね上がる。
- Windows clipboard の PNG 読み出しがサイズ確認前に生バッファ化を試みると、クラッシュや OOM の原因になる。

## 期待する挙動

- QR 検出は、安全上限を超える拡大候補を生成しない。
- Windows clipboard PNG 読み出しは、`GlobalSize` が上限超過なら即座に失敗する。
- それぞれに、再発防止用の focused test が追加される。

## 受け入れ条件

- [ ] QR 検出の拡大候補が安全上限を超える場合に生成されない。
- [ ] Windows clipboard PNG 読み出しが `GlobalSize` 超過を `unsafe.Slice` 前に拒否する。
- [ ] QR の focused test が追加される。
- [ ] clipboard の focused test または helper test が追加される。
- [ ] 関連する `go test` が実行され、結果が記録される。
