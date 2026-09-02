"use client";

import { useRef } from "react";

import { useHeartbeat } from "@/hooks/use-heartbeat";
import { getLatestFrame } from "@/lib/store";

export function Steer() {
  const markerRef = useRef<HTMLDivElement>(null);

  useHeartbeat(() => {
    const f = getLatestFrame();
    if (!f || !markerRef.current) return;
    const pct = ((f.steer + 1) / 2) * 100;
    markerRef.current.style.left = `${Math.min(100, Math.max(0, pct))}%`;
  });

  return (
    <div className="flex w-full items-center gap-3">
      <span className="w-4 text-[10px] font-black uppercase tracking-widest text-zinc-500">L</span>
      <div className="relative h-1.5 flex-1 rounded-full bg-zinc-800">
        <div className="absolute left-1/2 top-1/2 h-2.5 w-px -translate-x-1/2 -translate-y-1/2 bg-white/30" />
        <div
          ref={markerRef}
          className="absolute top-1/2 h-3 w-1.5 -translate-x-1/2 -translate-y-1/2 rounded-sm bg-amber-400 shadow-[0_0_8px_var(--color-amber-400)]"
          style={{ left: "50%" }}
        />
      </div>
      <span className="w-4 text-right text-[10px] font-black uppercase tracking-widest text-zinc-500">R</span>
    </div>
  );
}
