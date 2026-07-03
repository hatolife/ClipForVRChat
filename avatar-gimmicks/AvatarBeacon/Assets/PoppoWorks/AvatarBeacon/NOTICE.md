# NOTICE

This package includes material derived from YL-ATG_ForAvatar V0.0.3 by YozoraKurage.

Upstream project:
- https://github.com/YozoraKurage/YL-ATG

License:
- MIT License. See `LICENSES/YL-ATG-MIT.txt`.

Modifications in this tree:
- Renamed public asset path from `Assets/YozoLab/YL-ATG_ForAvatar` to `Assets/PoppoWorks/AvatarBeacon`.
- Renamed position parameters from `ATG/p/*` to `coord/*`.
- Renamed forward/rotation-vector parameters from `ATG/r/*` to `forward/*`.
- Removed auxiliary `ATG/SaveObject` and debug-only menu parameters that are not used for basis reconstruction.
- Changed the default Bone Proxy target from Head to Hips.
- Removed the visual-only arrow mesh/material assets.
- Normalized near-zero serialized Transform values in the prefab while preserving Contact/Constraint values used by the sensor graph.
- Adjusted prefab naming to `AvatarBeacon`.
