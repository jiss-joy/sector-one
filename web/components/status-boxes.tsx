"use client";

import { useRef } from "react";
import { useHeartbeat } from "@/hooks/use-heartbeat";
import { getLatestFrame } from "@/lib/store";

export function StatusBoxes({ className }: { className?: string }) {
  const absRef = useRef<HTMLDivElement>(null);
  const tcRef = useRef<HTMLDivElement>(null);

  useHeartbeat(() => {
    const f = getLatestFrame();
    if (!f) return;

    if (absRef.current) {
      absRef.current.dataset.enabled = f.abs_enabled ? "1" : "0";
      absRef.current.dataset.active = f.abs_in_action ? "1" : "0";
    }
    if (tcRef.current) {
      tcRef.current.dataset.enabled = f.tc_enabled ? "1" : "0";
      tcRef.current.dataset.active = f.tc_in_action ? "1" : "0";
    }
  });

  return (
    <div className={`flex gap-1.5 ${className}`}>
      <Box label="TC" color="bg-blue-600" ref={tcRef} />
      <Box label="ABS" color="bg-emerald-600" ref={absRef} />
    </div>
  );
}

function Box({ label, color, ref }: { label: string; color: string; ref: React.RefObject<HTMLDivElement | null> }) {
  return (
    <div
      ref={ref}
      data-enabled="0"
      data-active="0"
      className={`flex h-12 w-12 flex-col items-center justify-center rounded border border-white/5 ${color} opacity-10 data-[enabled=1]:opacity-50 data-[active=1]:opacity-100 transition-all duration-75 shadow-inner`}
    >
      <span className="text-[11px] font-black leading-none text-white/90">{label}</span>
    </div>
  );
}
