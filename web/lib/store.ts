import type { Connection, Frame } from "@/lib/types";

const RING_CAPACITY = 600;

let latestFrame: Frame | null = null;
const ringBuffer: Array<Frame | undefined> = new Array(RING_CAPACITY);
let length = 0;
let head = 0;

let connection: Connection = { state: "offline", source: "", subscribers: 0 };

const connectionListeners = new Set<() => void>();

export function pushFrame(frame: Frame) {
  ringBuffer[head] = frame;
  head = (head + 1) % RING_CAPACITY;
  if (length < RING_CAPACITY) length++;
  latestFrame = frame;
  if (frame.source && frame.source !== connection.source) {
    connection = { ...connection, source: frame.source };
    emitConnection();
  }
}

export function getLatestFrame(): Frame | null {
  return latestFrame;
}

export function ringBufferLength(): number {
  return length;
}

/** Oldest → newest. `fn` must not retain the frame past the call if you mutate later. */
export function forEachRingBuffer(fn: (f: Frame, i: number) => void) {
  const start = length < RING_CAPACITY ? 0 : head;
  for (let i = 0; i < length; i++) {
    const f = ringBuffer[(start + i) % RING_CAPACITY];
    if (f) fn(f, i);
  }
}

export function lastN(n: number, out: Frame[]): number {
  const take = Math.min(n, length);
  const start = (length < RING_CAPACITY ? 0 : head) + (length - take);
  for (let i = 0; i < take; i++) {
    out[i] = ringBuffer[(start + i) % RING_CAPACITY]!;
  }
  return take;
}

export function getConnection(): Connection {
  return connection;
}

export function setConnection(next: Partial<Connection>) {
  const merged: Connection = { ...connection, ...next };
  if (
    merged.state === connection.state &&
    merged.source === connection.source &&
    merged.subscribers === connection.subscribers
  ) {
    return;
  }
  connection = merged;
  emitConnection();
}

export function subscribe(fn: () => void) {
  connectionListeners.add(fn);
  return () => {
    connectionListeners.delete(fn);
  };
}

function emitConnection() {
  for (const listener of connectionListeners) listener();
}
