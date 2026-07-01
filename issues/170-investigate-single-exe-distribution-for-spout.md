# Spout同梱バイナリを単一exe配布に戻せるか調査する

## 問題

v0.1.8-rc13ではStream Camera(Spout)方式のために `ClipForVRChat.exe`、`spout-capture.exe`、`SpoutLibrary.dll` が分かれている。
v0.1.7までは利用者が主に `ClipForVRChat.exe` だけを改竄確認すればよかったが、複数バイナリになると確認対象が増え、PGP署名やハッシュ確認の手間が増える。

## 期待する挙動

ビルド済みのexe/dllを利用者から見て単一の `ClipForVRChat.exe` として配布できるかを調査し、実現方式、リスク、v0.1.8で採るべき方針を明確にする。
「1回のコマンドでビルドする」ことではなく、配布物と改竄確認対象を単一exeへ寄せることを主目的にする。

## 受け入れ条件

- [x] 現行の分離理由とRelease成果物構成を確認する。
- [x] `spout-capture.exe` と `SpoutLibrary.dll` を本体へ統合または埋め込みできる方式を比較する。
- [x] 方式ごとの改竄確認、ウイルス対策ソフト誤検知、ライセンス、クラッシュ影響、実機検証負荷を整理する。
- [x] v0.1.8で採るべき暫定方針と、将来対応する場合の実装チケット候補を提示する。

## 調査メモ

### 現行構成

- `v0.1.7` のRelease workflowは `ClipForVRChat.exe`、README、LICENSE、公開鍵URLをzipへ入れ、`ClipForVRChat.exe` に対する detached PGP署名 `.exe.asc` を添付していた。
- `v0.1.8-rc13` のRelease workflowは `Build Spout helper` で `spout-capture.exe` を作り、同じ出力フォルダの `SpoutLibrary.dll` も `dist/ClipForVRChat/` へコピーする。
- 現行の `.exe.asc` は `dist/ClipForVRChat/ClipForVRChat.exe` だけを対象にしている。zipのsha256はあるが、実行される `spout-capture.exe` と `SpoutLibrary.dll` をPGP署名で個別確認する導線はない。
- アプリ本体は `src/internal/appcore/spout.go` で `spout-capture.exe` を別プロセスとして呼び出す。これによりSpout/DirectX/OpenGL/DLLロードの失敗を本体プロセスから隔離している。
- Spout2 `2.007.015` の公式CMakeでは `SpoutLibrary` は `add_library(SpoutLibrary SHARED ...)` であり、公式コメントでも `SpoutLibrary.dll` を実行ファイルと同じフォルダに置く前提になっている。
- 一方でSpout2の下層には `Spout_static` があり、Spout SDKソースをC++アプリへ直接組み込む経路は存在する。

確認した一次ソース:

- `.github/workflows/release.yml`
- `scripts/build-windows-from-wsl.sh`
- `tools/spout-capture/CMakeLists.txt`
- `tools/spout-capture/main.cpp`
- Spout2 `2.007.015` の `SPOUTSDK/SpoutLibrary/CMakeLists.txt`
- Spout2 `2.007.015` の `SPOUTSDK/SpoutGL/CMakeLists.txt`
- Spout2 `2.007.015` の `SPOUTSDK/SpoutGL/readme.md`
- Spout2 `2.007.015` の `LICENSE`

### 方式比較

#### A. 現状維持し、署名だけ増やす

実装は最小。`spout-capture.exe.asc` と `SpoutLibrary.dll.asc`、またはzip全体のPGP署名をRelease assetへ追加すれば、改竄確認対象を明示できる。
ただし、利用者は複数ファイルを確認する必要があり、「exe1つに戻す」という要望は満たさない。

リスク:

- 改竄確認手順は増えたまま。
- Release asset数が増え、README/About/Release Notesの説明も増える。

#### B. `spout-capture.exe` に `SpoutLibrary.dll` を静的統合する

Spout2の `SpoutLibrary` は公式CMake上はDLL前提だが、下層の `Spout_static` を直接リンクしてhelperを作り直せば、配布物を `ClipForVRChat.exe` + `spout-capture.exe` の2バイナリへ減らせる可能性がある。
ただし本体exeとは別のままなので、改竄確認対象は最低2つ残る。

リスク:

- `tools/spout-capture/main.cpp` を `SpoutLibrary.h` のC互換DLL APIからSpoutGL/SpoutReceiver等のC++ SDK APIへ寄せる必要がある。
- Spout/DirectX/OpenGLまわりの実機回帰確認が必要。
- 目的である「単一exe配布」には届かない。

#### C. 本体 `ClipForVRChat.exe` にhelperとDLLを埋め込み、実行時に展開する

Release配布物は利用者から見て `ClipForVRChat.exe` 1つに戻せる。PGP署名も本体exe1つを確認すれば、埋め込まれた `spout-capture.exe` と `SpoutLibrary.dll` のバイト列まで署名対象に含められる。
実行時は本体が管理ディレクトリへhelper一式を展開し、既存どおり別プロセスで呼ぶ。プロセス分離の利点は維持できる。

実装案:

- Release workflowで `spout-capture.exe` と `SpoutLibrary.dll` を先にビルドする。
- その2ファイルと期待SHA-256をGoの埋め込み資産として本体ビルドへ渡す。
- 初回使用時または起動時に、`AppData` 配下などの管理ディレクトリへ `<version>/<sha256>/spout-capture.exe` と `SpoutLibrary.dll` を展開する。
- 展開前後にSHA-256を検証し、既存ファイルが一致すれば再利用、一致しなければ上書きまたは別ディレクトリへ展開する。
- 設定の `spoutHelperPath` が空または既定値の場合は埋め込みhelperを使い、ユーザーが明示パスを指定した場合だけ外部helperを使う。
- Release zipから `spout-capture.exe` と `SpoutLibrary.dll` を外し、`ClipForVRChat.exe` だけにPGP detached signatureを作る。

リスク:

- 実行時にexe/dllを展開する挙動は、環境によってウイルス対策ソフトの監視対象になりやすい。
- 展開先、上書き、削除、同時起動時の競合を設計する必要がある。
- `SpoutLibrary.dll` は通常のWindows DLLロードなので、helperと同じフォルダに正しいDLLを置く必要がある。
- 既存の「helperパスをユーザーが指定できる」設定との優先順位を整理する必要がある。

#### D. Spout処理をWails本体へ直接組み込む

理論上は本当の単一プロセス/単一exeにできるが、CGO/C++、DirectX/OpenGL、DLL/ライブラリリンク、Wailsプロセス安定性を同時に扱うことになる。
helper分離で避けていたクラッシュ影響が本体へ戻るため、v0.1.8のRC段階で採るには重い。

リスク:

- Spout/DirectX/OpenGLの失敗がアプリ全体クラッシュへ直結しやすい。
- Windows専用ビルドとGo/Wailsビルドの結合が強くなる。
- CI/Release、ローカルビルド、実機確認の範囲が大きく増える。

#### E. 汎用packerでexe/dllを1つにまとめる

技術的には自己展開型packerやDLL bundlerの系統で可能性はあるが、配布アプリとしては推奨しにくい。
ビルドの再現性、ライセンス、ウイルス対策ソフト誤検知、障害時の切り分け、署名対象の説明が難しくなる。

リスク:

- 誤検知リスクが高くなりやすい。
- CIでの検証と将来保守が外部packer依存になる。
- 何が実行時に展開/ロードされるかを利用者へ説明しづらい。

## 結論

v0.1.8で「利用者の改竄確認対象をv0.1.7相当に戻す」目的なら、C案の「本体exeにhelper/dllを埋め込み、実行時に検証付きで展開して呼ぶ」が最も現実的。
本体へSpoutを直接統合するD案より安全で、現行の別プロセス隔離を維持できる。

ただし、RC終盤でSpout実機確認範囲を増やす変更になるため、v0.1.8正式へ入れるかはリリース判断が必要。
すぐ正式化を優先するなら、A案としてzip全体または追加バイナリへのPGP署名を先に追加し、単一exe配布は次リリースでC案として実装するのが低リスク。

## 実装チケット候補

- 本体exeにSpout helper一式を埋め込み、検証付きで展開する。
- 埋め込みhelperの展開先、バージョン管理、同時起動時のロック、クリーンアップ方針を定義する。
- Release zipから `spout-capture.exe` と `SpoutLibrary.dll` を外し、zip内ファイル一覧と署名説明を単一exe前提へ戻す。
- 暫定対応を採る場合は、追加バイナリまたはzip全体のPGP署名をRelease workflowへ追加する。
