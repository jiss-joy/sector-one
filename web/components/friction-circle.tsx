"use client";

import { useRef } from "react";

import { useHeartbeat } from "@/hooks/use-heartbeat";
import { lastN } from "@/lib/store";
import { THRESHOLDS } from "@/lib/thresholds";
import type { Frame } from "@/lib/types";

const scratch: Frame[] = new Array(THRESHOLDS.FRICTION_CIRCLE_TAIL_COUNT);

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
  const readoutRef = useRef<HTMLSpanElement>(null);

  useHeartbeat(() => {
    const canvas = canvasRef.current;
    if (!canvas) return;
    const ctx = canvas.getContext("2d");
    if (!ctx) return;
    const { w, h, dpr } = resize(canvas);
    
    ctx.clearRect(0, 0, canvas.width, canvas.height);
    ctx.setTransform(dpr, 0, 0, dpr, 0, 0);

    const cx = w / 2;
    const cy = h / 2;
    const r = Math.min(w, h) * 0.4;

    // Grid (Rings)
    ctx.strokeStyle = "rgba(255,255,255,0.08)";
    ctx.lineWidth = 1;
    ctx.textAlign = "center";
    ctx.textBaseline = "middle";
    ctx.font = "8px monospace";
    ctx.fillStyle = "rgba(255,255,255,0.2)";

    for (const g of [1, 2]) {
      const ringR = r * (g / THRESHOLDS.FRICTION_CIRCLE_G_SCALE);
      ctx.beginPath();
      ctx.arc(cx, cy, ringR, 0, Math.PI * 2);
      ctx.stroke();
      
      // Ring Label
      ctx.fillText(`${g}G`, cx, cy - ringR - 6);
    }

    // Grid (Crosshair)
    ctx.beginPath();
    ctx.moveTo(cx - r, cy); ctx.lineTo(cx + r, cy);
    ctx.moveTo(cx, cy - r); ctx.lineTo(cx, cy + r);
    ctx.stroke();

    // Axis Labels
    ctx.font = "9px font-black tracking-tighter uppercase italic";
    ctx.fillStyle = "rgba(255,255,255,0.3)";
    ctx.fillText("Brake", cx, cy - r - 16);
    ctx.fillText("Accel", cx, cy + r + 16);
    ctx.textAlign = "left";
    ctx.fillText("Left", cx - r - 32, cy);
    ctx.textAlign = "right";
    ctx.fillText("Right", cx + r + 32, cy);

    const n = lastN(THRESHOLDS.FRICTION_CIRCLE_TAIL_COUNT, scratch);
    if (n === 0) return;

    const plot = (lat: number, lon: number) => {
      const x = cx + (lat / THRESHOLDS.FRICTION_CIRCLE_G_SCALE) * r;
      // Flip longitudinal axis if calibrated
      const gLong = THRESHOLDS.FRICTION_CIRCLE_FLIP_LONG ? -lon : lon;
      const y = cy - (gLong / THRESHOLDS.FRICTION_CIRCLE_G_SCALE) * r;
      return { x, y };
    };

    // Draw Tail (Fade)
    for (let i = 0; i < n - 1; i++) {
      const alpha = (i / n) * 0.4;
      ctx.beginPath();
      ctx.strokeStyle = `rgba(74,222,128,${alpha})`;
      ctx.lineWidth = 1.5;
      const p1 = plot(scratch[i].g_lat, scratch[i].g_long);
      const p2 = plot(scratch[i + 1].g_lat, scratch[i + 1].g_long);
      ctx.moveTo(p1.x, p1.y);
      ctx.lineTo(p2.x, p2.y);
      ctx.stroke();
    }

    // Live Dot
    const cur = scratch[n - 1];
    const p = plot(cur.g_lat, cur.g_long);
    
    // Glow
    ctx.shadowBlur = 10;
    ctx.shadowColor = "#4ade80";
    ctx.fillStyle = "#4ade80";
    ctx.beginPath();
    ctx.arc(p.x, p.y, 4, 0, Math.PI * 2);
    ctx.fill();
    ctx.shadowBlur = 0;

    // Current G Readout
    if (readoutRef.current) {
      const gTotal = Math.sqrt(cur.g_lat ** 2 + cur.g_long ** 2);
      readoutRef.current.textContent = `${gTotal.toFixed(2)} G`;
    }
  });

  return (
    <div className="relative flex flex-col items-center gap-2 rounded border border-white/5 bg-zinc-900/40 p-4">
      <div className="flex w-full items-center justify-between">
        <span className="text-[10px] font-bold tracking-widest text-zinc-500 uppercase italic">Friction Circle</span>
        <span ref={readoutRef} className="font-mono text-xs font-bold text-emerald-500 tabular-nums">0.00 G</span>
      </div>
      <div className="flex aspect-square w-full items-center justify-center p-4">
        <canvas ref={canvasRef} className="h-full w-full" />
      </div>
    </div>
  );
}
