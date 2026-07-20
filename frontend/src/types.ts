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

export interface Step {
  level: number;
  name: string;
  nameJa: string;
  pdu: string;
  addsHeader: boolean;
  headers: Record<string, string>;
  processing: string[];
  payload: string;
  headerBytes: number;
  totalBytes: number;
  structure: string;
  note: string;
  bitstream: string;
}

export interface EncapsulateRequest {
  message: string;
  srcIp: string;
  dstIp: string;
}
