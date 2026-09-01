/**
 * Shared configuration and magic thresholds for the telemetry dashboard.
 */

export const TELEMETRY_CONFIG = {
  // Thresholds for visual alerts
  THRESHOLD_LOCK_MIN_SPEED: 20, // km/h
  THRESHOLD_LOCK_WHEEL_RAD_S: 1, // rad/s
  THRESHOLD_SLIP_RATIO: 0.15,

  // UI Scaling
  FRICTION_CIRCLE_G_SCALE: 2.5, // G units to edge of circle
  TRACE_HISTORY_COUNT: 600, // 10s at 60Hz (matches store ring buffer)
  FRICTION_CIRCLE_TAIL_COUNT: 180, // ~3s at 60Hz

  // Left turn steer = Weight shifts to Right (FR/RR load up)
  // Braking = Dot moves Down (needs flip to move Up)
  FRICTION_CIRCLE_FLIP_LONG: true,
} as const;
