# Rebuild from scratch

This is a rebuild map of Sector One as it exists today. It is not the public README. Delete or gitignore this file before you publish.

No code. The current repo is the oracle. When a gate fails, diff the matching package. Do not treat this file as a second implementation.

## How to use this

Rebuild in a sibling directory or a new branch. Do not wipe this tree first.

Finish a phase only when its gate passes. Do not start the dash until Server-Sent Events (SSE) frames exist. Do not start the tray until replay and embed work.

`go run` serves embedded static files, not `web/` source. If the dash looks old, you did not rebuild-and-copy. That is not a React bug.

## Product

Local-only Assetto Corsa (2014) telemetry. One process on the rig.

The Go engine talks User Datagram Protocol (UDP) to the game on `127.0.0.1:9996`, optionally records raw packets to a `.bin`, replays that file without the game, and publishes frames over Hypertext Transfer Protocol (HTTP) on `127.0.0.1:8080`. The Next.js app is a static export glued into the executable. The browser is a paint surface. It does not own history, files, or the game session.

AC UDP 9996 → live loop → parsers → Frame JSON → SSE → dash. `--record` writes raw UDP from the live loop. `--replay` feeds the same parsers from a `.bin`. No JSON in the file.

Live and replay share the same parsers and the same publish path.

Windows is the ship target (one exe, system tray, auto-open browser). You develop on macOS with `--replay` and `--no-tray`.

## Hard rules

- Client ring holds 600 frames. Do not grow it “just in case.”
- One shared `requestAnimationFrame` loop. Widgets write Document Object Model (DOM) nodes and canvases through refs. No React state per SSE message.
- `useSyncExternalStore` is for connection chrome only, not the 60 Hz path.
- No Chart.js, Recharts, or any chart library.
- Record writer timestamps stay nanoseconds. Do not change the writer to milliseconds to make graphs easier.
- UDP `RTCarInfo` has TC/ABS as booleans (enabled / in action). Do not invent click counts.
- No fake tyre temperatures. Those are not on this UDP packet.
- Do not parse `.bin` in the browser. After-session viewing is `--replay`.
- No track map, replay scrubber, coaching, or file picker in Next.
- Do not invert wheel order in the parser to fix a label. Fix labels or a plot flip in the UI.

## Why `go run` showed the old UI

The engine serves whatever was in `src/ui/public` at **compile** time (`go:embed`). That directory is gitignored. Editing `web/` does nothing to `:8080` until you export and copy.

`scripts/build.sh` exists for that reason: Next static export → copy `web/out` into `src/ui/public` → `go build`. The Windows CI job does the same copy **before** `go build`, because a checkout has no `public/` folder.

`bun dev` on port 3000 hot-reloads source. The engine on 8080 does not. Two different servers, two different bundles.

`go generate` on the UI package points at the full build script. It only runs if you invoke generate. Do not assume `go run` rebuilds the dash.

## Phase 0 — Skeleton

**Goal.** Two trees that can compile and a place for the embed.

**Create.**

- Go module under `src/` (module path `sector-one`). `go.work` at the repo root so `go run ./src/cmd/engine` works from the root.
- Next app under `web/` with static export enabled and images unoptimized (required for `output: 'export'`).
- `src/ui` package that embeds a `public/` directory and returns an HTTP file server.
- Gitignore: `web/node_modules`, `web/.next`, `web/out`, `src/ui/public`, `cars.json`, `*.exe`.
- A one-line HTML page in `public/` so `go run` serves *something*.

**Contract.** Engine listens on `127.0.0.1:8080` and serves `/`. Next can still be developed on 3000 later.

**Pitfalls.** If embed has no files, `go:embed` of an empty or missing dir fails the build. Put a placeholder in `public/` from day one, even if gitignored locally.

**Gate.** `go run` prints that it is listening. A browser on 8080 shows the placeholder. Next `build` produces `web/out`.

## Phase 1 — UDP parsers

**Goal.** Bytes in, structs out. Highest-risk layer. Tests first.

**Create.** Package `internal/acudp`: handshake encode/parse, car-info parse. Tests with captured or fixture packets.

**Contract.**

- Handshake request is 12 bytes: identifier, version, operation. Operations you need: handshake (0), subscribe updates (1), dismiss (3). Spot subscribe exists in the protocol; you do not use it.
- Handshake reply is 408 or 808 bytes. Four UTF-16LE names (car, driver, track, config) plus two int32s. Kunos often terminates names at `%` or NUL. Trim that.
- Car packet is 328 bytes. First byte is identifier `a`. Reject any other length or identifier.
- Gear: 0 reverse, 1 neutral, 2+ forward. The dash later displays forward gears as `gear - 1`.
- Wheel arrays are Kunos order: front-left, front-right, rear-left, rear-right. [Likely] confirmed on one Le Mans recording (load moves right in a left-hander). Do not invert in the parser if a widget looks wrong — check labels.

**Pitfalls.** Handshake field widths change between 408 and 808; compute field size from packet length, do not hardcode one layout. Car packet has a size field you can ignore. Several leading bytes are unused. Read the existing tests and `carinfo.go` when you need offsets. Do not copy offsets into this document.

**Gate.** Tests parse a handshake and a car packet. Names are readable. Speed, RPM, gear, pedals, steer, loads, and the boolean flags match a known capture. You have not written HTTP yet.

## Phase 2 — Frame and specs

**Goal.** One JSON object the UI will consume for the rest of the project.

**Create.** Package `internal/physics`: `Frame`, `FromCar`, spec manager + `cars.json`.

**Contract.**

- JSON keys are snake_case. The TypeScript `Frame` type must match exactly.
- `FromCar` maps UDP names: gas → throttle; the three acceleration axes → `g_vert` / `g_lat` / `g_long`; session car name → `car`; live or replay → `source`.
- Spec manager is per car name. Defaults: 6000 RPM, 5000 N. Raise max RPM on a higher sample. If the limiter flag is set and RPM is positive, treat that RPM as max. Raise max load on a higher corner load. Persist the whole map to `cars.json` on change (process working directory).
- Unknown car name until a handshake is seen: `"unknown"`. Replay without a hello record stays on that key.

**Pitfalls.** `max_rpm` / `max_load` on the frame come from the spec manager, not the UDP packet. The dash RPM bar is `rpm / max_rpm`. If specs never update, the bar is a lie. Load bars in the current UI use a UI constant, not `max_load` — keep that split when you rebuild widgets.

**Gate.** Marshal a frame and assert the JSON keys the UI will type. Spec update raises RPM and load and rewrites `cars.json`.

## Phase 3 — Record format

**Goal.** Record raw UDP. Replay the files you already have.

**Create.** Package `internal/record`: header, write, read, tests.

**Contract.**

- File header: 16 bytes. Magic `S1REC` (padded), version 1. Reader validates magic and version only.
- Each record: kind (hello = 1, car = 2), timestamp (uint64 little-endian), payload length, payload.
- Writer timestamp is wall-clock **nanoseconds**.
- Payload is the raw UDP datagram, not a Frame.
- Unknown kinds: log and skip. Payload larger than 64 KiB: reject.

**Pitfalls.** Do not change the writer to milliseconds. Existing `lemans.bin` and `sleepDelta` assume nanoseconds. A “fix” here silently breaks every recording you already have. Live Frame `ts` is a different unit (Phase 5). Leave both.

**Gate.** Write a tiny file, read it back. Open `lemans.bin` and get hello + car records with a car name from the hello payload.

## Phase 4 — SSE

**Goal.** Fan-out frames to any number of tabs without stalling ingest.

**Create.** Package `internal/stream`: broadcaster + HTTP handler.

**Contract.**

- `GET /api/health` — JSON: ok, source (`live` or `replay`), subscriber count.
- `GET /api/telemetry` — `text/event-stream`. One SSE message per frame: `data:` + JSON + blank line. Flush immediately.
- Subscribe channel buffer is 1. On publish, if the client is still holding the previous frame, drop the old one and send the latest. A slow tab must never block the game loop.
- CORS: reflect a localhost Origin (Next on 3000), otherwise `*`. OPTIONS on both routes.

**Pitfalls.** If you buffer unbounded, a paused tab becomes a memory leak. If you block on send, UDP ingest stalls. Health is not telemetry; do not send frames on the health route.

**Gate.** Handler under test or a short `curl` of health. A throwaway EventSource client prints frames. Disconnecting a client decrements subscribers.

## Phase 5 — Live and replay

**Goal.** One `main` that can talk to AC or play a `.bin`, and publish on both paths.

**Create.** `cmd/engine`: flags, HTTP start, live loop, replay loop, logging helpers. Wire publish: spec update → `FromCar` → JSON → broadcaster.

**Flags.**

| Flag | Role |
| --- | --- |
| `--http` | Listen address. Default `127.0.0.1:8080`. |
| `--record` | Write a `.bin` while live. |
| `--replay` | Play a `.bin`. No UDP. |
| `--replay-rate` | Speed multiplier. `<= 0` means no inter-frame sleep. |
| `--record` and `--replay` together | Fatal. |

Tray flags come in Phase 9. Until then, block on context cancel (Ctrl-C).

**Live loop.**

1. Dial UDP to `127.0.0.1:9996`. Fail → log, sleep 2s, retry.
2. Send handshake. Read with a 5s deadline. Timeout is normal if Content Manager is on the home screen — wait and retry. Do not treat that as a crash.
3. Parse hello. Print car / driver / track / config. Optionally record the raw hello.
4. Send subscribe-updates.
5. Inner read loop, 5s deadline: parse car, optional record, calibrate, publish with `time.Now().UnixMilli()` as `ts`.
6. Read error (session ended) → dismiss, close, outer loop (new session). Cancel → dismiss, close, return.

**Replay loop.**

- No UDP. Read records. On hello: parse session, keep car name. On car: same calibrate + publish as live, but `ts` is the **record timestamp** (nanoseconds).
- `sleepDelta`: treat `now - prev` as a `time.Duration` (nanoseconds). Divide by rate. Cap any gap above 200ms (paused AC / Content Manager). First record and `rate <= 0`: no sleep.

**Pitfalls.**

- Live `ts` is milliseconds. Replay `ts` is nanoseconds. A “last 8 seconds” graph using `now - 8000` is 8 microseconds on replay and looks empty. **Leave this.** Time-axis traces are out of scope for the rewrite.
- Replay without a hello in the file: car name stays `"unknown"`.
- HTTP server is not shut down gracefully. In-flight SSE dies on process exit. Acceptable for this product.
- Next `engine.ts` hardcodes `127.0.0.1:8080`. Changing `--http` without changing that file breaks `bun dev`. Embedded UI is same-origin, so it does not care.

**Gate.** `--replay` of `lemans.bin` streams frames on `/api/telemetry`. `--record` a live session (or a short synthetic write), then `--replay` that file and get the same car name and plausible numbers. Ingest log every 2s (speed, gear, RPM, Hz) is enough observability.

## Phase 6 — Dash data path

**Goal.** Frames on screen without a pretty dash. Prove the paint model.

**Create.**

- Client-only home page. `useEffect` connects once.
- `lib/engine.ts`: EventSource on `/api/telemetry`, health poll every 2s. Hardcode `http://127.0.0.1:8080`.
- `lib/store.ts`: latest frame, ring of 600, connection object, connection listeners.
- `hooks/use-heartbeat.ts`: one global rAF loop; widgets register draw callbacks.
- `lib/types.ts`: `Frame` matching Go JSON; connection `offline | online | error`.

**Contract.**

- `pushFrame` writes the ring and latest. If `frame.source` changes, update connection source and notify listeners.
- Connection: EventSource `onopen` → online. EventSource `onerror` → **offline** (the socket retries). Health fetch network failure (connection refused) → **offline**. Health HTTP non-OK → **error**. Do not reuse older names `down` / `open`.
- `useSyncExternalStore` only in the header later. Widgets in this phase read `getLatestFrame()` inside the heartbeat callback and write a text node.
- Cleanup on unmount: close EventSource, clear health interval, set offline.

**Pitfalls.** If you `useState` every SSE message, the dash will jank and you will “fix” it with memo for the rest of the rewrite. If each widget starts its own rAF, you have N loops. Health `subscribers` is stored; the current header does not display it — optional.

**Gate.** With the engine replaying, a single RPM or speed text node updates smoothly. React DevTools does not re-render the page per packet. Disconnect the engine: pill path (even a console log) goes offline, not error.

## Phase 7 — Glance widgets

**Goal.** A dash you can read while driving. Calibrate before you decorate.

**Create.** One screen. Dark. No sidebar, no theme switcher, no settings page. shadcn leftovers (button, card, shine, border-beam) are not part of the product — do not rebuild them.

**Calibrate first** (from a known recording; current notes used `lemans.bin` / Adonis at Le Mans):

1. Left-hander: load must move to front-right / rear-right. If bars invert, swap **labels**, not parser order.
2. Hard brake: longitudinal G must move toward the labelled Brake end. One boolean flip in thresholds if it does not. Comment the session.
3. `LOAD_BAR_N`: peak single-corner load on that file plus about 20% headroom. Current value: 6000 N (peak sat near the engine default of 5000). A bar that is always 10% full is a broken scale.

Put the three results in a comment on `thresholds.ts`. Scatter these numbers across components and you will retune forever.

**Build widgets in this order.**

1. **Header.** Brand. Source string. Pill: Live / Offline / Error. This is the only `useSyncExternalStore` consumer.
2. **Cluster.** Car name (`car`, underscores to spaces). 20 RPM LEDs, fill = `rpm / max_rpm`. Color bands by index (blue → green → yellow → red). Dedicated limiter **dot** at ≥ 98% of max RPM — violent pulse on that dot only, not the whole bar. Do not drive the dot from the `engine_limiter` flag (the current UI ignores that field). Gear huge, speed beside it, one row. Gear label: 0 → R, 1 → N, else `gear - 1`.
3. **Pedals.** Clutch / brake / throttle, 0–1, percent labels.
4. **Steer.** Own horizontal bar. Input is −1…+1. Marker at `(steer + 1) / 2`. Not a trace.
5. **Status boxes.** TC and ABS. Three opacities: off, enabled, in action. No number.
6. **Laps.** Current / last / best, delta last − best after a completed lap, lap count, pit light. Times `MM:SS.mmm`. Zero or negative → placeholder dashes.
7. **Friction circle.** Labels Brake (top), Accel (bottom), Left, Right. 1G and 2G rings. Scale 2.5 G. Faded tail of the last 180 points. Live total G. Longitudinal flip on.
8. **Wheels.** Four corners, FL FR / RL RR. Load bar against `LOAD_BAR_N`. Newtons readout. LOCK if speed > 20 km/h and wheel rad/s nearly stopped. SLIP if absolute slip ratio > 0.15. LOCK wins. Short load sparkline (~90 samples). No temperatures.
9. **Traces.** Throttle, brake, clutch. X is **array index**, not time. Grid at 0 / 0.5 / 1. Steer is not on this chart. Entire ring (up to 600).

**Layout.**

- Left: friction circle, wheels.
- Center: cluster, steer, pedals.
- Right: status + laps, traces in a short fixed-height panel.

TC/ABS used to live large in the center and forced page scroll. Keep them small on the right.

**Fields you parse but do not draw:** `g_vert`, `engine_limiter`, `slip_angle`, `normalized_pos`, `max_load`. Keep them on the Frame. Do not invent widgets for them in this rewrite.

**Pitfalls.** Time-windowed traces already failed once: replay timestamps are nanoseconds, so `now - 8000` is nothing, and a `min-h-0` canvas can also size to zero. Index X is the honest rebuild. Pit-lane crawl looks dead — that is the data; do not add idle animation to hide it.

**Gate.** Replay looks glanceable at a 16:9 window. Car name appears after the hello record. Steer bar moves independently of traces. LOCK/SLIP are readable at a glance. You did not add Chart.js or tyre temps.

## Phase 8 — Embed and ship path

**Goal.** One process, no Bun on the rig.

**Create.**

- Build script: Next export → wipe and copy into `src/ui/public` → `go build` the engine.
- GitHub Actions Windows job: checkout, Go from `src/go.mod`, Bun, `bun install --frozen-lockfile`, `bun run build` in `web/`, copy `web/out` → `src/ui/public`, then `GOOS=windows` build. Publish the exe. `src/ui/public` is gitignored; if CI skips the copy, embed is empty or stale and the job is a lie.

**Contract.** Opening `http://127.0.0.1:8080` with no Next dev server shows the glance dash. Same EventSource URL, now same-origin.

**Pitfalls.** You already shipped a Windows exe that served the old dash because you committed tray on `main` and never re-exported the UI. Rebuild-and-copy is part of the product, not a local hack. After every dash change you intend to see on 8080, run the script (or the two copy steps) **before** `go run`.

**Gate.** Kill Next. Run the engine with `--replay`. Browser on 8080 is the new UI. CI workflow contains the copy step before `go build`.

## Phase 9 — Tray last

**Goal.** Double-click an exe on Windows. Console hidden. Browser opens. Quit from the menu.

**Create.** Only after Phase 8.

- Tray + icons (Windows `.ico`, others `.png`).
- Hide-console on Windows; no-op stub on other operating systems.
- Flag `--no-tray` for logs and for macOS day-to-day.
- Menu: Open dashboard, Quit. Auto-open browser on tray ready and on `--no-tray`.
- Live/replay run in a goroutine. Main thread blocks in the tray loop (or on context if `--no-tray`).

**Lifecycle.**

| Mode | After replay EOF | Browser |
| --- | --- | --- |
| Tray default | Process **stays up** (last frames remain on SSE) | Auto on tray ready |
| `--no-tray` | Process **exits** | Auto once |
| Live | Loops until signal / Quit | Same as above |

Icon click is **not** wired. Only the menu item opens the browser. Do not add `SetOnClick` in this rewrite unless you are done with everything above.

**Pitfalls.** Tray on macOS is awkward while iterating; use `--no-tray`. If you add tray before embed works, you will debug the wrong layer. Replay + tray staying alive is deliberate: the user still wants the last dash. Replay + `--no-tray` exiting is how you keep a log session finite.

**Gate.** Windows exe: tray icon, Open dashboard, Quit, console hidden. Mac: `go run --replay=… --no-tray` still prints logs and exits at EOF.

## After the rewrite — do not build these yet

These are intent, not the current product. If you add them during the rewrite, you are not reconstructing “until now.”

- Time-axis traces. Blocked on live milliseconds vs replay nanoseconds. Fix `ts` at **publish** time first (one unit on every Frame), then window the canvas. Do not change the record writer.
- Shared memory for tyre temps. Windows mapping `Local\acpmf_physics`. Stamp fields onto the same Frame, same SSE. Do not replace UDP. Mac replay will be empty for those fields unless you recorded them (you do not, today).
- TC/ABS click counts. Not in UDP. Shared memory physics values are 0–1, not “ABS 5.” Real clicks need per-car `electronics.ini`. Leave the boxes as booleans.
- Track map, scrubber, coaching, Next file picker.
- Tray icon click → open browser.
- Public README. Write that after the rewrite, from what you actually shipped.

## Suggested package map

When you get lost, this is the only architecture that matches the current repo.

- `internal/acudp` — bytes ↔ session / car. Stdlib only.
- `internal/record` — `.bin` header and records. Stdlib only.
- `internal/physics` — Frame, specs. Depends on acudp.
- `internal/stream` — broadcaster + `/api/*`. Stdlib only.
- `ui` — embed + static handler.
- `cmd/engine` — the only place that ties UDP/replay, specs, HTTP, and tray.

`cmd/engine` must stay thin. If parsers or SSE logic leak into `main`, you will not be able to test them.

## Tests you must write as you go

If you skip these, you will re-learn offsets on the live rig.

- Handshake encode and both reply sizes.
- Car-info length, identifier, a few known fields.
- Frame JSON keys.
- Record header + one record round-trip; reject bad magic / huge payload.
- Broadcaster: two subscribers; a slow one still receives latest, not a backlog.

There is no automated glance test. Phase 7’s gate is you, a replay, and your eyes.

## Done

You are done with the rewrite when:

1. `--replay` of `lemans.bin` drives the glance dash over SSE.
2. `:8080` with no Next process shows that dash (embed is current).
3. `--no-tray` still works on your Mac.
4. You did not add a time axis, shared memory, fake electronics clicks, or a README.

Then write the public README from this product, not from this file.
