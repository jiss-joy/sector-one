# Sector One

Local telemetry overlay for **Assetto Corsa (2014)**. One Windows process talks to the game over UDP, records or replays raw packets, and serves a glance dash at `http://127.0.0.1:8080`.

No account. No cloud. Not affiliated with Kunos Simulazioni or 505 Games.

## Why this exists

The game already computes the car. The UDP packet is unreadable while you drive. The usual fix is a chart library and React state at 60 Hz. That janks, and it is the wrong owner for history.

Sector One treats the browser as a paint surface:

- The Go engine is the clock. It handshakes AC, parses `RTCarInfo`, and publishes **one JSON frame** per packet over Server-Sent Events (SSE).
- A slow tab must not stall ingest. The broadcaster keeps a buffer of one and **drops the oldest**.
- The dash keeps the last **600** frames in a ring. Widgets write the DOM and canvases from a **single** `requestAnimationFrame` loop. Connection chrome is the only `useSyncExternalStore` subscriber. There is no `setState` per packet and no Chart.js.
- After a session, `--replay` of a `.bin` uses the **same parsers** as live. The file stores raw UDP, not JSON. The browser never reads the file.

## Architecture

```
Assetto Corsa  ──UDP 127.0.0.1:9996──►  cmd/sector-one
                                            │
                     --record .bin  ◄───────┤
                     --replay .bin  ───────►│
                                            ▼
                                      Frame JSON
                                            │
                              SSE  /api/telemetry
                                            ▼
                         embedded Next.js static dash
                              (also :3000 in dev)
```

| Path | Role |
| --- | --- |
| `cmd/sector-one` | Flags, live loop, replay, HTTP, system tray |
| `internal/acudp` | Handshake + 328-byte car packet |
| `internal/physics` | `Frame` + per-car `cars.json` spec (max RPM / load peaks) |
| `internal/record` | `.bin` header `S1REC` v1, raw payloads |
| `internal/stream` | SSE + health, drop-oldest fan-out |
| `ui` | `go:embed` of the static export |
| `web` | Next.js App Router, static export |

## What the dash shows

One 16:9 screen.

- **Center:** car name, 20 LED RPM bar (limiter dot at 98% of `max_rpm`), gear + speed, steer bar, clutch / brake / throttle
- **Left:** friction circle (1G / 2G, faded tail), tyre load + LOCK / SLIP
- **Right:** TC / ABS as dim / armed / in-action (booleans only — the UDP packet has no click counts), laps + pit, pedal traces (index on X, not time)

Tyre temperatures, fuel, and a track map are **not** in this version. They are not in `RTCarInfo`.

## See it / try it

Three ways, easiest first. You do **not** need Assetto Corsa for 1 or 2.

**1. Replay a recorded session (no game)**

From the [latest Release](https://github.com/jiss-joy/sector-one/releases/latest) download `sector-one.exe` (Windows) **and** `sample.bin`. Put them in the same folder. The `.bin` is raw UDP from a real session, not a video.

Windows (PowerShell or cmd — you must pass the flag; double-clicking the exe is live mode):

```bat
sector-one.exe --replay=sample.bin --replay-rate=4
```

macOS / Linux after [building](#build-from-source):

```bash
./sector-one --replay=sample.bin --replay-rate=4
```

A tray icon appears, the browser should open [http://127.0.0.1:8080](http://127.0.0.1:8080). Replay at `4` is 4× real time. When the file ends, the process exits. **Open dashboard** / **Quit** are on the tray menu; clicking the icon does nothing.

**2. Watch a drive**

<!-- Replace this line with the video URL or a GIF once you have it. -->
Screen recording of a live session: *(add link)*.

**3. Live next to Assetto Corsa (2014)**

Install the game, download `sector-one.exe`, double-click it (or run it with no flags). Start a session in Content Manager — home screen is fine while it waits; you need to be on track for car packets. UDP is `127.0.0.1:9996` on the same machine. To keep a copy of that session:

```bat
sector-one.exe --record=session.bin
```

## Build from source

You need [Go](https://go.dev/dl/) **1.26+** (see `go.mod`) and [Bun](https://bun.sh) **1.3.2**.

`ui/public` is gitignored. A fresh clone will **not** compile until you export the UI once.

```bash
git clone https://github.com/jiss-joy/sector-one.git
cd sector-one
cd web && bun install --frozen-lockfile && cd ..
./scripts/build.sh
```

That builds `web/out`, copies it into `ui/public`, and writes `./sector-one`.

Same flags as [See it / try it](#see-it--try-it). There is no `--no-tray`. `go run` from a terminal still prints logs in that terminal — you launched a console process. That is not the shipped exe.

`--record` and `--replay` together is a fatal error.

| Flag | Default | Meaning |
| --- | --- | --- |
| `--http` | `127.0.0.1:8080` | Listen address |
| `--record` | (off) | Write raw UDP to this `.bin` |
| `--replay` | (off) | Play a `.bin` instead of UDP |
| `--replay-rate` | `1` | Speed. `<= 0` means no inter-frame sleep |

Dashboard: [http://127.0.0.1:8080](http://127.0.0.1:8080). Health: `/api/health`. Stream: `/api/telemetry`.

### UI hot-reload (optional)

```bash
# terminal 1 — engine (must already be built / public/ present)
go run ./cmd/sector-one --replay="$HOME/Desktop/session.bin"

# terminal 2
cd web && bun dev
```

Next on `:3000` talks to the engine on `:8080` (`web/lib/engine.ts` is hardcoded). Changing `--http` without changing that file breaks this path. The embedded UI is same-origin and does not care.

After dash changes you want on `:8080`, run `./scripts/build.sh` again. `go run` embeds whatever was in `ui/public` at **compile** time.

### Tests

```bash
go test ./internal/...
```

Parser offsets, record round-trip, and SSE backpressure are tested. There is no automated glance test.

## Windows release

Push to `main` runs [`.github/workflows/release-windows.yml`](.github/workflows/release-windows.yml): Bun export → copy into `ui/public` → `GOOS=windows` exe → GitHub Release `latest`.

## Constraints worth reading in the code

- `.bin` record timestamps stay nanoseconds (replay pacing). `Frame.ts` is always milliseconds at publish.
- `cars.json` in the working directory is a peak cache (default 6000 RPM / 5000 N). It is gitignored.
- Wheel order is FL / FR / RL / RR. If a bar looks mirrored, fix the label, not the parser.

## License

[MIT](LICENSE). That covers this repo only. You do not get a right to redistribute Assetto Corsa or its assets.
