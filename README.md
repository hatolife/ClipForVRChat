# ClipForVRChat

ClipForVRChat is a Windows app for preparing images for VRChat.

It resizes images so they fit within VRChat's 2048x2048 limit, can upload them to a configured Discord webhook, and copies the resulting image URL to the clipboard.

## What It Does

- Resizes images to fit within 2048x2048 without changing the aspect ratio
- Keeps JPEG images as JPEG, and saves other image types as PNG
- Saves resized images locally when enabled
- Uploads resized images to Discord when enabled
- Copies the Discord image URL to the clipboard for single-image uploads
- Shows a URL list with thumbnails for multiple images

## Basic Usage

Run without arguments to process the image currently in the clipboard:

```txt
ClipForVRChat.exe
```

Run with one or more image files:

```txt
ClipForVRChat.exe image.png
ClipForVRChat.exe image1.png image2.jpg
```

Open the settings UI by passing the config file:

```txt
ClipForVRChat.exe config.json
```

## Settings

Settings are stored in `config.json` next to the exe.

The app is designed so normal users can edit settings from the UI instead of editing JSON by hand.

Important settings:

- Local save on/off
- Discord upload on/off
- Discord webhook URL
- Output directory
- JPEG quality
- UI display mode

## Discord Setup

Create a Discord webhook for the channel you want to use, then set that webhook URL in the settings UI.

Uploaded image URLs are Discord attachment direct links.

