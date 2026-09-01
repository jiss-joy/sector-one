"use client";

import { useRef } from "react";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { useHeartbeat } from "@/hooks/use-heartbeat";
import { TELEMETRY_CONFIG } from "@/lib/constants";
import { clamp01, gearLabel } from "@/lib/format";
import { getLatest } from "@/lib/store";

export function Readouts() {
  const speedRef = useRef<HTMLSpanElement>(null);
  const gearRef = useRef<HTMLSpanElement>(null);
  const rpmRef = useRef<HTMLSpanElement>(null);
  const rpmBarRef = useRef<HTMLDivElement>(null);
  const steerRef = useRef<HTMLDivElement>(null);

  useHeartbeat(() => {
    const f = getLatest();
    if (!f) return;
    if (speedRef.current) {
      speedRef.current.textContent = f.speed_kmh.toFixed(0);
    }
    if (gearRef.current) {
      gearRef.current.textContent = gearLabel(f.gear);
    }
    if (rpmRef.current) {
      rpmRef.current.textContent = f.rpm.toFixed(0);
    }
    if (rpmBarRef.current) {
      rpmBarRef.current.style.width = `${clamp01(f.rpm / TELEMETRY_CONFIG.RPM_MAX) * 100}%`;
    }
    if (steerRef.current) {
      const x = (clamp01((f.steer + 1) / 2) * 100).toFixed(1);
      steerRef.current.style.left = `calc(${x}% - 4px)`;
    }
  });

  return (
    <div className="grid grid-cols-3 gap-3">
      <Card size="sm">
        <CardHeader>
          <CardTitle className="text-muted-foreground">Speed</CardTitle>
        </CardHeader>
        <CardContent>
          <p className="font-mono text-4xl font-medium tabular-nums tracking-tight">
            <span ref={speedRef}>—</span>
            <span className="ml-1 text-sm text-muted-foreground">km/h</span>
          </p>
        </CardContent>
      </Card>
      <Card size="sm">
        <CardHeader>
          <CardTitle className="text-muted-foreground">Gear</CardTitle>
        </CardHeader>
        <CardContent>
          <p className="font-mono text-5xl font-medium tabular-nums">
            <span ref={gearRef}>—</span>
          </p>
        </CardContent>
      </Card>
      <Card size="sm">
        <CardHeader>
          <CardTitle className="text-muted-foreground">RPM</CardTitle>
        </CardHeader>
        <CardContent className="space-y-2">
          <p className="font-mono text-4xl font-medium tabular-nums">
            <span ref={rpmRef}>—</span>
          </p>
          <div className="h-1.5 overflow-hidden rounded-full bg-muted">
            <div ref={rpmBarRef} className="h-full w-0 bg-primary" />
          </div>
        </CardContent>
      </Card>
      <Card size="sm" className="col-span-3">
        <CardHeader>
          <CardTitle className="text-muted-foreground">Steer</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="relative h-2 rounded-full bg-muted">
            <div className="absolute top-1/2 left-1/2 h-3 w-px -translate-x-1/2 -translate-y-1/2 bg-foreground/40" />
            <div
              ref={steerRef}
              className="absolute top-1/2 left-1/2 size-2 -translate-y-1/2 rounded-full bg-primary"
            />
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
