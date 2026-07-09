# Mapping engines — read before building the Fly View / 3D View

Per TRD.md §4.1 and §4.4, this frontend uses TWO separate map rendering
engines for two separate views — they are not merged into one map layer:

- **mapcn.dev (MapLibre GL)** — primary Fly View operational map (live
  vehicle position, mission overlay, geofence, rally points).
  mapcn.dev ships as copy-in shadcn-style components, not an npm package —
  install its map primitives via its own CLI per https://www.mapcn.dev/docs
  after running `pnpm add maplibre-gl` in this directory.

- **CesiumJS** — separate, lazy-loaded 3D View for terrain-aware mission
  review. Install with:

      pnpm add cesium resium

  Do not attempt to render mapcn.dev components inside a Cesium viewer or
  vice versa — they are different rendering engines.

Also install `msgpackr` for decoding the binary WebSocket telemetry stream:

    pnpm add msgpackr
