import type { EncapsulateRequest, Layer, Step } from "./types";

const BASE = "http://localhost:8080";

export async function fetchLayers(): Promise<Layer[]> {
  const res = await fetch(`${BASE}/api/layers`);
  if (!res.ok) throw new Error(`GET /api/layers failed: ${res.status}`);
  return res.json();
}

export async function encapsulate(req: EncapsulateRequest): Promise<Step[]> {
  const res = await fetch(`${BASE}/api/encapsulate`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(req),
  });
  if (!res.ok) throw new Error(`POST /api/encapsulate failed: ${res.status}`);
  const data: { steps: Step[] } = await res.json();
  return data.steps;
}
