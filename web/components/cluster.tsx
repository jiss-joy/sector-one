"use client";

import { useRef } from "react";
import { useHeartbeat } from "@/hooks/use-heartbeat";
import { clamp01, getGearLabel } from "@/lib/format";
import { getLatestFrame } from "@/lib/store";

const LED_COUNT = 20;

export function Cluster() {
  const gearRef = useRef<HTMLSpanElement>(null);
  const speedRef = useRef<HTMLSpanElement>(null);
  const rpmNumRef = useRef<HTMLSpanElement>(null);
  const carRef = useRef<HTMLSpanElement>(null);
  const ledRefs = useRef<(HTMLDivElement | null)[]>([]);
  const limiterRef = useRef<HTMLDivElement>(null);

  useHeartbeat(() => {
    const frame = getLatestFrame();
    if (!frame) return;

    if (gearRef.current) {
      gearRef.current.textContent = getGearLabel(frame.gear);
    }
    if (speedRef.current) {
      speedRef.current.textContent = frame.speed_kmh.toFixed(0);
    }
    if (rpmNumRef.current) {
      rpmNumRef.current.textContent = frame.rpm.toFixed(0);
    }
    if (carRef.current && frame.car) {
      carRef.current.textContent = frame.car.replaceAll("_", " ");
    }

    // RPM LEDs
    const rpmPercentage = clamp01(frame.rpm / frame.max_rpm);
    const activeLeds = Math.floor(rpmPercentage * LED_COUNT);
    const isAtLimit = rpmPercentage >= 0.98;

    ledRefs.current.forEach((el, i) => {
      if (!el) return;
      if (i < activeLeds) {
        const base = "h-4 w-6 rounded-sm transition-all duration-75";
        let style = "";

        if (i < 4) {
          style = "bg-blue-500 shadow-[0_0_10px_var(--color-blue-500)]";
        } else if (i < 12) {
          style = "bg-emerald-500 shadow-[0_0_10px_var(--color-emerald-500)]";
        } else if (i < 18) {
          style = "bg-yellow-400 shadow-[0_0_10px_var(--color-yellow-400)]";
        } else {
          style = "bg-red-500 shadow-[0_0_15px_var(--color-red-500)]";
        }

        el.className = `${base} ${style}`;
      } else {
        el.className = "h-4 w-6 rounded-sm bg-zinc-800 shadow-none animate-none";
      }
    });

    if (limiterRef.current) {
      if (isAtLimit) {
        limiterRef.current.className = "h-4 w-4 rounded-full bg-red-600 shadow-[0_0_20px_var(--color-red-600)] animate-rev-limiter opacity-100";
      } else {
        limiterRef.current.className = "h-4 w-4 rounded-full bg-zinc-900 opacity-20 animate-none";
      }
    }
  });

  return (
    <div className="relative flex flex-col items-center gap-12 w-full max-w-2xl">
      <span
        ref={carRef}
        className="text-[14px] font-bold uppercase tracking-[0.25em] text-zinc-500"
      />
      {/* RPM LED Bar Section */}
      <div className="flex flex-col items-center gap-2">
        <div className="flex items-center gap-6">
          <div ref={limiterRef} className="h-4 w-4 rounded-full bg-zinc-900 opacity-20 transition-all duration-200" />
          <div className="flex items-baseline gap-1">
            <span ref={rpmNumRef} className="font-mono text-3xl font-bold text-white/80 tabular-nums">0</span>
            <span className="text-[10px] font-bold text-white/20 uppercase tracking-[0.2em]">RPM</span>
          </div>
          <div className="h-4 w-4 rounded-full bg-zinc-900 opacity-20 transition-all duration-200" />
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
      <div className="flex items-center justify-center gap-12 w-full">
        <div className="flex flex-col items-center">
          <span
            ref={gearRef}
            className="font-mono text-[180px] leading-none font-bold text-white tabular-nums drop-shadow-2xl"
          >
            N
          </span>
          <span className="text-[10px] font-bold tracking-[0.3em] text-white/20 uppercase mt-2">Gear</span>
        </div>

        <div className="h-32 w-px bg-white/5 mx-4" />

        <div className="flex flex-col items-center">
          <span
            ref={speedRef}
            className="font-mono text-8xl font-bold text-white/90 tabular-nums leading-none"
          >
            0
          </span>
          <span className="text-sm font-medium tracking-[0.2em] text-white/40 uppercase mt-4">
            KM/H
          </span>
        </div>
      </div>
    </div>
  );
}
