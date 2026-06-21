# Clipboard screenshot corruption investigation

## Problem

Images captured with a small Win+Shift+S selection can become visually corrupted after processing from the clipboard.

Observed examples:

- Thin colored artifacts near the left edge.
- Incorrect-looking output for small screenshot clipboard input.

Cases that do not reproduce:

- D&D with an image that does not need shrinking.
- Copying an image from Paint and processing from clipboard.

## Investigation notes

The current code reads clipboard image bytes via `golang.design/x/clipboard` and decodes them as an encoded image. On Windows, clipboard image data may be DIB/BGRA rather than PNG, especially for screenshot tools. If DIB bytes are treated as regular encoded image data, decoding or pixel layout may break.

## Desired behavior

- Clipboard image input should preserve the actual screenshot pixels.
- Small screenshots that do not need resizing should not change visually except for PNG encoding.

## Acceptance criteria

- Windows clipboard image reading handles screenshot clipboard data correctly.
- Add focused tests or a small fixture-based regression test if practical.

