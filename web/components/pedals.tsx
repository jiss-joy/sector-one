"use client";

import { useRef } from "react";

import { useHeartbeat } from "@/hooks/use-heartbeat";
import { clamp01 } from "@/lib/format";
import { getLatest } from "@/lib/store";

export function Pedals() {
  const thrBar = useRef<HTMLDivElement>(null);
  const brkBar = useRef<HTMLDivElement>(null);
  const cluBar = useRef<HTMLDivElement>(null);
  const thrLbl = useRef<HTMLSpanElement>(null);
  const brkLbl = useRef<HTMLSpanElement>(null);
  const cluLbl = useRef<HTMLSpanElement>(null);

  useHeartbeat(() => {
    const f = getLatest();
    if (!f) return;
    const t = clamp01(f.throttle);
    const b = clamp01(f.brake);
    const c = clamp01(f.clutch);
    if (thrBar.current) thrBar.current.style.height = `${t * 100}%`;
    if (brkBar.current) brkBar.current.style.height = `${b * 100}%`;
    if (cluBar.current) cluBar.current.style.height = `${c * 100}%`;
    if (thrLbl.current) thrLbl.current.textContent = (t * 100).toFixed(0);
    if (brkLbl.current) brkLbl.current.textContent = (b * 100).toFixed(0);
    if (cluLbl.current) cluLbl.current.textContent = (c * 100).toFixed(0);
  });

  return (
    <div className="flex h-44 items-end justify-around gap-6">
      <PedalCol barRef={cluBar} labelRef={cluLbl} name="CLU" color="bg-sky-500" />
      <PedalCol barRef={brkBar} labelRef={brkLbl} name="BRK" color="bg-red-500" />
      <PedalCol barRef={thrBar} labelRef={thrLbl} name="THR" color="bg-emerald-500" />
    </div>
  );
}

function PedalCol({
  barRef,
  labelRef,
  name,
  color,
}: {
  barRef: React.RefObject<HTMLDivElement>;
  labelRef: React.RefObject<HTMLSpanElement>;
  name: string;
  color: string;
}) {
  return (
    <div className="flex h-full flex-1 flex-col items-center gap-2">
      <span
        ref={labelRef}
        className="font-mono text-sm font-bold tabular-nums text-white/80"
      >
        0
      </span>
      <div className="relative w-full flex-1 overflow-hidden rounded-sm bg-zinc-800/50">
        <div ref={barRef} className={`absolute inset-x-0 bottom-0 ${color} shadow-[0_0_10px_rgba(0,0,0,0.5)] transition-all duration-75`} />
      </div>
      <span className="text-[10px] font-bold tracking-wider text-zinc-500 uppercase">
        {name}
      </span>
    </div>
  );
}
