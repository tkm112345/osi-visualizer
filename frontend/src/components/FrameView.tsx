import type { FramePart } from "../types";
import { pick, t, type Lang } from "../i18n";

interface Props {
  parts: FramePart[];
  lang: Lang;
}

const KIND_KEY = {
  header: "kindHeader",
  payload: "kindPayload",
  trailer: "kindTrailer",
} as const;

// FrameView は、その層での PDU（フレーム）の中身を区画ごとに縦並びで表示する。
// ヘッダ/ペイロード/トレーラの各区画に実データ（テンプレート由来）を示す。
export default function FrameView({ parts, lang }: Props) {
  return (
    <div className="frame-view">
      {parts.map((part, i) => (
        <div key={i} className={`frame-part ${part.kind}`}>
          <div className="frame-part-head">
            <span className="frame-part-label">{pick(part.label, lang)}</span>
            <span className="frame-part-tag">
              {t(KIND_KEY[part.kind], lang)} · {part.bytes} B
            </span>
          </div>
          <pre className="frame-part-detail">{pick(part.detail, lang)}</pre>
        </div>
      ))}
    </div>
  );
}
