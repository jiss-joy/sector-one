/**
 * Phase 6 Step 0 — calibration notes
 *
 * Session: lemans.bin / rss_gtm_adonis_v8_evo @ lemans_2017
 *
 * 1. Wheel labels: left-hander load moves to FR/RR. FL FR / RL RR order is correct.
 *    Do not invert in ParseCarInfo.
 * 2. Friction circle: hard brake only moved toward the labelled Brake end with
 *    FRICTION_CIRCLE_FLIP_LONG = true.
 * 3. LOAD_BAR_N: peak single-corner load on that file sat near the engine default
 *    of 5000 N. +20% headroom → 6000. Tick on the bar is this value.
 */

export const THRESHOLDS = {
  LOCK_MIN_SPEED_KMH: 20,
  LOCK_WHEEL_RAD_S: 1,
  SLIP_RATIO: 0.15,

  FRICTION_CIRCLE_G_SCALE: 2.5,
  FRICTION_CIRCLE_TAIL_COUNT: 180,
  FRICTION_CIRCLE_FLIP_LONG: true,

  TRACE_HISTORY_COUNT: 900,
  LOAD_BAR_N: 6000,
  LOAD_SPARK_COUNT: 90,
} as const;
