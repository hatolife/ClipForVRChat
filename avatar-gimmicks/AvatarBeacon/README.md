# AvatarBeacon 元ファイル package

AvatarBeacon は、VRChat アバターから OSC Avatar Parameters へ、アバター基準の位置と向きの値を送るためのアバターギミック元ファイルです。

この source zip は、GitHub Releaseへ添付する公開Assetであり、Release確認と手動での Unity package 作成に使います。
CI は `AvatarBeacon-vX.Y.Z-source.zip` を作成しますが、`.unitypackage` は作成しません。

Unity package を作る場合は、この package を Unity のアバタープロジェクトへコピーまたは展開し、次の場所に asset があることを確認してください。

```text
Assets/PoppoWorks/AvatarBeacon/
```

その後、必要に応じて Unity から `Assets/PoppoWorks/AvatarBeacon` を `.unitypackage` として export します。
GitHub Releaseへ公開添付する標準Assetは source zip です。`.unitypackage` を手動作成する場合の作業用ファイル名は次の形式にします。

```text
AvatarBeacon-vX.Y.Z.unitypackage
```

手動作成した `.unitypackage` を確認する場合は、必要に応じて隣に SHA-256 ファイルを作成します。
Windows PowerShell では次のコマンドを使います。

```powershell
Get-FileHash "AvatarBeacon-vX.Y.Z.unitypackage" -Algorithm SHA256 |
  ForEach-Object { "$($_.Hash.ToLower())  AvatarBeacon-vX.Y.Z.unitypackage" } |
  Set-Content -Encoding ASCII "AvatarBeacon-vX.Y.Z.unitypackage.sha256"
```

Unity利用者向けの導入手順は次のファイルに含めています。

```text
Assets/PoppoWorks/AvatarBeacon/README.md
```

YL-ATG 由来部分の MIT License と notice は次のファイルに含めています。

```text
Assets/PoppoWorks/AvatarBeacon/LICENSES/YL-ATG-MIT.txt
Assets/PoppoWorks/AvatarBeacon/NOTICE.md
```
