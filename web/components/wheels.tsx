"use client";

import { type Ref, useRef } from "react";

import { useHeartbeat } from "@/hooks/use-heartbeat";
import { isLocked, isSlipping } from "@/lib/format";
import { getLatestFrame, lastN } from "@/lib/store";
import { THRESHOLDS } from "@/lib/thresholds";
import type { Frame } from "@/lib/types";

const sparkScratch: Frame[] = new Array(THRESHOLDS.LOAD_SPARK_COUNT);
const CORNER_INDEX = { FL: 0, FR: 1, RL: 2, RR: 3 } as const;

function paintCorner(
  load: number,
  slip: number,
  wheel: number,
  speed: number,
  bar: HTMLDivElement | null,
  flag: HTMLSpanElement | null,
  loadEl: HTMLSpanElement | null,
) {
  if (bar) {
    bar.style.height = `${Math.min(100, (load / THRESHOLDS.LOAD_BAR_N) * 100)}%`;
  }
  if (loadEl) loadEl.textContent = load.toFixed(0);
  if (flag) {
    const lock = isLocked(speed, wheel);
    const slipF = isSlipping(slip);
    flag.textContent = lock ? "LOCK" : slipF ? "SLIP" : "";
    flag.className = `rounded px-1 py-0.5 text-xs font-black uppercase tracking-wide ${
      lock
        ? "bg-red-500 text-white"
        : slipF
          ? "bg-amber-400 text-black"
          : "bg-transparent text-transparent"
    }`;
  }
}

function paintSpark(canvas: HTMLCanvasElement | null, n: number, corner: number) {
  if (!canvas || n < 2) return;
  const ctx = canvas.getContext("2d");
  if (!ctx) return;
  const dpr = window.devicePixelRatio || 1;
  const w = canvas.clientWidth;
  const h = canvas.clientHeight;
  if (w < 2 || h < 1) return;
  const bw = Math.round(w * dpr);
  const bh = Math.round(h * dpr);
  if (canvas.width !== bw || canvas.height !== bh) {
    canvas.width = bw;
    canvas.height = bh;
  }
  ctx.setTransform(dpr, 0, 0, dpr, 0, 0);
  ctx.clearRect(0, 0, w, h);
  ctx.beginPath();
  ctx.strokeStyle = "rgba(52,211,153,0.9)";
  ctx.lineWidth = 1;
  ctx.lineJoin = "round";
  for (let i = 0; i < n; i++) {
    const x = (i / (n - 1)) * w;
    const y = (1 - Math.min(1, sparkScratch[i].load_n[corner] / THRESHOLDS.LOAD_BAR_N)) * h;
    if (i === 0) ctx.moveTo(x, y);
    else ctx.lineTo(x, y);
  }
  ctx.stroke();
}

export function Wheels() {
  const barFL = useRef<HTMLDivElement>(null);
  const barFR = useRef<HTMLDivElement>(null);
  const barRL = useRef<HTMLDivElement>(null);
  const barRR = useRef<HTMLDivElement>(null);
  const flagFL = useRef<HTMLSpanElement>(null);
  const flagFR = useRef<HTMLSpanElement>(null);
  const flagRL = useRef<HTMLSpanElement>(null);
  const flagRR = useRef<HTMLSpanElement>(null);
  const loadFL = useRef<HTMLSpanElement>(null);
  const loadFR = useRef<HTMLSpanElement>(null);
  const loadRL = useRef<HTMLSpanElement>(null);
  const loadRR = useRef<HTMLSpanElement>(null);
  const sparkFL = useRef<HTMLCanvasElement>(null);
  const sparkFR = useRef<HTMLCanvasElement>(null);
  const sparkRL = useRef<HTMLCanvasElement>(null);
  const sparkRR = useRef<HTMLCanvasElement>(null);

  useHeartbeat(() => {
    const f = getLatestFrame();
    if (!f) return;
    const load = f.load_n ?? [0, 0, 0, 0];
    const slip = f.slip_ratio ?? [0, 0, 0, 0];
    const wh = f.wheel_rad_s ?? [0, 0, 0, 0];
    paintCorner(load[0], slip[0], wh[0], f.speed_kmh, barFL.current, flagFL.current, loadFL.current);
    paintCorner(load[1], slip[1], wh[1], f.speed_kmh, barFR.current, flagFR.current, loadFR.current);
    paintCorner(load[2], slip[2], wh[2], f.speed_kmh, barRL.current, flagRL.current, loadRL.current);
    paintCorner(load[3], slip[3], wh[3], f.speed_kmh, barRR.current, flagRR.current, loadRR.current);

    const n = lastN(THRESHOLDS.LOAD_SPARK_COUNT, sparkScratch);
    paintSpark(sparkFL.current, n, CORNER_INDEX.FL);
    paintSpark(sparkFR.current, n, CORNER_INDEX.FR);
    paintSpark(sparkRL.current, n, CORNER_INDEX.RL);
    paintSpark(sparkRR.current, n, CORNER_INDEX.RR);
  });

  return (
    <div className="flex w-full flex-col items-center gap-3">
      <div className="flex w-full items-center justify-between px-1">
        <span className="text-xs font-bold uppercase tracking-widest text-zinc-500">Tyre Load</span>
        <span className="font-mono text-[10px] text-zinc-600">
          scale {THRESHOLDS.LOAD_BAR_N} N
        </span>
      </div>
      <div className="grid grid-cols-2 gap-x-8 gap-y-3">
        <MiniCorner name="FL" barRef={barFL} flagRef={flagFL} loadRef={loadFL} sparkRef={sparkFL} />
        <MiniCorner name="FR" barRef={barFR} flagRef={flagFR} loadRef={loadFR} sparkRef={sparkFR} />
        <MiniCorner name="RL" barRef={barRL} flagRef={flagRL} loadRef={loadRL} sparkRef={sparkRL} />
        <MiniCorner name="RR" barRef={barRR} flagRef={flagRR} loadRef={loadRR} sparkRef={sparkRR} />
      </div>
    </div>
  );
}

function MiniCorner({
  name,
  barRef,
  flagRef,
  loadRef,
  sparkRef,
}: {
  name: string;
  barRef: Ref<HTMLDivElement>;
  flagRef: Ref<HTMLSpanElement>;
  loadRef: Ref<HTMLSpanElement>;
  sparkRef: Ref<HTMLCanvasElement>;
}) {
  return (
    <div className="flex w-28 flex-col gap-1.5">
      <div className="flex h-5 items-center justify-between">
        <span className="text-xs font-black text-zinc-400">{name}</span>
        <span ref={flagRef} className="text-xs font-black uppercase text-transparent">
          {" "}
        </span>
      </div>
      <canvas ref={sparkRef} className="h-1 w-full" />
      <div className="relative h-16 w-full overflow-hidden rounded border border-white/5 bg-zinc-900/80">
        <div
          className="absolute inset-x-0 top-0 z-10 h-px bg-white/40"
          title={`${THRESHOLDS.LOAD_BAR_N} N`}
        />
        <div
          ref={barRef}
          className="absolute inset-x-0 bottom-0 bg-emerald-500/25 border-t border-emerald-400 transition-[height] duration-75"
          style={{ height: "0%" }}
        />
      </div>
      <div className="flex items-baseline justify-center gap-1">
        <span ref={loadRef} className="font-mono text-base font-bold tabular-nums text-white">
          0
        </span>
        <span className="text-[10px] font-bold uppercase text-zinc-500">N</span>
      </div>
    </div>
  );
}
