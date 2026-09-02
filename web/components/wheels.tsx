"use client";

import { Ref, useRef } from "react";

import { useHeartbeat } from "@/hooks/use-heartbeat";
import { isLocked, isSlipping } from "@/lib/format";
import { getLatestFrame } from "@/lib/store";

function paintCorner(
  load: number,
  slip: number,
  wheel: number,
  speed: number,
  bar: HTMLDivElement | null,
  flag: HTMLSpanElement | null,
  loadEl: HTMLSpanElement | null,
  maxLoad: number,
) {
  if (bar) bar.style.height = `${Math.min(100, (load / maxLoad) * 100)}%`;
  if (loadEl) loadEl.textContent = load.toFixed(0);
  if (flag) {
    const lock = isLocked(speed, wheel);
    const slipF = isSlipping(slip);
    flag.textContent = lock ? "LOCK" : slipF ? "SLIP" : "";
    flag.className = `font-mono text-[10px] tabular-nums ${
      lock ? "text-red-400" : slipF ? "text-amber-400" : "text-transparent"
    }`;
  }
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

  useHeartbeat(() => {
    const f = getLatestFrame();
    if (!f) return;
    const load = f.load_n ?? [0, 0, 0, 0];
    const slip = f.slip_ratio ?? [0, 0, 0, 0];
    const wh = f.wheel_rad_s ?? [0, 0, 0, 0];
    paintCorner(load[0], slip[0], wh[0], f.speed_kmh, barFL.current, flagFL.current, loadFL.current, f.max_load);
    paintCorner(load[1], slip[1], wh[1], f.speed_kmh, barFR.current, flagFR.current, loadFR.current, f.max_load);
    paintCorner(load[2], slip[2], wh[2], f.speed_kmh, barRL.current, flagRL.current, loadRL.current, f.max_load);
    paintCorner(load[3], slip[3], wh[3], f.speed_kmh, barRR.current, flagRR.current, loadRR.current, f.max_load);
  });

  return (
    <div className="flex flex-col gap-2">
      <span className="text-xs font-bold tracking-widest text-zinc-500 uppercase">Tyre Load</span>
      <div className="grid grid-cols-2 gap-2">
        <MiniCorner name="FL" barRef={barFL} flagRef={flagFL} loadRef={loadFL} />
        <MiniCorner name="FR" barRef={barFR} flagRef={flagFR} loadRef={loadFR} />
        <MiniCorner name="RL" barRef={barRL} flagRef={flagRL} loadRef={loadRL} />
        <MiniCorner name="RR" barRef={barRR} flagRef={flagRR} loadRef={loadRR} />
      </div>
    </div>
  );
}

function MiniCorner({
  name,
  barRef,
  flagRef,
  loadRef,
}: {
  name: string;
  barRef: Ref<HTMLDivElement>;
  flagRef: Ref<HTMLSpanElement>;
  loadRef: Ref<HTMLSpanElement>;
}) {
  return (
    <div className="relative flex h-24 w-20 flex-col items-center justify-between rounded border border-white/5 bg-zinc-900/50 p-2 overflow-hidden">
      {/* Background fill bar */}
      <div 
        ref={barRef} 
        className="absolute inset-x-0 bottom-0 bg-emerald-500/20 transition-all duration-75 border-t-2 border-emerald-500" 
        style={{ height: "0%" }} 
      />
      
      <div className="z-10 flex w-full justify-between items-start">
        <span className="text-[10px] font-bold text-zinc-500">{name}</span>
        <span ref={flagRef} className="text-[8px] font-bold uppercase tracking-tighter"></span>
      </div>

      <div className="z-10 flex flex-col items-center">
        <span ref={loadRef} className="font-mono text-sm font-bold text-white tabular-nums">0</span>
        <span className="text-[7px] font-bold text-white/20 uppercase">N</span>
      </div>
    </div>
  );
}
