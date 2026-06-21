# Icon quality and multi-size generation

## Problem

The current app icon looks poor and does not match the desired soft image-plus-clipboard direction.

## Desired behavior

- Use a cleaner icon inspired by the provided reference image.
- Keep an original high-resolution icon source in the repository.
- Generate the Wails `build/appicon.png` from the source.
- Prepare generation logic so multiple icon sizes can be produced from the original when needed.

## Acceptance criteria

- The app no longer uses the Wails default icon.
- `build/appicon.png` is generated from a maintained source asset.
- The icon generation command can emit multiple standard sizes.

