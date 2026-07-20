// バックエンドの JSON と対応する型定義。

// Text はバックエンドが返す多言語文字列（日本語・英語の両方）。
export interface Text {
  ja: string;
  en: string;
}

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
  label: Text;
  category: Text;
  family: string;
  transport: string;
  port: number;
  tls: boolean;
  l3Protocol: string;
  samplePayload: string;
  description: Text;
}

// FramePart は、ある層での PDU を構成する 1 区画（ヘッダ / ペイロード / トレーラ）。
export interface FramePart {
  label: Text;
  detail: Text;
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
  processing: Text[];
  payload: string;
  headerBytes: number;
  totalBytes: number;
  structure: string;
  note: Text;
  bitstream: string;
  frame: FramePart[];
}

export interface EncapsulateRequest {
  message: string;
  srcIp: string;
  dstIp: string;
  protocol: string;
}
