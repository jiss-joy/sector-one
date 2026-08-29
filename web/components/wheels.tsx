"use client";

import { useRef, type Ref } from "react";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { useRequestAnimationFrame } from "@/hooks/use-raf";
import { isLocked, isSlipping } from "@/lib/format";
import { getLatest } from "@/lib/store";

const LOAD_MAX = 8000;

function paintCorner(
  load: number,
  slip: number,
  wheel: number,
  speed: number,
  bar: HTMLDivElement | null,
  flag: HTMLSpanElement | null,
  loadEl: HTMLSpanElement | null,
) {
  if (bar) bar.style.height = `${Math.min(100, (load / LOAD_MAX) * 100)}%`;
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

  useRequestAnimationFrame(() => {
    const f = getLatest();
    if (!f) return;
    const load = f.load_n ?? [0, 0, 0, 0];
    const slip = f.slip_ratio ?? [0, 0, 0, 0];
    const wh = f.wheel_rad_s ?? [0, 0, 0, 0];
    paintCorner(load[0], slip[0], wh[0], f.speed_kmh, barFL.current, flagFL.current, loadFL.current);
    paintCorner(load[1], slip[1], wh[1], f.speed_kmh, barFR.current, flagFR.current, loadFR.current);
    paintCorner(load[2], slip[2], wh[2], f.speed_kmh, barRL.current, flagRL.current, loadRL.current);
    paintCorner(load[3], slip[3], wh[3], f.speed_kmh, barRR.current, flagRR.current, loadRR.current);
  });

  return (
    <Card size="sm">
      <CardHeader>
        <CardTitle className="text-muted-foreground">
          Load{" "}
          <span className="font-sans text-[10px] font-normal">not tyre temp</span>
        </CardTitle>
      </CardHeader>
      <CardContent className="grid grid-cols-2 gap-3">
        <Corner name="FL" barRef={barFL} flagRef={flagFL} loadRef={loadFL} />
        <Corner name="FR" barRef={barFR} flagRef={flagFR} loadRef={loadFR} />
        <Corner name="RL" barRef={barRL} flagRef={flagRL} loadRef={loadRL} />
        <Corner name="RR" barRef={barRR} flagRef={flagRR} loadRef={loadRR} />
      </CardContent>
    </Card>
  );
}

function Corner({
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
    <div className="space-y-1">
      <div className="flex items-center justify-between">
        <span className="text-[10px] uppercase text-muted-foreground">{name}</span>
        <span ref={flagRef} className="text-transparent">
          —
        </span>
      </div>
      <div className="relative h-16 overflow-hidden rounded-sm bg-muted">
        <div ref={barRef} className="absolute inset-x-0 bottom-0 bg-primary/80" />
      </div>
      <span
        ref={loadRef}
        className="block font-mono text-[10px] tabular-nums text-muted-foreground"
      >
        0
      </span>
    </div>
  );
}
