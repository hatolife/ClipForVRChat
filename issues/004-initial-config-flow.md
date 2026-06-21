# Initial config flow

## Problem

When no settings file exists, the app should not silently continue with defaults.

## Desired behavior

- If `config.json` is missing on normal launch, show the settings UI first.
- After the user saves settings, continue the original processing flow.

## Acceptance criteria

- First launch without args opens settings, then processes clipboard after save.
- First launch with image args opens settings, then processes those images after save.

