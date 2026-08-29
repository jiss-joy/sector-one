"use client";

import { useRef, type Ref } from "react";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { useRequestAnimationFrame } from "@/hooks/use-raf";
import { clamp01 } from "@/lib/format";
import { getLatest } from "@/lib/store";

export function Pedals() {
  const thrBar = useRef<HTMLDivElement>(null);
  const brkBar = useRef<HTMLDivElement>(null);
  const cluBar = useRef<HTMLDivElement>(null);
  const thrLbl = useRef<HTMLSpanElement>(null);
  const brkLbl = useRef<HTMLSpanElement>(null);
  const cluLbl = useRef<HTMLSpanElement>(null);

  useRequestAnimationFrame(() => {
    const f = getLatest();
    if (!f) return;
    const t = clamp01(f.throttle);
    const b = clamp01(f.brake);
    const c = clamp01(f.clutch);
    if (thrBar.current) thrBar.current.style.height = `${t * 100}%`;
    if (brkBar.current) brkBar.current.style.height = `${b * 100}%`;
    if (cluBar.current) cluBar.current.style.height = `${c * 100}%`;
    if (thrLbl.current) thrLbl.current.textContent = t.toFixed(2);
    if (brkLbl.current) brkLbl.current.textContent = b.toFixed(2);
    if (cluLbl.current) cluLbl.current.textContent = c.toFixed(2);
  });

  return (
    <Card size="sm" className="h-full">
      <CardHeader>
        <CardTitle className="text-muted-foreground">Pedals</CardTitle>
      </CardHeader>
      <CardContent className="flex h-40 items-end justify-around gap-4">
        <PedalCol barRef={thrBar} labelRef={thrLbl} name="throttle" color="bg-emerald-500" />
        <PedalCol barRef={brkBar} labelRef={brkLbl} name="brake" color="bg-red-500" />
        <PedalCol barRef={cluBar} labelRef={cluLbl} name="clutch" color="bg-sky-500" />
      </CardContent>
    </Card>
  );
}

function PedalCol({
  barRef,
  labelRef,
  name,
  color,
}: {
  barRef: Ref<HTMLDivElement>;
  labelRef: Ref<HTMLSpanElement>;
  name: string;
  color: string;
}) {
  return (
    <div className="flex h-full flex-1 flex-col items-center gap-1">
      <span
        ref={labelRef}
        className="font-mono text-[10px] tabular-nums text-muted-foreground"
      >
        0.00
      </span>
      <div className="relative w-full flex-1 overflow-hidden rounded-sm bg-muted">
        <div ref={barRef} className={`absolute inset-x-0 bottom-0 ${color}`} />
      </div>
      <span className="text-[10px] uppercase tracking-wide text-muted-foreground">
        {name}
      </span>
    </div>
  );
}
