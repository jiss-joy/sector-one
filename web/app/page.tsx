"use client";

import { ConnectEngine } from "@/components/connect-engine";
import { DashHeader } from "@/components/dash-header";
import { FrictionCircle } from "@/components/friction-circle";
import { Laps } from "@/components/laps";
import { Pedals } from "@/components/pedals";
import { Cluster } from "@/components/cluster";
import { Traces } from "@/components/traces";
import { Wheels } from "@/components/wheels";
import { StatusBoxes } from "@/components/status-boxes";

export default function Home() {
  return (
    <div className="mx-auto flex min-h-screen w-full max-w-[1600px] flex-col gap-4 p-4 antialiased selection:bg-primary/20">
      <div className="flex items-center justify-between gap-4 border-b border-white/5 pb-4">
        <DashHeader />
        <StatusBoxes />
        <ConnectEngine />
      </div>

      {/* Main Integrated Dash Layout */}
      <div className="grid flex-1 grid-cols-1 gap-6 lg:grid-cols-[1fr_2fr_1fr]">
        
        {/* Left Column: Technical Status */}
        <div className="flex flex-col justify-between gap-6 py-4">
          <div className="space-y-4">
            <FrictionCircle />
          </div>
          <Wheels />
        </div>

        {/* Center Column: Primary Driving Cluster */}
        <div className="flex flex-col items-center justify-center gap-8 py-4">
          <Cluster />
          <div className="w-full max-w-2xl px-8">
            <Pedals />
          </div>
        </div>

        {/* Right Column: Performance & Timing */}
        <div className="flex flex-col justify-between gap-6 py-4">
          <div className="space-y-6">
            <Laps />
          </div>
          <div className="mt-auto h-48 rounded border border-white/5 bg-zinc-900/20 p-2">
            <Traces />
          </div>
        </div>
      </div>
    </div>
  );
}
