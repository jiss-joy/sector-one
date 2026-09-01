import { TELEMETRY_CONFIG } from "./constants";

export function gearLabel(gear: number): string {
  if (gear === 0) return "R";
  if (gear === 1) return "N";
  return String(gear - 1);
}

export function formatLap(ms: number): string {
  if (ms <= 0) return "--:--.---";
  const m = Math.floor(ms / 60000);
  const s = Math.floor((ms % 60000) / 1000);
  const frac = ms % 1000;
  return `${String(m).padStart(2, "0")}:${String(s).padStart(2, "0")}.${String(frac).padStart(3, "0")}`;
}

export function clamp01(v: number): number {
  if (v < 0) return 0;
  if (v > 1) return 1;
  return v;
}

export const CORNERS = ["FL", "FR", "RL", "RR"] as const;

export function isLocked(speedKmh: number, wheelRadS: number): boolean {
  return (
    speedKmh > TELEMETRY_CONFIG.THRESHOLD_LOCK_MIN_SPEED &&
    Math.abs(wheelRadS) < TELEMETRY_CONFIG.THRESHOLD_LOCK_WHEEL_RAD_S
  );
}

export function isSlipping(slipRatio: number): boolean {
  return Math.abs(slipRatio) > TELEMETRY_CONFIG.THRESHOLD_SLIP_RATIO;
}
