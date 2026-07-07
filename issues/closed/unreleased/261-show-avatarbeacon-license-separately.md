# AvatarBeaconライセンスをOSS画面の最下部に分けて表示する

## 指示

> avatarbeaconのも最下部に分けて書きたい

## 文脈

#260でアプリ本体のOSSライセンス表示を実依存に合わせ、全項目の本文表示に対応した。AvatarBeacon側には `NOTICE.md` と `LICENSES/YL-ATG-MIT.txt` があるが、アプリ内OSS画面では通常のアプリ本体依存と分けて表示されていない。

## 解釈

AvatarBeaconの由来ライセンス表記を、通常のアプリ本体OSS一覧とは別のセクションとしてOSSライセンス画面の最下部に表示する。

## 問題

AvatarBeacon関連のライセンス表記をアプリ本体依存と同じ一覧に混ぜると、どの配布物に関係する表記か分かりにくい。

## 期待する挙動

OSSライセンス画面の下部にAvatarBeacon用の別枠があり、YL-ATG由来であることとMITライセンス本文を確認できる。

## 受け入れ条件

- AvatarBeacon関連ライセンスが通常OSS一覧とは別セクションで最下部に表示される。
- AvatarBeacon関連ライセンスに本文が表示される。
- 通常OSS一覧の重複、空欄、本文ありの検査が維持される。

## 対応内容

- `OSSLicense` に表示グループを追加し、AvatarBeacon関連ライセンスだけ `avatarBeacon` グループにした。
- OSSライセンス画面で通常OSS一覧とAvatarBeacon一覧を分け、AvatarBeaconを最下部に表示するようにした。
- AvatarBeacon / YL-ATG の由来説明とMITライセンス本文を表示できるようにした。
- AvatarBeaconライセンスが専用グループに入ることをテストで確認するようにした。
