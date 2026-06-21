# ClipForVRChat Specification

## Overview

ClipForVRChat is a Windows desktop application for VRChat users who need to resize external images to fit VRChat's image resolution limits, optionally upload them to Discord, and quickly copy usable image URLs.

The application is implemented in Go with Wails for the desktop UI.

## Goals

- Resize input images to fit within 2048x2048.
- Preserve aspect ratio.
- Avoid upscaling images that are already within the limit.
- Save resized images locally when configured.
- Upload resized images to a configured Discord webhook when configured.
- Copy a usable image URL or image data to the clipboard when possible.
- Provide a simple URL list UI for multi-image workflows.
- Keep settings in JSON next to the executable.

## Non-Goals

- Persistent upload history across app launches.
- Manual JSON editing as the primary settings workflow.
- Supporting arbitrary upload providers in the initial version.
- Long-term hosting guarantees beyond Discord attachment behavior.

## Application Name

`ClipForVRChat`

Executable name:

```txt
ClipForVRChat.exe
```

## Platform

Initial target:

- Windows

Technology:

- Go
- Wails

## Input

### No Arguments

When launched without arguments, the app reads the current clipboard image and processes it.

If the clipboard does not contain an image, the app shows an actionable error explaining that the user should copy an image to the clipboard or pass image files as arguments.

### Image File Arguments

When launched with one or more image file paths, each image is processed.

Supported initial input formats should include:

- JPEG
- PNG
- WebP if the selected Go image stack supports it cleanly

Unsupported files should produce actionable errors.

### Config File Argument

When launched with exactly one `.json` config file path, the app opens the settings UI for that config.

If image files and config files are mixed, the app should show an error explaining that settings editing and image processing should be launched separately.

## Image Processing

### Resize Rule

Images are resized so both width and height fit within:

```txt
2048x2048
```

The aspect ratio is preserved.

Images already within the limit are not enlarged.

### Output Format

- JPEG input is output as JPEG.
- Any non-JPEG input is output as PNG.
- Clipboard image input is output as PNG unless reliable source format detection is available.

### JPEG Quality

JPEG quality is configurable.

Default:

```txt
92
```

## Local Output

Local file output can be enabled or disabled.

When enabled, resized images are written to disk.

### Output Directory

The output directory is configurable.

If no output directory is configured:

- File input: output next to the input image.
- Clipboard input: output to an `output` directory next to the executable.

### Output Name

Default suffix:

```txt
_2048
```

Examples:

```txt
image.png -> image_2048.png
image.jpg -> image_2048.jpg
```

By default, existing files are not overwritten.

If the target path exists, the app should generate a numbered filename:

```txt
image_2048.png
image_2048_2.png
image_2048_3.png
```

## Discord Upload

Discord upload uses a Discord webhook.

The app uploads the resized image file or in-memory resized image data to the configured webhook.

After upload, the app extracts the Discord attachment direct URL from the webhook response.

If Discord upload is enabled but the webhook URL is missing or invalid, the app shows an actionable error telling the user to open settings and configure the Discord webhook URL.

## Clipboard Output

### Single Image

If exactly one image is processed and Discord upload succeeds, the uploaded image URL is copied to the clipboard automatically when enabled.

If local save and Discord upload are both disabled, a single processed image is copied to the clipboard as image data.

### Multiple Images

If multiple images are processed and Discord upload succeeds, the URL list UI is shown.

Clicking a list item copies that item's URL to the clipboard.

If local save and Discord upload are both disabled and multiple images are provided, the app shows an error explaining that it can only copy one processed image to the clipboard, and the user should enable local save, enable Discord upload, or process one image at a time.

## UI

The UI is built with Wails.

### URL List

The URL list is shown when:

- Multiple uploaded image URLs are available.
- The UI mode is set to `always`.
- An error needs to be shown in an actionable way.

Each URL list item includes:

- Thumbnail
- Source file name or clipboard label
- URL
- Click-to-copy behavior

After copying, the UI shows a short toast-like message:

```txt
コピーしました
```

### Settings UI

The settings UI is opened by passing the config file path:

```txt
ClipForVRChat.exe config.json
```

The settings UI should allow editing:

- Local save enabled
- Discord upload enabled
- Discord webhook URL
- Output directory
- Output suffix
- JPEG quality
- UI display mode
- Single-image automatic URL copy

## Settings

Settings are stored in JSON.

Default location:

```txt
config.json
```

The file is stored next to the executable.

Users are expected to use the settings UI. The JSON should still be readable and stable for debugging.

### Draft Schema

```json
{
  "image": {
    "maxWidth": 2048,
    "maxHeight": 2048,
    "suffix": "_2048",
    "overwrite": false,
    "jpegQuality": 92,
    "outputDirectory": ""
  },
  "output": {
    "saveLocal": true,
    "uploadDiscord": true,
    "showUi": "auto",
    "copySingleUrlToClipboard": true
  },
  "discord": {
    "webhookUrl": ""
  }
}
```

### UI Display Mode

Allowed values:

- `auto`
- `always`
- `never`

`auto` behavior:

- Single successful URL upload: do not show UI unless needed.
- Multiple successful URL uploads: show URL list UI.
- Error: show UI or a dialog with actionable guidance.

## Error Messages

Errors should explain what happened and what the user should do next.

Examples:

- Clipboard has no image: "クリップボードに画像がありません。画像をコピーしてから再実行するか、画像ファイルを exe にドラッグしてください。"
- Multiple image clipboard-only output is impossible: "複数画像はクリップボードへ直接保持できません。ローカル保存またはDiscord投稿をONにするか、1枚ずつ処理してください。"
- Missing webhook URL: "Discord投稿がONですがWebhook URLが未設定です。設定画面でWebhook URLを設定してください。"

## Open Questions

- Whether WebP input should be supported in the first implementation.
- Whether animated image inputs should be rejected or flattened to the first frame.
- Whether clipboard image output should preserve alpha when copied back to the clipboard.
- Whether the first release should include drag-and-drop onto the running UI in addition to file arguments.

