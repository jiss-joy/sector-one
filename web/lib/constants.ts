/**
 * Shared configuration and magic thresholds for the telemetry dashboard.
 */

export const TELEMETRY_CONFIG = {
  // Thresholds for visual alerts
  THRESHOLD_LOCK_MIN_SPEED: 20, // km/h
  THRESHOLD_LOCK_WHEEL_RAD_S: 1, // rad/s
  THRESHOLD_SLIP_RATIO: 0.15,

  // Gauge limits (calibrate these to your car)
  RPM_MAX: 8000,
  LOAD_MAX: 8000, // Newton

  // UI Scaling
  FRICTION_CIRCLE_G_SCALE: 2.5, // G units to edge of circle
  TRACE_HISTORY_COUNT: 600, // 10s at 60Hz (matches store ring buffer)
  FRICTION_CIRCLE_TAIL_COUNT: 180, // ~3s at 60Hz
} as const;
