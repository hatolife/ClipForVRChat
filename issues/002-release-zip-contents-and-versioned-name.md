# Release zip contents and versioned filename

## Problem

The Windows release zip should include user-facing files and its filename should include the version.

## Desired behavior

- Windows release zip contains:
  - `ClipForVRChat.exe`
  - `README.md`
  - `LICENSE`
- Release zip filename includes the version tag.

## Acceptance criteria

- A release for tag `vX.Y.Z` uploads `ClipForVRChat-vX.Y.Z-windows-amd64.zip`.
- The zip includes the executable, README, and license.

