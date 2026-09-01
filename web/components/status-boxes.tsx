"use client";

import { useRef } from "react";
import { useHeartbeat } from "@/hooks/use-heartbeat";
import { getLatest } from "@/lib/store";

export function StatusBoxes() {
  const absRef = useRef<HTMLDivElement>(null);
  const tcRef = useRef<HTMLDivElement>(null);

  useHeartbeat(() => {
    const f = getLatest();
    if (!f) return;

    if (absRef.current) {
      absRef.current.dataset.active = f.abs_in_action ? "1" : "0";
    }
    if (tcRef.current) {
      tcRef.current.dataset.active = f.tc_in_action ? "1" : "0";
    }
  });

  return (
    <div className="flex gap-2">
      <Box label="TC" color="bg-blue-600" ref={tcRef} />
      <Box label="ABS" color="bg-emerald-600" ref={absRef} />
    </div>
  );
}

function Box({ label, color, ref }: { label: string; color: string; ref: React.RefObject<HTMLDivElement> }) {
  return (
    <div
      ref={ref}
      className={`flex h-14 w-14 flex-col items-center justify-center rounded border border-white/10 ${color} opacity-20 data-[active=1]:opacity-100 transition-opacity duration-75`}
    >
      <span className="text-[10px] font-bold text-white/80">{label}</span>
      <span className="text-xl font-bold text-white">1</span>
    </div>
  );
}
