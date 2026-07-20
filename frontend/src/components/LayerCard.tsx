import { useState } from "react";
import type { Step } from "../types";
import FrameView from "./FrameView";

interface Props {
  step: Step;
  selected: boolean;
  active: boolean; // アニメーション中に現在処理しているレイヤー
  reached: boolean; // アニメーションで既に到達済みか
  onClick: () => void;
}

// L2〜L7 の色（下ほど暖色）。
const LEVEL_COLORS: Record<number, string> = {
  7: "#7c3aed",
  6: "#6366f1",
  5: "#0ea5e9",
  4: "#10b981",
  3: "#f59e0b",
  2: "#ef4444",
  1: "#6b7280",
};

export default function LayerCard({ step, selected, active, reached, onClick }: Props) {
  const [open, setOpen] = useState(false);
  const color = LEVEL_COLORS[step.level] ?? "#888";
  const inactive = step.active === false;
  const classes = ["layer-card"];
  if (selected) classes.push("selected");
  if (active) classes.push("active");
  if (!reached) classes.push("dim");
  if (!step.addsHeader) classes.push("no-header");
  if (inactive) classes.push("inactive");

  const hasFrame = !inactive && step.frame && step.frame.length > 0;

  return (
    <div
      className={classes.join(" ")}
      style={{ borderLeftColor: color }}
      onClick={onClick}
      role="button"
      tabIndex={0}
      onKeyDown={(e) => {
        if (e.key === "Enter" || e.key === " ") {
          e.preventDefault();
          onClick();
        }
      }}
    >
      <div className="layer-badge" style={{ background: color }}>
        L{step.level}
      </div>
      <div className="layer-body">
        <div className="layer-title">
          {step.name} <span className="layer-ja">{step.nameJa}</span>
        </div>
        <div className="layer-meta">
          {inactive ? (
            <span className="bytes muted">このシナリオでは使用しない</span>
          ) : (
            <>
              <span className="pdu">PDU: {step.pdu}</span>
              {step.headers.protocol && (
                <span className="proto-tag">{step.headers.protocol}</span>
              )}
              {step.addsHeader ? (
                <span className="bytes">累積 {step.totalBytes} B</span>
              ) : (
                <span className="bytes muted">ヘッダなし</span>
              )}
            </>
          )}
        </div>
        {hasFrame && (
          <div className="frame-accordion">
            <button
              className="frame-toggle"
              onClick={(e) => {
                e.stopPropagation();
                setOpen((v) => !v);
              }}
            >
              {open ? "▼ 実データを隠す" : "▶ この層での実データを見る"}
            </button>
            {open && <FrameView parts={step.frame} />}
          </div>
        )}
      </div>
    </div>
  );
}
