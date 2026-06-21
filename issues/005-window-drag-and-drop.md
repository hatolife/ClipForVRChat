# Window drag and drop

## Problem

The launched app window should support file drag and drop.

## Desired behavior

- Dropping one `.json` settings file opens the settings UI for that file.
- Dropping one or more images treats them as input files and runs normal processing.
- Dropping mixed settings and images shows an actionable error.

## Acceptance criteria

- Multi-image D&D works.
- Settings JSON D&D works.
- Mixed D&D does not process and explains what to do.

