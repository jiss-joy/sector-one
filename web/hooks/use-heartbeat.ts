"use client";

import { useEffect, useRef } from "react";

type Listener = () => void;

// Single global set of listeners to be called on every frame
const listeners = new Set<React.MutableRefObject<Listener>>();

let rafId: number | null = null;

/**
 * The single source of truth for the animation loop.
 * It iterates through all registered refs and calls their current functions.
 */
function tick() {
  listeners.forEach((ref) => {
    ref.current();
  });
  rafId = requestAnimationFrame(tick);
}

/**
 * useHeartbeat registers a callback to be executed once per frame
 * within a single shared requestAnimationFrame loop.
 */
export function useHeartbeat(draw: Listener) {
  const drawRef = useRef(draw);

  // Update the ref so the loop always has the latest logic
  useEffect(() => {
    drawRef.current = draw;
  });

  useEffect(() => {
    listeners.add(drawRef);

    // Start the loop if this is the first listener
    if (rafId === null) {
      rafId = requestAnimationFrame(tick);
    }

    return () => {
      listeners.add(drawRef);
      listeners.delete(drawRef);

      // Stop the loop if no one is listening anymore
      if (listeners.size === 0 && rafId !== null) {
        cancelAnimationFrame(rafId);
        rafId = null;
      }
    };
  }, []);
}
