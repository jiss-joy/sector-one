"use client";

import { useEffect, useRef } from "react";

// It creates a recursive loop that executes your callback exactly once per monitor refresh (usually 60Hz or 144Hz).
export function useRequestAnimationFrame(draw: () => void) {
  const ref = useRef(draw);

  useEffect(() => {
    ref.current = draw;
  });

  useEffect(() => {
    let id = 0;
    const tick = () => {
      ref.current();
      id = requestAnimationFrame(tick);
    };
    id = requestAnimationFrame(tick);
    return () => cancelAnimationFrame(id);
  }, []);
}
