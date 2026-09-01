"use client";

import { useRef } from "react";

import { useHeartbeat } from "@/hooks/use-heartbeat";
import { TELEMETRY_CONFIG } from "@/lib/constants";
import { lastN } from "@/lib/store";
import type { Frame } from "@/lib/types";

const scratch: Frame[] = new Array(TELEMETRY_CONFIG.FRICTION_CIRCLE_TAIL_COUNT);

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

export function FrictionCircle() {
  const canvasRef = useRef<HTMLCanvasElement>(null);

  useHeartbeat(() => {
    const canvas = canvasRef.current;
    if (!canvas) return;
    const ctx = canvas.getContext("2d");
    if (!ctx) return;
    const { w, h, dpr } = resize(canvas);
    
    // Clear entire raw buffer before transform
    ctx.clearRect(0, 0, canvas.width, canvas.height);
    
    ctx.setTransform(dpr, 0, 0, dpr, 0, 0);

    const cx = w / 2;
    const cy = h / 2;
    const r = Math.min(w, h) * 0.42;

    ctx.strokeStyle = "rgba(255,255,255,0.12)";
    ctx.lineWidth = 1;
    for (const g of [1, 2]) {
      ctx.beginPath();
      ctx.arc(cx, cy, r * (g / TELEMETRY_CONFIG.FRICTION_CIRCLE_G_SCALE), 0, Math.PI * 2);
      ctx.stroke();
    }
    ctx.beginPath();
    ctx.moveTo(cx - r, cy);
    ctx.lineTo(cx + r, cy);
    ctx.moveTo(cx, cy - r);
    ctx.lineTo(cx, cy + r);
    ctx.stroke();

    const n = lastN(TELEMETRY_CONFIG.FRICTION_CIRCLE_TAIL_COUNT, scratch);
    if (n === 0) return;

    const plot = (lat: number, lon: number) => {
      const x = cx + (lat / TELEMETRY_CONFIG.FRICTION_CIRCLE_G_SCALE) * r;
      // Flip longitudinal axis if calibrated
      const gLong = TELEMETRY_CONFIG.FRICTION_CIRCLE_FLIP_LONG ? -lon : lon;
      const y = cy - (gLong / TELEMETRY_CONFIG.FRICTION_CIRCLE_G_SCALE) * r;
      return { x, y };
    };

    ctx.beginPath();
    ctx.strokeStyle = "rgba(74,222,128,0.45)";
    ctx.lineWidth = 1;
    ctx.lineJoin = "round";
    for (let i = 0; i < n; i++) {
      const p = plot(scratch[i].g_lat, scratch[i].g_long);
      if (i === 0) ctx.moveTo(p.x, p.y);
      else ctx.lineTo(p.x, p.y);
    }
    ctx.stroke();

    const cur = scratch[n - 1];
    const p = plot(cur.g_lat, cur.g_long);
    ctx.fillStyle = "#4ade80";
    ctx.beginPath();
    ctx.arc(p.x, p.y, 4, 0, Math.PI * 2);
    ctx.fill();
  });

  return (
    <div className="relative flex flex-col gap-2 rounded border border-white/5 bg-zinc-900/40 p-4">
      <span className="text-[10px] font-bold tracking-widest text-zinc-500 uppercase">G-Force</span>
      <div className="flex aspect-square w-full items-center justify-center p-2">
        <canvas ref={canvasRef} className="h-full w-full" />
      </div>
    </div>
  );
}
