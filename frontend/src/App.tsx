import { useEffect, useRef, useState } from "react";
import { decapsulate, encapsulate, fetchProtocols } from "./api";
import type { Protocol, Step } from "./types";
import LayerStack from "./components/LayerStack";
import PacketDetail from "./components/PacketDetail";

type Host = "A" | "B";
interface Selection {
  host: Host;
  level: number;
}

const STEP_MS = 600;

export default function App() {
  const [message, setMessage] = useState("Hello");
  const [srcIp, setSrcIp] = useState("192.168.0.10");
  const [dstIp, setDstIp] = useState("93.184.216.34");

  const [protocols, setProtocols] = useState<Protocol[]>([]);
  const [protocol, setProtocol] = useState("http");

  const [encapSteps, setEncapSteps] = useState<Step[]>([]);
  const [decapSteps, setDecapSteps] = useState<Step[]>([]);
  const [activeA, setActiveA] = useState<number | null>(null);
  const [activeB, setActiveB] = useState<number | null>(null);
  const [selection, setSelection] = useState<Selection | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  const timers = useRef<number[]>([]);
  function clearTimers() {
    timers.current.forEach((t) => window.clearTimeout(t));
    timers.current = [];
  }
  useEffect(() => clearTimers, []);

  useEffect(() => {
    fetchProtocols()
      .then(setProtocols)
      .catch(() => setProtocols([]));
  }, []);

  async function handleSend() {
    setError(null);
    setLoading(true);
    setSelection(null);
    clearTimers();
    setActiveA(null);
    setActiveB(null);
    try {
      const req = { message, srcIp, dstIp, protocol };
      const [enc, dec] = await Promise.all([encapsulate(req), decapsulate(req)]);
      setEncapSteps(enc);
      setDecapSteps(dec);

      // フェーズ A: 送信ホストで L7 → L1（enc は L7→L1 順）
      enc.forEach((step, i) => {
        const t = window.setTimeout(() => {
          setActiveA(step.level);
          setActiveB(null);
        }, i * STEP_MS);
        timers.current.push(t);
      });
      // フェーズ B: 受信ホストで L1 → L7（dec は L1→L7 順）
      const offset = enc.length * STEP_MS;
      dec.forEach((step, j) => {
        const t = window.setTimeout(() => {
          setActiveA(null);
          setActiveB(step.level);
        }, offset + j * STEP_MS);
        timers.current.push(t);
      });
      const done = window.setTimeout(() => {
        setActiveA(null);
        setActiveB(null);
      }, offset + dec.length * STEP_MS);
      timers.current.push(done);
    } catch (e) {
      setError(e instanceof Error ? e.message : String(e));
      setEncapSteps([]);
      setDecapSteps([]);
    } finally {
      setLoading(false);
    }
  }

  const selectedProtocol = protocols.find((p) => p.key === protocol) ?? null;
  const started = encapSteps.length > 0;
  const activeSteps = selection?.host === "B" ? decapSteps : encapSteps;
  const selectedStep = selection
    ? activeSteps.find((s) => s.level === selection.level) ?? null
    : null;

  return (
    <div className="app">
      <header className="app-header">
        <h1>OSI 通信シミュレーター</h1>
        <p className="subtitle">
          送信ホスト（カプセル化 L7→L1）と受信ホスト（デカプセル化 L1→L7）で、
          データがどう処理されるかを可視化します。各レイヤーをクリックすると詳細が見られます。
        </p>
        <p className="sim-note">
          ⚠️ これは<strong>擬似シミュレーション</strong>です。実際にネットワークへパケットを送信することはありません。
          宛先 IP は表示上のラベルで、通信は一切発生しません。
        </p>
      </header>

      <section className="controls">
        <label>
          プロトコル
          <select value={protocol} onChange={(e) => setProtocol(e.target.value)}>
            {protocols.map((p) => (
              <option key={p.key} value={p.key}>
                {p.label}
              </option>
            ))}
          </select>
        </label>
        <label>
          メッセージ
          <input value={message} onChange={(e) => setMessage(e.target.value)} />
        </label>
        <label>
          送信元 IP（ホスト A）
          <input value={srcIp} onChange={(e) => setSrcIp(e.target.value)} />
        </label>
        <label>
          宛先 IP（ホスト B）
          <input value={dstIp} onChange={(e) => setDstIp(e.target.value)} />
        </label>
        <button className="send-btn" onClick={handleSend} disabled={loading}>
          {loading ? "処理中..." : "擬似送信 ▶ シミュレート"}
        </button>
      </section>

      {selectedProtocol && (
        <div className="proto-info">
          <span className="proto-chip">L7: {selectedProtocol.l7Name}</span>
          <span className="proto-arrow">→</span>
          <span className="proto-chip">
            L4: {selectedProtocol.transport === "ICMP" ? "なし (ICMP)" : selectedProtocol.transport}
            {selectedProtocol.port > 0 && ` :${selectedProtocol.port}`}
          </span>
          <span className="proto-arrow">→</span>
          <span className="proto-chip">L3: {selectedProtocol.l3Protocol}</span>
          {selectedProtocol.tls && <span className="proto-chip tls">L6: TLS 暗号化</span>}
          <span className="proto-desc">{selectedProtocol.description}</span>
        </div>
      )}

      {error && (
        <div className="error">
          エラー: {error}
          <br />
          バックエンド（http://localhost:8080）が起動しているか確認してください。
        </div>
      )}

      <main className="main">
        <div className="stacks">
          {!started ? (
            <p className="placeholder">「擬似送信」を押すと両ホストのレイヤースタックが表示されます。</p>
          ) : (
            <div className="two-hosts">
              <div className="host-col">
                <div className="host-title send">送信ホスト A ▼ カプセル化</div>
                <LayerStack
                  steps={encapSteps}
                  mode="encap"
                  selectedLevel={selection?.host === "A" ? selection.level : null}
                  activeLevel={activeA}
                  onSelect={(level) => setSelection({ host: "A", level })}
                />
              </div>

              <div className="wire" aria-hidden>
                <div className="wire-line" />
                <div className="wire-label">物理媒体（擬似）</div>
              </div>

              <div className="host-col">
                <div className="host-title recv">受信ホスト B ▲ デカプセル化</div>
                <LayerStack
                  steps={decapSteps}
                  mode="decap"
                  selectedLevel={selection?.host === "B" ? selection.level : null}
                  activeLevel={activeB}
                  onSelect={(level) => setSelection({ host: "B", level })}
                />
              </div>
            </div>
          )}
        </div>

        <aside className="detail-col">
          {selection && (
            <div className="detail-host-tag">
              {selection.host === "A" ? "送信ホスト A" : "受信ホスト B"}
            </div>
          )}
          <PacketDetail step={selectedStep} mode={selection?.host === "B" ? "decap" : "encap"} />
        </aside>
      </main>
    </div>
  );
}
