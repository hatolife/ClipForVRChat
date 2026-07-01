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

## 2026-07-01 追加調査スコープ

ユーザー方針としてC案を採用する。ただし将来の検証や障害切り分けが面倒になる可能性を考慮し、Releaseでは単一exe版を主導線にしつつ、分離バイナリ版もAssetsへ残す構成を調査する。

追加で、現在Goで実装している本体をC++/Rustなど別言語で書き直した場合に、Spout統合や単一exe配布の問題が解消するかも確認する。

また、VRChat Stream CameraをOBS等がキャプチャできる仕組みが現在のSpout方式と同じか、別方式ならその方式を採用することでバイナリ分離問題を解消できるかを調査する。

追加受け入れ条件:

- [x] C案採用時に、単一exe版と分離版を同時にCI/Release生成する構成案を整理する。
- [x] C++/Rust等への書き直しがSpout統合や単一exe化へ与える影響を整理する。
- [x] OBS等がVRChat Stream Cameraを受ける仕組みと、ClipForVRChatで採用可能な代替方式を確認する。
- [x] 追加調査の結論をこのissueへ追記する。

## 2026-07-01 追加調査結果

### C案採用時のRelease構成

C案を採用し、通常利用者向けの主導線は単一exe版にする。
ただし、将来の検証や障害切り分けのため、分離版も同じcommitからRelease assetとして生成し続ける。

推奨するRelease asset:

- `ClipForVRChat-vX.Y.Z-windows-amd64.exe`
  - 主導線。`spout-capture.exe` と `SpoutLibrary.dll` を埋め込んだ単一exe。
  - Release本文の「プログラムのダウンロード」はこのファイルへ向ける。
- `ClipForVRChat-vX.Y.Z-windows-amd64.exe.asc`
  - 主導線exeのPGP detached signature。
  - 利用者の改竄確認対象を原則この1ファイルに戻す。
- `ClipForVRChat-vX.Y.Z-windows-amd64.exe.sha256`
  - 主導線exeの破損確認用。
- `ClipForVRChat-vX.Y.Z-windows-amd64-separated.zip`
  - 検証・切り分け用。`ClipForVRChat.exe`、`spout-capture.exe`、`SpoutLibrary.dll`、README、LICENSE、`Spout2-LICENSE.txt` を含める。
  - Release本文では通常利用者向けではなく、helper単体確認や不具合切り分け用として説明する。
- `ClipForVRChat-vX.Y.Z-windows-amd64-separated.zip.sha256`
  - 分離版zipの破損確認用。
- `ClipForVRChat-vX.Y.Z-build-metadata.json`
  - 単一exe版と分離版の元commit、helper SHA-256、Spout2 version、Wails/Go/Node versionを記録する。

分離版の署名は必須ではないが、検証用として長く残すなら `separated.zip.asc` を追加してもよい。
ただし通常利用者の案内は単一exeの `.exe.asc` に寄せ、確認対象を増やさない。

実装方針:

- Release workflowでは、現行どおり先に `spout-capture.exe` と `SpoutLibrary.dll` をビルドする。
- 2ファイルのSHA-256を計算し、Goの `embed` 対象ディレクトリへコピーしてからWails本体をビルドする。
- 本体は既定設定では埋め込みhelperを `%APPDATA%` またはアプリ管理データ配下の `spout-helper/<app-version>/<helper-sha256>/` のようなディレクトリへ展開する。
- 展開前後にSHA-256を検証し、一致済みなら再利用する。一致しない既存ファイルは上書きではなく別hashディレクトリへ展開する。
- ユーザーが明示的に `spoutHelperPath` を指定した場合は外部helperを優先できるようにする。これにより、分離版zipや手元ビルドのhelperを使った切り分けを残せる。
- CIでは単一exe版の「外部helperなしで展開・`--version`・`--list-senders` 相当が通る」確認と、分離版zip内の `spout-capture.exe --help/--version/--list-senders` 確認を両方行う。

注意点:

- 直接 `.exe` をRelease assetにすると、README/LICENSE/Spout2 licenseを同梱するzipがなくなる。ライセンス表示はRelease本文、README、アプリ内OSS表示、または別添付のlicense資料で維持する。
- 実行時にexe/dllを展開する挙動はAV/EDRに監視されやすい。展開先をユーザーデータ配下に固定し、hash検証、診断ログ、失敗時の案内を明確にする。
- 「分離版」は通常導線に出さず、検証・切り分け用であることをRelease本文で明確にする。

### C++/Rust等への書き直し

- `src/internal/appcore/spout.go` は helper を別プロセスとして呼ぶ前提で、Wails 本体から Spout/DirectX/OpenGL の失敗を隔離している。
- `tools/spout-capture/main.cpp` は Spout sender の列挙と受信、PNG 保存だけに責務を絞っているため、rewrite の主目的は Spout の DLL/CRT/GUI ではなく「配布物の境界」をどこに置くかになる。
- Wails は公式に Windows で外部DLL不要、かつ最終成果物をassets込みの単一exeへまとめられるとしているが、これはWailsアプリ資産の話であり、Spout helper の別プロセス化は別問題。
- MSVC の `/MT` は静的CRTを使う設定なので、C++ helper のCRT依存を減らすことはできる。ただし `SpoutLibrary.dll` の問題は残るため、DLLを消すには `SpoutLibrary` ではなく `Spout_static` 経由でhelperを再実装する必要がある。
- Rust でも単体 exe は作れるが、Spout の受信・DirectX/OpenGL・WIC/COM 連携をどう組むかは結局別の unsafe FFI/ラッパー設計になるため、言語変更だけでは配布物の分割問題は解消しない。
- 現行のリリース構成では、`ClipForVRChat.exe` に helper/dll を埋め込んで実行時展開する案が、利用者から見た単一 exe と crash isolation を両立しやすい。
- したがって「Go/Wails 本体の C++/Rust フルリライト」は、この件に対する費用対効果が低い。やるなら helper のみを C++ で静的寄せする方が現実的だが、それでも最終配布の単一 exe 化は本体側の埋め込み/展開が必要。

### まとめ

1. C++ full rewrite: 単一 exe は作りやすくなる可能性はあるが、Spout 統合と GUI/配布基盤の作り直しが重すぎる。
2. Rust full rewrite: こちらも単一 exe は可能だが、Spout と Windows グラフィックス周りの FFI が残り、移行コストに見合いにくい。
3. C++ helper rewrite/static-link: `SpoutLibrary.dll` を減らすのには効くが、全体を単一 exe にする決定打ではない。
4. Rust helper rewrite/FFI/library options: 可能だが、Spout の実体は C++/Windows API 側に寄るので、実装複雑性は下がりにくい。
5. Recommendation: 本体は Go/Wails のまま、helper/dll を埋め込んで実行時展開する C 案を採る。helper だけを静的化するのは副次的改善として検討する。

### OBSとVRChat Stream Cameraの仕組み

VRChatのStream Camera本体はSpout方式と考えてよい。
VRChat WikiのUser Camera OSC表では `/usercamera/Mode` の `2` がStream、`/usercamera/Streaming` がSpout streamとして定義されている。

OBSでVRChatを映せる経路は大きく2種類ある。

1. OBS標準のWindow/Game/Display Capture
   - これは画面やウィンドウを取る経路で、VRChat Stream CameraのSpout senderを直接受信するものではない。
   - VRChatのゲーム画面やウィンドウを撮ることはできるが、Stream Camera映像そのものを静止画として受ける本命要件とは別物。
2. OBS Spout2 plugin等のSpout受信Source
   - これはSpout2互換プログラムからshared textureをimport/exportするpluginで、ClipForVRChatのhelperと同じくSpout受信系の仕組み。
   - plugin実装も `SpoutLibrary` をリンクし、Spout DLL類をplugin配下へコピーする構成になっている。

したがって、OBSでStream Cameraを直接受けられる場合も、実体はSpout受信であり、ClipForVRChatが既に採用している方向と同じ。
OBS標準の画面キャプチャへ戻せば追加DLL問題は避けやすいが、それは「Stream Camera映像そのもの」ではなく、v0.1.8-rc12以前の白画像・画面周辺キャプチャ問題へ戻る。
OBS Spout plugin方式を採用してもSpout native依存は残るため、exe/dll分割問題の根本解決にはならない。

### 追加調査で確認した外部一次情報

- VRChat Wiki OSC User Camera: https://wiki.vrchat.com/wiki/OSC#User_Camera
- VRChat OSC Overview: https://docs.vrchat.com/docs/osc-overview
- Wails Introduction: https://wails.io/docs/introduction/
- MSVC `/MD`, `/MT`, `/LD`: https://learn.microsoft.com/en-us/cpp/build/reference/md-mt-ld-use-run-time-library
- Spout2: https://github.com/leadedge/Spout2
- OBS win-capture source tree: https://github.com/obsproject/obs-studio/tree/master/plugins/win-capture
- OBS Spout2 plugin: https://github.com/Off-World-Live/obs-spout2-plugin

## 追加調査後の決定

- C案を採用する。
- 通常利用者向けRelease本文の主リンクは単一exe版にする。
- 分離版zipは同じRelease Assetsへ残し、将来の実機検証、helper単体検証、障害切り分けに使えるようにする。
- 本体のC++/Rust全面書き直しは、この問題の解決策として採らない。
- OBS標準キャプチャへの回帰は採らない。Stream Camera映像そのものを扱うため、Spout受信経路を維持する。
