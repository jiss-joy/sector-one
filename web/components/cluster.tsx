"use client";

import { useRef } from "react";
import { useHeartbeat } from "@/hooks/use-heartbeat";
import { clamp01, gearLabel } from "@/lib/format";
import { getLatest } from "@/lib/store";

const LED_COUNT = 20;

export function Cluster() {
  const gearRef = useRef<HTMLSpanElement>(null);
  const speedRef = useRef<HTMLSpanElement>(null);
  const rpmNumRef = useRef<HTMLSpanElement>(null);
  const ledRefs = useRef<(HTMLDivElement | null)[]>([]);

  useHeartbeat(() => {
    const f = getLatest();
    if (!f) return;

    if (gearRef.current) {
      gearRef.current.textContent = gearLabel(f.gear);
    }
    if (speedRef.current) {
      speedRef.current.textContent = f.speed_kmh.toFixed(0);
    }
    if (rpmNumRef.current) {
      rpmNumRef.current.textContent = f.rpm.toFixed(0);
    }

    // RPM LEDs
    const rpmPct = clamp01(f.rpm / f.max_rpm);
    const activeLeds = Math.floor(rpmPct * LED_COUNT);

    ledRefs.current.forEach((el, i) => {
      if (!el) return;
      if (i < activeLeds) {
        // Color segments like the screenshot
        if (i < 4) el.className = "h-4 w-6 rounded-sm bg-blue-500 shadow-[0_0_8px_rgba(59,130,246,0.5)]";
        else if (i < 8) el.className = "h-4 w-6 rounded-sm bg-emerald-500 shadow-[0_0_8px_rgba(16,185,129,0.5)]";
        else if (i < 14) el.className = "h-4 w-6 rounded-sm bg-yellow-400 shadow-[0_0_8px_rgba(250,204,21,0.5)]";
        else el.className = "h-4 w-6 rounded-sm bg-red-500 shadow-[0_0_12px_rgba(239,68,68,0.8)] animate-pulse";
      } else {
        el.className = "h-4 w-6 rounded-sm bg-zinc-800";
      }
    });
  });

  return (
    <div className="flex flex-col items-center gap-12">
      {/* RPM LED Bar Section */}
      <div className="flex flex-col items-center gap-2">
        <div className="flex items-baseline gap-1">
          <span ref={rpmNumRef} className="font-mono text-3xl font-bold text-white/80 tabular-nums">0</span>
          <span className="text-[10px] font-bold text-white/20 uppercase tracking-[0.2em]">RPM</span>
        </div>
        <div className="flex gap-1.5 pt-0">
          {Array.from({ length: LED_COUNT }).map((_, i) => (
            <div
              key={i}
              ref={(el) => {
                ledRefs.current[i] = el;
              }}
              className="h-4 w-6 rounded-sm bg-zinc-800 transition-colors duration-75"
            />
          ))}
        </div>
      </div>

      {/* Main Cluster */}
      <div className="flex flex-col items-center justify-center py-0">
        <span
          ref={gearRef}
          className="font-mono text-[200px] leading-none font-bold text-white tabular-nums drop-shadow-2xl"
        >
          N
        </span>
        <div className="flex flex-col items-center mt-12">
          <span
            ref={speedRef}
            className="font-mono text-7xl font-bold text-white/90 tabular-nums"
          >
            0
          </span>
          <span className="text-sm font-medium tracking-[0.2em] text-white/40 uppercase">
            KM/H
          </span>
        </div>
      </div>
    </div>
  );
}
