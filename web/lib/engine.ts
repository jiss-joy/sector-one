import { pushFrame, setConnection } from "@/lib/store";
import type { Frame } from "@/lib/types";

export const ENGINE = "http://127.0.0.1:8080";

export function connectEngine(): () => void {
  const eventSource = new EventSource(`${ENGINE}/api/telemetry`);

  eventSource.onopen = () => {
    setConnection({ state: "online" });
  };

  eventSource.onerror = () => {
    setConnection({ state: "offline" });
  };

  eventSource.onmessage = (ev) => {
    try {
      pushFrame(JSON.parse(ev.data) as Frame);
    } catch {
      // drop a bad line; do not tear down the socket
    }
  };

  const poll = window.setInterval(() => {
    void fetchHealth();
  }, 2000);
  void fetchHealth();

  return () => {
    eventSource.close();
    window.clearInterval(poll);
    setConnection({ state: "offline", subscribers: 0 });
  };
}

async function fetchHealth() {
  try {
    const res = await fetch(`${ENGINE}/api/health`);
    if (!res.ok) {
      setConnection({ state: "error" });
      return;
    }
    const j = (await res.json()) as {
      source?: string;
      subscribers?: number;
    };
    setConnection({
      state: "online",
      source: j.source ?? "",
      subscribers: typeof j.subscribers === "number" ? j.subscribers : 0,
    });
  } catch {
    setConnection({ state: "offline" });
  }
}
