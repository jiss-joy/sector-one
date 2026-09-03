"use client";

import { Ref, useRef } from "react";

import { useHeartbeat } from "@/hooks/use-heartbeat";
import { clamp01 } from "@/lib/format";
import { getLatestFrame } from "@/lib/store";

export function Pedals() {
  const throttleBar = useRef<HTMLDivElement>(null);
  const brakeBar = useRef<HTMLDivElement>(null);
  const clutchBar = useRef<HTMLDivElement>(null);
  const throttleLabel = useRef<HTMLSpanElement>(null);
  const brakeLabel = useRef<HTMLSpanElement>(null);
  const clutchLabel = useRef<HTMLSpanElement>(null);

  useHeartbeat(() => {
    const f = getLatestFrame();
    if (!f) return;
    const t = clamp01(f.throttle);
    const b = clamp01(f.brake);
    const c = clamp01(f.clutch);
    if (throttleBar.current) throttleBar.current.style.height = `${t * 100}%`;
    if (brakeBar.current) brakeBar.current.style.height = `${b * 100}%`;
    if (clutchBar.current) clutchBar.current.style.height = `${c * 100}%`;
    if (throttleLabel.current) throttleLabel.current.textContent = (t * 100).toFixed(0);
    if (brakeLabel.current) brakeLabel.current.textContent = (b * 100).toFixed(0);
    if (clutchLabel.current) clutchLabel.current.textContent = (c * 100).toFixed(0);
  });

  return (
    <div className="flex h-44 items-end justify-around gap-6">
      <PedalBar barRef={clutchBar} labelRef={clutchLabel} name="clutch" color="bg-sky-500" />
      <PedalBar barRef={brakeBar} labelRef={brakeLabel} name="brake" color="bg-red-500" />
      <PedalBar barRef={throttleBar} labelRef={throttleLabel} name="throttle" color="bg-emerald-500" />
    </div>
  );
}

function PedalBar({
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
