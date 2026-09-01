"use client";

import { useRef, type ReactNode } from "react";

import { useHeartbeat } from "@/hooks/use-heartbeat";
import { formatLap } from "@/lib/format";
import { getLatest } from "@/lib/store";

export function Laps() {
  const timeRef = useRef<HTMLSpanElement>(null);
  const lastRef = useRef<HTMLSpanElement>(null);
  const bestRef = useRef<HTMLSpanElement>(null);
  const countRef = useRef<HTMLSpanElement>(null);
  const deltaRef = useRef<HTMLSpanElement>(null);
  const pitRef = useRef<HTMLSpanElement>(null);

  useHeartbeat(() => {
    const f = getLatest();
    if (!f) return;
    if (timeRef.current) timeRef.current.textContent = formatLap(f.lap_time_ms);
    if (lastRef.current) lastRef.current.textContent = formatLap(f.last_lap_ms);
    if (bestRef.current) bestRef.current.textContent = formatLap(f.best_lap_ms);
    if (countRef.current) countRef.current.textContent = String(f.lap_count);
    if (deltaRef.current) {
      if (f.last_lap_ms > 0 && f.best_lap_ms > 0) {
        const d = f.last_lap_ms - f.best_lap_ms;
        const sign = d > 0 ? "+" : "";
        deltaRef.current.textContent = `${sign}${(d / 1000).toFixed(3)}`;
        deltaRef.current.className = `font-mono tabular-nums ${
          d > 0 ? "text-red-400" : d < 0 ? "text-emerald-400" : ""
        }`;
      } else {
        deltaRef.current.textContent = "—";
        deltaRef.current.className = "font-mono tabular-nums text-muted-foreground";
      }
    }
    if (pitRef.current) {
      pitRef.current.dataset.active = f.in_pit ? "1" : "0";
    }
  });

  return (
    <div className="flex flex-col gap-4">
      <div className="flex flex-col gap-2 rounded border border-white/5 bg-zinc-900/40 p-4">
        <span className="text-[10px] font-bold tracking-widest text-zinc-500 uppercase">Timing</span>
        <div className="space-y-1">
          <Row label="current">
            <span ref={timeRef} className="font-mono text-xl tabular-nums text-white">
              --:--.---
            </span>
          </Row>
          <Row label="delta">
            <span ref={deltaRef} className="font-mono text-lg tabular-nums text-muted-foreground">
              —
            </span>
          </Row>
          <Row label="last">
            <span ref={lastRef} className="font-mono tabular-nums">
              --:--.---
            </span>
          </Row>
          <Row label="best">
            <span ref={bestRef} className="font-mono tabular-nums text-emerald-400/80">
              --:--.---
            </span>
          </Row>
        </div>
      </div>
      
      <div className="flex items-center justify-between px-2">
        <div className="flex flex-col">
          <span className="text-[8px] font-bold text-zinc-500 uppercase">Lap</span>
          <span ref={countRef} className="font-mono text-2xl font-bold text-white">0</span>
        </div>
        <Light ref={pitRef} label="PIT" />
      </div>
    </div>
  );
}

function Row({
  label,
  children,
}: {
  label: string;
  children: ReactNode;
}) {
  return (
    <div className="flex items-baseline justify-between gap-4">
      <span className="text-[10px] uppercase tracking-wide text-muted-foreground">
        {label}
      </span>
      {children}
    </div>
  );
}

function Light({
  ref,
  label,
}: {
  ref: Ref<HTMLSpanElement>;
  label: string;
}) {
  return (
    <span
      ref={ref}
      data-active="0"
      className="inline-flex h-5 items-center rounded-full border border-border px-2 text-[10px] font-medium tracking-wide data-[active=1]:border-amber-400 data-[active=1]:bg-amber-400/20 data-[active=1]:text-amber-200"
    >
      {label}
    </span>
  );
}
