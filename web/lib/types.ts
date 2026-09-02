export type Frame = {
  ts: number;
  speed_kmh: number;
  rpm: number;
  gear: number;
  throttle: number;
  brake: number;
  clutch: number;
  steer: number;
  g_lat: number;
  g_long: number;
  g_vert: number;
  lap_time_ms: number;
  last_lap_ms: number;
  best_lap_ms: number;
  lap_count: number;
  abs_enabled: boolean;
  abs_in_action: boolean;
  tc_enabled: boolean;
  tc_in_action: boolean;
  in_pit: boolean;
  engine_limiter: boolean;
  wheel_rad_s: [number, number, number, number];
  slip_ratio: [number, number, number, number];
  slip_angle: [number, number, number, number];
  normalized_pos: number;
  load_n: [number, number, number, number];
  max_rpm: number;
  max_load: number;
  source: string;
};

export type ConnectionState = "offline" | "online" | "error";

export type Connection = {
  state: ConnectionState;
  source: string;
  subscribers: number;
};
