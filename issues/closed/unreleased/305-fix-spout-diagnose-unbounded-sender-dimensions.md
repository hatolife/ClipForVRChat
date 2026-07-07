# Spout diagnoseのsender寸法上限とoverflow対策を入れる

## 指示

> Task: fix immediate no-spec finding F10 from reports/security/2026-07-07T17-53-16.203Z/codex-security-findings-remediation-report.md.
> - Validate Spout sender width/height before allocating width*height*4 buffers in capture and diagnose paths.
> - Add clear max dimension/pixel/byte constants and overflow-safe multiplication.
> - Return a clear diagnostic error instead of throwing/overflowing.
> Build/check the helper if feasible in this environment; otherwise report why not. Edit files directly and report changed files/tests.

## 文脈

`reports/security/2026-07-07T17-53-16.203Z/codex-security-findings-remediation-report.md` の F10 は、Spout sender が返す幅/高さを無制限に信頼して `width * height * 4` バッファを確保している点を指摘している。

## 解釈

このissueでは、`tools/spout-capture/main.cpp` の capture / diagnose の両方で sender 寸法を検証し、上限超過や乗算オーバーフローを明示的な診断エラーとして返す。

## 問題

- sender が異常な幅/高さを返すと、巨大確保や乗算オーバーフローにつながる。
- diagnose 経路でも同じ `width * height * 4` 確保があり、異常値をそのまま使うと不安定になる。
- エラーが曖昧だと、切り分けや再実行が難しい。

## 期待する挙動

- Spout sender の width/height を、確保前に上限付きで検証する。
- ピクセル数と byte 数に明確な定数上限を設け、overflow-safe に計算する。
- 上限超過や異常値のときは、throw や未定義動作ではなく、読みやすい diagnostic error を返す。

## 受け入れ条件

- [x] capture / diagnose 両方で sender 寸法の検証が入る。
- [x] max dimension / pixel / byte の定数が明示される。
- [x] 乗算は overflow-safe に行われる。
- [x] 異常寸法のときは明確な diagnostic error が返る。
- [x] 可能なら helper の build / check を実施し、結果を記録する。
