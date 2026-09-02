"use client";

import { Activity, Signal, SignalLow, AlertCircle } from "lucide-react";
import { useSyncExternalStore } from "react";

import { getConnection, subscribe } from "@/lib/store";

export function Header() {
  const conn = useSyncExternalStore(subscribe, getConnection, getConnection);
  const isLive = conn.state === "online";
  const hasError = conn.state === "error";

  return (
    <div className="flex w-full items-center justify-between">
      <div className="flex items-center gap-4">
        <div className="flex size-8 items-center justify-center rounded-lg bg-primary/10 text-primary border border-primary/20">
          <Activity size={16} className={isLive ? "animate-pulse" : ""} />
        </div>
        <div>
          <div className="flex items-center gap-2">
            <h1 className="text-lg font-black tracking-tighter text-white uppercase italic">
              Sector <span className="text-primary">One</span>
            </h1>
          </div>
        </div>
      </div>

      <div className="flex items-center gap-6">
        {/* Source Info */}
        {conn.source && (
          <div className="flex flex-col items-end">
            <span className="text-[7px] font-bold text-zinc-600 uppercase tracking-[0.2em] leading-none mb-1">Active Machine</span>
            <span className="text-[10px] font-mono font-bold text-zinc-400 uppercase leading-none">{conn.source}</span>
          </div>
        )}

        {/* Connection Status Pill */}
        <div className={`flex items-center gap-2 rounded-full border px-2.5 py-1 transition-all duration-300 ${
          isLive 
            ? "border-emerald-500/20 bg-emerald-500/5 text-emerald-500/80 shadow-[0_0_10px_rgba(16,185,129,0.05)]" 
            : hasError 
              ? "border-red-500/20 bg-red-500/5 text-red-500/80"
              : "border-zinc-800 bg-zinc-900/50 text-zinc-600"
        }`}>
          {isLive ? (
            <Signal size={12} />
          ) : hasError ? (
            <AlertCircle size={12} />
          ) : (
            <SignalLow size={12} />
          )}
          <span className="text-[8px] font-black uppercase tracking-[0.2em]">
            {isLive ? "Live" : hasError ? "Error" : "Offline"}
          </span>
          {isLive && (
            <span className="flex h-1 w-1 relative">
              <span className="absolute inline-flex h-full w-full animate-ping rounded-full bg-emerald-400 opacity-75"></span>
              <span className="relative inline-flex h-1 w-1 rounded-full bg-emerald-500"></span>
            </span>
          )}
        </div>
      </div>
    </div>
  );
}
