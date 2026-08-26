"use client";

import { useRef } from "react";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { useRaf } from "@/hooks/use-raf";
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

  useRaf(() => {
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
    <Card size="sm">
      <CardHeader>
        <CardTitle className="text-muted-foreground">
          Traces{" "}
          <span className="font-sans text-[10px] font-normal">
            thr · brake · clutch · steer
          </span>
        </CardTitle>
      </CardHeader>
      <CardContent>
        <canvas ref={canvasRef} className="h-36 w-full" />
      </CardContent>
    </Card>
  );
}
