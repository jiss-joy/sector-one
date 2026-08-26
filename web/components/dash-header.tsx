"use client";

import { useSyncExternalStore } from "react";

import { Badge } from "@/components/ui/badge";
import { AnimatedShinyText } from "@/components/ui/animated-shiny-text";
import { BorderBeam } from "@/components/ui/border-beam";
import { getConn, subscribeConn } from "@/lib/store";

export function DashHeader() {
  const conn = useSyncExternalStore(subscribeConn, getConn, getConn);
  const live = conn.state === "open";

  return (
    <header className="relative overflow-hidden rounded-xl border bg-card px-4 py-3">
      {live ? (
        <BorderBeam
          size={80}
          duration={8}
          colorFrom="#22c55e"
          colorTo="#4ade80"
          borderWidth={1}
        />
      ) : null}
      <div className="flex flex-wrap items-center justify-between gap-3">
        <div>
          <p className="font-heading text-sm tracking-wide text-muted-foreground">
            SECTOR ONE
          </p>
          <h1 className="font-heading text-lg font-medium">Telemetry</h1>
        </div>
        <div className="flex items-center gap-2">
          <Badge variant={live ? "default" : "outline"}>
            {conn.state === "open"
              ? "connected"
              : conn.state === "error"
                ? "no engine"
                : "idle"}
          </Badge>
          {conn.source ? (
            <Badge variant="secondary">{conn.source}</Badge>
          ) : (
            <AnimatedShinyText className="mx-0 text-xs">
              waiting for :8080
            </AnimatedShinyText>
          )}
        </div>
      </div>
    </header>
  );
}
