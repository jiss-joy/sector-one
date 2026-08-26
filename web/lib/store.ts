import type { Conn, Frame } from "@/lib/types";

const CAP = 600;

let latest: Frame | null = null;
const ring: Array<Frame | undefined> = new Array(CAP);
let len = 0;
let head = 0;

let conn: Conn = { state: "down", source: "", subscribers: 0 };

const connListeners = new Set<() => void>();

export function pushFrame(f: Frame) {
  ring[head] = f;
  head = (head + 1) % CAP;
  if (len < CAP) len++;
  latest = f;
  if (f.source && f.source !== conn.source) {
    conn = { ...conn, source: f.source };
    emitConn();
  }
}

export function getLatest(): Frame | null {
  return latest;
}

export function ringLength(): number {
  return len;
}

/** Oldest → newest. `fn` must not retain the frame past the call if you mutate later. */
export function forEachRing(fn: (f: Frame, i: number) => void) {
  const start = len < CAP ? 0 : head;
  for (let i = 0; i < len; i++) {
    const f = ring[(start + i) % CAP];
    if (f) fn(f, i);
  }
}

export function lastN(n: number, out: Frame[]): number {
  const take = Math.min(n, len);
  const start = (len < CAP ? 0 : head) + (len - take);
  for (let i = 0; i < take; i++) {
    out[i] = ring[(start + i) % CAP]!;
  }
  return take;
}

export function getConn(): Conn {
  return conn;
}

export function setConn(next: Partial<Conn>) {
  const merged: Conn = { ...conn, ...next };
  if (
    merged.state === conn.state &&
    merged.source === conn.source &&
    merged.subscribers === conn.subscribers
  ) {
    return;
  }
  conn = merged;
  emitConn();
}

export function subscribeConn(fn: () => void) {
  connListeners.add(fn);
  return () => {
    connListeners.delete(fn);
  };
}

function emitConn() {
  for (const l of connListeners) l();
}
