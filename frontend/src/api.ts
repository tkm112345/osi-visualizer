import type { EncapsulateRequest, Layer, Step } from "./types";

const BASE = "http://localhost:8080";

export async function fetchLayers(): Promise<Layer[]> {
  const res = await fetch(`${BASE}/api/layers`);
  if (!res.ok) throw new Error(`GET /api/layers failed: ${res.status}`);
  return res.json();
}

async function postSteps(path: string, req: EncapsulateRequest): Promise<Step[]> {
  const res = await fetch(`${BASE}${path}`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(req),
  });
  if (!res.ok) throw new Error(`POST ${path} failed: ${res.status}`);
  const data: { steps: Step[] } = await res.json();
  return data.steps;
}

// 送信ホスト: L7 → L1 のカプセル化
export function encapsulate(req: EncapsulateRequest): Promise<Step[]> {
  return postSteps("/api/encapsulate", req);
}

// 受信ホスト: L1 → L7 のデカプセル化（擬似）
export function decapsulate(req: EncapsulateRequest): Promise<Step[]> {
  return postSteps("/api/decapsulate", req);
}
