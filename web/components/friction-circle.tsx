"use client";

import { useRef } from "react";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { useRequestAnimationFrame } from "@/hooks/use-raf";
import { lastN } from "@/lib/store";
import type { Frame } from "@/lib/types";

const TAIL = 180;
const scratch: Frame[] = new Array(TAIL);
const SCALE = 2.5; // G units to edge of circle

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

  useRequestAnimationFrame(() => {
    const canvas = canvasRef.current;
    if (!canvas) return;
    const ctx = canvas.getContext("2d");
    if (!ctx) return;
    const { w, h, dpr } = resize(canvas);
    ctx.setTransform(dpr, 0, 0, dpr, 0, 0);
    ctx.clearRect(0, 0, w, h);

    const cx = w / 2;
    const cy = h / 2;
    const r = Math.min(w, h) * 0.42;

    ctx.strokeStyle = "rgba(255,255,255,0.12)";
    ctx.lineWidth = 1;
    for (const g of [1, 2]) {
      ctx.beginPath();
      ctx.arc(cx, cy, r * (g / SCALE), 0, Math.PI * 2);
      ctx.stroke();
    }
    ctx.beginPath();
    ctx.moveTo(cx - r, cy);
    ctx.lineTo(cx + r, cy);
    ctx.moveTo(cx, cy - r);
    ctx.lineTo(cx, cy + r);
    ctx.stroke();

    const n = lastN(TAIL, scratch);
    if (n === 0) return;

    const plot = (lat: number, lon: number) => {
      const x = cx + (lat / SCALE) * r;
      // +g_long up. Flip this one line if a brake zone plots the wrong way.
      const y = cy - (lon / SCALE) * r;
      return { x, y };
    };

    ctx.beginPath();
    ctx.strokeStyle = "rgba(74,222,128,0.45)";
    ctx.lineWidth = 1;
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
    <Card size="sm" className="h-full">
      <CardHeader>
        <CardTitle className="text-muted-foreground">Friction circle</CardTitle>
      </CardHeader>
      <CardContent>
        <canvas ref={canvasRef} className="aspect-square w-full" />
      </CardContent>
    </Card>
  );
}
