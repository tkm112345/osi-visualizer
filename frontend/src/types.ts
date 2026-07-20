// バックエンドの JSON と対応する型定義。

export interface Layer {
  level: number;
  name: string;
  nameJa: string;
  pdu: string;
  protocols: string[];
  description: string;
  addsHeader: boolean;
}

export interface Protocol {
  key: string;
  l7Name: string;
  label: string;
  category: string;
  family: string;
  transport: string;
  port: number;
  tls: boolean;
  l3Protocol: string;
  samplePayload: string;
  description: string;
}

// FramePart は、ある層での PDU を構成する 1 区画（ヘッダ / ペイロード / トレーラ）。
export interface FramePart {
  label: string;
  detail: string;
  kind: "header" | "payload" | "trailer";
  bytes: number;
}

export interface Step {
  level: number;
  name: string;
  nameJa: string;
  pdu: string;
  addsHeader: boolean;
  active: boolean;
  headers: Record<string, string>;
  processing: string[];
  payload: string;
  headerBytes: number;
  totalBytes: number;
  structure: string;
  note: string;
  bitstream: string;
  frame: FramePart[];
}

export interface EncapsulateRequest {
  message: string;
  srcIp: string;
  dstIp: string;
  protocol: string;
}
