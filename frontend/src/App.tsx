import { useEffect, useRef, useState } from "react";
import { decapsulate, encapsulate, fetchProtocols } from "./api";
import type { Protocol, Step, Text } from "./types";
import { pick, t, type Lang } from "./i18n";
import LayerStack from "./components/LayerStack";
import PacketDetail from "./components/PacketDetail";

type Host = "A" | "B";
interface Selection {
  host: Host;
  level: number;
}

const STEP_MS = 600;

export default function App() {
  const [lang, setLang] = useState<Lang>("ja");
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
  const [playing, setPlaying] = useState(false);

  const timers = useRef<number[]>([]);
  function clearTimers() {
    timers.current.forEach((t) => window.clearTimeout(t));
    timers.current = [];
  }
  useEffect(() => clearTimers, []);

  // アニメーションを途中で停止する（現在の状態のまま固定される）。
  function handleStop() {
    clearTimers();
    setPlaying(false);
  }

  // アニメーションせずに両ホストのスタックだけ取得して表示する。
  async function loadStacks(msg: string, proto: string) {
    try {
      const req = { message: msg, srcIp, dstIp, protocol: proto };
      const [enc, dec] = await Promise.all([encapsulate(req), decapsulate(req)]);
      setEncapSteps(enc);
      setDecapSteps(dec);
      setError(null);
    } catch (e) {
      setError(e instanceof Error ? e.message : String(e));
    }
  }

  // 初回にプロトコル一覧を取得し、既定プロトコルのスタックを表示しておく。
  useEffect(() => {
    fetchProtocols()
      .then((list) => {
        setProtocols(list);
        const p = list.find((x) => x.key === "http") ?? list[0];
        if (p) {
          setProtocol(p.key);
          setMessage(p.samplePayload);
          loadStacks(p.samplePayload, p.key);
        }
      })
      .catch((e) => setError(e instanceof Error ? e.message : String(e)));
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  // プロトコル変更時: サンプルペイロードに差し替え、スタックを再表示する。
  function handleProtocolChange(key: string) {
    clearTimers();
    setPlaying(false);
    setActiveA(null);
    setActiveB(null);
    setSelection(null);
    setProtocol(key);
    const p = protocols.find((x) => x.key === key);
    const sample = p ? p.samplePayload : message;
    setMessage(sample);
    loadStacks(sample, key);
  }

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
      setPlaying(true);

      // 使用しない層（UART の L3〜L7 など）はアニメーションを飛ばす。
      const encActive = enc.filter((s) => s.active);
      const decActive = dec.filter((s) => s.active);

      // フェーズ A: 送信ホストで L7 → L1（enc は L7→L1 順）
      encActive.forEach((step, i) => {
        const t = window.setTimeout(() => {
          setActiveA(step.level);
          setActiveB(null);
        }, i * STEP_MS);
        timers.current.push(t);
      });
      // フェーズ B: 受信ホストで L1 → L7（dec は L1→L7 順）
      const offset = encActive.length * STEP_MS;
      decActive.forEach((step, j) => {
        const t = window.setTimeout(() => {
          setActiveA(null);
          setActiveB(step.level);
        }, offset + j * STEP_MS);
        timers.current.push(t);
      });
      const done = window.setTimeout(() => {
        setActiveA(null);
        setActiveB(null);
        setPlaying(false);
      }, offset + decActive.length * STEP_MS);
      timers.current.push(done);
    } catch (e) {
      setError(e instanceof Error ? e.message : String(e));
    } finally {
      setLoading(false);
    }
  }

  const selectedProtocol = protocols.find((p) => p.key === protocol) ?? null;
  const isSerial = selectedProtocol?.family === "serial";
  const started = encapSteps.length > 0;
  const activeSteps = selection?.host === "B" ? decapSteps : encapSteps;
  const selectedStep = selection
    ? activeSteps.find((s) => s.level === selection.level) ?? null
    : null;

  // ドロップダウン用にカテゴリを出現順で並べる（英語名をキーにして安定化）。
  const categories: Text[] = [];
  protocols.forEach((p) => {
    if (!categories.some((c) => c.en === p.category.en)) categories.push(p.category);
  });

  return (
    <div className="app">
      <header className="app-header">
        <div className="header-top">
          <h1>{t("title", lang)}</h1>
          <div className="lang-switch" role="group" aria-label={t("langLabel", lang)}>
            <button
              className={lang === "ja" ? "on" : ""}
              onClick={() => setLang("ja")}
            >
              日本語
            </button>
            <button
              className={lang === "en" ? "on" : ""}
              onClick={() => setLang("en")}
            >
              English
            </button>
          </div>
        </div>
        <p className="subtitle">{t("subtitle", lang)}</p>
        <p className="sim-note">{t("simNote", lang)}</p>
      </header>

      <section className="controls">
        <label>
          {t("protocol", lang)}
          <select value={protocol} onChange={(e) => handleProtocolChange(e.target.value)}>
            {categories.map((cat) => (
              <optgroup key={cat.en} label={pick(cat, lang)}>
                {protocols
                  .filter((p) => p.category.en === cat.en)
                  .map((p) => (
                    <option key={p.key} value={p.key}>
                      {pick(p.label, lang)}
                    </option>
                  ))}
              </optgroup>
            ))}
          </select>
        </label>
        <label className="msg-field">
          {t("message", lang)}
          <textarea value={message} onChange={(e) => setMessage(e.target.value)} rows={4} />
        </label>
        <label>
          {t("srcIp", lang)}
          <input value={srcIp} onChange={(e) => setSrcIp(e.target.value)} disabled={isSerial} />
        </label>
        <label>
          {t("dstIp", lang)}
          <input value={dstIp} onChange={(e) => setDstIp(e.target.value)} disabled={isSerial} />
        </label>
        <button className="send-btn" onClick={handleSend} disabled={loading}>
          {loading ? t("sending", lang) : t("send", lang)}
        </button>
        {playing && (
          <button className="stop-btn" onClick={handleStop}>
            {t("stop", lang)}
          </button>
        )}
      </section>

      {selectedProtocol && (
        <div className="proto-info">
          {isSerial ? (
            <>
              <span className="proto-chip">
                L2: {selectedProtocol.l7Name}
                {t("serialL2", lang)}
              </span>
              <span className="proto-arrow">→</span>
              <span className="proto-chip">{t("serialL1", lang)}</span>
              <span className="proto-chip warn">{t("serialNoIp", lang)}</span>
            </>
          ) : (
            <>
              <span className="proto-chip">L7: {selectedProtocol.l7Name}</span>
              <span className="proto-arrow">→</span>
              <span className="proto-chip">
                L4: {selectedProtocol.transport === "ICMP" ? t("l4None", lang) : selectedProtocol.transport}
                {selectedProtocol.port > 0 && ` :${selectedProtocol.port}`}
              </span>
              <span className="proto-arrow">→</span>
              <span className="proto-chip">L3: {selectedProtocol.l3Protocol}</span>
              {selectedProtocol.tls && <span className="proto-chip tls">{t("tlsChip", lang)}</span>}
            </>
          )}
          <span className="proto-desc">{pick(selectedProtocol.description, lang)}</span>
        </div>
      )}

      {error && (
        <div className="error">
          {t("errorPrefix", lang)}
          {error}
          <br />
          {t("errorHint", lang)}
        </div>
      )}

      <main className="main">
        <div className="stacks">
          {!started ? (
            <p className="placeholder">{t("loadingStacks", lang)}</p>
          ) : (
            <div className="two-hosts">
              <div className="host-col">
                <div className="host-title send">{t("hostASend", lang)}</div>
                <LayerStack
                  steps={encapSteps}
                  mode="encap"
                  selectedLevel={selection?.host === "A" ? selection.level : null}
                  activeLevel={activeA}
                  lang={lang}
                  onSelect={(level) => setSelection({ host: "A", level })}
                />
              </div>

              <div className="wire" aria-hidden>
                <div className="wire-line" />
                <div className="wire-label">{t("wire", lang)}</div>
              </div>

              <div className="host-col">
                <div className="host-title recv">{t("hostBRecv", lang)}</div>
                <LayerStack
                  steps={decapSteps}
                  mode="decap"
                  selectedLevel={selection?.host === "B" ? selection.level : null}
                  activeLevel={activeB}
                  lang={lang}
                  onSelect={(level) => setSelection({ host: "B", level })}
                />
              </div>
            </div>
          )}
        </div>

        <aside className="detail-col">
          {selection && (
            <div className="detail-host-tag">
              {selection.host === "A" ? t("hostATag", lang) : t("hostBTag", lang)}
            </div>
          )}
          <PacketDetail
            step={selectedStep}
            mode={selection?.host === "B" ? "decap" : "encap"}
            lang={lang}
          />
        </aside>
      </main>
    </div>
  );
}
