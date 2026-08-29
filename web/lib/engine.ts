import { pushFrame, setConn } from "@/lib/store";
import type { Frame } from "@/lib/types";

export const ENGINE = "";

export function connectEngine(): () => void {
  const eventSource = new EventSource(`${ENGINE}/api/telemetry`);

  eventSource.onopen = () => {
    setConn({ state: "open" });
  };

  eventSource.onerror = () => {
    setConn({ state: "error" });
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
    setConn({ state: "down", subscribers: 0 });
  };
}

async function fetchHealth() {
  try {
    const res = await fetch(`${ENGINE}/api/health`);
    if (!res.ok) {
      setConn({ state: "error" });
      return;
    }
    const j = (await res.json()) as {
      source?: string;
      subscribers?: number;
    };
    setConn({
      source: j.source ?? "",
      subscribers: typeof j.subscribers === "number" ? j.subscribers : 0,
    });
  } catch {
    setConn({ state: "error" });
  }
}
