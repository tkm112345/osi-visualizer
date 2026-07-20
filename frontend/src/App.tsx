import { useEffect, useRef, useState } from "react";
import { encapsulate } from "./api";
import type { Step } from "./types";
import LayerStack from "./components/LayerStack";
import PacketDetail from "./components/PacketDetail";

export default function App() {
  const [message, setMessage] = useState("Hello");
  const [srcIp, setSrcIp] = useState("192.168.0.10");
  const [dstIp, setDstIp] = useState("93.184.216.34");

  const [steps, setSteps] = useState<Step[]>([]);
  const [selectedLevel, setSelectedLevel] = useState<number | null>(null);
  const [activeLevel, setActiveLevel] = useState<number | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  const timers = useRef<number[]>([]);

  function clearTimers() {
    timers.current.forEach((t) => window.clearTimeout(t));
    timers.current = [];
  }

  useEffect(() => clearTimers, []);

  async function handleSend() {
    setError(null);
    setLoading(true);
    setSelectedLevel(null);
    clearTimers();
    try {
      const result = await encapsulate({ message, srcIp, dstIp });
      setSteps(result);
      // L7 → L1 へ順にハイライトするアニメーション。
      setActiveLevel(7);
      result.forEach((step, i) => {
        const t = window.setTimeout(() => {
          setActiveLevel(step.level);
          if (i === result.length - 1) {
            const done = window.setTimeout(() => setActiveLevel(null), 600);
            timers.current.push(done);
          }
        }, i * 600);
        timers.current.push(t);
      });
    } catch (e) {
      setError(e instanceof Error ? e.message : String(e));
      setSteps([]);
    } finally {
      setLoading(false);
    }
  }

  const selectedStep = steps.find((s) => s.level === selectedLevel) ?? null;

  return (
    <div className="app">
      <header className="app-header">
        <h1>OSI Model Visualizer</h1>
        <p className="subtitle">
          データが L7 → L1 へカプセル化されていく様子を可視化します。各レイヤーをクリックすると詳細が見られます。
        </p>
      </header>

      <section className="controls">
        <label>
          メッセージ
          <input value={message} onChange={(e) => setMessage(e.target.value)} />
        </label>
        <label>
          送信元 IP
          <input value={srcIp} onChange={(e) => setSrcIp(e.target.value)} />
        </label>
        <label>
          宛先 IP
          <input value={dstIp} onChange={(e) => setDstIp(e.target.value)} />
        </label>
        <button className="send-btn" onClick={handleSend} disabled={loading}>
          {loading ? "送信中..." : "送信 ↓ カプセル化"}
        </button>
      </section>

      {error && (
        <div className="error">
          エラー: {error}
          <br />
          バックエンド（http://localhost:8080）が起動しているか確認してください。
        </div>
      )}

      <main className="main">
        <div className="stack-col">
          {steps.length === 0 ? (
            <p className="placeholder">「送信」を押すとレイヤースタックが表示されます。</p>
          ) : (
            <>
              <div className="direction-label">▲ 上位層（アプリに近い）</div>
              <LayerStack
                steps={steps}
                selectedLevel={selectedLevel}
                activeLevel={activeLevel}
                onSelect={setSelectedLevel}
              />
              <div className="direction-label">▼ 下位層（物理媒体へ）</div>
            </>
          )}
        </div>
        <aside className="detail-col">
          <PacketDetail step={selectedStep} />
        </aside>
      </main>
    </div>
  );
}
