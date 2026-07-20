import type { FramePart } from "../types";

interface Props {
  parts: FramePart[];
}

const KIND_LABEL: Record<FramePart["kind"], string> = {
  header: "ヘッダ",
  payload: "ペイロード",
  trailer: "トレーラ",
};

// FrameView は、その層での PDU（フレーム）の中身を区画ごとに縦並びで表示する。
// ヘッダ/ペイロード/トレーラの各区画に実データ（テンプレート由来）を示す。
export default function FrameView({ parts }: Props) {
  return (
    <div className="frame-view">
      {parts.map((part, i) => (
        <div key={i} className={`frame-part ${part.kind}`}>
          <div className="frame-part-head">
            <span className="frame-part-label">{part.label}</span>
            <span className="frame-part-tag">
              {KIND_LABEL[part.kind]} · {part.bytes} B
            </span>
          </div>
          <pre className="frame-part-detail">{part.detail}</pre>
        </div>
      ))}
    </div>
  );
}
