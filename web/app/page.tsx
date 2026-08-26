"use client";

import { ConnectEngine } from "@/components/connect-engine";
import { DashHeader } from "@/components/dash-header";
import { FrictionCircle } from "@/components/friction-circle";
import { Laps } from "@/components/laps";
import { Pedals } from "@/components/pedals";
import { Readouts } from "@/components/readouts";
import { Traces } from "@/components/traces";
import { Wheels } from "@/components/wheels";

export default function Home() {
  return (
    <div className="mx-auto flex w-full max-w-6xl flex-1 flex-col gap-3 p-4">
      <ConnectEngine />
      <DashHeader />
      <div className="grid gap-3 lg:grid-cols-[1.4fr_1fr]">
        <div className="space-y-3">
          <Readouts />
          <Traces />
        </div>
        <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-1">
          <Pedals />
          <Laps />
        </div>
      </div>
      <div className="grid gap-3 md:grid-cols-2">
        <FrictionCircle />
        <Wheels />
      </div>
    </div>
  );
}
