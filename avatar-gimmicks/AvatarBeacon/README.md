# AvatarBeacon source package

AvatarBeacon is a VRChat avatar gimmick source package for sending avatar-based
position and forward-vector values through OSC Avatar Parameters.

This source zip is intended for release verification and manual Unity packaging.
CI does not create a `.unitypackage`. To make one, copy or extract this package
into a Unity avatar project, confirm the assets under:

```text
Assets/PoppoWorks/AvatarBeacon/
```

Then export `Assets/PoppoWorks/AvatarBeacon` from Unity as a `.unitypackage`.
Use this release file name pattern:

```text
AvatarBeacon-vX.Y.Z.unitypackage
```

For release verification, create a SHA-256 file next to the `.unitypackage`.
On Windows PowerShell:

```powershell
Get-FileHash "AvatarBeacon-vX.Y.Z.unitypackage" -Algorithm SHA256 |
  ForEach-Object { "$($_.Hash.ToLower())  AvatarBeacon-vX.Y.Z.unitypackage" } |
  Set-Content -Encoding ASCII "AvatarBeacon-vX.Y.Z.unitypackage.sha256"
```

The Unity-facing setup guide is included at:

```text
Assets/PoppoWorks/AvatarBeacon/README.md
```

YL-ATG derived parts are covered by the MIT license and notice included in:

```text
Assets/PoppoWorks/AvatarBeacon/LICENSES/YL-ATG-MIT.txt
Assets/PoppoWorks/AvatarBeacon/NOTICE.md
```
