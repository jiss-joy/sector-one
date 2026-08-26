"use client";

import { useEffect, useRef } from "react";

/** One rAF loop. `draw` is stored after commit so the loop is not restarted. */
export function useRaf(draw: () => void) {
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
