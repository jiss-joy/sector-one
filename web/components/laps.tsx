"use client";

import { Ref, useRef } from "react";

import { useHeartbeat } from "@/hooks/use-heartbeat";
import { formatGap, formatLap } from "@/lib/format";
import { getCompletedLaps, getLatestFrame } from "@/lib/store";

const ROW_COUNT = 5;

export function Laps() {
  const timeRef = useRef<HTMLDivElement>(null);
  const bestRef = useRef<HTMLDivElement>(null);
  const pitRef = useRef<HTMLSpanElement>(null);
  const headRef = useRef<HTMLTableSectionElement>(null);
  const rowRefs = useRef<(HTMLTableRowElement | null)[]>([]);
  const lapRefs = useRef<(HTMLTableCellElement | null)[]>([]);
  const driverRefs = useRef<(HTMLTableCellElement | null)[]>([]);
  const timeRefs = useRef<(HTMLTableCellElement | null)[]>([]);
  const gapRefs = useRef<(HTMLTableCellElement | null)[]>([]);

  useHeartbeat(() => {
    const f = getLatestFrame();
    if (f) {
      if (timeRef.current) timeRef.current.textContent = formatLap(f.lap_time_ms);
      if (bestRef.current) bestRef.current.textContent = formatLap(f.best_lap_ms);
      if (pitRef.current) {
        pitRef.current.hidden = !f.in_pit;
        pitRef.current.dataset.active = f.in_pit ? "1" : "0";
      }
    }

    const laps = getCompletedLaps();
    if (headRef.current) {
      headRef.current.hidden = laps.length === 0;
    }
    for (let i = 0; i < ROW_COUNT; i++) {
      const row = laps[laps.length - 1 - i];
      const empty = !row;
      if (rowRefs.current[i]) rowRefs.current[i]!.hidden = empty;
      if (lapRefs.current[i]) lapRefs.current[i]!.textContent = empty ? "" : String(row.lap);
      if (driverRefs.current[i]) driverRefs.current[i]!.textContent = empty ? "" : row.driver || "—";
      if (timeRefs.current[i]) timeRefs.current[i]!.textContent = empty ? "" : formatLap(row.timeMs);
      if (gapRefs.current[i]) {
        const el = gapRefs.current[i]!;
        el.textContent = empty ? "" : formatGap(row.gapMs);
        el.className = gapClass(empty ? null : row.gapMs);
      }
    }
  });

  return (
    <div className="flex flex-col gap-3">
      <div className="flex items-end justify-between gap-3 px-0.5">
        <div>
          <div className="flex items-center gap-2">
            <span className="text-[10px] font-bold tracking-widest text-zinc-500 uppercase">Current</span>
            <Light ref={pitRef} label="PIT" />
          </div>
          <div ref={timeRef} className="font-mono text-xl tabular-nums text-white">
            --:--.---
          </div>
        </div>
        <div className="text-right">
          <span className="text-[10px] font-bold tracking-widest text-zinc-500 uppercase">Best</span>
          <div ref={bestRef} className="font-mono text-sm tabular-nums text-emerald-400/80">
            --:--.---
          </div>
        </div>
      </div>

      <table className="w-full border-collapse text-left">
        <caption className="sr-only">Last five completed laps, gap versus previous lap</caption>
        <thead ref={headRef} hidden>
          <tr className="text-[10px] font-bold tracking-widest text-zinc-500 uppercase">
            <th className="pb-1 pr-2 font-bold">Lap</th>
            <th className="pb-1 pr-2 font-bold">Driver</th>
            <th className="pb-1 pr-2 font-bold">Time</th>
            <th className="pb-1 font-bold text-right">Gap</th>
          </tr>
        </thead>
        <tbody className="font-mono text-sm tabular-nums">
          {Array.from({ length: ROW_COUNT }, (_, i) => (
            <tr
              key={i}
              ref={(el) => {
                rowRefs.current[i] = el;
              }}
              hidden
              className="border-t border-white/5"
            >
              <td
                ref={(el) => {
                  lapRefs.current[i] = el;
                }}
                className="py-1 pr-2 text-zinc-400"
              />
              <td
                ref={(el) => {
                  driverRefs.current[i] = el;
                }}
                className="max-w-[7rem] truncate py-1 pr-2 text-zinc-300"
              />
              <td
                ref={(el) => {
                  timeRefs.current[i] = el;
                }}
                className="py-1 pr-2 text-white"
              />
              <td
                ref={(el) => {
                  gapRefs.current[i] = el;
                }}
                className="py-1 text-right text-zinc-500"
              />
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}

function gapClass(ms: number | null): string {
  const base = "py-1 text-right tabular-nums";
  if (ms === null || ms === 0) return `${base} text-zinc-500`;
  if (ms < 0) return `${base} text-emerald-400`;
  return `${base} text-red-400`;
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
      hidden
      data-active="0"
      className="inline-flex h-5 items-center rounded-full border border-amber-400 bg-amber-400/20 px-2 text-[10px] font-medium tracking-wide text-amber-200"
    >
      {label}
    </span>
  );
}
