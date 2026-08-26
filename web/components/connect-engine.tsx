"use client";

import { useEffect } from "react";

import { connectEngine } from "@/lib/engine";

export function ConnectEngine() {
  useEffect(() => connectEngine(), []);
  return null;
}
