import type { Step } from "../types";

interface Props {
  step: Step | null;
  mode: "encap" | "decap";
}

export default function PacketDetail({ step, mode }: Props) {
  if (!step) {
    return (
      <div className="detail empty">
        <p>レイヤーをクリックすると、そのレイヤーで何が起きているかが表示されます。</p>
      </div>
    );
  }

  if (step.active === false) {
    return (
      <div className="detail">
        <h2>
          L{step.level} {step.name}
          <span className="detail-ja">{step.nameJa}</span>
        </h2>
        <p className="detail-note">{step.note}</p>
        <p className="muted">この通信シナリオでは、この層は使用されません。</p>
      </div>
    );
  }

  const headerLabel =
    mode === "encap" ? "このレイヤーで付与するヘッダ" : "このレイヤーで解析・除去するヘッダ";

  const headerEntries = Object.entries(step.headers);

  return (
    <div className="detail">
      <h2>
        L{step.level} {step.name}
        <span className="detail-ja">{step.nameJa}</span>
      </h2>

      <div className="detail-row">
        <span className="label">PDU</span>
        <span>{step.pdu}</span>
      </div>

      <p className="detail-note">{step.note}</p>

      {step.addsHeader ? (
        <section>
          <h3>{headerLabel}</h3>
          {headerEntries.length > 0 ? (
            <table className="header-table">
              <tbody>
                {headerEntries.map(([k, v]) => (
                  <tr key={k}>
                    <td className="field">{k}</td>
                    <td className="value">{v}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          ) : (
            <p className="muted">（データ本体。まだヘッダは付いていません）</p>
          )}
          {step.headerBytes > 0 && (
            <p className="muted">ヘッダサイズ: {step.headerBytes} B</p>
          )}
        </section>
      ) : (
        <section>
          <h3>このレイヤーの処理（ヘッダは付与しない）</h3>
          <ul className="processing-list">
            {step.processing.map((p) => (
              <li key={p}>{p}</li>
            ))}
          </ul>
        </section>
      )}

      {step.bitstream && (
        <section>
          <h3>ビット列（先頭）</h3>
          <code className="bitstream">{step.bitstream}</code>
        </section>
      )}

      <section>
        <h3>カプセル化構造</h3>
        <code className="structure">{step.structure}</code>
        <p className="muted">累積サイズ: {step.totalBytes} B</p>
      </section>
    </div>
  );
}
