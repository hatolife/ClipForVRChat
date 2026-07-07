# Windows C++のmin/max macro衝突対策を再発防止として記録する

## 指示

> stdのmin,maxみたいなあるあるのやつの対応は(std::min)()のようにするのが定番だと思う これ再発防止でどこか書いておいて

## 文脈

`v0.1.8-rc3` のRelease workflowで、Windows header由来の `max` macroが `std::numeric_limits<...>::max()` に干渉し、Spout helperのMSVCビルドが失敗した。
暫定修正として `windows.h` より前に `NOMINMAX` を定義してRelease workflowは成功したが、今後Windows向けC++コードで同種のmacro衝突を避ける運用を明文化する必要がある。

## 解釈

今後の実装修正・レビューで参照できるよう、プロジェクト作業ルールにWindows C++の `min` / `max` macro衝突対策を追加する。
標準ライブラリの `std::min` / `std::max` / `std::numeric_limits<T>::max()` は、必要に応じて `(std::min)(...)` / `(std::max)(...)` / `(std::numeric_limits<T>::max)()` のようにparenthesized callで書く方針を明記する。

## 問題

- `windows.h` や一部依存headerが `min` / `max` macroを定義すると、`std::max(...)` や `std::numeric_limits<T>::max()` がMSVCで壊れる。
- Linuxローカルテストでは再現しにくく、Release workflowのWindows buildで初めて失敗し得る。
- `NOMINMAX` だけに依存すると、include順や依存headerによって再発する可能性が残る。

## 期待する挙動

Windows向けC++コードでは、標準ライブラリの `min` / `max` 名とWindows macroが衝突し得ることを前提に実装・レビューする。
再発防止ルールがAGENTS.mdに残り、今後の作業者が確認できる。

## 受け入れ条件

- [x] AGENTS.mdにWindows C++の `min` / `max` macro衝突対策が記載される。
- [x] `(std::min)(...)` / `(std::max)(...)` / `(std::numeric_limits<T>::max)()` の例が記載される。
- [x] `NOMINMAX` は有効な対策だが、include順や依存headerに注意が必要なことが分かる。

## 作業メモ

2026-07-08: Windows向けC++でMSVC build時の `min` / `max` macro衝突を防ぐために、AGENTS.mdの作業ルールを「Windows header include時の注意なし」から「`std::min` / `std::max` / `std::numeric_limits<T>::max()` は必要に応じてparenthesized callを使い、`NOMINMAX` はinclude順に注意してMSVC/Windows CIで確認する」方針へ変更した。
