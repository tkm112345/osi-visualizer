import type { Text } from "./types";

export type Lang = "ja" | "en";

// pick はバックエンドが返す Text から表示言語を選ぶ。
export function pick(t: Text | undefined, lang: Lang): string {
  if (!t) return "";
  return t[lang] || t.ja || t.en;
}

// UI 文言（フロントエンド固定文字列）の辞書。
export const ui = {
  title: { ja: "OSI 通信シミュレーター", en: "OSI Model Visualizer" },
  subtitle: {
    ja: "送信ホスト（カプセル化 L7→L1）と受信ホスト（デカプセル化 L1→L7）で、データがどう処理されるかを可視化します。各レイヤーをクリックすると詳細が見られます。",
    en: "Visualizes how data is processed on the sending host (encapsulation L7→L1) and the receiving host (decapsulation L1→L7). Click a layer to see details.",
  },
  simNote: {
    ja: "⚠️ これは擬似シミュレーションです。実際にネットワークへパケットを送信することはありません。宛先 IP は表示上のラベルで、通信は一切発生しません。",
    en: "⚠️ This is a pseudo-simulation. No packets are ever sent to the network. The destination IP is only a label; no communication happens.",
  },
  protocol: { ja: "プロトコル", en: "Protocol" },
  message: { ja: "メッセージ / ペイロード（テンプレート）", en: "Message / payload (template)" },
  srcIp: { ja: "送信元 IP（ホスト A）", en: "Source IP (host A)" },
  dstIp: { ja: "宛先 IP（ホスト B）", en: "Destination IP (host B)" },
  send: { ja: "擬似送信 ▶ シミュレート", en: "Simulate send ▶" },
  sending: { ja: "処理中...", en: "Working..." },
  stop: { ja: "■ 停止", en: "■ Stop" },
  langLabel: { ja: "言語", en: "Language" },
  hostASend: { ja: "送信ホスト A ▼ カプセル化", en: "Sending host A ▼ Encapsulation" },
  hostBRecv: { ja: "受信ホスト B ▲ デカプセル化", en: "Receiving host B ▲ Decapsulation" },
  wire: { ja: "物理媒体（擬似）", en: "Physical medium (simulated)" },
  loadingStacks: { ja: "レイヤースタックを読み込み中...", en: "Loading layer stack..." },
  hostATag: { ja: "送信ホスト A", en: "Sending host A" },
  hostBTag: { ja: "受信ホスト B", en: "Receiving host B" },
  errorPrefix: { ja: "エラー: ", en: "Error: " },
  errorHint: {
    ja: "バックエンド（http://localhost:8080）が起動しているか確認してください。",
    en: "Check that the backend (http://localhost:8080) is running.",
  },
  serialL2: { ja: " フレーム", en: " frame" },
  serialL1: { ja: "L1: 信号線", en: "L1: Signal lines" },
  serialNoIp: { ja: "L3〜L7 は使わない（IP なし）", en: "L3–L7 unused (no IP)" },
  l4None: { ja: "なし (ICMP)", en: "none (ICMP)" },
  tlsChip: { ja: "L6: TLS 暗号化", en: "L6: TLS encryption" },
  // LayerCard
  notUsed: { ja: "このシナリオでは使用しない", en: "Not used in this scenario" },
  noHeader: { ja: "ヘッダなし", en: "No header" },
  cumulative: { ja: "累積", en: "total" },
  seeData: { ja: "▶ この層での実データを見る", en: "▶ Show actual data at this layer" },
  hideData: { ja: "▼ 実データを隠す", en: "▼ Hide actual data" },
  // FrameView
  kindHeader: { ja: "ヘッダ", en: "Header" },
  kindPayload: { ja: "ペイロード", en: "Payload" },
  kindTrailer: { ja: "トレーラ", en: "Trailer" },
  // PacketDetail
  clickHint: {
    ja: "レイヤーをクリックすると、そのレイヤーで何が起きているかが表示されます。",
    en: "Click a layer to see what happens at that layer.",
  },
  layerUnusedNote: {
    ja: "この通信シナリオでは、この層は使用されません。",
    en: "This layer is not used in this communication scenario.",
  },
  headerAdd: { ja: "このレイヤーで付与するヘッダ", en: "Header added at this layer" },
  headerRemove: { ja: "このレイヤーで解析・除去するヘッダ", en: "Header parsed/removed at this layer" },
  dataBodyNoHeader: {
    ja: "（データ本体。まだヘッダは付いていません）",
    en: "(Data body. No header yet.)",
  },
  headerSize: { ja: "ヘッダサイズ: ", en: "Header size: " },
  processingNoHeader: {
    ja: "このレイヤーの処理（ヘッダは付与しない）",
    en: "Processing at this layer (no header added)",
  },
  bitstream: { ja: "ビット列（先頭）", en: "Bit stream (start)" },
  encStructure: { ja: "カプセル化構造", en: "Encapsulation structure" },
  cumulativeSize: { ja: "累積サイズ: ", en: "Cumulative size: " },
} as const;

export function t(key: keyof typeof ui, lang: Lang): string {
  return ui[key][lang];
}
