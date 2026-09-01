"use client";

import { useRef } from "react";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { useHeartbeat } from "@/hooks/use-heartbeat";
import { clamp01 } from "@/lib/format";
import { forEachRing } from "@/lib/store";

function resize(canvas: HTMLCanvasElement) {
  const dpr = window.devicePixelRatio || 1;
  const w = canvas.clientWidth;
  const h = canvas.clientHeight;
  const bw = Math.round(w * dpr);
  const bh = Math.round(h * dpr);
  if (canvas.width !== bw || canvas.height !== bh) {
    canvas.width = bw;
    canvas.height = bh;
  }
  return { w, h, dpr };
}

export function Traces() {
  const canvasRef = useRef<HTMLCanvasElement>(null);

  useHeartbeat(() => {
    const canvas = canvasRef.current;
    if (!canvas) return;
    const ctx = canvas.getContext("2d");
    if (!ctx) return;
    const { w, h, dpr } = resize(canvas);
    ctx.setTransform(dpr, 0, 0, dpr, 0, 0);
    ctx.clearRect(0, 0, w, h);

    let n = 0;
    forEachRing(() => {
      n++;
    });
    if (n < 2) return;

    const series: { stroke: string; pick: (f: { throttle: number; brake: number; clutch: number; steer: number }) => number }[] =
      [
        { stroke: "#22c55e", pick: (f) => clamp01(f.throttle) },
        { stroke: "#ef4444", pick: (f) => clamp01(f.brake) },
        { stroke: "#38bdf8", pick: (f) => clamp01(f.clutch) },
        { stroke: "#f59e0b", pick: (f) => clamp01((f.steer + 1) / 2) },
      ];

    for (const ser of series) {
      ctx.beginPath();
      ctx.strokeStyle = ser.stroke;
      ctx.lineWidth = 1.25;
      let i = 0;
      forEachRing((f) => {
        const x = (i / (n - 1)) * w;
        const y = (1 - ser.pick(f)) * h;
        if (i === 0) ctx.moveTo(x, y);
        else ctx.lineTo(x, y);
        i++;
      });
      ctx.stroke();
    }
  });

  return (
    <div className="flex h-full flex-col gap-2 p-2">
      <div className="flex items-center justify-between px-1">
        <span className="text-[10px] font-bold tracking-widest text-zinc-500 uppercase">Input History</span>
        <div className="flex gap-2">
          <div className="flex items-center gap-1">
            <div className="h-1.5 w-1.5 rounded-full bg-[#22c55e]" />
            <span className="text-[8px] font-bold text-zinc-500">THR</span>
          </div>
          <div className="flex items-center gap-1">
            <div className="h-1.5 w-1.5 rounded-full bg-[#ef4444]" />
            <span className="text-[8px] font-bold text-zinc-500">BRK</span>
          </div>
        </div>
      </div>
      <canvas ref={canvasRef} className="flex-1 w-full" />
    </div>
  );
}
