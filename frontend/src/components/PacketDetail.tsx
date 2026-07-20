import type { Step } from "../types";
import { pick, t, type Lang } from "../i18n";

interface Props {
  step: Step | null;
  mode: "encap" | "decap";
  lang: Lang;
}

export default function PacketDetail({ step, mode, lang }: Props) {
  if (!step) {
    return (
      <div className="detail empty">
        <p>{t("clickHint", lang)}</p>
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
        <p className="detail-note">{pick(step.note, lang)}</p>
        <p className="muted">{t("layerUnusedNote", lang)}</p>
      </div>
    );
  }

  const headerLabel = mode === "encap" ? t("headerAdd", lang) : t("headerRemove", lang);

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

      <p className="detail-note">{pick(step.note, lang)}</p>

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
            <p className="muted">{t("dataBodyNoHeader", lang)}</p>
          )}
          {step.headerBytes > 0 && (
            <p className="muted">
              {t("headerSize", lang)}
              {step.headerBytes} B
            </p>
          )}
        </section>
      ) : (
        <section>
          <h3>{t("processingNoHeader", lang)}</h3>
          <ul className="processing-list">
            {step.processing.map((p, i) => (
              <li key={i}>{pick(p, lang)}</li>
            ))}
          </ul>
        </section>
      )}

      {step.bitstream && (
        <section>
          <h3>{t("bitstream", lang)}</h3>
          <code className="bitstream">{step.bitstream}</code>
        </section>
      )}

      <section>
        <h3>{t("encStructure", lang)}</h3>
        <code className="structure">{step.structure}</code>
        <p className="muted">
          {t("cumulativeSize", lang)}
          {step.totalBytes} B
        </p>
      </section>
    </div>
  );
}
